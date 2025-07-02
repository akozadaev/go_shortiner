package report

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"go_shurtiner/internal/app/model"
	"go_shurtiner/pkg/report/helper"
	fileModel "go_shurtiner/pkg/report/model"
	"time"
)

var reportHeader = map[string][]string{
	"ru": {
		"Логин/Email",
		"ФИО",
		"Оригинальная ссылка",
		"Укороченная ссылка",
		"Дата создания",
	},
	"en": {
		"Login/Email",
		"Full name",
		"Original link",
		"Shortened link",
		"Creation date",
	},
}

var colWidthOpts = []helper.ColWidthOpts{
	{
		Min:   1,
		Max:   1,
		Width: 25.0,
	},
	{
		Min:   2,
		Max:   2,
		Width: 30.0,
	},
	{
		Min:   3,
		Max:   3,
		Width: 25.0,
	},
	{
		Min:   4,
		Max:   4,
		Width: 40.0,
	},
	{
		Min:   5,
		Max:   5,
		Width: 15.0,
	},
}

type StatReport interface {
	GenerateReport(data *[]model.PreparedReport) error
}

func NewStatReport() StatReport {
	return &statReport{}
}

type statReport struct {
}

func (r *statReport) GenerateReport(data *[]model.PreparedReport) error {
	var err error
	sheets := []string{
		"Статистика",
	}

	timezone := helper.DefaultTimeZone
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		loc = time.Local
	}

	file := excelize.NewFile()

	defer func() {
		err = file.Close()
	}()

	var fileInfo fileModel.File
	fileInfo.File = file
	fileInfo.Path = fmt.Sprintf("/tmp/%s.xlsx", uuid.New().String())
	currentDate := time.Now().In(loc).Format("02-01-2006_at_15:04_UTC-07:00")
	fileInfo.FileName = "./reports/report_" + currentDate + ".xlsx"

	_, err = helper.CreateSheets(file, sheets)
	if err != nil {
		return err
	}

	styleHeader, err := helper.GetReportExcelizeHederStyle(file)
	stylesExcelize, err := helper.GetReportExcelizeStyles(file)
	streamWriter, err := file.NewStreamWriter(sheets[0])

	if err = helper.SetColWidth(streamWriter, colWidthOpts); err != nil {
		return err
	}
	headerRows, err := helper.SetSheetHeader(streamWriter, reportHeader, styleHeader)
	rowInd := 1
	for _, row := range *data {
		createdAtDate, _ := time.Parse("02.01.2006", row.Timestamp.In(loc).Format("02.01.2006"))

		fullNameDefault := "Анонимный пользователь"
		fullName := fullNameDefault
		if len(row.UserEmail) != 0 {
			fullName = row.UserFullName
		}

		rowInd++

		var fullNameCell any
		if fullName == fullNameDefault {
			fullNameCell = excelize.Cell{Value: fullName, StyleID: stylesExcelize[helper.ReportStringGrayStyle]}
		} else {
			fullNameCell = excelize.Cell{Value: fullName, StyleID: stylesExcelize[helper.ReportStringStyle]}
		}
		var cell string
		cell, err = excelize.CoordinatesToCellName(1, rowInd)
		if err = streamWriter.SetRow(cell,
			[]interface{}{
				excelize.Cell{Value: row.UserEmail, StyleID: stylesExcelize[helper.ReportStringStyle]},
				fullNameCell,
				excelize.Cell{Value: row.Source, StyleID: stylesExcelize[helper.ReportStringStyle]},
				excelize.Cell{Value: row.Shortened, StyleID: stylesExcelize[helper.ReportStringStyle]},
				excelize.Cell{Value: createdAtDate, StyleID: stylesExcelize[helper.ReportDateStyle]},
			},
			excelize.RowOpts{
				Height: helper.ReportDefaultRowSize,
				Hidden: false,
			}); err != nil {
			return err
		}
		if err != nil {
			return err
		}
	}

	if err = helper.AutoFilter(file, sheets[0], headerRows, rowInd); err != nil {
		return err
	}

	if err = streamWriter.Flush(); err != nil {
		return err
	}

	defer func() {
		err = fileInfo.File.SaveAs(fileInfo.FileName)
	}()

	return err
}
