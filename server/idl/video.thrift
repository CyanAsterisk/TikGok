namespace go video

include "base.thrift"

struct douyin_feed_request {
    1: required i64 latest_time, // Optional parameter, limit the latest submission timestamp of the returned video, accurate to seconds, and leave it blank to indicate the current time
    2: required i64 user_id,     // User Id
}

struct douyin_feed_response {
    1: required base.douyin_base_response base_resp,
    2: required list<base.Video> video_list,         // Video list
    3: required i64 next_time,                       // In the video returned this time, publish the earliest time as the latest_time in the next request
}

struct douyin_publish_action_request {
    1: required i64 user_id,      // User Id
    2: required string play_url,  // Video play url
    3: required string cover_url, // Video cover url
    4: required string title,     // Video title
}

struct douyin_publish_action_response {
    1: required base.douyin_base_response base_resp,
}

struct douyin_publish_list_request {
    1: required i64 user_id, // User id
}

struct douyin_publish_list_response {
    1: required base.douyin_base_response base_resp,
    2: required list<base.Video> video_list,         // List of videos posted by the user
}

struct douyin_get_video_request {
    1: required i64 video_id, // Video id
}

struct douyin_get_video_response {
    1: required base.douyin_base_response base_resp,
    2: required base.Video video,
}

service VideoService {
    douyin_feed_response Feed(1: douyin_feed_request req),
    douyin_publish_action_response PublishVideo(1: douyin_publish_action_request req),
    douyin_publish_list_response VideoList(1: douyin_publish_list_request req),
    douyin_get_video_response GetVideo(1: douyin_get_video_request req),
}