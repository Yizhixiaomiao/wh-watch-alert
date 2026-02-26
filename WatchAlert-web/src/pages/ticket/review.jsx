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
} from "antd"
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    CheckOutlined,
} from "@ant-design/icons"
import {
    getReviews,
    submitReview,
    getReviewers,
    assignReviewers,
} from "../../api/ticket_review"
import { getUserList } from "../../api/user"
import { HandleApiError, FormatTime } from "../../utils/lib"

const { TextArea } = Input

export const TicketReview = () => {
    const [form] = Form.useForm()
    const [assignForm] = Form.useForm()
    const [loading, setLoading] = useState(false)
    const [reviews, setReviews] = useState([])
    const [reviewers, setReviewers] = useState([])
    const [users, setUsers] = useState([])
    const [submitModalVisible, setSubmitModalVisible] = useState(false)
    const [assignModalVisible, setAssignModalVisible] = useState(false)
    const [currentReview, setCurrentReview] = useState(null)
    const [ticketId, setTicketId] = useState(null)

    useEffect(() => {
        const urlParams = new URLSearchParams(window.location.search)
        const id = urlParams.get('ticketId')
        if (id) {
            setTicketId(id)
            fetchReviews()
            fetchReviewers()
        }
        fetchUsers()
    }, [])

    const fetchReviews = async () => {
        setLoading(true)
        try {
            const res = await getReviews({ ticketId, page: 1, size: 100 })
            if (res && res.data) {
                setReviews(res.data.list || [])
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setLoading(false)
        }
    }

    const fetchReviewers = async () => {
        try {
            const res = await getReviewers({ page: 1, size: 100 })
            if (res && res.data) {
                setReviewers(res.data.list || [])
            }
        } catch (error) {
            HandleApiError(error)
        }
    }

    const fetchUsers = async () => {
        try {
            const res = await getUserList({ page: 1, size: 100 })
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
            setSubmitModalVisible(false)
            form.resetFields()
            setCurrentReview(null)
            fetchReviews()
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleAssignReviewers = async (values) => {
        try {
            await assignReviewers({
                ticketId,
                reviewerIds: values.reviewerIds,
            })
            setAssignModalVisible(false)
            assignForm.resetFields()
            fetchReviews()
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

    const getReviewerName = (reviewerId) => {
        const reviewer = reviewers.find(r => r.reviewerId === reviewerId)
        if (reviewer) return reviewer.userName

        const user = users.find(u => u.userid === reviewerId)
        return user ? (user.username || reviewerId) : reviewerId
    }

    const getStatusTag = (status) => {
        const statusMap = {
            pending: { color: "default", text: "待评审" },
            completed: { color: "success", text: "已完成" },
        }
        const config = statusMap[status] || { color: "default", text: status }
        return <Tag color={config.color}>{config.text}</Tag>
    }

    const columns = [
        {
            title: "评委",
            dataIndex: "reviewerId",
            key: "reviewerId",
            render: (reviewerId) => getReviewerName(reviewerId),
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

    return (
        <div style={{ padding: "24px" }}>
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
                title="工单评审"
                extra={
                    <Button
                        type="primary"
                        icon={<PlusOutlined />}
                        onClick={() => setAssignModalVisible(true)}
                    >
                        分配评委
                    </Button>
                }
            >
                <Table
                    columns={columns}
                    dataSource={reviews}
                    loading={loading}
                    rowKey="reviewId"
                    pagination={false}
                />
            </Card>

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
                            options={reviewers.map(r => ({
                                label: `${r.userName} (${r.department})`,
                                value: r.reviewerId,
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