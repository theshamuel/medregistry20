package service

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/pkg/errors"
	"github.com/theshamuel/medregistry20/app/utils"
	"log"
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
	// Not implemented yet
	res, err := utils.ConvertExcelFileToBytes(f)
	if err != nil {
		log.Printf("[ERROR] cannot convert file to bytes")
		return nil, err
	}

	return res, errors.New("not implemented yet")
}
