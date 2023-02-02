namespace go sociality

include "base.thrift"

struct douyin_relation_action_request {
    1: required i64 user_id,    // User Id
    2: required i64 to_user_id, // The other party's user id
    3: required i8 action_type, // 1-Follow, 2-Unfollow
}

struct douyin_relation_action_response {
    1: required base.douyin_base_response base_resp,
}

struct douyin_relation_follow_list_request {
    1: required i64 viewer_id, // User id of viewer,set to zero when unclear
    2: required i64 owner_id,  // User id of owner.
}

struct douyin_relation_follow_list_response {
    1: required base.douyin_base_response base_resp,
    2: required list<base.User> user_list,           // List of user information
}

struct douyin_relation_follower_list_request {
    1: required i64 viewer_id, // User id of viewer,set to zero when unclear
    2: required i64 owner_id,  // User id of owner.
}

struct douyin_relation_follower_list_response {
    1: required base.douyin_base_response base_resp,
    2: required list<base.User> user_list,           // List of user information
}

struct douyin_relation_friend_list_request {
    1: required i64 viewer_id, // User id of viewer,set to zero when unclear
    2: required i64 owner_id,  // User id of owner.
}

struct douyin_relation_friend_list_response {
    1: required base.douyin_base_response base_resp,
    2: required list<base.FriendUser> user_list,     // List of user information
}

struct douyin_check_follow_request {
    1: required i64 user_id
    2: required i64 to_user_id
}

struct douyin_check_follow_response {
    1: required base.douyin_base_response base_resp
    2: required bool check
}

struct douyin_get_follower_count_request {
    1: required i64 user_id
}

struct douyin_get_follower_count_response {
    1: required base.douyin_base_response base_resp
    2: required i64 count
}

struct douyin_get_following_count_request {
    1: required i64 user_id
}

struct douyin_get_following_count_response {
    1: required base.douyin_base_response base_resp
    2: required i64 count
}

service SocialityService {
    douyin_relation_action_response Action(1: douyin_relation_action_request req),
    douyin_relation_follow_list_response FollowingList(1: douyin_relation_follow_list_request req),
    douyin_relation_follower_list_response FollowerList(1: douyin_relation_follower_list_request req),
    douyin_relation_friend_list_response FriendList(1: douyin_relation_friend_list_request req),
    douyin_check_follow_response CheckFollow(1: douyin_check_follow_request req),
    douyin_get_follower_count_response GetFollowerCount(1: douyin_get_follower_count_request req),
    douyin_get_following_count_response GetFollowingCount(1: douyin_get_following_count_request req),
}