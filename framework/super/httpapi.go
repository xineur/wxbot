package super

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/gorilla/websocket"

	"github.com/tidwall/gjson"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

// private
// 生成指定长度的随机字符序列
func generateRandomChars(length int) string {
	rand.Seed(time.Now().UnixNano())

	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}

	return string(result)
}

func getId() string {
	// 生成当前时间的时间戳字符串
	timestamp := time.Now().UnixNano()
	timestampStr := fmt.Sprintf("%d", timestamp)

	// 生成随机字符序列
	randomChars := generateRandomChars(5) // 生成长度为5的随机字符序列

	// 合并时间戳和随机字符序列
	result := timestampStr + randomChars
	return result
}

// private end

func (f *Framework) msgFormat(msg string) string {
	buff := bytes.NewBuffer(make([]byte, 0, len(msg)*2))
	for _, r := range msg {
		if unicode.Is(unicode.Han, r) || unicode.IsLetter(r) {
			buff.WriteString(string(r))
			continue
		}
		switch utf8.RuneLen(r) {
		case 2, 3:
			buff.WriteString(`[emoji=\u`)
			buff.WriteString(fmt.Sprintf("%04x", r) + `]`)
		case 4:
			r1, r2 := utf16.EncodeRune(r)
			buff.WriteString(`[emoji=\u`)
			buff.WriteString(strconv.FormatInt(int64(r1), 16))
			buff.WriteString(`\u`)
			buff.WriteString(strconv.FormatInt(int64(r2), 16) + `]`)
		default:
			buff.WriteString(string(r))
		}
	}
	return buff.String()
}

func getDataBySocket(socketConn *websocket.Conn, params interface{}, handleFunc func(dataResp *RespData) (bool, error)) (bool, error) {
	var dataResp RespData
	jsonParams, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}

	err = socketConn.WriteMessage(websocket.TextMessage, []byte(jsonParams))
	if err != nil {
		log.Fatal(err)
	}

	for {
		// 接收服务器的响应
		msgType, message, err := socketConn.ReadMessage()
		if err != nil {
			return false, fmt.Errorf("读取消息出错: %s", err.Error())
		}
		log.Debug("接收到消息: %s", string(message))
		if msgType == websocket.TextMessage {

			messageStr := string(message)
			dataResp.Content = gjson.Get(messageStr, "content").String()
			dataResp.Type = gjson.Get(messageStr, "type").Int()
			dataResp.ID = gjson.Get(messageStr, "id").String()
			dataResp.Receiver = gjson.Get(messageStr, "receiver").String()
			dataResp.Sender = gjson.Get(messageStr, "sender").String()
			dataResp.Srvid = gjson.Get(messageStr, "srvid").Int()
			dataResp.Status = gjson.Get(messageStr, "status").String()
			dataResp.Time = gjson.Get(messageStr, "time").String()
		}

		// 调用处理函数进行处理
		if ok, err := handleFunc(&dataResp); err != nil {
			return false, fmt.Errorf("处理数据出错: %s", err.Error())
		} else if ok {
			// 处理成功，继续进行其他操作
			return true, nil
		}
	}
}

// TODO)) Socket 方法已废弃
func (f *Framework) GetRobotInfo() (*robot.User, error) {

	// params := SendFormat{
	// 	ID:      getId(),
	// 	Type:    PERSONAL_INFO,
	// 	Content: "personal info",
	// 	WxID:    "ROOT",
	// }

	var contentInfo ContentInfo
	contentInfo.WxID = f.BotWxId
	contentInfo.WxCode = f.BotWxId
	contentInfo.WxName = ""

	// _, err := getDataBySocket(f.SocketConn, params, func(resp *RespData) (bool, error) {
	// 	log.Debug("获取当前用户信息: %s", resp.Content)
	// 	if resp.Type != PERSONAL_INFO || resp.Status != "SUCCSESSED" {
	// 		return false, nil
	// 	}

	// 	contentInfo.WxCode = gjson.Get(resp.Content, "wx_code").String()
	// 	contentInfo.WxID = gjson.Get(resp.Content, "wx_id").String()
	// 	contentInfo.WxName = gjson.Get(resp.Content, "wx_name").String()

	// 	// 进行处理逻辑
	// 	return true, nil
	// })

	// if err != nil {
	// 	log.Fatal("接收获取当前用户信息失败: %s", err)
	// }

	return &robot.User{
		WxId:         contentInfo.WxID,
		WxNum:        contentInfo.WxCode,
		Nick:         contentInfo.WxName,
		Country:      "",
		Province:     "",
		City:         "",
		AvatarMinUrl: "",
		AvatarMaxUrl: "",
	}, nil
}

