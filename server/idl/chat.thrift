namespace go chat
include "base.thrift"

struct douyin_message_get_chat_history_request {
    1: i64 user_id // User Id
    2: i64 to_user_id // The other party's user id
}

struct douyin_message_get_chat_history_response {
    1: base.douyin_base_response base_resp
    2: list<base.Message> message_list // Message list
}

struct douyin_message_action_request {
    1: i64 user_id // User Id
    2: i64 to_user_id // The other party's user id
    3: i8 action_type // 1- Send a message
    4: string content // Message content
}

struct douyin_message_action_response {
    1: base.douyin_base_response base_resp
}

struct douyin_message_get_latest_request {
    1: i64 user_id // User Id
    2: i64 to_user_id // The other party's user id
}

struct douyin_message_get_latest_response {
    1: base.douyin_base_response base_resp
    2: string message // Latest chat messages with this friend
    3: i64 msgType // message type, 0 => the message received by the current requesting user, 1 => the message sent by the current requesting user
}

struct douyin_message_batch_get_latest_request {
    1: i64 user_id // User Id
    2: list<i64> to_user_ids // The other party's user ids
}

struct douyin_message_batch_get_latest_response {
    1: base.douyin_base_response base_resp
    2: list<string> message // Latest chat messages with this friend
    3: list<i64> msgType // message type, 0 => the message received by the current requesting user, 1 => the message sent by the current requesting user
}

service ChatService {
    douyin_message_get_chat_history_response GetChatHistory(1: douyin_message_get_chat_history_request req)
    douyin_message_action_response SentMessage(1: douyin_message_action_request req)
    douyin_message_get_latest_response GetLatestMessage(1: douyin_message_get_latest_request req)
    douyin_message_batch_get_latest_response BatchGetLatestMessage(1: douyin_message_batch_get_latest_request req)
}