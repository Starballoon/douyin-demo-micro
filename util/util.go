package util

import (
	"bytes"
	"crypto"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"math"
	"regexp"
	"time"
)

var reg = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*@[a-zA-Z0-9]+(\.[a-zA-Z0-9]+)+`)

func ValidEmail(email string) bool {
	// 在前端非邮箱字符串都是合法的……
	return reg.MatchString(email)
}

func LoginCheck(username, password string) bool {
	// TODO format check of password
	return ValidEmail(username) && len(password) >= 6
}

// MaxI golang 为什么不给一个MaxI比较整数……
func MaxI(a, b int) int {
	return int(math.Max(float64(a), float64(b)))
}

func MaxI64(a, b int64) int64 {
	return int64(math.Max(float64(a), float64(b)))
}

func ValidTime(t int64) bool {
	return math.Abs(float64(time.Now().Year()-time.UnixMilli(t).Year())) < 50.0
}

func Encrypt(password string) []byte {
	hash := crypto.SHA512.New()
	hash.Write([]byte(password))
	return hash.Sum(nil)
}

func HijackRequestBody(c *gin.Context) {
	bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	fmt.Println(string(bodyBytes))
}
