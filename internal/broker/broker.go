package broker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var (
	emailVerificationQueue = "queue:verification:email"
	cashoutQueue           = "queue:cashout"
	depositQueue           = "queue:deposit"
)

var (
	ErrInternal = errors.New("error internal")
)

type BrokerInterface interface {
	WriteVerificationTask(ctx context.Context, email string) error
	WriteCashoutTask(ctx context.Context, machineId uuid.UUID, email string, accId uuid.UUID,
		amount int, newMoney int) error
	WriteDepositTask(ctx context.Context, machineId uuid.UUID, email string, accId uuid.UUID,
		amount int, newMoney int) error
}

type Broker struct {
	rdb *redis.Client
}

type Deps struct {
	RDB *redis.Client
}

func NewBroker(deps Deps) *Broker {
	return &Broker{
		rdb: deps.RDB,
	}
}

func (b *Broker) writeTask(ctx context.Context, key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		logrus.Errorf("error marshaling data when write task into %s: %s", key, err)
		return ErrInternal
	}
	_, err = b.rdb.RPush(ctx, key, string(jsonData)).Result()
	if err != nil {
		logrus.Errorf("error whire transfer task to redis: %s", err)
		return ErrInternal
	}
	return nil
}

func (b *Broker) WriteVerificationTask(ctx context.Context, email string) error {
	return b.writeTask(ctx, emailVerificationQueue, sendEmailVerificationMessageTask{
		Email: email,
	})
}

func (b *Broker) WriteCashoutTask(ctx context.Context, machineId uuid.UUID,
	email string, accId uuid.UUID, amount int, newMoney int) error {
	data := cashoutTask{
		MachineId: machineId,
		Email:     email,
		AccId:     accId,
		Amount:    amount,
		NewMoney:  newMoney,
	}
	return b.writeTask(ctx, cashoutQueue, data)
}

func (b *Broker) WriteDepositTask(ctx context.Context, machineId uuid.UUID, email string,
	accId uuid.UUID, amount int, newMoney int) error {
	data := depositTask{
		MachineId: machineId,
		Email:     email,
		AccId:     accId,
		Amount:    amount,
		NewMoney:  newMoney,
	}
	return b.writeTask(ctx, depositQueue, data)
}
