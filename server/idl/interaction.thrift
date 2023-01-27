namespace go interaction
include "base.thrift"

struct douyin_favorite_action_request {
    1: i64 user_id // User Id
    2: i64 video_id // Video Id
    3: i8 action_type // 1-like, 2-unlike
}

struct douyin_favorite_action_response {
   1: base.douyin_base_response base_resp
}

struct douyin_favorite_list_request {
    1: i64 user_id // User id
}

struct douyin_favorite_list_response {
    1: base.douyin_base_response base_resp
    2: list<base.Video> video_list // List of videos posted by the user
}

struct douyin_comment_action_request {
    1: i64 user_id // User Id
    2: i64 video_id // Video Id
    3: i8 action_type // 1-like, 2-unlike
    4: string comment_text // The content of the comment filled by the user, used when action_type=1
    5: i64 comment_id // The comment id to be deleted is used when action_type=2
}

struct douyin_comment_action_response {
    1: base.douyin_base_response base_resp
    2: base.Comment comment // The comment successfully returns the comment content, no need to re-pull the entire list
}

struct douyin_comment_list_request {
    1: i64 video_id // Video Id
}

struct douyin_comment_list_response {
    1: base.douyin_base_response base_resp
    2: list<base.Comment> comment_list // List of comments
}

service InteractionServer {
    douyin_favorite_action_response Favorite(1: douyin_favorite_action_request req)
    douyin_favorite_list_response FavoriteList(1: douyin_favorite_list_request req)
    douyin_comment_action_response Comment(1: douyin_comment_action_request req)
    douyin_comment_list_response CommentList(1: douyin_comment_list_request req)
}