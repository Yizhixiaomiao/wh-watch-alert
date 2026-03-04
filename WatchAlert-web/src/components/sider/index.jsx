import React, {useEffect, useState} from 'react';
import {
    UserOutlined,
    BellOutlined,
    PieChartOutlined,
    NotificationOutlined,
    CalendarOutlined,
    DashboardOutlined,
    DeploymentUnitOutlined,
    AreaChartOutlined,
    FileDoneOutlined,
    SettingOutlined,
    ExceptionOutlined,
    ApiOutlined, TeamOutlined, DownOutlined, LogoutOutlined,
    FileTextOutlined,
    ThunderboltOutlined
} from '@ant-design/icons';
import {Link, useNavigate} from 'react-router-dom';
import {Menu, Layout, Typography, Dropdown, Space, message, Spin, theme, Popover, Avatar, Divider} from 'antd';
import logoIcon from "../../img/logo.svg";
import {getUserInfo, getUserPermissions} from "../../api/user";
import {getTenantList} from "../../api/tenant";

const { SubMenu } = Menu;
const { Sider } = Layout;

// 菜单路径到权限Key的映射（使用英文key，对应数据库中权限的key字段）
// 与 src/components/index.jsx 中的 ROUTE_PERMISSION_MAP 保持一致
const MENU_PERMISSION_MAP = {
    '/': null,  // 首页不需要权限
    '/ruleGroup': 'ruleGroupList',
    '/tmplType/Metrics/group': 'ruleTmplGroupList',
    '/subscribes': 'listSubscribe',
    '/alert/simulator': 'ruleList', // 使用通用权限替代
    '/faultCenter': 'faultCenterList',
    '/ticket': 'ticketList',
    '/ticket/create': 'ticketCreate',
    '/ticket/repair': 'ticketList', // 人工报修可能使用工单查看权限
    '/ticket/review': 'ticketList', // 工单评审使用工单查看权限
    '/ticket/statistics': 'ticketGetStatistics',
    '/ticket/sla': 'ticketList', // SLA策略使用工单查看权限
    '/ticket/template': 'ticketList', // 工单模板使用工单查看权限
    '/ticket/workHours': 'workHoursList',
    '/knowledge': 'knowledgeList',
    '/knowledge/create': 'knowledgeCreate',
    '/assignment-rule': 'assignmentRuleList',
    '/noticeObjects': 'noticeList',
    '/noticeTemplate': 'noticeTemplateList',
    '/noticeRecords': 'noticeRecordList',
    '/dutyManage': 'dutyManageList',
    '/probing': 'listProbing',
    '/onceProbing': 'createProbing', // 及时拨测可能使用创建拨测规则权限
    '/datasource': 'dataSourceList',
    '/folders': 'listFolder',
    '/user': 'userList',
    '/userRole': 'roleList',
    '/tenants': 'getTenantList',
    '/auditLog': 'userList', // 日志审计使用用户列表权限作为替代
    '/settings': 'getSystemSetting',
};

