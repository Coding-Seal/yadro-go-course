package xkcd

import "time"

type Comic struct {
	ID               int
	Date             time.Time
	News             string
	SafeTitle        string
	ImgURL           string
	Title            string
	Transcription    string
	AltTranscription string
	Link             string
}
