package profile

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type mockSkillStorage struct {
	skillsUser  SkillsByUser
	kind        string
	isFoundData bool
	err         error
}

func (ms *mockSkillStorage) GetSkills(ctx context.Context, id string, kind string) (*SkillsByUser, error) {
	var rs SkillsByUser
	if ms.err != nil {
		return nil, ms.err
	}
	if ms.kind != "technical" && ms.kind != "soft" {
		return nil, ErrInvalidKindOfSkill
	}
	rs.UserID = ms.skillsUser.UserID
	for _, v := range ms.skillsUser.Skills {
		if v.SkillInfo.Kind == kind {
			rs.Skills = append(rs.Skills, v)
		}
	}
	return &rs, nil
}
func (ms *mockSkillStorage) UpdateProfileSkills(id string, set skillset, skills []Skill) error {
	if ms.isFoundData {
		return ms.err
	} else {
		return mongo.ErrNoDocuments
	}
}

func TestGetSkillsByUserID(t *testing.T) {
	t.Run("should return 200 and technical skill", func(t *testing.T) {
		id := "999999999999999999991"
		kind := "technical"
		objIdTechSkill, _ := primitive.ObjectIDFromHex("5e201c51e09c2c084c88a790")
		objIdSoftSkill, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		skill := SkillsByUser{
			UserID: id,
			Skills: []SkillNameScore{
				{
					SkillInfo: SkillInfo{
						ID:          objIdTechSkill,
						Name:        "HTML",
						Description: "mock HTML description",
						Logo:        "https://www.svgrepo.com/show/452228/html-5.svg",
						Kind:        "technical",
					},
					Score: 100,
				},
				{
					SkillInfo: SkillInfo{
						ID:          objIdSoftSkill,
						Name:        "Adaptability",
						Description: "Adaptability is the ability to adjust to new situations and challenges.",
						Logo:        "",
						Kind:        "soft",
					},
					Score: 100,
				},
			},
		}
		mock := &mockSkillStorage{
			skillsUser: skill,
			kind:       kind,
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.GET("/profile/skills", app.NewGinHandler(handler.GetSkillsByUserID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/profile/skills?kind="+kind, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "success",
			"message": "",
			"data": {
				"id": "999999999999999999991",
		 		"skills": [
			 		{
						"skill" : {
							"id": "5e201c51e09c2c084c88a790",
							"name": "HTML",
							"description": "mock HTML description",
							"logo": "https://www.svgrepo.com/show/452228/html-5.svg",
							"kind": "technical"
						},
						"score": 100
					}
				]
			}
		}`
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "technical", mock.kind)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 200 and soft skill", func(t *testing.T) {
		id := "999999999999999999991"
		kind := "soft"
		objIdTechSkill, _ := primitive.ObjectIDFromHex("5e201c51e09c2c084c88a790")
		objIdSoftSkill, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		skill := SkillsByUser{
			UserID: id,
			Skills: []SkillNameScore{
				{
					SkillInfo: SkillInfo{
						ID:          objIdTechSkill,
						Name:        "HTML",
						Description: "mock HTML description",
						Logo:        "https://www.svgrepo.com/show/452228/html-5.svg",
						Kind:        "technical",
					},
					Score: 100,
				},
				{
					SkillInfo: SkillInfo{
						ID:          objIdSoftSkill,
						Name:        "Adaptability",
						Description: "Adaptability is the ability to adjust to new situations and challenges.",
						Logo:        "",
						Kind:        "soft",
					},
					Score: 100,
				},
			},
		}
		mock := &mockSkillStorage{
			skillsUser: skill,
			kind:       kind,
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.GET("/profile/skills", app.NewGinHandler(handler.GetSkillsByUserID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/profile/skills?kind="+kind, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
					"status": "success",
					"message": "",
					"data": {
						"id": "999999999999999999991",
				 		"skills": [
					 		{
								"skill" : {
									"id": "64e17f4ae098346113ae4f61",
									"name": "Adaptability",
									"description": "Adaptability is the ability to adjust to new situations and challenges.",
									"logo": "",
									"kind": "soft"
								},
								"score": 100
							}
						]
					}
				}`
		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, "soft", mock.kind)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 400 when kind is not technical or soft", func(t *testing.T) {
		id := "999999999999999999991"
		kind := "otherkind"
		objIdTechSkill, _ := primitive.ObjectIDFromHex("5e201c51e09c2c084c88a790")
		objIdSoftSkill, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		skill := SkillsByUser{
			UserID: id,
			Skills: []SkillNameScore{
				{
					SkillInfo: SkillInfo{
						ID:          objIdTechSkill,
						Name:        "HTML",
						Description: "mock HTML description",
						Logo:        "https://www.svgrepo.com/show/452228/html-5.svg",
						Kind:        "technical",
					},
					Score: 100,
				},
				{
					SkillInfo: SkillInfo{
						ID:          objIdSoftSkill,
						Name:        "Adaptability",
						Description: "Adaptability is the ability to adjust to new situations and challenges.",
						Logo:        "",
						Kind:        "soft",
					},
					Score: 100,
				},
			},
		}
		mock := &mockSkillStorage{
			skillsUser: skill,
			kind:       kind,
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.GET("/profile/skills", app.NewGinHandler(handler.GetSkillsByUserID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/profile/skills?kind="+kind, nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
					"status": "error",
					"message": "This kind of skill does not exist."
				}`
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 500 when service error", func(t *testing.T) {
		kind := "technical"
		mock := &mockSkillStorage{
			err:  errors.New("service error"),
			kind: kind,
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.GET("/profile/skills", app.NewGinHandler(handler.GetSkillsByUserID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/profile/skills?kind="+kind, nil)

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

func TestUpdateTechnicalSkill(t *testing.T) {
	t.Run("should return 200 and message success", func(t *testing.T) {
		var skill []Skill
		objId, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		skill = append(skill, Skill{
			SkillID: objId,
			Score:   100,
		})
		mock := &mockSkillStorage{
			isFoundData: true,
		}
		jsonBody, _ := json.Marshal(skill)
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.POST("/users/technical", app.NewGinHandler(handler.UpdateTechnicalSkill, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/technical", bytes.NewBuffer(jsonBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
					"status": "success",
					"message": "",
					"data": "updated technical skill"
				}`
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 400 when bad request", func(t *testing.T) {
		handler := NewSkillHandler(&mockSkillStorage{})

		engine := gin.New()
		engine.POST("/users/technical", app.NewGinHandler(handler.UpdateTechnicalSkill, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/technical", bytes.NewBuffer([]byte{}))

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
		var skill []Skill
		objId, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		skill = append(skill, Skill{
			SkillID: objId,
			Score:   100,
		})
		jsonBody, _ := json.Marshal(skill)
		mock := &mockSkillStorage{
			isFoundData: false,
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.POST("/users/technical", app.NewGinHandler(handler.UpdateTechnicalSkill, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/technical", bytes.NewBuffer(jsonBody))

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
		var skill []Skill
		objId, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		skill = append(skill, Skill{
			SkillID: objId,
			Score:   100,
		})
		jsonBody, _ := json.Marshal(skill)
		mock := &mockSkillStorage{
			isFoundData: true,
			err:         errors.New("service error"),
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.POST("/users/technical", app.NewGinHandler(handler.UpdateTechnicalSkill, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/technical", bytes.NewBuffer(jsonBody))

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

func TestUpdateSoftlSkill(t *testing.T) {
	t.Run("should return 200 and message success", func(t *testing.T) {
		var skill []Skill
		objId, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		skill = append(skill, Skill{
			SkillID: objId,
			Score:   100,
		})
		mock := &mockSkillStorage{
			isFoundData: true,
		}
		jsonBody, _ := json.Marshal(skill)
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.POST("/users/soft", app.NewGinHandler(handler.UpdateSoftSkill, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/soft", bytes.NewBuffer(jsonBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
					"status": "success",
					"message": "",
					"data": "updated soft skill"
				}`
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 400 when bad request", func(t *testing.T) {
		handler := NewSkillHandler(&mockSkillStorage{})

		engine := gin.New()
		engine.POST("/users/soft", app.NewGinHandler(handler.UpdateSoftSkill, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/soft", bytes.NewBuffer([]byte{}))

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
		var skill []Skill
		objId, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		skill = append(skill, Skill{
			SkillID: objId,
			Score:   100,
		})
		jsonBody, _ := json.Marshal(skill)
		mock := &mockSkillStorage{
			isFoundData: false,
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.POST("/users/soft", app.NewGinHandler(handler.UpdateSoftSkill, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/soft", bytes.NewBuffer(jsonBody))

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
		var skill []Skill
		objId, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		skill = append(skill, Skill{
			SkillID: objId,
			Score:   100,
		})
		jsonBody, _ := json.Marshal(skill)
		mock := &mockSkillStorage{
			isFoundData: true,
			err:         errors.New("service error"),
		}
		handler := NewSkillHandler(mock)

		engine := gin.New()
		engine.POST("/users/soft", app.NewGinHandler(handler.UpdateSoftSkill, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users/soft", bytes.NewBuffer(jsonBody))

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
