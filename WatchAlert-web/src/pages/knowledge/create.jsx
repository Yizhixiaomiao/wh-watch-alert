"use client"

import { useState, useEffect } from "react"
import {
    Card,
    Form,
    Input,
    Select,
    Button,
    Space,
    message,
    Row,
    Col,
    Tag,
} from "antd"
import {
    SaveOutlined,
    ArrowLeftOutlined,
} from "@ant-design/icons"
import {
    createKnowledge,
    getKnowledgeCategories,
} from "../../api/knowledge"
import { HandleApiError } from "../../utils/lib"
import { clearCacheByUrl } from "../../utils/http"
import { useNavigate } from "react-router-dom"

const { TextArea } = Input

export const KnowledgeCreate = () => {
    const navigate = useNavigate()
    const [form] = Form.useForm()
    const [loading, setLoading] = useState(false)
    const [categories, setCategories] = useState([])
    const [tagInput, setTagInput] = useState("")
    const [tags, setTags] = useState([])

    useEffect(() => {
        fetchCategories()
    }, [])

    const fetchCategories = async () => {
        try {
            const res = await getKnowledgeCategories({
                isActive: true,
                page: 1,
                size: 100,
            })
            if (res && res.data) {
                setCategories(res.data.list || [])
            }
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleSubmit = async (values) => {
        setLoading(true)
        try {
            await createKnowledge({
                ...values,
                tags,
                status: "published",
            })
            message.success("知识创建成功")
            // 清除缓存
            clearCacheByUrl('/api/w8t/knowledge')
            clearCacheByUrl('/api/w8t/knowledge/list')
            navigate("/knowledge")
        } catch (error) {
            HandleApiError(error)
        } finally {
            setLoading(false)
        }
    }

    const handleAddTag = () => {
        if (tagInput && !tags.includes(tagInput)) {
            setTags([...tags, tagInput])
            setTagInput("")
        }
    }

    const handleRemoveTag = (tagToRemove) => {
        setTags(tags.filter((tag) => tag !== tagToRemove))
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

            <Card title="创建知识">
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={handleSubmit}
                    initialValues={{
                        status: "published",
                    }}
                >
                    <Row gutter={16}>
                        <Col span={24}>
                            <Form.Item
                                name="title"
                                label="标题"
                                rules={[{ required: true, message: "请输入标题" }]}
                            >
                                <Input placeholder="请输入知识标题" />
                            </Form.Item>
                        </Col>
                    </Row>

                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item
                                name="category"
                                label="分类"
                                rules={[{ required: true, message: "请选择分类" }]}
                            >
                                <Select placeholder="请选择分类">
                                    {categories.map((cat) => (
                                        <Select.Option key={cat.categoryId} value={cat.name}>
                                            {cat.name}
                                        </Select.Option>
                                    ))}
                                </Select>
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item label="标签">
                                <Space.Compact style={{ width: "100%" }}>
                                    <Input
                                        value={tagInput}
                                        onChange={(e) => setTagInput(e.target.value)}
                                        onPressEnter={handleAddTag}
                                        placeholder="输入标签后按回车"
                                    />
                                    <Button type="primary" onClick={handleAddTag}>
                                        添加
                                    </Button>
                                </Space.Compact>
                                <div style={{ marginTop: 8 }}>
                                    {tags.map((tag) => (
                                        <Tag
                                            key={tag}
                                            closable
                                            onClose={() => handleRemoveTag(tag)}
                                        >
                                            {tag}
                                        </Tag>
                                    ))}
                                </div>
                            </Form.Item>
                        </Col>
                    </Row>

                    <Form.Item
                        name="content"
                        label="内容"
                        rules={[{ required: true, message: "请输入内容" }]}
                    >
                        <TextArea
                            rows={12}
                            placeholder="请输入知识内容（支持HTML）"
                        />
                    </Form.Item>

                    <Form.Item name="sourceTicket" label="来源工单">
                        <Input placeholder="请输入来源工单ID（可选）" />
                    </Form.Item>

                    <Form.Item>
                        <Space>
                            <Button
                                type="primary"
                                htmlType="submit"
                                loading={loading}
                                icon={<SaveOutlined />}
                            >
                                保存
                            </Button>
                            <Button onClick={() => navigate("/knowledge")}>
                                取消
                            </Button>
                        </Space>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    )
}

export default KnowledgeCreate