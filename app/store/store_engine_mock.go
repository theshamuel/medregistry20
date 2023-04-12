// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package store

import (
	"github.com/theshamuel/medregistry20/app/store/model"
	"sync"
)

// Ensure, that EngineInterfaceMock does implement EngineInterface.
// If this is not the case, regenerate this file with moq.
var _ EngineInterface = &EngineInterfaceMock{}

// EngineInterfaceMock is a mock implementation of EngineInterface.
//
//	func TestSomethingThatUsesEngineInterface(t *testing.T) {
//
//		// make and configure a mocked EngineInterface
//		mockedEngineInterface := &EngineInterfaceMock{
//			CloseFunc: func() error {
//				panic("mock out the Close method")
//			},
//			CompanyDetailFunc: func() (model.Company, error) {
//				panic("mock out the CompanyDetail method")
//			},
//			FindClientByIDFunc: func(id string) (model.Client, error) {
//				panic("mock out the FindClientByID method")
//			},
//			FindDoctorByIDFunc: func(id string) (model.Doctor, error) {
//				panic("mock out the FindDoctorByID method")
//			},
//			FindVisitByIDFunc: func(id string) (model.Visit, error) {
//				panic("mock out the FindVisitByID method")
//			},
//			FindVisitsByClientIDSinceTillFunc: func(clientID string, startDateEventStr string, endDateEventStr string) ([]model.Visit, error) {
//				panic("mock out the FindVisitsByClientIDSinceTill method")
//			},
//			FindVisitsByDoctorSinceTillFunc: func(doctorID string, startDateEvent string, endDateEvent string) ([]model.Visit, error) {
//				panic("mock out the FindVisitsByDoctorSinceTill method")
//			},
//			GetNalogSpravkaSeqFunc: func() (int, error) {
//				panic("mock out the GetNalogSpravkaSeq method")
//			},
//			IncrementNalogSpravkaSeqFunc: func(idx int) error {
//				panic("mock out the IncrementNalogSpravkaSeq method")
//			},
//		}
//
//		// use mockedEngineInterface in code that requires EngineInterface
//		// and then make assertions.
//
//	}
type EngineInterfaceMock struct {
	// CloseFunc mocks the Close method.
	CloseFunc func() error

	// CompanyDetailFunc mocks the CompanyDetail method.
	CompanyDetailFunc func() (model.Company, error)

	// FindClientByIDFunc mocks the FindClientByID method.
	FindClientByIDFunc func(id string) (model.Client, error)

	// FindDoctorByIDFunc mocks the FindDoctorByID method.
	FindDoctorByIDFunc func(id string) (model.Doctor, error)

	// FindVisitByIDFunc mocks the FindVisitByID method.
	FindVisitByIDFunc func(id string) (model.Visit, error)

	// FindVisitsByClientIDSinceTillFunc mocks the FindVisitsByClientIDSinceTill method.
	FindVisitsByClientIDSinceTillFunc func(clientID string, startDateEventStr string, endDateEventStr string) ([]model.Visit, error)

	// FindVisitsByDoctorSinceTillFunc mocks the FindVisitsByDoctorSinceTill method.
	FindVisitsByDoctorSinceTillFunc func(doctorID string, startDateEvent string, endDateEvent string) ([]model.Visit, error)

	// GetNalogSpravkaSeqFunc mocks the GetNalogSpravkaSeq method.
	GetNalogSpravkaSeqFunc func() (int, error)

	// IncrementNalogSpravkaSeqFunc mocks the IncrementNalogSpravkaSeq method.
	IncrementNalogSpravkaSeqFunc func(idx int) error

	// calls tracks calls to the methods.
	calls struct {
		// Close holds details about calls to the Close method.
		Close []struct {
		}
		// CompanyDetail holds details about calls to the CompanyDetail method.
		CompanyDetail []struct {
		}
		// FindClientByID holds details about calls to the FindClientByID method.
		FindClientByID []struct {
			// ID is the id argument value.
			ID string
		}
		// FindDoctorByID holds details about calls to the FindDoctorByID method.
		FindDoctorByID []struct {
			// ID is the id argument value.
			ID string
		}
		// FindVisitByID holds details about calls to the FindVisitByID method.
		FindVisitByID []struct {
			// ID is the id argument value.
			ID string
		}
		// FindVisitsByClientIDSinceTill holds details about calls to the FindVisitsByClientIDSinceTill method.
		FindVisitsByClientIDSinceTill []struct {
			// ClientID is the clientID argument value.
			ClientID string
			// StartDateEventStr is the startDateEventStr argument value.
			StartDateEventStr string
			// EndDateEventStr is the endDateEventStr argument value.
			EndDateEventStr string
		}
		// FindVisitsByDoctorSinceTill holds details about calls to the FindVisitsByDoctorSinceTill method.
		FindVisitsByDoctorSinceTill []struct {
			// DoctorID is the doctorID argument value.
			DoctorID string
			// StartDateEvent is the startDateEvent argument value.
			StartDateEvent string
			// EndDateEvent is the endDateEvent argument value.
			EndDateEvent string
		}
		// GetNalogSpravkaSeq holds details about calls to the GetNalogSpravkaSeq method.
		GetNalogSpravkaSeq []struct {
		}
		// IncrementNalogSpravkaSeq holds details about calls to the IncrementNalogSpravkaSeq method.
		IncrementNalogSpravkaSeq []struct {
			// Idx is the idx argument value.
			Idx int
		}
	}
	lockClose                         sync.RWMutex
	lockCompanyDetail                 sync.RWMutex
	lockFindClientByID                sync.RWMutex
	lockFindDoctorByID                sync.RWMutex
	lockFindVisitByID                 sync.RWMutex
	lockFindVisitsByClientIDSinceTill sync.RWMutex
	lockFindVisitsByDoctorSinceTill   sync.RWMutex
	lockGetNalogSpravkaSeq            sync.RWMutex
	lockIncrementNalogSpravkaSeq      sync.RWMutex
}

