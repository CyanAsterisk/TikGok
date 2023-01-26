namespace go user

struct User {
    1: i64 id // User id
    2: string name // Username
    3: i64 follow_count // Total number of followings
    4: i64 follower_count // Total number of followers
    5: bool is_follow // true-followed, false-not followed
}

struct douyin_base_response {
    1: i32 status_code // Status code, 0-success, other values-failure
    2: string status_msg // Return status description
}


struct douyin_user_register_request {
    1: string username // Username, up to 32 characters
    2: string password // Password, up to 32 characters
}

struct douyin_user_register_response {
    1: douyin_base_response base_resp
    2: i64 user_id // User id
    3: string token // User authentication token
}

struct douyin_user_login_request {
    1: string username // Username, up to 32 characters
    2: string password // Password, up to 32 characters
}

struct douyin_user_login_response {
    1: douyin_base_response base_resp
    2: i64 user_id // User id
    3: string token // User authentication token
}

struct douyin_user_request {
    1: i64 user_id // User id
    2: string token // User authentication token
}

struct douyin_user_response {
    1: douyin_base_response base_resp
    2: User user // User Information
}

service UserService {
    douyin_user_register_response Register(1: douyin_user_register_request req)
    douyin_user_login_response Login(1: douyin_user_login_request req)
    douyin_user_response GetUserInfo(1: douyin_user_request req)
}