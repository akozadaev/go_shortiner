package helper

import (
	"errors"
	"fmt"
	"github.com/xuri/excelize/v2"
	"os"
)

const DefaultLanguage = "ru"
const DefaultTimeZone = "Europe/Moscow"
const ReportStringStyle string = "string"

// ReportStringGrayStyle добавляет серый текст
const ReportStringGrayStyle string = "string_gray"
const ReportDateStyle string = "date"
const ReportDateTimeStyle string = "date_time"
const ReportTimeStyle string = "time"
const ReportTimeAsStringStyle string = "time_string"
const ReportTimeMinuteStyle string = "time_minute"
const ReportDecimalStyle string = "decimal"
const ReportIntegerStyle string = "integer"
const fontFamily = "Times New Roman"
const ReportDefaultRowSize = 14.0
const reportDefaultFontSize = 14.0

type ColWidthOpts struct {
	Min, Max int
	Width    float64
}

func CreateSheets(file *excelize.File, sheets []string) (map[int]string, error) {
	sheetMap := map[int]string{}
	if len(sheets) == 0 {
		return sheetMap, errors.New("no sheet names were transferred")
	}
	var sheet string
	var name string
	for _, sheet = range sheets {
		name = sheet
		sheetId, err := file.NewSheet(name)
		if err != nil {
			return nil, err
		}
		sheetMap[sheetId] = name
	}

	_ = file.DeleteSheet("Sheet1")

	return sheetMap, nil
}

func GetReportExcelizeHederStyle(file *excelize.File) (int, error) {
	styleID, err := file.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Size: reportDefaultFontSize, Family: fontFamily}, Alignment: &excelize.Alignment{
		Vertical:    "center",
		WrapText:    true,
		ShrinkToFit: true,
	}})

	return styleID, err
}

// Создание форматов и стилей для ячеек отчета.
// Важно! Чтобы стиль применился корректно, нужно, чтобы он совпал с типом записываемого в ячейку значения.
// https://xuri.me/excelize/en/style.html#number_format
func GetReportExcelizeStyles(file *excelize.File) (map[string]int, error) {
	styles := make(map[string]int)

	fmtDate := "dd.mm.yyyy"
	dateStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &fmtDate, Font: &excelize.Font{Bold: false, Size: reportDefaultFontSize, Family: fontFamily},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		}})
	if err != nil {
		return nil, err
	}
	styles[ReportDateStyle] = dateStyle

	fmtDateTime := "dd.mm.yyyy hh:mm"
	dateTimeStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &fmtDateTime, Font: &excelize.Font{Bold: false, Size: reportDefaultFontSize, Family: fontFamily},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		}})
	if err != nil {
		return nil, err
	}
	styles[ReportDateTimeStyle] = dateTimeStyle

	fmtTimeMinute := "hh:mm"
	dateTimeMinuteStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &fmtTimeMinute, Font: &excelize.Font{Bold: false, Size: reportDefaultFontSize, Family: fontFamily},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		}})
	if err != nil {
		return nil, err
	}
	styles[ReportTimeMinuteStyle] = dateTimeMinuteStyle

	fmtTime := "hh:mm:ss"
	timeStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &fmtTime, Font: &excelize.Font{Bold: false, Size: reportDefaultFontSize, Family: fontFamily},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		}})
	if err != nil {
		return nil, err
	}
	styles[ReportTimeStyle] = timeStyle

	fmtTimeAsString := "@" // 49
	timeAsStringStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &fmtTimeAsString, Font: &excelize.Font{Bold: false, Size: reportDefaultFontSize, Family: fontFamily},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		}})
	if err != nil {
		return nil, err
	}
	styles[ReportTimeAsStringStyle] = timeAsStringStyle

	fmtStr := "@" // 49
	stringStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &fmtStr, Font: &excelize.Font{Bold: false, Size: reportDefaultFontSize, Family: fontFamily}})
	if err != nil {
		return nil, err
	}
	styles[ReportStringStyle] = stringStyle

	fmtStrGray := "g"
	stringGrayStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &fmtStrGray, Font: &excelize.Font{Bold: false, Size: reportDefaultFontSize, Family: fontFamily, Color: "777777"}})
	if err != nil {
		return nil, err
	}
	styles[ReportStringGrayStyle] = stringGrayStyle

	fmtNum := "0" //1
	numStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &fmtNum, Font: &excelize.Font{Bold: false, Size: reportDefaultFontSize, Family: fontFamily},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		}})
	if err != nil {
		return nil, err
	}
	styles[ReportIntegerStyle] = numStyle

	decStr := "0.00" //2
	decStyle, err := file.NewStyle(&excelize.Style{CustomNumFmt: &decStr, Font: &excelize.Font{Bold: false, Size: reportDefaultFontSize, Family: fontFamily},
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		}})
	if err != nil {
		return nil, err
	}
	styles[ReportDecimalStyle] = decStyle

	_ = file.SetDefaultFont(fontFamily)

	return styles, nil
}

// Настройка ширины столбцов для листа для которого создан передаваемый в качестве параметра StreamWriter
func SetColWidth(streamWriter *excelize.StreamWriter, opts []ColWidthOpts) error {
	for _, opt := range opts {
		err := streamWriter.SetColWidth(opt.Min, opt.Max, opt.Width)
		if err != nil {
			return err
		}
	}

	return nil
}

func SetHeader(streamWrite *excelize.StreamWriter, head []interface{}) error {
	if err := streamWrite.SetRow("A1", head, excelize.RowOpts{
		Height: 60,
		Hidden: false,
	}); err != nil {
		return err
	}
	return nil
}

func SetSheetHeader(streamWrite *excelize.StreamWriter, headers any, styleHeader int) (int, error) {
	var headerRows []string
	dataType := fmt.Sprintf("%T", headers)
	f := func() (int, error) {
		var head = make([]interface{}, 0)
		for _, headerRow := range headerRows {
			head = append(head, excelize.Cell{StyleID: styleHeader, Value: headerRow})
		}

		err := SetHeader(streamWrite, head)
		if err != nil {
			return 0, err
		}

		return len(head), nil
	}

	switch dataType {
	case "map[string][]string":
		headerMap := headers.(map[string][]string)
		rows, ok := headerMap[os.Getenv("LANGUAGE")]
		if !ok {
			rows = headerMap[DefaultLanguage]
		}
		headerRows = rows

	case "[]string":
		headerRows = headers.([]string)
	default:
		return 0, errors.New("incorrect header data type")
	}
	cnt, err := f()
	return cnt, err
}

// AutoFilter Добавление автофильтров на указанном листе от ячейки A1 до ячейки с индекксами colInd, rowInd
func AutoFilter(file *excelize.File, sheetName string, colInd, rowInd int) error {
	cell, err := excelize.CoordinatesToCellName(colInd, rowInd)
	if err != nil {
		return err
	}
	if err := file.AutoFilter(sheetName, fmt.Sprintf("A1:%s", cell), []excelize.AutoFilterOptions{}); err != nil {
		return err
	}
	return nil
}
