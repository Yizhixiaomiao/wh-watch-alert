"use client"

import { useState, useEffect } from "react"
import {
    Button,
    Table,
    Space,
    Input,
    Modal,
    Form,
    InputNumber,
    message,
    Popconfirm,
    Tag,
} from "antd"
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
} from "@ant-design/icons"
import { TableWithPagination } from "../../utils/TableWithPagination"
import { clearCacheByUrl } from "../../utils/http"

const { TextArea } = Input

export const WorkHoursStandard = () => {
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

    const fetchList = async () => {
        setLoading(true)
        try {
            // 添加时间戳防止缓存
            const response = await fetch(`/api/w8t/work-hours/standard/list?_t=${Date.now()}`, {
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
            message.error('获取工时标准列表失败')
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
            const response = await fetch('/api/w8t/work-hours/standard/delete', {
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
                clearCacheByUrl('/api/w8t/work-hours')
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
            const url = isEdit ? '/api/w8t/work-hours/standard/update' : '/api/w8t/work-hours/standard/create'
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
                clearCacheByUrl('/api/w8t/work-hours')
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
            title: '工时类型',
            dataIndex: 'type',
            key: 'type',
            width: 150,
        },
        {
            title: '标准工时（小时）',
            dataIndex: 'standardHours',
            key: 'standardHours',
            width: 150,
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
                        description="确定要删除这个工时标准吗？"
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
                <h2 style={{ margin: 0 }}>工时标准管理</h2>
                <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={handleAdd}
                >
                    添加工时标准
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
                title={isEdit ? '编辑工时标准' : '添加工时标准'}
                open={modalVisible}
                onOk={handleSubmit}
                onCancel={() => setModalVisible(false)}
                width={600}
            >
                <Form form={form} layout="vertical">
                    <Form.Item
                        label="工时类型"
                        name="type"
                        rules={[{ required: true, message: '请输入工时类型' }]}
                    >
                        <Input placeholder="例如：代码审查、测试、部署等" />
                    </Form.Item>
                    <Form.Item
                        label="标准工时（小时）"
                        name="standardHours"
                        rules={[{ required: true, message: '请输入标准工时' }]}
                    >
                        <InputNumber
                            min={0}
                            step={0.5}
                            precision={1}
                            style={{ width: '100%' }}
                            placeholder="请输入标准工时"
                        />
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

export default WorkHoursStandard
