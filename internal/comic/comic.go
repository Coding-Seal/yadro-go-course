package comic

import "time"

type Comic struct {
	ID       int            `json:"id"`
	Title    string         `json:"title"`
	Date     time.Time      `json:"date"`
	ImgURL   string         `json:"img_url"`
	Keywords map[string]int `json:"keywords"`
}
