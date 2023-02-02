namespace go interaction

include "base.thrift"

struct douyin_favorite_action_request {
    1: required i64 user_id,    // User Id
    2: required i64 video_id,   // Video Id
    3: required i8 action_type, // 1-like, 2-unlike
}

struct douyin_favorite_action_response {
    1: required base.douyin_base_response base_resp,
}

struct douyin_favorite_list_request {
    1: required i64 viewer_id, // User id of viewer,set to zero when unclear
    2: required i64 owner_id,  // User id of owner.
}

struct douyin_favorite_list_response {
    1: required base.douyin_base_response base_resp,
    2: required list<base.Video> video_list,         // List of videos posted by the user
}

struct douyin_favorite_count_request {
    1: required i64 video_id, // Video Id
}

struct douyin_favorite_count_response {
    1: required base.douyin_base_response base_resp,
    2: required i64 count,
}

struct douyin_check_favorite_request {
    1: required i64 user_id
    2: required i64 video_id
}

struct douyin_check_favorite_response {
    1: required base.douyin_base_response base_resp
    2: required bool check
}

struct douyin_comment_action_request {
    1: required i64 user_id,         // User Id
    2: required i64 video_id,        // Video Id
    3: required i8 action_type,      // 1-like, 2-unlike
    4: required string comment_text, // The content of the comment filled by the user, used when action_type=1
    5: required i64 comment_id,      // The comment id to be deleted is used when action_type=2
}

struct douyin_comment_action_response {
    1: required base.douyin_base_response base_resp,
    2: required base.Comment comment,                // The comment successfully returns the comment content, no need to re-pull the entire list
}

struct douyin_comment_list_request {
    1: required i64 video_id, // Video Id
}

struct douyin_comment_list_response {
    1: required base.douyin_base_response base_resp,
    2: required list<base.Comment> comment_list,     // List of comments
}

struct douyin_comment_count_request {
    1: required i64 video_id, // Video Id
}

struct douyin_comment_count_response {
    1: required base.douyin_base_response base_resp,
    2: required i64 count,
}

service InteractionServer {
    douyin_favorite_action_response Favorite(1: douyin_favorite_action_request req),
    douyin_favorite_list_response FavoriteList(1: douyin_favorite_list_request req),
    douyin_favorite_count_response FavoriteCount(1: douyin_favorite_count_request req),
    douyin_comment_action_response Comment(1: douyin_comment_action_request req),
    douyin_comment_list_response CommentList(1: douyin_comment_list_request req),
    douyin_comment_count_response CommentCount(1: douyin_comment_count_request req),
    douyin_check_favorite_response CheckFavorite(1: douyin_check_favorite_request req),
}