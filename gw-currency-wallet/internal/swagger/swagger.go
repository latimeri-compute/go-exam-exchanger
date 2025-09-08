package swagger

type ErrorResponse struct {
	Err string `json:"error"`
}

type ErrorUnauthorizedResponse struct {
	Err string `json:"error" example:"Unauthorized"`
}

type ErrorInsufficientFunds struct {
	Err string `json:"error" example:"Insufficient funds or invalid currencies"`
}

type InternalErrorResponse struct {
	Err string `json:"error" example:"Internal error"`
}

type MessageResponse struct {
	Mes string `json:"message"`
}

type ReturnUpdatedBalance struct {
	Message string
}
