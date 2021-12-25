package fakedb

import (
	"context"

	db "github.com/otmosina/simplebank/db/sqlc"
)

type Store interface {
	GetAccount(ctx context.Context, id int64) (db.Account, error)
}

type MemStore struct {
	data map[int64]db.Account
}

func (store *MemStore) GetAccount(ctx context.Context, id int64) (db.Account, error) {
	return store.data[id], nil
}

// type newStore Store

func Test(s Store) error {
	_, err := s.GetAccount(context.Background(), 1)
	return err
}

func mainmain() {
	mem := &MemStore{}
	Test(mem)
}
