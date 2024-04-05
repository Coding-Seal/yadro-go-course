package xkcd

type Comic struct {
	ID       int      `json:"-"`
	Title    string   `json:"-"`
	ImgURL   string   `json:"url"`
	Keywords []string `json:"keywords"`
}
