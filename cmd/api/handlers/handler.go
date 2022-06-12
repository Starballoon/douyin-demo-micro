package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

func successResp(c *gin.Context) {
	c.JSON(http.StatusOK, Response{StatusCode: 0})
}

func failResp(c *gin.Context, code int, err error) {
	c.JSON(http.StatusOK, Response{
		StatusCode: int32(code),
		StatusMsg:  err.Error()})
}

func loginError(c *gin.Context, code int, err error) {
	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{
			StatusCode: int32(code),
			StatusMsg:  err.Error()},
	})
}
