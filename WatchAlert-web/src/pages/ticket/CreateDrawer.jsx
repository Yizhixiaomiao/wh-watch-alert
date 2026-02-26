"use client"

import { useState, useEffect } from "react"
import {
    Form,
    Input,
    Select,
    Button,
    Space,
    Row,
    Col,
    Drawer,
} from "antd"
import { SaveOutlined } from "@ant-design/icons"
import { createTicket } from "../../api/ticket"
import { HandleApiError } from "../../utils/lib"
import { getUserList } from "../../api/user"
import { FaultCenterList } from "../../api/faultCenter"

const { TextArea } = Input

export const CreateTicketDrawer = ({ visible, onClose, onSuccess }) => {
    const [form] = Form.useForm()
    const [submitting, setSubmitting] = useState(false)
    const [userList, setUserList] = useState([])
    const [faultCenterList, setFaultCenterList] = useState([])

    useEffect(() => {
        if (visible) {
            fetchUserList()
            fetchFaultCenterList()
            // 重置表单
            form.resetFields()
        }
    }, [visible])

    const fetchUserList = async () => {
        try {
            const res = await getUserList({})
            if (res && res.data) {
                const users = Array.isArray(res.data) ? res.data : []
                setUserList(users)
            }
        } catch (error) {
            console.error("获取用户列表失败:", error)
        }
    }

    const fetchFaultCenterList = async () => {
        try {
            const res = await FaultCenterList({ page: 1, size: 100 })
            if (res && res.data) {
                const faults = res.data.list || res.data || []
                setFaultCenterList(faults)
            }
        } catch (error) {
            console.error("获取故障中心列表失败:", error)
        }
    }

    // 提交表单
    const handleSubmit = async (values) => {
        setSubmitting(true)
        try {
            const res = await createTicket(values)
            form.resetFields()
            onClose()
            if (onSuccess) {
                onSuccess(res.data)
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setSubmitting(false)
        }
    }

    return (
        <Drawer
            title="创建工单"
            open={visible}
            onClose={onClose}
            width={600}
            footer={null}
        >
            <Form
                form={form}
                layout="vertical"
                onFinish={handleSubmit}
                initialValues={{
                    type: "Fault",
                    priority: "P2",
                    severity: "Medium",
                }}
            >
                <Form.Item
                    name="title"
                    label="工单标题"
                    rules={[{ required: true, message: "请输入工单标题" }]}
                >
                    <Input placeholder="请输入工单标题" />
                </Form.Item>

                <Row gutter={16}>
                    <Col span={12}>
                        <Form.Item
                            name="type"
                            label="工单类型"
                            rules={[{ required: true, message: "请选择工单类型" }]}
                        >
                            <Select
                                placeholder="请选择工单类型"
                                options={[
                                    { value: "Alert", label: "告警工单" },
                                    { value: "Fault", label: "故障工单" },
                                    { value: "Change", label: "变更工单" },
                                    { value: "Query", label: "咨询工单" },
                                ]}
                            />
                        </Form.Item>
                    </Col>
                    <Col span={12}>
                        <Form.Item
                            name="priority"
                            label="优先级"
                            rules={[{ required: true, message: "请选择优先级" }]}
                        >
                            <Select
                                placeholder="请选择优先级"
                                options={[
                                    { value: "P0", label: "P0 - 最高优先级" },
                                    { value: "P1", label: "P1 - 高优先级" },
                                    { value: "P2", label: "P2 - 中优先级" },
                                    { value: "P3", label: "P3 - 低优先级" },
                                    { value: "P4", label: "P4 - 最低优先级" },
                                ]}
                            />
                        </Form.Item>
                    </Col>
                </Row>

                <Form.Item name="severity" label="严重程度">
                    <Select
                        placeholder="请选择严重程度"
                        options={[
                            { value: "Critical", label: "严重" },
                            { value: "High", label: "高" },
                            { value: "Medium", label: "中" },
                            { value: "Low", label: "低" },
                        ]}
                    />
                </Form.Item>

                <Form.Item
                    name="description"
                    label="工单描述"
                    rules={[{ required: true, message: "请输入工单描述" }]}
                >
                    <TextArea
                        rows={6}
                        placeholder="请详细描述问题或需求..."
                    />
                </Form.Item>

                <Form.Item name="assignedTo" label="指定处理人">
                    <Select
                        placeholder="请选择处理人（可选）"
                        allowClear
                        showSearch
                        filterOption={(input, option) =>
                            (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                        }
                        options={userList.map(user => ({
                            label: `${user.username || user.userid} ${user.email ? '(' + user.email + ')' : ''}`,
                            value: user.userid
                        }))}
                    />
                </Form.Item>

                <Form.Item name="faultCenterId" label="关联故障中心">
                    <Select
                        placeholder="请选择故障中心（可选）"
                        allowClear
                        showSearch
                        filterOption={(input, option) =>
                            (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                        }
                        options={faultCenterList.map(fault => ({
                            label: fault.name || fault.title || fault.id,
                            value: fault.id
                        }))}
                    />
                </Form.Item>

                <Form.Item style={{ marginTop: 24 }}>
                    <Space>
                        <Button
                            type="primary"
                            htmlType="submit"
                            icon={<SaveOutlined />}
                            loading={submitting}
                            style={{ backgroundColor: "#000000" }}
                        >
                            创建工单
                        </Button>
                        <Button onClick={onClose}>
                            取消
                        </Button>
                    </Space>
                </Form.Item>
            </Form>
        </Drawer>
    )
}
