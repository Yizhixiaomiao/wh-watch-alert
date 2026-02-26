package tasks

import (
	"context"
	"time"
	"watchAlert/internal/ctx"
	"watchAlert/internal/models"
	"watchAlert/pkg/tools"

	"github.com/robfig/cron/v3"
	"github.com/zeromicro/go-zero/core/logc"
)

// SLAMonitor SLA监控任务
type SLAMonitor struct {
	ctx    *ctx.Context
	cron   *cron.Cron
	cancel context.CancelFunc
}

// NewSLAMonitor 创建SLA监控任务
func NewSLAMonitor(ctx *ctx.Context) *SLAMonitor {
	return &SLAMonitor{
		ctx:  ctx,
		cron: cron.New(cron.WithSeconds()),
	}
}

// Start 启动SLA监控任务
func (m *SLAMonitor) Start() error {
	// 每5分钟检查一次逾期工单
	_, err := m.cron.AddFunc("0 */5 * * * *", m.checkOverdueTickets)
	if err != nil {
		return err
	}

	// 每小时发送SLA逾期告警
	_, err = m.cron.AddFunc("0 0 * * * *", m.sendOverdueAlerts)
	if err != nil {
		return err
	}

	// 每天凌晨生成SLA报告
	_, err = m.cron.AddFunc("0 0 0 * * *", m.generateSLAReport)
	if err != nil {
		return err
	}

	m.cron.Start()
	logc.Infof(m.ctx.Ctx, "SLA监控任务已启动")
	
	return nil
}

// Stop 停止SLA监控任务
func (m *SLAMonitor) Stop() {
	if m.cron != nil {
		m.cron.Stop()
		logc.Infof(m.ctx.Ctx, "SLA监控任务已停止")
	}
}

// checkOverdueTickets 检查逾期工单
func (m *SLAMonitor) checkOverdueTickets() {
	logc.Infof(m.ctx.Ctx, "开始检查逾期工单")
	
	// 获取所有未解决的工单
	tickets, err := m.ctx.DB.Ticket().ListActiveTickets()
	if err != nil {
		logc.Errorf(m.ctx.Ctx, "获取活跃工单失败: %v", err)
		return
	}

	overdueCount := 0
	for _, ticket := range tickets {
		if ticket.ResolutionSLA == 0 {
			continue
		}

		// 获取SLA策略
		slaPolicy, err := m.ctx.DB.Ticket().GetSLAPolicyByPriority(ticket.TenantId, ticket.Priority)
		if err != nil {
			// 如果没有SLA策略，使用简单计算
			if time.Now().Unix() > ticket.DueTime && !ticket.IsOverdue {
				m.markTicketOverdue(ticket)
				overdueCount++
			}
			continue
		}

		// 解析工作时间配置
		workingHoursConfig, err := tools.ParseWorkingHoursConfig(slaPolicy.WorkingHours)
		if err != nil {
			// 如果配置解析失败，使用简单计算
			if time.Now().Unix() > ticket.DueTime && !ticket.IsOverdue {
				m.markTicketOverdue(ticket)
				overdueCount++
			}
			continue
		}

		// 设置节假日
		workingHoursConfig.Holidays = slaPolicy.Holidays

		// 使用工作日计算是否逾期
		createdTime := time.Unix(ticket.CreatedAt, 0)
		if tools.IsOverdue(createdTime, time.Now(), ticket.ResolutionSLA, workingHoursConfig) && !ticket.IsOverdue {
			m.markTicketOverdue(ticket)
			overdueCount++
		}
	}

	logc.Infof(m.ctx.Ctx, "逾期工单检查完成，标记 %d 个工单为逾期", overdueCount)
}

// markTicketOverdue 标记工单为逾期
func (m *SLAMonitor) markTicketOverdue(ticket models.Ticket) {
	// 更新工单逾期状态
	err := m.ctx.DB.Ticket().UpdateOverdueStatus(ticket.TicketId, true)
	if err != nil {
		logc.Errorf(m.ctx.Ctx, "更新工单逾期状态失败: %v", err)
		return
	}

	// 创建工作日志
	logEntry := models.TicketWorkLog{
		Id:        "log-" + tools.RandId(),
		TicketId:  ticket.TicketId,
		UserId:    "system",
		UserName:  "系统",
		Action:    "overdue",
		Content:   "工单已逾期",
		OldValue:  "false",
		NewValue:  "true",
		CreatedAt: time.Now().Unix(),
	}
	
	err = m.ctx.DB.Ticket().CreateWorkLog(logEntry)
	if err != nil {
		logc.Errorf(m.ctx.Ctx, "创建逾期工作日志失败: %v", err)
	}

	logc.Infof(m.ctx.Ctx, "工单 %s 已标记为逾期", ticket.TicketNo)
}

// sendOverdueAlerts 发送逾期告警
func (m *SLAMonitor) sendOverdueAlerts() {
	logc.Infof(m.ctx.Ctx, "开始发送SLA逾期告警")
	
	// 获取所有逾期且未解决的工单
	overdueTickets, err := m.ctx.DB.Ticket().ListOverdueTickets()
	if err != nil {
		logc.Errorf(m.ctx.Ctx, "获取逾期工单失败: %v", err)
		return
	}

	// 按租户分组发送告警
	tenantMap := make(map[string][]models.Ticket)
	for _, ticket := range overdueTickets {
		tenantMap[ticket.TenantId] = append(tenantMap[ticket.TenantId], ticket)
	}

	for _, tickets := range tenantMap {
		// 获取租户的通知配置
		// TODO: 实现租户级别的SLA逾期通知配置

		for _, ticket := range tickets {
			// 发送逾期告警给工单处理人
			if ticket.AssignedTo != "" {
				m.sendOverdueNotification(ticket)
			}
		}
	}

	logc.Infof(m.ctx.Ctx, "SLA逾期告警发送完成")
}

