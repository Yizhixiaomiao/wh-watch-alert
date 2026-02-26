import { useState, useEffect } from "react"
import { 
    Input, 
    Button, 
    Card, 
    List, 
    Tag, 
    Space, 
    Typography, 
    Empty,
    Modal,
    message
} from "antd"
import { 
    SearchOutlined, 
    PhoneOutlined, 
    ClockCircleOutlined,
    CheckCircleOutlined,
    ExclamationCircleOutlined,
    SyncOutlined
} from "@ant-design/icons"
import { mobileQueryTicket } from "../../../api/ticket"

const { Text, Title } = Typography
const { Search } = Input

export const MobileTicketQuery = () => {
    const [loading, setLoading] = useState(false)
    const [tickets, setTickets] = useState([])
    const [searchPhone, setSearchPhone] = useState('')
    const [searchTicketNo, setSearchTicketNo] = useState('')
    const [searchMode, setSearchMode] = useState('phone')
    const [selectedTicket, setSelectedTicket] = useState(null)

    useEffect(() => {
        // 从URL参数获取默认搜索条件
        const urlParams = new URLSearchParams(window.location.search)
        const phone = urlParams.get('phone')
        const ticketNo = urlParams.get('ticketNo')
        
        if (phone) {
            setSearchPhone(phone)
            setSearchMode('phone')
            handleSearch('phone', phone)
        } else if (ticketNo) {
            setSearchTicketNo(ticketNo)
            setSearchMode('ticketNo')
            handleSearch('ticketNo', ticketNo)
        }
    }, [])

    const handleSearch = async (mode, value) => {
        if (!value) {
            message.warning('请输入搜索条件')
            return
        }

        try {
            setLoading(true)
            
            const params = mode === 'phone' 
                ? { contactPhone: value }
                : { ticketNo: value }

            const res = await mobileQueryTicket(params)
            
            if (res && res.data) {
                setTickets(Array.isArray(res.data) ? res.data : [res.data])
                if (!Array.isArray(res.data) && res.data.ticketId) {
                    setTickets([res.data])
                }
            } else {
                setTickets([])
            }
        } catch (error) {
            message.error('查询失败，请重试')
            setTickets([])
        } finally {
            setLoading(false)
        }
    }

    const getStatusColor = (status) => {
        const colorMap = {
            'Pending': 'orange',
            'Assigned': 'blue',
            'Processing': 'cyan',
            'Resolved': 'purple',
            'Closed': 'green',
            'Cancelled': 'red',
            'Escalated': 'magenta'
        }
        return colorMap[status] || 'default'
    }

    const getStatusText = (status) => {
        const textMap = {
            'Pending': '待处理',
            'Assigned': '已分配',
            'Processing': '处理中',
            'Resolved': '待验证',
            'Closed': '已完成',
            'Cancelled': '已取消',
            'Escalated': '已升级'
        }
        return textMap[status] || status
    }

    const getPriorityColor = (priority) => {
        const colorMap = {
            'P0': 'red',
            'P1': 'orange',
            'P2': 'gold',
            'P3': 'green',
            'P4': 'blue'
        }
        return colorMap[priority] || 'default'
    }

    const getPriorityText = (priority) => {
        const textMap = {
            'P0': 'P0 - 最高',
            'P1': 'P1 - 高',
            'P2': 'P2 - 中',
            'P3': 'P3 - 低',
            'P4': 'P4 - 最低'
        }
        return textMap[priority] || priority
    }

    const showDetail = (ticket) => {
        // 跳转到PC端详情页面
        window.location.href = `/ticket/detail/${ticket.ticketId}`
    }

    const renderTicketItem = (ticket) => (
        <List.Item
            key={ticket.ticketId}
            onClick={() => showDetail(ticket)}
            style={{ 
                background: 'white', 
                marginBottom: 8, 
                borderRadius: 8,
                cursor: 'pointer',
                padding: 16
            }}
        >
            <div style={{ width: '100%' }}>
                {/* 工单号和状态 */}
                <div style={{ 
                    display: 'flex', 
                    justifyContent: 'space-between', 
                    alignItems: 'center',
                    marginBottom: 8
                }}>
                    <Text strong style={{ fontSize: 14 }}>
                        {ticket.ticketNo}
                    </Text>
                    <Tag color={getStatusColor(ticket.status)}>
                        {getStatusText(ticket.status)}
                    </Tag>
                </div>

                {/* 标题 */}
                <div style={{ marginBottom: 8 }}>
                    <Text ellipsis style={{ fontSize: 13, lineHeight: 1.4 }}>
                        {ticket.title}
                    </Text>
                </div>

                {/* 优先级和时间 */}
                <div style={{ 
                    display: 'flex', 
                    justifyContent: 'space-between', 
                    alignItems: 'center'
                }}>
                    <Tag color={getPriorityColor(ticket.priority)} size="small">
                        {getPriorityText(ticket.priority)}
                    </Tag>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        <ClockCircleOutlined style={{ marginRight: 4 }} />
                        {new Date(ticket.createdAt * 1000).toLocaleString()}
                    </Text>
                </div>

                {/* 处理人和预计时间 */}
                {(ticket.assignedTo || ticket.estimateTime) && (
                    <div style={{ marginTop: 8, fontSize: 12 }}>
                        {ticket.assignedTo && (
                            <div style={{ marginBottom: 4 }}>
                                <Text type="secondary">处理人：{ticket.assignedTo}</Text>
                            </div>
                        )}
                        {ticket.estimateTime && (
                            <div>
                                <Text type="secondary">
                                    <SyncOutlined style={{ marginRight: 4 }} />
                                    预计处理时间：{ticket.estimateTime}
                                </Text>
                            </div>
                        )}
                    </div>
                )}
            </div>
        </List.Item>
    )

    return (
        <div style={{ 
            maxWidth: '100%', 
            margin: '0 auto', 
            padding: '16px',
            background: '#f5f5f5',
            minHeight: '100vh'
        }}>
            {/* 头部 */}
            <div style={{ 
                textAlign: 'center', 
                padding: '20px 0',
                background: 'white',
                marginBottom: '16px',
                borderRadius: '8px'
            }}>
                <Title level={3} style={{ margin: 0, color: '#1890ff' }}>
                    工单进度查询
                </Title>
            </div>

            {/* 搜索区域 */}
            <Card style={{ marginBottom: 16 }}>
                <div style={{ marginBottom: 16 }}>
                    <Space>
                        <Button 
                            type={searchMode === 'phone' ? 'primary' : 'default'}
                            icon={<PhoneOutlined />}
                            onClick={() => setSearchMode('phone')}
                            size="small"
                        >
                            手机号查询
                        </Button>
                        <Button 
                            type={searchMode === 'ticketNo' ? 'primary' : 'default'}
                            icon={<SearchOutlined />}
                            onClick={() => setSearchMode('ticketNo')}
                            size="small"
                        >
                            工单号查询
                        </Button>
                    </Space>
                </div>

                {searchMode === 'phone' ? (
                    <Search
                        placeholder="请输入报修时使用的手机号"
                        value={searchPhone}
                        onChange={(e) => setSearchPhone(e.target.value)}
                        onSearch={() => handleSearch('phone', searchPhone)}
                        enterButton={<SearchOutlined />}
                        loading={loading}
                    />
                ) : (
                    <Search
                        placeholder="请输入工单号（如：TK20240101123456）"
                        value={searchTicketNo}
                        onChange={(e) => setSearchTicketNo(e.target.value)}
                        onSearch={() => handleSearch('ticketNo', searchTicketNo)}
                        enterButton={<SearchOutlined />}
                        loading={loading}
                    />
                )}
            </Card>

            {/* 查询结果 */}
            {tickets.length > 0 && (
                <Card title={
                    <Space>
                        <Text>查询结果</Text>
                        <Tag>{tickets.length} 条记录</Tag>
                    </Space>
                }>
                    <List
                        dataSource={tickets}
                        renderItem={renderTicketItem}
                        style={{ padding: 0 }}
                    />
                </Card>
            )}

            {/* 无结果提示 */}
            {!loading && tickets.length === 0 && (searchPhone || searchTicketNo) && (
                <Card>
                    <Empty 
                        description="未找到相关工单"
                        image={Empty.PRESENTED_IMAGE_SIMPLE}
                    >
                        <Space direction="vertical">
                            <Text type="secondary">
                                请检查手机号或工单号是否正确
                            </Text>
                            <Button 
                                type="primary" 
                                onClick={() => window.location.href = '/mobile/ticket/create'}
                            >
                                新建报修单
                            </Button>
                        </Space>
                    </Empty>
                </Card>
            )}

            {/* 提示用户查看详情 */}
                <div style={{ textAlign: 'center', padding: '40px' }}>
                    <Space direction="vertical">
                        <Text type="secondary">
                            点击工单可查看完整详情（含图片）
                        </Text>
                        <Button 
                            type="primary"
                            onClick={() => window.location.href = `/ticket/detail/${selectedTicket.ticketId}`}
                        >
                            查看详情
                        </Button>
                    </Space>
                </div>
        </div>
    )
}