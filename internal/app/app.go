package app

import (
	"maps"
	"yadro-go-course/pkg/comic"
	"yadro-go-course/pkg/database"
	"yadro-go-course/pkg/words"
	"yadro-go-course/pkg/xkcd"
)

type App struct {
	fetcher *xkcd.Fetcher
	stemmer *words.Stemmer
	db      *database.JsonDB
	comics  map[int]*comic.Comic
}

func (a *App) LoadComics() {
	maps.Copy(a.comics, a.db.Read())
}
func (a *App) SaveComics() {
	a.db.Save(a.comics)
}
func (a *App) FetchRemainingComics(lastID int) {
	var missingComics []int
	for id := 1; id <= lastID; id++ {
		if _, ok := a.comics[id]; !ok {
			missingComics = append(missingComics, id)
		}
	}
}
func (a *App) GetLastComicID() int {

}
func (a *App) RefreshDB(lastID int) int {
}
func (a *App) PrintComics(num int) {

}
