package xkcd

import "time"

type Comic struct {
	ID               int
	Title            string
	Date             time.Time
	ImgURL           string
	News             string
	SafeTitle        string
	Transcription    string
	AltTranscription string
	Link             string
}
