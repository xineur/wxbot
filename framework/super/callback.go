package super

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

	"github.com/gorilla/websocket"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

const (
	eventGroupChat           = "EventGroupChat"           // 群聊消息事件
	eventPrivateChat         = "EventPrivateChat"         // 私聊消息事件
	eventDeviceCallback      = "EventDeviceCallback"      // 设备回调事件
	eventFriendVerify        = "EventFrieneVerify"        // 好友请求事件
	eventGroupNameChange     = "EventGroupNameChange"     // 群名称变动事件
	eventGroupMemberAdd      = "EventGroupMemberAdd"      // 群成员增加事件
	eventGroupMemberDecrease = "EventGroupMemberDecrease" // 群成员减少事件
	eventInvitedInGroup      = "EventInvitedInGroup"      // 被邀请入群事件
	eventQRCodePayment       = "EventQRcodePayment"       // 面对面收款事件
	eventDownloadFile        = "EventDownloadFile"        // 文件下载结束事件
	eventGroupEstablish      = "EventGroupEstablish"      // 创建新的群聊事件
)

type Framework struct {
	BotWxId     string          // 机器人微信ID
	BotNickname string          // 机器人微信名称
	ApiUrl      string          // http api地址
	ApiToken    string          // http api鉴权token
	SocketConn  *websocket.Conn // socket 实例
}

var contentInfo ContentInfo

func New(botWxId, botNickname, apiUrl, apiToken string, socketConn *websocket.Conn, port uint) *Framework {
	go SocketCallback(apiUrl, port)
	_, message, err := socketConn.ReadMessage()
	if err != nil {
		log.Printf("err")
	}
	log.Printf("message: %s", message)
	return &Framework{
		BotWxId:     botWxId,
		BotNickname: botNickname,
		ApiUrl:      apiUrl,
		ApiToken:    apiToken,
		SocketConn:  socketConn,
	}
}

func SocketCallback(apiUrl string, port uint) {
	url := strings.Replace(apiUrl, "http", "ws", 1)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Failed to create new WebSocket connection: %s", err.Error())
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("[Super] 接受socket消息失败: %s", err.Error())
		}

		typeValue := gjson.Get(string(message), "type").Int()
		typeInt := int(typeValue)

		switch typeInt {
		case RECV_TXT_MSG, RECV_PIC_MSG, NEW_FRIEND_REQUEST, AGREE_TO_FRIEND_REQUEST:
			url := fmt.Sprintf("%s:%d/wxbot/callback", "http://127.0.0.1", port)
			params := gjson.Parse(string(message))
			if err := NewRequest().Post(url).SetBody(params.Value()).Do().Err; err != nil {
				log.Errorf("[Super] 发送消息失败: %v", err.Error())
			}
			break
		}
	}
}

func (f *Framework) Callback(ctx *gin.Context, handler func(*robot.Event, robot.IFramework)) {
	recv, err := ctx.GetRawData()
	if err != nil {
		log.Errorf("[Super] 接收回调错误, error: %v", err)
		return
	}
	handler(buildEvent(string(recv), f), f)
}

func buildEvent(resp string, f *Framework) *robot.Event {
	var event robot.Event
	switch gjson.Get(resp, "type").Int() {
	// RECV_TXT_MSG, RECV_PIC_MSG, NEW_FRIEND_REQUEST, AGREE_TO_FRIEND_REQUEST
	case PERSONAL_INFO:
		contentInfo.WxCode = gjson.Get(resp, "content.wx_code").String()
		contentInfo.WxID = gjson.Get(resp, "content.wx_id").String()
		contentInfo.WxName = gjson.Get(resp, "content.wx_name").String()
	case RECV_TXT_MSG:
		WxID := gjson.Get(resp, "wxid").String()
		// 公总号消息 不做处理
		if strings.Contains(WxID, "gh_") {
			break
		}

		if strings.Contains(WxID, "@chatroom") {
			// 群聊消息
			content := gjson.Get(resp, "content").String()
			// TODO: 如果用户名被用到, 尝试从缓存中获取
			event = robot.Event{
				Type:          robot.EventPrivateChat,
				FromGroup:     WxID,
				FromGroupName: "",
				FromWxId:      gjson.Get(resp, "id1").String(),
				FromName:      "",

				FromUniqueID:   WxID,
				FromUniqueName: "",
				Message: &robot.Message{
					Id:      gjson.Get(resp, "id").String(),
					Type:    gjson.Get(resp, "type").Int(),
					Content: extractTextAfterBot(content),
				},
			}

			if strings.Contains(content, fmt.Sprintf("@%s", f.BotNickname)) {
				event.IsAtMe = true
			}
			for _, data := range robot.GetBot().Groups() {
				if data.WxId == event.FromGroup {
					event.FromGroupName = data.Nick
					event.FromUniqueName = data.Nick
					break
				}
			}

			log.Printf("1%s1: 1%s1 -- %d", extractTextAfterBot(content), f.BotNickname, strings.Contains(content, fmt.Sprintf("@%s", contentInfo.WxName)))
		} else {
			// 个人消息
			// TODO: 如果用户名被用到, 尝试从缓存中获取
			event = robot.Event{
				Type:           robot.EventPrivateChat,
				IsAtMe:         true,
				FromUniqueID:   WxID,
				FromUniqueName: "",
				FromWxId:       WxID,
				FromName:       "",
				Message: &robot.Message{
					Id:      gjson.Get(resp, "id").String(),
					Type:    gjson.Get(resp, "type").Int(),
					Content: gjson.Get(resp, "content").String(),
				},
			}
		}
	case RECV_PIC_MSG:
		WxID := gjson.Get(resp, "wxid").String()
		// 公总号消息 不做处理
		if strings.Contains(WxID, "gh_") {
			break
		}
		log.Printf("接受到图片消息: %s", resp)
	}
	event.RobotWxId = contentInfo.WxID
	event.RawMessage = resp
	return &event
}

func extractTextAfterBot(content string) string {
	pattern := fmt.Sprintf("^\\s*?@.* (.*?)$")
	// 编译正则表达式
	reg := regexp.MustCompile(pattern)

	// 查找匹配的文本
	match := reg.FindStringSubmatch(content)

	if len(match) >= 2 {
		// 第一个捕获组即为我们需要的文本
		text := match[1]
		return text
	}

	return "" // 如果找不到对应的内容，返回空字符串
}
