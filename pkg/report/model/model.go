package model

import "github.com/xuri/excelize/v2"

type File struct {
	File     *excelize.File
	Path     string
	FileName string
	FileSize int64
}
