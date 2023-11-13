package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/theshamuel/medregistry20/app/utils"
	"log"
	"strings"
	"time"
)

type ReportContractReq struct {
	ClientID  string
	DoctorID  string
	VisitID   string
	DateEvent string
}

func (s *DataStore) BuildReportContract(req ReportContractReq) ([]byte, error) {
	f, err := excelize.OpenFile(s.ReportPath + "/templateContract.xlsx")
	if err != nil {
		log.Printf("[ERROR] cannot read template templateContract.xlsx #%v", err)
		return nil, err
	}

	sheetName := f.GetSheetName(2)

	client, err := s.Engine.FindClientByID(req.ClientID)
	if err != nil {
		log.Printf("[ERROR] cannot get client details #%v", err)
		return nil, err
	}

	clientFullName := fmt.Sprintf("%s %s %s",
		caser.String(strings.ToTitle(client.Surname)),
		caser.String(strings.ToTitle(client.Firstname)),
		caser.String(strings.ToTitle(client.Middlename)))

	clientFullNameShort := fmt.Sprintf("%s %s. %s.",
		caser.String(strings.ToTitle(client.Surname)),
		client.Firstname[:2],
		client.Middlename[:2])

	contractDate, err := time.Parse("2006-01-02", req.DateEvent)
	if err != nil {
		log.Printf("[ERROR] cannot get contract date #%v", err)
		return nil, err
	}
	//org, err := s.Engine.CompanyDetail()
	//if err != nil {
	//	log.Printf("[ERROR] cannot org details #%v", err)
	//	return nil, err
	//}
	//
	//visit, err := s.Engine.FindVisitByID(req.VisitID)
	//if err != nil {
	//	log.Printf("[ERROR] cannot get visit #%v", err)
	//	return nil, err
	//}
	//
	contractNumber, err := s.Engine.GetSeq("contractNum")
	if err != nil {
		log.Printf("[ERROR] cannot get contractNum #%v", err)
		return nil, err
	}

	//Fill up client FIO
	clientNameCell := f.GetCellValue(sheetName, "R6")
	clientNameCell = strings.ReplaceAll(clientNameCell, "[clientName]", clientFullName)
	f.SetCellStr(sheetName, "R6", clientNameCell)

	clientNameCell = f.GetCellValue(sheetName, "L73")
	clientNameCell = strings.ReplaceAll(clientNameCell, "[clientName]", clientFullNameShort)
	f.SetCellStr(sheetName, "L73", clientNameCell)

	//Fill up contract number
	contractNumCell := f.GetCellValue(sheetName, "H8")
	contractNumCell = strings.ReplaceAll(contractNumCell, "[contractNum]", fmt.Sprintf("%d", contractNumber))
	f.SetCellStr(sheetName, "H8", contractNumCell)

	//Fill up contract date
	contractDateCell := f.GetCellValue(sheetName, "F6")
	contractDateCell = strings.ReplaceAll(contractDateCell, "[date]", time.Now().Format("02.01.2006"))
	f.SetCellStr(sheetName, "F6", contractDateCell)

	contractDateCell = f.GetCellValue(sheetName, "A9")
	contractDateCell = strings.ReplaceAll(contractDateCell, "[day]", fmt.Sprintf("%02d", contractDate.Day()))
	f.SetCellStr(sheetName, "A9", contractDateCell)

	contractDateCell = f.GetCellValue(sheetName, "C9")
	contractDateCell = strings.ReplaceAll(contractDateCell, "[month]", utils.GetMonthWordByOrderNumber(contractDate.Month()))
	f.SetCellStr(sheetName, "C9", contractDateCell)

	contractDateCell = f.GetCellValue(sheetName, "F9")
	contractDateCell = strings.ReplaceAll(contractDateCell, "[year]", fmt.Sprintf("%d", contractDate.Year()))
	f.SetCellStr(sheetName, "F9", contractDateCell)

	res, err := utils.ConvertExcelFileToBytes(f)
	if err != nil {
		log.Printf("[ERROR] cannot convert file to bytes")
		return nil, err
	}

	err = s.Engine.IncrementSeq(contractNumber+1, "contractNum")
	if err != nil {
		log.Printf("[ERROR] cannot increment number of sequence for nalog spravka report")
	}

	return res, nil
}
