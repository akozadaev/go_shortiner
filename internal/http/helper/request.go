package helper

import "bytes"

func RequestHasJsonArray(json []byte) bool {
	json = bytes.TrimLeft(json, " \t\r\n")
	return len(json) > 0 && json[0] == '['
}
