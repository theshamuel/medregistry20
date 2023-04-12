package service

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/theshamuel/medregistry20/app/store"
	"github.com/theshamuel/medregistry20/app/store/model"
	"github.com/theshamuel/medregistry20/app/utils"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type DataStore struct {
	Engine     store.EngineInterface
	ReportPath string
}

type ReportNalogSpravkaReq struct {
	ClientID              string
	DateFrom              string
	DateTo                string
	PayerFIO              string
	RelationClientToPayer string
	GenderOfPayer         string
	IsClientSelfPayer     bool
}

func (s *DataStore) BuildReportPeriodByDoctorBetweenDateEvent(doctorID string, startDateEvent, endDateEvent string) ([]byte, error) {
	visits, _ := s.Engine.FindVisitsByDoctorSinceTill(doctorID, startDateEvent, endDateEvent)
	f, err := excelize.OpenFile(s.ReportPath + "/templateReportOfWorkPeriodByDoctor.xlsx")
	if err != nil {
		log.Printf("[ERROR] Cannot read template templateReportOfWorkPeriodByDoctor.xlsx #%v", err)
		return nil, err
	}
	data := model.ProcessDataDoctorSalaryRecord(visits)
	numRecordInt := 6
	sheetName := f.GetSheetName(1)
	doctorNameCellStyle, err := f.NewStyle(`{"alignment":{"horizontal":"center", "vertical":"center"},
												"font":{"bold":true,"family":"Times New Roman","size":14,"color":"#000000"},
												"border":[{"type":"left","color":"000000","style":2},{"type":"top","color":"000000","style":2},{"type":"bottom","color":"000000","style":2},{"type":"right","color":"0000000","style":2}]}`)
	if err != nil {
		log.Printf("[ERROR] cannot create cell style for doctor name %#v", err)
	}
	serviceCellStyle, err := f.NewStyle(`{"alignment":{"horizontal":"left", "vertical":"left"},
												"font":{"bold":false,"family":"Times New Roman","size":12,"color":"#000000"},
												"border":[{"type":"left","color":"000000","style":1},{"type":"top","color":"000000","style":1},{"type":"bottom","color":"000000","style":1},{"type":"right","color":"0000000","style":1}]}`)
	if err != nil {
		log.Printf("[ERROR] cannot create cell style for service %#v", err)
	}
	totalSumCellStyle, err := f.NewStyle(`{"alignment":{"horizontal":"right", "vertical":"center"},
												"font":{"bold":true,"family":"Times New Roman","size":14,"color":"#000000"}}`)
	if err != nil {
		log.Printf("[ERROR] cannot create cell style for total sum %#v", err)
	}
	numberTitle := f.GetCellValue(sheetName, "B1")
	numberTitle = strings.ReplaceAll(numberTitle, "[number]", time.Now().Format("06/0201"))
	f.SetCellStr(sheetName, "B1", numberTitle)
	periodTitle := f.GetCellValue(sheetName, "B2")
	periodTitle = strings.ReplaceAll(periodTitle, "[startDate]", startDateEvent)
	periodTitle = strings.ReplaceAll(periodTitle, "[endDate]", endDateEvent)
	f.SetCellStr(sheetName, "B2", periodTitle)
	for _, valDsr := range data {
		numRecordSrt := strconv.Itoa(numRecordInt)
		f.MergeCell(sheetName, "A"+numRecordSrt, "E"+numRecordSrt)
		f.SetCellStr(sheetName, "A"+numRecordSrt, valDsr.DoctorName)
		f.SetCellStyle(sheetName, "A"+numRecordSrt, "E"+numRecordSrt, doctorNameCellStyle)
		numRecordInt = numRecordInt + 1
		for _, service := range valDsr.Services {
			numRecordSrt = strconv.Itoa(numRecordInt)
			f.MergeCell(sheetName, "A"+numRecordSrt, "B"+numRecordSrt)
			f.SetCellStr(sheetName, "A"+numRecordSrt, service.ServiceName)
			f.SetCellInt(sheetName, "C"+numRecordSrt, service.Count)
			f.SetCellInt(sheetName, "D"+numRecordSrt, service.DoctorRate)
			f.SetCellFormula(sheetName, "E"+numRecordSrt, "C"+numRecordSrt+"*"+"D"+numRecordSrt)
			f.SetCellStyle(sheetName, "A"+numRecordSrt, "E"+numRecordSrt, serviceCellStyle)
			numRecordInt = numRecordInt + 1
		}
		numRecordSrt = strconv.Itoa(numRecordInt)
		f.MergeCell(sheetName, "A"+numRecordSrt, "D"+numRecordSrt)
		f.SetCellStr(sheetName, "D"+numRecordSrt, "ИТОГО")
		f.SetCellFormula(sheetName, "E"+numRecordSrt, "SUM(E"+strconv.Itoa(numRecordInt-len(valDsr.Services))+":E"+strconv.Itoa(numRecordInt-1)+")")
		f.SetCellStyle(sheetName, "A"+numRecordSrt, "E"+numRecordSrt, totalSumCellStyle)
		numRecordInt = numRecordInt + 1
	}

	if err != nil {
		log.Printf("[ERROR] can not save report templateReportOfWorkPeriodByDoctor.xlsx")
	}

	res, err := ConvertExcelFileToBytes(f)
	if err != nil {
		log.Printf("[ERROR] Cannot convert file to bytes")
	}
	//В том, что он (она) оплатил(а) медицинские услуги стоимостью
	return res, nil
}

