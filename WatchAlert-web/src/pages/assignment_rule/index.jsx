"use client"

import { useState, useEffect } from "react"
import {
    Card,
    Table,
    Button,
    Space,
    Tag,
    Modal,
    Form,
    Input,
    Select,
    Switch,
    InputNumber,
    message,
} from "antd"
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    ThunderboltOutlined,
} from "@ant-design/icons"
import {
    getAssignmentRules,
    createAssignmentRule,
    updateAssignmentRule,
    deleteAssignmentRule,
} from "../../api/assignment_rule"
import { HandleApiError, FormatTime } from "../../utils/lib"
import { getUserList } from "../../api/user"
import { getDutyManagerList as getDutyList } from "../../api/duty"
import { clearCacheByUrl } from "../../utils/http"

const { TextArea } = Input

export const AssignmentRule = () => {
    const [form] = Form.useForm()
    const [loading, setLoading] = useState(false)
    const [rules, setRules] = useState([])
    const [users, setUsers] = useState([])
    const [duties, setDuties] = useState([])
    const [modalVisible, setModalVisible] = useState(false)
    const [editingRule, setEditingRule] = useState(null)

    useEffect(() => {
        fetchRules()
        fetchUsers()
        fetchDuties()
    }, [])

    const fetchRules = async () => {
        setLoading(true)
        try {
            const res = await getAssignmentRules({ page: 1, size: 100 }, { skipCache: true })
            if (res && res.data) {
                setRules(res.data.list || [])
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setLoading(false)
        }
    }

    const fetchUsers = async () => {
        try {
            const res = await getUserList({ page: 1, size: 100 })
            if (res && res.data) {
                setUsers(Array.isArray(res.data) ? res.data : [])
            }
        } catch (error) {
            console.error("获取用户列表失败:", error)
        }
    }

    const fetchDuties = async () => {
        try {
            const res = await getDutyList({ page: 1, size: 100 })
            if (res && res.data) {
                setDuties(Array.isArray(res.data) ? res.data : [])
            }
        } catch (error) {
            console.error("获取值班组列表失败:", error)
        }
    }

    const handleCreate = async (values) => {
        try {
            const res = await createAssignmentRule({
                ...values,
            })
            message.success('分配规则创建成功')
            setModalVisible(false)
            form.resetFields()
            setEditingRule(null)
            clearCacheByUrl('/api/w8t/assignment-rule')
            clearCacheByUrl('/api/w8t/assignment-rule/list')
            await fetchRules()
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleUpdate = async (values) => {
        try {
            const res = await updateAssignmentRule({
                ...values,
                ruleId: editingRule.ruleId,
            })
            if (res && res.code === 200) {
                message.success('分配规则更新成功')
                setModalVisible(false)
                form.resetFields()
                setEditingRule(null)
                clearCacheByUrl('/api/w8t/assignment-rule')
                await fetchRules()
            }
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleDelete = (rule) => {
        Modal.confirm({
            title: "确认删除",
            content: "确定要删除这个分配规则吗？",
            onOk: async () => {
                try {
                    const res = await deleteAssignmentRule({
                        ruleId: rule.ruleId,
                    })
                    if (res && res.code === 200) {
                        message.success('分配规则删除成功')
                        clearCacheByUrl('/api/w8t/assignment-rule')
                        await fetchRules()
                    }
                } catch (error) {
                    HandleApiError(error)
                }
            },
        })
    }

    const openCreateModal = () => {
        setEditingRule(null)
        form.resetFields()
        setModalVisible(true)
    }

    const openEditModal = (rule) => {
        setEditingRule(rule)
        form.setFieldsValue(rule)
        setModalVisible(true)
    }

    const columns = [
        {
            title: "规则名称",
            dataIndex: "name",
            key: "name",
            render: (name) => <span style={{ fontWeight: 500 }}>{name}</span>,
        },
        {
            title: "规则类型",
            dataIndex: "ruleType",
            key: "ruleType",
            render: (ruleType) => {
                const typeMap = {
                    alert_type: { color: "blue", text: "告警类型" },
                    duty_schedule: { color: "green", text: "值班表" },
                }
                const config = typeMap[ruleType] || { color: "default", text: ruleType }
                return <Tag color={config.color}>{config.text}</Tag>
            },
        },
        {
            title: "匹配条件",
            key: "condition",
            render: (_, record) => (
                <Space size="small">
                    {record.alertType && <Tag color="blue">{record.alertType}</Tag>}
                    {record.dataSource && <Tag color="cyan">{record.dataSource}</Tag>}
                    {record.severity && <Tag color="red">{record.severity}</Tag>}
                </Space>
            ),
        },
        {
            title: "分配类型",
            dataIndex: "assignmentType",
            key: "assignmentType",
            render: (assignmentType) => {
                const typeMap = {
                    user: { text: "用户" },
                    group: { text: "组" },
                    duty: { text: "值班组" },
                }
                const config = typeMap[assignmentType] || { text: assignmentType }
                return <Tag>{config.text}</Tag>
            },
        },
        {
            title: "分配目标",
            key: "target",
            render: (_, record) => {
                if (record.assignmentType === "user") {
                    const user = users.find(u => u.userid === record.targetUserId)
                    return user ? user.username : record.targetUserId
                } else if (record.assignmentType === "group") {
                    return record.targetGroupId
                } else if (record.assignmentType === "duty") {
                    const duty = duties.find(d => d.id === record.targetDutyId)
                    return duty ? duty.name : record.targetDutyId
                }
                return "-"
            },
        },
        {
            title: "优先级",
            dataIndex: "priority",
            key: "priority",
            sorter: (a, b) => a.priority - b.priority,
            render: (priority) => (
                <Tag color={priority >= 100 ? "red" : "default"}>
                    {priority}
                </Tag>
            ),
        },
        {
            title: "状态",
            dataIndex: "enabled",
            key: "enabled",
            render: (enabled) => (
                <Tag color={enabled ? "success" : "default"}>
                    {enabled ? "启用" : "禁用"}
                </Tag>
            ),
        },
        {
            title: "操作",
            key: "action",
            render: (_, record) => (
                <Space>
                    <Button
                        type="link"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => openEditModal(record)}
                    >
                        编辑
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        danger
                        icon={<DeleteOutlined />}
                        onClick={() => handleDelete(record)}
                    >
                        删除
                    </Button>
                </Space>
            ),
        },
    ]

    return (
        <div style={{ padding: "24px" }}>
            <Card
                title="智能派单规则"
                extra={
                    <Button
                        type="primary"
                        icon={<PlusOutlined />}
                        onClick={openCreateModal}
                    >
                        添加规则
                    </Button>
                }
            >
                <Table
                    columns={columns}
                    dataSource={rules}
                    loading={loading}
                    rowKey="ruleId"
                    pagination={false}
                />
            </Card>

            {/* 添加/编辑弹窗 */}
            <Modal
                title={editingRule ? "编辑规则" : "添加规则"}
                open={modalVisible}
                onCancel={() => {
                    setModalVisible(false)
                    form.resetFields()
                    setEditingRule(null)
                }}
                onOk={() => form.submit()}
                width={700}
            >
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={editingRule ? handleUpdate : handleCreate}
                    initialValues={{
                        ruleType: "duty_schedule",
                        assignmentType: "duty",
                        priority: 1,
                        enabled: true,
                    }}
                >
                    <Form.Item
                        name="name"
                        label="规则名称"
                        rules={[{ required: true, message: "请输入规则名称" }]}
                    >
                        <Input placeholder="请输入规则名称" />
                    </Form.Item>

                    <Form.Item
                        name="ruleType"
                        label="规则类型"
                        rules={[{ required: true, message: "请选择规则类型" }]}
                    >
                        <Select placeholder="请选择规则类型">
                            <Select.Option value="alert_type">告警类型规则</Select.Option>
                            <Select.Option value="duty_schedule">值班表规则</Select.Option>
                        </Select>
                    </Form.Item>

                    <Form.Item noStyle shouldUpdate={(prevValues, currentValues) => {
                        return prevValues.ruleType !== currentValues.ruleType
                    }}>
                        {({ getFieldValue }) =>
                            getFieldValue("ruleType") === "alert_type" && (
                                <>
                                    <Form.Item
                                        name="alertType"
                                        label="告警类型"
                                    >
                                        <Input placeholder="请输入告警类型，如：CPU、Memory、Disk" />
                                    </Form.Item>
                                    <Form.Item
                                        name="dataSource"
                                        label="数据源类型"
                                    >
                                        <Input placeholder="请输入数据源类型" />
                                    </Form.Item>
                                    <Form.Item
                                        name="severity"
                                        label="告警级别"
                                    >
                                        <Select placeholder="请选择告警级别">
                                            <Select.Option value="Critical">Critical</Select.Option>
                                            <Select.Option value="High">High</Select.Option>
                                            <Select.Option value="Medium">Medium</Select.Option>
                                            <Select.Option value="Low">Low</Select.Option>
                                        </Select>
                                    </Form.Item>
                                </>
                            )
                        }
                    </Form.Item>

                    <Form.Item
                        name="assignmentType"
                        label="分配类型"
                        rules={[{ required: true, message: "请选择分配类型" }]}
                    >
                        <Select placeholder="请选择分配类型">
                            <Select.Option value="user">分配给用户</Select.Option>
                            <Select.Option value="group">分配给组</Select.Option>
                            <Select.Option value="duty">分配给值班组</Select.Option>
                        </Select>
                    </Form.Item>

                    <Form.Item noStyle shouldUpdate={(prevValues, currentValues) => {
                        return prevValues.assignmentType !== currentValues.assignmentType
                    }}>
                        {({ getFieldValue }) => {
                            if (getFieldValue("assignmentType") === "user") {
                                return (
                                    <Form.Item
                                        name="targetUserId"
                                        label="选择用户"
                                        rules={[{ required: true, message: "请选择用户" }]}
                                    >
                                        <Select
                                            placeholder="请选择用户"
                                            showSearch
                                            filterOption={(input, option) =>
                                                (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                                            }
                                            options={users.map(user => ({
                                                label: user.username || user.userid,
                                                value: user.userid,
                                            }))}
                                        />
                                    </Form.Item>
                                )
                            } else if (getFieldValue("assignmentType") === "group") {
                                return (
                                    <Form.Item
                                        name="targetGroupId"
                                        label="目标组"
                                        rules={[{ required: true, message: "请输入目标组" }]}
                                    >
                                        <Input placeholder="请输入目标组ID" />
                                    </Form.Item>
                                )
                            } else if (getFieldValue("assignmentType") === "duty") {
                                return (
                                    <Form.Item
                                        name="targetDutyId"
                                        label="选择值班组"
                                        rules={[{ required: true, message: "请选择值班组" }]}
                                    >
                                        <Select
                                            placeholder="请选择值班组"
                                            showSearch
                                            filterOption={(input, option) =>
                                                (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                                            }
                                            options={duties.map(duty => ({
                                                label: duty.name || duty.id,
                                                value: duty.id,
                                            }))}
                                        />
                                    </Form.Item>
                                )
                            }
                            return null
                        }}
                    </Form.Item>

                    <Form.Item
                        name="priority"
                        label="优先级"
                        rules={[{ required: true, message: "请输入优先级" }]}
                    >
                        <InputNumber
                            min={1}
                            max={1000}
                            style={{ width: "100%" }}
                            placeholder="数值越小优先级越高"
                        />
                    </Form.Item>

                    <Form.Item
                        name="enabled"
                        label="是否启用"
                        valuePropName="checked"
                    >
                        <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    )
}

export default AssignmentRule