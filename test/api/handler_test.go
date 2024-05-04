package api_test

import (
	"benzinga-webhook-receiver/internal/api"
	"benzinga-webhook-receiver/pkg/config"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Handler", func() {
	BeforeEach(func() {
		// Load configurations for testing
		config.LoadConfig("test.handler.config.yaml")
	})

	Context("POST /log", func() {
		It("should handle the request correctly", func() {
			router := gin.Default()
			router.POST("/log", api.HandleLogPost)

			// Complete JSON payload including all necessary fields
			requestBody := strings.NewReader(`{
				"user_id": 1,
				"total": 1.65,
				"title": "delectus aut autem",
				"meta": {
					"logins": [
						{
							"time": "2020-08-08T01:52:50Z",
							"ip": "0.0.0.0"
						}
					],
					"phone_numbers": {
						"home": "555-1212",
						"mobile": "123-5555"
					}
				},
				"completed": false
			}`)
			request := httptest.NewRequest("POST", "/log", requestBody)
			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(200))
		})

		It("should reject invalid JSON format", func() {
			router := gin.Default()
			router.POST("/log", api.HandleLogPost)

			// Provide a JSON string with syntax errors
			requestBody := strings.NewReader(`{"user_id": 1, "total": 1.65, "title":}`)
			request := httptest.NewRequest("POST", "/log", requestBody)
			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(http.StatusBadRequest)) // Check for status code 400
		})

		It("should return the correct response content", func() {
			router := gin.Default()
			requestBody := strings.NewReader(`{"user_id": 1, "total": 1.65, "title": "Test", "meta": {}, "completed": true}`)
			request := httptest.NewRequest("POST", "/log", requestBody)
			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			Expect(recorder.Code).To(Equal(200))
			Expect(recorder.Body.String()).To(ContainSubstring("logged"))
		})
	})
})
