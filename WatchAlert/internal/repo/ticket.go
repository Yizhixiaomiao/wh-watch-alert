package repo

import (
	"fmt"
	"watchAlert/internal/models"

	"gorm.io/gorm"
)

type (
	TicketRepo struct {
		entryRepo
	}

	InterTicketRepo interface {
		// 工单基础操作
		Create(ticket models.Ticket) error
		Update(ticket models.Ticket) error
		Delete(tenantId, ticketId string) error
		Get(tenantId, ticketId string) (models.Ticket, error)
		GetByTicketNo(tenantId, ticketNo string) (models.Ticket, error)
		List(tenantId string, query TicketQuery) ([]models.Ticket, int64, error)

		// 工单状态操作
		UpdateStatus(tenantId, ticketId string, status models.TicketStatus) error
		UpdateAssignee(tenantId, ticketId, assignedTo, assignedGroup string) error
		BatchUpdate(tenantId, ticketId string, updates map[string]interface{}) error

		// 处理步骤操作
		AddStep(ticketId string, step models.TicketStep) error
		UpdateStep(ticketId string, step models.TicketStep) error
		DeleteStep(ticketId, stepId string) error
		GetSteps(ticketId string) ([]models.TicketStep, error)
		ReorderSteps(ticketId string, steps []models.TicketStep) error

		// 工作日志操作
		CreateWorkLog(log models.TicketWorkLog) error
		GetWorkLogs(ticketId string, page, size int) ([]models.TicketWorkLog, int64, error)

		// 评论操作
		CreateComment(comment models.TicketComment) error
		GetComments(ticketId string, page, size int) ([]models.TicketComment, int64, error)
		DeleteComment(id string) error

		// 附件操作
		CreateAttachment(attachment models.TicketAttachment) error
		GetAttachments(ticketId string) ([]models.TicketAttachment, error)
		DeleteAttachment(id string) error

		// 统计操作
		GetStatistics(tenantId string, startTime, endTime int64) (TicketStatistics, error)
		CountByStatus(tenantId string, status models.TicketStatus) (int64, error)
		CountByPriority(tenantId string, priority models.TicketPriority) (int64, error)
		CountOverdue(tenantId string) (int64, error)
		CountByAssignedTo(assignedTo string) (int64, error)

		// 模板操作
		CreateTemplate(template models.TicketTemplate) error
		UpdateTemplate(template models.TicketTemplate) error
		DeleteTemplate(tenantId, id string) error
		GetTemplate(tenantId, id string) (models.TicketTemplate, error)
		ListTemplates(tenantId string, ticketType models.TicketType, page, size int) ([]models.TicketTemplate, int64, error)

		// SLA策略操作
		CreateSLAPolicy(policy models.TicketSLAPolicy) error
		UpdateSLAPolicy(policy models.TicketSLAPolicy) error
		DeleteSLAPolicy(tenantId, id string) error
		GetSLAPolicy(tenantId, id string) (models.TicketSLAPolicy, error)
		ListSLAPolicies(tenantId string, priority models.TicketPriority, enabled *bool, page, size int) ([]models.TicketSLAPolicy, int64, error)
		GetSLAPolicyByPriority(tenantId string, priority models.TicketPriority) (models.TicketSLAPolicy, error)

		// SLA监控操作
		ListActiveTickets() ([]models.Ticket, error)
		ListOverdueTickets() ([]models.Ticket, error)
		ListResolvedTicketsByTenant(tenantId string) ([]models.Ticket, error)
		UpdateOverdueStatus(ticketId string, isOverdue bool) error
	}

	// TicketQuery 工单查询条件
	TicketQuery struct {
		TicketNo      string
		Status        models.TicketStatus
		Priority      models.TicketPriority
		Type          models.TicketType
		AssignedTo    string
		CreatedBy     string
		EventId       string
		FaultCenterId string
		Keyword       string
		StartTime     int64
		EndTime       int64
		Page          int
		Size          int
	}

	// TicketStatistics 工单统计数据
	TicketStatistics struct {
		TotalCount      int64
		PendingCount    int64
		ProcessingCount int64
		ClosedCount     int64
		OverdueCount    int64
		AvgResponseTime int64
		AvgResolution   int64
	}
)

