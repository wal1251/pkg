// Package transaction содержит удобную обертку для работы с транзакциями
//
// Для использования обертки вам потребуется создать ent.Client , который создаст на основе вашей необходимые клиенты
// После вы должны реализовать структуру Tx[T]
//
// Для этого есть два подхода
// 1. Функциональный:
//
//	func MakeTxFn(client *ent.Client) func(ctx context.Context, options *sql.TxOptions) (*Tx[ent.Client], error) {
//		return func(ctx context.Context, options *sql.TxOptions) (*Tx[ent.Client], error) {
//			transaction, err := client.BeginTx(ctx, options)
//			if err != nil {
//				logs.FromContext(ctx).Err(err).Msg("failed to open transaction")
//
//				return nil, err
//			}
//			return &Tx[ent.Client]{
//				Transaction: transaction,
//				client:      client,
//			}, nil
//		}
//	}
//
// 2. Процедурный:
//
//	func NewTx(ctx context.Context, client *ent.Client, options *sql.TxOptions) (*Tx[ent.Client], error) {
//		transaction, err := client.BeginTx(ctx, options)
//		if err != nil {
//			logs.FromContext(ctx).Err(err).Msg("failed to open transaction")
//
//			return nil, err
//		}
//		return &Tx[ent.Client]{
//			Transaction: transaction,
//		}, nil
//	}
//
// После данной имплементации вы можете использовать взаимодействии с БД:
//
//	tx, err := MakeTxFn(s.Client)(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
//
//	if err != nil {
//		return nil, err
//	}
//	defer tx.Done(ctx)
//
//	entToken, err := tx.
//		GetClient().
//		Token.
//		Query().
//		WithAccessToken().
//		Where(token.UserIDEQ(userID), token.AccessTokenIDIsNil()).
//		Only(ctx)
//	if err != nil && !ent.IsNotFound(err) {
//		return nil, err
//	}
package transaction

import (
	"context"

	"github.com/wal1251/pkg/core/logs"
)

type (
	Transaction interface {
		Rollback() error
		Commit() error
	}

	Tx[T any] struct {
		Transaction
		Client    *T
		cancelled bool
	}
)

func (t *Tx[T]) CancelOnError(ctx context.Context, err error) error {
	if err != nil {
		return t.Cancel(ctx, err)
	}

	return nil
}

func (t *Tx[T]) Cancel(ctx context.Context, err error) error {
	if t.cancelled {
		return err
	}

	t.cancelled = true

	logger := logs.FromContext(ctx)
	if err := t.Rollback(); err != nil {
		logger.Err(err).Msg("transaction rollback failed")
	}
	logger.Warn().Err(err).Msg("transaction cancelled due to error")

	return err
}

func (t *Tx[T]) Done(ctx context.Context) {
	if t.cancelled {
		return
	}
	if err := t.Commit(); err != nil {
		logs.FromContext(ctx).Err(err).Msg("transaction commit failed")
	}
}
