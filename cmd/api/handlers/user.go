package handlers

import (
	"douyin-demo-micro/cmd/api/rpc"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/util"
	"errors"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User user.User `json:"user"`
}

type LogIn struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

func Register(ctx *gin.Context) {
	var login LogIn
	if err := ctx.ShouldBind(&login); err != nil {
		return
	}
	_ = rpc.CreateUser(ctx, &user.CreateUserRequest{
		Req: &user.NameRequest{
			Username: login.Username,
			Password: login.Password,
		}})
}

func UserInfo(c *gin.Context) {
	user, err := findUser(c)
	if err != nil {
		loginError(c, 1, err)
		return
	}
	c.JSON(http.StatusOK, UserResponse{
		Response: Response{StatusCode: 0},
		User:     *user,
	})
}

func AuthMidware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		// Token signing key, default signing method HS256
		Key:        []byte(util.JWTKEY),
		Timeout:    time.Hour,
		MaxRefresh: time.Hour,
		// token field in JSON
		IdentityKey: util.IDENTITY,
		// 前端有些token放在query，有些用post放在request body的form里
		TokenLookup: "query: token, form: token",
		// 检查身份，合法返回User，不合法返回错误
		// Authenticator,PayloadFunc,LoginResponse三个函数在LoginHandler中依次调用
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login LogIn
			if err := c.ShouldBind(&login); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			tmpUser, err := rpc.CheckUser(c, &user.CheckUserRequest{
				Req: &user.NameRequest{
					Username: login.Username,
					Password: login.Password,
				}})
			if err != nil {
				return nil, jwt.ErrFailedAuthentication
			}
			c.Set(util.JWT_TMP_KEY, tmpUser)
			return tmpUser, nil
		},
		// Authenticator返回的第一个参数作为PayloadFunc的第一个入参
		// PayloadFunc的返回值会加密保存到token中
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*user.User); ok {
				return jwt.MapClaims{
					// TODO token里加不加邮箱之类的字段？
					util.IDENTITY: v.Id,
				}
			}
			return jwt.MapClaims{}
		},
		LoginResponse: func(c *gin.Context, code int, message string, time time.Time) {
			// TODO 同时给token和user_id可以考虑校验一下
			// 接口要求返回user_id，但是这个中间件的LoginHandler在调用LoginResponse前
			// 不会把token绑定到context去，所以拿不到，只能在Authenticator里绑定
			userID := int64(0)
			if tmpUser, ok := c.Get(util.JWT_TMP_KEY); ok {
				userID = tmpUser.(*user.User).Id
				c.Set(util.JWT_TMP_KEY, struct{}{})
			}
			c.JSON(http.StatusOK, UserLoginResponse{
				Response: Response{StatusCode: 0},
				UserId:   userID,
				Token:    message,
			})
		},
		// 下面两次参数在检查授权时依次调用，前者返回身份标识并作为后者的第一个入参
		// 参考该库的middlewareImpl函数，两者可以合并为controller.findUser
		// 根据token判别身份合法性，这里只考虑id合法，中间件包括有效时间等合法
		IdentityHandler: func(c *gin.Context) interface{} {
			id, err := identify(c)
			if err != nil {
				return int64(0)
			}
			return id
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			id, ok := data.(int64)
			if !ok || id < 1 {
				return false
			}
			if _, err := rpc.FindUser(c, &user.FindUserRequest{Req: &user.IdRequest{UserId: id}}); err == nil {
				return true
			}
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			loginError(c, code, errors.New(message))
		},
	})
}

func AuthWithHelper() (*jwt.GinJWTMiddleware, func(c *gin.Context), error) {
	auth, err := AuthMidware()
	return auth, func(c *gin.Context) {
		// Resolve identity for /feed because without authentication
		claims, err := auth.GetClaimsFromJWT(c)
		if err != nil {
			return
		}
		if claims[util.IDENTITY] != nil {
			// jwt库会检查context中绑定的JWT_PAYLOAD这一map
			c.Set("JWT_PAYLOAD", jwt.MapClaims{util.IDENTITY: claims[util.IDENTITY]})
		}
	}, err
}

func identify(c *gin.Context) (int64, error) {
	claims := jwt.ExtractClaims(c)
	id, ok := claims[util.IDENTITY]
	if !ok {
		return 0, util.ErrIllegalArguments
	}
	if v, ok := id.(float64); !ok || v < 1 {
		return 0, util.ErrIllegalArguments
	}
	return int64(claims[util.IDENTITY].(float64)), nil
}

func findUser(c *gin.Context) (*user.User, error) {
	id, err := identify(c)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}
	if user, err := rpc.FindUser(c, &user.FindUserRequest{Req: &user.IdRequest{UserId: id}}); err == nil {
		return user, nil
	}
	return nil, util.ErrInternalError
}
