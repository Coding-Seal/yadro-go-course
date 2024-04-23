package app

import (
	"log"
	"maps"
)

func (a *App) LoadComics() error {
	comics, err := a.db.Read()

	if err != nil {
		return err
	}

	maps.Copy(a.comics, comics)

	return nil
}
func (a *App) SaveComics() {
	err := a.db.SaveAll(a.comics)
	if err != nil {
		log.Println("error saving comics: ", err)
	}
}
