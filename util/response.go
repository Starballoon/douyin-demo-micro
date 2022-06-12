package util

import (
	"douyin-demo-micro/kitex_gen/user"
	"time"
)

func SuccessResp() *user.BaseResp {
	return &user.BaseResp{
		StatusCode:    0,
		StatusMessage: "Success",
		ServiceTime:   time.Now().Unix(),
	}
}

func FailResp(code int64, err error) *user.BaseResp {
	return &user.BaseResp{
		StatusCode:    code,
		StatusMessage: err.Error(),
		ServiceTime:   time.Now().Unix(),
	}
}
