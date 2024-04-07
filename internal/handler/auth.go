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

func (h *Handler) RefreshTokens(ctx context.Context,
	request RefreshTokensRequestObject) (RefreshTokensResponseObject, error) {
	return RefreshTokens500JSONResponse{
		Message: "Not implemented",
	}, nil
}

func (h *Handler) Logout(ctx context.Context,
	request LogoutRequestObject) (LogoutResponseObject, error) {
	return Logout500JSONResponse{
		Message: "Not implemented",
	}, nil
}

func (h *Handler) LogoutAll(ctx context.Context,
	request LogoutAllRequestObject) (LogoutAllResponseObject, error) {
	return LogoutAll500JSONResponse{
		Message: "Not implemented",
	}, nil
}

func (h *Handler) GetMe(ctx context.Context,
	request GetMeRequestObject) (GetMeResponseObject, error) {
	return GetMe500JSONResponse{
		Message: "Not implemented",
	}, nil
}
