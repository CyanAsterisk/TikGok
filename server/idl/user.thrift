namespace go user

include "base.thrift"

struct douyin_user_register_request {
    1: required string username, // Username, up to 32 characters
    2: required string password, // Password, up to 32 characters
}

struct douyin_user_register_response {
    1: required base.douyin_base_response base_resp,
    2: required i64 user_id,                         // User id
    3: required string token,                        // User authentication token
}

struct douyin_user_login_request {
    1: required string username, // Username, up to 32 characters
    2: required string password, // Password, up to 32 characters
}

struct douyin_user_login_response {
    1: required base.douyin_base_response base_resp,
    2: required i64 user_id,                         // User id
    3: required string token,                        // User authentication token
}

struct douyin_user_request {
    1: required i64 viewer_id, // User id of viewer,set to zero when unclear
    2: required i64 owner_id,  // User id of owner.
}

struct douyin_user_response {
    1: required base.douyin_base_response base_resp,
    2: required base.User user,                      // User Information
}

service UserService {
    douyin_user_register_response Register(1: douyin_user_register_request req),
    douyin_user_login_response Login(1: douyin_user_login_request req),
    douyin_user_response GetUserInfo(1: douyin_user_request req),
}