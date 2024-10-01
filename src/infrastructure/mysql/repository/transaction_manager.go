package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/shinya-ac/server1Q1A/application/transaction"
	"github.com/shinya-ac/server1Q1A/infrastructure/mysql/db"
)

type transactionManager struct{}

func NewTransactionManager() transaction.TransactionManager {
	return &transactionManager{}
}

func (tm *transactionManager) RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	dbcon := db.GetDB()

	tx, err := dbcon.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %w", err)
	}

	err = fn(ctx)
	if err != nil {
		log.Printf("エラー発生、ロールバックを実行します: %v\n", err)
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("トランザクションエラー: %v, ロールバックエラー: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("トランザクションのコミットに失敗しました: %w", err)
	}

	return nil
}