// Close calls CloseFunc.
func (mock *EngineInterfaceMock) Close() error {
	if mock.CloseFunc == nil {
		panic("EngineInterfaceMock.CloseFunc: method is nil but EngineInterface.Close was just called")
	}
	callInfo := struct {
	}{}
	mock.lockClose.Lock()
	mock.calls.Close = append(mock.calls.Close, callInfo)
	mock.lockClose.Unlock()
	return mock.CloseFunc()
}

// CloseCalls gets all the calls that were made to Close.
// Check the length with:
//
//	len(mockedEngineInterface.CloseCalls())
func (mock *EngineInterfaceMock) CloseCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockClose.RLock()
	calls = mock.calls.Close
	mock.lockClose.RUnlock()
	return calls
}

// CompanyDetail calls CompanyDetailFunc.
func (mock *EngineInterfaceMock) CompanyDetail() (model.Company, error) {
	if mock.CompanyDetailFunc == nil {
		panic("EngineInterfaceMock.CompanyDetailFunc: method is nil but EngineInterface.CompanyDetail was just called")
	}
	callInfo := struct {
	}{}
	mock.lockCompanyDetail.Lock()
	mock.calls.CompanyDetail = append(mock.calls.CompanyDetail, callInfo)
	mock.lockCompanyDetail.Unlock()
	return mock.CompanyDetailFunc()
}

// CompanyDetailCalls gets all the calls that were made to CompanyDetail.
// Check the length with:
//
//	len(mockedEngineInterface.CompanyDetailCalls())
func (mock *EngineInterfaceMock) CompanyDetailCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockCompanyDetail.RLock()
	calls = mock.calls.CompanyDetail
	mock.lockCompanyDetail.RUnlock()
	return calls
}

// FindClientByID calls FindClientByIDFunc.
func (mock *EngineInterfaceMock) FindClientByID(id string) (model.Client, error) {
	if mock.FindClientByIDFunc == nil {
		panic("EngineInterfaceMock.FindClientByIDFunc: method is nil but EngineInterface.FindClientByID was just called")
	}
	callInfo := struct {
		ID string
	}{
		ID: id,
	}
	mock.lockFindClientByID.Lock()
	mock.calls.FindClientByID = append(mock.calls.FindClientByID, callInfo)
	mock.lockFindClientByID.Unlock()
	return mock.FindClientByIDFunc(id)
}

// FindClientByIDCalls gets all the calls that were made to FindClientByID.
// Check the length with:
//
//	len(mockedEngineInterface.FindClientByIDCalls())
func (mock *EngineInterfaceMock) FindClientByIDCalls() []struct {
	ID string
} {
	var calls []struct {
		ID string
	}
	mock.lockFindClientByID.RLock()
	calls = mock.calls.FindClientByID
	mock.lockFindClientByID.RUnlock()
	return calls
}

