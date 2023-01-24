namespace go api

struct douyin_user_register_request {
    1: string username(api.query="username", api.vd="len($)>0 && len($)<33") // Username, up to 32 characters
    2: string password(api.query="password", api.vd="len($)>0 && len($)<33") // Password, up to 32 characters
}

struct douyin_user_register_response {
    1: i32 status_code // Status code, 0-success, other values-failure
    2: string status_msg // Return status description
    3: i64 user_id // User id
    4: string token // User authentication token
}

struct douyin_user_login_request {
    1: string username(api.query="username", api.vd="len($)>0 && len($)<33") // Username, up to 32 characters
    2: string password(api.query="password", api.vd="len($)>0 && len($)<33") // Password, up to 32 characters
}

struct douyin_user_login_response {
    1: i32 status_code // Status code, 0-success, other values-failure
    2: string status_msg // Return status description
    3: i64 user_id // User id
    4: string token // User authentication token
}

struct douyin_user_request {
    1: i64 user_id(api.query="user_id", api.vd="len($)>0 && len($)<20") // User id
    2: string token(api.query="token") // User authentication token
}

struct douyin_user_response {
    1: i32 status_code // Status code, 0-success, other values-failure
    2: string status_msg // Return status description
    3: User user // User Information
}

struct User {
    1: i64 id // User id
    2: string name // Username
    3: i64 follow_count // Total number of followings
    4: i64 follower_count // Total number of followers
    5: bool is_follow // true-followed, false-not followed
}

struct douyin_feed_request {
    1: i64 latest_time(api.query="latest_time") // Optional parameter, limit the latest submission timestamp of the returned video, accurate to seconds, and leave it blank to indicate the current time
    2: string token(api.query="token") // Optional parameter, login user settings
}

struct douyin_feed_response {
    1: i32 status_code // Status code, 0-success, other values-failure
    2: string status_msg // Return status description
    3: Video video_list // Video list
    4: i64 next_time // In the video returned this time, publish the earliest time as the latest_time in the next request
}

struct Video {
    1: i64 id // Video unique identifier
    2: User author // Video author information
    3: string play_url // Video play URL
    4: string cover_url // Video cover address
    5: i64 favorite_count // Total number of likes for the video
    6: i64 comment_count // Total number of comments on the video
    7: bool is_favorite // true-liked, false-not liked
    8: string title // Video title
}

struct douyin_publish_action_request {
    1: string token(api.form="token") // User authentication token
    2: byte data(api.form="data") // Video data
    3: string title(api.form="title", api.vd="len($)>0 && len($)<33") // Video title
}

struct douyin_publish_action_response {
    1: i32 status_code // Status code, 0-success, other values-failure
    2: string status_msg // Return status description
}

struct douyin_publish_list_request {
    1: i64 user_id(api.query="user_id", api.vd="len($)>0 && len($)<20") // User id
    2: string token(api.query="token") // User authentication token
}

struct douyin_publish_list_response {
    1: i32 status_code // Status code, 0-success, other values-failure
    2: string status_msg // Return status description
    3: list<Video> video_list // List of videos posted by the user
}

service ApiService {
    douyin_user_register_response Register(1: douyin_user_register_request req)(api.post="/douyin/user/register");
    douyin_user_login_response Login(1: douyin_user_login_request req)(api.post="/douyin/user/login");
    douyin_user_request GetUserInfo(1: douyin_user_request req)(api.get="/douyin/user");

    douyin_feed_response Feed (1: douyin_feed_request req)(api.get="/douyin/feed");
    douyin_publish_action_response PublishVideo (1: douyin_publish_action_request req)(api.post="/douyin/publish/action");
    douyin_publish_list_response VideoList (1: douyin_publish_list_request req)(api.get="/douyin/publish/list");
}