// 完整的菜单定义（包含所有可能的菜单项）
const allMenuItems = [
    { key: '1', path: '/', icon: <AreaChartOutlined />, label: '概览', permission: null },
    {
        key: '2',
        icon: <BellOutlined />,
        label: '告警管理',
        children: [
            { key: '2-1', path: '/ruleGroup', label: '告警规则', permission: 'ruleGroupList' },
            { key: '2-5', path: '/tmplType/Metrics/group', label: '规则模版', permission: 'ruleTmplGroupList' },
            { key: '2-6', path: '/subscribes', label: '告警订阅', permission: 'listSubscribe' },
            { key: '2-7', path: '/alert/simulator', label: '告警模拟器', permission: 'ruleList' }
        ]
    },
    { key: '12', path: '/faultCenter', icon: <ExceptionOutlined />, label: '故障中心', permission: 'faultCenterList' },
    {
        key: '13',
        icon: <FileTextOutlined />,
        label: '工单管理',
        children: [
            { key: '13-1', path: '/ticket', label: '工单列表', permission: 'ticketList' },
            { key: '13-2', path: '/ticket/create', label: '创建工单', permission: 'ticketCreate' },
            { key: '13-3', path: '/ticket/repair', label: '人工报修', permission: 'ticketList' },
            { key: '13-4', path: '/ticket/review', label: '工单评审', permission: 'ticketList' },
            { key: '13-5', path: '/ticket/statistics', label: '工单统计', permission: 'ticketGetStatistics' },
            { key: '13-6', path: '/ticket/sla', label: 'SLA策略', permission: 'ticketList' },
            { key: '13-7', path: '/ticket/template', label: '工单模板', permission: 'ticketList' },
            { key: '13-8', path: '/ticket/workHours', label: '工时标准', permission: 'workHoursList' }
        ]
    },
    {
        key: '14',
        icon: <FileTextOutlined />,
        label: '知识库',
        children: [
            { key: '14-1', path: '/knowledge', label: '知识列表', permission: 'knowledgeList' },
            { key: '14-2', path: '/knowledge/create', label: '创建知识', permission: 'knowledgeCreate' }
        ]
    },
    {
        key: '15',
        icon: <ThunderboltOutlined />,
        label: '智能派单',
        children: [
            { key: '15-1', path: '/assignment-rule', label: '派单规则', permission: 'assignmentRuleList' }
        ]
    },
    {
        key: '3',
        icon: <NotificationOutlined />,
        label: '通知管理',
        children: [
            { key: '3-1', path: '/noticeObjects', label: '通知对象', permission: 'noticeList' },
            { key: '3-2', path: '/noticeTemplate', label: '通知模版', permission: 'noticeTemplateList' },
            { key: '3-3', path: '/noticeRecords', label: '通知记录', permission: 'noticeRecordList' }
        ]
    },
    { key: '4', path: '/dutyManage', icon: <CalendarOutlined />, label: '值班中心', permission: 'dutyManageList' },
    {
        key: '11',
        icon: <ApiOutlined />,
        label: '网络分析',
        children: [
            { key: '11-1', path: '/probing', label: '拨测任务', permission: 'listProbing' },
            { key: '11-2', path: '/onceProbing', label: '及时拨测', permission: 'createProbing' }
        ]
    },
    { key: '6', path: '/datasource', icon: <PieChartOutlined />, label: '数据源', permission: 'dataSourceList' },
    { key: '8', path: '/folders', icon: <DashboardOutlined />, label: '仪表盘', permission: 'listFolder' },
    {
        key: '5',
        icon: <UserOutlined />,
        label: '人员组织',
        children: [
            { key: '5-1', path: '/user', label: '用户管理', permission: 'userList' },
            { key: '5-2', path: '/userRole', label: '角色管理', permission: 'roleList' }
        ]
    },
    { key: '7', path: '/tenants', icon: <DeploymentUnitOutlined />, label: '租户管理', permission: 'getTenantList' },
    { key: '9', path: '/auditLog', icon: <FileDoneOutlined />, label: '日志审计', permission: 'userList' },
    { key: '10', path: '/settings', icon: <SettingOutlined />, label: '系统设置', permission: 'getSystemSetting' }
];

