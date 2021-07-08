// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	v1 "github.com/marmotedu/api/apiserver/v1"

	srvv1 "github.com/marmotedu/iam/internal/apiserver/service/v1"
	"github.com/marmotedu/iam/internal/apiserver/store"
	_ "github.com/marmotedu/iam/pkg/validator"
)

func TestUserHandler_ChangePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	user := &v1.User{
		Password: "$2a$10$KqZhl5WStpa2K.ddEyzyf.zXllEXP4gIG8xQUgMhU1ZvMUn/Ta5um",
	}
	mockFactory := store.NewMockFactory(ctrl)
	mockUserStore := store.NewMockUserStore(ctrl)
	mockUserStore.EXPECT().Get(gomock.Any(), gomock.Eq("colin"), gomock.Any()).Return(user, nil)
	mockFactory.EXPECT().Users().Return(mockUserStore)

	mockService := srvv1.NewMockService(ctrl)
	mockUserSrv := srvv1.NewMockUserSrv(ctrl)
	mockUserSrv.EXPECT().ChangePassword(gomock.Any(), gomock.Any()).Return(nil)
	mockService.EXPECT().Users().Return(mockUserSrv)

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	body := bytes.NewBufferString(`{"oldPassword":"Admin@2020","newPassword":"Colin@2021"}`)
	c.Request, _ = http.NewRequest("PUT", "/v1/users/colin/change_password", body)
	c.Params = []gin.Param{{Key: "name", Value: "colin"}}
	c.Request.Header.Set("Content-Type", "application/json")

	type fields struct {
		srv   srvv1.Service
		store store.Factory
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "default",
			fields: fields{
				srv:   mockService,
				store: mockFactory,
			},
			args: args{
				c: c,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{
				srv:   tt.fields.srv,
				store: tt.fields.store,
			}
			u.ChangePassword(tt.args.c)
		})
	}
}
