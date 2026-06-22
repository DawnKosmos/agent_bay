package cdc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pglogrepl"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgproto3"
)

type Config struct {
	ConnString  string
	SlotName    string
	Publication string
	Timeout     time.Duration
}

type Replicator struct {
	config    Config
	relations map[uint32]*pglogrepl.RelationMessage
	handlers  map[string]*handlers
	tables    []TableConfig
}

type handlers struct {
	onInsert func(rel *pglogrepl.RelationMessage, tuple *pglogrepl.TupleData, lsn pglogrepl.LSN)
	onUpdate func(rel *pglogrepl.RelationMessage, oldTuple, newTuple *pglogrepl.TupleData, lsn pglogrepl.LSN)
	onDelete func(rel *pglogrepl.RelationMessage, tuple *pglogrepl.TupleData, lsn pglogrepl.LSN)
}

type Handlers[T any] struct {
	OnInsert func(e InsertEvent[T])
	OnUpdate func(e UpdateEvent[T])
	OnDelete func(e DeleteEvent[T])
}

func New(cfg Config) *Replicator {
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}
	return &Replicator{
		config:    cfg,
		relations: make(map[uint32]*pglogrepl.RelationMessage),
		handlers:  make(map[string]*handlers),
	}
}

func (r *Replicator) getOrCreateHandlers(table string) *handlers {
	if h, ok := r.handlers[table]; ok {
		return h
	}
	h := &handlers{}
	r.handlers[table] = h
	return h
}

func Register[T any](r *Replicator, cfg TableConfig, h Handlers[T]) {
	r.tables = append(r.tables, cfg)
	ih := r.getOrCreateHandlers(cfg.Name)
	if h.OnInsert != nil {
		ih.onInsert = func(rel *pglogrepl.RelationMessage, tuple *pglogrepl.TupleData, lsn pglogrepl.LSN) {
			event := InsertEvent[T]{
				New:       MapTuple[T](rel, tuple),
				TableName: cfg.Name,
				LSN:       lsn,
			}
			h.OnInsert(event)
		}
	}
	if h.OnUpdate != nil {
		ih.onUpdate = func(rel *pglogrepl.RelationMessage, oldTuple, newTuple *pglogrepl.TupleData, lsn pglogrepl.LSN) {
			event := UpdateEvent[T]{
				Old:       MapTuple[T](rel, oldTuple),
				New:       MapTuple[T](rel, newTuple),
				TableName: cfg.Name,
				LSN:       lsn,
			}
			h.OnUpdate(event)
		}
	}
	if h.OnDelete != nil {
		ih.onDelete = func(rel *pglogrepl.RelationMessage, tuple *pglogrepl.TupleData, lsn pglogrepl.LSN) {
			event := DeleteEvent[T]{
				Old:       MapTuple[T](rel, tuple),
				TableName: cfg.Name,
				LSN:       lsn,
			}
			h.OnDelete(event)
		}
	}
}

func OnInsert[T any](r *Replicator, table string, handler func(e InsertEvent[T])) {
	h := r.getOrCreateHandlers(table)
	h.onInsert = func(rel *pglogrepl.RelationMessage, tuple *pglogrepl.TupleData, lsn pglogrepl.LSN) {
		event := InsertEvent[T]{
			New:       MapTuple[T](rel, tuple),
			TableName: table,
			LSN:       lsn,
		}
		handler(event)
	}
}

func OnUpdate[T any](r *Replicator, table string, handler func(e UpdateEvent[T])) {
	h := r.getOrCreateHandlers(table)
	h.onUpdate = func(rel *pglogrepl.RelationMessage, oldTuple, newTuple *pglogrepl.TupleData, lsn pglogrepl.LSN) {
		event := UpdateEvent[T]{
			Old:       MapTuple[T](rel, oldTuple),
			New:       MapTuple[T](rel, newTuple),
			TableName: table,
			LSN:       lsn,
		}
		handler(event)
	}
}

func OnDelete[T any](r *Replicator, table string, handler func(e DeleteEvent[T])) {
	h := r.getOrCreateHandlers(table)
	h.onDelete = func(rel *pglogrepl.RelationMessage, tuple *pglogrepl.TupleData, lsn pglogrepl.LSN) {
		event := DeleteEvent[T]{
			Old:       MapTuple[T](rel, tuple),
			TableName: table,
			LSN:       lsn,
		}
		handler(event)
	}
}

