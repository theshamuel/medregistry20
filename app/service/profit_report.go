package service

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/theshamuel/medregistry20/app/utils"
	"log"
	"strings"
)

func (s *DataStore) BuildReportOfPeriodProfit(startDateEvent, endDateEvent string) ([]byte, error) {
	f, err := excelize.OpenFile(s.ReportPath + "/templateProfitReport.xlsx")
	if err != nil {
		log.Printf("[ERROR] cannot read template templateProfitReport.xlsx #%v", err)
		return nil, err
	}

	sheetName := f.GetSheetName(1)

	doctorName := "dummy"

	//Fill up doctor name
	clientNameCell := f.GetCellValue(sheetName, "XX")
	clientNameCell = strings.ReplaceAll(clientNameCell, "[doctorName]", doctorName)
	f.SetCellStr(sheetName, "XX", clientNameCell)

	res, err := utils.ConvertExcelFileToBytes(f)
	if err != nil {
		log.Printf("[ERROR] cannot convert file to bytes")
		return nil, err
	}

	return res, nil
}
