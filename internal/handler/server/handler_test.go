package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Pklerik/gophKeep/internal/auth"
	"github.com/Pklerik/gophKeep/internal/cryptography"
	"github.com/Pklerik/gophKeep/internal/models"
	repoMocks "github.com/Pklerik/gophKeep/internal/repository/mocks"
	"github.com/Pklerik/gophKeep/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*gomock.Controller, *repoMocks.MockUserRepositoryInterface, *repoMocks.MockSecretRepositoryInterface) {
	ctrl := gomock.NewController(t)
	userRepo := repoMocks.NewMockUserRepositoryInterface(ctrl)
	secretRepo := repoMocks.NewMockSecretRepositoryInterface(ctrl)
	return ctrl, userRepo, secretRepo
}

func TestRegister_MethodNotAllowed(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	rr := httptest.NewRecorder()

	h.Register(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestRegister_Success(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	// Expect GetUserByUsername -> not found (return error)
	ur.EXPECT().GetUserByUsername("newuser").Return(nil, errors.New("not found"))
	ur.EXPECT().CreateUser(gomock.Any(), "newuser", gomock.Any()).Return(nil)

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	auth.SetSecretKey("handler-test")

	body := models.AuthRequest{Username: "newuser", Password: "strongpass"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.Register(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var resp models.AuthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

func TestLogin_Success(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	hash := cryptography.HashPassword("password1")
	user := &models.User{ID: "u1", Username: "u1", PasswordHash: hash}

	ur.EXPECT().GetUserByUsername("u1").Return(user, nil)

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	auth.SetSecretKey("handler-test")

	body := models.AuthRequest{Username: "u1", Password: "password1"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp models.AuthResponse
	err := json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
}

func TestCreateGetListUpdateDeleteSecret_Flow(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	// For authentication: create a token for user "uid1"
	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	// CreateSecret expectations
	sr.EXPECT().CreateSecret(gomock.Any()).DoAndReturn(func(s *models.Secret) error {
		s.ID = "sid1"
		return nil
	})

	// GetSecret expectation
	sr.EXPECT().GetSecretByID("sid1").Return(&models.Secret{ID: "sid1", UserID: "uid1", Title: "t"}, nil)

	// ListSecrets expectation
	sr.EXPECT().GetUserSecrets("uid1").Return([]models.Secret{{ID: "sid1", UserID: "uid1"}}, nil)

	// UpdateSecret expectation
	sr.EXPECT().GetSecretByID("sid1").Return(&models.Secret{ID: "sid1", UserID: "uid1"}, nil)
	sr.EXPECT().UpdateSecret(gomock.Any()).Return(nil)

	// DeleteSecret expectation
	sr.EXPECT().DeleteSecret("sid1", "uid1").Return(nil)

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	// CreateSecret
	createReq := map[string]interface{}{"type": models.SecretTypeText, "title": "t", "data": "d", "metadata": "m"}
	cb, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/secret", bytes.NewReader(cb))
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	h.CreateSecret(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	// GetSecret - use chi context for path parameter
	req = httptest.NewRequest(http.MethodGet, "/secret/sid1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "sid1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	h.GetSecret(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// ListSecrets
	req = httptest.NewRequest(http.MethodGet, "/secrets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	h.ListSecrets(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// UpdateSecret
	updateReq := map[string]interface{}{"type": models.SecretTypeText, "title": "t2", "data": "d2", "metadata": "m2"}
	ub, _ := json.Marshal(updateReq)
	req = httptest.NewRequest(http.MethodPut, "/secret/sid1", bytes.NewReader(ub))
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("id", "sid1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	h.UpdateSecret(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// DeleteSecret
	req = httptest.NewRequest(http.MethodDelete, "/secret/sid1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	ctx = chi.NewRouteContext()
	ctx.URLParams.Add("id", "sid1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	h.DeleteSecret(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
}

func TestExtractUserID_MissingAuth(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	uid, err := h.extractUserID(req)
	// Current implementation returns ("", nil) for missing header; assert that
	assert.Equal(t, "", uid)
	assert.NoError(t, err)
}

func TestRegister_InvalidJSON(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte("notjson")))
	rr := httptest.NewRecorder()

	h.Register(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestRegister_CreateUserError(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	ur.EXPECT().GetUserByUsername("uerr").Return(nil, errors.New("not found"))
	ur.EXPECT().CreateUser(gomock.Any(), "uerr", gomock.Any()).Return(errors.New("dbfail"))

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	body := models.AuthRequest{Username: "uerr", Password: "strongpass"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.Register(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLogin_InvalidJSON(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader([]byte("badjson")))
	rr := httptest.NewRecorder()

	h.Login(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestLogin_AuthenticationFailed(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	ur.EXPECT().GetUserByUsername("nouser").Return(nil, errors.New("not found"))

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	body := models.AuthRequest{Username: "nouser", Password: "pw"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(b))
	rr := httptest.NewRecorder()

	h.Login(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestCreateSecret_UnauthorizedAndInvalidBody(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	// Unauthorized: invalid token
	req := httptest.NewRequest(http.MethodPost, "/secret", bytes.NewReader([]byte("{}")))
	req.Header.Set("Authorization", "Bearer badtoken")
	rr := httptest.NewRecorder()
	h.CreateSecret(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	// Invalid body with valid token
	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")
	req = httptest.NewRequest(http.MethodPost, "/secret", bytes.NewReader([]byte("notjson")))
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	h.CreateSecret(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetUpdateDelete_MissingID(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	// GetSecret missing id
	req := httptest.NewRequest(http.MethodGet, "/secret", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	h.GetSecret(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// UpdateSecret missing id
	req = httptest.NewRequest(http.MethodPut, "/secret", bytes.NewReader([]byte("{}")))
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	h.UpdateSecret(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	// DeleteSecret missing id
	req = httptest.NewRequest(http.MethodDelete, "/secret", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr = httptest.NewRecorder()
	h.DeleteSecret(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateSecret_ServiceError(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	// auth token
	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	// service repo will return error
	sr.EXPECT().CreateSecret(gomock.Any()).Return(errors.New("create fail"))

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	createReq := map[string]interface{}{"type": models.SecretTypeText, "title": "t", "data": "d", "metadata": "m"}
	cb, _ := json.Marshal(createReq)
	req := httptest.NewRequest(http.MethodPost, "/secret", bytes.NewReader(cb))
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	h.CreateSecret(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetSecret_NotFound(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	sr.EXPECT().GetSecretByID("nosuch").Return(nil, errors.New("not found"))

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodGet, "/secret/nosuch", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "nosuch")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	h.GetSecret(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestListSecrets_Error(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	sr.EXPECT().GetUserSecrets("uid1").Return(nil, errors.New("list fail"))

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodGet, "/secrets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	h.ListSecrets(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateSecret_UpdateFail(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	// GetSecret called by service GetSecret
	sr.EXPECT().GetSecretByID("sid1").Return(&models.Secret{ID: "sid1", UserID: "uid1"}, nil)
	sr.EXPECT().UpdateSecret(gomock.Any()).Return(errors.New("update fail"))

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	updateReq := map[string]interface{}{"type": models.SecretTypeText, "title": "t2", "data": "d2", "metadata": "m2"}
	ub, _ := json.Marshal(updateReq)
	req := httptest.NewRequest(http.MethodPut, "/secret/sid1", bytes.NewReader(ub))
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "sid1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	h.UpdateSecret(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteSecret_Fail(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	sr.EXPECT().DeleteSecret("sid1", "uid1").Return(errors.New("delete fail"))

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodDelete, "/secret/sid1", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", "sid1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	h.DeleteSecret(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestMethodNotAllowed_ForHandlers(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	// Login expects POST
	req := httptest.NewRequest(http.MethodGet, "/login", nil)
	rr := httptest.NewRecorder()
	h.Login(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	// CreateSecret expects POST
	req = httptest.NewRequest(http.MethodGet, "/secret", nil)
	rr = httptest.NewRecorder()
	h.CreateSecret(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	// GetSecret expects GET
	req = httptest.NewRequest(http.MethodPost, "/secret", nil)
	rr = httptest.NewRecorder()
	h.GetSecret(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	// ListSecrets expects GET
	req = httptest.NewRequest(http.MethodPost, "/secrets", nil)
	rr = httptest.NewRecorder()
	h.ListSecrets(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	// UpdateSecret expects PUT
	req = httptest.NewRequest(http.MethodPost, "/secret?id=1", nil)
	rr = httptest.NewRecorder()
	h.UpdateSecret(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)

	// DeleteSecret expects DELETE
	req = httptest.NewRequest(http.MethodPost, "/secret?id=1", nil)
	rr = httptest.NewRecorder()
	h.DeleteSecret(rr, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}

func TestListSecrets_NilResults(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	sr.EXPECT().GetUserSecrets("uid1").Return(nil, nil)

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodGet, "/secrets", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	h.ListSecrets(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	// body should be an empty array
	assert.JSONEq(t, "[]", rr.Body.String())
}

func TestUpdateSecret_InvalidJSON(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodPut, "/secret?id=sid1", bytes.NewReader([]byte("badjson")))
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	h.UpdateSecret(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

type errReadCloser struct{ err error }

func (e errReadCloser) Read(p []byte) (int, error) { return 0, e.err }
func (e errReadCloser) Close() error               { return nil }

func TestCreateSecret_ReadError(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodPost, "/secret", errReadCloser{err: errors.New("boom")})
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	h.CreateSecret(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateSecret_ReadError(t *testing.T) {
	ctrl, ur, sr := setup(t)
	defer ctrl.Finish()

	auth.SetSecretKey("handler-test")
	token, _, _ := auth.GenerateToken("uid1")

	uh := service.NewUserService(ur)
	sh := service.NewSecretService(sr)
	h := NewHandler(uh, sh)

	req := httptest.NewRequest(http.MethodPut, "/secret?id=sid1", errReadCloser{err: errors.New("boom")})
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()
	h.UpdateSecret(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
