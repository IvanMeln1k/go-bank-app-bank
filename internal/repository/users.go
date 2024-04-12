package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/IvanMeln1k/go-bank-app-bank/internal/domain"
	"github.com/IvanMeln1k/go-bank-app-bank/pkg/transactions"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type UsersRepository struct {
	db        *sqlx.DB
	ctxGetter transactions.CtxGetterInterface
}

func NewUsersRepository(db *sqlx.DB, ctxGetter transactions.CtxGetterInterface) *UsersRepository {
	return &UsersRepository{
		db:        db,
		ctxGetter: ctxGetter,
	}
}

func (r *UsersRepository) Create(ctx context.Context, user domain.User) (uuid.UUID, error) {
	var id uuid.UUID
	tx := r.ctxGetter.TrOrDb(ctx, r.db)

	query := fmt.Sprintf(`INSERT INTO %s (id, surname, name, patronyc, email, hash_password) VALUES
		((SELECT gen_random_uuid()), $1, $2, $3, $4, $5) RETURNING id`, usersTable)
	row := tx.QueryRowxContext(ctx, query, user.Surname, user.Name, user.Patronyc, user.Email,
		user.Password)
	if err := row.Scan(&id); err != nil {
		logrus.Errorf("error insert user into db when creating user: %s", err)
		return id, ErrInternal
	}

	return id, nil
}

func (r *UsersRepository) get(ctx context.Context, field string,
	data interface{}) (domain.User, error) {
	var user domain.User

	query := fmt.Sprintf(`SELECT * FROM %s u WHERE %s=$1`, usersTable, field)
	if err := r.db.Get(&user, query, data); err != nil {
		logrus.Errorf("error getting user from db by %s: %s", field, err)
		if errors.Is(sql.ErrNoRows, err) {
			return user, ErrUserNotFound
		}
		return user, ErrInternal
	}

	return user, nil
}

func (r *UsersRepository) Get(ctx context.Context, id uuid.UUID) (domain.User, error) {
	return r.get(ctx, "id", id)
}

func (r *UsersRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return r.get(ctx, "email", email)
}

func (r *UsersRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s u WHERE id=$1", usersTable)
	tx := r.ctxGetter.TrOrDb(ctx, r.db)

	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		logrus.Errorf("error delete user from db by id: %s", err)
		return ErrInternal
	}

	return nil
}

func (r *UsersRepository) Update(ctx context.Context, id uuid.UUID,
	data domain.UserUpdate) (domain.User, error) {
	var user domain.User
	tx := r.ctxGetter.TrOrDb(ctx, r.db)

	values := []interface{}{}
	names := []string{}
	argId := 1

	addField := func(field string, value interface{}) {
		values = append(values, value)
		names = append(names, fmt.Sprintf("%s=$%d", field, argId))
		argId++
	}

	if data.Surname != nil {
		addField("surname", *data.Surname)
	}
	if data.Name != nil {
		addField("name", *data.Name)
	}
	if data.Patronyc != nil {
		addField("patronyc", *data.Patronyc)
	}
	if data.Email != nil {
		addField("email", *data.Email)
	}
	if data.Password != nil {
		addField("hash_password", *data.Password)
	}
	if data.Verified != nil {
		addField("verified", *data.Verified)
	}
	values = append(values, id)

	querySet := strings.Join(names, ", ")
	query := fmt.Sprintf("UPDATE %s u SET %s WHERE id=$%d RETURNING u.*", usersTable, querySet, argId)

	row := tx.QueryRowxContext(ctx, query, values...)
	if err := row.StructScan(&user); err != nil {
		logrus.Errorf("error updating user into db by id: %s", err)
		return user, ErrInternal
	}

	return user, nil
}
