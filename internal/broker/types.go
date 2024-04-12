package broker

import "github.com/google/uuid"

type sendEmailVerificationMessageTask struct {
	Email string `json:"email"`
}

type transferTask struct {
	EmailFrom string    `json:"emailFrom"`
	EmailTo   string    `json:"emailTo"`
	AccIdFrom uuid.UUID `json:"accIdFrom"`
	AccIdTo   uuid.UUID `json:"accIdTo"`
	Amount    int       `json:"amount"`
}

type cashoutTask struct {
	MachineId uuid.UUID `json:"machineId"`
	Email     string    `json:"email"`
	AccId     uuid.UUID `json:"accId"`
	Amount    int       `json:"amount"`
	NewMoney  int       `json:"new_money"`
}

type depositTask struct {
	MachineId uuid.UUID `json:"machineId"`
	Email     string    `json:"email"`
	AccId     uuid.UUID `json:"accId"`
	Amount    int       `json:"amount"`
	NewMoney  int       `json:"new_money"`
}
