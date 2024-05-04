package api_test

import (
	"benzinga-webhook-receiver/internal/api"
	"benzinga-webhook-receiver/pkg/config"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	ginkgo "github.com/onsi/ginkgo"
	gomega "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Handler", func() {
	ginkgo.BeforeEach(func() {
		config.LoadConfig("test.handler.yaml")
	})

	ginkgo.Context("POST /log", func() {
		ginkgo.It("should handle the request correctly", func() {
			router := gin.Default()
			router.POST("/log", api.HandleLogPost)

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

			gomega.Expect(recorder.Code).To(gomega.Equal(200))
		})

		ginkgo.It("should reject invalid JSON format", func() {
			router := gin.Default()
			router.POST("/log", api.HandleLogPost)

			requestBody := strings.NewReader(`{"user_id": 1, "total": 1.65, "title":}`)
			request := httptest.NewRequest("POST", "/log", requestBody)
			request.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, request)

			gomega.Expect(recorder.Code).To(gomega.Equal(http.StatusBadRequest))
		})

		ginkgo.It("should return the correct response content", func() {
			router := gin.Default()
			router.POST("/log", api.HandleLogPost)
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
			gomega.Expect(recorder.Code).To(gomega.Equal(200))
			gomega.Expect(recorder.Body.String()).To(gomega.ContainSubstring("logged"))
		})
	})
})
