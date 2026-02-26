"use client"

import { useState, useEffect } from "react"
import {
    Card,
    Table,
    Button,
    Space,
    Tag,
    Input,
    Select,
    Row,
    Col,
    Statistic,
    message,
    Dropdown,
    Popconfirm,
    Modal,
    Form,
    List,
    Switch,
    Typography,
} from "antd"
import {
    PlusOutlined,
    SearchOutlined,
    EyeOutlined,
    LikeOutlined,
    FileTextOutlined,
    DeleteOutlined,
    DownOutlined,
    DownloadOutlined,
    CheckCircleOutlined,
    InboxOutlined,
} from "@ant-design/icons"
import {
    getKnowledges,
    getKnowledgeCategories,
    likeKnowledge,
    deleteKnowledge,
    updateKnowledge,
    createKnowledgeCategory,
    updateKnowledgeCategory,
    deleteKnowledgeCategory,
} from "../../api/knowledge"
import { HandleApiError, FormatTime } from "../../utils/lib"
import { useNavigate } from "react-router-dom"
import { clearCacheByUrl } from "../../utils/http"
import * as XLSX from 'xlsx'

const { Text } = Typography
const { Search } = Input

export const KnowledgeList = () => {
    const navigate = useNavigate()
    const [loading, setLoading] = useState(false)
    const [knowledges, setKnowledges] = useState([])
    const [categories, setCategories] = useState([])
    const [filters, setFilters] = useState({})
    const [likedSet, setLikedSet] = useState(new Set())
    const [selectedRowKeys, setSelectedRowKeys] = useState([])
    const [batchActionLoading, setBatchActionLoading] = useState(false)
    const [categoryModalVisible, setCategoryModalVisible] = useState(false)
    const [editCategoryModalVisible, setEditCategoryModalVisible] = useState(false)
    const [editingCategory, setEditingCategory] = useState(null)
    const [categoryForm] = Form.useForm()

    useEffect(() => {
        fetchKnowledges()
        fetchCategories()
    }, [])

    const fetchKnowledges = async (params = {}) => {
        setLoading(true)
        try {
            const res = await getKnowledges({
                page: 1,
                size: 100,
                ...params,
                ...filters,
            }, { skipCache: true })
            if (res && res.data) {
                setKnowledges(res.data.list || [])
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setLoading(false)
        }
    }

    const fetchCategories = async () => {
        try {
            const res = await getKnowledgeCategories({
                page: 1,
                size: 100,
            }, { skipCache: true })
            if (res && res.data) {
                setCategories(res.data.list || [])
            }
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleOpenCategoryModal = (category = null) => {
        setEditingCategory(category)
        if (category) {
            categoryForm.setFieldsValue({
                name: category.name,
                description: category.description,
                isActive: category.isActive !== false,
            })
        } else {
            categoryForm.resetFields()
        }
        setEditCategoryModalVisible(true)
    }

    const handleCreateOrUpdateCategory = async (values) => {
        try {
            if (editingCategory) {
                await updateKnowledgeCategory({
                    categoryId: editingCategory.categoryId,
                    name: values.name,
                    description: values.description,
                    isActive: values.isActive,
                })
            } else {
                await createKnowledgeCategory({
                    name: values.name,
                    description: values.description,
                    isActive: values.isActive !== false,
                })
            }
            setEditCategoryModalVisible(false)
            categoryForm.resetFields()
            setEditingCategory(null)
            clearCacheByUrl('/api/w8t/knowledge/category/list')
            await fetchCategories()
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleDeleteCategory = (categoryId) => {
        Modal.confirm({
            title: '确认删除',
            content: '确定要删除这个分类吗？删除后该分类下的知识将变为未分类状态。',
            onOk: async () => {
                try {
                    await deleteKnowledgeCategory({ categoryId })
                    clearCacheByUrl('/api/w8t/knowledge/category/list')
                    await fetchCategories()
                    await fetchKnowledges()
                    message.success('分类删除成功')
                } catch (error) {
                    HandleApiError(error)
                }
            },
        })
    }

    const handleInitDefaultCategories = async () => {
        Modal.confirm({
            title: '初始化默认分类',
            content: '确定要初始化默认分类吗？这将创建以下分类：故障处理、网络问题、系统优化、安全防护、其他问题。',
            onOk: async () => {
                try {
                    const defaultCategories = [
                        { name: '故障处理', description: '系统故障和异常处理相关知识' },
                        { name: '网络问题', description: '网络配置和故障排查相关知识' },
                        { name: '系统优化', description: '系统性能优化和调优相关知识' },
                        { name: '安全防护', description: '系统安全和防护相关知识' },
                        { name: '其他问题', description: '其他类型的问题和解决方案' },
                    ]

                    for (const category of defaultCategories) {
                        await createKnowledgeCategory({
                            name: category.name,
                            description: category.description,
                        })
                    }

                    clearCacheByUrl('/api/w8t/knowledge/category/list')
                    message.success('默认分类初始化成功')
                    await fetchCategories()
                } catch (error) {
                    HandleApiError(error)
                }
            },
        })
    }

    const handleSearch = (keyword) => {
        setFilters({ ...filters, keyword })
        fetchKnowledges({ keyword })
    }

    const handleLike = async (knowledge) => {
        try {
            const res = await likeKnowledge({
                knowledgeId: knowledge.knowledgeId,
            })
            if (res && res.data) {
                const liked = res.data.liked
                const newLikedSet = new Set(likedSet)
                if (liked) {
                    newLikedSet.add(knowledge.knowledgeId)
                } else {
                    newLikedSet.delete(knowledge.knowledgeId)
                }
                setLikedSet(newLikedSet)
                clearCacheByUrl('/api/w8t/knowledge/list')
                clearCacheByUrl('/api/w8t/knowledge')
                await fetchKnowledges()
            }
        } catch (error) {
            HandleApiError(error)
        }
    }

    const handleBatchDelete = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请先选择要删除的知识')
            return
        }

        setBatchActionLoading(true)
        try {
            let successCount = 0
            let failCount = 0

            for (const knowledgeId of selectedRowKeys) {
                try {
                    await deleteKnowledge({ knowledgeId })
                    successCount++
                } catch (error) {
                    failCount++
                    console.error(`删除知识 ${knowledgeId} 失败:`, error)
                }
            }

            if (successCount > 0) {
                message.success(`成功删除 ${successCount} 条知识${failCount > 0 ? `，失败 ${failCount} 条` : ''}`)
                clearCacheByUrl('/api/w8t/knowledge/list')
                clearCacheByUrl('/api/w8t/knowledge')
                await fetchKnowledges()
                setSelectedRowKeys([])
            } else {
                message.error(`删除失败，共 ${failCount} 条知识`)
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setBatchActionLoading(false)
        }
    }

    const handleBatchArchive = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请先选择要归档的知识')
            return
        }

        setBatchActionLoading(true)
        try {
            let successCount = 0
            let failCount = 0

            for (const knowledgeId of selectedRowKeys) {
                try {
                    await updateKnowledge({ knowledgeId, status: 'archived' })
                    successCount++
                } catch (error) {
                    failCount++
                    console.error(`归档知识 ${knowledgeId} 失败:`, error)
                }
            }

            if (successCount > 0) {
                message.success(`成功归档 ${successCount} 条知识${failCount > 0 ? `，失败 ${failCount} 条` : ''}`)
                clearCacheByUrl('/api/w8t/knowledge/list')
                clearCacheByUrl('/api/w8t/knowledge')
                await fetchKnowledges()
                setSelectedRowKeys([])
            } else {
                message.error(`归档失败，共 ${failCount} 条知识`)
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setBatchActionLoading(false)
        }
    }

    const handleBatchPublish = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请先选择要发布的知识')
            return
        }

        setBatchActionLoading(true)
        try {
            let successCount = 0
            let failCount = 0

            for (const knowledgeId of selectedRowKeys) {
                try {
                    await updateKnowledge({ knowledgeId, status: 'published' })
                    successCount++
                } catch (error) {
                    failCount++
                    console.error(`发布知识 ${knowledgeId} 失败:`, error)
                }
            }

            if (successCount > 0) {
                message.success(`成功发布 ${successCount} 条知识${failCount > 0 ? `，失败 ${failCount} 条` : ''}`)
                clearCacheByUrl('/api/w8t/knowledge/list')
                clearCacheByUrl('/api/w8t/knowledge')
                await fetchKnowledges()
                setSelectedRowKeys([])
            } else {
                message.error(`发布失败，共 ${failCount} 条知识`)
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setBatchActionLoading(false)
        }
    }

    const handleExport = async () => {
        message.loading({ content: "正在导出...", key: "export" })

        try {
            const params = {
                page: 1,
                size: 10000,
                ...filters,
            }
            const res = await getKnowledges(params, { skipCache: true })

            if (!res || !res.data || !res.data.list || res.data.list.length === 0) {
                message.warning({ content: "没有数据可导出", key: "export" })
                return
            }

            const exportData = res.data.list.map(k => {
                return {
                    "标题": k.title || "",
                    "分类": k.category || "",
                    "标签": k.tags && k.tags.length > 0 ? k.tags.join(', ') : "",
                    "状态": k.status === 'draft' ? "草稿" : (k.status === 'published' ? "已发布" : (k.status === 'archived' ? "已归档" : k.status)),
                    "浏览量": k.viewCount || 0,
                    "点赞数": k.likeCount || 0,
                    "使用数": k.useCount || 0,
                    "创建时间": FormatTime(k.createdAt),
                    "更新时间": FormatTime(k.updatedAt),
                }
            })

            const ws = XLSX.utils.json_to_sheet(exportData)
            const wb = XLSX.utils.book_new()
            XLSX.utils.book_append_sheet(wb, ws, "知识列表")

            const colWidths = [
                { wch: 40 },
                { wch: 15 },
                { wch: 30 },
                { wch: 12 },
                { wch: 12 },
                { wch: 12 },
                { wch: 12 },
                { wch: 20 },
                { wch: 20 },
            ]
            ws['!cols'] = colWidths

            const fileName = `知识列表_${new Date().toLocaleDateString().replace(/\//g, '-')}.xlsx`
            XLSX.writeFile(wb, fileName)

            message.success({ content: `成功导出 ${exportData.length} 条数据`, key: "export" })
        } catch (error) {
            message.error({ content: "导出失败", key: "export" })
            HandleApiError(error)
        }
    }

    const batchOperationMenu = {
        items: [
            {
                key: "batchPublish",
                label: "批量发布",
                icon: <CheckCircleOutlined />,
                onClick: handleBatchPublish,
                disabled: selectedRowKeys.length === 0,
            },
            {
                key: "batchArchive",
                label: "批量归档",
                icon: <InboxOutlined />,
                onClick: handleBatchArchive,
                disabled: selectedRowKeys.length === 0,
            },
            {
                key: "batchDelete",
                label: "批量删除",
                icon: <DeleteOutlined />,
                onClick: handleBatchDelete,
                disabled: selectedRowKeys.length === 0,
                danger: true,
            },
            {
                type: 'divider',
            },
            {
                key: "batchExport",
                label: "导出Excel",
                icon: <DownloadOutlined />,
                onClick: handleExport,
            },
        ],
    }

    const columns = [
        {
            title: "标题",
            dataIndex: "title",
            key: "title",
            render: (title, record) => (
                <a
                    onClick={(e) => {
                        e.stopPropagation()
                        navigate(`/knowledge/detail/${record.knowledgeId}`)
                    }}
                    style={{ fontWeight: 500, cursor: 'pointer' }}
                >
                    {title}
                </a>
            ),
        },
        {
            title: "分类",
            dataIndex: "category",
            key: "category",
            render: (category) => <Tag color="blue">{category}</Tag>,
        },
        {
            title: "标签",
            dataIndex: "tags",
            key: "tags",
            render: (tags) =>
                tags && tags.length > 0 ? (
                    <>
                        {tags.slice(0, 3).map((tag, index) => (
                            <Tag key={index} color="cyan">
                                {tag}
                            </Tag>
                        ))}
                        {tags.length > 3 && <Tag>+{tags.length - 3}</Tag>}
                    </>
                ) : (
                    "-"
                ),
        },
        {
            title: "状态",
            dataIndex: "status",
            key: "status",
            render: (status) => {
                const statusMap = {
                    draft: { color: "default", text: "草稿" },
                    published: { color: "success", text: "已发布" },
                    archived: { color: "default", text: "已归档" },
                }
                const config = statusMap[status] || { color: "default", text: status }
                return <Tag color={config.color}>{config.text}</Tag>
            },
        },
        {
            title: "统计",
            key: "stats",
            render: (_, record) => (
                <Space>
                    <span>
                        <EyeOutlined /> {record.viewCount}
                    </span>
                    <span>
                        <LikeOutlined /> {record.likeCount}
                    </span>
                    <span>
                        <FileTextOutlined /> {record.useCount}
                    </span>
                </Space>
            ),
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
                    <Button
                        type="link"
                        size="small"
                        onClick={() => navigate(`/knowledge/detail/${record.knowledgeId}`)}
                    >
                        查看
                    </Button>
                    <Button
                        type="link"
                        size="small"
                        icon={<LikeOutlined />}
                        onClick={() => handleLike(record)}
                    >
                        {likedSet.has(record.knowledgeId) ? "取消点赞" : "点赞"}
                    </Button>
                </Space>
            ),
        },
    ]

    const totalLikes = knowledges.reduce((sum, k) => sum + k.likeCount, 0)
    const totalViews = knowledges.reduce((sum, k) => sum + k.viewCount, 0)
    const totalUses = knowledges.reduce((sum, k) => sum + k.useCount, 0)

    return (
        <div style={{ padding: "24px" }}>
            <Row gutter={24} style={{ marginBottom: 24 }}>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="知识总数"
                            value={knowledges.length}
                            suffix="篇"
                            prefix={<FileTextOutlined />}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="总浏览量"
                            value={totalViews}
                            suffix="次"
                            prefix={<EyeOutlined />}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="总点赞数"
                            value={totalLikes}
                            suffix="次"
                            prefix={<LikeOutlined />}
                        />
                    </Card>
                </Col>
                <Col span={6}>
                    <Card>
                        <Statistic
                            title="总使用数"
                            value={totalUses}
                            suffix="次"
                            prefix={<FileTextOutlined />}
                        />
                    </Card>
                </Col>
            </Row>

            <Card
                title="知识库"
                extra={
                    <Space>
                        <Select
                            placeholder="筛选分类"
                            style={{ width: 150 }}
                            allowClear
                            onChange={(value) => {
                                setFilters({ ...filters, category: value })
                                fetchKnowledges({ category: value })
                            }}
                        >
                            {categories.map((cat) => (
                                <Select.Option key={cat.categoryId} value={cat.name}>
                                    {cat.name}
                                </Select.Option>
                            ))}
                        </Select>
                        <Search
                            placeholder="搜索知识"
                            allowClear
                            onSearch={handleSearch}
                            style={{ width: 250 }}
                        />
                        <Dropdown menu={batchOperationMenu}>
                            <Button>
                                批量操作 {selectedRowKeys.length > 0 && `(${selectedRowKeys.length})`} <DownOutlined />
                            </Button>
                        </Dropdown>
                        <Button onClick={() => setCategoryModalVisible(true)}>
                            管理分类
                        </Button>
                        <Button
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={() => navigate("/knowledge/create")}
                        >
                            添加知识
                        </Button>
                    </Space>
                }
            >
                <Table
                    columns={columns}
                    dataSource={knowledges}
                    loading={loading}
                    rowKey="knowledgeId"
                    rowSelection={{
                        selectedRowKeys,
                        onChange: setSelectedRowKeys,
                    }}
                    pagination={{
                        showSizeChanger: true,
                        showTotal: (total) => `共 ${total} 条`,
                    }}
                />
            </Card>

            {/* 分类管理 Modal */}
            <Modal
                title="管理分类"
                open={categoryModalVisible}
                onCancel={() => {
                    setCategoryModalVisible(false)
                    categoryForm.resetFields()
                    setEditingCategory(null)
                }}
                footer={null}
                width={800}
            >
                <div style={{ marginBottom: 16 }}>
                    <Button type="primary" icon={<PlusOutlined />} onClick={() => handleOpenCategoryModal()}>
                        添加分类
                    </Button>
                    {categories.length === 0 && (
                        <Button
                            style={{ marginLeft: 8 }}
                            onClick={handleInitDefaultCategories}
                        >
                            初始化默认分类
                        </Button>
                    )}
                </div>
                {categories.length > 0 ? (
                    <List
                        dataSource={categories}
                        renderItem={(category) => (
                            <List.Item
                                actions={[
                                    <Button
                                        type="link"
                                        size="small"
                                        onClick={() => handleOpenCategoryModal(category)}
                                    >
                                        编辑
                                    </Button>,
                                    <Button
                                        type="link"
                                        size="small"
                                        danger
                                        onClick={() => handleDeleteCategory(category.categoryId)}
                                    >
                                        删除
                                    </Button>,
                                ]}
                            >
                                <List.Item.Meta
                                    title={
                                        <Space>
                                            <Text strong>{category.name}</Text>
                                            {category.isActive !== false && (
                                                <Tag color="green">启用</Tag>
                                            )}
                                            {category.isActive === false && (
                                                <Tag color="red">禁用</Tag>
                                            )}
                                        </Space>
                                    }
                                    description={category.description || '暂无描述'}
                                />
                            </List.Item>
                        )}
                    />
                ) : (
                    <div style={{ textAlign: 'center', padding: '40px 0', color: '#999' }}>
                        <InboxOutlined style={{ fontSize: '48px', marginBottom: '16px' }} />
                        <p>暂无分类</p>
                        <p style={{ fontSize: '14px' }}>点击"添加分类"创建新分类，或点击"初始化默认分类"快速创建常用分类</p>
                    </div>
                )}
            </Modal>

            {/* 添加/编辑分类 Modal */}
            <Modal
                title={editingCategory ? "编辑分类" : "添加分类"}
                open={editCategoryModalVisible}
                onCancel={() => {
                    setEditCategoryModalVisible(false)
                    categoryForm.resetFields()
                    setEditingCategory(null)
                }}
                onOk={() => categoryForm.submit()}
            >
                <Form form={categoryForm} layout="vertical" onFinish={handleCreateOrUpdateCategory}>
                    <Form.Item
                        name="name"
                        label="分类名称"
                        rules={[{ required: true, message: "请输入分类名称" }]}
                    >
                        <Input placeholder="请输入分类名称" />
                    </Form.Item>
                    <Form.Item
                        name="description"
                        label="分类描述"
                    >
                        <Input.TextArea rows={3} placeholder="请输入分类描述" />
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

export default KnowledgeList