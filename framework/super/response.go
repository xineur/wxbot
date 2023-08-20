package super

const (
	HEART_BEAT              = 5005 // 检测心跳
	RECV_TXT_MSG            = 1    // 接受到好友消息文字
	RECV_PIC_MSG            = 3    // 接受到好友消息图片
	USER_LIST               = 5000 // 发送获取好友列表
	GET_USER_LIST_SUCCESS   = 5000 // 获取到好友列表
	GET_USER_LIST_FAIL      = 5002
	TXT_MSG                 = 555  // 发送文本消息
	PIC_MSG                 = 500  // 发送图片消息
	AT_MSG                  = 550  // 发送@消息
	CHAT_ROOM_MEMBER        = 5010 // 获取群好友列表及回调
	CHAT_ROOM_MEMBER_NICK   = 5020 // 获取群成员 nickname 用于@群成员
	PERSONAL_INFO           = 6500 // 获取微信个人信息及回调
	DEBUG_SWITCH            = 6000 // debug模式
	PERSONAL_DETAIL         = 6550 // 获取好友详细信息
	DESTROY_ALL             = 9999
	NEW_FRIEND_REQUEST      = 37    // 微信好友请求消息
	AGREE_TO_FRIEND_REQUEST = 10000 // 同意微信好友请求消息
)

// 通用的发送数据格式
type SendFormat struct {
	ID      string `json:"id"`
	Type    int    `json:"type"`
	Content string `json:"content"`
	WxID    string `json:"wxid"`
}

type SendGroupMember struct {
	ID   string `json:"id"`
	Type int    `json:"type"`
	Room string `json:"room"`
}

type PicMsg struct {
	ID   string `json:"id"`
	WxID string `json:"wxid"`
	Path string `json:"path"`
	Type int    `json:"type"`
}

type Para struct {
	ID       string `json:"id"`
	Type     int    `json:"type"`
	WxID     string `json:"wxid"`
	RoomID   string `json:"roomid"`
	Content  string `json:"content"`
	Nickname string `json:"nickname"`
	Ext      string `json:"ext"`
}

type RequestBody struct {
	Para Para `json:"para"`
}

// 当前用户信息
// "content":"{\"wx_code\":\"lzj_mts\",\"wx_id\":\"wxid_jdqjbiyk0pa022\",\"wx_name\":\"一只bot\"}","id":"16894905698650043001gS3f","receiver":"CLIENT","sender":"SERVER","srvid":1,"status":"SUCCSESSED","time":"2023-07-16 06:56:09","type":6500
type ContentInfo struct {
	WxCode string `json:"wx_code"`
	WxID   string `json:"wx_id"`
	WxName string `json:"wx_name"`
}

// 用户信息
// {"headimg":"","name":"苏某人的日常","node":158977008,"remarks":"","wxcode":"gh_10a9a9fa9be5","wxid":"gh_10a9a9fa9be5"}
type FriendInfo struct {
	HeadImg string `json:"headimg"`
	Name    string `json:"name"`
	Node    int64  `json:"node"`
	Remarks string `json:"remarks"`
	WxCode  string `json:"wxcode"`
	WxID    string `json:"wxid"`
}

type RespData struct {
	Content  string `json:"content"`
	ID       string `json:"id"`
	Receiver string `json:"receiver"`
	Sender   string `json:"sender"`
	Srvid    int64  `json:"srvid"`
	Status   string `json:"status"`
	Time     string `json:"time"`
	Type     int64  `json:"type"`
}

