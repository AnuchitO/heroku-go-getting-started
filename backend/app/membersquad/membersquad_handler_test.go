package membersquad

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

type mockStorage struct {
	Storage
	members []Member
	err     error
}

func (m *mockStorage) GetByID(ctx context.Context, id string) (*Member, error) {
	if m.err != nil {
		return nil, m.err
	}

	var rsl Member
	for _, member := range m.members {
		if member.ID == id {
			rsl = member
		}
	}

	return &rsl, nil
}

func (m *mockStorage) UpdateMySquad(ctx context.Context, idUser string, update []Squad) error {
	if m.err != nil {
		return m.err
	}

	for index, member := range m.members {
		if member.ID == idUser {
			m.members[index].MySquads = update
			break
		}
	}

	return nil
}

func (m *mockStorage) GetAllBySquadId(ctx context.Context, squadId primitive.ObjectID) ([]Member, error) {
	if m.err != nil {
		return nil, m.err
	}

	var rsl []Member
	for _, member := range m.members {
		for _, squad := range member.MySquads {
			if squad.SquadID == squadId {
				rsl = append(rsl, member)
			}
		}
	}

	return rsl, nil
}

func fmtJsonBody(body map[string]any) *bytes.Buffer {
	rsl, _ := json.Marshal(body)
	return bytes.NewBuffer(rsl)
}

func makeHexObjId(s string) primitive.ObjectID {
	id, _ := primitive.ObjectIDFromHex(s)
	return id
}

