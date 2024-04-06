package comic

type Comic struct {
	ID       int      `json:"-"`
	ImgURL   string   `json:"url"`
	Keywords []string `json:"keywords"`
}
