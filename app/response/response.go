package response

type Response struct {
	Code          int    `json:"code"`
	Message       string `json:"message"`
	Ok            bool   `json:"ok"`
	InternalError bool   `json:"internal_error"`
	Body          any    `json:"body,omitempty"`
}
