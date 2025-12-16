// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package config

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/playbymail/ottoapp/backend/domains"
	"github.com/playbymail/ottoapp/backend/stores/sqlite"
	"github.com/playbymail/ottoapp/backend/stores/sqlite/sqlc"
)

type Service struct {
	db *sqlite.DB
}

func New(db *sqlite.DB) (*Service, error) {
	return &Service{db: db}, nil
}

func (s *Service) CreateKeyValue(key, value string) error {
	now := time.Now().UTC()
	createdAt, updatedAt := now.Unix(), now.Unix()
	err := s.db.Queries().CreateConfigKeyValue(s.db.Context(), sqlc.CreateConfigKeyValueParams{
		Key:       key,
		Value:     value,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return errors.Join(domains.ErrDatabaseError, err)
	}
	return nil
}

func (s *Service) ReadKeyValue(key string) (value string, err error) {
	value, err = s.db.Queries().ReadConfigKeyValue(s.db.Context(), key)
	if err != nil {
		value = ""
	}
	return value, err
}

func (s *Service) UpdateKeyValue(key, value string) error {
	now := time.Now().UTC()
	createdAt, updatedAt := now.Unix(), now.Unix()
	_, err := s.db.Queries().UpsertConfigKeyValue(s.db.Context(), sqlc.UpsertConfigKeyValueParams{
		Key:       key,
		Value:     value,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
	if err != nil {
		return errors.Join(domains.ErrDatabaseError, err)
	}
	return nil
}

func (s *Service) UpdateKeyValuePairs(args ...string) error {
	if len(args) == 0 {
		return nil
	} else if len(args)%2 != 0 {
		return fmt.Errorf("invalid set of key+value pairs")
	}
	now := time.Now().UTC()
	createdAt, updatedAt := now.Unix(), now.Unix()

	// start transaction
	ctx := s.db.Context()
	tx, err := s.db.Stdlib().BeginTx(ctx, nil)
	if err != nil {
		log.Printf("config: begin %v\n", err)
		return errors.Join(domains.ErrDatabaseError, err)
	}
	defer tx.Rollback() // rollback if we return early; harmless after commit
	qtx := s.db.Queries().WithTx(tx)

	for n := 0; n < len(args); n += 2 {
		_, err := qtx.UpsertConfigKeyValue(ctx, sqlc.UpsertConfigKeyValueParams{
			Key:       args[n],
			Value:     args[n+1],
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
		if err != nil {
			return errors.Join(domains.ErrDatabaseError, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("config: commit %v\n", err)
		return errors.Join(domains.ErrDatabaseError, err)
	}

	return nil
}

func (s *Service) DeleteKeyValue(key string) error {
	err := s.db.Queries().DeleteConfigKeyValue(s.db.Context(), key)
	if err != nil {
		return errors.Join(domains.ErrDatabaseError, err)
	}
	return nil
}
