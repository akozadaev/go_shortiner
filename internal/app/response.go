package app

import (
	"go_shurtiner/internal/app/model"
	"go_shurtiner/pkg/mediator"
	"math/rand"
)

type Link struct {
	Source    string `json:"source"`
	Shortened string `json:"shortened"`
}

type LinkResponse struct {
	Data []Link `json:"data"`
}

func NewLinkResponse(sourcesLinks []model.CreateLink, host string) *LinkResponse {
	var hashStr string
	links := make([]Link, 0)
	for _, sl := range sourcesLinks {
		rnd := rand.Intn(2)
		switch rnd {
		case 0:
			md5 := new(mediator.Md5)
			hashStr = md5.Generate(sl.Source)
		case 1:
			base62 := new(mediator.Base62)
			hashStr = base62.Generate(sl.Source)
		}
		link := Link{
			Source:    sl.Source,
			Shortened: hashStr,
		}
		links = append(links, link)
	}

	return &LinkResponse{
		Data: links,
	}
}
