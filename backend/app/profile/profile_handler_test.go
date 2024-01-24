package profile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type mockProfileStorage struct {
	profile []Profile
	err     error
}

func (s *mockProfileStorage) GetByID(ctx context.Context, id string) (*Profile, error) {
	if s.err != nil {
		return nil, s.err
	}
	if s.profile == nil {
		return nil, mongo.ErrNoDocuments
	}
	return &s.profile[0], nil
}

func TestGetUser(t *testing.T) {
	t.Run("Should return 200 and user", func(t *testing.T) {
		currentTime := time.Now()
		profile := []Profile{{
			ID:              "999999999999999999991",
			Email:           "ariskill@arise.tech",
			EmployeeID:      "66666",
			FirstName:       "Ariser1",
			LastName:        "skills",
			JobRole:         "backend",
			MySquads:        nil,
			CreatedAt:       currentTime,
			UpdatedAt:       currentTime,
			CreatedBy:       "system",
			UpdatedBy:       "someone",
			SoftSkills:      nil,
			TechnicalSkills: nil,
		}}
		profile[0].AboutMe = ""

		mockStorage := &mockProfileStorage{profile: profile}
		handler := NewProfileHandler(mockStorage)

		engine := gin.New()
		engine.GET("/user", app.NewGinHandler(handler.User, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/user", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		bodywantJson := map[string]any{
			"sub":             "999999999999999999991",
			"employeeId":      "66666",
			"email":           "ariskill@arise.tech",
			"givenName":       "Ariser1",
			"familyName":      "skills",
			"jobRole":         "backend",
			"createdAt":       currentTime,
			"updatedAt":       currentTime,
			"createdBy":       "system",
			"updatedBy":       "someone",
			"softSkills":      nil,
			"technicalSkills": nil,
			"aboutMe":         "",
			"squadId":         nil,
			"socialMedias":    nil,
			"tags":            nil,
		}
		wantRaw := map[string]any{
			"status":  "success",
			"message": "",
			"data":    bodywantJson,
		}
		wantJson, _ := json.Marshal(wantRaw)
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, string(wantJson), resp)
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		err := errors.New("UnexpectedError")

		mockStorage := &mockProfileStorage{err: err}
		handler := NewProfileHandler(mockStorage)

		engine := gin.New()
		engine.GET("/user", app.NewGinHandler(handler.User, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/user", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "error",
			"message": "UnexpectedError"
		}`
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("User not found", func(t *testing.T) {
		mockStorage := &mockProfileStorage{
			profile: nil,
		}
		handler := NewProfileHandler(mockStorage)

		engine := gin.New()
		engine.GET("/user", app.NewGinHandler(handler.User, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/user", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "error",
			"message": %q
		}`, mongo.ErrNoDocuments)
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}
