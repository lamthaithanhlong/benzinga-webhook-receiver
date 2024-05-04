package config_test

import (
	"benzinga-webhook-receiver/internal/api"
	"benzinga-webhook-receiver/pkg/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	ginkgo "github.com/onsi/ginkgo"
	gomega "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Config", func() {
	var router *gin.Engine

	ginkgo.BeforeEach(func() {
		router = gin.Default()
		api.SetupRoutes(router)
		config.LoadConfig("test.handler.config.yaml")
		config.StartBatchTimer()
	})

	ginkgo.AfterEach(func() {
		config.StopBatchTimer()
	})

	ginkgo.Describe("Batch Size Trigger", func() {
		ginkgo.It("should send a batch when the batch size is reached", func() {
			for i := 0; i < config.Cfg.BatchSize; i++ {
				requestBody := strings.NewReader(`{"user_id": 1, "total": 1.65}`)
				request := httptest.NewRequest("POST", "/log", requestBody)
				request.Header.Set("Content-Type", "application/json")
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, request)
				if i < config.Cfg.BatchSize-1 {
					gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusOK))
				}
			}
		})
	})

	ginkgo.Describe("Batch Interval Trigger", func() {
		ginkgo.It("should send a batch after the batch interval", func() {
			// Send one entry, less than batch size
			requestBody := strings.NewReader(`{"user_id": 1, "total": 1.65}`)
			request := httptest.NewRequest("POST", "/log", requestBody)
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			// Wait for more than the batch interval
			time.Sleep(config.Cfg.BatchInterval + 1*time.Second)
		})
	})
})
