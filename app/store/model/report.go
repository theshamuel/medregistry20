package model

import (
	"errors"
	"log"
	"strings"
)

type DoctorSalaryRecord struct {
	ID         string
	DoctorName string
	Services   map[string]ServiceReport
}

type ServiceReport struct {
	ServiceName string
	Count       int
	DoctorRate  int
}

func ProcessDataDoctorSalaryRecord(visits []Visit) map[string]*DoctorSalaryRecord {
	var res = make(map[string]*DoctorSalaryRecord)
	for _, visit := range visits {
		if !visit.DoctorExcludedFromReports {
			//Because of in total sum for visit pcr could be included once we're using this flag
			addedPCRPrice := false
			dsr := res[visit.DoctorID]
			if dsr == nil {
				dsr = &DoctorSalaryRecord{ID: visit.DoctorID, DoctorName: visit.DoctorName, Services: make(map[string]ServiceReport)}
			}
			for _, service := range visit.Services {
				if addedPCRPrice && service.Category == "pcr" || service.Category == "mazok"{
					continue
				}
				//Cut off a tail of ID. ID=_id+MEDREG+Random(int). It is necessary for available in visit duplication of services in grid.
				service.ID = strings.Split(service.ID, "MEDREG")[0]
				if _, ok := dsr.Services[service.ID]; !ok {
					//Check if this service should be paid (exclude analizys)
					if service.DoctorPay > 0  {
						personalRate, err := calcPersonalRate(service, dsr.ID)
						if err != nil {
							log.Printf("[DEBUG] not personal rate for %s; %s", dsr.ID, err.Error())
						}
						srv := ServiceReport{ServiceName: service.Name, Count: 1, DoctorRate: service.DoctorPay}
						if personalRate != nil {
							srv.DoctorRate = personalRate.DoctorSalary
						}
						if service.Category == "pcr" && addedPCRPrice {
							srv.DoctorRate = 0
						}
						dsr.Services[service.ID] = srv
					}
				} else {
					incCountService(dsr, service.ID)
				}
			}
			//Check that is doctor need by paid for services
			if len(dsr.Services) > 0 {
				res[dsr.ID] = dsr
			}
		}
	}
	return res
}

func calcPersonalRate(s Service, doctorID string) (*PersonalRate, error) {
	if s.PersonalRates != nil || len(s.PersonalRates) == 0 {
		return nil, errors.New("can't found personal rate")
	}
	for _, pr := range s.PersonalRates {
		if pr.DoctorID == doctorID {
			if pr.DoctorPayType == "" || pr.DoctorPayType == "sum" {
				pr.DoctorSalary = pr.DoctorPay
			} else {
				pr.DoctorSalary = s.Price * pr.DoctorPay / 100
			}
			return &pr, nil
		}
	}
	return nil, errors.New("can't found personal rate")
}

//Increment count of equals services through all doctor's visits
func incCountService(doctorSalaryRecord *DoctorSalaryRecord, serviceID string){
	tmp := doctorSalaryRecord.Services[serviceID]
	tmp.Count = tmp.Count + 1
	doctorSalaryRecord.Services[serviceID] = tmp
}