package controllers_test

import (
	"ExamenMeLiMutante/controllers"
	_ "ExamenMeLiMutante/controllers"
	"ExamenMeLiMutante/models"
	"ExamenMeLiMutante/server"
	servicesMock "ExamenMeLiMutante/test/mocks/services"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Mutant", func() {
	var ctrl *gomock.Controller
	var mutantsService *servicesMock.MockIMutantService
	var mutantsReq models.MutantRequest
	var mutantsRes models.MutantResponse
	var statsRes *models.MutantsStats
	var dna []string

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mutantsService = servicesMock.NewMockIMutantService(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when a request is successfully processed", func() {
		It("should return status code ok and isMutant = true", func() {
			dna = []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}
			mutantsReq = models.MutantRequest{DnaChain: dna}
			mutantsRes = models.MutantResponse{IsMutant: true}
			mutantsService.EXPECT().VerifyMutant(mutantsReq).Return(mutantsRes)
			jsonBody, _ := json.Marshal(mutantsReq)
			e := echo.New()
			server.SetupValidator(e)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := controllers.NewMutantController(mutantsService)
			h.VerifyMutantStatus(c)

			Expect(rec.Code).Should(Equal(http.StatusOK))
		})

		It("should return status code forbidden and isMutant = false", func() {
			dna = []string{"ATGC", "CGTA", "ATGC", "CGTG"}
			mutantsReq = models.MutantRequest{DnaChain: dna}
			mutantsRes = models.MutantResponse{IsMutant: false}
			mutantsService.EXPECT().VerifyMutant(mutantsReq).Return(mutantsRes)
			jsonBody, _ := json.Marshal(mutantsReq)
			e := echo.New()
			server.SetupValidator(e)
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			h := controllers.NewMutantController(mutantsService)
			h.VerifyMutantStatus(c)

			Expect(rec.Code).Should(Equal(http.StatusForbidden))
		})

		Context("when a dna request has not NxN dimensions", func() {
			It("should return status code bad request", func() {
				dna = []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACT"}
				mutantsReq = models.MutantRequest{DnaChain: dna}
				jsonBody, _ := json.Marshal(mutantsReq)
				e := echo.New()
				server.SetupValidator(e)
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				h := controllers.NewMutantController(mutantsService)
				err := h.VerifyMutantStatus(c)
				log.Info(err)
				Expect(rec.Code).Should(Equal(http.StatusBadRequest))
			})
		})

		Context("when a dna request has invalid characters", func() {
			It("should return status code bad request", func() {
				dna = []string{"ATGCGA", "CAGTGC", "TTATGT", "AGAAGG", "CCCCTA", "TCACX"}
				mutantsReq = models.MutantRequest{DnaChain: dna}
				jsonBody, _ := json.Marshal(mutantsReq)
				e := echo.New()
				server.SetupValidator(e)
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				h := controllers.NewMutantController(mutantsService)
				err := h.VerifyMutantStatus(c)
				log.Info(err)
				Expect(rec.Code).Should(Equal(http.StatusBadRequest))
			})
		})

		Context("when a dna request body is not correctly modeled", func() {
			It("should return status code bad request", func() {
				dna := []int{123, 1234, 12345}
				invalidReq := struct{ Dna []int }{dna}
				jsonBody, _ := json.Marshal(invalidReq)
				e := echo.New()
				server.SetupValidator(e)
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(jsonBody)))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				h := controllers.NewMutantController(mutantsService)
				err := h.VerifyMutantStatus(c)
				log.Info(err)
				Expect(rec.Code).Should(Equal(http.StatusBadRequest))
			})
		})

		Context("when a status request is requested", func() {
			It("should return status code ok and status body", func() {
				stats := models.MutantsStats{
					CountHumanDna:  100,
					CountMutantDna: 10,
					Ratio:          0.1,
				}
				statsRes = &stats
				mutantsService.EXPECT().GetSubjectsStats().Return(statsRes, nil)
				e := echo.New()
				server.SetupValidator(e)
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				h := controllers.NewMutantController(mutantsService)
				h.GetMutantStats(c)
				Expect(rec.Code).Should(Equal(http.StatusOK))
			})

			It("should return status code internal server error and no status body", func() {
				stats := models.MutantsStats{}
				statsRes = &stats
				mutantsService.EXPECT().GetSubjectsStats().Return(nil, fmt.Errorf("Error: %s", "errorMock"))
				e := echo.New()
				server.SetupValidator(e)
				req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
				req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				rec := httptest.NewRecorder()
				c := e.NewContext(req, rec)
				h := controllers.NewMutantController(mutantsService)
				h.GetMutantStats(c)
				Expect(rec.Code).Should(Equal(http.StatusInternalServerError))
			})
		})
	})
})
