package skill

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type mockStorage struct {
	Storage
	skills []Skill
	skill  Skill
	kind   string
	err    error
}

func (m *mockStorage) GetByKind(ctx context.Context, kind string) ([]Skill, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.kind = kind
	return m.skills, nil
}

func (m *mockStorage) GetByID(ctx context.Context, id string) (Skill, error) {
	if m.err != nil {
		return Skill{}, m.err
	}
	return m.skill, nil
}

// TODO: eliminate all the gin.Context and use app.Context instead

func TestGetSkillsByKind(t *testing.T) {
	t.Run("should return 200 and skills", func(t *testing.T) {
		id, _ := primitive.ObjectIDFromHex("5e201c51e09c2c084c88a790")
		sks := []Skill{
			{
				ID:          id,
				Name:        "Go",
				Description: "Go is an efficient, statically typed, and concurrent programming language.",
				Logo:        "base64",
				Kind:        "technical",
			},
		}
		mock := &mockStorage{
			skills: sks,
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.GET("/skills/kind/:kindtype", app.NewGinHandler(handler.GetSkillsByKind, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/skills/kind/technical", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "success",
			"message": "",
			"data": [
				{
					"id": "5e201c51e09c2c084c88a790",
					"name": "Go",
					"description": "Go is an efficient, statically typed, and concurrent programming language.",
					"logo": "base64",
					"kind": "technical"
				}
			]
		}`
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "technical", mock.kind)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 500 when service error", func(t *testing.T) {
		mock := &mockStorage{
			err: errors.New("service error"),
		}
		handler := NewSkillHandler(mock)

		rec := httptest.NewRecorder()
		engine := gin.New()
		c := gin.CreateTestContextOnly(rec, engine)
		engine.GET("/skills/kind/:kindtype", app.NewGinHandler(handler.GetSkillsByKind, zap.NewNop()))
		c.Request, _ = http.NewRequest(http.MethodGet, "/skills/kind/soft", nil)

		engine.HandleContext(c)

		resp := rec.Body.String()
		want := `{
			"status": "error",
			"message": "service error"
		}`
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}

func TestGetSkillByID(t *testing.T) {
	t.Run("should return 200 and skills when found the skillID", func(t *testing.T) {
		skillID := "5e201c51e09c2c084c88a790" // #nosec
		id, _ := primitive.ObjectIDFromHex(skillID)
		sk := Skill{
			ID:          id,
			Name:        "Go",
			Description: "Go is an efficient",
			Logo:        "base64",
			Kind:        "technical",
		}
		mock := &mockStorage{
			skill: sk,
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.GET("/skills/:id", app.NewGinHandler(handler.SkillByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/skills/"+skillID, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "success",
			"message": "",
			"data": {
				"id": "5e201c51e09c2c084c88a790",
				"name": "Go",
				"description": "Go is an efficient",
				"logo": "base64",
				"kind": "technical"
			}
		}`
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 404 when no found the skill", func(t *testing.T) {
		skillID := "5e201c51e09c2c084c88a790" // #nosec
		mock := &mockStorage{
			err: mongo.ErrNoDocuments,
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.GET("/skills/:id", app.NewGinHandler(handler.SkillByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/skills/"+skillID, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "error",
			"message": "mongo: no documents in result"
		}`
		assert.Equal(t, 404, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 500 when error from query skill", func(t *testing.T) {
		skillID := "5e201c51e09c2c084c88a790" // #nosec
		mock := &mockStorage{
			err: errors.New("db error"),
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.GET("/skills/:id", app.NewGinHandler(handler.SkillByID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/skills/"+skillID, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "error",
			"message": "db error"
		}`
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}
