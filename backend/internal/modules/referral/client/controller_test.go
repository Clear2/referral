package client

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type fakeControllerService struct {
	registration *RegistrationView
}

func (service *fakeControllerService) GenerateInvitation(context.Context, int) (string, error) {
	return "", nil
}

func (service *fakeControllerService) Register(context.Context, RegisterInput) (*RegistrationView, error) {
	return service.registration, nil
}

func (service *fakeControllerService) Dashboard(context.Context, int) (*DashboardView, error) {
	return nil, nil
}

type fakeSessionStarter struct {
	userID int
}

func (starter *fakeSessionStarter) StartSession(_ *gin.Context, userID int) error {
	starter.userID = userID
	return nil
}

func TestRegisterWithNameAndEmailStartsInviteeSession(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeControllerService{registration: &RegistrationView{Invitee: UserView{ID: 12}}}
	sessions := &fakeSessionStarter{}
	controller := &Controller{service: service, sessions: sessions}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/api/v1/referrals/register", strings.NewReader(`{
		"code":"ABCDEFGH",
		"name":"Alice",
		"email":"alice@example.com"
	}`))
	request.Header.Set("Content-Type", "application/json")
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = request

	controller.Register(ctx)

	if sessions.userID != 12 {
		t.Fatalf("session user ID = %d, want 12", sessions.userID)
	}
	if recorder.Code != 200 {
		t.Fatalf("response status = %d, want 200", recorder.Code)
	}
}
