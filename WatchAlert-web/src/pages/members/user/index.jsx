import {Input, Table, Button, Popconfirm, Tooltip, Space, Modal, Select, Input as TextArea, message, Tag, Form, Typography, Dropdown} from 'antd';
import React, { useState, useEffect } from 'react';
import UserCreateModal from './UserCreateModal';
import UserChangePass from './UserChangePass';
import UserPermissionsModal from './UserPermissionsModal';
import { UserActivityLogs } from './UserActivityLogs';
import { deleteUser, getUserList, updateUserStatus, batchUserOperation } from '../../../api/user';
import {CopyOutlined, DeleteOutlined, EditOutlined, PlusOutlined, LockOutlined, UnlockOutlined, StopOutlined, CheckCircleOutlined, AppstoreOutlined, DownOutlined, HistoryOutlined} from "@ant-design/icons";
import {HandleShowTotal, HandleApiError} from "../../../utils/lib";
import {Link} from "react-router-dom";
import {copyToClipboard} from "../../../utils/copyToClipboard";
import { clearCacheByUrl } from "../../../utils/http"

const { Search } = Input;
const { Text } = Typography;

export const User = () => {
    const [selectedRow, setSelectedRow] = useState(null); // 当前选中行
    const [updateVisible, setUpdateVisible] = useState(false); // 更新弹窗可见性
    const [changeVisible, setChangeVisible] = useState(false); // 修改密码弹窗可见性
    const [visible, setVisible] = useState(false); // 创建弹窗可见性
    const [list, setList] = useState([]); // 用户列表
    const [height, setHeight] = useState(window.innerHeight); // 动态表格高度

    // 用户状态管理相关状态
    const [statusModalVisible, setStatusModalVisible] = useState(false); // 状态管理弹窗可见性
    const [statusUser, setStatusUser] = useState(null); // 当前要修改状态的用户
    const [statusForm] = Form.useForm(); // 状态表单

    // 用户权限查看相关状态
    const [permissionsModalVisible, setPermissionsModalVisible] = useState(false); // 权限查看弹窗可见性
    const [permissionsUser, setPermissionsUser] = useState(null); // 当前要查看权限的用户

    // 批量操作相关状态
    const [selectedRowKeys, setSelectedRowKeys] = useState([]); // 选中的行
    const [batchStatusModalVisible, setBatchStatusModalVisible] = useState(false); // 批量状态修改弹窗
    const [batchRoleModalVisible, setBatchRoleModalVisible] = useState(false); // 批量角色分配弹窗
    const [batchForm] = Form.useForm(); // 批量操作表单
    const [roles, setRoles] = useState([]);
    const [tenants, setTenants] = useState([]);
    const [selectedTenant, setSelectedTenant] = useState('');
    const [activityLogsVisible, setActivityLogsVisible] = useState(false);
    const [activityLogsUser, setActivityLogsUser] = useState(null);

    // 表格列定义
    const columns = [
        {
            title: '用户名',
            dataIndex: 'username',
            key: 'username',
            render: (text, record) => (
                <div style={{ display: 'flex', flexDirection: 'column' }}>
                    {text}
                    <Tooltip title="点击复制 ID">
                        <span
                            style={{
                                color: '#8c8c8c',     // 灰色字体
                                fontSize: '12px',
                                cursor: 'pointer',
                                userSelect: 'none',
                                display: 'inline-block',
                                maxWidth: '200px',
                                overflow: 'hidden',
                                textOverflow: 'ellipsis',
                                whiteSpace: 'nowrap'
                            }}
                            onClick={() => copyToClipboard(record.userid)}
                        >
                            {record.userid}
                            <CopyOutlined style={{ marginLeft: 8 }} />
                        </span>
                    </Tooltip>
                </div>
            ),
        },
        {
            title: '邮箱',
            dataIndex: 'email',
            key: 'email',
            render: (text) => text || '-',
        },
        {
            title: '手机号',
            dataIndex: 'phone',
            key: 'phone',
            render: (text) => text || '-',
        },
        {
            title: '创建人',
            dataIndex: 'create_by',
            key: 'create_by',
        },
        {
            title: '创建时间',
            dataIndex: 'create_at',
            key: 'create_at',
            render: (text) => new Date(text * 1000).toLocaleString(),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: '100px',
            render: (text) => {
                const statusConfig = {
                    enabled: { color: 'green', text: '启用', icon: <CheckCircleOutlined /> },
                    disabled: { color: 'red', text: '禁用', icon: <StopOutlined /> },
                    locked: { color: 'gold', text: '锁定', icon: <LockOutlined /> }
                }
                const config = statusConfig[text] || statusConfig.enabled
                return (
                    <Tag
                        color={config.color}
                        style={{ fontSize: '12px', fontWeight: '500' }}
                        icon={config.icon}
                    >
                        {config.text}
                    </Tag>
                )
            },
        },
        {
            title: '操作',
            dataIndex: 'operation',
            fixed: 'right',
            width: 320,
            render: (_, record) => (
                list.length >= 1 && (
                    <div>
                        <Button
                            type="link"
                            onClick={() => openChangePassModal(record)}
                            disabled={record.create_by === 'LDAP'}
                        >
                            重置密码
                        </Button>
                        <Space size="middle">
                            <Tooltip title="查看权限">
                                <Button
                                    type="text"
                                    icon={<AppstoreOutlined />}
                                    onClick={() => openPermissionsModal(record)}
                                    style={{ color: "#722ed1" }}
                                />
                            </Tooltip>
                            <Tooltip title="活动日志">
                                <Button
                                    type="text"
                                    icon={<HistoryOutlined />}
                                    onClick={() => openActivityLogsModal(record)}
                                    style={{ color: "#fa8c16" }}
                                />
                            </Tooltip>
                            <Tooltip title={getStatusTooltip(record.status)}>
                                <Button
                                    type="text"
                                    icon={getStatusIcon(record.status)}
                                    onClick={() => openStatusModal(record)}
                                    disabled={record.userid === 'admin'}
                                    style={{ color: getStatusColor(record.status) }}
                                />
                            </Tooltip>
                            <Tooltip title="更新">
                                <Button
                                    type="text"
                                    icon={<EditOutlined />}
                                    onClick={() => handleUpdateModalOpen(record)}
                                    style={{ color: "#1677ff" }}
                                />
                            </Tooltip>
                            <Tooltip title="删除">
                                <Popconfirm
                                    title="确定要删除此用户吗?"
                                    onConfirm={() => handleDelete(record)}
                                    okText="确定"
                                    cancelText="取消"
                                    placement="left"
                                >
                                    <Button type="text" icon={<DeleteOutlined />} style={{ color: "#ff4d4f" }} />
                                </Popconfirm>
                            </Tooltip>
                        </Space>
                    </div>
                )
            ),
        },
    ];

    // 动态调整表格高度
    useEffect(() => {
        const handleResize = () => setHeight(window.innerHeight);
        window.addEventListener('resize', handleResize);
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    // 加载用户列表
    const handleList = async (skipCache = false) => {
        try {
            const params = {}
            if (selectedTenant) {
                params.tenantId = selectedTenant
            }
            const res = await getUserList(params, { skipCache });
            setList(res.data);
        } catch (error) {
            console.error(error);
        }
    };

    // 删除用户
    const handleDelete = async (record) => {
        try {
            await deleteUser({ userid: record.userid });
            clearCacheByUrl('/api/w8t/user')
            clearCacheByUrl('/api/w8t/user/userList')
            await handleList(true);
        } catch (error) {
            console.error(error);
        }
    };

    // 打开更新用户弹窗
    const handleUpdateModalOpen = (record) => {
        setSelectedRow(record);
        setUpdateVisible(true);
    };

    // 打开重置密码弹窗
    const openChangePassModal = (record) => {
        setSelectedRow(record); // 动态绑定当前选中用户
        setChangeVisible(true);
    };

    // 打开状态管理弹窗
    const openStatusModal = (record) => {
        setStatusUser(record)
        statusForm.setFieldsValue({
            status: record.status || 'enabled',
            reason: ''
        })
        setStatusModalVisible(true)
    }

    // 打开权限查看弹窗
    const openPermissionsModal = (record) => {
        setPermissionsUser(record)
        setPermissionsModalVisible(true)
    }

    const openActivityLogsModal = (record) => {
        setActivityLogsUser(record)
        setActivityLogsVisible(true)
    }

    // 批量操作处理函数
    const handleBatchStatus = async (values) => {
        try {
            await batchUserOperation({
                userIds: selectedRowKeys,
                operation: 'status',
                status: values.status,
                statusReason: values.reason
            })
            setBatchStatusModalVisible(false)
            setSelectedRowKeys([])
            clearCacheByUrl('/api/w8t/user')
            clearCacheByUrl('/api/w8t/user/list')
            await handleList(true)
        } catch (error) {
            console.error("批量修改状态失败:", error)
        }
    }

    const handleBatchRole = async (values) => {
        try {
            await batchUserOperation({
                userIds: selectedRowKeys,
                operation: 'role',
                role: values.role
            })
            setBatchRoleModalVisible(false)
            setSelectedRowKeys([])
            clearCacheByUrl('/api/w8t/user')
            clearCacheByUrl('/api/w8t/user/list')
            await handleList(true)
        } catch (error) {
            console.error("批量分配角色失败:", error)
        }
    }

    const handleBatchDelete = async () => {
        try {
            await batchUserOperation({
                userIds: selectedRowKeys,
                operation: 'delete'
            })
            setSelectedRowKeys([])
            clearCacheByUrl('/api/w8t/user')
            clearCacheByUrl('/api/w8t/user/list')
            await handleList(true)
        } catch (error) {
            console.error("批量删除失败:", error)
        }
    }

    useEffect(() => {
        const fetchRoles = async () => {
            try {
                const { getRoleList } = await import('../../../api/role');
                const res = await getRoleList();
                setRoles(res.data || []);
            } catch (error) {
                console.error("获取角色列表失败:", error);
            }
        };

        const fetchTenants = async () => {
            try {
                const { getTenantList } = await import('../../../api/tenant');
                const res = await getTenantList({});
                setTenants(res.data || []);
            } catch (error) {
                console.error("获取租户列表失败:", error);
            }
        };

        fetchRoles();
        fetchTenants();
    }, []);

    const handleTenantChange = async (value) => {
        setSelectedTenant(value)
        try {
            const params = {}
            if (value) {
                params.tenantId = value
            }
            const res = await getUserList(params, { skipCache: true })
            setList(res.data);
        } catch (error) {
            console.error(error);
        }
    }

    // 处理状态更新
    const handleStatusUpdate = async (values) => {
        try {
            await updateUserStatus({
                userid: statusUser.userid,
                status: values.status,
                statusReason: values.reason
            })
            setStatusModalVisible(false)
            clearCacheByUrl('/api/w8t/user')
            clearCacheByUrl('/api/w8t/user/list')
            await handleList(true)
            message.success(`用户 ${statusUser.username} 状态已更新为 ${values.status}`)
        } catch (error) {
            console.error("更新用户状态失败:", error)
            HandleApiError(error)
        }
    }

    // 获取状态提示信息
    const getStatusTooltip = (status) => {
        const tooltips = {
            enabled: '点击禁用用户',
            disabled: '点击启用用户',
            locked: '点击解锁用户'
        }
        return tooltips[status] || '点击设置用户状态'
    }

    // 获取状态图标
    const getStatusIcon = (status) => {
        const icons = {
            enabled: <UnlockOutlined />,
            disabled: <StopOutlined />,
            locked: <LockOutlined />
        }
        return icons[status] || <UnlockOutlined />
    }

    // 获取状态颜色
    const getStatusColor = (status) => {
        const colors = {
            enabled: '#52c41a',
            disabled: '#ff4d4f',
            locked: '#faad14'
        }
        return colors[status] || '#52c41a'
    }

    // 搜索用户
    const onSearch = async (value) => {
        try {
            const params = {
                query: value
            }
            if (selectedTenant) {
                params.tenantId = selectedTenant
            }
            const res = await getUserList(params, { skipCache: true })
            setList(res.data);
        } catch (error) {
            console.error(error);
        }
    };

    // 初始化加载用户列表
    useEffect(() => {
        handleList(true);
    }, []);

    return (
        <>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Space>
                    <Search
                        allowClear
                        placeholder="输入搜索关键字"
                        onSearch={onSearch}
                        style={{ width: 300 }}
                    />
                    <Select
                        allowClear
                        placeholder="筛选租户"
                        style={{ width: 200 }}
                        onChange={handleTenantChange}
                        options={tenants.map(t => ({ label: t.name, value: t.id }))}
                    />
                </Space>
                <Space>
                    {selectedRowKeys.length > 0 && (
                        <Dropdown
                            menu={{
                                items: [
                                    {
                                        key: 'status',
                                        label: '批量修改状态',
                                        onClick: () => {
                                            batchForm.setFieldsValue({ status: 'enabled', reason: '' })
                                            setBatchStatusModalVisible(true)
                                        }
                                    },
                                    {
                                        key: 'role',
                                        label: '批量分配角色',
                                        onClick: () => {
                                            batchForm.setFieldsValue({ role: '' })
                                            setBatchRoleModalVisible(true)
                                        }
                                    },
                                    {
                                        key: 'delete',
                                        label: '批量删除',
                                        danger: true,
                                        onClick: () => {
                                            Modal.confirm({
                                                title: '确认批量删除',
                                                content: `确定要删除选中的 ${selectedRowKeys.length} 个用户吗？`,
                                                okText: '确定',
                                                cancelText: '取消',
                                                onOk: handleBatchDelete
                                            })
                                        }
                                    }
                                ]
                            }}
                        >
                            <Button type="default" icon={<DownOutlined />}>
                                批量操作 ({selectedRowKeys.length})
                            </Button>
                        </Dropdown>
                    )}
                    <Button
                        type="primary"
                        onClick={() => setVisible(true)}
                        style={{
                            backgroundColor: '#000000'
                        }}
                        icon={<PlusOutlined />}
                    >
                        创建
                    </Button>
                </Space>
            </div>

            {/* 用户创建弹窗 */}
            <UserCreateModal
                visible={visible}
                onClose={() => setVisible(false)}
                type="create"
                handleList={handleList}
            />

            {/* 用户更新弹窗 */}
            <UserCreateModal
                visible={updateVisible}
                onClose={() => setUpdateVisible(false)}
                selectedRow={selectedRow}
                type="update"
                handleList={handleList}
            />

            {/* 重置密码弹窗 */}
            {selectedRow && (
                <UserChangePass
                    visible={changeVisible}
                    onClose={() => setChangeVisible(false)}
                    userid={selectedRow.userid}
                    username={selectedRow.username}
                />
            )}

            {/* 用户状态管理弹窗 */}
            <Modal
                title="管理用户状态"
                visible={statusModalVisible}
                onCancel={() => setStatusModalVisible(false)}
                footer={null}
                width={600}
            >
                <Form form={statusForm} onFinish={handleStatusUpdate}>
                    <div style={{ marginBottom: 16 }}>
                        <Text>用户：{statusUser?.username} ({statusUser?.userid})</Text>
                    </div>

                    <Form.Item
                        name="status"
                        label="用户状态"
                        rules={[{ required: true, message: '请选择用户状态' }]}
                    >
                        <Select
                            placeholder="请选择用户状态"
                            value={statusUser?.status || 'enabled'}
                            onChange={(value) => {
                                statusForm.setFieldsValue({ status: value })
                            }}
                        >
                            <Select.Option value="enabled" label="启用" />
                            <Select.Option value="disabled" label="禁用" />
                            <Select.Option value="locked" label="锁定" />
                        </Select>
                    </Form.Item>

                    <Form.Item
                        name="reason"
                        label="变更原因"
                        rules={[{ required: true, message: '请输入变更原因' }]}
                    >
                        <Input.TextArea
                            placeholder="请输入状态变更原因..."
                            rows={4}
                        />
                    </Form.Item>

                    <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
                        <Space>
                            <Button onClick={() => setStatusModalVisible(false)}>
                                取消
                            </Button>
                            <Button type="primary" htmlType="submit">
                                确认修改
                            </Button>
                        </Space>
                    </div>
                </Form>
            </Modal>

            <UserPermissionsModal
                visible={permissionsModalVisible}
                onClose={() => setPermissionsModalVisible(false)}
                userid={permissionsUser?.userid}
                username={permissionsUser?.username}
            />

            <Modal
                title={`用户活动日志 - ${activityLogsUser?.username}`}
                visible={activityLogsVisible}
                onCancel={() => setActivityLogsVisible(false)}
                footer={null}
                width={1000}
            >
                <UserActivityLogs
                    userId={activityLogsUser?.userid}
                    username={activityLogsUser?.username}
                />
            </Modal>

            <Modal
                title={`批量修改状态 (${selectedRowKeys.length} 个用户)`}
                visible={batchStatusModalVisible}
                onCancel={() => setBatchStatusModalVisible(false)}
                footer={null}
                width={600}
            >
                <Form form={batchForm} onFinish={handleBatchStatus}>
                    <div style={{ marginBottom: 16 }}>
                        <Text>已选择 {selectedRowKeys.length} 个用户</Text>
                    </div>

                    <Form.Item
                        name="status"
                        label="用户状态"
                        rules={[{ required: true, message: '请选择用户状态' }]}
                    >
                        <Select placeholder="请选择用户状态">
                            <Select.Option value="enabled" label="启用" />
                            <Select.Option value="disabled" label="禁用" />
                            <Select.Option value="locked" label="锁定" />
                        </Select>
                    </Form.Item>

                    <Form.Item
                        name="reason"
                        label="变更原因"
                        rules={[{ required: true, message: '请输入变更原因' }]}
                    >
                        <Input.TextArea
                            placeholder="请输入状态变更原因..."
                            rows={4}
                        />
                    </Form.Item>

                    <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
                        <Space>
                            <Button onClick={() => setBatchStatusModalVisible(false)}>
                                取消
                            </Button>
                            <Button type="primary" htmlType="submit">
                                确认修改
                            </Button>
                        </Space>
                    </div>
                </Form>
            </Modal>

            <Modal
                title={`批量分配角色 (${selectedRowKeys.length} 个用户)`}
                visible={batchRoleModalVisible}
                onCancel={() => setBatchRoleModalVisible(false)}
                footer={null}
                width={600}
            >
                <Form form={batchForm} onFinish={handleBatchRole}>
                    <div style={{ marginBottom: 16 }}>
                        <Text>已选择 {selectedRowKeys.length} 个用户</Text>
                    </div>

                    <Form.Item
                        name="role"
                        label="角色"
                        rules={[{ required: true, message: '请选择角色' }]}
                    >
                        <Select placeholder="请选择角色">
                            {roles.map(role => (
                                <Select.Option key={role.roleid} value={role.roleid}>
                                    {role.rolename}
                                </Select.Option>
                            ))}
                        </Select>
                    </Form.Item>

                    <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
                        <Space>
                            <Button onClick={() => setBatchRoleModalVisible(false)}>
                                取消
                            </Button>
                            <Button type="primary" htmlType="submit">
                                确认分配
                            </Button>
                        </Space>
                    </div>
                </Form>
            </Modal>

            {/* 用户表格 */}
            <div style={{ overflowX: 'auto', marginTop: 10 }}>
                <Table
                    columns={columns}
                    dataSource={list}
                    rowKey="userid"
                    rowSelection={{
                        selectedRowKeys,
                        onChange: (newSelectedRowKeys) => {
                            setSelectedRowKeys(newSelectedRowKeys)
                        },
                        getCheckboxProps: (record) => ({
                            disabled: record.userid === 'admin',
                        }),
                    }}
                    scroll={{
                        y: height - 280, // 动态设置滚动高度
                        x: 'max-content', // 水平滚动
                    }}
                    style={{
                        backgroundColor: "#fff",
                        borderRadius: "8px",
                        overflow: "hidden",
                    }}
                    pagination={{
                        showTotal: HandleShowTotal,
                        pageSizeOptions: ['10'],
                    }}
                />
            </div>
        </>
    );
};
