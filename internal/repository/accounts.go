package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/IvanMeln1k/go-bank-app-bank/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountsRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

func (r *AccountRepository) Create(ctx context.Context, userId uuid.UUID,
	account domain.Account) (uuid.UUID, error) {
	var id uuid.UUID

	query := fmt.Sprintf(`INSERT INTO %s a (id, money, user_id) VALUES
		((SELECT gen_random_uuid()), $1, $2) RETURNING a.id`, accountsTable)
	row := r.db.QueryRow(query)
	if err := row.Scan(&id); err != nil {
		logrus.Errorf("error insert into db account: %s", err)
		return id, ErrInternal
	}

	return id, nil
}

func (r *AccountRepository) Get(ctx context.Context, id uuid.UUID) (domain.Account, error) {
	var account domain.Account

	query := fmt.Sprintf(`SELECT * FROM %s a WHERE id=$1`, accountsTable)
	if err := r.db.Get(&account, query, id); err != nil {
		logrus.Errorf("error select account from db by id: %s", err)
		if errors.Is(sql.ErrNoRows, err) {
			return account, ErrAccountNotFound
		}
		return account, ErrInternal
	}

	return account, nil
}

func (r *AccountRepository) GetAll(ctx context.Context, userId uuid.UUID) ([]domain.Account, error) {
	var accounts []domain.Account

	query := fmt.Sprintf(`SELECT * FROM %s a WHERE user_id=$1`, accountsTable)
	if err := r.db.Select(&accounts, query, userId); err != nil {
		logrus.Errorf("error select accounts from db by user_id: %s", err)
		if errors.Is(sql.ErrNoRows, err) {
			return accounts, ErrAccountNotFound
		}
		return accounts, ErrInternal
	}

	return accounts, nil
}

func (r *AccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s a WHERE id=$1`, accountsTable)
	if _, err := r.db.Exec(query); err != nil {
		logrus.Errorf("error delete account from db by id: %s", err)
		return ErrInternal
	}

	return nil
}

func (r *AccountRepository) Update(ctx context.Context, id uuid.UUID,
	data domain.AccountUpdate) (domain.Account, error) {
	var account domain.Account

	values := []interface{}{}
	names := []string{}
	argId := 1

	addProperty := func(field string, value interface{}) {
		values = append(values, value)
		names = append(names, fmt.Sprintf("%s=$%d", field, argId))
		argId++
	}

	if data.Money != nil {
		addProperty("money", *data.Money)
	}

	values = append(values, id)
	setQuery := strings.Join(names, ", ")
	query := fmt.Sprintf(`UPDATE %s a SET %s WHERE id=$%d RETURNING a.*`,
		accountsTable, setQuery, argId)
	if err := r.db.Get(&account, query, values...); err != nil {
		logrus.Errorf("error update account into db by id: %s", err)
		return account, ErrInternal
	}

	return account, nil
}