// ObjectInfoResp 对象可以是好友、群、公众号
type ObjectInfoResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		Wxid                   string `json:"wxid"`
		WxNum                  string `json:"wxNum"`
		Nick                   string `json:"nick"`
		Remark                 string `json:"remark"`
		NickBrief              string `json:"nickBrief"`
		NickWhole              string `json:"nickWhole"`
		RemarkBrief            string `json:"remarkBrief"`
		RemarkWhole            string `json:"remarkWhole"`
		EnBrief                string `json:"enBrief"`
		EnWhole                string `json:"enWhole"`
		V3                     string `json:"v3"`
		V4                     string `json:"v4"`
		Sign                   string `json:"sign"`
		Country                string `json:"country"`
		Province               string `json:"province"`
		City                   string `json:"city"`
		MomentsBackgroudImgUrl string `json:"momentsBackgroudImgUrl"`
		AvatarMinUrl           string `json:"avatarMinUrl"`
		AvatarMaxUrl           string `json:"avatarMaxUrl"`
		Sex                    string `json:"sex"`
		MemberNum              int    `json:"memberNum"`
	} `json:"result"`
	Wxid      string `json:"wxid"`
	Port      int    `json:"port"`
	Pid       int    `json:"pid"`
	Flag      string `json:"flag"`
	Timestamp string `json:"timestamp"`
}

// FriendsListResp 获取好友列表响应
type FriendsListResp struct {
	Content []struct {
		HeadImg string `json:"headimg"`
		Name    string `json:"name"`
		Node    int    `json:"node"`
		Remarks string `json:"remarks"`
		WxCode  string `json:"wxcode"`
		WxId    string `json:"wxid"`
	} `json:"content"`
	ID   string `json:"id"`
	Type int    `json:"type"`
}

// GroupListResp 获取群组列表响应
type GroupListResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result []struct {
		Wxid                   string `json:"wxid"`
		WxNum                  string `json:"wxNum"`
		Nick                   string `json:"nick"`
		Remark                 string `json:"remark"`
		NickBrief              string `json:"nickBrief"`
		NickWhole              string `json:"nickWhole"`
		RemarkBrief            string `json:"remarkBrief"`
		RemarkWhole            string `json:"remarkWhole"`
		EnBrief                string `json:"enBrief"`
		EnWhole                string `json:"enWhole"`
		V3                     string `json:"v3"`
		Sign                   string `json:"sign"`
		Country                string `json:"country"`
		Province               string `json:"province"`
		City                   string `json:"city"`
		MomentsBackgroudImgUrl string `json:"momentsBackgroudImgUrl"`
		AvatarMinUrl           string `json:"avatarMinUrl"`
		AvatarMaxUrl           string `json:"avatarMaxUrl"`
		Sex                    string `json:"sex"`
		MemberNum              int    `json:"memberNum"`
	} `json:"result"`
	Wxid      string `json:"wxid"`
	Port      int    `json:"port"`
	Pid       int    `json:"pid"`
	Flag      string `json:"flag"`
	Timestamp string `json:"timestamp"`
}

// GroupMemberListResp 获取群成员列表响应
type GroupMemberListResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result []struct {
		Wxid      string `json:"wxid"`
		GroupNick string `json:"groupNick"`
	} `json:"result"`
	Wxid      string `json:"wxid"`
	Port      int    `json:"port"`
	Pid       int    `json:"pid"`
	Flag      string `json:"flag"`
	Timestamp string `json:"timestamp"`
}

// SubscriptionListResp 获取订阅号列表响应
type SubscriptionListResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result []struct {
		Wxid                   string `json:"wxid"`
		WxNum                  string `json:"wxNum"`
		Nick                   string `json:"nick"`
		Remark                 string `json:"remark"`
		NickBrief              string `json:"nickBrief"`
		NickWhole              string `json:"nickWhole"`
		RemarkBrief            string `json:"remarkBrief"`
		RemarkWhole            string `json:"remarkWhole"`
		EnBrief                string `json:"enBrief"`
		EnWhole                string `json:"enWhole"`
		V3                     string `json:"v3"`
		Sign                   string `json:"sign"`
		Country                string `json:"country"`
		Province               string `json:"province"`
		City                   string `json:"city"`
		MomentsBackgroudImgUrl string `json:"momentsBackgroudImgUrl"`
		AvatarMinUrl           string `json:"avatarMinUrl"`
		AvatarMaxUrl           string `json:"avatarMaxUrl"`
		Sex                    string `json:"sex"`
		MemberNum              int    `json:"memberNum"`
	} `json:"result"`
	Wxid      string `json:"wxid"`
	Port      int    `json:"port"`
	Pid       int    `json:"pid"`
	Flag      string `json:"flag"`
	Timestamp string `json:"timestamp"`
}
