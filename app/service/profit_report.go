package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/theshamuel/medregistry20/app/utils"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"log"
	"strings"
)

type ProfitReportRecord struct {
	Name string
	Sum  int
}

type ProfitReportTbl struct {
	Records  []*ProfitReportRecord
	TotalSum int
}

func (s *DataStore) BuildReportOfPeriodProfit(startDateEvent, endDateEvent string) ([]byte, error) {
	f, err := excelize.OpenFile(s.ReportPath + "/templateProfitReport.xlsx")
	if err != nil {
		log.Printf("[ERROR] cannot read template templateProfitReport.xlsx #%v", err)
		return nil, err
	}

	sheetName := f.GetSheetName(1)

	//Fill up title
	titleCell := f.GetCellValue(sheetName, "A2")
	titleCell = strings.ReplaceAll(titleCell, "[startDate]", startDateEvent)
	titleCell = strings.ReplaceAll(titleCell, "[endDate]", endDateEvent)
	f.SetCellStr(sheetName, "A2", titleCell)

	doctorsData, err := s.buildDoctorProfitData(startDateEvent, endDateEvent)
	if err != nil {
		log.Printf("[ERROR] cannot buld doctor profit table in profit report #%v", err)
		return nil, err
	}
	var cellCode string
	var cellValue string
	for i, record := range doctorsData.Records {
		f.SetRowVisible(f.GetSheetName(1), 35-i, true)
		//Fill up doctor name
		cellCode = fmt.Sprintf("%s%d", "A", 35-i)
		cellValue = record.Name
		f.SetCellStr(sheetName, cellCode, cellValue)

		//Fill up doctor's total
		cellCode = fmt.Sprintf("%s%d", "B", 35-i)
		f.SetCellInt(sheetName, cellCode, record.Sum)
	}

	//Fill up total sum
	cellCode = "B36"
	f.SetCellInt(sheetName, cellCode, doctorsData.TotalSum)

	res, err := utils.ConvertExcelFileToBytes(f)
	if err != nil {
		log.Printf("[ERROR] cannot convert file to bytes")
		return nil, err
	}

	return res, nil
}

func (s *DataStore) buildDoctorProfitData(fromDate, toDate string) (*ProfitReportTbl, error) {
	res := &ProfitReportTbl{}
	res.Records = make([]*ProfitReportRecord, 0)
	doctors, err := s.Engine.FindDoctors()
	if err != nil {
		return nil, err
	}
	caser := cases.Title(language.Russian)

	for _, doctor := range doctors {
		if doctor.Surname == "Лаборатория" {
			continue
		}
		if doctor.Surname == "Здоровенко" {
			doctor.Middlename = doctor.Middlename + " (лаборатория)"
		}
		visits, err := s.Engine.FindVisitsByDoctorSinceTill(doctor.ID, fromDate, toDate)
		if err != nil {
			return nil, err
		}
		var total int
		for _, visit := range visits {
			total += visit.TotalSum
		}

		res.Records = append(res.Records, &ProfitReportRecord{
			Name: fmt.Sprintf("%s %s %s", caser.String(strings.ToLower(doctor.Surname)), caser.String(strings.ToLower(doctor.FirstName)), caser.String(strings.ToLower(doctor.Middlename))),
			Sum:  total,
		})
		res.TotalSum += total
	}

	return res, nil
}
