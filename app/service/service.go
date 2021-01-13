package service

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/theshamuel/medregistry20/app/store"
	"github.com/theshamuel/medregistry20/app/store/model"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type DataStore struct {
	Engine store.EngineInterface
	ReportPath string
}

func (s *DataStore) BuildReportPeriodByDoctorBetweenDateEvent(doctorID string, startDateEvent, endDateEvent string) ([]byte, error) {
	visits, _ := s.Engine.FindVisitByDoctorSinceTill(doctorID, startDateEvent, endDateEvent)
	f, err := excelize.OpenFile(s.ReportPath + "/templateReportOfWorkPeriodByDoctor.xlsx")
	if err != nil {
		log.Printf("[ERROR] Cannot read template templateReportOfWorkPeriodByDoctor.xlsx #%v", err)
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

	res, err := ConvertExcellFileToBytes(f)
	if err != nil {
		log.Printf("[ERROR] Cannot convert file to bytes")
	}

	return res, nil
}

func (s *DataStore) BuildReportVisitResult(visitID string) ([]byte, error) {
	visit, _ := s.Engine.FindVisitById(visitID)
	client, _ := s.Engine.FindClientById(visit.ClientID)
	doctor, _ := s.Engine.FindDoctorById(visit.DoctorID)
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

	res, err := ConvertExcellFileToBytes(f)
	if err != nil {
		log.Printf("[ERROR] Cannot convert file to bytes")
	}
	return res, nil
}

func ConvertExcellFileToBytes(f *excelize.File) ([]byte, error) {
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

	res, _ := ioutil.ReadFile(fileName)
	return res, nil
}


func (s *DataStore) Close() error {
	return s.Engine.Close()
}