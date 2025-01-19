package mediator

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
)

func (r *Base62) Generate(str string) string {
	const chars string = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	base := chars
	var nb uint64
	lbase := len(base)
	le := len(str)
	for i := 0; i < le; i++ {
		mult := 1
		for j := 0; j < le-i-1; j++ {
			mult *= lbase
		}
		nb += uint64(strings.IndexByte(base, str[i]) * mult)
	}
	return strconv.FormatUint(nb, 10)
}

func (r *Md5) Generate(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])

}
