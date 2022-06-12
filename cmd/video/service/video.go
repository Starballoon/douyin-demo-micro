package service

import (
	"context"
	"douyin-demo-micro/cmd/video/dal"
	"douyin-demo-micro/cmd/video/rpc"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/kitex_gen/video"
	"douyin-demo-micro/util"
	"time"
)

func CreateVideo(ctx context.Context, req *video.CreateVideoRequest) error {
	if len(req.VideoFilename) == 0 || len(req.CoverFilename) == 0 ||
		len(req.Title) == 0 || req.AuthorId == 0 {
		return util.ErrIllegalArguments
	}
	user1, err := rpc.FindUser(ctx, &user.FindUserRequest{Req: &user.IdRequest{UserId: req.AuthorId}})
	if err != nil {
		return util.ErrInternalError
	}
	if user1 == nil || user1.Id == 0 {
		return util.ErrIllegalArguments
	}
	err = dal.CreateVideo(ctx, &dal.Video{
		AuthorID:      req.AuthorId,
		VideoFilename: req.VideoFilename,
		CoverFilename: req.CoverFilename,
		Title:         req.Title,
	})
	if err != nil {
		return err
	}
	return nil
}

func Feed(ctx context.Context, req *video.FeedRequest) ([]*video.Video, *int64, error) {
	var user1 *user.User
	if req.UserId != nil {
		user1, _ = rpc.FindUser(ctx, &user.FindUserRequest{Req: &user.IdRequest{UserId: *req.UserId}})
	}

	videos, err := dal.Feed(ctx, req.LatestTime)
	if err != nil {
		return nil, nil, err
	}

	authorMap, err := getAuthors(ctx, videos, req.UserId)
	if err != nil {
		return nil, nil, err
	}

	favorites := make([]int64, 0)
	if user1 != nil {
		favorites, _ = getFavorites(ctx, videos, user1.Id)
	}

	result := mergeVideoAuthors(videos, authorMap, favorites)
	nextTime := oldestTime(videos).UnixMilli()
	return result, &nextTime, nil
}

func PublishList(ctx context.Context, req *video.PublishListRequest) ([]*video.Video, error) {
	if req.Req.UserId == 0 {
		return nil, util.ErrInternalError
	}
	user1, err := rpc.FindUser(ctx, &user.FindUserRequest{Req: &user.IdRequest{UserId: req.Req.UserId}})
	if err != nil || user1 == nil || user1.Id == 0 {
		return nil, util.ErrInternalError
	}

	videos, err := dal.QueryVideosByAuthorID(ctx, req.Req.UserId)
	if err != nil {
		return nil, err
	}

	authorMap := make(map[int64]*user.User)
	authorMap[user1.Id] = user1

	favorites, _ := getFavorites(ctx, videos, user1.Id)

	result := mergeVideoAuthors(videos, authorMap, favorites)
	return result, nil
}

func FavoriteList(ctx context.Context, req *video.FavoriteListRequest) ([]*video.Video, error) {
	if req.Req.UserId == 0 {
		return nil, util.ErrInternalError
	}
	user1, err := rpc.FindUser(ctx, &user.FindUserRequest{Req: &user.IdRequest{UserId: req.Req.UserId}})
	if err != nil || user1 == nil || user1.Id == 0 {
		return nil, util.ErrInternalError
	}

	// 默认约束不允许一个人同时喜欢一个视频多次
	videoIDs, _ := dal.FavoritesByID(ctx, user1.Id)
	videos, err := dal.QueryVideosByIDs(ctx, videoIDs)
	if err != nil {
		return nil, err
	}

	authorMap, err := getAuthors(ctx, videos, &user1.Id)
	if err != nil {
		return nil, err
	}

	favorites := videoIDs

	result := mergeVideoAuthors(videos, authorMap, favorites)
	return result, nil
}

func CheckVideo(ctx context.Context, req *video.CheckVideoRequest) error {
	if req.VideoId == 0 {
		return util.ErrIllegalArguments
	}
	err := dal.CheckVideo(ctx, req.VideoId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCommentCount(ctx context.Context, req *video.UpdateCommentCountRequest) error {
	if req.VideoId == 0 {
		return util.ErrIllegalArguments
	}
	if req.Delta == 0 {
		return nil
	}
	err := dal.UpdateCommentCount(ctx, req.VideoId, req.Delta)
	if err != nil {
		return err
	}
	return nil
}

// dedupID video author ID 去重
func dedupID(videos []*dal.Video) []int64 {
	authorIDMap := make(map[int64]struct{})
	for _, v := range videos {
		authorIDMap[v.AuthorID] = struct{}{}
	}
	authorIDs := make([]int64, 0, len(authorIDMap))
	for authorID := range authorIDMap {
		authorIDs = append(authorIDs, authorID)
	}
	return authorIDs
}

func getAuthors(ctx context.Context, videos []*dal.Video, userID *int64) (map[int64]*user.User, error) {
	authorIDs := dedupID(videos)
	authors, err := rpc.MGetUser(ctx, &user.MGetUserRequest{
		UserIds: authorIDs,
		UserId:  userID,
	})
	if err != nil {
		return nil, util.ErrInternalError
	}
	authorMap := make(map[int64]*user.User)
	for _, author := range authors {
		authorMap[author.Id] = author
	}
	return authorMap, nil
}

func getFavorites(ctx context.Context, videos []*dal.Video, userID int64) ([]int64, error) {
	videoIDs := make([]int64, 0, len(videos))
	for _, v := range videos {
		videoIDs = append(videoIDs, int64(v.ID))
	}
	favorites, err := dal.FavoritesByIDs(ctx, videoIDs, userID)
	if err != nil {
		return nil, util.ErrInternalError
	}
	return favorites, nil
}

func mergeVideoAuthors(videos []*dal.Video, authorMap map[int64]*user.User, favorites []int64) []*video.Video {
	favoriteMap := make(map[int64]struct{})
	if favorites != nil {
		for _, favorite := range favorites {
			favoriteMap[favorite] = struct{}{}
		}
	}
	result := make([]*video.Video, 0, len(videos))
	for _, video := range videos {
		author, _ := authorMap[video.AuthorID]
		_, favorite := favoriteMap[int64(video.ID)]
		result = append(result, mergeVideoAuthor(video, author, favorite))
	}
	return result
}

func mergeVideoAuthor(v *dal.Video, author *user.User, favorite bool) *video.Video {
	return &video.Video{
		Id:            int64(v.ID),
		Author:        author,
		PlayUrl:       v.VideoFilename,
		CoverUrl:      v.CoverFilename,
		FavoriteCount: v.FavoriteCount,
		CommentCount:  v.CommentCount,
		IsFavorite:    favorite,
		Title:         v.Title,
	}
}

func oldestTime(videos []*dal.Video) time.Time {
	nextTime := time.Now()
	for _, video := range videos {
		if video.Model.CreatedAt.Before(nextTime) {
			nextTime = video.Model.CreatedAt
		}
	}
	return nextTime
}
