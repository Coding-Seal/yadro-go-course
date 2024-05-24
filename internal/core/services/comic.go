package services

import "yadro-go-course/internal/core/ports"

type Comic struct { // TODO: add logging
	ports.ComicsRepo
}

func NewComicService(repo ports.ComicsRepo) *Comic {
	return &Comic{repo}
}
