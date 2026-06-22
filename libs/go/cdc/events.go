package cdc

import "github.com/jackc/pglogrepl"

type InsertEvent[T any] struct {
	New       T
	TableName string
	LSN       pglogrepl.LSN
}

type UpdateEvent[T any] struct {
	Old       T
	New       T
	TableName string
	LSN       pglogrepl.LSN
}

type DeleteEvent[T any] struct {
	Old       T
	TableName string
	LSN       pglogrepl.LSN
}

type BeginEvent struct {
	Xid uint32
	LSN pglogrepl.LSN
}

type CommitEvent struct {
	LSN pglogrepl.LSN
}
