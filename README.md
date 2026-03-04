# WatchAlert 项目开发记录

## 2026-03-04 告警工单AI处理建议和知识库同步功能

### 功能概述

本次更新实现了告警工单创建时自动获取AI处理建议，以及将处理建议和处理步骤同步到知识库的功能。

### 修改的文件

| 文件 | 修改内容 |
|------|----------|
| `WatchAlert/internal/services/alert_ticket.go` | 重构AI提示词构建逻辑，添加告警事件详情合并 |
| `WatchAlert/internal/services/ticket.go` | 添加 `SyncTreatmentSuggestionToKnowledge` 和 `SyncStepToKnowledge` 方法 |
| `WatchAlert/internal/models/ticket.go` | Ticket模型添加 `TreatmentSuggestion` 字段，TicketStep模型添加 `KnowledgeId` 字段 |
| `WatchAlert/internal/types/ticket.go` | 添加 `RequestTicketTreatmentSuggestionSync` 和 `RequestTicketStepSync` 请求类型 |
| `WatchAlert/api/ticket.go` | 添加API路由和控制器方法 |

### 新增API接口

#### 1. 同步处理建议到知识库

**接口:** `POST /api/w8t/ticket/treatment-suggestion/sync`

**请求参数:**
```json
{
  "ticketId": "tk-xxx",        // 必填，工单ID
  "category": "故障处理",       // 必填，知识分类
  "tags": ["告警", "数据库"],   // 可选，知识标签
  "publishNow": true           // 可选，是否立即发布，默认false
}
```

**响应:**
```json
{
  "knowledgeId": "kn-xxx"
}
```

#### 2. 同步处理步骤到知识库

**接口:** `POST /api/w8t/ticket/step/sync`

**请求参数:**
```json
{
  "ticketId": "tk-xxx",        // 必填，工单ID
  "stepId": "step-xxx",        // 必填，步骤ID
  "category": "故障处理",       // 必填，知识分类
  "tags": ["排查步骤", "数据库"], // 可选，知识标签
  "publishNow": true           // 可选，是否立即发布，默认false
}
```

**响应:**
```json
{
  "knowledgeId": "kn-xxx"
}
```

### 核心功能说明

#### 1. 告警触发工单创建时获取AI处理建议

当告警触发工单创建时，系统会自动：
1. 检查AI配置是否启用
2. 获取AI客户端
3. 构建AI提示词，包含：
   - 规则名称 (RuleName)
   - 触发条件 (SearchQL)
   - 标签信息
   - 告警注解
   - 事件详情（数据源类型、严重程度、触发时间等）
4. 调用AI获取处理建议
5. 将建议写入工单的 `TreatmentSuggestion` 字段

#### 2. 处理建议同步到知识库

- 支持将工单的处理建议同步到知识库
- 如果已同步，再次调用会更新已有知识内容
- 同步后会记录工作日志

#### 3. 处理步骤同步到知识库

- 支持将单个处理步骤同步到知识库
- 同步后会更新步骤的 `KnowledgeId` 字段
- 知识内容包含：关联工单信息、步骤标题、步骤描述、处理方法、处理结果

### 数据库变更

需要在 `ticket` 表添加字段：
```sql
ALTER TABLE ticket ADD COLUMN treatment_suggestion TEXT;
```

需要在 `ticket_step` 表添加字段：
```sql
ALTER TABLE ticket_step ADD COLUMN knowledge_id VARCHAR(255);
CREATE INDEX idx_knowledge_id ON ticket_step(knowledge_id);
```

### 提交记录

- 提交ID: `04131a0`
- 提交信息: `feat: 添加告警工单AI处理建议和知识库同步功能`