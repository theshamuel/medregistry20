package service

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/theshamuel/medregistry20/app/store/model"
	"github.com/theshamuel/medregistry20/app/utils"
	"log"
	"strings"
	"time"
)

type ReportContractReq struct {
	ClientID  string
	VisitID   string
	DateEvent string
}

type ReportContractServiceRecord struct {
	ID       string
	Name     string
	Quantity int
	Price    int
	Discount float32
	Sum      int
}

type ReportContractServiceTbl struct {
	Records  []*ReportContractServiceRecord
	TotalSum int
}

func (s *DataStore) BuildReportContract(req ReportContractReq) ([]byte, error) {
	f, err := excelize.OpenFile(s.ReportPath + "/templateContract.xlsx")
	if err != nil {
		log.Printf("[ERROR] cannot read template templateContract.xlsx #%v", err)
		return nil, err
	}

	sheetName := f.GetSheetName(1)

	visit, err := s.Engine.FindVisitByID(req.VisitID)
	if err != nil {
		log.Printf("[ERROR] cannot get visit #%v", err)
		return nil, err
	}

	contractDate := visit.DateEvent

	client, err := s.Engine.FindClientByID(req.ClientID)
	if err != nil {
		log.Printf("[ERROR] cannot get client #%v", err)
		return nil, err
	}

	isMiddlename := client.Middlename != ""
	clientFullName := ""
	clientFullNameShort := ""
	if isMiddlename {
		clientFullName = fmt.Sprintf("%s %s %s",
			caser.String(strings.ToTitle(client.Surname)),
			caser.String(strings.ToTitle(client.Firstname)),
			caser.String(strings.ToTitle(client.Middlename)))

		clientFullNameShort = fmt.Sprintf("%s %s. %s.",
			caser.String(strings.ToTitle(client.Surname)),
			client.Firstname[:2],
			client.Middlename[:2])
	} else {
		clientFullName = fmt.Sprintf("%s %s",
			caser.String(strings.ToTitle(client.Surname)),
			caser.String(strings.ToTitle(client.Firstname)))

		clientFullNameShort = fmt.Sprintf("%s %s.",
			caser.String(strings.ToTitle(client.Surname)),
			client.Firstname[:2])
	}

	contractNumber, err := s.Engine.GetSeq("contractNum")
	if err != nil {
		log.Printf("[ERROR] cannot get contractNum #%v", err)
		return nil, err
	}

	//Fill up client info
	// Name
	clientNameCell := f.GetCellValue(sheetName, "R6")
	clientNameCell = strings.ReplaceAll(clientNameCell, "[clientName]", clientFullName)
	f.SetCellStr(sheetName, "R6", clientNameCell)

	clientNameCell = f.GetCellValue(sheetName, "N107")
	clientNameCell = strings.ReplaceAll(clientNameCell, "[clientName]", clientFullName)
	f.SetCellStr(sheetName, "N107", clientNameCell)

	clientNameCell = f.GetCellValue(sheetName, "A11")
	clientNameCell = strings.ReplaceAll(clientNameCell, "[clientName]", clientFullName)
	f.SetCellStr(sheetName, "A11", clientNameCell)

	clientNameCell = f.GetCellValue(sheetName, "J117")
	clientNameCell = strings.ReplaceAll(clientNameCell, "[clientName]", clientFullNameShort)
	f.SetCellStr(sheetName, "J117", clientNameCell)

	clientNameCell = f.GetCellValue(sheetName, "J125")
	clientNameCell = strings.ReplaceAll(clientNameCell, "[clientName]", clientFullNameShort)
	f.SetCellStr(sheetName, "J125", clientNameCell)

	// Birthday
	clientBDCell := f.GetCellValue(sheetName, "L108")
	clientBDCell = strings.ReplaceAll(clientBDCell, "[clientBirthday]", client.Birthday.Format("02.01.2006"))
	f.SetCellStr(sheetName, "L108", clientBDCell)

	// Passport
	clientPassportCell := f.GetCellValue(sheetName, "A11")
	clientPassportCell = strings.ReplaceAll(clientPassportCell, "[passportNumber]", fmt.Sprintf("%s %s", client.PassportSerial, client.PassportNumber))
	passportDate, _ := time.Parse("2006-01-02", client.PassportDate)
	clientPassportCell = strings.ReplaceAll(clientPassportCell, "[passportDate]", passportDate.Format("02.01.2006"))
	clientPassportCell = strings.ReplaceAll(clientPassportCell, "[passportPlace]", client.PassportPlace)

	f.SetCellStr(sheetName, "A11", clientPassportCell)

	// Address
	clientAddressCell := f.GetCellValue(sheetName, "A11")
	clientAddressCell = strings.ReplaceAll(clientAddressCell, "[clientAddress]", client.Address)
	f.SetCellStr(sheetName, "A11", clientAddressCell)

	clientAddressCell = f.GetCellValue(sheetName, "M109")
	clientAddressCell = strings.ReplaceAll(clientAddressCell, "[clientAddress]", client.Address)
	f.SetCellStr(sheetName, "M109", clientAddressCell)

	// Phone
	clientPhoneCell := f.GetCellValue(sheetName, "A11")
	clientPhoneCell = strings.ReplaceAll(clientPhoneCell, "[clientPhone]", utils.GetPhoneValue(client.Phone))
	f.SetCellStr(sheetName, "A11", clientPhoneCell)

	clientPhoneCell = f.GetCellValue(sheetName, "A183")
	clientPhoneCell = strings.ReplaceAll(clientPhoneCell, "[clientPhone]", utils.GetPhoneValue(client.Phone))
	f.SetCellStr(sheetName, "A183", clientPhoneCell)

	clientPhoneCell = f.GetCellValue(sheetName, "N110")
	clientPhoneCell = strings.ReplaceAll(clientPhoneCell, "[clientPhone]", utils.GetPhoneValue(client.Phone))
	f.SetCellStr(sheetName, "N110", clientPhoneCell)

	// Email
	clientEmailCell := f.GetCellValue(sheetName, "P111")
	clientEmailCell = strings.ReplaceAll(clientEmailCell, "[clientEmail]", utils.GetEmailValue(client.Email))
	f.SetCellStr(sheetName, "P111", clientEmailCell)

	clientEmailCell = f.GetCellValue(sheetName, "A183")
	clientEmailCell = strings.ReplaceAll(clientEmailCell, "[clientEmail]", utils.GetEmailValue(client.Email))
	f.SetCellStr(sheetName, "A183", clientEmailCell)

	//Fill up contract number
	contractNumCell := f.GetCellValue(sheetName, "H8")
	contractNumCell = strings.ReplaceAll(contractNumCell, "[contractNum]", fmt.Sprintf("%d", contractNumber))
	f.SetCellStr(sheetName, "H8", contractNumCell)

	contractNumCell = f.GetCellValue(sheetName, "A123")
	contractNumCell = strings.ReplaceAll(contractNumCell, "[contractNum]", fmt.Sprintf("%d", contractNumber))
	f.SetCellStr(sheetName, "A123", contractNumCell)

	//Fill up contract date
	contractDateCell := f.GetCellValue(sheetName, "F6")
	prepDate, _ := time.Parse("2006-01-02", req.DateEvent)
	contractDateCell = strings.ReplaceAll(contractDateCell, "[date]", prepDate.Format("02.01.2006"))
	f.SetCellStr(sheetName, "F6", contractDateCell)

	contractDateCell = f.GetCellValue(sheetName, "A123")
	contractDateCell = strings.ReplaceAll(contractDateCell, "[contractDate]", contractDate.Format("02.01.2006"))
	f.SetCellStr(sheetName, "A123", contractDateCell)

	contractDateCell = f.GetCellValue(sheetName, "A121")
	contractDateCell = strings.ReplaceAll(contractDateCell, "[contractDate]", contractDate.Format("02.01.2006"))
	f.SetCellStr(sheetName, "A121", contractDateCell)

	contractDateCell = f.GetCellValue(sheetName, "A9")
	contractDateCell = strings.ReplaceAll(contractDateCell, "[day]", fmt.Sprintf("%02d", contractDate.Day()))
	f.SetCellStr(sheetName, "A9", contractDateCell)

	contractDateCell = f.GetCellValue(sheetName, "C9")
	contractDateCell = strings.ReplaceAll(contractDateCell, "[month]", utils.GetMonthWordByOrderNumber(contractDate.Month()))
	f.SetCellStr(sheetName, "C9", contractDateCell)

	contractDateCell = f.GetCellValue(sheetName, "F9")
	contractDateCell = strings.ReplaceAll(contractDateCell, "[year]", fmt.Sprintf("%d", contractDate.Year()))
	f.SetCellStr(sheetName, "F9", contractDateCell)

	// Fill up services table and total sum
	serviceTable, err := buildContractServiceTbl(&visit)
	if err != nil {
		log.Printf("[ERROR] cannot buld services table in contract report #%v", err)
		return nil, err
	}
	startIndex := 22

	// For case when services more than 6 we need to unhide necessary count of hidden rows
	for i := 0; i < len(serviceTable.Records)-6; i++ {
		f.SetRowVisible(f.GetSheetName(1), 26+i, true)
	}
	serviceNameCellStyle, _ := f.NewStyle(`{
										 "alignment":
											{
												"horizontal":"left",
												"vertical":"center", 
												"wrap_text":true
											},
										 "font":
											{
												"bold":false, 
												"family":"Times New Roman", 
												"size":10
											}
										}`)
	for i, k := range serviceTable.Records {
		// For case when services more than 5 we need to fill up the bottom line of services table but
		// the index of this row is 61, so I need this condition below to recalculate startIndex
		if i == len(serviceTable.Records)-1 && len(serviceTable.Records) > 5 {
			startIndex = 61 - i
		}
		serviceName := k.Name
		if strings.Contains(serviceName, "[") && strings.Contains(serviceName, "]") {
			serviceName = fmt.Sprintf("%s%s", serviceName[0:strings.Index(serviceName, "[")-1], serviceName[strings.Index(serviceName, "]")+1:])
		}
		f.SetCellStyle(sheetName, fmt.Sprintf("C%d", startIndex+i), fmt.Sprintf("C%d", startIndex+i), serviceNameCellStyle)

		f.SetCellStr(sheetName, fmt.Sprintf("B%d", startIndex+i), fmt.Sprintf("%d", i+1))
		f.SetCellStr(sheetName, fmt.Sprintf("C%d", startIndex+i), serviceName)
		f.SetCellStr(sheetName, fmt.Sprintf("O%d", startIndex+i), fmt.Sprintf("%d", k.Quantity))
		f.SetCellStr(sheetName, fmt.Sprintf("Q%d", startIndex+i), fmt.Sprintf("%d", k.Price))
		f.SetCellStr(sheetName, fmt.Sprintf("S%d", startIndex+i), fmt.Sprintf("%d", k.Sum))
	}

	if len(serviceTable.Records) > 0 {
		totalSumCell := f.GetCellValue(sheetName, "A123")
		totalSumCell = strings.ReplaceAll(totalSumCell, "[totalSum]", fmt.Sprintf("%d", serviceTable.TotalSum))
		f.SetCellStr(sheetName, "A123", totalSumCell)
	}

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

func buildContractServiceTbl(visit *model.Visit) (*ReportContractServiceTbl, error) {
	var res = &ReportContractServiceTbl{}
	res.Records = make([]*ReportContractServiceRecord, 0)

	if len(visit.Services) == 0 {
		return nil, errors.New("failed to build contract services list for the report")
	}

	for _, service := range visit.Services {
		rec := &ReportContractServiceRecord{
			ID:       service.ID,
			Name:     service.Name,
			Quantity: 1,
			Discount: 0,
			Price:    service.Price,
		}

		if service.Discount > 0 {
			rec.Discount = float32(service.Discount) / 100.00
		}

		isServiceAdded := false
		for _, existedService := range res.Records {
			if existedService.ID == rec.ID {
				existedService.Quantity = existedService.Quantity + 1
				isServiceAdded = true
				break
			}
		}

		if !isServiceAdded {
			res.Records = append(res.Records, rec)
		}
	}

	for _, record := range res.Records {
		if record.Discount > 0 {
			record.Sum = int(float32(record.Quantity*record.Price) * (1.00 - record.Discount))
		} else {
			record.Sum = record.Quantity * record.Price
		}
		res.TotalSum = res.TotalSum + record.Sum
	}

	return res, nil
}
