namespace go comment

include "user.thrift"

struct Comment{
    1:i64 id
    2:user.User user
    3:string content
    4:string create_date
}

struct CreateCommentRequest{
    1:i64 reviewer_id
    2:i64 video_id
    3:string content
}

struct CreateCommentResponse{
    1:Comment comment
    2:user.BaseResp resp
}

struct DeleteCommentRequest{
    1:i64 user_id
    2:i64 comment_id
}

struct CommentListRequest{
    1:i64 video_id
}

struct CommentListResponse{
    1:list<Comment> comment_list
    2:user.BaseResp resp
}

service CommentService{
    CreateCommentResponse CreateComment(1:CreateCommentRequest req)
    user.BaseResp DeleteComment(1:DeleteCommentRequest req)
    CommentListResponse CommentList(1:CommentListRequest req)
}