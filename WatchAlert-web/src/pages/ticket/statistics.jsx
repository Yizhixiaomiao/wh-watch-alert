"use client"

import { useState, useEffect } from "react"
import {
    Card,
    Row,
    Col,
    Statistic,
    Select,
    DatePicker,
    Button,
    Table,
    Tag,
    Space,
    message,
    Spin,
    Typography,
    Tabs,
    Divider
} from "antd"
import {
    BarChartOutlined,
    LineChartOutlined,
    PieChartOutlined,
    DownloadOutlined,
    ReloadOutlined
} from "@ant-design/icons"
import { getTicketStatistics, getTicketList } from "../../api/ticket"
import { NoticeMetricChart } from "../chart/noticeMetricChart"
import ReactECharts from "echarts-for-react"
import * as XLSX from 'xlsx'
import dayjs from 'dayjs'

const { RangePicker } = DatePicker
const { Option } = Select
const { Title, Text } = Typography
const { TabPane } = Tabs

export const TicketStatistics = () => {
    const [loading, setLoading] = useState(false)
    const [dateRange, setDateRange] = useState([
        dayjs().subtract(30, 'day'),
        dayjs()
    ])
    const [statistics, setStatistics] = useState(null)
    const [trendPeriod, setTrendPeriod] = useState('day')

    useEffect(() => {
        fetchStatistics()
    }, [dateRange])

    const fetchStatistics = async () => {
        try {
            setLoading(true)
            const params = {
                startTime: dateRange[0].unix(),
                endTime: dateRange[1].unix()
            }
            const res = await getTicketStatistics(params)
            if (res && res.data) {
                setStatistics(res.data)
            }
        } catch (error) {
            message.error("获取统计数据失败")
            console.error(error)
        } finally {
            setLoading(false)
        }
    }

    const handleExport = () => {
        if (!statistics) {
            message.warning("暂无数据可导出")
            return
        }

        const data = [
            ["统计指标", "数值"],
            ["总工单数", statistics.totalCount],
            ["待处理工单", statistics.pendingCount],
            ["处理中工单", statistics.processingCount],
            ["已关闭工单", statistics.closedCount],
            ["逾期工单", statistics.overdueCount],
            ["平均响应时间(秒)", statistics.avgResponseTime],
            ["平均解决时间(秒)", statistics.avgResolution],
            ["SLA达成率(%)", (statistics.slaRate * 100).toFixed(2)]
        ]

        const ws = XLSX.utils.aoa_to_sheet(data)
        const wb = XLSX.utils.book_new()
        XLSX.utils.book_append_sheet(wb, ws, "工单统计")
        XLSX.writeFile(wb, `工单统计_${dayjs().format('YYYY-MM-DD')}.xlsx`)
        message.success("导出成功")
    }

    const priorityOption = {
        tooltip: {
            trigger: 'item',
            formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: {
            orient: 'vertical',
            right: 10,
            top: 'center'
        },
        series: [
            {
                name: '优先级',
                type: 'pie',
                radius: ['50%', '70%'],
                data: statistics?.priorityStats ? Object.entries(statistics.priorityStats).map(([key, value]) => ({
                    name: key,
                    value: value
                })) : [],
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)'
                    }
                }
            }
        ]
    }

    const statusOption = {
        tooltip: {
            trigger: 'item',
            formatter: '{a} <br/>{b}: {c} ({d}%)'
        },
        legend: {
            orient: 'vertical',
            right: 10,
            top: 'center'
        },
        series: [
            {
                name: '状态',
                type: 'pie',
                radius: ['50%', '70%'],
                data: statistics?.statusStats ? Object.entries(statistics.statusStats).map(([key, value]) => ({
                    name: key,
                    value: value
                })) : [],
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)'
                    }
                }
            }
        ]
    }

    const trendOption = {
        tooltip: {
            trigger: 'axis'
        },
        legend: {
            data: ['工单数', '已解决']
        },
        xAxis: {
            type: 'category',
            data: statistics?.trendData?.map(item => item.date) || []
        },
        yAxis: {
            type: 'value'
        },
        series: [
            {
                name: '工单数',
                type: 'line',
                data: statistics?.trendData?.map(item => item.count) || [],
                smooth: true
            },
            {
                name: '已解决',
                type: 'line',
                data: statistics?.trendData?.map(item => item.resolved) || [],
                smooth: true
            }
        ]
    }

    const responseTimeOption = {
        tooltip: {
            trigger: 'axis'
        },
        xAxis: {
            type: 'category',
            data: statistics?.userStats?.map(item => item.userName) || []
        },
        yAxis: {
            type: 'value',
            name: '响应时间(秒)'
        },
        series: [
            {
                name: '平均响应时间',
                type: 'bar',
                data: statistics?.userStats?.map(item => item.avgResponseTime) || []
            }
        ]
    }

    const resolutionTimeOption = {
        tooltip: {
            trigger: 'axis'
        },
        xAxis: {
            type: 'category',
            data: statistics?.userStats?.map(item => item.userName) || []
        },
        yAxis: {
            type: 'value',
            name: '解决时间(秒)'
        },
        series: [
            {
                name: '平均解决时间',
                type: 'bar',
                data: statistics?.userStats?.map(item => item.avgResolution) || []
            }
        ]
    }

    const userColumns = [
        {
            title: '用户名',
            dataIndex: 'userName',
            key: 'userName'
        },
        {
            title: '处理工单数',
            dataIndex: 'ticketCount',
            key: 'ticketCount',
            render: (text) => <Tag color="blue">{text}</Tag>
        },
        {
            title: '平均响应时间',
            dataIndex: 'avgResponseTime',
            key: 'avgResponseTime',
            render: (text) => <span>{Math.round(text / 60)} 分钟</span>
        },
        {
            title: '平均解决时间',
            dataIndex: 'avgResolution',
            key: 'avgResolution',
            render: (text) => <span>{Math.round(text / 60)} 分钟</span>
        },
        {
            title: 'SLA达成率',
            dataIndex: 'slaRate',
            key: 'slaRate',
            render: (text) => (
                <Tag color={text >= 0.9 ? 'green' : text >= 0.7 ? 'orange' : 'red'}>
                    {(text * 100).toFixed(2)}%
                </Tag>
            )
        }
    ]

    const formatTime = (seconds) => {
        if (seconds < 60) return `${seconds}秒`
        if (seconds < 3600) return `${Math.round(seconds / 60)}分钟`
        if (seconds < 86400) return `${Math.round(seconds / 3600)}小时`
        return `${Math.round(seconds / 86400)}天`
    }

    return (
        <div style={{ padding: '24px' }}>
            <Row gutter={[16, 16]} style={{ marginBottom: '24px' }}>
                <Col span={24}>
                    <Card>
                        <Row gutter={[16, 16]} align="middle">
                            <Col flex="auto">
                                <Title level={4} style={{ margin: 0 }}>工单统计报表</Title>
                            </Col>
                            <Col>
                                <RangePicker
                                    value={dateRange}
                                    onChange={setDateRange}
                                    format="YYYY-MM-DD"
                                    style={{ marginRight: '16px' }}
                                />
                            </Col>
                            <Col>
                                <Select
                                    value={trendPeriod}
                                    onChange={setTrendPeriod}
                                    style={{ width: 120, marginRight: '16px' }}
                                >
                                    <Option value="day">按天</Option>
                                    <Option value="week">按周</Option>
                                    <Option value="month">按月</Option>
                                </Select>
                            </Col>
                            <Col>
                                <Button
                                    icon={<ReloadOutlined />}
                                    onClick={fetchStatistics}
                                    loading={loading}
                                >
                                    刷新
                                </Button>
                            </Col>
                            <Col>
                                <Button
                                    type="primary"
                                    icon={<DownloadOutlined />}
                                    onClick={handleExport}
                                >
                                    导出报表
                                </Button>
                            </Col>
                        </Row>
                    </Card>
                </Col>
            </Row>

            <Spin spinning={loading}>
                <Row gutter={[16, 16]}>
                    <Col xs={24} sm={12} md={6}>
                        <Card>
                            <Statistic
                                title="总工单数"
                                value={statistics?.totalCount || 0}
                                prefix={<BarChartOutlined />}
                                valueStyle={{ color: '#1890ff' }}
                            />
                        </Card>
                    </Col>
                    <Col xs={24} sm={12} md={6}>
                        <Card>
                            <Statistic
                                title="待处理"
                                value={statistics?.pendingCount || 0}
                                prefix={<LineChartOutlined />}
                                valueStyle={{ color: '#faad14' }}
                            />
                        </Card>
                    </Col>
                    <Col xs={24} sm={12} md={6}>
                        <Card>
                            <Statistic
                                title="处理中"
                                value={statistics?.processingCount || 0}
                                prefix={<PieChartOutlined />}
                                valueStyle={{ color: '#52c41a' }}
                            />
                        </Card>
                    </Col>
                    <Col xs={24} sm={12} md={6}>
                        <Card>
                            <Statistic
                                title="逾期工单"
                                value={statistics?.overdueCount || 0}
                                valueStyle={{ color: '#f5222d' }}
                            />
                        </Card>
                    </Col>
                </Row>

                <Row gutter={[16, 16]} style={{ marginTop: '24px' }}>
                    <Col xs={24} sm={12}>
                        <Card title="优先级分布">
                            <ReactECharts option={priorityOption} style={{ height: '300px' }} />
                        </Card>
                    </Col>
                    <Col xs={24} sm={12}>
                        <Card title="状态分布">
                            <ReactECharts option={statusOption} style={{ height: '300px' }} />
                        </Card>
                    </Col>
                </Row>

                <Row gutter={[16, 16]} style={{ marginTop: '24px' }}>
                    <Col span={24}>
                        <Card title="工单趋势">
                            <ReactECharts option={trendOption} style={{ height: '350px' }} />
                        </Card>
                    </Col>
                </Row>

                <Row gutter={[16, 16]} style={{ marginTop: '24px' }}>
                    <Col xs={24} sm={12}>
                        <Card title="平均响应时间">
                            <ReactECharts option={responseTimeOption} style={{ height: '300px' }} />
                        </Card>
                    </Col>
                    <Col xs={24} sm={12}>
                        <Card title="平均解决时间">
                            <ReactECharts option={resolutionTimeOption} style={{ height: '300px' }} />
                        </Card>
                    </Col>
                </Row>

                <Row gutter={[16, 16]} style={{ marginTop: '24px' }}>
                    <Col span={24}>
                        <Card title="用户效率统计">
                            <Table
                                columns={userColumns}
                                dataSource={statistics?.userStats || []}
                                rowKey="userId"
                                pagination={false}
                                size="middle"
                            />
                        </Card>
                    </Col>
                </Row>

                <Row gutter={[16, 16]} style={{ marginTop: '24px' }}>
                    <Col span={24}>
                        <Card title="效率指标">
                            <Row gutter={[16, 16]}>
                                <Col xs={24} sm={12} md={6}>
                                    <Statistic
                                        title="平均响应时间"
                                        value={formatTime(statistics?.avgResponseTime || 0)}
                                    />
                                </Col>
                                <Col xs={24} sm={12} md={6}>
                                    <Statistic
                                        title="平均解决时间"
                                        value={formatTime(statistics?.avgResolution || 0)}
                                    />
                                </Col>
                                <Col xs={24} sm={12} md={6}>
                                    <Statistic
                                        title="SLA达成率"
                                        value={((statistics?.slaRate || 0) * 100).toFixed(2)}
                                        suffix="%"
                                        valueStyle={{
                                            color: (statistics?.slaRate || 0) >= 0.9 ? '#52c41a' :
                                                   (statistics?.slaRate || 0) >= 0.7 ? '#faad14' : '#f5222d'
                                        }}
                                    />
                                </Col>
                                <Col xs={24} sm={12} md={6}>
                                    <Statistic
                                        title="已关闭工单"
                                        value={statistics?.closedCount || 0}
                                        valueStyle={{ color: '#52c41a' }}
                                    />
                                </Col>
                            </Row>
                        </Card>
                    </Col>
                </Row>
            </Spin>
        </div>
    )
}