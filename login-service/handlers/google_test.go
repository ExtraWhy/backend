package handlers_test

import(
	"github.com/gin-gonic/gin"
	"login-service/handlers"
	"testing"
)

func TestGoogleCallback(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		c *gin.Context
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handlers.GoogleCallback(tt.c)
		})
	}
}

