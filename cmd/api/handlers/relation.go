package handlers

import (
	"douyin-demo-micro/cmd/api/rpc"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []*user.User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(ctx *gin.Context) {
	user1, err := findUser(ctx)
	if err != nil {
		failResp(ctx, 1, err)
		return
	}
	// 没有这个字段
	// userID, err := strconv.Atoi(ctx.Query("user_id"))
	action, err := strconv.Atoi(ctx.Query("action_type"))
	// 1-关注， 2-取消关注
	if err != nil || !(action == util.CREATE_FOLLOW || action == util.DELETE_FOLLOW) {
		failResp(ctx, 1, util.ErrIllegalArguments)
		return
	}

	leaderID, err := strconv.Atoi(ctx.Query("to_user_id"))
	// 不允许自己关注自己
	if err != nil || user1.Id == int64(leaderID) {
		failResp(ctx, 1, util.ErrIllegalArguments)
		return
	}
	user2, err := rpc.FindUser(ctx, &user.FindUserRequest{Req: &user.IdRequest{UserId: int64(leaderID)}})
	if err != nil && user2.Id > 0 {
		failResp(ctx, 1, util.ErrIllegalArguments)
		return
	}
	switch action {
	case util.CREATE_FOLLOW:
		err = rpc.Follow(ctx, &user.FollowRequest{LeaderId: int64(leaderID), FollowerId: user1.Id})
	case util.DELETE_FOLLOW:
		err = rpc.Unfollow(ctx, &user.UnfollowRequest{LeaderId: int64(leaderID), FollowerId: user1.Id})
	}
	if err != nil {
		failResp(ctx, 1, err)
		return
	}
	successResp(ctx)
}

// FollowList all users have same follow list
func FollowList(ctx *gin.Context) {
	// 如果token中的id和request body中的user_id字段不同，该表现出什么行为？
	user1, err := findUser(ctx)
	if err != nil {
		failResp(ctx, 1, err)
		return
	}
	leaders, err := rpc.FollowList(ctx, &user.FollowListRequest{Req: &user.IdRequest{UserId: user1.Id}})
	if err != nil {
		ctx.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  util.ErrInternalError.Error(),
			},
		})
	}
	ctx.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		//UserList: []service.User{DemoUser},
		UserList: leaders,
	})
}

// FollowerList all users have same follower list
func FollowerList(ctx *gin.Context) {
	// 如果token中的id和request body中的user_id字段不同，该表现出什么行为？
	user1, err := findUser(ctx)
	if err != nil {
		failResp(ctx, 1, err)
		return
	}
	// TODO 这里会拉取别人的粉丝吗？
	followers, err := rpc.FollowerList(ctx, &user.FollowerListRequest{Req: &user.IdRequest{UserId: user1.Id}})
	if err != nil {
		ctx.JSON(http.StatusOK, UserListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  util.ErrInternalError.Error(),
			},
		})
	}
	ctx.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		//UserList: []service.User{DemoUser},
		UserList: followers,
	})
}
