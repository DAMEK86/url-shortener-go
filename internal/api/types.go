package api

type ErrorResponse struct {
	ErrorID      string `json:"errorId"`
	ErrorMessage string `json:"errorMessage"`
}

const (
	InvalidParameterErrorID = "INVALID_PARAMETER"
	EntityNotFoundErrorID   = "ENTITY_NOT_FOUND"
	UnexpectedErrorID       = "UNEXPECTED_ERROR"
)

type UrlRequest struct {
	Url string `json:"url" description:"URL to shorten"`
}

type UrlResponse struct {
	Url string `json:"url" description:"shorten URL"`
}
