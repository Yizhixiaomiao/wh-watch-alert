package sender

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"watchAlert/pkg/tools"

	"github.com/bytedance/sonic"
	"github.com/zeromicro/go-zero/core/logc"
)

type (
	// WeChatWorkSender 企业微信应用发送策略
	WeChatWorkSender struct{}

	WeChatWorkResponse struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
)

func NewWeChatWorkSender() SendInter {
	return &WeChatWorkSender{}
}

func (w *WeChatWorkSender) Send(params SendParams) error {
	return w.post(params.Hook, params.Content)
}

func (w *WeChatWorkSender) Test(params SendParams) error {
	parsedURL, err := url.Parse(params.Hook)
	if err != nil {
		return err
	}
	queryParams := parsedURL.Query()
	agentid := queryParams.Get("agentid")
	if agentid == "" {
		agentid = "0"
	}

	testContent := fmt.Sprintf(`{
		"touser": "@all",
		"msgtype": "text",
		"agentid": %s,
		"text": {
			"content": "%s"
		}
	}`, agentid, RobotTestContent)

	return w.post(params.Hook, testContent)
}

func (w *WeChatWorkSender) post(hook, content string) error {
	logc.Info(context.Background(), fmt.Sprintf("[WeChatWork] Original hook: %s", hook))

	logc.Info(context.Background(), fmt.Sprintf("[WeChatWork] Test content: %s", content))
	hookURL, err := w.parseHook(hook)
	if err != nil {
		return err
	}

	logc.Info(context.Background(), fmt.Sprintf("[WeChatWork] Parsed hook URL: %s", hookURL))

	res, err := tools.Post(nil, hookURL, bytes.NewReader([]byte(content)), 30)
	if err != nil {
		return err
	}

	var response WeChatWorkResponse
	if err := tools.ParseReaderBody(res.Body, &response); err != nil {
		return errors.New(fmt.Sprintf("Error unmarshalling WeChatWork response: %s", err.Error()))
	}
	if response.ErrCode != 0 {
		return errors.New(fmt.Sprintf("WeChatWork API error: %d - %s", response.ErrCode, response.ErrMsg))
	}

	return nil
}

// parseHook 解析 Hook URL
// 支持以下格式：
// 1. https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=ACCESS_TOKEN&agentid=AGENTID&touser=TOUSER
// 2. 包含 corp_id 和 corp_secret 的完整配置，需要先获取 access_token
func (w *WeChatWorkSender) parseHook(hook string) (string, error) {
	logc.Info(context.Background(), fmt.Sprintf("[WeChatWork parseHook] Input hook: %s", hook))

	if hook == "" {
		return "", errors.New("WeChatWork hook is empty")
	}

	parsedURL, err := url.Parse(hook)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Invalid WeChatWork hook URL: %s", err.Error()))
	}

	queryParams := parsedURL.Query()
	logc.Info(context.Background(), fmt.Sprintf("[WeChatWork parseHook] Query params: %v", queryParams))

	// 检查是否已经包含 access_token
	accessToken := queryParams.Get("access_token")
	if accessToken == "" {
		corpID := queryParams.Get("corp_id")
		corpSecret := queryParams.Get("corp_secret")

		logc.Info(context.Background(), fmt.Sprintf("[WeChatWork parseHook] No access_token found, corp_id=%s, corp_secret=%s", corpID, corpSecret))

		if corpID == "" || corpSecret == "" {
			return "", errors.New("WeChatWork hook must contain either access_token or corp_id and corp_secret")
		}

		// 获取 access_token
		accessToken, err = w.getAccessToken(corpID, corpSecret)
		if err != nil {
			logc.Error(context.Background(), fmt.Sprintf("[WeChatWork parseHook] Failed to get access_token: %v", err))
			return "", err
		}

		logc.Info(context.Background(), fmt.Sprintf("[WeChatWork parseHook] Got access_token: %s", accessToken))

		// 清空所有查询参数，只保留必要的参数
		newQuery := url.Values{}
		newQuery.Set("access_token", accessToken)

		// 保留 agentid 和 touser 参数
		agentid := queryParams.Get("agentid")
		touser := queryParams.Get("touser")

		if agentid != "" {
			newQuery.Set("agentid", agentid)
			logc.Info(context.Background(), fmt.Sprintf("[WeChatWork parseHook] Added agentid: %s", agentid))
		}
		if touser != "" {
			newQuery.Set("touser", touser)
			logc.Info(context.Background(), fmt.Sprintf("[WeChatWork parseHook] Added touser: %s", touser))
		}

		parsedURL.RawQuery = newQuery.Encode()
		logc.Info(context.Background(), fmt.Sprintf("[WeChatWork parseHook] Final query params: %v", newQuery))
	} else {
		logc.Info(context.Background(), fmt.Sprintf("[WeChatWork parseHook] access_token already present: %s", accessToken))
	}

	finalURL := parsedURL.String()
	logc.Info(context.Background(), fmt.Sprintf("[WeChatWork parseHook] Final URL: %s", finalURL))

	return finalURL, nil
}

// getAccessToken 获取企业微信应用的 access_token
func (w *WeChatWorkSender) getAccessToken(corpID, corpSecret string) (string, error) {
	getTokenURL := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", corpID, corpSecret)

	res, err := tools.Get(nil, getTokenURL, 30)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Failed to get access_token: %s", err.Error()))
	}

	var tokenResponse struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := tools.ParseReaderBody(res.Body, &tokenResponse); err != nil {
		return "", errors.New(fmt.Sprintf("Failed to parse access_token response: %s", err.Error()))
	}

	if tokenResponse.ErrCode != 0 {
		return "", errors.New(fmt.Sprintf("Get access_token failed: %d - %s", tokenResponse.ErrCode, tokenResponse.ErrMsg))
	}

	return tokenResponse.AccessToken, nil
}

// GetSendMsg 发送内容
func (w *WeChatWorkSender) GetSendMsg(content string) map[string]any {
	msg := make(map[string]any)
	if content == "" {
		return msg
	}
	err := sonic.Unmarshal([]byte(content), &msg)
	if err != nil {
		return msg
	}
	return msg
}
