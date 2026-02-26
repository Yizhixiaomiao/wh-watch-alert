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
    Switch,
    message,
    Popconfirm,
} from "antd"
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    UserOutlined,
} from "@ant-design/icons"
import {
    getReviewers,
    createReviewer,
    updateReviewer,
    deleteReviewer,
} from "../../api/ticket_review"
import { HandleApiError, FormatTime } from "../../utils/lib"
import { clearCacheByUrl } from "../../utils/http"

export const ReviewerManage = () => {
    const [form] = Form.useForm()
    const [loading, setLoading] = useState(false)
    const [reviewers, setReviewers] = useState([])
    const [modalVisible, setModalVisible] = useState(false)
    const [editingReviewer, setEditingReviewer] = useState(null)

    useEffect(() => {
        fetchReviewers()
    }, [])

    const fetchReviewers = async () => {
        setLoading(true)
        try {
            const res = await getReviewers({ page: 1, size: 100 }, { skipCache: true })
            if (res && res.data) {
                setReviewers(res.data.list || [])
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setLoading(false)
        }
    }

    const handleCreate = async (values) => {
        try {
            await createReviewer({
                ...values,
            })
            message.success('评委创建成功')
            setModalVisible(false)
            form.resetFields()
            setEditingReviewer(null)
            clearCacheByUrl('/api/w8t/ticket/reviewer')
            clearCacheByUrl('/api/w8t/ticket/reviewer/list')
            await fetchReviewers()
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleUpdate = async (values) => {
        try {
            await updateReviewer({
                ...values,
                reviewerId: editingReviewer.reviewerId,
            })
            message.success('评委更新成功')
            setModalVisible(false)
            form.resetFields()
            setEditingReviewer(null)
            clearCacheByUrl('/api/w8t/ticket/reviewer')
            clearCacheByUrl('/api/w8t/ticket/reviewer/list')
            await fetchReviewers()
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleDelete = async (reviewerId) => {
        try {
            await deleteReviewer({
                reviewerId,
            })
            message.success('评委删除成功')
            clearCacheByUrl('/api/w8t/ticket/reviewer')
            clearCacheByUrl('/api/w8t/ticket/reviewer/list')
            await fetchReviewers()
        } catch (error) {
            HandleApiError(error)
        }
    }

    const openModal = (reviewer = null) => {
        setEditingReviewer(reviewer)
        if (reviewer) {
            form.setFieldsValue({
                userName: reviewer.userName,
                email: reviewer.email,
                phone: reviewer.phone,
                department: reviewer.department,
                specialty: reviewer.specialty,
                isActive: reviewer.isActive,
            })
        } else {
            form.resetFields()
        }
        setModalVisible(true)
    }

    const handleSubmit = async (values) => {
        if (editingReviewer) {
            await handleUpdate(values)
        } else {
            await handleCreate(values)
        }
    }

    const getStatusTag = (isActive) => {
        return isActive ? (
            <Tag color="success">启用</Tag>
        ) : (
            <Tag color="default">禁用</Tag>
        )
    }

    const columns = [
        {
            title: "评委ID",
            dataIndex: "reviewerId",
            key: "reviewerId",
            width: 200,
            ellipsis: true,
        },
        {
            title: "姓名",
            dataIndex: "userName",
            key: "userName",
        },
        {
            title: "邮箱",
            dataIndex: "email",
            key: "email",
            ellipsis: true,
        },
        {
            title: "电话",
            dataIndex: "phone",
            key: "phone",
        },
        {
            title: "部门",
            dataIndex: "department",
            key: "department",
        },
        {
            title: "专业领域",
            dataIndex: "specialty",
            key: "specialty",
        },
        {
            title: "状态",
            dataIndex: "isActive",
            key: "isActive",
            render: (isActive) => getStatusTag(isActive),
        },
        {
            title: "创建时间",
            dataIndex: "createdAt",
            key: "createdAt",
            render: (createdAt) => FormatTime(createdAt),
        },
        {
            title: "操作",
            key: "action",
            width: 150,
            fixed: 'right',
            render: (_, record) => (
                <Space>
                    <Button
                        type="link"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => openModal(record)}
                    >
                        编辑
                    </Button>
                    <Popconfirm
                        title="确认删除"
                        description="确定要删除这个评委吗？"
                        onConfirm={() => handleDelete(record.reviewerId)}
                        okText="确定"
                        cancelText="取消"
                    >
                        <Button
                            type="link"
                            size="small"
                            danger
                            icon={<DeleteOutlined />}
                        >
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ]

    return (
        <div style={{ padding: "24px" }}>
            <Card
                title="评委管理"
                extra={
                    <Button
                        type="primary"
                        icon={<PlusOutlined />}
                        onClick={() => openModal()}
                    >
                        新增评委
                    </Button>
                }
            >
                <Table
                    columns={columns}
                    dataSource={reviewers}
                    loading={loading}
                    rowKey="reviewerId"
                    scroll={{ x: 1200 }}
                    pagination={{
                        showSizeChanger: true,
                        showTotal: (total) => `共 ${total} 条`,
                    }}
                />
            </Card>

            <Modal
                title={editingReviewer ? "编辑评委" : "新增评委"}
                open={modalVisible}
                onCancel={() => {
                    setModalVisible(false)
                    form.resetFields()
                    setEditingReviewer(null)
                }}
                onOk={() => form.submit()}
                width={600}
                destroyOnClose
            >
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={handleSubmit}
                >
                    <Form.Item
                        name="userName"
                        label="评委姓名"
                        rules={[
                            { required: true, message: "请输入评委姓名" },
                            { max: 50, message: "姓名不能超过50个字符" }
                        ]}
                    >
                        <Input
                            prefix={<UserOutlined />}
                            placeholder="请输入评委姓名"
                        />
                    </Form.Item>
                    <Form.Item
                        name="email"
                        label="邮箱"
                        rules={[
                            { type: 'email', message: "请输入有效的邮箱地址" },
                            { max: 100, message: "邮箱不能超过100个字符" }
                        ]}
                    >
                        <Input placeholder="请输入邮箱地址" />
                    </Form.Item>
                    <Form.Item
                        name="phone"
                        label="电话"
                        rules={[
                            { pattern: /^1[3-9]\d{9}$/, message: "请输入有效的手机号码" }
                        ]}
                    >
                        <Input placeholder="请输入手机号码" />
                    </Form.Item>
                    <Form.Item
                        name="department"
                        label="部门"
                        rules={[
                            { max: 50, message: "部门不能超过50个字符" }
                        ]}
                    >
                        <Input placeholder="请输入所属部门" />
                    </Form.Item>
                    <Form.Item
                        name="specialty"
                        label="专业领域"
                        rules={[
                            { max: 100, message: "专业领域不能超过100个字符" }
                        ]}
                    >
                        <Input placeholder="请输入专业领域，如：网络、数据库、应用等" />
                    </Form.Item>
                    <Form.Item
                        name="isActive"
                        label="状态"
                        valuePropName="checked"
                        initialValue={true}
                    >
                        <Switch checkedChildren="启用" unCheckedChildren="禁用" />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    )
}

export default ReviewerManage
