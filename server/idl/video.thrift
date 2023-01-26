namespace go video

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

struct douyin_base_response{
    1: i32 status_code // Status code, 0-success, other values-failure
    2: string status_msg // Return status description
}

struct User {
    1: i64 id // User id
    2: string name // Username
    3: i64 follow_count // Total number of followings
    4: i64 follower_count // Total number of followers
    5: bool is_follow // true-followed, false-not followed
}

struct douyin_feed_request {
    1: i64 latest_time // Optional parameter, limit the latest submission timestamp of the returned video, accurate to seconds, and leave it blank to indicate the current time
    2: string token // Optional parameter, login user settings
}

struct douyin_feed_response {
    1: douyin_base_response base_resp
    2: Video video_list // Video list
    3: i64 next_time // In the video returned this time, publish the earliest time as the latest_time in the next request
}

struct douyin_publish_action_request {
    1: string token // User authentication token
    2: byte data // Video data
    3: string title // Video title
}

struct douyin_publish_action_response {
    1: douyin_base_response base_resp
}

struct douyin_publish_list_request {
    1: i64 user_id // User id
    2: string token // User authentication token
}

struct douyin_publish_list_response {
    1: douyin_base_response base_resp
    2: list<Video> video_list // List of videos posted by the user
}

service VideoService {
    douyin_feed_response Feed (1: douyin_feed_request req)
    douyin_publish_action_response PublishVideo (1: douyin_publish_action_request req)
    douyin_publish_list_response VideoList (1: douyin_publish_list_request req)
}