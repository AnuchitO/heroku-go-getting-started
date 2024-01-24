package profile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type mockUserStorage struct {
	isFoundData bool
	err         error
}

func (ms *mockUserStorage) AboutMeUpdate(id string, about aboutme) error {
	if ms.err != nil {
		return ms.err
	}

	if !ms.isFoundData {
		return mongo.ErrNoDocuments
	}
	return nil
}

func TestUpdateAboutMe(t *testing.T) {
	t.Run("should return 200 and message success", func(t *testing.T) {
		aboutme := aboutme{
			AboutMe:     "this is a unit test",
			SocialMedia: []string{"facebook", "ig", "other"},
			Tags:        []string{"testTags", "testTags2"},
		}
		jsonBody, _ := json.Marshal(aboutme)
		mock := &mockUserStorage{
			isFoundData: true,
		}
		handler := NewUserHandler(mock)

		engine := gin.New()
		engine.PUT("/profile", app.NewGinHandler(handler.UpdateAboutMe, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(jsonBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
						"status": "success",
						"message": ""
					}`
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 400 when bad request", func(t *testing.T) {
		handler := NewUserHandler(&mockUserStorage{})

		engine := gin.New()
		engine.PUT("/profile", app.NewGinHandler(handler.UpdateAboutMe, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer([]byte{}))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
						"status": "error",
						"message": "EOF"
					}`
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 404 when no row modified", func(t *testing.T) {
		aboutme := aboutme{
			AboutMe:     "this is a unit test",
			SocialMedia: []string{"facebook", "ig", "other"},
			Tags:        []string{"testTags", "testTags2"},
		}
		jsonBody, _ := json.Marshal(aboutme)
		mock := &mockUserStorage{
			isFoundData: false,
		}
		handler := NewUserHandler(mock)

		engine := gin.New()
		engine.PUT("/profile", app.NewGinHandler(handler.UpdateAboutMe, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(jsonBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
						"status": "error",
						"message": %q
					}`, mongo.ErrNoDocuments)
		assert.Equal(t, 404, rec.Code)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 500 when service error", func(t *testing.T) {
		aboutme := aboutme{
			AboutMe:     "this is a unit test",
			SocialMedia: []string{"facebook", "ig", "other"},
			Tags:        []string{"testTags", "testTags2"},
		}
		jsonBody, _ := json.Marshal(aboutme)
		mock := &mockUserStorage{
			isFoundData: true,
			err:         errors.New("service error"),
		}
		handler := NewUserHandler(mock)

		engine := gin.New()
		engine.PUT("/profile", app.NewGinHandler(handler.UpdateAboutMe, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewBuffer(jsonBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
					"status": "error",
					"message": "service error"
				}`
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}
