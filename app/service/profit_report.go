package service

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/theshamuel/medregistry20/app/utils"
	"log"
	"strings"
)

type ProfitReportRecord struct {
	Name string
	Sum  int
}

type ProfitReportTbl struct {
	Records  []*ReportContractServiceRecord
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

	doctors, err := s.Engine.FindDoctors()
	if err != nil {
		log.Printf("[ERROR] cannot retrieve doctors from API server v1 #%v", err)
		return nil, err
	}

	for _, doctor := range doctors {
		log.Println(doctor.Surname)
	}
	//recordFileIndex := 5
	doctorsData, err := buildDoctorProfitData()
	if err != nil {
		log.Printf("[ERROR] cannot buld doctor profit table in profit report #%v", err)
		return nil, err
	}

	for i := 0; i < len(doctorsData.Records)-1; i++ {
		f.SetRowVisible(f.GetSheetName(1), 26+i, true)
	}

	res, err := utils.ConvertExcelFileToBytes(f)
	if err != nil {
		log.Printf("[ERROR] cannot convert file to bytes")
		return nil, err
	}

	return res, nil
}

func buildDoctorProfitData() (*ProfitReportTbl, error) {
	res := &ProfitReportTbl{}
	res.Records = make([]*ReportContractServiceRecord, 0)

	return res, nil
}
