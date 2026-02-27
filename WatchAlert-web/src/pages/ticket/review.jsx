"use client"

import { useState, useEffect } from "react"
import {
    Card,
    Table,
    Button,
    Space,
    Tag,
    Rate,
    Modal,
    Form,
    InputNumber,
    Input,
    Select,
    message,
    Row,
    Col,
    Statistic,
    Tabs,
    Empty,
} from "antd"
import {
    PlusOutlined,
    EditOutlined,
    CheckOutlined,
    FileTextOutlined,
} from "@ant-design/icons"
import {
    getReviews,
    submitReview,
    assignReviewers,
} from "../../api/ticket_review"
import { getUserList } from "../../api/user"
import { getTicketList } from "../../api/ticket"
import { HandleApiError, FormatTime } from "../../utils/lib"
import { clearCacheByUrl } from "../../utils/http"
import { useNavigate } from "react-router-dom"

const { TextArea } = Input

export const TicketReview = () => {
    const navigate = useNavigate()
    const [form] = Form.useForm()
    const [assignForm] = Form.useForm()
    const [loading, setLoading] = useState(false)
    const [resolvedTickets, setResolvedTickets] = useState([])
    const [reviews, setReviews] = useState([])
    const [users, setUsers] = useState([])
    const [submitModalVisible, setSubmitModalVisible] = useState(false)
    const [assignModalVisible, setAssignModalVisible] = useState(false)
    const [currentReview, setCurrentReview] = useState(null)
    const [selectedTicketId, setSelectedTicketId] = useState(null)
    const [activeTab, setActiveTab] = useState("pending")

    useEffect(() => {
        fetchUsers()
        fetchResolvedTickets()
    }, [])

    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search)
        const id = urlParams.get('ticketId')
        if (id) {
            setSelectedTicketId(id)
            fetchReviews(id)
        }
    }, [])

    useEffect(() => {
        if (selectedTicketId) {
            fetchReviews(selectedTicketId)
        }
    }, [selectedTicketId])

    const fetchResolvedTickets = async () => {
        setLoading(true)
        try {
            const res = await getTicketList({ status: "Resolved", page: 1, size: 100 }, { skipCache: true })
            if (res && res.data) {
                setResolvedTickets(res.data.list || [])
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setLoading(false)
        }
    }

    const fetchReviews = async (ticketId) => {
        if (!ticketId) return
        setLoading(true)
        try {
            const res = await getReviews({ ticketId, page: 1, size: 100 }, { skipCache: true })
            if (res && res.data) {
                setReviews(res.data.list || [])
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setLoading(false)
        }
    }

    const fetchUsers = async () => {
        try {
            const res = await getUserList({ page: 1, size: 100 }, { skipCache: true })
            if (res && res.data) {
                setUsers(Array.isArray(res.data) ? res.data : [])
            }
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleSubmitReview = async (values) => {
        try {
            await submitReview({
                reviewId: currentReview.reviewId,
                rating: values.rating,
                workHours: values.workHours,
                comment: values.comment,
            })
            message.success("评审提交成功")
            setSubmitModalVisible(false)
            form.resetFields()
            setCurrentReview(null)
            clearCacheByUrl('/api/w8t/ticket/review')
            if (selectedTicketId) {
                await fetchReviews(selectedTicketId)
            }
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleAssignReviewers = async (values) => {
        if (!selectedTicketId) {
            message.warning("请先选择一个工单")
            return
        }
        try {
            await assignReviewers({
                ticketId: selectedTicketId,
                reviewerIds: values.reviewerIds,
            })
            message.success("评委分配成功")
            setAssignModalVisible(false)
            assignForm.resetFields()
            clearCacheByUrl('/api/w8t/ticket/review')
            await fetchReviews(selectedTicketId)
        } catch (error) {
            HandleApiError(error)
        }
    }

    const openSubmitModal = (review) => {
        setCurrentReview(review)
        form.setFieldsValue({
            rating: review.rating || 0,
            workHours: review.workHours || 0,
            comment: review.comment || "",
        })
        setSubmitModalVisible(true)
    }

    const getUserName = (userId) => {
        const user = users.find(u => u.userid === userId)
        return user ? (user.username || userId) : userId
    }

    const getStatusTag = (status) => {
        const statusMap = {
            pending: { color: "default", text: "待评审" },
            completed: { color: "success", text: "已完成" },
        }
        const config = statusMap[status] || { color: "default", text: status }
        return <Tag color={config.color}>{config.text}</Tag>
    }

    const handleSelectTicket = (ticketId) => {
        setSelectedTicketId(ticketId)
        setActiveTab("review")
        // 更新URL
        const url = new URL(window.location)
        url.searchParams.set('ticketId', ticketId)
        window.history.pushState({}, '', url)
        fetchReviews(ticketId)
    }

    // 待评审工单列表
    const pendingTicketsColumns = [
        {
            title: "工单编号",
            dataIndex: "ticketNo",
            key: "ticketNo",
            width: 150,
        },
        {
            title: "标题",
            dataIndex: "title",
            key: "title",
            ellipsis: true,
        },
        {
            title: "优先级",
            dataIndex: "priority",
            key: "priority",
            width: 100,
            render: (priority) => {
                const priorityMap = {
                    P0: { color: "red", text: "P0" },
                    P1: { color: "orange", text: "P1" },
                    P2: { color: "blue", text: "P2" },
                    P3: { color: "green", text: "P3" },
                    P4: { color: "default", text: "P4" },
                }
                const config = priorityMap[priority] || { color: "default", text: priority }
                return <Tag color={config.color}>{config.text}</Tag>
            }
        },
        {
            title: "处理人",
            dataIndex: "assignedTo",
            key: "assignedTo",
            width: 120,
            render: (userId) => getUserName(userId),
        },
        {
            title: "解决时间",
            dataIndex: "resolvedAt",
            key: "resolvedAt",
            width: 170,
            render: (time) => FormatTime(time),
        },
        {
            title: "操作",
            key: "action",
            width: 100,
            render: (_, record) => (
                <Button
                    type="primary"
                    size="small"
                    onClick={() => handleSelectTicket(record.ticketId)}
                >
                    评审
                </Button>
            ),
        },
    ]

    // 评审记录列表
    const reviewColumns = [
        {
            title: "评委",
            dataIndex: "reviewerId",
            key: "reviewerId",
            render: (reviewerId) => getUserName(reviewerId),
        },
        {
            title: "状态",
            dataIndex: "status",
            key: "status",
            render: (status) => getStatusTag(status),
        },
        {
            title: "评分",
            dataIndex: "rating",
            key: "rating",
            render: (rating) => (rating ? <Rate disabled value={rating} /> : "-"),
        },
        {
            title: "工时（小时）",
            dataIndex: "workHours",
            key: "workHours",
            render: (workHours) => workHours || "-",
        },
        {
            title: "评语",
            dataIndex: "comment",
            key: "comment",
            ellipsis: true,
            render: (comment) => comment || "-",
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
            render: (_, record) => (
                <Space>
                    {record.status === "pending" && (
                        <Button
                            type="primary"
                            size="small"
                            icon={<CheckOutlined />}
                            onClick={() => openSubmitModal(record)}
                        >
                            提交评审
                        </Button>
                    )}
                    {record.status === "completed" && (
                        <Button
                            size="small"
                            icon={<EditOutlined />}
                            onClick={() => openSubmitModal(record)}
                        >
                            查看
                        </Button>
                    )}
                </Space>
            ),
        },
    ]

    const completedReviews = reviews.filter(r => r.status === "completed")
    const avgRating = completedReviews.length > 0
        ? completedReviews.reduce((sum, r) => sum + r.rating, 0) / completedReviews.length
        : 0
    const totalWorkHours = completedReviews.reduce((sum, r) => sum + r.workHours, 0)
    const avgWorkHours = completedReviews.length > 0
        ? totalWorkHours / completedReviews.length
        : 0

    const selectedTicket = resolvedTickets.find(t => t.ticketId === selectedTicketId)

    return (
        <div style={{ padding: "24px" }}>
            <Tabs
                activeKey={activeTab}
                onChange={setActiveTab}
                items={[
                    {
                        key: "pending",
                        label: "待评审工单",
                        icon: <FileTextOutlined />,
                        children: (
                            <Card>
                                <Table
                                    columns={pendingTicketsColumns}
                                    dataSource={resolvedTickets}
                                    loading={loading}
                                    rowKey="ticketId"
                                    pagination={{
                                        showSizeChanger: true,
                                        showTotal: (total) => `共 ${total} 条`,
                                    }}
                                    locale={{
                                        emptyText: (
                                            <Empty
                                                description="暂无待评审工单"
                                                subDescription="工单状态变为「已解决」后，可以在此进行评审"
                                            />
                                        )
                                    }}
                                />
                            </Card>
                        ),
                    },
                    {
                        key: "review",
                        label: selectedTicket ? `评审: ${selectedTicket.title}` : "评审详情",
                        children: selectedTicketId ? (
                            <>
                                <Row gutter={24} style={{ marginBottom: 24 }}>
                                    <Col span={6}>
                                        <Card>
                                            <Statistic
                                                title="评审总数"
                                                value={reviews.length}
                                                suffix="个"
                                            />
                                        </Card>
                                    </Col>
                                    <Col span={6}>
                                        <Card>
                                            <Statistic
                                                title="已完成"
                                                value={completedReviews.length}
                                                suffix="个"
                                                valueStyle={{ color: "#3f8600" }}
                                            />
                                        </Card>
                                    </Col>
                                    <Col span={6}>
                                        <Card>
                                            <Statistic
                                                title="平均评分"
                                                value={avgRating}
                                                precision={1}
                                                suffix="/ 5"
                                                valueStyle={{ color: "#cf1322" }}
                                            />
                                        </Card>
                                    </Col>
                                    <Col span={6}>
                                        <Card>
                                            <Statistic
                                                title="平均工时"
                                                value={avgWorkHours}
                                                precision={1}
                                                suffix="小时"
                                                valueStyle={{ color: "#1890ff" }}
                                            />
                                        </Card>
                                    </Col>
                                </Row>

                                <Card
                                    title={`工单: ${selectedTicket?.title || ''}`}
                                    extra={
                                        <Space>
                                            <Button onClick={() => navigate(`/ticket/detail/${selectedTicketId}`)}>
                                                查看工单详情
                                            </Button>
                                            <Button
                                                type="primary"
                                                icon={<PlusOutlined />}
                                                onClick={() => setAssignModalVisible(true)}
                                            >
                                                分配评委
                                            </Button>
                                        </Space>
                                    }
                                >
                                    <Table
                                        columns={reviewColumns}
                                        dataSource={reviews}
                                        loading={loading}
                                        rowKey="reviewId"
                                        pagination={false}
                                        locale={{
                                            emptyText: (
                                                <Empty
                                                    description="暂无评审记录"
                                                    subDescription="点击「分配评委」添加评审人员"
                                                />
                                            )
                                        }}
                                    />
                                </Card>
                            </>
                        ) : (
                            <Card>
                                <Empty
                                    description="请先选择一个待评审的工单"
                                    subDescription="从「待评审工单」标签页选择工单进行评审"
                                />
                            </Card>
                        ),
                    },
                ]}
            />

            {/* 提交评审弹窗 */}
            <Modal
                title="提交评审"
                open={submitModalVisible}
                onCancel={() => {
                    setSubmitModalVisible(false)
                    form.resetFields()
                    setCurrentReview(null)
                }}
                onOk={() => form.submit()}
                width={600}
            >
                <Form form={form} layout="vertical" onFinish={handleSubmitReview}>
                    <Form.Item
                        name="rating"
                        label="评分"
                        rules={[{ required: true, message: "请选择评分" }]}
                    >
                        <Rate />
                    </Form.Item>
                    <Form.Item
                        name="workHours"
                        label="工时（小时）"
                        rules={[{ required: true, message: "请输入工时" }]}
                    >
                        <InputNumber min={0} step={0.5} style={{ width: "100%" }} />
                    </Form.Item>
                    <Form.Item name="comment" label="评语">
                        <TextArea
                            rows={4}
                            placeholder="请输入评语和改进建议"
                        />
                    </Form.Item>
                </Form>
            </Modal>

            {/* 分配评委弹窗 */}
            <Modal
                title="分配评委"
                open={assignModalVisible}
                onCancel={() => {
                    setAssignModalVisible(false)
                    assignForm.resetFields()
                }}
                onOk={() => assignForm.submit()}
            >
                <Form form={assignForm} layout="vertical" onFinish={handleAssignReviewers}>
                    <Form.Item
                        name="reviewerIds"
                        label="选择评委"
                        rules={[{ required: true, message: "请选择评委" }]}
                    >
                        <Select
                            mode="multiple"
                            placeholder="请选择评委"
                            options={users.map(u => ({
                                label: u.username || u.userid,
                                value: u.userid,
                            }))}
                            filterOption={(input, option) =>
                                option.label.toLowerCase().includes(input.toLowerCase())
                            }
                        />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    )
}

export default TicketReview