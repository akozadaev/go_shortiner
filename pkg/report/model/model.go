package model

import "github.com/xuri/excelize/v2"

// File структура файла
type File struct {
	File     *excelize.File
	Path     string
	FileName string
	FileSize int64
}
