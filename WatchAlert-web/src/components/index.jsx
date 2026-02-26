 "use client"

import { useState, useEffect, useRef } from "react"
import {Layout, theme, Button, Typography, Spin, Result, message} from "antd"
import { LeftOutlined, LoadingOutlined } from "@ant-design/icons"
import "./index.css"
import { ComponentSider } from "./sider"
import Auth from "../utils/Auth"
import {getTenantList} from "../api/tenant";
import {getUserInfo, getUserPermissions} from "../api/user";

// 菜单路径到权限Key的映射（使用英文key，对应数据库中权限的key字段）
// 与 src/components/sider/index.jsx 中的 MENU_PERMISSION_MAP 保持一致
const ROUTE_PERMISSION_MAP = {
    '/ruleGroup': 'ruleGroupList',
    '/tmplType/Metrics/group': 'ruleTmplGroupList',
    '/subscribes': 'listSubscribe',
    '/alert/simulator': 'ruleList', // 使用通用权限替代
    '/faultCenter': 'faultCenterList',
    '/ticket': 'ticketList',
    '/ticket/create': 'ticketCreate',
    '/ticket/repair': 'ticketList', // 人工报修可能使用工单查看权限
    '/ticket/statistics': 'ticketGetStatistics',
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

const Components = (props) => {
    const { name, c } = props
    const { Content } = Layout
    const [loading, setLoading] = useState(true)
    const [authorization, setAuthorization] = useState(null)
    const [tenantId, setTenantId] = useState(null)
    const [error, setError] = useState(false)
    const [isRendered, setIsRendered] = useState(false)
    const [userPermissions, setUserPermissions] = useState({})
    const [userInfo, setUserInfo] = useState(null)
    const contentRef = useRef(null)
    // Add permission error state
    const [permissionError, setPermissionError] = useState(false)

    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken()

// Check route permission based on current path
 const checkRoutePermission = (path, permissions, userInfo) => {
        const permissionKey = ROUTE_PERMISSION_MAP[path];
        if (!permissionKey) {
            return true;
        }
        
        // Admin users have all permissions
        if (userInfo?.role === 'admin') {
            return true;
        }
        
        // 检查是否已成功获取用户权限
        const hasPermissionsLoaded = permissions && Object.keys(permissions).length > 0;
        
        if (hasPermissionsLoaded) {
            // permsList API返回扁平数组，直接检查权限
            const allPerms = Array.isArray(permissions) ? permissions : Object.values(permissions).flat();
            return allPerms.some(perm => perm.key === permissionKey);
        }
        
        // 权限未加载时，默认拒绝访问（等待权限加载完成）
        return false;
    };

    // Fetch user info and permissions
    useEffect(() => {
        let isMounted = true

        const checkAuthAndTenant = async () => {
            try {
                const auth = localStorage.getItem("Authorization")
                const tenant = localStorage.getItem("TenantID")

                if (!auth) {
                    setError(true)
                    setLoading(false)
                    return
                }

                // Get user info
                try {
                    const userRes = await getUserInfo();
                    if (userRes.data) {
                        setUserInfo(userRes.data);
                        
                        // Get user permissions
                        await fetchUserPermissions(userRes.data.userid, userRes.data.role);
                    }
                } catch (userErr) {
                    console.error("Failed to fetch user info:", userErr);
                    setError(true);
                    setLoading(false);
                    return;
                }

                // If has auth but no tenant, try to get it
                if (auth && !tenant) {
                    try {
                        const userRes = await getUserInfo()
                        if (userRes.data?.userid) {
                            await fetchTenantList(userRes.data.userid)
                        }
                    } catch (err) {
                        console.error("Failed to fetch user info:", err)
                        setError(true)
                    }
                }

                // Re-check tenant info
                const updatedTenant = localStorage.getItem("TenantID")
                if (isMounted) {
                    setAuthorization(auth)
                    setTenantId(updatedTenant)
                    setError(!updatedTenant)
                }
            } catch (error) {
                console.error("Error accessing localStorage:", error)
                if (isMounted) {
                    setError(true)
                }
            }
            
            setLoading(false);
        }

        // Delay check logic
        const delayCheck = setTimeout(() => {
            checkAuthAndTenant()
        }, 500) // Delay 500 milliseconds

        return () => {
            isMounted = false
            clearTimeout(delayCheck)
        }
    }, [])

    // Fetch user permissions
    const fetchUserPermissions = async (userid, role) => {
        // Admin users get all access, no need to fetch permissions
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
                // If getting permissions fails, log the issue and set empty permissions
                // The permission check functions will handle empty permissions appropriately
                console.log("获取用户权限失败，将使用角色默认权限:", res.data || res.msg);
                setUserPermissions({});
            }
        } catch (error) {
            // If API call fails, log the error and set empty permissions
            console.log("获取用户权限失败，将使用角色默认权限:", error.message);
            setUserPermissions({});
        }
    }

    // Check permission after loading user info and permissions
    useEffect(() => {
        if (!loading && userInfo && userPermissions) {
            // Get current path
            const currentPath = window.location.pathname;
            
            // Check if user has permission for this route
            if (!checkRoutePermission(currentPath, userPermissions, userInfo)) {
                setPermissionError(true);
            }
        }
    }, [loading, userInfo, userPermissions]);

    // Listen for render completion
    useEffect(() => {
        if (loading || error || permissionError || !authorization || !tenantId) return

        const observer = new MutationObserver(() => {
            const timer = setTimeout(() => {
                setIsRendered(true)
                observer.disconnect()
            }, 300)

            return () => clearTimeout(timer)
        })

        if (contentRef.current) {
            observer.observe(contentRef.current, {
                childList: true,
                subtree: true,
                attributes: true,
                characterData: true
            })
        }

        const maxWaitTimer = setTimeout(() => {
            setIsRendered(true)
            observer.disconnect()
        }, 1500)

        return () => {
            observer.disconnect()
            clearTimeout(maxWaitTimer)
        }
    }, [loading, error, permissionError, authorization, tenantId])

    const goBackPage = () => {
        window.history.back()
    }

    const fetchTenantList = async (userid) => {
        const auth = localStorage.getItem("Authorization")
        if (!auth) {
            console.error("Authorization token is missing")
            setError(true)
            return
        }

        try {
            const res = await getTenantList({ userId: userid })

            if (!res?.data || !Array.isArray(res.data) || res.data.length === 0) {
                console.error("No tenant data available")
                setError(true)
                return
            }

            const tenantOptions = res.data.map((tenant, index) => ({
                label: tenant.name,
                value: tenant.id,
                index: index,
            }))

            // Set first tenant as default
            const firstTenant = tenantOptions[0]
            localStorage.setItem("TenantName", firstTenant.label)
            localStorage.setItem("TenantID", firstTenant.value)
            localStorage.setItem("TenantIndex", firstTenant.index)

            return tenantOptions
        } catch (error) {
            console.error("Failed to fetch tenant list:", error)
            setError(true)
            throw error
        }
    }

    // Full screen loading component
    const FullScreenLoading = () => (
        <div
            style={{
                position: 'fixed',
                top: 0,
                left: 0,
                right: 0,
                bottom: 0,
                display: 'flex',
                justifyContent: 'center',
                alignItems: 'center',
                flexDirection: 'column',
                backgroundColor: 'white',
                zIndex: 9999,
                transition: 'opacity 0.3s ease-out'
            }}
        >
            <Spin indicator={<LoadingOutlined style={{ fontSize: 40 }} spin />} size="large" />
            <Typography.Text type="secondary" style={{ marginTop: 16 }}>
                {loading ? "正在验证用户信息..." : "页面准备中..."}
            </Typography.Text>
        </div>
    )

    // Error component
    const ErrorScreen = () => (
        <div
            style={{
                height: "100vh",
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                background: "#f0f2f5",
            }}
        >
            <Result
                status="error"
                title={!authorization ? "用户无效" : "租户无效"}
                subTitle={!authorization ? "请先登录系统" : "未获取到有效的租户"}
                extra={[
                    <Button
                        type="primary"
                        key="login"
                        onClick={() => {
                            localStorage.clear()
                            window.location.href = "/login"
                        }}
                    >
                        返回登录
                    </Button>,
                ]}
            />
        </div>
    )

    // Permission error component
    const PermissionErrorScreen = () => (
        <div
            style={{
                height: "100vh",
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                background: "#f0f2f5",
            }}
        >
            <Result
                status="403"
                title="无访问权限"
                subTitle="您没有权限访问此页面"
                extra={[
                    <Button
                        type="primary"
                        key="back"
                        onClick={() => {
                            window.history.back();
                        }}
                    >
                        返回上一页
                    </Button>,
                    <Button
                        key="home"
                        onClick={() => {
                            window.location.href = "/";
                        }}
                    >
                        返回首页
                    </Button>,
                ]}
            />
        </div>
    )

    // Priority: Loading > Permission Error > Auth Error > Render Content
    if (loading) {
        return <FullScreenLoading />
    }

    if (permissionError) {
        return <PermissionErrorScreen />;
    }

    if (error || !authorization || !tenantId) {
        return <ErrorScreen />
    }

    // Check permission again before rendering content to ensure real-time check
    const currentPath = window.location.pathname;
    if (!checkRoutePermission(currentPath, userPermissions, userInfo)) {
        return <PermissionErrorScreen />;
    }

    // Main content area
    return (
        <>
            {/* Only hide loading interface when content is fully rendered */}
            {!isRendered && <FullScreenLoading />}

            <Layout
                style={{
                    height: "100vh",
                    overflow: "hidden",
                    background: "#000000",
                    opacity: isRendered ? 1 : 0, // Add fade-in effect
                    transition: 'opacity 0.3s ease-in'
                }}
                ref={contentRef}
            >
                <Layout style={{ background: "transparent", marginTop: "16px" }}>
                    {/* Sidebar */}
                    <div
                        style={{
                            width: "220px",
                            borderRadius: borderRadiusLG,
                            overflow: "hidden",
                            boxShadow: "0 2px 8px rgba(0, 0, 0, 0.06)",
                            height: "100%",
                            background: "#000000",
                        }}
                    >
                        <div style={{ height: "100%", overflow: "auto", padding: "16px 0", marginLeft: "15px" }}>
                            <ComponentSider />
                        </div>
                    </div>

                    {/* Content area */}
                    <Layout style={{ background: "transparent", padding: "0 16px 0px 16px" }}>
                        <Content
                            style={{
                                background: colorBgContainer,
                                borderRadius: borderRadiusLG,
                                padding: "0",
                                height: "calc(100vh - 32px)",
                                overflow: "hidden",
                                boxShadow: "0 2px 8px rgba(0, 0, 0, 0.06)",
                            }}
                        >
                            {/* Page header */}
                            {name !== "off" && (
                                <div
                                    style={{
                                        padding: "16px 24px",
                                        borderBottom: "1px solid #f0f0f0",
                                        display: "flex",
                                        alignItems: "center",
                                        gap: "8px",
                                    }}
                                >
                                    <Button type="text" icon={<LeftOutlined />} onClick={goBackPage} style={{ padding: "4px" }} />
                                    <Typography.Title level={4} style={{ margin: 0, fontSize: "16px" }}>
                                        {name}
                                    </Typography.Title>
                                </div>
                            )}

                            {/* Main content */}
                            <div
                                style={{
                                    padding: name !== "off" ? "24px" : "0",
                                    height: name !== "off" ? "calc(100% - 53px)" : "100%",
                                    overflow: "auto",
                                }}
                            >
                                {c}
                            </div>
                        </Content>

                        {/* Footer */}
                        <div
                            style={{
                                textAlign: "center",
                                color: "#B1B1B1",
                                fontSize: "12px",
                            }}
                        >
                            Wh-Ops-Alert 卫华监控告警平台!
                        </div>
                    </Layout>
                </Layout>
            </Layout>
        </>
    )
}

export const ComponentsContent = Auth(Components)
export default ComponentsContent