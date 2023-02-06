namespace go base

struct douyin_base_response {
    1: i32 status_code // Status code, 0-success, other values-failure
    2: string status_msg // Return status description
}

struct Comment {
    1: i64 id // Video comment id
    2: User user // Comment user information
    3: string content // Comment's content
    4: string create_date // Comment release date, format mm-dd
}

struct User {
    1: i64 id // User id
    2: string name // Username
    3: i64 follow_count // Total number of followings
    4: i64 follower_count // Total number of followers
    5: bool is_follow // true-followed, false-not followed
}

struct SocialInfo{
     1: i64 follow_count // Total number of followings
     2: i64 follower_count // Total number of followers
     3: bool is_follow // true-followed, false-not followed
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

struct InteractInfo{
    1: i64 favorite_count // Total number of likes for the video
    2: i64 comment_count // Total number of comments on the video
    3: bool is_favorite // true-liked, false-not liked
}

struct FriendUser {
    1: i64 id // User id
    2: string name // Username
    3: i64 follow_count // Total number of followings
    4: i64 follower_count // Total number of followers
    5: bool is_follow // true-followed, false-not followed
    7: string message // Latest chat messages with this friend
    8: i64 msgType // message type, 0 => the message received by the current requesting user, 1 => the message sent by the current requesting user
}

struct Message {
    1: i64 id // Message id
    2: i64 to_user_id // The id of the recipient of the message
    3: i64 from_user_id // The id of the sender of the message
    4: string content // Message content
    5: string create_time // Message creation time
}

struct LatestMsg{
     1: string message // Latest chat messages with this friend
     2: i64 msgType // message type, 0 => the message received by the current requesting user, 1 => the message sent by the current requesting user
}