func (r *Replicator) Start(ctx context.Context) error {
	replConfig, err := pgconn.ParseConfig(r.config.ConnString)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	replConfig.RuntimeParams["replication"] = "database"

	conn, err := pgconn.ConnectConfig(ctx, replConfig)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close(ctx)

	sysIdent, err := pglogrepl.IdentifySystem(ctx, conn)
	if err != nil {
		return fmt.Errorf("identify system: %w", err)
	}

	err = pglogrepl.StartReplication(ctx, conn, r.config.SlotName, sysIdent.XLogPos, pglogrepl.StartReplicationOptions{
		PluginArgs: []string{"proto_version '1'", fmt.Sprintf("publication_names '%s'", r.config.Publication)},
	})
	if err != nil {
		return fmt.Errorf("start replication: %w", err)
	}

	clientLSN := sysIdent.XLogPos
	nextStatusUpdate := time.Now().Add(r.config.Timeout)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if time.Now().After(nextStatusUpdate) {
			err = pglogrepl.SendStandbyStatusUpdate(ctx, conn, pglogrepl.StandbyStatusUpdate{WALWritePosition: clientLSN})
			if err != nil {
				return fmt.Errorf("send status update: %w", err)
			}
			nextStatusUpdate = time.Now().Add(r.config.Timeout)
		}

		ctxTimeout, cancel := context.WithTimeout(ctx, r.config.Timeout)
		msg, err := conn.ReceiveMessage(ctxTimeout)
		cancel()

		if err != nil {
			if pgconn.Timeout(err) {
				continue
			}
			return fmt.Errorf("receive message: %w", err)
		}

		copyData, ok := msg.(*pgproto3.CopyData)
		if !ok {
			continue
		}

		switch copyData.Data[0] {
		case pglogrepl.PrimaryKeepaliveMessageByteID:
			pka, err := pglogrepl.ParsePrimaryKeepaliveMessage(copyData.Data[1:])
			if err != nil {
				log.Printf("parse keepalive: %v", err)
				continue
			}
			if pka.ReplyRequested {
				nextStatusUpdate = time.Time{}
			}

		case pglogrepl.XLogDataByteID:
			xld, err := pglogrepl.ParseXLogData(copyData.Data[1:])
			if err != nil {
				log.Printf("parse xlog data: %v", err)
				continue
			}

			r.processWALData(xld)
			clientLSN = xld.WALStart + pglogrepl.LSN(len(xld.WALData))
		}
	}
}

func (r *Replicator) processWALData(xld pglogrepl.XLogData) {
	logicalMsg, err := pglogrepl.Parse(xld.WALData)
	if err != nil {
		log.Printf("parse logical message: %v", err)
		return
	}

	switch msg := logicalMsg.(type) {
	case *pglogrepl.RelationMessage:
		r.relations[msg.RelationID] = msg

	case *pglogrepl.InsertMessage:
		rel, ok := r.relations[msg.RelationID]
		if !ok {
			return
		}
		tableName := fmt.Sprintf("%s.%s", rel.Namespace, rel.RelationName)
		if h, ok := r.handlers[tableName]; ok && h.onInsert != nil {
			h.onInsert(rel, msg.Tuple, xld.WALStart)
		}

	case *pglogrepl.UpdateMessage:
		rel, ok := r.relations[msg.RelationID]
		if !ok {
			return
		}
		tableName := fmt.Sprintf("%s.%s", rel.Namespace, rel.RelationName)
		if h, ok := r.handlers[tableName]; ok && h.onUpdate != nil {
			h.onUpdate(rel, msg.OldTuple, msg.NewTuple, xld.WALStart)
		}

	case *pglogrepl.DeleteMessage:
		rel, ok := r.relations[msg.RelationID]
		if !ok {
			return
		}
		tableName := fmt.Sprintf("%s.%s", rel.Namespace, rel.RelationName)
		if h, ok := r.handlers[tableName]; ok && h.onDelete != nil {
			h.onDelete(rel, msg.OldTuple, xld.WALStart)
		}
	}
}
