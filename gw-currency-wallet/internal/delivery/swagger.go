package delivery

type errorResponse struct {
	Err string `json:"error"`
}

type errorUnauthorizedResponse struct {
	Err string `json:"error" example:"Unauthorized"`
}

type errorInsufficientFunds struct {
	Err string `json:"error" example:"Insufficient funds or invalid currencies"`
}

type internalErrorResponse struct {
	Err string `json:"error" example:"Internal error"`
}

type messageResponse struct {
	Mes string `json:"message"`
}

var jwtHead string = "BEARER "
