package requests

// Shortener constructs a request to shorten a given link
type Shortener struct {
	Link string `json:"link"`
}
