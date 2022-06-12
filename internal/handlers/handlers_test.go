package handlers

import (
	"testing"

	"github.com/golang/mock/gomock"

	_ "github.com/Sur0vy/cows_health.git/internal/mocks"
)

func TestBaseHandler_AddFarm(t *testing.T) {

}

func TestBaseHandler_DelFarm(t *testing.T) {

}

func TestBaseHandler_GetFarms(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	var id int64 = 1
	mockMale := mock.NewMockMale(ctl)
	gomock.InOrder(
		mockMale.EXPECT().Get(id).Return(nil),
	)

	user := NewUser(mockMale)
	err := user.GetUserInfo(id)
	if err != nil {
		t.Errorf("user.GetUserInfo err: %v", err)
	}
}

func TestBaseHandler_Login(t *testing.T) {

}

func TestBaseHandler_Logout(t *testing.T) {

}

func TestBaseHandler_Register(t *testing.T) {

}

func TestBaseHandler_ResponseBadRequest(t *testing.T) {

}

func TestNewBaseHandler(t *testing.T) {

}

func Test_getIDFromJSON(t *testing.T) {

}
