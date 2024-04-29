package services

import "yadro-go-course/internal/core/ports"

type Comic struct { // TODO: add logging
	ports.ComicsRepo
}

var _ ports.ComicService = (*Comic)(nil)
