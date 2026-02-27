"use client"

import { useState, useEffect } from "react"
import {
    Button,
    Tag,
    Space,
    Input,
    Dropdown,
    message,
    Tooltip,
    Popconfirm,
    Radio,
    DatePicker,
    Spin,
    Modal,
    Select,
} from "antd"
import {
    PlusOutlined,
    DeleteOutlined,
    DownOutlined,
    DownloadOutlined,
    EditOutlined,
    EyeOutlined,
    CheckCircleOutlined,
    CloseCircleOutlined,
} from "@ant-design/icons"
import { useNavigate } from "react-router-dom"
import {
    getTicketList,
    deleteTicket,
    claimTicket,
    closeTicket,
    assignTicket,
} from "../../api/ticket"
import { HandleApiError, HandleShowTotal, FormatTime } from "../../utils/lib"
import { VirtualList } from "../../components/VirtualList"
import { TableWithPagination } from "../../utils/TableWithPagination"
import { getUserList } from "../../api/user"
import { getUserInfo } from "../../api/user"
import { FaultCenterList } from "../../api/faultCenter"
import * as XLSX from 'xlsx'
import { CreateTicketDrawer } from "./CreateDrawer"
import { clearCacheByUrl } from "../../utils/http"

const { Search } = Input

export const TicketList = () => {
    const navigate = useNavigate()
    const [list, setList] = useState([])
    const [selectStatus, setSelectStatus] = useState("all")
    const [faultCenterList, setFaultCenterList] = useState([])
    const [userList, setUserList] = useState([])
    const [createDrawerVisible, setCreateDrawerVisible] = useState(false)
    const [dateRange, setDateRange] = useState(null)
    const [loading, setLoading] = useState(false)
    const [keyword, setKeyword] = useState("")
    const [selectedRowKeys, setSelectedRowKeys] = useState([])
    const [batchActionLoading, setBatchActionLoading] = useState(false)
    const [assignModalVisible, setAssignModalVisible] = useState(false)
    const [assignUserId, setAssignUserId] = useState("")
    const [pagination, setPagination] = useState({
        index: 1,
        size: 20,
        total: 0,
    })
    const [userInfo, setUserInfo] = useState(null)
    const [hideCreateButton, setHideCreateButton] = useState(false)

    // 状态映射
    const statusMap = {
        Pending: { color: "default", text: "待处理" },
        Assigned: { color: "blue", text: "已分配" },
        Processing: { color: "processing", text: "处理中" },
        Verifying: { color: "purple", text: "验证中" },
        Resolved: { color: "success", text: "已解决" },
        Closed: { color: "default", text: "已关闭" },
        Cancelled: { color: "error", text: "已取消" },
        Escalated: { color: "warning", text: "已升级" },
    }

    // 优先级映射
    const priorityMap = {
        P0: { color: "red", text: "P0-最高" },
        P1: { color: "orange", text: "P1-高" },
        P2: { color: "blue", text: "P2-中" },
        P3: { color: "green", text: "P3-低" },
        P4: { color: "default", text: "P4-最低" },
    }

    // 工单类型映射
    const typeMap = {
        Alert: { text: "告警工单" },
        Fault: { text: "故障工单" },
        Change: { text: "变更工单" },
        Query: { text: "咨询工单" },
    }

    useEffect(() => {
        handleList(pagination.index, pagination.size)
        fetchFaultCenterList()
        fetchUserList()
    }, [])

    useEffect(() => {
        handleList(1, pagination.size)
    }, [selectStatus])

    // 获取故障中心列表
    const fetchFaultCenterList = async () => {
        try {
            const res = await FaultCenterList({ page: 1, size: 1000 })
            if (res && res.data) {
                const faults = res.data.list || res.data || []
                setFaultCenterList(faults)
            }
        } catch (error) {
            console.error("获取故障中心列表失败:", error)
        }
    }

    // 获取用户信息
    const fetchUserInfo = async () => {
        try {
            const res = await getUserInfo()
            if (res && res.data) {
                setUserInfo(res.data)
                const userRole = res.data.role || res.data.userRole
                console.log("fetchUserInfo - 用户信息:", res.data)
                console.log("fetchUserInfo - 用户角色:", userRole)
                // viewer 和 on_call 角色隐藏创建工单按钮
                if (userRole === 'viewer' || userRole === 'on_call') {
                    setHideCreateButton(true)
                }
            }
        } catch (error) {
            console.error("获取用户信息失败:", error)
        }
    }

    // 获取用户列表
    const fetchUserList = async () => {
        try {
            const res = await getUserList({})
            if (res && res.data) {
                const users = Array.isArray(res.data) ? res.data : []
                setUserList(users)
            }
        } catch (error) {
            console.error("获取用户列表失败:", error)
        }
    }

    useEffect(() => {
        const auth = localStorage.getItem('Authorization')
        if (!auth) {
            navigate('/login')
            return
        }

        // 从 localStorage 获取用户信息
        const userInfoStr = localStorage.getItem('userInfo')
        if (userInfoStr) {
            try {
                const userData = JSON.parse(userInfoStr)
                setUserInfo(userData)
                const userRole = userData.role || userData.userRole
                console.log("useEffect - 用户信息:", userData)
                console.log("useEffect - 用户角色:", userRole)
                // viewer 和 on_call 角色隐藏创建工单按钮
                if (userRole === 'viewer' || userRole === 'on_call') {
                    setHideCreateButton(true)
                    console.log("useEffect - 设置隐藏创建工单按钮")
                }
            } catch (e) {
                console.error("解析用户信息失败:", e)
            }
        } else {
            console.log("useEffect - localStorage 中没有 userInfo")
        }

        fetchUserInfo()
        handleList(pagination.index, pagination.size)
        fetchFaultCenterList()
        fetchUserList()
    }, [])

    // 获取工单列表
    const handleList = async (index, size, startTime, endTime, keyword = "") => {
        setLoading(true)
        try {
            const params = {
                page: index,
                size: size,
                status: selectStatus === "all" ? "" : selectStatus,
            }
            if (startTime && endTime) {
                params.startTime = startTime
                params.endTime = endTime
            }
            if (keyword) {
                params.keyword = keyword
            }
            const res = await getTicketList(params, { skipCache: true })
            if (res && res.data) {
                setPagination({
                    index: res.data.index || index,
                    size: res.data.size || size,
                    total: res.data.total || 0,
                })
                setList(res.data.list || [])
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setLoading(false)
        }
    }

    // 搜索
    const onSearch = async (value) => {
        handleList(1, pagination.size, "", "", value)
    }

    // 状态切换
    const changeStatus = ({ target: { value } }) => {
        setPagination({ ...pagination, index: 1 })
        setSelectStatus(value)
    }

    const canCreateTicket = () => {
        if (!userInfo) {
            console.log("检查创建工单权限 - 用户信息未加载")
            return false
        }
        const userRole = userInfo.role || userInfo.userRole
        const hasPermission = userRole && userRole !== 'viewer' && userRole !== 'on_call'
        console.log("检查创建工单权限:", { 
            role: userInfo.role, 
            userRole: userInfo.userRole, 
            userRole_used: userRole, 
            canCreate: hasPermission 
        })
        return hasPermission
    }

    // 删除工单
    const handleDelete = async (ticketId) => {
        try {
            const res = await deleteTicket({ ticketId })
            message.success('工单删除成功')
            clearCacheByUrl('/api/w8t/ticket')
            clearCacheByUrl('/api/w8t/ticket/list')
            await handleList(pagination.index, pagination.size)
        } catch (error) {
            HandleApiError(error)
        }
    }

    // 认领工单
    const handleClaim = async (ticketId) => {
        try {
            await claimTicket({ ticketId })
            message.success('工单认领成功')
            clearCacheByUrl('/api/w8t/ticket')
            clearCacheByUrl('/api/w8t/ticket/list')
            await handleList(pagination.index, pagination.size)
        } catch (error) {
            HandleApiError(error)
        }
    }

    // 导出Excel
    const handleExport = async () => {
        try {
            message.loading({ content: "正在导出...", key: "export" })

            const params = {
                page: 1,
                size: 10000,
                status: selectStatus === "all" ? "" : selectStatus,
            }
            // 如果选择了日期范围，添加时间参数
            if (dateRange && dateRange.length === 2) {
                params.startTime = dateRange[0].startOf('day').unix()
                params.endTime = dateRange[1].endOf('day').unix()
            }
            const res = await getTicketList(params)

            if (!res || !res.data || !res.data.list || res.data.list.length === 0) {
                message.warning({ content: "没有数据可导出", key: "export" })
                return
            }

            const exportData = res.data.list.map(ticket => {
                const fault = faultCenterList.find(f => f.id === ticket.faultCenterId)
                const creator = userList.find(u => u.userid === ticket.createdBy)
                const assignee = userList.find(u => u.userid === ticket.assignedTo)
                return {
                    "标题": ticket.title || "",
                    "类型": typeMap[ticket.type]?.text || ticket.type || "",
                    "故障中心": fault ? (fault.name || fault.title || "-") : "-",
                    "优先级": priorityMap[ticket.priority]?.text || ticket.priority || "",
                    "状态": statusMap[ticket.status]?.text || ticket.status || "",
                    "创建人": creator ? (creator.username || creator.userid) : (ticket.createdBy || "-"),
                    "处理人": assignee ? (assignee.username || assignee.userid) : (ticket.assignedTo || "-"),
                    "创建时间": FormatTime(ticket.createdAt),
                    "更新时间": FormatTime(ticket.updatedAt),
                }
            })

            const ws = XLSX.utils.json_to_sheet(exportData)
            const wb = XLSX.utils.book_new()
            XLSX.utils.book_append_sheet(wb, ws, "工单列表")

            const colWidths = [
                { wch: 30 },  // 标题
                { wch: 12 },  // 类型
                { wch: 20 },  // 故障中心
                { wch: 15 },  // 优先级
                { wch: 12 },  // 状态
                { wch: 15 },  // 创建人
                { wch: 15 },  // 处理人
                { wch: 20 },  // 创建时间
                { wch: 20 },  // 更新时间
            ]
            ws['!cols'] = colWidths

            const fileName = `工单列表_${new Date().toLocaleDateString().replace(/\//g, '-')}.xlsx`
            XLSX.writeFile(wb, fileName)

            message.success({ content: `成功导出 ${exportData.length} 条数据`, key: "export" })
        } catch (error) {
            message.error({ content: "导出失败", key: "export" })
            HandleApiError(error)
        }
    }

    const handleBatchDelete = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请先选择要删除的工单')
            return
        }

        setBatchActionLoading(true)
        try {
            let successCount = 0
            let failCount = 0

            for (const ticketId of selectedRowKeys) {
                try {
                    await deleteTicket({ ticketId })
                    successCount++
                } catch (error) {
                    failCount++
                    console.error(`删除工单 ${ticketId} 失败:`, error)
                }
            }

            if (successCount > 0) {
                message.success(`成功删除 ${successCount} 个工单${failCount > 0 ? `，失败 ${failCount} 个` : ''}`)
                clearCacheByUrl('/api/w8t/ticket')
                await handleList(pagination.index, pagination.size)
                setSelectedRowKeys([])
            } else {
                message.error(`删除失败，共 ${failCount} 个工单`)
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setBatchActionLoading(false)
        }
    }

    const handleBatchClose = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请先选择要关闭的工单')
            return
        }

        setBatchActionLoading(true)
        try {
            let successCount = 0
            let failCount = 0

            for (const ticketId of selectedRowKeys) {
                try {
                    await closeTicket({ ticketId })
                    successCount++
                } catch (error) {
                    failCount++
                    console.error(`关闭工单 ${ticketId} 失败:`, error)
                }
            }

            if (successCount > 0) {
                message.success(`成功关闭 ${successCount} 个工单${failCount > 0 ? `，失败 ${failCount} 个` : ''}`)
                clearCacheByUrl('/api/w8t/ticket')
                clearCacheByUrl('/api/w8t/ticket/list')
                await handleList(pagination.index, pagination.size)
                setSelectedRowKeys([])
            } else {
                message.error(`关闭失败，共 ${failCount} 个工单`)
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setBatchActionLoading(false)
        }
    }

    const handleBatchAssign = async () => {
        if (selectedRowKeys.length === 0) {
            message.warning('请先选择要指派的工单')
            return
        }
        setAssignModalVisible(true)
    }

    const confirmBatchAssign = async () => {
        if (!assignUserId) {
            message.warning('请选择指派人员')
            return
        }

        setBatchActionLoading(true)
        try {
            let successCount = 0
            let failCount = 0

            for (const ticketId of selectedRowKeys) {
                try {
                    await assignTicket({ ticketId, assignedTo: assignUserId })
                    successCount++
                } catch (error) {
                    failCount++
                    console.error(`指派工单 ${ticketId} 失败:`, error)
                }
            }

            if (successCount > 0) {
                message.success(`成功指派 ${successCount} 个工单${failCount > 0 ? `，失败 ${failCount} 个` : ''}`)
                clearCacheByUrl('/api/w8t/ticket')
                clearCacheByUrl('/api/w8t/ticket/list')
                await handleList(pagination.index, pagination.size)
                setSelectedRowKeys([])
                setAssignModalVisible(false)
                setAssignUserId("")
            } else {
                message.error(`指派失败，共 ${failCount} 个工单`)
            }
        } catch (error) {
            HandleApiError(error)
        } finally {
            setBatchActionLoading(false)
        }
    }

    // 批量操作菜单
    const batchOperationMenu = {
        items: [
            {
                key: "batchAssign",
                label: "批量指派",
                icon: <EditOutlined />,
                onClick: handleBatchAssign,
                disabled: selectedRowKeys.length === 0,
            },
            {
                key: "batchDelete",
                label: "批量删除",
                icon: <DeleteOutlined />,
                onClick: handleBatchDelete,
                disabled: selectedRowKeys.length === 0,
            },
            {
                key: "batchClose",
                label: "批量关闭",
                icon: <CheckCircleOutlined />,
                onClick: handleBatchClose,
                disabled: selectedRowKeys.length === 0,
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

    // 日期范围变化处理
    const handleDateRangeChange = (dates) => {
        setDateRange(dates)
        if (dates && dates.length === 2) {
            // 开始时间：当天的 00:00:00
            const startTime = dates[0].startOf('day').unix()
            // 结束时间：当天的 23:59:59
            const endTime = dates[1].endOf('day').unix()
            handleList(1, pagination.size, startTime, endTime)
        } else {
            handleList(1, pagination.size)
        }
    }

    // 表格列定义
    const columns = [
        {
            title: "标题",
            dataIndex: "title",
            key: "title",
            width: 200,
            render: (text, record) => (
                <a onClick={() => navigate(`/ticket/detail/${record.ticketId}`)}>
                    <Tooltip placement="topLeft" title={text}>
                        <div style={{
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap',
                        }}>
                            {text}
                        </div>
                    </Tooltip>
                </a>
            ),
        },
        {
            title: "类型",
            dataIndex: "type",
            key: "type",
            width: 110,
            render: (type) => typeMap[type]?.text || type,
        },
        {
            title: "故障中心",
            dataIndex: "faultCenterId",
            key: "faultCenterId",
            width: 160,
            render: (faultCenterId) => {
                if (!faultCenterId) return "-"
                const fault = faultCenterList.find(f => f.id === faultCenterId)
                return fault ? fault.name || fault.title || faultCenterId : faultCenterId
            },
        },
        {
            title: "优先级",
            dataIndex: "priority",
            key: "priority",
            width: 110,
            render: (priority) => (
                <Tag color={priorityMap[priority]?.color}>
                    {priorityMap[priority]?.text || priority}
                </Tag>
            ),
        },
        {
            title: "状态",
            dataIndex: "status",
            key: "status",
            width: 110,
            render: (status) => (
                <Tag color={statusMap[status]?.color}>
                    {statusMap[status]?.text || status}
                </Tag>
            ),
        },
        {
            title: "创建人",
            dataIndex: "createdBy",
            key: "createdBy",
            width: 110,
            render: (userId) => {
                if (!userId) return "-"
                const user = userList.find(u => u.userid === userId)
                return user ? (user.username || userId) : userId
            },
        },
        {
            title: "处理人",
            dataIndex: "assignedTo",
            key: "assignedTo",
            width: 110,
            render: (userId) => {
                if (!userId) return "-"
                const user = userList.find(u => u.userid === userId)
                const displayName = user ? (user.username || userId) : userId
                return (
                    <Tag style={{
                        borderRadius: "12px",
                        padding: "0 10px",
                        fontSize: "12px",
                        fontWeight: "500",
                        display: "inline-flex",
                        alignItems: "center",
                        gap: "4px",
                    }}>
                        {displayName}
                    </Tag>
                )
            },
        },
        {
            title: "创建时间",
            dataIndex: "createdAt",
            key: "createdAt",
            width: 170,
            render: (time) => {
                return (
                    <div style={{ display: "flex", alignItems: "center", gap: "6px" }}>
                        <span>{FormatTime(time)}</span>
                    </div>
                )
            },
        },
        {
            title: "操作",
            dataIndex: "operation",
            fixed: "right",
            width: 120,
            render: (_, record) => (
                <Space size="middle">
                    <Tooltip title="查看">
                        <Button
                            type="link"
                            icon={<EyeOutlined />}
                            style={{ color: "#1677ff" }}
                            onClick={() => navigate(`/ticket/detail/${record.ticketId}`)}
                        />
                    </Tooltip>
                    {(record.status === "Pending" || record.status === "Assigned") && (
                        <Tooltip title="认领">
                            <Button
                                type="text"
                                icon={<EditOutlined />}
                                style={{ color: "#615454" }}
                                onClick={() => handleClaim(record.ticketId)}
                            />
                        </Tooltip>
                    )}
                    <Tooltip title="删除">
                        <Popconfirm
                            title="确定要删除此工单吗?"
                            onConfirm={() => handleDelete(record.ticketId)}
                            okText="确定"
                            cancelText="取消"
                            placement="left"
                        >
                            <Button
                                type="text"
                                icon={<DeleteOutlined />}
                                style={{ color: "red" }}
                            />
                        </Popconfirm>
                    </Tooltip>
                </Space>
            ),
        },
    ]

    return (
        <div style={{ display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
            <div style={{
                background: '#fff',
                borderRadius: '8px',
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                overflow: 'hidden'
            }}>
                <div style={{
                    display: "flex",
                    justifyContent: "space-between",
                    marginBottom: "20px",
                    alignItems: "center"
                }}>
                    <div style={{ display: "flex", gap: "10px" }}>
                        <Radio.Group
                            options={[
                                { label: "全部", value: "all" },
                                { label: "待处理", value: "Pending" },
                                { label: "处理中", value: "Processing" },
                                { label: "已解决", value: "Resolved" },
                                { label: "已关闭", value: "Closed" },
                            ]}
                            defaultValue={selectStatus}
                            onChange={changeStatus}
                            optionType="button"
                        />

                        <Search
                            allowClear
                            placeholder="输入搜索关键字"
                            onSearch={onSearch}
                            style={{ width: 300 }}
                        />
                    </div>
                    <div style={{ display: "flex", gap: "10px" }}>
                        <DatePicker.RangePicker
                            value={dateRange}
                            onChange={handleDateRangeChange}
                            placeholder={["开始日期", "结束日期"]}
                            style={{ width: 300 }}
                        />

                        <Dropdown menu={batchOperationMenu}>
                            <Button>
                                批量操作 {selectedRowKeys.length > 0 && `(${selectedRowKeys.length})`} <DownOutlined />
                            </Button>
                        </Dropdown>

                        <Button
                            type="primary"
                            icon={<PlusOutlined />}
                            onClick={() => setCreateDrawerVisible(true)}
                            style={{ backgroundColor: "#000000", display: hideCreateButton ? 'none' : '' }}
                        >
                            创建工单
                        </Button>
                    </div>
                </div>

                <CreateTicketDrawer
                    visible={createDrawerVisible}
                    onClose={() => setCreateDrawerVisible(false)}
                    onSuccess={() => {
                        setCreateDrawerVisible(false)
                        handleList(pagination.index, pagination.size)
                    }}
                />

                <Modal
                    title="批量指派工单"
                    open={assignModalVisible}
                    onOk={confirmBatchAssign}
                    onCancel={() => {
                        setAssignModalVisible(false)
                        setAssignUserId("")
                    }}
                    confirmLoading={batchActionLoading}
                    okText="确定指派"
                    cancelText="取消"
                >
                    <div style={{ marginTop: 20 }}>
                        <span style={{ marginRight: 10 }}>选择指派人员：</span>
                        <Select
                            style={{ width: 300 }}
                            placeholder="请选择指派人员"
                            value={assignUserId || undefined}
                            onChange={setAssignUserId}
                            showSearch
                            filterOption={(input, option) =>
                                (option?.label ?? '').toLowerCase().includes(input.toLowerCase())
                            }
                        >
                            {userList.map(user => (
                                <Select.Option
                                    key={user.userid}
                                    value={user.userid}
                                    label={user.username}
                                >
                                    {user.username}
                                </Select.Option>
                            ))}
                        </Select>
                    </div>
                </Modal>

                <TableWithPagination
                    columns={columns}
                    dataSource={list}
                    pagination={pagination}
                    onPageChange={(page, pageSize) => {
                        setPagination({ ...pagination, index: page, size: pageSize })
                        handleList(page, pageSize)
                    }}
                    onPageSizeChange={(current, pageSize) => {
                        setPagination({ ...pagination, index: current, size: pageSize })
                        handleList(current, pageSize)
                    }}
                    scrollY={'calc(100vh - 300px)'}
                    rowKey={record => record.ticketId}
                    showTotal={HandleShowTotal}
                    rowSelection={{
                        selectedRowKeys,
                        onChange: setSelectedRowKeys,
                    }}
                />
            </div>
        </div>
    )
}
