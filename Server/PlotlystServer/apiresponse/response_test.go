package apiresponse

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestErrorWritesStructuredContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	Error(context, http.StatusBadRequest, "invalid_request", "Request is invalid", "field detail")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}
	var response ErrorBody
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unable to decode response: %v", err)
	}
	if response.Error.Code != "invalid_request" || response.Error.Message != "Request is invalid" {
		t.Fatalf("unexpected error response: %+v", response)
	}
}
