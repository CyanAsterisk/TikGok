namespace go sociality

struct User {
    1: i64 id // User id
    2: string name // Username
    3: i64 follow_count // Total number of followings
    4: i64 follower_count // Total number of followers
    5: bool is_follow // true-followed, false-not followed
}

struct douyin_base_response{
     1: i32 status_code // Status code, 0-success, other values-failure
     2: string status_msg // Return status description
}

struct douyin_relation_action_request {
    1: i64 user_id // User Id
    2: i64 to_user_id // The other party's user id
    3: i32 action_type // 1-Follow, 2-Unfollow
}

struct douyin_relation_action_response {
   1: douyin_base_response base_resp
}

struct douyin_relation_follow_list_request {
    1: i64 user_id // User id
}

struct douyin_relation_follow_list_response {
    1: douyin_base_response base_resp
    2: list<User> user_list // List of user information
}

struct douyin_relation_follower_list_request {
    1: i64 user_id // User id
}

struct douyin_relation_follower_list_response {
    1: douyin_base_response base_resp
    2: list<User> user_list // List of user information
}

struct douyin_relation_friend_list_request {
    1: i64 user_id // User id
}

struct douyin_relation_friend_list_response {
    1: douyin_base_response base_resp
    2: list<User> user_list // List of user information
}

service SocialityService {
    douyin_relation_action_response Action(1: douyin_relation_action_request req)
    douyin_relation_follow_list_response FollowingList(1: douyin_relation_follow_list_request req)
    douyin_relation_follower_list_response FollowerList(1: douyin_relation_follower_list_request req)
    douyin_relation_friend_list_response FriendList(1: douyin_relation_friend_list_request req)
}