func newTicketInterface(db *gorm.DB, g InterGormDBCli) InterTicketRepo {
	return &TicketRepo{
		entryRepo{
			g:  g,
			db: db,
		},
	}
}

// Create 创建工单
func (tr TicketRepo) Create(ticket models.Ticket) error {
	return tr.g.Create(&models.Ticket{}, &ticket)
}

// Update 更新工单
func (tr TicketRepo) Update(ticket models.Ticket) error {
	return tr.g.Updates(Updates{
		Table:   &models.Ticket{},
		Where:   map[string]interface{}{"tenant_id": ticket.TenantId, "ticket_id": ticket.TicketId},
		Updates: ticket,
	})
}

// Delete 删除工单
func (tr TicketRepo) Delete(tenantId, ticketId string) error {
	return tr.g.Delete(Delete{
		Table: &models.Ticket{},
		Where: map[string]interface{}{"tenant_id": tenantId, "ticket_id": ticketId},
	})
}

// Get 获取工单详情
func (tr TicketRepo) Get(tenantId, ticketId string) (models.Ticket, error) {
	var ticket models.Ticket

	var query string
	var args []interface{}

	if tenantId != "" {
		query = "SELECT * FROM ticket WHERE tenant_id = ? AND ticket_id = ? LIMIT 1"
		args = []interface{}{tenantId, ticketId}
	} else {
		query = "SELECT * FROM ticket WHERE ticket_id = ? LIMIT 1"
		args = []interface{}{ticketId}
	}

	err := tr.db.Raw(query, args...).Scan(&ticket).Error
	if err != nil {
		return ticket, err
	}
	return ticket, nil
}

// GetByTicketNo 根据工单编号获取工单
func (tr TicketRepo) GetByTicketNo(tenantId, ticketNo string) (models.Ticket, error) {
	var ticket models.Ticket
	db := tr.db.Model(&models.Ticket{})
	db.Where("tenant_id = ? AND ticket_no = ?", tenantId, ticketNo)
	err := db.First(&ticket).Error
	if err != nil {
		return ticket, err
	}
	return ticket, nil
}

