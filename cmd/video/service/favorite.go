package service

import (
	"context"
	"douyin-demo-micro/cmd/video/dal"
	"douyin-demo-micro/cmd/video/rpc"
	"douyin-demo-micro/kitex_gen/user"
	"douyin-demo-micro/kitex_gen/video"
	"douyin-demo-micro/util"
)

func Favorite(ctx context.Context, req *video.FavoriteRequest) error {
	if req.UserId == 0 || req.VideoId == 0 {
		return util.ErrIllegalArguments
	}
	if user, err := rpc.FindUser(ctx, &user.FindUserRequest{Req: &user.IdRequest{UserId: req.UserId}}); err != nil || user.Id == 0 {
		return util.ErrIllegalArguments
	}
	if err := dal.CheckVideo(ctx, req.VideoId); err != nil {
		return util.ErrIllegalArguments
	}

	favorite, err := dal.QueryUnscoped(ctx, req.VideoId, req.UserId)
	if err != nil {
		return err
	}

	if favorite.ID > 0 && favorite.Model.DeletedAt.Valid {
		if err = dal.RecoverFavor(ctx, favorite); err != nil {
			return err
		}
		if err = dal.UpdateFavoriteCount(ctx, req.VideoId, 1); err != nil {
			return util.ErrInternalError
		}
	} else if favorite.ID == 0 {
		if err = dal.CreateFavor(ctx, &dal.Favorite{VideoID: req.VideoId, UserID: req.UserId}); err != nil {
			return err
		}
		if err = dal.UpdateFavoriteCount(ctx, req.VideoId, 1); err != nil {
			return util.ErrInternalError
		}
	}
	return nil
}

func Unfavorite(ctx context.Context, req *video.UnfavoriteRequest) error {
	if req.UserId == 0 || req.VideoId == 0 {
		return util.ErrIllegalArguments
	}
	if user, err := rpc.FindUser(ctx, &user.FindUserRequest{Req: &user.IdRequest{UserId: req.UserId}}); err != nil || user.Id == 0 {
		return util.ErrIllegalArguments
	}
	if err := dal.CheckVideo(ctx, req.VideoId); err != nil {
		return util.ErrIllegalArguments
	}

	favorite, err := dal.QueryUnscoped(ctx, req.VideoId, req.UserId)
	if err != nil {
		return err
	}

	if favorite.ID > 0 && !favorite.Model.DeletedAt.Valid {
		if err = dal.DeleteFavor(ctx, favorite); err != nil {
			return util.ErrInternalError
		}
		if err = dal.UpdateFavoriteCount(ctx, req.VideoId, -1); err != nil {
			return util.ErrInternalError
		}
	}
	return nil
}