// FindDoctorByID calls FindDoctorByIDFunc.
func (mock *EngineInterfaceMock) FindDoctorByID(id string) (model.Doctor, error) {
	if mock.FindDoctorByIDFunc == nil {
		panic("EngineInterfaceMock.FindDoctorByIDFunc: method is nil but EngineInterface.FindDoctorByID was just called")
	}
	callInfo := struct {
		ID string
	}{
		ID: id,
	}
	mock.lockFindDoctorByID.Lock()
	mock.calls.FindDoctorByID = append(mock.calls.FindDoctorByID, callInfo)
	mock.lockFindDoctorByID.Unlock()
	return mock.FindDoctorByIDFunc(id)
}

// FindDoctorByIDCalls gets all the calls that were made to FindDoctorByID.
// Check the length with:
//
//	len(mockedEngineInterface.FindDoctorByIDCalls())
func (mock *EngineInterfaceMock) FindDoctorByIDCalls() []struct {
	ID string
} {
	var calls []struct {
		ID string
	}
	mock.lockFindDoctorByID.RLock()
	calls = mock.calls.FindDoctorByID
	mock.lockFindDoctorByID.RUnlock()
	return calls
}

// FindVisitByID calls FindVisitByIDFunc.
func (mock *EngineInterfaceMock) FindVisitByID(id string) (model.Visit, error) {
	if mock.FindVisitByIDFunc == nil {
		panic("EngineInterfaceMock.FindVisitByIDFunc: method is nil but EngineInterface.FindVisitByID was just called")
	}
	callInfo := struct {
		ID string
	}{
		ID: id,
	}
	mock.lockFindVisitByID.Lock()
	mock.calls.FindVisitByID = append(mock.calls.FindVisitByID, callInfo)
	mock.lockFindVisitByID.Unlock()
	return mock.FindVisitByIDFunc(id)
}

// FindVisitByIDCalls gets all the calls that were made to FindVisitByID.
// Check the length with:
//
//	len(mockedEngineInterface.FindVisitByIDCalls())
func (mock *EngineInterfaceMock) FindVisitByIDCalls() []struct {
	ID string
} {
	var calls []struct {
		ID string
	}
	mock.lockFindVisitByID.RLock()
	calls = mock.calls.FindVisitByID
	mock.lockFindVisitByID.RUnlock()
	return calls
}

// FindVisitsByClientIDSinceTill calls FindVisitsByClientIDSinceTillFunc.
func (mock *EngineInterfaceMock) FindVisitsByClientIDSinceTill(clientID string, startDateEventStr string, endDateEventStr string) ([]model.Visit, error) {
	if mock.FindVisitsByClientIDSinceTillFunc == nil {
		panic("EngineInterfaceMock.FindVisitsByClientIDSinceTillFunc: method is nil but EngineInterface.FindVisitsByClientIDSinceTill was just called")
	}
	callInfo := struct {
		ClientID          string
		StartDateEventStr string
		EndDateEventStr   string
	}{
		ClientID:          clientID,
		StartDateEventStr: startDateEventStr,
		EndDateEventStr:   endDateEventStr,
	}
	mock.lockFindVisitsByClientIDSinceTill.Lock()
	mock.calls.FindVisitsByClientIDSinceTill = append(mock.calls.FindVisitsByClientIDSinceTill, callInfo)
	mock.lockFindVisitsByClientIDSinceTill.Unlock()
	return mock.FindVisitsByClientIDSinceTillFunc(clientID, startDateEventStr, endDateEventStr)
}

// FindVisitsByClientIDSinceTillCalls gets all the calls that were made to FindVisitsByClientIDSinceTill.
// Check the length with:
//
//	len(mockedEngineInterface.FindVisitsByClientIDSinceTillCalls())
func (mock *EngineInterfaceMock) FindVisitsByClientIDSinceTillCalls() []struct {
	ClientID          string
	StartDateEventStr string
	EndDateEventStr   string
} {
	var calls []struct {
		ClientID          string
		StartDateEventStr string
		EndDateEventStr   string
	}
	mock.lockFindVisitsByClientIDSinceTill.RLock()
	calls = mock.calls.FindVisitsByClientIDSinceTill
	mock.lockFindVisitsByClientIDSinceTill.RUnlock()
	return calls
}

