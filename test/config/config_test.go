package config_test

import (
	"benzinga-webhook-receiver/internal/api"
	"benzinga-webhook-receiver/pkg/config"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

var _ = Describe("Config", func() {
	var router *gin.Engine

	BeforeEach(func() {
		router = gin.Default()
		api.SetupRoutes(router)
		config.LoadConfig("test/setting/test.handler.config.yaml")
		config.StartBatchTimer()
	})

	AfterEach(func() {
		config.StopBatchTimer() // If you implement a stop functionality in your timer logic
	})

	Describe("Batch Size Trigger", func() {
		It("should send a batch when the batch size is reached", func() {
			for i := 0; i < config.Cfg.BatchSize; i++ {
				requestBody := strings.NewReader(`{"user_id": 1, "total": 1.65}`)
				request := httptest.NewRequest("POST", "/log", requestBody)
				request.Header.Set("Content-Type", "application/json")
				recorder := httptest.NewRecorder()
				router.ServeHTTP(recorder, request)
				if i < config.Cfg.BatchSize-1 {
					Expect(recorder.Code).To(Equal(http.StatusOK))
				}
			}
			// Additional check to confirm batch was sent could be implemented here
		})
	})

	Describe("Batch Interval Trigger", func() {
		It("should send a batch after the batch interval", func() {
			// Send one entry, less than batch size
			requestBody := strings.NewReader(`{"user_id": 1, "total": 1.65}`)
			request := httptest.NewRequest("POST", "/log", requestBody)
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			// Wait for more than the batch interval
			time.Sleep(config.Cfg.BatchInterval + 1*time.Second)

			// Check here to confirm that the batch has been sent
			// This might require you to check some form of log or output depending on your implementation
		})
	})
})
