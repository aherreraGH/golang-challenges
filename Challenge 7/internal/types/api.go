package types

// Handle requests with a simple object, JSON from the caller.
type PostBookRequest struct {
	Who  string `json:"who"`
	Code string `json:"code"`
}
