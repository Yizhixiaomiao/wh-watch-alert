"use client"

import { useState, useEffect } from "react"
import {
    Card,
    Button,
    Space,
    Tag,
    Row,
    Col,
    Divider,
    message,
    Typography,
    Descriptions,
    Statistic,
} from "antd"
import {
    LikeOutlined,
    EyeOutlined,
    FileTextOutlined,
    EditOutlined,
    ArrowLeftOutlined,
} from "@ant-design/icons"
import {
    getKnowledge,
    likeKnowledge,
    saveKnowledgeToTicket,
} from "../../api/knowledge"
import { HandleApiError, FormatTime } from "../../utils/lib"
import { useNavigate, useParams } from "react-router-dom"

const { Title, Text, Paragraph } = Typography

export const KnowledgeDetail = () => {
    const navigate = useNavigate()
    const { id } = useParams()
    const [loading, setLoading] = useState(false)
    const [knowledge, setKnowledge] = useState(null)
    const [liked, setLiked] = useState(false)

    useEffect(() => {
        fetchKnowledge()
    }, [id])

    const fetchKnowledge = async () => {
        setLoading(true)
        try {
            const res = await getKnowledge({ knowledgeId: id })
            if (res && res.data) {
                setKnowledge(res.data)
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setLoading(false)
        }
    }

    const handleLike = async () => {
        try {
            const res = await likeKnowledge({
                knowledgeId: id,
            })
            if (res && res.data) {
                setLiked(res.data.liked)
                fetchKnowledge()
            }
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleSaveToTicket = async () => {
        const urlParams = new URLSearchParams(window.location.search)
        const ticketId = urlParams.get('ticketId')
        
        if (!ticketId) {
            message.error('请从工单页面打开此知识')
            return
        }

        try {
            await saveKnowledgeToTicket({
                knowledgeId: id,
                ticketId,
            })
            message.success('知识已添加到工单')
        } catch (error) {
            HandleApiError(error)
        }
    }

    if (loading || !knowledge) {
        return <div>加载中...</div>
    }

    return (
        <div style={{ padding: "24px" }}>
            <Button
                icon={<ArrowLeftOutlined />}
                onClick={() => navigate("/knowledge")}
                style={{ marginBottom: 16 }}
            >
                返回列表
            </Button>

            <Card loading={loading}>
                <div style={{ marginBottom: 24 }}>
                    <Title level={2}>{knowledge.title}</Title>
                    <Space>
                        <Tag color="blue">{knowledge.category}</Tag>
                        {knowledge.tags && knowledge.tags.map((tag, index) => (
                            <Tag key={index} color="cyan">
                                {tag}
                            </Tag>
                        ))}
                        <Tag
                            color={
                                knowledge.status === "published"
                                    ? "success"
                                    : knowledge.status === "draft"
                                    ? "default"
                                    : "default"
                            }
                        >
                            {knowledge.status === "published"
                                ? "已发布"
                                : knowledge.status === "draft"
                                ? "草稿"
                                : "已归档"}
                        </Tag>
                    </Space>
                </div>

                <Descriptions column={2} bordered style={{ marginBottom: 24 }}>
                    <Descriptions.Item label="创建时间">
                        {FormatTime(knowledge.createdAt)}
                    </Descriptions.Item>
                    <Descriptions.Item label="更新时间">
                        {FormatTime(knowledge.updatedAt)}
                    </Descriptions.Item>
                    <Descriptions.Item label="来源工单">
                        {knowledge.sourceTicket ? (
                            <a href={`/ticket/detail/${knowledge.sourceTicket}`}>
                                {knowledge.sourceTicket}
                            </a>
                        ) : (
                            "-"
                        )}
                    </Descriptions.Item>
                    <Descriptions.Item label="作者">
                        {knowledge.authorId || "-"}
                    </Descriptions.Item>
                </Descriptions>

                <Row gutter={24} style={{ marginBottom: 24 }}>
                    <Col span={8}>
                        <Card size="small">
                            <Statistic
                                title="浏览次数"
                                value={knowledge.viewCount}
                                prefix={<EyeOutlined />}
                            />
                        </Card>
                    </Col>
                    <Col span={8}>
                        <Card size="small">
                            <Statistic
                                title="点赞次数"
                                value={knowledge.likeCount}
                                prefix={<LikeOutlined />}
                            />
                        </Card>
                    </Col>
                    <Col span={8}>
                        <Card size="small">
                            <Statistic
                                title="使用次数"
                                value={knowledge.useCount}
                                prefix={<FileTextOutlined />}
                            />
                        </Card>
                    </Col>
                </Row>

                <Divider />

                <div
                    style={{
                        minHeight: 300,
                        padding: 16,
                        background: "#fafafa",
                        borderRadius: 4,
                    }}
                >
                    <div dangerouslySetInnerHTML={{ __html: knowledge.content }} />
                </div>

                <Divider />

                <Space>
                    <Button
                        type={liked ? "default" : "primary"}
                        icon={<LikeOutlined />}
                        onClick={handleLike}
                    >
                        {liked ? "取消点赞" : "点赞"}
                    </Button>
                    <Button icon={<EditOutlined />}>
                        编辑
                    </Button>
                    {new URLSearchParams(window.location.search).get('ticketId') && (
                        <Button type="primary" onClick={handleSaveToTicket}>
                            添加到工单
                        </Button>
                    )}
                </Space>
            </Card>
        </div>
    )
}

export default KnowledgeDetail