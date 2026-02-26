package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
	"time"
	"watchAlert/pkg/tools"

	"github.com/zeromicro/go-zero/core/logc"
)

// WeChatMiniProgramSender 企业微信小程序发送策略
type WeChatMiniProgramSender struct {
	tokenCache *TokenCache
}

// WeChatMiniProgramConfig 企业微信小程序配置
type WeChatMiniProgramConfig struct {
	CorpId  string `json:"corpId"`  // 企业ID
	AgentId int64  `json:"agentId"` // 应用ID
	Secret  string `json:"secret"`  // 应用密钥
	ToUser  string `json:"toUser"`  // 接收者用户ID，多个用|分隔
	ToParty string `json:"toParty"` // 接收者部门ID，多个用|分隔
	ToTag   string `json:"toTag"`   // 接收者标签ID，多个用|分隔
}

// TokenCache access_token 缓存
type TokenCache struct {
	token     string
	expiresAt int64
	mu        sync.RWMutex
}

// WeChatAccessTokenResponse 获取access_token响应
type WeChatAccessTokenResponse struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// WeChatSendMessageResponse 发送消息响应
type WeChatSendMessageResponse struct {
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	InvalidUser  string `json:"invaliduser"`
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

// WeChatMiniProgramMessage 企业微信小程序消息
type WeChatMiniProgramMessage struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentId int64  `json:"agentid"`
	Text    struct {
		Content string `json:"content"`
	} `json:"text"`
	Markdown struct {
		Content string `json:"content"`
	} `json:"markdown,omitempty"`
	TextCard struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		URL         string `json:"url"`
		BtnTxt      string `json:"btntxt,omitempty"`
	} `json:"textcard,omitempty"`
	Safe int `json:"safe"` // 表示是否是保密消息，0表示可对外分享，1表示不能分享且内容显示水印，默认为0
}

func NewWeChatMiniProgramSender() SendInter {
	return &WeChatMiniProgramSender{
		tokenCache: &TokenCache{},
	}
}

func (w *WeChatMiniProgramSender) Send(params SendParams) error {
	return w.postMessage(params.Hook, params.Content)
}

func (w *WeChatMiniProgramSender) Test(params SendParams) error {
	return w.postMessage(params.Hook, WechatTestContent)
}

// postMessage 发送消息
func (w *WeChatMiniProgramSender) postMessage(hook, content string) error {
	// 解析配置
	var config WeChatMiniProgramConfig
	if err := json.Unmarshal([]byte(hook), &config); err != nil {
		return fmt.Errorf("解析企业微信小程序配置失败: %v", err)
	}

	// 获取access_token
	accessToken, err := w.getAccessToken(config.CorpId, config.Secret)
	if err != nil {
		return fmt.Errorf("获取access_token失败: %v", err)
	}

	// 解析消息内容
	var messageContent map[string]interface{}
	if err := json.Unmarshal([]byte(content), &messageContent); err != nil {
		return fmt.Errorf("解析消息内容失败: %v", err)
	}

	// 构建消息
	msgType := "text"
	if mt, ok := messageContent["msgtype"].(string); ok {
		msgType = mt
	}

	message := WeChatMiniProgramMessage{
		ToUser:  config.ToUser,
		ToParty: config.ToParty,
		ToTag:   config.ToTag,
		MsgType: msgType,
		AgentId: config.AgentId,
		Safe:    0,
	}

	switch msgType {
	case "text":
		if contentStr, ok := messageContent["content"].(string); ok {
			message.Text.Content = contentStr
		}
	case "markdown":
		if contentStr, ok := messageContent["content"].(string); ok {
			message.Markdown.Content = contentStr
		}
	case "textcard":
		if title, ok := messageContent["title"].(string); ok {
			message.TextCard.Title = title
		}
		if desc, ok := messageContent["description"].(string); ok {
			message.TextCard.Description = desc
		}
		if url, ok := messageContent["url"].(string); ok {
			message.TextCard.URL = url
		}
		if btnTxt, ok := messageContent["btntxt"].(string); ok {
			message.TextCard.BtnTxt = btnTxt
		}
	default:
		return fmt.Errorf("不支持的消息类型: %s", msgType)
	}

	// 序列化消息
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	// 发送消息
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", accessToken)
	res, err := tools.Post(nil, url, bytes.NewReader(messageJSON), 10)
	if err != nil {
		return fmt.Errorf("发送消息失败: %v", err)
	}

	// 解析响应
	var response WeChatSendMessageResponse
	if err := tools.ParseReaderBody(res.Body, &response); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	if response.ErrCode != 0 {
		errMsg := fmt.Sprintf("企业微信API错误: %s (errcode: %d)", response.ErrMsg, response.ErrCode)
		if response.InvalidUser != "" {
			errMsg += fmt.Sprintf(", 无效用户: %s", response.InvalidUser)
		}
		if response.InvalidParty != "" {
			errMsg += fmt.Sprintf(", 无效部门: %s", response.InvalidParty)
		}
		if response.InvalidTag != "" {
			errMsg += fmt.Sprintf(", 无效标签: %s", response.InvalidTag)
		}
		return fmt.Errorf(errMsg)
	}

	logc.Info(nil, fmt.Sprintf("企业微信小程序消息发送成功: %s", messageJSON))
	return nil
}

// getAccessToken 获取access_token
func (w *WeChatMiniProgramSender) getAccessToken(corpId, secret string) (string, error) {
	w.tokenCache.mu.RLock()
	// 检查缓存是否有效（提前5分钟过期）
	if w.tokenCache.token != "" && w.tokenCache.expiresAt > time.Now().Unix()+300 {
		token := w.tokenCache.token
		w.tokenCache.mu.RUnlock()
		return token, nil
	}
	w.tokenCache.mu.RUnlock()

	// 获取新的access_token
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpId, secret)
	res, err := tools.Get(nil, url, 10)
	if err != nil {
		return "", fmt.Errorf("请求access_token失败: %v", err)
	}

	var response WeChatAccessTokenResponse
	if err := tools.ParseReaderBody(res.Body, &response); err != nil {
		return "", fmt.Errorf("解析access_token响应失败: %v", err)
	}

	if response.ErrCode != 0 {
		return "", fmt.Errorf("获取access_token失败: %s (errcode: %d)", response.ErrMsg, response.ErrCode)
	}

	// 缓存token
	w.tokenCache.mu.Lock()
	w.tokenCache.token = response.AccessToken
	w.tokenCache.expiresAt = time.Now().Unix() + int64(response.ExpiresIn)
	w.tokenCache.mu.Unlock()

	return response.AccessToken, nil
}