export const ComponentSider = () => {
    const navigate = useNavigate();
    const [selectedMenuKey, setSelectedMenuKey] = useState('');
    const [userInfo, setUserInfo] = useState(null)
    const [loading, setLoading] = useState(true)
    const [tenantList, setTenantList] = useState([])
    const [getTenantStatus, setTenantStatus] = useState(null)
    const [userPermissions, setUserPermissions] = useState({})

    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken()

    const handleMenuClick = (key, path) => {
        if (path) {
            setSelectedMenuKey(key);
            navigate(path);
        }
    };

    // 检查是否有权限访问某个菜单项
    const hasPermission = (permissionKey) => {
        // Admin用户有所有权限
        if (userInfo?.role === 'admin') {
            return true;
        }

        // 如果没有指定权限要求，则允许访问
        if (!permissionKey) {
            return true;
        }

        // 检查用户是否有该权限
        const allPerms = Array.isArray(userPermissions) ? userPermissions : Object.values(userPermissions).flat();
        return allPerms.some(perm => perm.key === permissionKey);
    };

    // 根据权限过滤菜单项
    const filterMenuItems = (items) => {
        return items
            .map(item => {
                if (item.children) {
                    // 过滤子菜单
                    const filteredChildren = item.children.filter(child =>
                        hasPermission(child.permission)
                    );

                    // 如果过滤后没有子菜单，则不显示父菜单
                    if (filteredChildren.length === 0) {
                        return null;
                    }

                    return {
                        ...item,
                        children: filteredChildren
                    };
                }

                // 检查单个菜单项的权限
                if (!hasPermission(item.permission)) {
                    return null;
                }

                return item;
            })
            .filter(item => item !== null);
    };

    const renderMenuItems = (items) => {
        return items.map(item => {
            if (item.children) {
                return (
                    <SubMenu key={item.key} icon={item.icon} title={item.label}>
                        {item.children.map(child => (
                            <Menu.Item
                                key={child.key}
                                onClick={() => handleMenuClick(child.key, child.path)}
                            >
                                {child.label}
                            </Menu.Item>
                        ))}
                    </SubMenu>
                );
            }
            return (
                <Menu.Item
                    key={item.key}
                    icon={item.icon}
                    onClick={() => handleMenuClick(item.key, item.path)}
                >
                    {item.label}
                </Menu.Item>
            );
        });
    };

    const handleLogout = () => {
        localStorage.clear()
        navigate("/login")
    }

    const userMenu = (
        <Menu mode="vertical">
            <Menu.Item key="profile" icon={<UserOutlined />}>
                <Link to="/profile">个人信息</Link>
            </Menu.Item>
            <Menu.Divider />
            <Menu.Item key="logout" icon={<LogoutOutlined />} onClick={handleLogout} danger>
                退出登录
            </Menu.Item>
        </Menu>
    )

    useEffect(() => {
        fetchUserInfo()
    }, [])

    const fetchUserInfo = async () => {
        try {
            const res = await getUserInfo()
            setUserInfo(res.data)

            if (res.data.userid) {
                await fetchTenantList(res.data.userid)
                // 获取用户权限
                await fetchUserPermissions(res.data.userid, res.data.role)
            }

            setLoading(false)
        } catch (error) {
            console.error("Failed to fetch user info:", error)
            window.localStorage.removeItem("Authorization")
            navigate("/login")
        }
    }

    const fetchUserPermissions = async (userid, role) => {
        // Admin用户不需要获取权限
        if (role === 'admin') {
            setUserPermissions({});
            return;
        }

        try {
            const params = {
                userid: userid,
            };
            const res = await getUserPermissions(params);
            if (res.code === 200 && res.data) {
                setUserPermissions(res.data);
            } else {
                console.log("获取用户权限失败，将使用角色默认权限:", res.data || res.msg);
                setUserPermissions({});
            }
        } catch (error) {
            console.log("获取用户权限失败，将使用角色默认权限:", error.message);
            setUserPermissions({});
        }
    }

    const fetchTenantList = async (userid) => {
        try {
            const params = {
                userId: userid,
            }
            const res = await getTenantList(params)

            if (res.data === null || res.data.length === 0) {
                message.error("该用户没有可用租户")
                return
            }

            const opts = res.data.map((key, index) => ({
                label: key.name,
                value: key.id,
                index: index,
            }))

            setTenantList(opts)

            if (getTenantName() === null && opts.length > 0) {
                localStorage.setItem("TenantName", opts[0].label)
                localStorage.setItem("TenantID", opts[0].value)
                localStorage.setItem("TenantIndex", opts[0].index)
            }

            setTenantStatus(true)
        } catch (error) {
            console.error("Failed to fetch tenant list:", error)
            localStorage.clear()
            message.error("获取租户错误, 退出登录")
        }
    }

    const getTenantName = () => {
        return localStorage.getItem("TenantName")
    }

    const getTenantIndex = () => {
        return localStorage.getItem("TenantIndex")
    }

    const changeTenant = (c) => {
        localStorage.setItem("TenantIndex", c.key)
        if (c.item.props.name) {
            localStorage.setItem("TenantName", c.item.props.name)
        }
        if (c.item.props.value) {
            localStorage.setItem("TenantID", c.item.props.value)
        }

        setSelectedMenuKey('1')
        navigate('/')
        window.location.reload();
    }

    const tenantMenu = (
        <Menu selectable defaultSelectedKeys={[getTenantIndex()]} onSelect={changeTenant}>
            {tenantList.map((item) => (
                <Menu.Item key={item.index} name={item.label} value={item.value}>
                    {item.label}
                </Menu.Item>
            ))}
        </Menu>
    )

    if (loading || !getTenantStatus) {
        return (
            <div
                style={{
                    height: "100vh",
                    display: "flex",
                    justifyContent: "center",
                    alignItems: "center",
                    background: colorBgContainer,
                }}
            >
                <Spin tip="加载中..." size="large" />
            </div>
        )
    }

    // 根据权限过滤菜单
    const filteredMenuItems = filterMenuItems(allMenuItems);

    return (
        <Sider
            style={{
                overflow: 'hidden',
                height: '100%',
                background: '#000',
                borderRadius: '12px',
                display: 'flex',
                flexDirection: 'column',
                position: 'relative',
            }}
            theme="dark"
        >
            {/* 顶部Logo和租户选择区域 */}
            <div style={{
                padding: '16px 16px 0',
                position: 'sticky',
                top: 0,
                zIndex: 1,
            }}>
                <div style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    marginBottom: 16,
                    marginTop: '-70px',
                }}>
                    <img
                        src={logoIcon || "/placeholder.svg"}
                        alt="WatchAlert Logo"
                        style={{ width: "220px", height: "200px", borderRadius: "8px" }}
                    />
                </div>

                <Dropdown overlay={tenantMenu} trigger={["click"]} placement="bottomLeft">
                    <div style={{
                        display: 'flex',
                        marginTop: '-40px',
                        alignItems: 'center',
                        padding: '8px 12px',
                        borderRadius: '4px',
                        cursor: 'pointer',
                        background: 'rgba(255, 255, 255, 0.1)',
                        marginBottom: '16px',
                        ':hover': {
                            background: 'rgba(255, 255, 255, 0.2)',
                        }
                    }}>
                        <TeamOutlined style={{color: '#fff', fontSize: '14px', marginRight: '8px'}}/>
                        <Typography.Text
                            style={{color: '#fff', fontSize: '14px', flex: 1, overflow: 'hidden', textOverflow: 'ellipsis'}}>
                            {getTenantName()}
                        </Typography.Text>
                        <DownOutlined style={{color: '#fff', fontSize: '12px'}}/>
                    </div>
                </Dropdown>
            </div>

            <Divider style={{margin: '0', background: 'rgba(255, 255, 255, 0.1)'}}/>

            {/* 主内容，预留底部空间 */}
            <div
                style={{
                    textAlign:'left',
                    alignItems: 'flex-start',
                    overflowY: 'auto',
                    flex: 1,
                    height: '76vh',
                    paddingBottom: 70, // 预留底部空间
                }}
            >
                <Menu
                    theme="dark"
                    mode="inline"
                    selectedKeys={[selectedMenuKey]}
                    style={{ background: 'transparent'}}
                >
                    {renderMenuItems(filteredMenuItems)}
                </Menu>
            </div>

            {/* 绝对定位底部用户信息 */}
            <div style={{
                position: 'absolute',
                left: 0,
                bottom: 0,
                width: '100%',
                padding: '10px',
                borderTop: '1px solid rgba(255, 255, 255, 0.1)',
                background: '#000',
            }}>
                <Popover content={userMenu} trigger="click" placement="topRight">
                    <div style={{
                        display: "flex",
                        alignItems: "center",
                        cursor: "pointer",
                        padding: '8px',
                        borderRadius: '4px',
                        width: '100%',
                    }}>
                        <Avatar
                            style={{
                                backgroundColor: "#7265e6",
                                display: "flex",
                                alignItems: "center",
                                justifyContent: "center",
                            }}
                            size="default"
                        >
                            {userInfo.username ? userInfo.username.charAt(0).toUpperCase() : ""}
                        </Avatar>
                        <div style={{marginLeft: "12px", overflow: 'hidden'}}>
                            <Typography.Text style={{color: "#FFFFFFA6", display: 'block'}}>
                                {userInfo.username}
                            </Typography.Text>
                        </div>
                    </div>
                </Popover>
            </div>
        </Sider>
    );
};