package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/theshamuel/medregistry20/app/utils"
	"log"
	"strconv"
	"strings"
	"time"
)

type ReportNalogSpravkaReq struct {
	ClientID              string
	DateFrom              string
	DateTo                string
	PayerFIO              string
	RelationClientToPayer string
	GenderOfPayer         string
	IsClientSelfPayer     bool
}

func (s *DataStore) BuildReportNalogSpravka(req ReportNalogSpravkaReq) ([]byte, error) {
	f, err := excelize.OpenFile(s.ReportPath + "/templateNalogSpravka.xlsx")
	if err != nil {
		log.Printf("[ERROR] cannot read template templateNalogSpravka.xlsx #%v", err)
		return nil, err
	}
	sheetName := f.GetSheetName(1)

	visits, err := s.Engine.FindVisitsByClientIDSinceTill(req.ClientID, req.DateFrom, req.DateTo)
	if err != nil {
		log.Printf("[ERROR] cannot get client visits #%v", err)
		return nil, err
	}

	commonCellRedStyle, _ := f.NewStyle(`{"alignment":{"horizontal":"left", "vertical":"center"},
										 "font":{"bold":true, "underline": "single", "family":"Times New Roman", "size":10, "color":"#FF0000" }
										}`)

	if len(visits) == 0 {
		visitDatesCell := f.GetCellValue(sheetName, "D19")
		visitDatesCell = strings.ReplaceAll(visitDatesCell, "[visitDates]", "ВИЗИТЫ ЗА УКАЗАННЫЙ ПЕРИОД ОТСУТСТВУЮТ")
		f.SetCellStr(sheetName, "D19", visitDatesCell)

		visitDatesCell = f.GetCellValue(sheetName, "C53")
		visitDatesCell = strings.ReplaceAll(visitDatesCell, "[visitDates]", "ВИЗИТЫ ЗА УКАЗАННЫЙ ПЕРИОД ОТСУТСТВУЮТ")
		f.SetCellStr(sheetName, "C53", visitDatesCell)

		f.SetCellStyle(sheetName, "D19", "D19", commonCellRedStyle)
		f.SetCellStyle(sheetName, "C53", "C53", commonCellRedStyle)
		res, _ := utils.ConvertExcelFileToBytes(f)
		return res, nil
	}

	superTotalSum := 0.0
	superTotalSumStr := ""
	var visitDatesStr strings.Builder
	for i, v := range visits {
		superTotalSum = superTotalSum + v.TotalSumWithPenny
		if i < len(visits)-1 {
			visitDatesStr.WriteString(v.DateEvent.Format("02.01.2006"))
			visitDatesStr.WriteString(", ")
		} else {
			visitDatesStr.WriteString(v.DateEvent.Format("02.01.2006"))
		}
	}
	if (int(superTotalSum*100) % 100) == 0 {
		superTotalSumStr = fmt.Sprintf("%d", int(superTotalSum))
	} else {
		superTotalSumStr = fmt.Sprintf("%.2f", superTotalSum)
	}

	numberOfNalogSpravka, err := s.Engine.GetNalogSpravkaSeq()
	if err != nil {
		log.Printf("[ERROR] cannot get numberOfNalogSpravka #%v", err)
		return nil, err
	}

	company, _ := s.Engine.CompanyDetail()

	payerName := visits[0].ClientSurname + " " + visits[0].ClientName + " " + visits[0].ClientMiddlename
	if !req.IsClientSelfPayer {
		payerName = caser.String(strings.ToLower(req.PayerFIO))
	}

	//Fill up payer FIO
	payerFIOCell := f.GetCellValue(sheetName, "F13")
	payerFIOCell = strings.ReplaceAll(payerFIOCell, "[payerFio]", payerName)
	f.SetCellStr(sheetName, "F13", payerFIOCell)

	payerFIOCell = f.GetCellValue(sheetName, "F40")
	payerFIOCell = strings.ReplaceAll(payerFIOCell, "[payerFio]", payerName)
	f.SetCellStr(sheetName, "F40", payerFIOCell)

	//Fill up license cell
	licenseCell := f.GetCellValue(sheetName, "A30")
	licenseCell = strings.ReplaceAll(licenseCell, "[orgLicenceNumber]", company.License)
	f.SetCellStr(sheetName, "A30", licenseCell)

	//Fill up number
	numberCell := f.GetCellValue(sheetName, "B9")
	numberCell = strings.ReplaceAll(numberCell, "[number]", strconv.Itoa(numberOfNalogSpravka))
	f.SetCellStr(sheetName, "B9", numberCell)

	numberCell = f.GetCellValue(sheetName, "B37")
	numberCell = strings.ReplaceAll(numberCell, "[number]", strconv.Itoa(numberOfNalogSpravka))
	f.SetCellStr(sheetName, "B37", numberCell)

	//Fill up total sum
	totalSumCell := f.GetCellValue(sheetName, "H17")
	totalSumCell = strings.ReplaceAll(totalSumCell, "[totalSum]", superTotalSumStr)
	f.SetCellStr(sheetName, "H17", totalSumCell)

	totalSumAndCurrencyCell := f.GetCellValue(sheetName, "J42")
	totalSumAndCurrencyCell = strings.ReplaceAll(totalSumAndCurrencyCell, "[totalSumAndCurrency]", superTotalSumStr+" "+utils.GetCurrencySuffix(int(superTotalSum)%10))
	f.SetCellStr(sheetName, "J42", totalSumAndCurrencyCell)

	//Fill up visit dates
	visitDatesCell := f.GetCellValue(sheetName, "D19")
	visitDatesCell = strings.ReplaceAll(visitDatesCell, "[visitDates]", visitDatesStr.String())
	f.SetCellStr(sheetName, "D19", visitDatesCell)

	visitDatesCell = f.GetCellValue(sheetName, "C53")
	visitDatesCell = strings.ReplaceAll(visitDatesCell, "[visitDates]", visitDatesStr.String())
	f.SetCellStr(sheetName, "C53", visitDatesCell)

	//Fill up release date
	today := time.Now()
	releaseDateCell := f.GetCellValue(sheetName, "FC21")
	releaseDateCell = strings.ReplaceAll(releaseDateCell, "[releaseDate]", today.Format("01.02.2006"))
	f.SetCellStr(sheetName, "F21", releaseDateCell)

	releaseDateWithMonthWordCell := f.GetCellValue(sheetName, "C38")
	releaseDateWithMonthWordCell = strings.ReplaceAll(releaseDateWithMonthWordCell, "[dayOfReleaseDate]", strconv.Itoa(today.Day()))
	releaseDateWithMonthWordCell = strings.ReplaceAll(releaseDateWithMonthWordCell, "[monthOfReleaseDateWord]", utils.GetMonthWordByOrderNumber(today.Month()))
	releaseDateWithMonthWordCell = strings.ReplaceAll(releaseDateWithMonthWordCell, "[yearOfReleaseDate]", strconv.Itoa(today.Year()))
	f.SetCellStr(sheetName, "C38", releaseDateWithMonthWordCell)

	//Fill up gender
	genderCellNormalStyle, _ := f.NewStyle(`{"alignment":{"horizontal":"left", "vertical":"center"},
										 "font":{"bold":true, "underline": "single", "family":"Times New Roman", "size":10, "color":"#000000" }
										}`)
	familyRelationCellNormalStyle, _ := f.NewStyle(`{"alignment":{"horizontal":"left", "vertical":"center"},
										 "font":{"bold":true, "underline": "single", "family":"Times New Roman", "size":10, "color":"#000000" }
										}`)
	familyRelationCellRedlStyle, _ := f.NewStyle(`{"alignment":{"horizontal":"left", "vertical":"center"},
										 "font":{"bold":true, "underline": "single", "family":"Times New Roman", "size":10, "color":"#FF0000" }
										}`)

	if req.IsClientSelfPayer {
		if visits[0].ClientGender == "woman" {
			f.SetCellStyle(sheetName, "D42", "D42", genderCellNormalStyle)
			f.SetCellStyle(sheetName, "D47", "D47", familyRelationCellNormalStyle)
		} else {
			f.SetCellStyle(sheetName, "C42", "C42", genderCellNormalStyle)
			f.SetCellStyle(sheetName, "C47", "C47", familyRelationCellNormalStyle)
		}
	} else {
		if req.GenderOfPayer == "female" {
			f.SetCellStyle(sheetName, "D42", "D42", genderCellNormalStyle)
		} else {
			f.SetCellStyle(sheetName, "C42", "C42", genderCellNormalStyle)
		}

		switch req.RelationClientToPayer {
		case "spouse":
			f.SetCellStyle(sheetName, "E47", "E47", familyRelationCellNormalStyle)
		case "son":
			f.SetCellStyle(sheetName, "F47", "F47", familyRelationCellNormalStyle)
		case "daughter":
			f.SetCellStyle(sheetName, "G47", "G47", familyRelationCellNormalStyle)
		case "mother":
			f.SetCellStyle(sheetName, "H47", "H47", familyRelationCellNormalStyle)
		case "father":
			f.SetCellStyle(sheetName, "I47", "I47", familyRelationCellNormalStyle)
		default:
			f.SetCellStyle(sheetName, "C47", "I47", familyRelationCellRedlStyle)
		}
	}

	res, err := utils.ConvertExcelFileToBytes(f)
	if err != nil {
		log.Printf("[ERROR] cannot convert file to bytes")
		return nil, err
	}

	err = s.Engine.IncrementNalogSpravkaSeq(numberOfNalogSpravka + 1)
	if err != nil {
		log.Printf("[ERROR] cannot increment number of sequence for nalog spravka report")
	}
	return res, nil
}
