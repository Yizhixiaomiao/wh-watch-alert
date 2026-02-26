import React, { useState, useEffect } from 'react';
import { Table, Select, Button, Space, Tag, Typography } from 'antd';
import { getUserActivityLogs } from '../../../api/user';
import {HandleShowTotal} from "../../../utils/lib";

const { Text } = Typography;

export const UserActivityLogs = ({ userId, username }) => {
    const [logs, setLogs] = useState([]);
    const [loading, setLoading] = useState(false);
    const [filters, setFilters] = useState({
        action: '',
        resourceType: '',
        page: 1,
        size: 10
    });

    const actionOptions = [
        { label: '登录', value: 'login' },
        { label: '登出', value: 'logout' },
        { label: '创建', value: 'create' },
        { label: '更新', value: 'update' },
        { label: '删除', value: 'delete' },
        { label: '查看', value: 'view' },
        { label: '导出', value: 'export' }
    ];

    const resourceTypeOptions = [
        { label: '用户', value: 'user' },
        { label: '租户', value: 'tenant' },
        { label: '工单', value: 'ticket' },
        { label: '规则', value: 'rule' },
        { label: '告警', value: 'alert' }
    ];

    const columns = [
        {
            title: '时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 180,
            render: (text) => new Date(text * 1000).toLocaleString()
        },
        {
            title: '用户',
            dataIndex: 'userName',
            key: 'userName',
            width: 120
        },
        {
            title: '操作',
            dataIndex: 'action',
            key: 'action',
            width: 100,
            render: (text) => {
                const actionMap = {
                    login: { text: '登录', color: 'green' },
                    logout: { text: '登出', color: 'default' },
                    create: { text: '创建', color: 'blue' },
                    update: { text: '更新', color: 'orange' },
                    delete: { text: '删除', color: 'red' },
                    view: { text: '查看', color: 'default' },
                    export: { text: '导出', color: 'purple' }
                };
                const config = actionMap[text] || { text: text, color: 'default' };
                return <Tag color={config.color}>{config.text}</Tag>;
            }
        },
        {
            title: '资源类型',
            dataIndex: 'resourceType',
            key: 'resourceType',
            width: 100
        },
        {
            title: '资源名称',
            dataIndex: 'resourceName',
            key: 'resourceName',
            width: 150,
            render: (text) => text || '-'
        },
        {
            title: 'IP地址',
            dataIndex: 'ipAddress',
            key: 'ipAddress',
            width: 150,
            render: (text) => text || '-'
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            render: (text) => (
                <Tag color={text === 'success' ? 'green' : 'red'}>
                    {text === 'success' ? '成功' : '失败'}
                </Tag>
            )
        },
        {
            title: '错误信息',
            dataIndex: 'errorMessage',
            key: 'errorMessage',
            render: (text) => text ? <Text type="danger">{text}</Text> : '-'
        }
    ];

    const fetchLogs = async () => {
        setLoading(true);
        try {
            const params = {
                userid: userId,
                username: username,
                action: filters.action,
                resourceType: filters.resourceType,
                page: filters.page,
                size: filters.size
            };
            const res = await getUserActivityLogs(params);
            setLogs(res.data || []);
        } catch (error) {
            console.error("获取用户活动日志失败:", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        fetchLogs();
    }, [userId, username, filters]);

    const handleFilterChange = (key, value) => {
        setFilters(prev => ({
            ...prev,
            [key]: value,
            page: 1
        }));
    };

    return (
        <div>
            <Space style={{ marginBottom: 16 }}>
                <Select
                    placeholder="筛选操作类型"
                    style={{ width: 150 }}
                    allowClear
                    options={actionOptions}
                    onChange={(value) => handleFilterChange('action', value)}
                />
                <Select
                    placeholder="筛选资源类型"
                    style={{ width: 150 }}
                    allowClear
                    options={resourceTypeOptions}
                    onChange={(value) => handleFilterChange('resourceType', value)}
                />
                <Button onClick={fetchLogs}>刷新</Button>
            </Space>

            <Table
                columns={columns}
                dataSource={logs}
                loading={loading}
                rowKey="id"
                pagination={{
                    current: filters.page,
                    pageSize: filters.size,
                    total: logs.length,
                    showTotal: HandleShowTotal,
                    onChange: (page) => handleFilterChange('page', page)
                }}
                scroll={{
                    y: 400,
                    x: 'max-content'
                }}
            />
        </div>
    );
};