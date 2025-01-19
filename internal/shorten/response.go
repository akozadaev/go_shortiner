package shorten

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

type Link struct {
	Source    string `json:"source"`
	Shortened string `json:"shortened"`
}

type LinkResponse struct {
	Data Link `json:"data"`
}

func NewLinkResponse(source string) *LinkResponse {
	var hash = sha256.Sum256([]byte(source))
	var hashStr, _ = fmt.Printf("%x", hash)
	return &LinkResponse{
		Data: Link{
			Source:    source,
			Shortened: strconv.Itoa(hashStr),
		},
	}
}