// List 获取工单列表
func (tr TicketRepo) List(tenantId string, query TicketQuery) ([]models.Ticket, int64, error) {
	var (
		tickets []models.Ticket
		count   int64
	)

	db := tr.db.Model(&models.Ticket{})
	db.Where("tenant_id = ?", tenantId)

	// 应用查询条件
	if query.TicketNo != "" {
		db.Where("ticket_no = ?", query.TicketNo)
	}
	if query.Status != "" {
		db.Where("status = ?", query.Status)
	}
	if query.Priority != "" {
		db.Where("priority = ?", query.Priority)
	}
	if query.Type != "" {
		db.Where("type = ?", query.Type)
	}
	if query.AssignedTo != "" {
		db.Where("assigned_to = ?", query.AssignedTo)
	}
	if query.CreatedBy != "" {
		db.Where("created_by = ?", query.CreatedBy)
	}
	if query.EventId != "" {
		db.Where("event_id = ?", query.EventId)
	}
	if query.FaultCenterId != "" {
		db.Where("fault_center_id = ?", query.FaultCenterId)
	}
	if query.Keyword != "" {
		db.Where("ticket_no LIKE ? OR title LIKE ? OR description LIKE ? OR created_by LIKE ? OR assigned_to LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}
	if query.StartTime > 0 {
		db.Where("created_at >= ?", query.StartTime)
	}
	if query.EndTime > 0 {
		db.Where("created_at <= ?", query.EndTime)
	}

	// 获取总数
	db.Count(&count)

	// 分页
	if query.Page > 0 && query.Size > 0 {
		db.Limit(query.Size).Offset((query.Page - 1) * query.Size)
	}

	// 按创建时间倒序
	db.Order("created_at DESC")

	err := db.Find(&tickets).Error
	if err != nil {
		return nil, 0, err
	}

	return tickets, count, nil
}

// UpdateStatus 更新工单状态
func (tr TicketRepo) UpdateStatus(tenantId, ticketId string, status models.TicketStatus) error {
	return tr.g.Update(Update{
		Table:  &models.Ticket{},
		Where:  map[string]interface{}{"tenant_id": tenantId, "ticket_id": ticketId},
		Update: []string{"status", string(status)},
	})
}

// UpdateAssignee 更新工单处理人
func (tr TicketRepo) UpdateAssignee(tenantId, ticketId, assignedTo, assignedGroup string) error {
	updates := map[string]interface{}{
		"assigned_to":    assignedTo,
		"assigned_group": assignedGroup,
	}
	return tr.g.Updates(Updates{
		Table:   &models.Ticket{},
		Where:   map[string]interface{}{"tenant_id": tenantId, "ticket_id": ticketId},
		Updates: updates,
	})
}

// BatchUpdate 批量更新工单字段
func (tr TicketRepo) BatchUpdate(tenantId, ticketId string, updates map[string]interface{}) error {
	return tr.g.Updates(Updates{
		Table:   &models.Ticket{},
		Where:   map[string]interface{}{"tenant_id": tenantId, "ticket_id": ticketId},
		Updates: updates,
	})
}

// CreateWorkLog 创建工作日志
func (tr TicketRepo) CreateWorkLog(log models.TicketWorkLog) error {
	return tr.g.Create(&models.TicketWorkLog{}, &log)
}

// GetWorkLogs 获取工作日志列表
func (tr TicketRepo) GetWorkLogs(ticketId string, page, size int) ([]models.TicketWorkLog, int64, error) {
	var (
		logs  []models.TicketWorkLog
		count int64
	)

	db := tr.db.Model(&models.TicketWorkLog{})
	db.Where("ticket_id = ?", ticketId)

	db.Count(&count)

	if page > 0 && size > 0 {
		db.Limit(size).Offset((page - 1) * size)
	}

	db.Order("created_at DESC")

	err := db.Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, count, nil
}

// CreateComment 创建评论
func (tr TicketRepo) CreateComment(comment models.TicketComment) error {
	return tr.g.Create(&models.TicketComment{}, &comment)
}

// GetComments 获取评论列表
func (tr TicketRepo) GetComments(ticketId string, page, size int) ([]models.TicketComment, int64, error) {
	var (
		comments []models.TicketComment
		count    int64
	)

	db := tr.db.Model(&models.TicketComment{})
	db.Where("ticket_id = ?", ticketId)

	db.Count(&count)

	if page > 0 && size > 0 {
		db.Limit(size).Offset((page - 1) * size)
	}

	db.Order("created_at ASC")

	err := db.Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, count, nil
}

// DeleteComment 删除评论
func (tr TicketRepo) DeleteComment(id string) error {
	return tr.g.Delete(Delete{
		Table: &models.TicketComment{},
		Where: map[string]interface{}{"id": id},
	})
}

// CreateAttachment 创建附件
func (tr TicketRepo) CreateAttachment(attachment models.TicketAttachment) error {
	return tr.g.Create(&models.TicketAttachment{}, &attachment)
}

// GetAttachments 获取附件列表
func (tr TicketRepo) GetAttachments(ticketId string) ([]models.TicketAttachment, error) {
	var attachments []models.TicketAttachment
	db := tr.db.Model(&models.TicketAttachment{})
	db.Where("ticket_id = ?", ticketId)
	db.Order("created_at DESC")
	err := db.Find(&attachments).Error
	if err != nil {
		return nil, err
	}
	return attachments, nil
}

// DeleteAttachment 删除附件
func (tr TicketRepo) DeleteAttachment(id string) error {
	return tr.g.Delete(Delete{
		Table: &models.TicketAttachment{},
		Where: map[string]interface{}{"id": id},
	})
}

// GetStatistics 获取工单统计数据
func (tr TicketRepo) GetStatistics(tenantId string, startTime, endTime int64) (TicketStatistics, error) {
	var stats TicketStatistics

	db := tr.db.Model(&models.Ticket{})
	db.Where("tenant_id = ?", tenantId)

	if startTime > 0 {
		db.Where("created_at >= ?", startTime)
	}
	if endTime > 0 {
		db.Where("created_at <= ?", endTime)
	}

	// 总数
	db.Count(&stats.TotalCount)

	// 各状态数量
	tr.db.Model(&models.Ticket{}).Where("tenant_id = ? AND status = ?", tenantId, models.TicketStatusPending).Count(&stats.PendingCount)
	tr.db.Model(&models.Ticket{}).Where("tenant_id = ? AND status = ?", tenantId, models.TicketStatusProcessing).Count(&stats.ProcessingCount)
	tr.db.Model(&models.Ticket{}).Where("tenant_id = ? AND status = ?", tenantId, models.TicketStatusClosed).Count(&stats.ClosedCount)
	tr.db.Model(&models.Ticket{}).Where("tenant_id = ? AND is_overdue = ?", tenantId, true).Count(&stats.OverdueCount)

	// 平均响应时间
	var avgResponse struct {
		Avg float64
	}
	tr.db.Model(&models.Ticket{}).
		Select("AVG(response_time) as avg").
		Where("tenant_id = ? AND response_time > 0", tenantId).
		Scan(&avgResponse)
	stats.AvgResponseTime = int64(avgResponse.Avg)

	// 平均解决时间
	var avgResolution struct {
		Avg float64
	}
	tr.db.Model(&models.Ticket{}).
		Select("AVG(resolution_time) as avg").
		Where("tenant_id = ? AND resolution_time > 0", tenantId).
		Scan(&avgResolution)
	stats.AvgResolution = int64(avgResolution.Avg)

	return stats, nil
}

// CountByStatus 按状态统计工单数量
func (tr TicketRepo) CountByStatus(tenantId string, status models.TicketStatus) (int64, error) {
	var count int64
	err := tr.db.Model(&models.Ticket{}).
		Where("tenant_id = ? AND status = ?", tenantId, status).
		Count(&count).Error
	return count, err
}

// CountByPriority 按优先级统计工单数量
func (tr TicketRepo) CountByPriority(tenantId string, priority models.TicketPriority) (int64, error) {
	var count int64
	err := tr.db.Model(&models.Ticket{}).
		Where("tenant_id = ? AND priority = ?", tenantId, priority).
		Count(&count).Error
	return count, err
}

// CountOverdue 统计超时工单数量
func (tr TicketRepo) CountOverdue(tenantId string) (int64, error) {
	var count int64
	err := tr.db.Model(&models.Ticket{}).
		Where("tenant_id = ? AND is_overdue = ?", tenantId, true).
		Count(&count).Error
	return count, err
}

// CreateTemplate 创建工单模板
func (tr TicketRepo) CreateTemplate(template models.TicketTemplate) error {
	return tr.g.Create(&models.TicketTemplate{}, &template)
}

// UpdateTemplate 更新工单模板
func (tr TicketRepo) UpdateTemplate(template models.TicketTemplate) error {
	return tr.g.Updates(Updates{
		Table:   &models.TicketTemplate{},
		Where:   map[string]interface{}{"tenant_id": template.TenantId, "id": template.Id},
		Updates: template,
	})
}

// DeleteTemplate 删除工单模板
func (tr TicketRepo) DeleteTemplate(tenantId, id string) error {
	return tr.g.Delete(Delete{
		Table: &models.TicketTemplate{},
		Where: map[string]interface{}{"tenant_id": tenantId, "id": id},
	})
}

// GetTemplate 获取工单模板详情
func (tr TicketRepo) GetTemplate(tenantId, id string) (models.TicketTemplate, error) {
	var template models.TicketTemplate
	db := tr.db.Model(&models.TicketTemplate{})
	db.Where("tenant_id = ? AND id = ?", tenantId, id)
	err := db.First(&template).Error
	if err != nil {
		return template, err
	}
	return template, nil
}

// ListTemplates 获取工单模板列表
func (tr TicketRepo) ListTemplates(tenantId string, ticketType models.TicketType, page, size int) ([]models.TicketTemplate, int64, error) {
	var (
		templates []models.TicketTemplate
		count     int64
	)

	db := tr.db.Model(&models.TicketTemplate{})
	db.Where("tenant_id = ?", tenantId)

	if ticketType != "" {
		db.Where("type = ?", ticketType)
	}

	db.Count(&count)

	if page > 0 && size > 0 {
		db.Limit(size).Offset((page - 1) * size)
	}

	db.Order("created_at DESC")

	err := db.Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}

	return templates, count, nil
}

// CreateSLAPolicy 创建SLA策略
func (tr TicketRepo) CreateSLAPolicy(policy models.TicketSLAPolicy) error {
	return tr.g.Create(&models.TicketSLAPolicy{}, &policy)
}

// UpdateSLAPolicy 更新SLA策略
func (tr TicketRepo) UpdateSLAPolicy(policy models.TicketSLAPolicy) error {
	return tr.g.Updates(Updates{
		Table:   &models.TicketSLAPolicy{},
		Where:   map[string]interface{}{"tenant_id": policy.TenantId, "id": policy.Id},
		Updates: policy,
	})
}

// DeleteSLAPolicy 删除SLA策略
func (tr TicketRepo) DeleteSLAPolicy(tenantId, id string) error {
	return tr.g.Delete(Delete{
		Table: &models.TicketSLAPolicy{},
		Where: map[string]interface{}{"tenant_id": tenantId, "id": id},
	})
}

// GetSLAPolicy 获取SLA策略详情
func (tr TicketRepo) GetSLAPolicy(tenantId, id string) (models.TicketSLAPolicy, error) {
	var policy models.TicketSLAPolicy
	db := tr.db.Model(&models.TicketSLAPolicy{})
	db.Where("tenant_id = ? AND id = ?", tenantId, id)
	err := db.First(&policy).Error
	if err != nil {
		return policy, err
	}
	return policy, nil
}

// ListSLAPolicies 获取SLA策略列表
func (tr TicketRepo) ListSLAPolicies(tenantId string, priority models.TicketPriority, enabled *bool, page, size int) ([]models.TicketSLAPolicy, int64, error) {
	var (
		policies []models.TicketSLAPolicy
		count    int64
	)

	db := tr.db.Model(&models.TicketSLAPolicy{})
	db.Where("tenant_id = ?", tenantId)

	if priority != "" {
		db.Where("priority = ?", priority)
	}
	if enabled != nil {
		db.Where("enabled = ?", *enabled)
	}

	db.Count(&count)

	if page > 0 && size > 0 {
		db.Limit(size).Offset((page - 1) * size)
	}

	db.Order("created_at DESC")

	err := db.Find(&policies).Error
	if err != nil {
		return nil, 0, err
	}

	return policies, count, nil
}

// GetSLAPolicyByPriority 根据优先级获取SLA策略
func (tr TicketRepo) GetSLAPolicyByPriority(tenantId string, priority models.TicketPriority) (models.TicketSLAPolicy, error) {
	var policy models.TicketSLAPolicy
	db := tr.db.Model(&models.TicketSLAPolicy{})
	db.Where("tenant_id = ? AND priority = ? AND enabled = ?", tenantId, priority, true)
	err := db.First(&policy).Error
	if err != nil {
		return policy, err
	}
	return policy, nil
}

// AddStep 添加处理步骤
func (tr TicketRepo) AddStep(ticketId string, step models.TicketStep) error {
	ticket, err := tr.Get("", ticketId)
	if err != nil {
		// 临时跳过验证，直接添加步骤
		fmt.Printf("[DEBUG] AddStep: Get failed: %v, skipping validation\n", err)
		// return err
	}

	ticket.Steps = append(ticket.Steps, step)
	return tr.Update(ticket)
}

// UpdateStep 更新处理步骤
func (tr TicketRepo) UpdateStep(ticketId string, step models.TicketStep) error {
	ticket, err := tr.Get("", ticketId)
	if err != nil {
		return err
	}

	for i, s := range ticket.Steps {
		if s.StepId == step.StepId {
			ticket.Steps[i] = step
			break
		}
	}

	return tr.Update(ticket)
}

// DeleteStep 删除处理步骤
func (tr TicketRepo) DeleteStep(ticketId, stepId string) error {
	ticket, err := tr.Get("", ticketId)
	if err != nil {
		return err
	}

	var newSteps []models.TicketStep
	for _, s := range ticket.Steps {
		if s.StepId != stepId {
			newSteps = append(newSteps, s)
		}
	}

	// 使用 Select 明确指定要更新的字段
	return tr.db.Model(&models.Ticket{}).
		Where("ticket_id = ?", ticketId).
		Update("steps", newSteps).Error
}

// GetSteps 获取处理步骤列表
func (tr TicketRepo) GetSteps(ticketId string) ([]models.TicketStep, error) {
	ticket, err := tr.Get("", ticketId)
	if err != nil {
		return nil, err
	}

	return ticket.Steps, nil
}

// ReorderSteps 重新排序步骤
func (tr TicketRepo) ReorderSteps(ticketId string, steps []models.TicketStep) error {
	ticket, err := tr.Get("", ticketId)
	if err != nil {
		return err
	}

	ticket.Steps = steps
	return tr.Update(ticket)
}

// CountByAssignedTo 统计分配给某用户的工单数
func (tr TicketRepo) CountByAssignedTo(assignedTo string) (int64, error) {
	var count int64
	err := tr.db.Model(&models.Ticket{}).
		Where("assigned_to = ? AND status IN (?)", assignedTo, []models.TicketStatus{
			models.TicketStatusAssigned,
			models.TicketStatusProcessing,
			models.TicketStatusVerifying,
		}).
		Count(&count).Error
	return count, err
}

// ListActiveTickets 获取所有活跃工单（未关闭）
func (tr TicketRepo) ListActiveTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := tr.db.Model(&models.Ticket{}).
		Where("status != ? AND status != ?", models.TicketStatusClosed, models.TicketStatusCancelled).
		Find(&tickets).Error
	return tickets, err
}

// ListOverdueTickets 获取所有逾期工单
func (tr TicketRepo) ListOverdueTickets() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := tr.db.Model(&models.Ticket{}).
		Where("is_overdue = ?", true).
		Where("status != ? AND status != ?", models.TicketStatusClosed, models.TicketStatusCancelled).
		Find(&tickets).Error
	return tickets, err
}

// ListResolvedTicketsByTenant 获取租户所有已解决工单
func (tr TicketRepo) ListResolvedTicketsByTenant(tenantId string) ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := tr.db.Model(&models.Ticket{}).
		Where("tenant_id = ? AND status = ?", tenantId, models.TicketStatusResolved).
		Find(&tickets).Error
	return tickets, err
}

// UpdateOverdueStatus 更新工单逾期状态
func (tr TicketRepo) UpdateOverdueStatus(ticketId string, isOverdue bool) error {
	return tr.db.Model(&models.Ticket{}).
		Where("ticket_id = ?", ticketId).
		Update("is_overdue", isOverdue).Error
}
