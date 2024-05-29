package fetcher

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testComic = Comic{
	Day:              0,
	Month:            0,
	Year:             0,
	ID:               0,
	News:             "testStr",
	SafeTitle:        "testStr",
	ImgURL:           "testStr",
	Title:            "testStr",
	Transcription:    "testStr",
	AltTranscription: "testStr",
	Link:             "testStr",
}

func Test_Happy_ComicHandler(t *testing.T) {
	h := ComicHandler(5)
	r := httptest.NewRequest("Get", "/", nil)
	r.SetPathValue("id", "1")

	w := httptest.NewRecorder()
	h(w, r)
	resp := w.Result()
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	var comic Comic

	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&comic))

	testComic.ID = 1
	assert.Equal(t, testComic, comic)
}

func Test_Unhappy_ComicHandler(t *testing.T) {
	h := ComicHandler(5)
	r := httptest.NewRequest("Get", "/", nil)
	r.SetPathValue("id", "0")

	w := httptest.NewRecorder()
	h(w, r)
	resp := w.Result()
	assert.Equal(t, resp.StatusCode, http.StatusNotFound)
}

func TestLastIDHandler(t *testing.T) {
	h := LastIDHandler(5)
	r := httptest.NewRequest("Get", "/", nil)
	w := httptest.NewRecorder()
	h(w, r)
	resp := w.Result()
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	var comic Comic

	assert.NoError(t, json.NewDecoder(resp.Body).Decode(&comic))
	assert.Equal(t, 5, comic.ID)
}

func TestNewMockXKCD(t *testing.T) {
	assert.NotNil(t, NewMockXKCD(5))
}
