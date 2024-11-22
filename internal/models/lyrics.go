package models

// request
type Song struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

// response
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
