package handler

import "context"

func (h *Handler) Deposit(ctx context.Context,
	request DepositRequestObject) (DepositResponseObject, error) {
	return Deposit500JSONResponse{
		Message: "Not implemented",
	}, nil
}

func (h *Handler) CashOut(ctx context.Context,
	request CashOutRequestObject) (CashOutResponseObject, error) {
	return CashOut500JSONResponse{
		Message: "Not implemented",
	}, nil
}
