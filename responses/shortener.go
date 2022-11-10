package responses

// Shortener a basic shortener response struct
type Shortener struct {
	Link  string `json:"link"`
	Short string `json:"short"`
}
