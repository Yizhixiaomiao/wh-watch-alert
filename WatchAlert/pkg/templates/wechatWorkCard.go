package templates

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
	models2 "watchAlert/internal/models"
	"watchAlert/pkg/tools"
)

func wechatWorkTemplate(alert models2.AlertCurEvent, noticeTmpl models2.NoticeTemplateExample, hook string) string {
	alert = prepareAlertData(alert)

	data := tools.ConvertStructToMap(alert)

	variableMap := map[string]string{
		"Title":          "rule_name",
		"Severity":       "severity",
		"Status":         "status",
		"TriggeredAt":    "first_trigger_time_format",
		"RuleName":       "rule_name",
		"DatasourceName": "datasource_id",
		"Content":        "annotations",
		"Labels":         "labels",
		"Fingerprint":    "fingerprint",
		"TenantID":       "tenantId",
		"DashboardURL":   "dashboard_url",
		"IsRecovered":    "is_recovered",
	}

	templateContent := noticeTmpl.Template
	if alert.IsRecovered && noticeTmpl.TemplateRecover != "" {
		templateContent = noticeTmpl.TemplateRecover
	} else if !alert.IsRecovered && noticeTmpl.TemplateFiring != "" {
		templateContent = noticeTmpl.TemplateFiring
	}

	re := regexp.MustCompile(`\{\{(.*?)\}\}`)
	templateStr := re.ReplaceAllStringFunc(templateContent, func(match string) string {
		variable := strings.TrimSpace(match[2 : len(match)-2])
		if fieldName, ok := variableMap[variable]; ok {
			return "${" + fieldName + "}"
		}
		return match
	})

	content := tools.ParserVariables(templateStr, data)

	toUser := "@all"
	agentID := 0

	if hook != "" {
		parsedURL, err := url.Parse(hook)
		if err == nil {
			queryParams := parsedURL.Query()
			if touser := queryParams.Get("touser"); touser != "" {
				toUser = touser
			}
			if agentid := queryParams.Get("agentid"); agentid != "" {
				if id, err := strconv.Atoi(agentid); err == nil {
					agentID = id
				}
			}
		}
	}

	t := models2.WeChatWorkMsgTemplate{
		ToUser:  toUser,
		MsgType: "text",
		AgentID: agentID,
		Text: &models2.WeChatWorkText{
			Content: content,
		},
	}

	return tools.JsonMarshalToString(t)
}
