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
	"time"

	"gitdev.devops.krungthai.com/aster/ariskill/app"
	"gitdev.devops.krungthai.com/aster/ariskill/app/squad"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type mockSquadStorage struct {
	squad *squad.Squad
	time  time.Time
	err   error
}

func middleware(profileID string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("profileID", profileID)
		c.Next()
	}
}

func (ms *mockSquadStorage) GetOneByID(ctx context.Context, id string) (*squad.Squad, error) {
	if ms.err != nil {
		return nil, ms.err
	}
	if ms.squad != nil {
		return ms.squad, nil
	}
	return nil, mongo.ErrNoDocuments
}

func (ms *mockSquadStorage) UpdateByID(ctx context.Context, id string, updateSquad *squad.Squad) (*squad.Squad, error) {
	var sq squad.Squad
	if ms.err != nil {
		return nil, ms.err
	}
	skID, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
	sq = squad.Squad{
		Name:        "Aster",
		Description: "Aster",
		CreatedAt:   ms.time,
		SkillsRatings: []squad.SkillRatings{
			{
				SkillId: skID,
				Ratings: []squad.SkillRating{
					{
						UserId: "uID",
						Score:  2,
					},
				},
			},
		},
	}
	sq.Id, _ = primitive.ObjectIDFromHex("64e17f4ae098346113ae4f62")
	return &sq, nil
}

func TestGetUserSkillRatingBySquadID(t *testing.T) {
	t.Run("should return 200 and user skill in squad", func(t *testing.T) {
		profileID := "999999999999999999991"
		objID1, _ := primitive.ObjectIDFromHex("5e201c51e09c2c084c88a790")
		objSkillID1, _ := primitive.ObjectIDFromHex("64e17f43e098346113ae4f57")
		objSkillID2, _ := primitive.ObjectIDFromHex("64e17f43e098346113ae4f59")
		currentTime := time.Now()
		squad := &squad.Squad{
			Id:          objID1,
			Name:        "Mock",
			Description: "data is mock",
			CreatedAt:   currentTime,
			SkillsRatings: []squad.SkillRatings{
				{
					SkillId: objSkillID1,
					Ratings: []squad.SkillRating{
						{
							UserId: "999999999999999999991",
							Score:  100,
						},
						{
							UserId: "123456789666269492155",
							Score:  50,
						},
					},
				},
				{
					SkillId: objSkillID2,
					Ratings: []squad.SkillRating{
						{
							UserId: "999999999999999999991",
							Score:  80,
						},
						{
							UserId: "123456789666269492155",
							Score:  87,
						},
					},
				},
			},
		}

		mock := &mockSquadStorage{
			squad: squad,
		}

		handler := NewSquadHandler(mock)

		engine := gin.New()
		engine.Use(middleware(profileID))
		engine.GET("/squads/:squadID/rate", app.NewGinHandler(handler.GetUserSkillRatingBySquadID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/5e201c51e09c2c084c88a790/rate", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := fmt.Sprintf(`{
			"status": "success",
			"message": "",
			"data": {
				"id": "5e201c51e09c2c084c88a790",
				"name": "Mock",
				"desc": "data is mock",
				"createdAt": %q,
				"skills": [
					{
						"skid": "64e17f43e098346113ae4f57",
						"score": 100
					},
					{
						"skid": "64e17f43e098346113ae4f59",
						"score": 80
					}
				]
			}
		}`, currentTime.Format(time.RFC3339Nano))
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 404 when squad not found", func(t *testing.T) {
		profileID := "999999999999999999991"
		mock := &mockSquadStorage{
			squad: nil,
		}

		handler := NewSquadHandler(mock)

		engine := gin.New()
		engine.Use(middleware(profileID))
		engine.GET("/squads/:id/rate", app.NewGinHandler(handler.GetUserSkillRatingBySquadID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/55555555555/rate", nil)

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
		profileID := "999999999999999999991"

		mock := &mockSquadStorage{
			err: errors.New("mock error 500"),
		}

		handler := NewSquadHandler(mock)

		engine := gin.New()
		engine.Use(middleware(profileID))
		engine.GET("/squads/:id/rate", app.NewGinHandler(handler.GetUserSkillRatingBySquadID, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/squads/55555555555/rate", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "error",
			"message": "mock error 500"
		}`
		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}

func TestRateSkill(t *testing.T) {
	t.Run("should return 200 and message success", func(t *testing.T) {
		objId, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		rateSkills := RateSkill{
			Ratings: []Rating{
				{
					SkillId: objId,
					Score:   2,
				},
			},
		}
		skID, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		sq := squad.Squad{
			Name:         "Aster",
			TeamleadMail: "",
			Description:  "Aster",
			CreatedAt:    time.Date(2023, 10, 25, 12, 0, 0, 0, time.UTC),
			SkillsRatings: []squad.SkillRatings{
				{
					SkillId: skID,
					Ratings: []squad.SkillRating{
						{
							UserId: "uID",
							Score:  1,
						},
					},
				},
			},
		}
		sq.Id, _ = primitive.ObjectIDFromHex("64e17f4ae098346113ae4f62")

		mock := &mockSquadStorage{
			time:  time.Date(2023, 10, 25, 12, 0, 0, 0, time.UTC),
			squad: &sq,
		}
		jsonBody, _ := json.Marshal(rateSkills)
		handler := NewSquadHandler(mock)

		engine := gin.New()
		engine.POST("/squads/:id/rate", app.NewGinHandler(handler.RateSkills, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/squads/:id/rate", bytes.NewBuffer(jsonBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status":"success",
			"message":"",
			"data":{
			   "id":"64e17f4ae098346113ae4f62",
			   "name":"Aster",
			   "teamleadMail":"",
			   "desc":"Aster",
			   "createdAt":"2023-10-25T12:00:00Z",
			   "skillsRatings":[
				  {
					 "skid":"64e17f4ae098346113ae4f61",
					 "ratings":[
						{
						   "uid":"uID",
						   "score":2
						}
					 ]
				  }
			   ]
			}
		 }`
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 400 when bad request", func(t *testing.T) {
		handler := NewSquadHandler(&mockSquadStorage{})

		engine := gin.New()
		engine.POST("/squads/:id/rate", app.NewGinHandler(handler.RateSkills, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/squads/:id/rate", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
					"status": "error",
					"message": "invalid request"
				}`
		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
	})
	t.Run("should return 404 when no row modified", func(t *testing.T) {
		jsonBody, _ := json.Marshal(squad.Squad{})
		mock := &mockSquadStorage{}
		handler := NewSquadHandler(mock)

		engine := gin.New()
		engine.POST("/squads/:id/rate", app.NewGinHandler(handler.RateSkills, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/squads/:id/rate", bytes.NewBuffer(jsonBody))

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
		objId, _ := primitive.ObjectIDFromHex("64e17f4ae098346113ae4f61")
		rateSkills := RateSkill{
			Ratings: []Rating{
				{
					SkillId: objId,
					Score:   2,
				},
			},
		}
		jsonBody, _ := json.Marshal(rateSkills)
		mock := &mockSquadStorage{
			err: errors.New("service error"),
		}
		handler := NewSquadHandler(mock)

		engine := gin.New()
		engine.POST("/squads/:id/rate", app.NewGinHandler(handler.RateSkills, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/squads/:id/rate", bytes.NewBuffer(jsonBody))

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
