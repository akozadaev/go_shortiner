package shorten

import (
	"crypto/md5"
	"encoding/hex"
	"go_shurtiner/internal/shorten/model"
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

		hashStr = GetMD5Hash(sl.Source)
		link := Link{
			Source:    sl.Source,
			Shortened: host + hashStr,
		}
		links = append(links, link)
	}

	return &LinkResponse{
		Data: links,
	}
}

func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
