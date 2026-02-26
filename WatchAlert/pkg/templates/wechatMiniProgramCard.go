package templates

import (
	"fmt"
	models2 "watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

func wechatMiniProgramTemplate(alert models2.AlertCurEvent, noticeTmpl models2.NoticeTemplateExample) string {
	Title := ParserTemplate("Title", alert, noticeTmpl.Template)
	Event := ParserTemplate("Event", alert, noticeTmpl.Template)
	Footer := ParserTemplate("Footer", alert, noticeTmpl.Template)

	// 默认使用markdown格式
	t := models2.WeChatMiniProgramMsgTemplate{
		MsgType: "markdown",
		Markdown: &models2.WeChatMiniProgramMarkdown{
			Content: fmt.Sprintf("**%s**\n\n%s\n\n%s", Title, Event, Footer),
		},
	}

	return tools.JsonMarshalToString(t)
}
