package fetcher

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

func NewMockXKCD(lastID int) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /info.0.json", LastIDHandler(lastID))
	mux.HandleFunc("GET /{id}/info.0.json", ComicHandler(lastID))

	return httptest.NewServer(mux)
}

type Comic struct {
	Day              int        `json:"day,string"`
	Month            time.Month `json:"month,string"`
	Year             int        `json:"year,string"`
	ID               int        `json:"num"`
	News             string     `json:"news"`
	SafeTitle        string     `json:"safe_title"`
	ImgURL           string     `json:"img"`
	Title            string     `json:"title"`
	Transcription    string     `json:"transcript"`
	AltTranscription string     `json:"alt"`
	Link             string     `json:"link"`
}

func LastIDHandler(lastID int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		comic := Comic{
			Day:              0,
			Month:            0,
			Year:             0,
			ID:               lastID,
			News:             "testStr",
			SafeTitle:        "testStr",
			ImgURL:           "testStr",
			Title:            "testStr",
			Transcription:    "testStr",
			AltTranscription: "testStr",
			Link:             "testStr",
		}
		_ = json.NewEncoder(w).Encode(comic)
	}
}

func ComicHandler(lastID int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			panic(err)
		}

		if id > lastID || id < 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		comic := Comic{
			Day:              0,
			Month:            0,
			Year:             0,
			ID:               id,
			News:             "testStr",
			SafeTitle:        "testStr",
			ImgURL:           "testStr",
			Title:            "testStr",
			Transcription:    "testStr",
			AltTranscription: "testStr",
			Link:             "testStr",
		}
		_ = json.NewEncoder(w).Encode(comic)
	}
}
