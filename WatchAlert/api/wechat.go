package api

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"watchAlert/internal/types"
	"watchAlert/pkg/tools"
)

// WechatConfig 微信JS-SDK配置
type WechatConfig struct {
	AppId     string `json:"appId"`
	Timestamp int64  `json:"timestamp"`
	NonceStr  string `json:"nonceStr"`
	Signature string `json:"signature"`
}

// WechatTicketStatusUpdate 微信工单状态更新通知
type WechatTicketStatusUpdate struct {
	OpenId   string `json:"openId"`
	TicketId string `json:"ticketId"`
	TicketNo string `json:"ticketNo"`
	Status   string `json:"status"`
	Title    string `json:"title"`
	Url      string `json:"url"`
}

var WechatController = new(wechatControllerStruct)

type wechatControllerStruct struct{}

/*
微信相关 API
/api/w8t/wechat
*/
func (wc wechatControllerStruct) API(gin *gin.RouterGroup) {
	wechat := gin.Group("wechat")
	{
		wechat.GET("config", WechatController.GetWechatConfig)
		wechat.POST("template-message", WechatController.SendTemplateMessage)
		wechat.GET("user", WechatController.GetWechatUserInfo)
	}
}

// GetWechatConfig 获取微信JS-SDK配置
func (wc wechatControllerStruct) GetWechatConfig(ctx *gin.Context) {
	url := ctx.Query("url")
	if url == "" {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("url参数不能为空")
		})
		return
	}

	// 这里需要配置实际的微信公众号信息
	// 建议从环境变量或配置文件中读取
	appId := "your_wechat_app_id"         // 微信公众号AppId
	appSecret := "your_wechat_app_secret" // 微信公众号AppSecret

	// 获取access_token
	accessToken, err := wc.getAccessToken(appId, appSecret)
	if err != nil {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("获取access_token失败: %v", err)
		})
		return
	}

	// 获取jsapi_ticket
	jsapiTicket, err := wc.getJsapiTicket(accessToken)
	if err != nil {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("获取jsapi_ticket失败: %v", err)
		})
		return
	}

	// 生成签名
	timestamp := time.Now().Unix()
	nonceStr := tools.RandId()[:16]
	signature := wc.calculateSignature(jsapiTicket, nonceStr, timestamp, url)

	config := WechatConfig{
		AppId:     appId,
		Timestamp: timestamp,
		NonceStr:  nonceStr,
		Signature: signature,
	}

	Service(ctx, func() (interface{}, interface{}) {
		return config, nil
	})
}

// SendTemplateMessage 发送微信模板消息
func (wc wechatControllerStruct) SendTemplateMessage(ctx *gin.Context) {
	r := new(types.WechatTemplateMessage)
	BindJson(ctx, r)

	err := wc.sendTemplateMessage(r)
	if err != nil {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("发送模板消息失败: %v", err)
		})
		return
	}

	Service(ctx, func() (interface{}, interface{}) {
		return map[string]string{"message": "模板消息发送成功"}, nil
	})
}

// GetWechatUserInfo 获取微信用户信息
func (wc wechatControllerStruct) GetWechatUserInfo(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("code参数不能为空")
		})
		return
	}

	// 这里需要配置实际的微信公众号信息
	appId := "your_wechat_app_id"
	appSecret := "your_wechat_app_secret"

	// 通过code获取access_token
	tokenResp, err := wc.getWebAccessToken(appId, appSecret, code)
	if err != nil {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("获取access_token失败: %v", err)
		})
		return
	}

	// 获取用户信息
	userInfo, err := wc.getUserInfo(tokenResp.AccessToken, tokenResp.OpenId)
	if err != nil {
		Service(ctx, func() (interface{}, interface{}) {
			return nil, fmt.Errorf("获取用户信息失败: %v", err)
		})
		return
	}

	Service(ctx, func() (interface{}, interface{}) {
		return userInfo, nil
	})
}

// 内部方法

