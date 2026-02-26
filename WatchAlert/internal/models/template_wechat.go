package models

type WeChatMsgTemplate struct {
	MsgType  string         `json:"msgtype"`
	MarkDown WeChatMarkDown `json:"markdown"`
}

type WeChatMarkDown struct {
	Content string `json:"content"`
}

// WeChatMiniProgramMsgTemplate 企业微信小程序消息模板
type WeChatMiniProgramMsgTemplate struct {
	MsgType  string                     `json:"msgtype"`
	Text     *WeChatMiniProgramText     `json:"text,omitempty"`
	Markdown *WeChatMiniProgramMarkdown `json:"markdown,omitempty"`
	TextCard *WeChatMiniProgramTextCard `json:"textcard,omitempty"`
}

// WeChatMiniProgramText 文本消息
type WeChatMiniProgramText struct {
	Content string `json:"content"`
}

// WeChatMiniProgramMarkdown Markdown消息
type WeChatMiniProgramMarkdown struct {
	Content string `json:"content"`
}

// WeChatMiniProgramTextCard 文本卡片消息
type WeChatMiniProgramTextCard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	BtnTxt      string `json:"btntxt,omitempty"`
}

// WeChatWorkMsgTemplate 企业微信应用消息模板
type WeChatWorkMsgTemplate struct {
	ToUser   string              `json:"touser"`
	ToParty  string              `json:"toparty"`
	ToTag    string              `json:"totag"`
	MsgType  string              `json:"msgtype"`
	AgentID  int                 `json:"agentid"`
	Text     *WeChatWorkText     `json:"text,omitempty"`
	MarkDown *WeChatWorkMarkdown `json:"markdown,omitempty"`
}

// WeChatWorkText 文本消息
type WeChatWorkText struct {
	Content string `json:"content"`
}

// WeChatWorkMarkdown Markdown消息
type WeChatWorkMarkdown struct {
	Content string `json:"content"`
}
