namespace go interaction

struct User {
    1: i64 id // User id
    2: string name // Username
    3: i64 follow_count // Total number of followings
    4: i64 follower_count // Total number of followers
    5: bool is_follow // true-followed, false-not followed
}

struct Comment {
    1: i64 id // Video comment id
    2: User user // Comment user information
    3: string content // Comment's content
    4: string create_date // Comment release date, format mm-dd
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

struct douyin_base_response{
     1: i32 status_code // Status code, 0-success, other values-failure
     2: string status_msg // Return status description
}

struct douyin_favorite_action_request {
    1: i64 user_id // User Id
    2: i64 video_id // Video Id
    3: i8 action_type // 1-like, 2-unlike
}

struct douyin_favorite_action_response {
   1: douyin_base_response base_resp
}

struct douyin_favorite_list_request {
    1: i64 user_id // User id
}

struct douyin_favorite_list_response {
    1: douyin_base_response base_resp
    2: list<Video> video_list // List of videos posted by the user
}

struct douyin_comment_action_request {
    1: i64 user_id // User Id
    2: i64 video_id // Video Id
    3: i8 action_type // 1-like, 2-unlike
    4: string comment_text // The content of the comment filled by the user, used when action_type=1
    5: i64 comment_id // The comment id to be deleted is used when action_type=2
}

struct douyin_comment_action_response {
    1: douyin_base_response base_resp
    2: Comment comment // The comment successfully returns the comment content, no need to re-pull the entire list
}

struct douyin_comment_list_request {
    1: i64 video_id // Video Id
}

struct douyin_comment_list_response {
    1: douyin_base_response base_resp
    2: list<Comment> comment_list // List of comments
}

service InteractionServer {
    douyin_favorite_action_response Favorite(1: douyin_favorite_action_request req)
    douyin_favorite_list_response FavoriteList(1: douyin_favorite_list_request req)
    douyin_comment_action_response Comment(1: douyin_comment_action_request req)
    douyin_comment_list_response CommentList(1: douyin_comment_list_request req)
}