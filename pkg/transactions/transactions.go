package transactions

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

var (
	ErrCreateTr  = errors.New("error create transaction")
	ErrCommitTr  = errors.New("error commit transaction")
	ErrInvalidTr = errors.New("invalid transaction")
)

type CtxKey string

const (
	keyTx CtxKey = "tx"
)

type ManagerInterface interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type CtxManagerInterface interface {
	Tr(ctx context.Context) (*sqlx.Tx, error)
}

type CtxGetterInterface interface {
	TrOrDb(ctx context.Context, db *sqlx.DB) sqlx.ExtContext
}

type Manager struct {
	db *sqlx.DB
}

func NewManager(db *sqlx.DB) *Manager {
	return &Manager{
		db: db,
	}
}

func (m *Manager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	hasExternalTx := true
	tx, ok := ctx.Value(keyTx).(*sqlx.Tx)

	if !ok || tx == nil {
		hasExternalTx = false
		newTx, err := m.db.Beginx()
		if err != nil {
			return ErrCreateTr
		}
		ctx = context.WithValue(ctx, keyTx, newTx)
		tx = newTx
	}

	err := fn(ctx)

	if !hasExternalTx && err != nil {
		tx.Rollback()
		return err
	}

	if !hasExternalTx {
		err := tx.Commit()
		if err != nil {
			return ErrCommitTr
		}
	}

	return nil
}

type CtxManager struct {
}

func NewCtxManager() *CtxManager {
	return &CtxManager{}
}

func (cm *CtxManager) Tr(ctx context.Context) (*sqlx.Tx, error) {
	tx, ok := ctx.Value(keyTx).(*sqlx.Tx)
	if !ok {
		return nil, ErrInvalidTr
	}
	return tx, nil
}

type CtxGetter struct {
	ctxManager CtxManagerInterface
}

func NewCtxGetter(ctxManager CtxManagerInterface) *CtxGetter {
	return &CtxGetter{
		ctxManager: ctxManager,
	}
}

func (tg *CtxGetter) TrOrDb(ctx context.Context, db *sqlx.DB) sqlx.ExtContext {
	tx, err := tg.ctxManager.Tr(ctx)
	if err != nil {
		return db
	}
	return tx
}
