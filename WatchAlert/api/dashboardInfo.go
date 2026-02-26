package api

import (
	"sort"
	"watchAlert/internal/ctx"
	"watchAlert/internal/middleware"
	"watchAlert/internal/models"
	"watchAlert/internal/types"
	"watchAlert/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logc"
)

type dashboardInfoController struct{}

var DashboardInfoController = new(dashboardInfoController)

func (dashboardInfoController dashboardInfoController) API(gin *gin.RouterGroup) {
	system := gin.Group("system")
	system.Use(
		middleware.Auth(),
		middleware.ParseTenant(),
	)
	{
		system.GET("getDashboardInfo", dashboardInfoController.GetDashboardInfo)
	}
}

func (dashboardInfoController dashboardInfoController) GetDashboardInfo(context *gin.Context) {
	var c = ctx.DO()

	tid, _ := context.Get("TenantID")
	tidString := tid.(string)

	faultCenter, err := c.DB.FaultCenter().Get(tidString, context.Query("faultCenterId"), "")
	if err != nil {
		logc.Error(c.Ctx, err.Error())
		return
	}

	response.Success(context, types.ResponseDashboardInfo{
		CountAlertRules:   getRuleNumber(c, tidString),
		FaultCenterNumber: getFaultCenterNumber(c, tidString),
		UserNumber:        getUserNumber(c),
		CurAlertList:      getAlertList(c, faultCenter),
		AlarmDistribution: types.AlarmDistribution{
			P0: getAlarmDistribution(c, faultCenter, "P0"),
			P1: getAlarmDistribution(c, faultCenter, "P1"),
			P2: getAlarmDistribution(c, faultCenter, "P2"),
		},
	}, "success")
}

func getRuleNumber(ctx *ctx.Context, tenantId string) int64 {
	list, _, err := ctx.DB.Rule().List(tenantId, "", "", "", "", models.Page{
		Index: 0,
		Size:  10000,
	})
	if err != nil {
		return 0
	}
	return int64(len(list))
}

// getFaultCenterNumber 获取故障中心总数
func getFaultCenterNumber(ctx *ctx.Context, tenantId string) int64 {
	list, err := ctx.DB.FaultCenter().List(tenantId, "")
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
		return 0
	}
	return int64(len(list))
}

func getUserNumber(ctx *ctx.Context) int64 {
	list, err := ctx.DB.User().List("", "", "")
	if err != nil {
		logc.Error(ctx.Ctx, err.Error())
		return 0
	}
	return int64(len(list))
}

// getAlertList 获取当前告警列表，返回最近10个规则的告警（按规则最近触发时间倒序，每个规则只显示一次）
func getAlertList(ctx *ctx.Context, faultCenter models.FaultCenter) []types.AlertList {
	eventsMap, err := ctx.Redis.Alert().GetAllEvents(models.BuildAlertEventCacheKey(faultCenter.TenantId, faultCenter.ID))
	if err != nil {
		logc.Error(ctx.Ctx, "Failed to get events from cache:", err.Error())
		return nil
	}

	// 如果没有事件，直接返回空列表
	if len(eventsMap) == 0 {
		return []types.AlertList{}
	}

	// 按规则名分组，保留每个规则最近的一条告警
	ruleLatestEvent := make(map[string]*models.AlertCurEvent)
	for _, event := range eventsMap {
		ruleName := event.RuleName
		if existing, exists := ruleLatestEvent[ruleName]; exists {
			// 如果已存在该规则，比较触发时间，保留最近的
			if event.FirstTriggerTime > existing.FirstTriggerTime {
				ruleLatestEvent[ruleName] = event
			}
		} else {
			ruleLatestEvent[ruleName] = event
		}
	}

	// 将 map 转换为数组以便排序
	var events []*models.AlertCurEvent
	for _, event := range ruleLatestEvent {
		events = append(events, event)
	}

	// 按触发时间倒序排序
	sort.Slice(events, func(i, j int) bool {
		return events[i].FirstTriggerTime > events[j].FirstTriggerTime
	})

	// 最多返回10条
	limit := 10
	if len(events) < limit {
		limit = len(events)
	}

	var list []types.AlertList
	for i := 0; i < limit; i++ {
		list = append(list, types.AlertList{
			Severity:      events[i].Severity,
			RuleName:      events[i].RuleName,
			FaultCenterId: events[i].FaultCenterId,
		})
	}

	return list
}

// getAlarmDistribution 获取告警分布
func getAlarmDistribution(ctx *ctx.Context, faultCenter models.FaultCenter, severity string) int64 {
	events, err := ctx.Redis.Alert().GetAllEvents(models.BuildAlertEventCacheKey(faultCenter.TenantId, faultCenter.ID))
	if err != nil {
		return 0
	}

	var number int64
	for _, event := range events {
		if event.Severity == severity {
			number++
		}
	}
	return number
}