func TestAddMemberSquad(t *testing.T) {
	gin.SetMode("test")
	t.Run("should return 200 when include is true", func(t *testing.T) {
		sqId, _ := primitive.ObjectIDFromHex("650bfb051ac125739cfb7a3e")
		reqBody := map[string]any{
			"sqid":    sqId,
			"include": true,
			"members": []map[string]any{
				{
					"uid":  "3",
					"role": "member",
				},
			},
		}

		mock := &mockStorage{}
		handler := NewMemberSquadHandler(mock)

		engine := gin.New()
		engine.PUT("/member-squads/members", app.NewGinHandler(handler.AddMemberSquad, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/member-squads/members", fmtJsonBody(reqBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "success",
			"message": ""
		}`

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 200 when include is false", func(t *testing.T) {
		sqId, _ := primitive.ObjectIDFromHex("650bfb051ac125739cfb7a3e")
		reqBody := map[string]any{
			"sqid":    sqId,
			"include": false,
			"members": []map[string]any{
				{
					"uid":  "2",
					"role": "member",
				},
			},
		}

		mock := &mockStorage{}
		handler := NewMemberSquadHandler(mock)

		engine := gin.New()
		engine.PUT("/member-squads/members", app.NewGinHandler(handler.AddMemberSquad, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/member-squads/members", fmtJsonBody(reqBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "success",
			"message": ""
		}`

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 400 when bind error (invalid request body format)", func(t *testing.T) {
		mock := &mockStorage{}
		handler := NewMemberSquadHandler(mock)

		engine := gin.New()
		engine.PUT("/member-squads/members", app.NewGinHandler(handler.AddMemberSquad, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/member-squads/members", nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "error",
			"message": "Request is invalid format"
		}`

		assert.Equal(t, 400, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should return 500 when add error (member is not found)", func(t *testing.T) {
		sqId, _ := primitive.ObjectIDFromHex("650bfb051ac125739cfb7a3e")
		reqBody := map[string]any{
			"sqid":    sqId,
			"include": false,
			"members": []map[string]any{
				{
					"uid":  "1000",
					"role": "role",
				},
			},
		}

		mock := &mockStorage{err: errors.New("mongo: no documents in result")}
		handler := NewMemberSquadHandler(mock)

		engine := gin.New()
		engine.PUT("/member-squads/members", app.NewGinHandler(handler.AddMemberSquad, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/member-squads/members", fmtJsonBody(reqBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "error",
			"message": "mongo: no documents in result"
		}`

		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("should be success when add squad to each user", func(t *testing.T) {
		sqId, _ := primitive.ObjectIDFromHex("000000000000000000000000")
		reqBody := map[string]any{
			"sqid":    sqId,
			"include": false,
			"members": []map[string]any{
				{
					"uid":  "1",
					"role": "member",
				},
				{
					"uid":  "2",
					"role": "member",
				},
				{
					"uid":  "3",
					"role": "member",
				},
			},
		}

		mockMember := []Member{
			{
				ID: "1",
				MySquads: []Squad{
					{
						SquadID: makeHexObjId("000000000000000000000001"),
						Role:    "member",
					},
					{
						SquadID: makeHexObjId("000000000000000000000003"),
						Role:    "member",
					},
				},
			},
			{
				ID: "2",
				MySquads: []Squad{
					{
						SquadID: makeHexObjId("000000000000000000000000"),
						Role:    "member",
					},
				},
			},
			{
				ID: "3",
				MySquads: []Squad{
					{
						SquadID: makeHexObjId("000000000000000000000002"),
						Role:    "member",
					},
				},
			},
		}
		mock := &mockStorage{members: mockMember}
		handler := NewMemberSquadHandler(mock)

		engine := gin.New()
		engine.PUT("/member-squads/members", app.NewGinHandler(handler.AddMemberSquad, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/member-squads/members", fmtJsonBody(reqBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "success",
			"message": ""
		}`

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)

		for _, member := range mock.members {
			isMemberInSquad := false
			for _, squad := range member.MySquads {
				if squad.SquadID == sqId {
					isMemberInSquad = true
					break
				}
			}

			assert.Equal(t, true, isMemberInSquad, "member is not found in squad")
		}
	})

	t.Run("should return 500 when add error (member is not found)", func(t *testing.T) {
		sqId, _ := primitive.ObjectIDFromHex("000000000000001000000000")
		reqBody := map[string]any{
			"sqid":    sqId,
			"include": false,
			"members": []map[string]any{
				{
					"uid": "10000",
				},
			},
		}

		mock := &mockStorage{err: mongo.ErrNoDocuments}
		handler := NewMemberSquadHandler(mock)

		engine := gin.New()
		engine.PUT("/member-squads/members", app.NewGinHandler(handler.AddMemberSquad, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPut, "/member-squads/members", fmtJsonBody(reqBody))

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()
		want := `{
			"status": "error",
			"message": "mongo: no documents in result"
		}`

		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}

func TestDeleteMemberSquad(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		type sqIds struct {
			id string
		}
		sq := sqIds{
			id: "650bfb051ac125739cfb7a3e",
		}
		mock := &mockStorage{}
		handler := NewMemberSquadHandler(mock)

		engine := gin.New()
		engine.DELETE("/member-squads/:squadID/members", app.NewGinHandler(handler.DeleteMemberSquad, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/member-squads/%v/members", sq.id), nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := `{
			"status": "success",
			"message": ""
		}`

		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("Delete error - Squad ID invalid format", func(t *testing.T) {
		sqId := "1"
		mock := &mockStorage{err: primitive.ErrInvalidHex}
		handler := NewMemberSquadHandler(mock)

		engine := gin.New()
		engine.DELETE("/member-squads/:squadID/members", app.NewGinHandler(handler.DeleteMemberSquad, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/member-squads/%v/members", sqId), nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := fmt.Sprintf(`{
			"status": "error",
			"message": %q
		}`, primitive.ErrInvalidHex.Error())

		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})

	t.Run("Delete error - No user has this squad", func(t *testing.T) {
		sqId := "000000000000000000000000"
		mock := &mockStorage{err: mongo.ErrNoDocuments}
		handler := NewMemberSquadHandler(mock)

		engine := gin.New()
		engine.DELETE("/member-squads/:squadID/members", app.NewGinHandler(handler.DeleteMemberSquad, zap.NewNop()))
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/member-squads/%v/members", sqId), nil)

		engine.ServeHTTP(rec, req)

		resp := rec.Body.String()

		want := fmt.Sprintf(`{
			"status": "error",
			"message": %q
		}`, mongo.ErrNoDocuments.Error())

		assert.Equal(t, 500, rec.Code)
		assert.JSONEq(t, want, resp)
	})
}
