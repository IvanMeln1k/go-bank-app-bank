package handler

import "context"

func (h *Handler) CreateAccount(ctx context.Context,
	request CreateAccountRequestObject) (CreateAccountResponseObject, error) {
	return CreateAccount500JSONResponse{
		Message: "Not implemented",
	}, nil
}

func (h *Handler) GetAllAccounts(ctx context.Context,
	request GetAllAccountsRequestObject) (GetAllAccountsResponseObject, error) {
	return GetAllAccounts500JSONResponse{
		Message: "Not implemented",
	}, nil
}

func (h *Handler) GetAccountInfo(ctx context.Context, request GetAccountInfoRequestObject) (GetAccountInfoResponseObject, error) {
	return GetAccountInfo500JSONResponse{
		Message: "Not implemented",
	}, nil
}

func (h *Handler) DeleteAccount(ctx context.Context, request DeleteAccountRequestObject) (DeleteAccountResponseObject, error) {
	return DeleteAccount500JSONResponse{
		Message: "Not implemented",
	}, nil
}

func (h *Handler) Transfer(ctx context.Context, request TransferRequestObject) (TransferResponseObject, error) {
	return Transfer500JSONResponse{
		Message: "Not implemented",
	}, nil
}