// FindVisitsByDoctorSinceTill calls FindVisitsByDoctorSinceTillFunc.
func (mock *EngineInterfaceMock) FindVisitsByDoctorSinceTill(doctorID string, startDateEvent string, endDateEvent string) ([]model.Visit, error) {
	if mock.FindVisitsByDoctorSinceTillFunc == nil {
		panic("EngineInterfaceMock.FindVisitsByDoctorSinceTillFunc: method is nil but EngineInterface.FindVisitsByDoctorSinceTill was just called")
	}
	callInfo := struct {
		DoctorID       string
		StartDateEvent string
		EndDateEvent   string
	}{
		DoctorID:       doctorID,
		StartDateEvent: startDateEvent,
		EndDateEvent:   endDateEvent,
	}
	mock.lockFindVisitsByDoctorSinceTill.Lock()
	mock.calls.FindVisitsByDoctorSinceTill = append(mock.calls.FindVisitsByDoctorSinceTill, callInfo)
	mock.lockFindVisitsByDoctorSinceTill.Unlock()
	return mock.FindVisitsByDoctorSinceTillFunc(doctorID, startDateEvent, endDateEvent)
}

// FindVisitsByDoctorSinceTillCalls gets all the calls that were made to FindVisitsByDoctorSinceTill.
// Check the length with:
//
//	len(mockedEngineInterface.FindVisitsByDoctorSinceTillCalls())
func (mock *EngineInterfaceMock) FindVisitsByDoctorSinceTillCalls() []struct {
	DoctorID       string
	StartDateEvent string
	EndDateEvent   string
} {
	var calls []struct {
		DoctorID       string
		StartDateEvent string
		EndDateEvent   string
	}
	mock.lockFindVisitsByDoctorSinceTill.RLock()
	calls = mock.calls.FindVisitsByDoctorSinceTill
	mock.lockFindVisitsByDoctorSinceTill.RUnlock()
	return calls
}

// GetNalogSpravkaSeq calls GetNalogSpravkaSeqFunc.
func (mock *EngineInterfaceMock) GetNalogSpravkaSeq() (int, error) {
	if mock.GetNalogSpravkaSeqFunc == nil {
		panic("EngineInterfaceMock.GetNalogSpravkaSeqFunc: method is nil but EngineInterface.GetNalogSpravkaSeq was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetNalogSpravkaSeq.Lock()
	mock.calls.GetNalogSpravkaSeq = append(mock.calls.GetNalogSpravkaSeq, callInfo)
	mock.lockGetNalogSpravkaSeq.Unlock()
	return mock.GetNalogSpravkaSeqFunc()
}

// GetNalogSpravkaSeqCalls gets all the calls that were made to GetNalogSpravkaSeq.
// Check the length with:
//
//	len(mockedEngineInterface.GetNalogSpravkaSeqCalls())
func (mock *EngineInterfaceMock) GetNalogSpravkaSeqCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetNalogSpravkaSeq.RLock()
	calls = mock.calls.GetNalogSpravkaSeq
	mock.lockGetNalogSpravkaSeq.RUnlock()
	return calls
}

// IncrementNalogSpravkaSeq calls IncrementNalogSpravkaSeqFunc.
func (mock *EngineInterfaceMock) IncrementNalogSpravkaSeq(idx int) error {
	if mock.IncrementNalogSpravkaSeqFunc == nil {
		panic("EngineInterfaceMock.IncrementNalogSpravkaSeqFunc: method is nil but EngineInterface.IncrementNalogSpravkaSeq was just called")
	}
	callInfo := struct {
		Idx int
	}{
		Idx: idx,
	}
	mock.lockIncrementNalogSpravkaSeq.Lock()
	mock.calls.IncrementNalogSpravkaSeq = append(mock.calls.IncrementNalogSpravkaSeq, callInfo)
	mock.lockIncrementNalogSpravkaSeq.Unlock()
	return mock.IncrementNalogSpravkaSeqFunc(idx)
}

// IncrementNalogSpravkaSeqCalls gets all the calls that were made to IncrementNalogSpravkaSeq.
// Check the length with:
//
//	len(mockedEngineInterface.IncrementNalogSpravkaSeqCalls())
func (mock *EngineInterfaceMock) IncrementNalogSpravkaSeqCalls() []struct {
	Idx int
} {
	var calls []struct {
		Idx int
	}
	mock.lockIncrementNalogSpravkaSeq.RLock()
	calls = mock.calls.IncrementNalogSpravkaSeq
	mock.lockIncrementNalogSpravkaSeq.RUnlock()
	return calls
}