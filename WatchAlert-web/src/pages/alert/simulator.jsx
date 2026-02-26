import { useState } from "react"
import { Button, Card, Form, Input, Select, Switch, InputNumber, message, Space, Descriptions, Tag, Alert as AntAlert } from "antd"
import { PlayCircleOutlined, ReloadOutlined, DeleteOutlined, ThunderboltOutlined } from "@ant-design/icons"
import axios from "axios"

const { Option } = Select

export const AlertSimulator = () => {
    const [form] = Form.useForm()
    const [loading, setLoading] = useState(false)
    const [lastAlert, setLastAlert] = useState(null)
    const [alertHistory, setAlertHistory] = useState([])

    const severityColors = {
        Critical: "red",
        Warning: "orange",
        Info: "blue",
    }

    const statusColors = {
        pre_alert: "orange",
        alerting: "red",
        recovered: "green",
    }

    const statusText = {
        pre_alert: "预告警",
        alerting: "告警中",
        recovered: "已恢复",
    }

    const createMockAlert = async (values) => {
        setLoading(true)
        try {
            const response = await axios.post("/api/w8t/debug/createMockAlert", {
                rule_name: values.ruleName,
                severity: values.severity,
                labels: {
                    instance: values.instance || "localhost:9090",
                    job: values.job || "node_exporter",
                    ...values.customLabels,
                },
                auto_create_ticket: values.autoCreateTicket,
                auto_recover: values.autoRecover,
                recover_after: values.recoverAfter,
                duration: values.duration,
                tenant_id: values.tenantId,
                fault_center_id: values.faultCenterId,
            })

            if (response.data.success) {
                message.success("模拟告警创建成功")
                const alertData = response.data.alert
                setLastAlert(alertData)
                setAlertHistory([alertData, ...alertHistory].slice(0, 10))

                if (values.autoCreateTicket) {
                    message.info("系统将根据规则自动创建工单")
                }
            }
        } catch (error) {
            message.error("创建失败: " + (error.response?.data?.message || error.message))
        } finally {
            setLoading(false)
        }
    }

    const recoverAlert = async (eventId) => {
        try {
            const response = await axios.post("/api/w8t/debug/recoverMockAlert", {
                event_id: eventId,
            })

            if (response.data.success) {
                message.success("告警已恢复")
                setAlertHistory(alertHistory.map(alert =>
                    alert.event_id === eventId ? { ...alert, status: "recovered" } : alert
                ))
                if (lastAlert && lastAlert.event_id === eventId) {
                    setLastAlert({ ...lastAlert, status: "recovered" })
                }
            }
        } catch (error) {
            message.error("恢复失败: " + (error.response?.data?.message || error.message))
        }
    }

    const cleanupAlerts = async () => {
        try {
            const response = await axios.post("/api/w8t/debug/cleanupMockAlerts", {
                tenant_id: form.getFieldValue("tenantId") || "demo-tenant-001",
            })

            if (response.data.success) {
                message.success("已清理所有模拟告警数据")
                setLastAlert(null)
                setAlertHistory([])
            }
        } catch (error) {
            message.error("清理失败: " + (error.response?.data?.message || error.message))
        }
    }

    return (
        <div style={{ padding: "24px" }}>
            <div style={{ marginBottom: "24px" }}>
                <h2 style={{ display: "flex", alignItems: "center", gap: "8px" }}>
                    <ThunderboltOutlined style={{ color: "#ff4d4f" }} />
                    告警模拟器
                </h2>
                <p style={{ color: "#666" }}>创建模拟告警，测试告警转工单和通知推送功能</p>
            </div>

            <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "24px" }}>
                <Card title="创建模拟告警" bordered={false}>
                    <Form
                        form={form}
                        layout="vertical"
                        onFinish={createMockAlert}
                        initialValues={{
                            ruleName: "CPU High Usage",
                            severity: "Critical",
                            instance: "localhost:9090",
                            job: "node_exporter",
                            autoCreateTicket: true,
                            autoRecover: true,
                            recoverAfter: 60,
                            duration: 30,
                            tenantId: "default",
                            faultCenterId: "default",
                        }}
                    >
                        <Form.Item
                            label="规则名称"
                            name="ruleName"
                            rules={[{ required: true, message: "请输入规则名称" }]}
                            tooltip="建议使用英文，避免中文编码问题"
                        >
                            <Input placeholder="例如: CPU High Usage" />
                        </Form.Item>

                        <Form.Item
                            label="严重程度"
                            name="severity"
                            rules={[{ required: true, message: "请选择严重程度" }]}
                        >
                            <Select>
                                <Option value="Critical">Critical - 严重</Option>
                                <Option value="Warning">Warning - 警告</Option>
                                <Option value="Info">Info - 信息</Option>
                            </Select>
                        </Form.Item>

                        <Form.Item label="标签信息">
                            <Space direction="vertical" style={{ width: "100%" }}>
                                <Form.Item name="instance" noStyle>
                                    <Input placeholder="Instance (例如: localhost:9090)" />
                                </Form.Item>
                                <Form.Item name="job" noStyle>
                                    <Input placeholder="Job (例如: node_exporter)" />
                                </Form.Item>
                            </Space>
                        </Form.Item>

                        <Form.Item
                            label="租户ID"
                            name="tenantId"
                        >
                            <Input placeholder="默认: default" />
                        </Form.Item>

                        <Form.Item
                            label="故障中心ID"
                            name="faultCenterId"
                        >
                            <Input placeholder="默认: default" />
                        </Form.Item>

                        <Form.Item
                            label="持续时间（秒）"
                            name="duration"
                            tooltip="告警持续多久后触发状态转换"
                        >
                            <InputNumber min={1} max={600} style={{ width: "100%" }} />
                        </Form.Item>

                        <Form.Item
                            label="自动创建工单"
                            name="autoCreateTicket"
                            valuePropName="checked"
                            tooltip="告警触发后自动创建工单"
                        >
                            <Switch />
                        </Form.Item>

                        <Form.Item
                            label="自动恢复"
                            name="autoRecover"
                            valuePropName="checked"
                        >
                            <Switch />
                        </Form.Item>

                        <Form.Item
                            label="恢复时间（秒）"
                            name="recoverAfter"
                            tooltip="自动恢复的等待时间"
                        >
                            <InputNumber min={1} max={3600} style={{ width: "100%" }} />
                        </Form.Item>

                        <Form.Item>
                            <Space>
                                <Button
                                    type="primary"
                                    icon={<PlayCircleOutlined />}
                                    loading={loading}
                                    htmlType="submit"
                                >
                                    创建告警
                                </Button>
                                <Button
                                    icon={<DeleteOutlined />}
                                    onClick={cleanupAlerts}
                                >
                                    清理数据
                                </Button>
                            </Space>
                        </Form.Item>
                    </Form>
                </Card>

                <Card title="告警状态" bordered={false}>
                    {lastAlert ? (
                        <Space direction="vertical" style={{ width: "100%" }}>
                            <AntAlert
                                message="最新告警"
                                description={`告警ID: ${lastAlert.event_id}`}
                                type="info"
                                showIcon
                            />

                            <Descriptions bordered size="small">
                                <Descriptions.Item label="规则名称" span={3}>
                                    {lastAlert.rule_name}
                                </Descriptions.Item>
                                <Descriptions.Item label="严重程度" span={2}>
                                    <Tag color={severityColors[lastAlert.severity]}>
                                        {lastAlert.severity}
                                    </Tag>
                                </Descriptions.Item>
                                <Descriptions.Item label="状态" span={1}>
                                    <Tag color={statusColors[lastAlert.status]}>
                                        {statusText[lastAlert.status] || lastAlert.status}
                                    </Tag>
                                </Descriptions.Item>
                                <Descriptions.Item label="Event ID" span={3}>
                                    <code>{lastAlert.event_id}</code>
                                </Descriptions.Item>
                            </Descriptions>

                            <Button
                                icon={<ReloadOutlined />}
                                onClick={() => recoverAlert(lastAlert.event_id)}
                                disabled={lastAlert.status === "recovered"}
                                block
                            >
                                恢复告警
                            </Button>

                            {alertHistory.length > 0 && (
                                <div>
                                    <h4 style={{ marginTop: "16px" }}>历史记录</h4>
                                    {alertHistory.map((alert, index) => (
                                        <Card
                                            key={index}
                                            size="small"
                                            style={{ marginBottom: "8px" }}
                                        >
                                            <Space>
                                                <Tag color={severityColors[alert.severity]}>
                                                    {alert.severity}
                                                </Tag>
                                                <Tag color={statusColors[alert.status]}>
                                                    {statusText[alert.status] || alert.status}
                                                </Tag>
                                                <span>{alert.rule_name}</span>
                                            </Space>
                                        </Card>
                                    ))}
                                </div>
                            )}
                        </Space>
                    ) : (
                        <div style={{ textAlign: "center", padding: "40px", color: "#999" }}>
                            暂无告警数据
                        </div>
                    )}
                </Card>
            </div>

            <Card
                title="功能说明"
                bordered={false}
                style={{ marginTop: "24px" }}
            >
                <Space direction="vertical" size="middle">
                    <div>
                        <strong>告警流程：</strong>
                        <ol style={{ marginLeft: "20px" }}>
                            <li>创建模拟告警后，系统会将其推送到故障中心</li>
                            <li>等待持续时间后，告警状态从"预告警"转换为"告警中"</li>
                            <li>如果启用"自动创建工单"，系统会根据告警转工单规则自动创建工单</li>
                            <li>系统会发送告警通知（通过配置的通知渠道）</li>
                            <li>如果启用"自动恢复"，告警会在指定时间后自动恢复</li>
                        </ol>
                    </div>
                    <div>
                        <strong>工单自动创建条件：</strong>
                        <ul style={{ marginLeft: "20px" }}>
                            <li>告警状态为"告警中"（alerting）</li>
                            <li>匹配到有效的告警转工单规则</li>
                            <li>规则配置为自动创建工单</li>
                        </ul>
                    </div>
                    <div>
                        <strong>注意事项：</strong>
                        <ul style={{ marginLeft: "20px" }}>
                            <li>模拟告警使用 "mock-rule-*" 前缀标识</li>
                            <li>使用"清理数据"按钮可以删除所有模拟告警和相关工单</li>
                            <li>建议先配置好告警转工单规则再进行测试</li>
                        </ul>
                    </div>
                </Space>
            </Card>
        </div>
    )
}