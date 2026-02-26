"use client"

import { useState, useEffect } from "react"
import {
    Button,
    Table,
    Space,
    Input,
    Modal,
    Form,
    message,
    Popconfirm,
    Tag,
    Select,
    InputNumber,
} from "antd"
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
} from "@ant-design/icons"
import { TableWithPagination } from "../../utils/TableWithPagination"

const { TextArea } = Input

export const TicketSlaPolicy = () => {
    const [list, setList] = useState([])
    const [loading, setLoading] = useState(false)
    const [modalVisible, setModalVisible] = useState(false)
    const [isEdit, setIsEdit] = useState(false)
    const [currentRecord, setCurrentRecord] = useState(null)
    const [form] = Form.useForm()
    const [pagination, setPagination] = useState({
        index: 1,
        size: 20,
        total: 0,
    })

    // SLA状态映射
    const statusMap = {
        Active: { color: "success", text: "启用" },
        Inactive: { color: "default", text: "禁用" },
    }

    // 优先级映射
    const priorityMap = {
        P0: { color: "red", text: "P0-最高" },
        P1: { color: "orange", text: "P1-高" },
        P2: { color: "blue", text: "P2-中" },
        P3: { color: "green", text: "P3-低" },
        P4: { color: "default", text: "P4-最低" },
    }

    const fetchList = async () => {
        setLoading(true)
        try {
            const response = await fetch('/api/w8t/ticket/sla/list', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('Authorization')}`,
                    'TenantID': localStorage.getItem('TenantID'),
                }
            })
            const data = await response.json()
            if (data.code === 200) {
                const dataList = Array.isArray(data.data) ? data.data : (data.data?.list || [])
                setList(dataList)
                setPagination({
                    ...pagination,
                    total: dataList.length || 0,
                })
            } else {
                message.error(data.msg || '获取失败')
            }
        } catch (error) {
            message.error('获取SLA策略列表失败')
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchList()
    }, [])

    const handleAdd = () => {
        setIsEdit(false)
        setCurrentRecord(null)
        form.resetFields()
        setModalVisible(true)
    }

    const handleEdit = (record) => {
        setIsEdit(true)
        setCurrentRecord(record)
        form.setFieldsValue(record)
        setModalVisible(true)
    }

    const handleDelete = async (id) => {
        try {
            const response = await fetch('/api/w8t/ticket/sla/delete', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('Authorization')}`,
                    'TenantID': localStorage.getItem('TenantID'),
                },
                body: JSON.stringify({ id })
            })
            const data = await response.json()
            if (data.code === 200) {
                message.success('删除成功')
                fetchList()
            } else {
                message.error(data.msg || '删除失败')
            }
        } catch (error) {
            message.error('删除失败')
        }
    }

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields()
            const url = isEdit ? '/api/w8t/ticket/sla/update' : '/api/w8t/ticket/sla/create'
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem('Authorization')}`,
                    'TenantID': localStorage.getItem('TenantID'),
                },
                body: JSON.stringify(values)
            })
            const data = await response.json()
            if (data.code === 200) {
                message.success(isEdit ? '更新成功' : '创建成功')
                setModalVisible(false)
                fetchList()
            } else {
                message.error(data.msg || '操作失败')
            }
        } catch (error) {
            message.error('操作失败')
        }
    }

    const columns = [
        {
            title: '序号',
            dataIndex: 'index',
            key: 'index',
            width: 80,
            render: (text, record, index) => (pagination.index - 1) * pagination.size + index + 1,
        },
        {
            title: '策略名称',
            dataIndex: 'name',
            key: 'name',
            width: 150,
        },
        {
            title: '优先级',
            dataIndex: 'priority',
            key: 'priority',
            width: 120,
            render: (priority) => (
                <Tag color={priorityMap[priority]?.color}>
                    {priorityMap[priority]?.text || priority}
                </Tag>
            ),
        },
        {
            title: '响应时间（小时）',
            dataIndex: 'responseTimeHours',
            key: 'responseTimeHours',
            width: 150,
        },
        {
            title: '解决时间（小时）',
            dataIndex: 'resolutionTimeHours',
            key: 'resolutionTimeHours',
            width: 150,
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (status) => (
                <Tag color={statusMap[status]?.color}>
                    {statusMap[status]?.text || status}
                </Tag>
            ),
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
            ellipsis: true,
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (text) => text ? new Date(text * 1000).toLocaleString('zh-CN') : '-',
        },
        {
            title: '操作',
            key: 'action',
            width: 150,
            render: (_, record) => (
                <Space size="small">
                    <Button
                        type="link"
                        icon={<EditOutlined />}
                        size="small"
                        onClick={() => handleEdit(record)}
                    >
                        编辑
                    </Button>
                    <Popconfirm
                        title="确认删除"
                        description="确定要删除这个SLA策略吗？"
                        onConfirm={() => handleDelete(record.id)}
                        okText="确定"
                        cancelText="取消"
                    >
                        <Button
                            type="link"
                            danger
                            icon={<DeleteOutlined />}
                            size="small"
                        >
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ]

    return (
        <div style={{ padding: '24px' }}>
            <div style={{ marginBottom: '16px', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <h2 style={{ margin: 0 }}>SLA策略管理</h2>
                <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={handleAdd}
                >
                    添加SLA策略
                </Button>
            </div>

            <TableWithPagination
                columns={columns}
                dataSource={list}
                loading={loading}
                pagination={pagination}
                setPagination={setPagination}
                rowKey="id"
            />

            <Modal
                title={isEdit ? '编辑SLA策略' : '添加SLA策略'}
                open={modalVisible}
                onOk={handleSubmit}
                onCancel={() => setModalVisible(false)}
                width={600}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        label="策略名称"
                        name="name"
                        rules={[{ required: true, message: '请输入策略名称' }]}
                    >
                        <Input placeholder="请输入SLA策略名称" />
                    </Form.Item>
                    <Form.Item
                        label="优先级"
                        name="priority"
                        rules={[{ required: true, message: '请选择优先级' }]}
                    >
                        <Select placeholder="请选择优先级">
                            <Select.Option value="P0">P0-最高</Select.Option>
                            <Select.Option value="P1">P1-高</Select.Option>
                            <Select.Option value="P2">P2-中</Select.Option>
                            <Select.Option value="P3">P3-低</Select.Option>
                            <Select.Option value="P4">P4-最低</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                        label="响应时间（小时）"
                        name="responseTimeHours"
                        rules={[{ required: true, message: '请输入响应时间' }]}
                    >
                        <InputNumber
                            min={0}
                            step={1}
                            style={{ width: '100%' }}
                            placeholder="请输入响应时间（小时）"
                        />
                    </Form.Item>
                    <Form.Item
                        label="解决时间（小时）"
                        name="resolutionTimeHours"
                        rules={[{ required: true, message: '请输入解决时间' }]}
                    >
                        <InputNumber
                            min={0}
                            step={1}
                            style={{ width: '100%' }}
                            placeholder="请输入解决时间（小时）"
                        />
                    </Form.Item>
                    <Form.Item
                        label="状态"
                        name="status"
                        rules={[{ required: true, message: '请选择状态' }]}
                    >
                        <Select placeholder="请选择状态">
                            <Select.Option value="Active">启用</Select.Option>
                            <Select.Option value="Inactive">禁用</Select.Option>
                        </Select>
                    </Form.Item>
                    <Form.Item
                        label="描述"
                        name="description"
                    >
                        <TextArea rows={4} placeholder="请输入描述" />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    )
}

export default TicketSlaPolicy