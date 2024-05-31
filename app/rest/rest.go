package rest

import (
	"context"
	"fmt"
	"github.com/didip/tollbooth/v6"
	"github.com/didip/tollbooth_chi"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/theshamuel/medregistry20/app/service"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Rest struct {
	DataService restInterface
	Version     string
	URI         string
	ReportPath  string
	httpServer  *http.Server
	lock        sync.Mutex
}

type restInterface interface {
	BuildReportPeriodByDoctorBetweenDateEvent(doctorID string, startDateEvent, endDateEvent string) ([]byte, error)
	BuildReportVisitResult(visitID string) ([]byte, error)
	BuildReportNalogSpravka(req service.ReportNalogSpravkaReq) ([]byte, error)
	BuildReportContract(req service.ReportContractReq) ([]byte, error)
	BuildReportOfPeriodProfit(startDateEvent, endDateEvent string) ([]byte, error)
}

// Run http server
func (r *Rest) Run(port int) {
	log.Printf("[INFO] Run http server on port %d", port)
	r.lock.Lock()
	r.httpServer = r.buildHTTPServer(port, r.routes())
	//TODO: theshamuel insert errorLogger
	r.lock.Unlock()
	err := r.httpServer.ListenAndServe()
	log.Printf("[WARN] http server terminated, %s", err)
}

// Shutdown http server
func (r *Rest) Shutdown() {
	log.Println("[WARN] shutdown http server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r.lock.Lock()
	if r.httpServer != nil {
		if err := r.httpServer.Shutdown(ctx); err != nil {
			log.Printf("[ERROR] http shutdown error, %s", err)
		}
		log.Println("[DEBUG] shutdown http server completed")
	}
	r.lock.Unlock()
}

func (r *Rest) buildHTTPServer(port int, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      120 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func (r *Rest) routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Throttle(1000), middleware.RealIP, middleware.Recoverer, middleware.Logger)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-XSRF-Token", "X-JWT"},
		ExposedHeaders:   []string{"Authorization"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	//health check api
	router.Use(corsMiddleware.Handler)
	router.Route("/", func(api chi.Router) {
		api.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(5, nil)))
		api.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			_, err := w.Write([]byte(fmt.Sprintln("pong")))
			if err != nil {
				log.Printf("[ERROR] cannot write response #%v", err)
			}
		})
	})

	router.Route("/api/v2/", func(rapi chi.Router) {
		//app api
		rapi.Group(func(api chi.Router) {
			api.Use(middleware.Timeout(30 * time.Second))
			api.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(50, nil)))
			api.Use(middleware.NoCache)
			api.Get("/reports/file/reportPeriodByDoctor/{doctorId}/{startDateEvent}/{endDateEvent}/{fileReportName}", r.reportPeriodByDoctorBetweenDateEvent)
			api.Get("/reports/file/reportVisitResult/{visitId}/{fileReportName}", r.reportVisitResult)
			api.Get("/reports/file/reportNalogSpravka/{clientID}/{dateEventFrom}/{dateEventTo}/{payerFio}/{genderOfPayer}/{relationPayerToClient}/{isClientSelfPayer}/{fileReportName}", r.reportNalogSpravka)
			api.Get("/reports/file/reportContract/{clientID}/{visitID}/{dateEvent}/{fileReportName}", r.reportContract)
			api.Get("/reports/file/reportPeriodProfit/{startDateEvent}/{endDateEvent}/{fileReportName}", r.reportOfPeriodProfit)
		})
	})

	return router
}

func (r *Rest) reportPeriodByDoctorBetweenDateEvent(w http.ResponseWriter, req *http.Request) {
	doctorID := chi.URLParam(req, "doctorId")
	startDateEvent := chi.URLParam(req, "startDateEvent")
	endDateEvent := chi.URLParam(req, "endDateEvent")
	file, err := r.DataService.BuildReportPeriodByDoctorBetweenDateEvent(doctorID, startDateEvent, endDateEvent)
	if err != nil {
		log.Printf("[ERROR] cannot build report by doctor in period response %#v", err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(file)
	if err != nil {
		log.Printf("[ERROR] cannot write by doctor in period response %#v", err)
		return
	}
}

func (r *Rest) reportVisitResult(w http.ResponseWriter, req *http.Request) {
	visitID := chi.URLParam(req, "visitId")
	log.Printf("[INFO] reportPeriodByDoctorBetweenDateEvent params visitId=%s", visitID)

	file, err := r.DataService.BuildReportVisitResult(visitID)
	if err != nil {
		log.Printf("[ERROR] cannot build report by visit %#v", err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(file)
	if err != nil {
		log.Printf("[ERROR] cannot write visit result response %#v", err)
		return
	}
}

func (r *Rest) reportNalogSpravka(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] reportNalogSpravka")
	clientID := chi.URLParam(req, "clientID")
	dateEventFrom := chi.URLParam(req, "dateEventFrom")
	dateEventTo := chi.URLParam(req, "dateEventTo")
	payerFIO := chi.URLParam(req, "payerFio")
	genderOfPayer := chi.URLParam(req, "genderOfPayer")
	relationPayerToClient := chi.URLParam(req, "relationPayerToClient")
	//TODO: add error checker
	isClientSelfPayer, _ := strconv.ParseBool(chi.URLParam(req, "isClientSelfPayer"))

	reportReq := service.ReportNalogSpravkaReq{
		ClientID:              clientID,
		DateFrom:              dateEventFrom,
		DateTo:                dateEventTo,
		PayerFIO:              payerFIO,
		GenderOfPayer:         genderOfPayer,
		RelationClientToPayer: relationPayerToClient,
		IsClientSelfPayer:     isClientSelfPayer,
	}

	file, err := r.DataService.BuildReportNalogSpravka(reportReq)
	if err != nil {
		log.Printf("[ERROR] cannot build report by visit %#v", err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(file)
	if err != nil {
		log.Printf("[ERROR] cannot write visit result response %#v", err)
		return
	}
}

func (r *Rest) reportContract(w http.ResponseWriter, req *http.Request) {
	log.Printf("[INFO] reportContract")
	clientID := chi.URLParam(req, "clientID")
	visitID := chi.URLParam(req, "visitID")
	dateEvent := chi.URLParam(req, "dateEvent")

	reportReq := service.ReportContractReq{
		ClientID:  clientID,
		VisitID:   visitID,
		DateEvent: dateEvent,
	}

	file, err := r.DataService.BuildReportContract(reportReq)
	if err != nil {
		log.Printf("[ERROR] cannot build report contract %#v", err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(file)
	if err != nil {
		log.Printf("[ERROR] cannot write contract response %#v", err)
		return
	}
}

func (r *Rest) reportOfPeriodProfit(w http.ResponseWriter, req *http.Request) {
	startDateEvent := chi.URLParam(req, "startDateEvent")
	endDateEvent := chi.URLParam(req, "endDateEvent")
	file, err := r.DataService.BuildReportOfPeriodProfit(startDateEvent, endDateEvent)
	if err != nil {
		log.Printf("[ERROR] cannot build report of profit for repiod %#v", err)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = w.Write(file)
	if err != nil {
		log.Printf("[ERROR] cannot write by doctor in period response %#v", err)
		return
	}
}
