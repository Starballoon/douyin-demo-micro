package dal

import (
	"context"
	"douyin-demo-micro/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

var email = []string{"test1", "bytedance.com"}
var password = "12345678"

func TestMain(m *testing.M) {
	InitDB()
}

func TestCreateUser(t *testing.T) {
	err := CreateUser(context.Background(), &User{
		EmailAccount: email[0],
		EmailDomain:  email[1],
		Password:     util.Encrypt(password),
		Name:         "",
	})
	assert.Nil(t, err)
}

func TestQueryUserByAccount(t *testing.T) {
	err := CreateUser(context.Background(), &User{
		EmailAccount: email[0],
		EmailDomain:  email[1],
		Password:     util.Encrypt(password),
		// TODO 新用户名生成策略
		Name: "NewUser",
	})
	assert.Nil(t, err)

	user, err := QueryUserByAccount(context.Background(), email)
	assert.Nil(t, err)
	assert.NotEqual(t, user.ID, 0)
}