func (s *DataStore) BuildReportVisitResult(visitID string) ([]byte, error) {
	visit, _ := s.Engine.FindVisitByID(visitID)
	client, _ := s.Engine.FindClientByID(visit.ClientID)
	doctor, _ := s.Engine.FindDoctorByID(visit.DoctorID)
	company, _ := s.Engine.CompanyDetail()

	f, err := excelize.OpenFile(s.ReportPath + "/templateVisitResult.xlsx")
	if err != nil {
		log.Printf("[ERROR] Cannot read template templateVisitResult.xlsx #%v", err)
		return nil, err
	}
	sheetName := f.GetSheetName(1)

	//Fill up license cell [Left side]
	licenseCell := f.GetCellValue(sheetName, "E3")
	licenseCell = strings.ReplaceAll(licenseCell, "[orgLicenceNumber]", company.License)
	f.SetCellStr(sheetName, "E3", licenseCell)

	//Fill up license cell [Right side]
	licenseCell = f.GetCellValue(sheetName, "N3")
	licenseCell = strings.ReplaceAll(licenseCell, "[orgLicenceNumber]", company.License)
	f.SetCellStr(sheetName, "N3", licenseCell)

	//Fill up doctor type cell [Left side]
	doctorTypeCell := f.GetCellValue(sheetName, "A6")
	doctorTypeCell = strings.ReplaceAll(doctorTypeCell, "[doctorType]", strings.ToUpper(doctor.PositionGenitive))
	f.SetCellStr(sheetName, "A6", doctorTypeCell)

	//Fill up doctor type cell [Right side]
	doctorTypeCell = f.GetCellValue(sheetName, "J6")
	doctorTypeCell = strings.ReplaceAll(doctorTypeCell, "[doctorType]", strings.ToUpper(doctor.PositionGenitive))
	f.SetCellStr(sheetName, "J6", doctorTypeCell)

	//Fill up client full name cell [Left side]
	fioClientCell := f.GetCellValue(sheetName, "A8")
	fioClientCell = strings.ReplaceAll(fioClientCell, "[fioClient]",
		strings.Title(strings.ToLower(client.Surname))+" "+
			strings.Title(strings.ToLower(client.Firstname))+" "+
			strings.Title(strings.ToLower(client.Middlename)))
	f.SetCellStr(sheetName, "A8", fioClientCell)

	//Fill up client full name cell [Right side]
	fioClientCell = f.GetCellValue(sheetName, "J8")
	fioClientCell = strings.ReplaceAll(fioClientCell, "[fioClient]",
		strings.Title(strings.ToLower(client.Surname))+" "+
			strings.Title(strings.ToLower(client.Firstname))+" "+
			strings.Title(strings.ToLower(client.Middlename)))
	f.SetCellStr(sheetName, "J8", fioClientCell)

	//Fill up client birth date cell [Left side]
	birthDateClientCell := f.GetCellValue(sheetName, "A9")
	birthDateClientCell = strings.ReplaceAll(birthDateClientCell, "[birthDateClient]",
		client.Birthday.Format("02.01.2006"))
	f.SetCellStr(sheetName, "A9", birthDateClientCell)

	//Fill up client birth date cell [Right side]
	birthDateClientCell = f.GetCellValue(sheetName, "J9")
	birthDateClientCell = strings.ReplaceAll(birthDateClientCell, "[birthDateClient]",
		client.Birthday.Format("02.01.2006"))
	f.SetCellStr(sheetName, "J9", birthDateClientCell)

	//Fill up client age cell [Left side]
	ageClientCell := f.GetCellValue(sheetName, "A9")
	ageClientCell = strings.ReplaceAll(ageClientCell, "[ageClient]", "полных "+
		strconv.Itoa(client.Age)+" "+client.AgePostfix)
	f.SetCellStr(sheetName, "A9", ageClientCell)

	//Fill up client age cell [Right side]
	ageClientCell = f.GetCellValue(sheetName, "J9")
	ageClientCell = strings.ReplaceAll(ageClientCell, "[ageClient]", "полных "+
		strconv.Itoa(client.Age)+" "+client.AgePostfix)
	f.SetCellStr(sheetName, "J9", ageClientCell)

	//Fill up client date of visit cell [Left side]
	dateEventCell := f.GetCellValue(sheetName, "A10")
	dateEventCell = strings.ReplaceAll(dateEventCell, "[dateEvent]", visit.DateEvent.Format("02.01.2006"))
	f.SetCellStr(sheetName, "A10", dateEventCell)

	//Fill up client date of visit cell [Right side]
	dateEventCell = f.GetCellValue(sheetName, "J10")
	dateEventCell = strings.ReplaceAll(dateEventCell, "[dateEvent]", visit.DateEvent.Format("02.01.2006"))
	f.SetCellStr(sheetName, "J10", dateEventCell)

	//Fill up client diagnosis cell [Left side]
	diagnosisCell := f.GetCellValue(sheetName, "A12")
	diagnosisCell = strings.ReplaceAll(diagnosisCell, "[diagnosis]", visit.Diagnosis)
	f.SetCellStr(sheetName, "A12", diagnosisCell)

	//Fill up client diagnosis cell [Right side]
	diagnosisCell = f.GetCellValue(sheetName, "J12")
	diagnosisCell = strings.ReplaceAll(diagnosisCell, "[diagnosis]", visit.Diagnosis)
	f.SetCellStr(sheetName, "J12", diagnosisCell)

	//Fill up client additional examination cell [Left side]
	additionalExaminationCell := f.GetCellValue(sheetName, "A19")
	additionalExaminationCell = strings.ReplaceAll(additionalExaminationCell, "[additionalExamination]",
		visit.AdditionalExamination)
	f.SetCellStr(sheetName, "A19", additionalExaminationCell)

	//Fill up client additional examination cell [Right side]
	additionalExaminationCell = f.GetCellValue(sheetName, "J19")
	additionalExaminationCell = strings.ReplaceAll(additionalExaminationCell, "[additionalExamination]",
		visit.AdditionalExamination)
	f.SetCellStr(sheetName, "J19", additionalExaminationCell)

	//Fill up client therapy cell [Left side]
	therapyCell := f.GetCellValue(sheetName, "A25")
	therapyCell = strings.ReplaceAll(therapyCell, "[therapy]", visit.Therapy)
	f.SetCellStr(sheetName, "A25", therapyCell)

	//Fill up client therapy cell [Right side]
	therapyCell = f.GetCellValue(sheetName, "J25")
	therapyCell = strings.ReplaceAll(therapyCell, "[therapy]", visit.Therapy)
	f.SetCellStr(sheetName, "J25", therapyCell)

	//Fill up client therapy cell [Left side]
	fioDoctorCell := f.GetCellValue(sheetName, "A40")
	fioDoctorCell = strings.ReplaceAll(fioDoctorCell, "[fioDoctor]",
		strings.Title(strings.ToLower(doctor.Surname))+" "+
			strings.Title(strings.ToLower(doctor.FirstName))+" "+
			strings.Title(strings.ToLower(doctor.Middlename)))
	f.SetCellStr(sheetName, "A40", fioDoctorCell)

	//Fill up client therapy cell [Right side]
	fioDoctorCell = f.GetCellValue(sheetName, "J40")
	fioDoctorCell = strings.ReplaceAll(fioDoctorCell, "[fioDoctor]",
		strings.Title(strings.ToLower(doctor.Surname))+" "+
			strings.Title(strings.ToLower(doctor.FirstName))+" "+
			strings.Title(strings.ToLower(doctor.Middlename)))
	f.SetCellStr(sheetName, "J40", fioDoctorCell)

	res, err := ConvertExcelFileToBytes(f)
	if err != nil {
		log.Printf("[ERROR] cannot convert file to bytes")
	}
	return res, nil
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
		res, _ := ConvertExcelFileToBytes(f)
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
		payerName = req.PayerFIO
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
	genderCellNormalStyle, _ := f.NewStyle(`{"alignment":{"horizontal":"left", "vertical":"left"},
										 "font":{"bold":true, "underline": "single", "family":"Times New Roman", "size":10, "color":"#000000" }
										}`)
	familyRelationCellNormalStyle, _ := f.NewStyle(`{"alignment":{"horizontal":"left", "vertical":"left"},
										 "font":{"bold":true, "underline": "single", "family":"Times New Roman", "size":10, "color":"#000000" }
										}`)
	familyRelationCellRedlStyle, _ := f.NewStyle(`{"alignment":{"horizontal":"left", "vertical":"left"},
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
		if visits[0].ClientGender == "woman" {
			f.SetCellStyle(sheetName, "D42", "D42", genderCellNormalStyle)
			f.SetCellStyle(sheetName, "D47", "D47", familyRelationCellNormalStyle)
		} else {
			f.SetCellStyle(sheetName, "C42", "C42", genderCellNormalStyle)
			f.SetCellStyle(sheetName, "C47", "C47", familyRelationCellNormalStyle)
		}

		switch req.RelationClientToPayer {
		case "spouse":
			f.SetCellStyle(sheetName, "E47", "E47", familyRelationCellNormalStyle)
			break
		case "son":
			f.SetCellStyle(sheetName, "F47", "F47", familyRelationCellNormalStyle)
			break
		case "daughter":
			f.SetCellStyle(sheetName, "G47", "G47", familyRelationCellNormalStyle)
			break
		case "mother":
			f.SetCellStyle(sheetName, "H47", "H47", familyRelationCellNormalStyle)
			break
		case "father":
			f.SetCellStyle(sheetName, "I47", "I47", familyRelationCellNormalStyle)
			break
		default:
			f.SetCellStyle(sheetName, "C47", "I47", familyRelationCellRedlStyle)
			break
		}
	}

	res, err := ConvertExcelFileToBytes(f)
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

func ConvertExcelFileToBytes(f *excelize.File) ([]byte, error) {
	fileName := os.TempDir() + "/" + strconv.FormatInt(time.Now().Unix(), 10) + ".xlsx"
	err := f.SaveAs(fileName)
	if err != nil {
		log.Printf("[ERROR] Cannot create temporary file #%v", err)
		return nil, err
	}
	defer func() {
		err := os.Remove(fileName)
		if err != nil {
			log.Printf("[WARN] failed to remove %s from FS: %s", fileName, err)
		}
	}()

	res, _ := os.ReadFile(fileName)
	return res, nil
}

func (s *DataStore) Close() error {
	return s.Engine.Close()
}