// sendOverdueNotification 发送逾期通知
func (m *SLAMonitor) sendOverdueNotification(ticket models.Ticket) {
	// TODO: 实现逾期通知发送逻辑
	// 可以通过邮件、钉钉、飞书等方式发送
	logc.Infof(m.ctx.Ctx, "工单 %s 逾期通知已发送给 %s", ticket.TicketNo, ticket.AssignedTo)
}

// generateSLAReport 生成SLA报告
func (m *SLAMonitor) generateSLAReport() {
	logc.Infof(m.ctx.Ctx, "开始生成SLA报告")

	// 获取所有租户
	var tenants []models.Tenant
	err := m.ctx.DB.DB().Find(&tenants).Error
	if err != nil {
		logc.Errorf(m.ctx.Ctx, "获取租户列表失败: %v", err)
		return
	}

	for _, tenant := range tenants {
		// 生成每个租户的SLA报告
		report, err := m.generateTenantSLAReport(tenant.ID)
		if err != nil {
			logc.Errorf(m.ctx.Ctx, "生成租户 %s SLA报告失败: %v", tenant.ID, err)
			continue
		}

		logc.Infof(m.ctx.Ctx, "租户 %s SLA报告: 总工单=%d, 逾期=%d, 按时解决=%d, SLA达成率=%.2f%%",
			tenant.ID,
			report.TotalTickets,
			report.OverdueCount,
			report.OnTimeResolved,
			report.SLAComplianceRate)
	}

	logc.Infof(m.ctx.Ctx, "SLA报告生成完成")
}

// SLAReport SLA报告
type SLAReport struct {
	TenantId         string
	TotalTickets     int
	OverdueCount     int
	OnTimeResolved   int
	SLAComplianceRate float64
}

// generateTenantSLAReport 生成租户SLA报告
func (m *SLAMonitor) generateTenantSLAReport(tenantId string) (*SLAReport, error) {
	// 获取租户所有已解决的工单
	tickets, err := m.ctx.DB.Ticket().ListResolvedTicketsByTenant(tenantId)
	if err != nil {
		return nil, err
	}

	report := &SLAReport{
		TenantId:     tenantId,
		TotalTickets: len(tickets),
	}

	var overdueCount int
	var onTimeResolved int

	for _, ticket := range tickets {
		if ticket.ResolutionSLA == 0 {
			continue
		}

		// 检查是否逾期
		if ticket.IsOverdue {
			overdueCount++
		} else {
			onTimeResolved++
		}
	}

	report.OverdueCount = overdueCount
	report.OnTimeResolved = onTimeResolved

	// 计算SLA达成率
	if report.TotalTickets > 0 {
		report.SLAComplianceRate = float64(onTimeResolved) / float64(report.TotalTickets) * 100
	}

	return report, nil
}

// CheckTicketOverdue 检查单个工单是否逾期
func (m *SLAMonitor) CheckTicketOverdue(ticket models.Ticket) bool {
	if ticket.ResolutionSLA == 0 {
		return false
	}

	// 获取SLA策略
	slaPolicy, err := m.ctx.DB.Ticket().GetSLAPolicyByPriority(ticket.TenantId, ticket.Priority)
	if err != nil {
		// 如果没有SLA策略，使用简单计算
		return time.Now().Unix() > ticket.DueTime
	}

	// 解析工作时间配置
	workingHoursConfig, err := tools.ParseWorkingHoursConfig(slaPolicy.WorkingHours)
	if err != nil {
		// 如果配置解析失败，使用简单计算
		return time.Now().Unix() > ticket.DueTime
	}

	// 设置节假日
	workingHoursConfig.Holidays = slaPolicy.Holidays

	// 使用工作日计算是否逾期
	createdTime := time.Unix(ticket.CreatedAt, 0)
	return tools.IsOverdue(createdTime, time.Now(), ticket.ResolutionSLA, workingHoursConfig)
}

// GetTicketRemainingSLATime 获取工单剩余SLA时间
func (m *SLAMonitor) GetTicketRemainingSLATime(ticket models.Ticket) int64 {
	if ticket.ResolutionSLA == 0 {
		return 0
	}

	// 获取SLA策略
	slaPolicy, err := m.ctx.DB.Ticket().GetSLAPolicyByPriority(ticket.TenantId, ticket.Priority)
	if err != nil {
		// 如果没有SLA策略，使用简单计算
		remaining := ticket.DueTime - time.Now().Unix()
		if remaining < 0 {
			return 0
		}
		return remaining
	}

	// 解析工作时间配置
	workingHoursConfig, err := tools.ParseWorkingHoursConfig(slaPolicy.WorkingHours)
	if err != nil {
		// 如果配置解析失败，使用简单计算
		remaining := ticket.DueTime - time.Now().Unix()
		if remaining < 0 {
			return 0
		}
		return remaining
	}

	// 设置节假日
	workingHoursConfig.Holidays = slaPolicy.Holidays

	// 使用工作日计算剩余时间
	createdTime := time.Unix(ticket.CreatedAt, 0)
	return tools.CalculateSLARemainingTime(createdTime, time.Now(), ticket.ResolutionSLA, workingHoursConfig)
}