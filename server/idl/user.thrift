namespace go user
include "base.thrift"


struct douyin_user_register_request {
    1: string username // Username, up to 32 characters
    2: string password // Password, up to 32 characters
}

struct douyin_user_register_response {
    1: base.douyin_base_response base_resp
    2: i64 user_id // User id
    3: string token // User authentication token
}

struct douyin_user_login_request {
    1: string username // Username, up to 32 characters
    2: string password // Password, up to 32 characters
}

struct douyin_user_login_response {
    1: base.douyin_base_response base_resp
    2: i64 user_id // User id
    3: string token // User authentication token
}

struct douyin_user_request {
    1: i64 user_id // User id
}

struct douyin_user_response {
    1: base.douyin_base_response base_resp
    2: base.User user // User Information
}

service UserService {
    douyin_user_register_response Register(1: douyin_user_register_request req)
    douyin_user_login_response Login(1: douyin_user_login_request req)
    douyin_user_response GetUserInfo(1: douyin_user_request req)
}