// WebAccessTokenResponse 网页授权access_token响应
type WebAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenId       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionId      string `json:"unionid"`
}

// WechatUserInfo 微信用户信息
type WechatUserInfo struct {
	OpenId     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        int    `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadImgUrl string `json:"headimgurl"`
	UnionId    string `json:"unionid"`
}

// getWebAccessToken 获取网页授权access_token
func (wc wechatControllerStruct) getWebAccessToken(appId, appSecret, code string) (*WebAccessTokenResponse, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		appId, appSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    int64  `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
		OpenId       string `json:"openid"`
		Scope        string `json:"scope"`
		UnionId      string `json:"unionid"`
		ErrCode      int    `json:"errcode"`
		ErrMsg       string `json:"errmsg"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信API错误: %s", result.ErrMsg)
	}

	return &WebAccessTokenResponse{
		AccessToken:  result.AccessToken,
		ExpiresIn:    result.ExpiresIn,
		RefreshToken: result.RefreshToken,
		OpenId:       result.OpenId,
		Scope:        result.Scope,
		UnionId:      result.UnionId,
	}, nil
}

// getUserInfo 获取微信用户信息
func (wc wechatControllerStruct) getUserInfo(accessToken, openId string) (*WechatUserInfo, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN",
		accessToken, openId)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		OpenId     string `json:"openid"`
		Nickname   string `json:"nickname"`
		Sex        int    `json:"sex"`
		Province   string `json:"province"`
		City       string `json:"city"`
		Country    string `json:"country"`
		HeadImgUrl string `json:"headimgurl"`
		UnionId    string `json:"unionid"`
		ErrCode    int    `json:"errcode"`
		ErrMsg     string `json:"errmsg"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if result.ErrCode != 0 {
		return nil, fmt.Errorf("微信API错误: %s", result.ErrMsg)
	}

	return &WechatUserInfo{
		OpenId:     result.OpenId,
		Nickname:   result.Nickname,
		Sex:        result.Sex,
		Province:   result.Province,
		City:       result.City,
		Country:    result.Country,
		HeadImgUrl: result.HeadImgUrl,
		UnionId:    result.UnionId,
	}, nil
}

// getAccessToken 获取微信access_token
func (wc wechatControllerStruct) getAccessToken(appId, appSecret string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		appId, appSecret)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int64  `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信API错误: %s", result.ErrMsg)
	}

	return result.AccessToken, nil
}

// getJsapiTicket 获取微信jsapi_ticket
func (wc wechatControllerStruct) getJsapiTicket(accessToken string) (string, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi",
		accessToken)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Ticket  string `json:"ticket"`
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if result.ErrCode != 0 {
		return "", fmt.Errorf("微信API错误: %s", result.ErrMsg)
	}

	return result.Ticket, nil
}

// calculateSignature 计算微信JS-SDK签名
func (wc wechatControllerStruct) calculateSignature(jsapiTicket, nonceStr string, timestamp int64, url string) string {
	// 注意：URL需要去除#及其后面的部分
	if idx := strings.Index(url, "#"); idx != -1 {
		url = url[:idx]
	}

	// 排序参数
	params := []string{
		"jsapi_ticket=" + jsapiTicket,
		"noncestr=" + nonceStr,
		"timestamp=" + fmt.Sprintf("%d", timestamp),
		"url=" + url,
	}
	sort.Strings(params)

	// 拼接字符串
	stringToSign := strings.Join(params, "&")

	// SHA1加密
	hash := sha1.Sum([]byte(stringToSign))
	return fmt.Sprintf("%x", hash)
}

// sendTemplateMessage 发送微信模板消息（示例）
func (wc wechatControllerStruct) sendTemplateMessage(msg *types.WechatTemplateMessage) error {
	// 这里需要实现实际的模板消息发送逻辑
	// 由于需要access_token，这里只是示例结构
	return nil
}

// jsonToString 将io.Reader转换为字符串
func jsonToString(body []byte) (string, error) {
	return string(body), nil
}
