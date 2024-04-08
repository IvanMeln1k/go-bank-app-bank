package handler

import "context"

func (h *Handler) SignUp(ctx context.Context,
	request SignUpRequestObject) (SignUpResponseObject, error) {
	return SignUp500JSONResponse{
		Message: "Not implemented",
	}, nil
}

func (h *Handler) SignIn(ctx context.Context,
	request SignInRequestObject) (SignInResponseObject, error) {
	return SignIn500JSONResponse{
		Message: "Not implemented",
	}, nil
}

func (h *Handler) GetMe(ctx context.Context,
	request GetMeRequestObject) (GetMeResponseObject, error) {
	return GetMe500JSONResponse{
		Message: "Not implemented",
	}, nil
}
