package storage

import (
	"context"
	"database/sql"
	"fmt"
)

type TxFunc func(*sql.Tx) error

func WithTransaction(ctx context.Context, db *sql.DB, fn TxFunc) error {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("işlem başlatılamadı: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {

		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("işlem geri alınırken hata oluştu: %v (orijinal hata: %w)", rbErr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("işlem tamamlanırken hata oluştu: %w", err)
	}

	return nil
}
