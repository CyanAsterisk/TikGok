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

struct douyin_get_favorite_video_id_list_request {
    1: required i64 user_id,
}

struct douyin_get_favorite_video_id_list_response {
    1: required base.douyin_base_response base_resp,
    2: required list<i64> video_id_list,                // List of videos posted by the user
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

struct douyin_get_comment_list_request {
    1: required i64 video_id, // Video Id
}

struct douyin_get_comment_list_response {
    1: required base.douyin_base_response base_resp,
    2: required list<base.Comment> comment_list,     // List of comments
}
struct douyin_get_interact_info_request{
    1: required i64 video_id, // Video Id
    2: required i64 viewer_id, // viewer_id,
}

struct douyin_get_interact_info_response {
    1: required base.douyin_base_response base_resp,
    2: required base.InteractInfo interact_info,
}

struct douyin_batch_get_interact_info_request{
    1: required list<i64> video_id_list, // Video Id list.
    2: required i64 viewer_id, // viewer_id,
}

struct douyin_batch_get_interact_info_response {
    1: required base.douyin_base_response base_resp,
    2: required list<base.InteractInfo> interact_info_list,
}

service InteractionServer {
    douyin_favorite_action_response Favorite(1: douyin_favorite_action_request req),
    douyin_get_favorite_video_id_list_response GetFavoriteVideoIdList(1: douyin_get_favorite_video_id_list_request req),
    douyin_comment_action_response Comment(1: douyin_comment_action_request req),
    douyin_get_comment_list_response GetCommentList(1: douyin_get_comment_list_request req),
    douyin_get_interact_info_response GetInteractInfo (1: douyin_get_interact_info_request req),
    douyin_batch_get_interact_info_response BatchGetInteractInfo (1: douyin_batch_get_interact_info_request req),
}