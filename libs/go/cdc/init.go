package cdc

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

type ReplicaIdentity string

const (
	ReplicaIdentityDefault ReplicaIdentity = "DEFAULT"
	ReplicaIdentityFull    ReplicaIdentity = "FULL"
	ReplicaIdentityNothing ReplicaIdentity = "NOTHING"
)

type TableConfig struct {
	Name            string
	ReplicaIdentity ReplicaIdentity
	Columns         []string
}

type InitOptions struct {
	DropIfExists bool
}

func (r *Replicator) Init(ctx context.Context, opts InitOptions) error {
	db, err := pgx.Connect(ctx, r.config.ConnString)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer db.Close(ctx)

	if opts.DropIfExists {
		if err := dropSlotIfExists(ctx, db, r.config.SlotName); err != nil {
			return err
		}
		if err := dropPublicationIfExists(ctx, db, r.config.Publication); err != nil {
			return err
		}
	}

	for _, table := range r.tables {
		if table.ReplicaIdentity != "" {
			_, err := db.Exec(ctx, fmt.Sprintf("ALTER TABLE %s REPLICA IDENTITY %s", table.Name, table.ReplicaIdentity))
			if err != nil {
				return fmt.Errorf("set replica identity for %s: %w", table.Name, err)
			}
		}
	}

	if err := createPublicationIfNotExists(ctx, db, r.config.Publication, r.tables); err != nil {
		return err
	}

	if err := createSlotIfNotExists(ctx, db, r.config.SlotName); err != nil {
		return err
	}

	return nil
}

func dropSlotIfExists(ctx context.Context, db *pgx.Conn, slotName string) error {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = $1)", slotName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check slot: %w", err)
	}
	if exists {
		_, err = db.Exec(ctx, "SELECT pg_drop_replication_slot($1)", slotName)
		if err != nil {
			return fmt.Errorf("drop slot: %w", err)
		}
	}
	return nil
}

func dropPublicationIfExists(ctx context.Context, db *pgx.Conn, pubName string) error {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_publication WHERE pubname = $1)", pubName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check publication: %w", err)
	}
	if exists {
		_, err = db.Exec(ctx, fmt.Sprintf("DROP PUBLICATION %s", pubName))
		if err != nil {
			return fmt.Errorf("drop publication: %w", err)
		}
	}
	return nil
}

func createPublicationIfNotExists(ctx context.Context, db *pgx.Conn, pubName string, tables []TableConfig) error {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_publication WHERE pubname = $1)", pubName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check publication: %w", err)
	}
	if exists {
		return nil
	}

	var tableDefs []string
	for _, t := range tables {
		if len(t.Columns) > 0 {
			tableDefs = append(tableDefs, fmt.Sprintf("%s (%s)", t.Name, strings.Join(t.Columns, ", ")))
		} else {
			tableDefs = append(tableDefs, t.Name)
		}
	}

	sql := fmt.Sprintf("CREATE PUBLICATION %s FOR TABLE %s", pubName, strings.Join(tableDefs, ", "))
	_, err = db.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("create publication: %w", err)
	}
	return nil
}

func createSlotIfNotExists(ctx context.Context, db *pgx.Conn, slotName string) error {
	var exists bool
	err := db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM pg_replication_slots WHERE slot_name = $1)", slotName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check slot: %w", err)
	}
	if !exists {
		_, err = db.Exec(ctx, "SELECT pg_create_logical_replication_slot($1, 'pgoutput')", slotName)
		if err != nil {
			return fmt.Errorf("create slot: %w", err)
		}
	}
	return nil
}

func (r *Replicator) Drop(ctx context.Context) error {
	db, err := pgx.Connect(ctx, r.config.ConnString)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer db.Close(ctx)

	if err := dropSlotIfExists(ctx, db, r.config.SlotName); err != nil {
		return err
	}
	if err := dropPublicationIfExists(ctx, db, r.config.Publication); err != nil {
		return err
	}
	return nil
}
