import {Modal, Descriptions, Tag, Empty, Spin} from 'antd';
import React, {useState, useEffect} from 'react';
import {getUserPermissions} from '../../../api/user';
import {AppstoreOutlined} from '@ant-design/icons';

const moduleNames = {
    'rule': '告警规则',
    'ruleGroup': '告警规则组',
    'ruleTmpl': '规则模板',
    'ruleTmplGroup': '规则模板组',
    'silence': '静默规则',
    'dutyManage': '值班管理',
    'calendar': '值班表',
    'tenant': '租户管理',
    'user': '用户管理',
    'role': '角色管理',
    'datasource': '数据源',
    'notice': '通知对象',
    'noticeTemplate': '通知模板',
    'subscribe': '告警订阅',
    'event': '告警事件',
    'dashboard': '仪表盘',
    'probing': '网络拨测',
    'knowledge': '知识库',
    'knowledgeCategory': '知识库分类',
    'ticket': '工单管理',
    'ticketTemplate': '工单模板',
    'ticketSLAPolicy': 'SLA策略',
    'ticketStep': '工单步骤',
    'ticketReview': '工单评审',
    'workHours': '工时标准',
    'assignmentRule': '智能派单',
    'faultCenter': '故障中心',
    'setting': '系统设置',
    '其他': '其他'
};

const UserPermissionsModal = ({visible, onClose, userid, username}) => {
    const [permissions, setPermissions] = useState({});
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        if (visible && userid) {
            fetchPermissions();
        }
    }, [visible, userid]);

    const fetchPermissions = async () => {
        setLoading(true);
        try {
            const res = await getUserPermissions({userid});
            setPermissions(res.data || {});
        } catch (error) {
            console.error('获取用户权限失败:', error);
        } finally {
            setLoading(false);
        }
    };

    const totalPermissions = Object.values(permissions).reduce((sum, perms) => sum + perms.length, 0);

    return (
        <Modal
            title={
                <div style={{display: 'flex', alignItems: 'center', gap: '8px'}}>
                    <AppstoreOutlined />
                    <span>用户权限 - {username}</span>
                </div>
            }
            visible={visible}
            onCancel={onClose}
            footer={null}
            width={800}
        >
            <Spin spinning={loading}>
                <Descriptions bordered column={1} size="small">
                    <Descriptions.Item label="用户ID">
                        <span style={{fontFamily: 'monospace', color: '#1890ff'}}>{userid}</span>
                    </Descriptions.Item>
                    <Descriptions.Item label="权限总数">
                        <Tag color="blue" style={{margin: 0}}>
                            {totalPermissions} 个权限
                        </Tag>
                    </Descriptions.Item>
                </Descriptions>

                {Object.keys(permissions).length === 0 ? (
                    <div style={{padding: '40px 0', textAlign: 'center'}}>
                        <Empty description="暂无权限数据" />
                    </div>
                ) : (
                    <div style={{marginTop: '16px'}}>
                        {Object.entries(permissions).map(([module, perms]) => (
                            <div key={module} style={{marginBottom: '16px'}}>
                                <div style={{
                                    marginBottom: '8px',
                                    fontWeight: '600',
                                    color: '#262626',
                                    display: 'flex',
                                    alignItems: 'center',
                                    gap: '8px'
                                }}>
                                    <Tag color="geekblue" style={{margin: 0}}>
                                        {moduleNames[module] || module}
                                    </Tag>
                                    <span style={{fontSize: '12px', color: '#8c8c8c'}}>
                                        ({perms.length} 个)
                                    </span>
                                </div>
                                <div style={{
                                    display: 'grid',
                                    gridTemplateColumns: 'repeat(auto-fill, minmax(240px, 1fr))',
                                    gap: '8px'
                                }}>
                                    {perms.map((perm, index) => (
                                        <Tag
                                            key={index}
                                            style={{
                                                margin: 0,
                                                fontSize: '12px',
                                                lineHeight: '22px',
                                                borderColor: '#d9d9d9'
                                            }}
                                        >
                                            {perm.name}
                                        </Tag>
                                    ))}
                                </div>
                            </div>
                        ))}
                    </div>
                )}
            </Spin>
        </Modal>
    );
};

export default UserPermissionsModal;