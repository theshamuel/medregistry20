package rest

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/theshamuel/medregistry20/app/service"
	"github.com/theshamuel/medregistry20/app/store"
	"github.com/theshamuel/medregistry20/app/store/model"
	"go.uber.org/goleak"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRest_Shutdown(t *testing.T) {
	srv := Rest{}
	done := make(chan bool)

	go func() {
		time.Sleep(200 * time.Millisecond)
		srv.Shutdown()
		close(done)
	}()

	st := time.Now()
	srv.Run(8888)
	assert.True(t, time.Since(st).Seconds() < 1, "should take about 1s")
	<-done
}

func TestRest_Run(t *testing.T) {
	srv := Rest{}
	port := generateRndPort()
	go func() {
		srv.Run(port)
	}()

	waitHTTPServer(port)

	client := http.Client{}

	resp, err := client.Get(fmt.Sprintf("http://localhost:%d/ping", port))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	srv.Shutdown()
}

func TestRest_Ping(t *testing.T) {
	ts, _, _, teardown := startHTTPServer()
	defer teardown()

	res, code := getRequest(t, ts.URL+"/ping")
	assert.Equal(t, "pong\n", res)
	assert.Equal(t, http.StatusOK, code)
}

func TestRest_ReportVisitResult(t *testing.T) {
	ts, _, engineMock, teardown := startHTTPServer()
	defer teardown()
	_, code := getRequest(t, ts.URL+"/api/v2/reports/file/reportVisitResult/1/test.xlsx")
	assert.Equal(t, http.StatusOK, code)

	if len(engineMock.FindVisitByIDCalls()) != 1 {
		t.Errorf("[ERROR] ReportVisitResult was called %d times, FindVisitByIDCalls was called %d times", len(engineMock.FindVisitByIDCalls()),
			len(engineMock.FindVisitByIDCalls()))
	}
	if len(engineMock.FindClientByIDCalls()) != 1 {
		t.Errorf("[ERROR] ReportVisitResult was called %d times, FindClientByIDCalls was called %d times", len(engineMock.FindClientByIDCalls()),
			len(engineMock.FindClientByIDCalls()))
	}
	if len(engineMock.FindDoctorByIDCalls()) != 1 {
		t.Errorf("[ERROR] ReportVisitResult was called %d times, FindDoctorByIDCalls was called %d times", len(engineMock.FindDoctorByIDCalls()),
			len(engineMock.FindDoctorByIDCalls()))
	}
	if len(engineMock.CompanyDetailCalls()) != 1 {
		t.Errorf("[ERROR] ReportVisitResult was called %d times, CompanyDetailCalls was called %d times", len(engineMock.CompanyDetailCalls()),
			len(engineMock.CompanyDetailCalls()))
	}
}

func startHTTPServer() (ts *httptest.Server, rest *Rest, engineMock *store.EngineInterfaceMock, gracefulTeardown func()) {
	engineMock = &store.EngineInterfaceMock{
		FindVisitByDoctorSinceTillFunc: func(doctorID string, startDateEvent string, endDateEvent string) ([]model.Visit, error) {
			visits := []model.Visit{
				{
					ID:                        "1",
					DoctorID:                  "1",
					DoctorName:                "Alex Alex",
					DoctorExcludedFromReports: false,
				},
				{
					ID:                        "2",
					DoctorID:                  "1",
					DoctorName:                "Alex Alex",
					DoctorExcludedFromReports: false,
				},
			}
			return visits, nil
		},
		FindVisitByIDFunc: func(id string) (model.Visit, error) {
			return model.Visit{}, nil
		},
		FindClientByIDFunc: func(id string) (model.Client, error) {
			return model.Client{}, nil
		},
		FindDoctorByIDFunc: func(id string) (model.Doctor, error) {
			return model.Doctor{}, nil
		},
		CompanyDetailFunc: func() (model.Company, error) {
			return model.Company{}, nil
		},
	}
	service := &service.DataStore{
		Engine: engineMock,
	}

	rest = &Rest{
		DataService: service,
		Version:     "test",
		URI:         "http://localhost:8888/api/v1/",
		ReportPath:  "/srv/reportPath",
	}
	ts = httptest.NewServer(rest.routes())
	gracefulTeardown = func() {
		ts.Close()
	}
	return ts, rest, engineMock, gracefulTeardown
}

func generateRndPort() (port int) {
	for i := 0; i < 10; i++ {
		port = 40000 + int(rand.Int31n(10000))
		if ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port)); err == nil {
			_ = ln.Close()
			break
		}
	}
	return port
}

func waitHTTPServer(port int) {
	for i := 0; i < 10; i++ {
		time.Sleep(time.Second * 1)
		conn, _ := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), time.Millisecond*10)
		if conn != nil {
			_ = conn.Close()
			break
		}
	}
}

func getRequest(t *testing.T, url string) (data string, statusCode int) {
	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	return string(body), resp.StatusCode
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
