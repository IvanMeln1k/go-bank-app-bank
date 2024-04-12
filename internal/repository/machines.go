package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/IvanMeln1k/go-bank-app-bank/internal/domain"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/transactions"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type MachinesRepository struct {
	db        *sqlx.DB
	CtxGetter transactions.CtxGetterInterface
}

func NewMachinesRepository(db *sqlx.DB, CtxGetter transactions.CtxGetterInterface) *MachinesRepository {
	return &MachinesRepository{
		db:        db,
		CtxGetter: CtxGetter,
	}
}

func (r *MachinesRepository) Get(ctx context.Context, id uuid.UUID) (domain.Machine, error) {
	var machine domain.Machine

	query := fmt.Sprintf(`SELECT * FROM %s m WHERE id=$1`, machinesTable)
	if err := r.db.Get(&machine, query, id); err != nil {
		logrus.Errorf("error select machine from db by id: %s", err)
		if errors.Is(sql.ErrNoRows, err) {
			return machine, ErrMachineNotFound
		}
		return machine, ErrInternal
	}

	return machine, nil
}