// TODO: 未完成
func (f *Framework) GetMemePictures(msg *robot.Message) string {
	log.Printf("GetMemePictures %s", msg.Content)
	log.Printf("[Super] 不支持获取图片")
	return ("[Super] 不支持获取图片")
}

func (f *Framework) SendText(toWxId, text string) error {
	apiUrl := fmt.Sprintf("%s/api/sendtxtmsg", f.ApiUrl)

	params := map[string]interface{}{
		"para": map[string]interface{}{
			"id":       getId(),
			"type":     TXT_MSG,
			"wxid":     toWxId,
			"roomid":   "null",
			"content":  text,
			"nickname": "NULL",
			"ext":      "null",
		},
	}

	if err := NewRequest().Get(apiUrl).SetBody(params).SetHeader("Content-Type", "application/json").Do().Err; err != nil {
		log.Errorf("[Super] SendText error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendTextAndAt(toGroupWxId, toWxId, toWxName, text string) error {

	apiUrl := fmt.Sprintf("%s/api/sendtxtmsg", f.ApiUrl)

	params := map[string]interface{}{
		"para": map[string]interface{}{
			"id":       getId(),
			"type":     AT_MSG,
			"content":  text,
			"wxid":     toWxId,
			"roomid":   toGroupWxId,
			"nickname": toWxName,
			"ext":      "null",
		},
	}

	if err := NewRequest().Get(apiUrl).SetBody(params).SetHeader("Content-Type", "application/json").Do().Err; err != nil {
		log.Errorf("[Super] SendTextAndAt error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendImage(toWxId, path string) error {
	apiUrl := fmt.Sprintf("%s/api/sendpic", f.ApiUrl)

	params := map[string]interface{}{
		"para": map[string]interface{}{
			"id":       getId(),
			"type":     PIC_MSG,
			"wxid":     toWxId,
			"content":  "C:\\Users\\86176\\Downloads\\preview.png",
			"ext":      "null",
			"nickname": "null",
		},
	}

	if err := NewRequest().Get(apiUrl).SetBody(params).SetHeader("Content-Type", "application/json").Do().Err; err != nil {
		log.Errorf("[Super] SendImage error: %v", err.Error())
		return err
	}
	return nil
}

func (f *Framework) SendShareLink(toWxId, title, desc, imageUrl, jumpUrl string) error {
	log.Printf("[Super] 不支持发送链接")
	return errors.New("[Super] 不支持发送链接")
}

// TODO: 待完成
func (f *Framework) SendFile(toWxId, path string) error {
	log.Printf("[Super] 不支持发送文件")
	return errors.New("[Super] 不支持发送文件")
}

func (f *Framework) SendVideo(toWxId, path string) error {
	log.Printf("[Super] 不支持发送视频")
	return errors.New("[Super] 不支持发送视频")
}

func (f *Framework) SendEmoji(toWxId, path string) error {
	log.Printf("[Super] 不支持发送emoji")
	return errors.New("[Super] 不支持发送emoji")
}

func (f *Framework) SendMusic(toWxId, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	log.Printf("[Super] 不支持发送音乐")
	return errors.New("[Super] 不支持发送音乐")
}

func (f *Framework) SendMiniProgram(toWxId, ghId, title, content, imagePath, jumpPath string) error {
	log.Printf("[Super] 不支持发送小程序")
	return errors.New("[Super] 不支持发送小程序")
}

func (f *Framework) SendMessageRecord(toWxId, title string, dataList []map[string]interface{}) error {
	log.Printf("[Super] 不支持发送消息记录")
	return errors.New("[Super] 不支持发送消息记录")
}

func (f *Framework) SendMessageRecordXML(toWxId, xmlStr string) error {
	log.Printf("[Super] 不支持发送消息记录xml")
	return errors.New("[Super] 不支持发送消息记录xml")
}

func (f *Framework) SendFavorites(toWxId, favoritesId string) error {
	log.Printf("[Super] 不支持发送收藏夹")
	return errors.New("[Super] 不支持发送收藏夹")
}

func (f *Framework) SendXML(toWxId, xmlStr string) error {
	log.Printf("[Super] 不支持发送xml")
	return errors.New("[Super] 不支持发送xml")
}

func (f *Framework) SendBusinessCard(toWxId, targetWxId string) error {
	log.Printf("[Super] 不支持发送名片")
	return errors.New("[Super] 不支持发送名片")
}

func (f *Framework) AgreeFriendVerify(v3, v4, scene string) error {
	log.Printf("[Super] 不支持同意好友请求")
	return errors.New("[Super] 不支持同意好友请求")
}

func (f *Framework) InviteIntoGroup(groupWxId, wxId string, typ int) error {
	log.Printf("[Super] 不支持拉好友入群聊")
	return errors.New("[Super] 不支持拉好友入群聊")
}

// TODO: 未完成
func (f *Framework) GetObjectInfo(wxId string) (*robot.User, error) {
	params := SendFormat{
		ID:      getId(),
		Type:    PERSONAL_DETAIL,
		Content: "op:personal info",
		WxID:    wxId,
	}

	var contentInfo ContentInfo

	_, err := getDataBySocket(f.SocketConn, params, func(dataResp *RespData) (bool, error) {
		// if resp, ok := dataResp.(*RobotInfoResp); ok {
		// 	// 处理 *RobotInfoResp 类型的情况
		// 	if (resp.Type == PERSONAL_INFO) && (resp.Status == "SUCCSESSED") {
		// 		err := json.Unmarshal(resp.Content, &contentInfo)
		// 		if err != nil {
		// 			return false, fmt.Errorf("解析content字段出错: %s", err.Error())
		// 		}
		// 		// 进行处理逻辑
		// 		return true, nil
		// 	}
		// }
		return false, nil
	})

	if err != nil {
		log.Fatal("获取对象信息: %s", err)
	}
	return &robot.User{
		WxId:                    contentInfo.WxID,
		WxNum:                   contentInfo.WxCode,
		Nick:                    contentInfo.WxName,
		Remark:                  contentInfo.WxName,
		NickBrief:               contentInfo.WxName,
		NickWhole:               contentInfo.WxName,
		RemarkBrief:             contentInfo.WxName,
		RemarkWhole:             contentInfo.WxName,
		EnBrief:                 contentInfo.WxName,
		EnWhole:                 contentInfo.WxName,
		V3:                      contentInfo.WxName,
		V4:                      contentInfo.WxName,
		Sign:                    contentInfo.WxName,
		Country:                 contentInfo.WxName,
		Province:                contentInfo.WxName,
		City:                    contentInfo.WxName,
		MomentsBackgroundImgUrl: contentInfo.WxName,
		AvatarMinUrl:            contentInfo.WxName,
		AvatarMaxUrl:            contentInfo.WxName,
		Sex:                     contentInfo.WxName,
		MemberNum:               0,
	}, nil
}

func (f *Framework) GetFriends(isRefresh bool) ([]*robot.User, error) {
	apiUrl := fmt.Sprintf("%s/api/getcontactlist", f.ApiUrl)

	params := map[string]interface{}{
		"para": map[string]interface{}{
			"id":       getId(),
			"type":     USER_LIST,
			"wxid":     "null",
			"roomid":   "null",
			"content":  "null",
			"ext":      "null",
			"nickname": "null",
		},
	}

	var dataResp FriendsListResp
	if err := NewRequest().Get(apiUrl).SetBody(params).SetHeader("Content-Type", "application/json").SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[Super] GetFriends error: %v", err.Error())
		return nil, err
	}

	var friendsInfoList []*robot.User
	for _, res := range dataResp.Content {
		if strings.Contains(res.WxId, "@chatroom") || strings.Contains(res.WxId, "gh_") {
			continue
		}
		friendsInfoList = append(friendsInfoList, &robot.User{
			WxId:                    res.WxId,
			WxNum:                   res.WxCode,
			Nick:                    res.Name,
			Remark:                  res.Remarks,
			NickBrief:               res.Name,
			NickWhole:               res.Name,
			RemarkBrief:             "",
			RemarkWhole:             "",
			EnBrief:                 "",
			EnWhole:                 "",
			V3:                      "",
			Sign:                    "",
			Country:                 "",
			Province:                "",
			City:                    "",
			MomentsBackgroundImgUrl: "",
			AvatarMinUrl:            res.HeadImg,
			AvatarMaxUrl:            res.HeadImg,
			Sex:                     "",
			MemberNum:               0,
		})
	}

	// 过滤系统用户
	var SystemUserWxId = map[string]struct{}{"medianote": {}, "newsapp": {}, "fmessage": {}, "floatbottle": {}}
	var filteredFriendInfo []*robot.User
	for i := range friendsInfoList {
		if _, ok := SystemUserWxId[friendsInfoList[i].WxId]; !ok {
			filteredFriendInfo = append(filteredFriendInfo, friendsInfoList[i])
		}
	}

	return filteredFriendInfo, nil
}

func (f *Framework) GetGroups(isRefresh bool) ([]*robot.User, error) {
	apiUrl := fmt.Sprintf("%s/api/getcontactlist", f.ApiUrl)

	params := map[string]interface{}{
		"para": map[string]interface{}{
			"id":       getId(),
			"type":     USER_LIST,
			"wxid":     "null",
			"roomid":   "null",
			"content":  "null",
			"ext":      "null",
			"nickname": "null",
		},
	}

	var dataResp FriendsListResp
	if err := NewRequest().Get(apiUrl).SetBody(params).SetHeader("Content-Type", "application/json").SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[Super] GetGroups error: %v", err.Error())
		return nil, err
	}

	var groupInfoList []*robot.User
	for _, res := range dataResp.Content {
		if strings.Contains(res.WxId, "@chatroom") {

			groupInfoList = append(groupInfoList, &robot.User{
				WxId:                    res.WxId,
				WxNum:                   res.WxCode,
				Nick:                    res.Name,
				Remark:                  res.Remarks,
				NickBrief:               res.Name,
				NickWhole:               res.Name,
				RemarkBrief:             "",
				RemarkWhole:             "",
				EnBrief:                 "",
				EnWhole:                 "",
				V3:                      "",
				Sign:                    "",
				Country:                 "",
				Province:                "",
				City:                    "",
				MomentsBackgroundImgUrl: "",
				AvatarMinUrl:            res.HeadImg,
				AvatarMaxUrl:            res.HeadImg,
				Sex:                     "",
				MemberNum:               0,
			})

		}
	}

	return groupInfoList, nil
}

// TODO: 未完成
func (f *Framework) GetGroupMembers(groupWxId string, isRefresh bool) ([]*robot.User, error) {
	params := map[string]string{
		"id":   getId(),
		"type": strconv.Itoa(CHAT_ROOM_MEMBER_NICK),
		"room": groupWxId,
	}

	apiUrl := fmt.Sprintf("%s/api/getmembernick", f.ApiUrl)
	var dataResp GroupMemberListResp
	if err := NewRequest().Get(apiUrl).SetQueryParams(params).SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[Super] GetGroupMembers error: %v", err.Error())
		return nil, err
	}
	var groupMemberInfoList []*robot.User
	for _, res := range dataResp.Result {
		log.Printf("[Super] GetGroupMembers response", res)
		groupMemberInfoList = append(groupMemberInfoList, &robot.User{
			WxId: res.Wxid,
			Nick: res.GroupNick,
		})
	}
	return groupMemberInfoList, nil
}

func (f *Framework) GetMPs(isRefresh bool) ([]*robot.User, error) {
	apiUrl := fmt.Sprintf("%s/api/getcontactlist", f.ApiUrl)

	params := map[string]interface{}{
		"para": map[string]interface{}{
			"id":       getId(),
			"type":     USER_LIST,
			"wxid":     "null",
			"roomid":   "null",
			"content":  "null",
			"ext":      "null",
			"nickname": "null",
		},
	}

	var dataResp FriendsListResp
	if err := NewRequest().Get(apiUrl).SetBody(params).SetHeader("Content-Type", "application/json").SetSuccessResult(&dataResp).Do().Err; err != nil {
		log.Errorf("[Super] GetMPs error: %v", err.Error())
		return nil, err
	}

	var subscriptionInfoList []*robot.User
	for _, res := range dataResp.Content {
		if strings.Contains(res.WxId, "gh_") {
			subscriptionInfoList = append(subscriptionInfoList, &robot.User{
				WxId:                    res.WxId,
				WxNum:                   res.WxCode,
				Nick:                    res.Name,
				Remark:                  res.Remarks,
				NickBrief:               res.Name,
				NickWhole:               res.Name,
				RemarkBrief:             "",
				RemarkWhole:             "",
				EnBrief:                 "",
				EnWhole:                 "",
				V3:                      "",
				Sign:                    "",
				Country:                 "",
				Province:                "",
				City:                    "",
				MomentsBackgroundImgUrl: "",
				AvatarMinUrl:            res.HeadImg,
				AvatarMaxUrl:            res.HeadImg,
				Sex:                     "",
				MemberNum:               0,
			})

		}
	}

	return subscriptionInfoList, nil
}
