/* eslint-disable react-hooks/exhaustive-deps */
import { Modal, Form, Input, Button, Tree, Input as SearchInput, Tag, Space, Divider, Tooltip } from 'antd'
import { SearchOutlined, FolderOutlined, FileOutlined, CheckOutlined, CloseOutlined, AppstoreOutlined } from '@ant-design/icons'
import React, { useEffect, useState, useMemo } from 'react'
import { createRole, updateRole } from '../../../api/role'
import { getPermissionsList } from '../../../api/permissions'
import { clearCacheByUrl } from '../../../utils/http'

const MyFormItemContext = React.createContext([])

function toArr(str) {
    return Array.isArray(str) ? str : [str]
}

const MyFormItem = ({ name, ...props }) => {
    const prefixPath = React.useContext(MyFormItemContext)
    const concatName = name !== undefined ? [...prefixPath, ...toArr(name)] : undefined
    return <Form.Item name={concatName} {...props} />
}

const UserRoleCreateModal = ({ visible, onClose, selectedRow, type, handleList }) => {
    const [form] = Form.useForm()
    const [permissionsData, setPermissionsData] = useState([])
    const [treeData, setTreeData] = useState([])
    const [checkedKeys, setCheckedKeys] = useState([])
    const [halfCheckedKeys, setHalfCheckedKeys] = useState([])
    const [expandedKeys, setExpandedKeys] = useState([])
    const [searchValue, setSearchValue] = useState('')
    const [disabledPermission, setDisabledPermission] = useState(false)
    const [spaceValue, setSpaceValue] = useState('')

    const getAllLeafKeys = useMemo(() => {
        const keys = []
        const collectLeafKeys = (nodes) => {
            nodes.forEach(node => {
                if (node.isLeaf) {
                    keys.push(node.key)
                } else if (node.children) {
                    collectLeafKeys(node.children)
                }
            })
        }
        collectLeafKeys(treeData)
        return keys
    }, [treeData])

    const selectedCount = useMemo(() => {
        return checkedKeys.filter(key => !key.startsWith('category-') && !key.startsWith('sub-')).length
    }, [checkedKeys])

    const totalCount = useMemo(() => {
        return getAllLeafKeys.length
    }, [getAllLeafKeys])

    const handleInputChange = (e) => {
        const newValue = e.target.value.replace(/\s/g, '')
        setSpaceValue(newValue)
    }

    const handleKeyPress = (e) => {
        if (e.key === ' ') {
            e.preventDefault()
        }
    }

    useEffect(() => {
        if (selectedRow && treeData.length > 0) {
            const selectedKeys = selectedRow.permissions.map(p => {
                if (typeof p === 'string') return p
                if (typeof p === 'object' && p.key) return p.key
                return null
            }).filter(k => k !== null)

            const permissionKeyToTitle = new Map()
            const collectKeyMap = (nodes) => {
                nodes.forEach(node => {
                    if (node.isLeaf) {
                        permissionKeyToTitle.set(node.permissionKey, node.key)
                    }
                    if (node.children) {
                        collectKeyMap(node.children)
                    }
                })
            }
            collectKeyMap(treeData)

            const mappedKeys = selectedKeys.map(k => permissionKeyToTitle.get(k) || k)

            setCheckedKeys(mappedKeys)
            form.setFieldsValue({
                id: selectedRow.id,
                name: selectedRow.name,
                description: selectedRow.description,
                permissions: mappedKeys,
            })
            if (selectedRow.name === 'admin') {
                setDisabledPermission(true)
            } else {
                setDisabledPermission(false)
            }
        }
    }, [selectedRow, form, treeData])

    const handleFormSubmit = async (values) => {
        const titleToPermissionInfo = new Map()
        const collectTitleMap = (nodes) => {
            nodes.forEach(node => {
                if (node.isLeaf) {
                    titleToPermissionInfo.set(node.key, {
                        key: node.permissionKey,
                        api: node.api
                    })
                }
                if (node.children) {
                    collectTitleMap(node.children)
                }
            })
        }
        collectTitleMap(treeData)

        const leafKeys = checkedKeys.filter(key => !key.startsWith('category-') && !key.startsWith('sub-'))
        const mappedPermissions = leafKeys.map(key => titleToPermissionInfo.get(key))

        if (type === 'create') {
            try {
                const params = {
                    ...values,
                    permissions: mappedPermissions
                }
                await createRole(params)
                clearCacheByUrl('/api/w8t/role')
                clearCacheByUrl('/api/w8t/role/list')
                await handleList(true)
            } catch (error) {
                console.error(error)
            }
        }

        if (type === 'update') {
            const params = {
                ...values,
                id: selectedRow.id,
                permissions: mappedPermissions
            }
            await updateRole(params)
            clearCacheByUrl('/api/w8t/role')
            clearCacheByUrl('/api/w8t/role/list')
            await handleList(true)
        }

        onClose()
    }

    const fetchData = async () => {
        try {
            const response = await getPermissionsList()
            const data = response.data
            setPermissionsData(data)
            buildTreeData(data)
        } catch (error) {
            console.error(error)
        }
    }

    useEffect(() => {
        fetchData()
    }, [])

    const buildTreeData = (data) => {
        const categoryOrder = [
            '告警管理',
            '故障中心',
            '工单管理',
            '知识库',
            '通知管理',
            '值班管理',
            '拨测规则',
            '告警订阅',
            '监控资源',
            '系统管理'
        ]

        const subCategoryOrder = {
            '告警管理': ['告警规则', '规则组', '规则模板', '模板组', '静默规则', '告警事件'],
            '故障中心': ['故障中心'],
            '工单管理': ['工单操作', '工单模板', 'SLA策略', '工单步骤', '工单评审', '工时标准', '智能派单'],
            '知识库': ['知识库', '知识库分类'],
            '通知管理': ['通知对象', '通知模板', '通知记录'],
            '值班管理': ['值班管理', '值班表'],
            '拨测规则': ['拨测规则'],
            '告警订阅': ['告警订阅'],
            '监控资源': ['数据源', '仪表盘'],
            '系统管理': ['用户管理', '角色管理', '租户管理', '系统配置']
        }

        const categoryMap = {}

        data.forEach(item => {
            const { category, subCategory, key, api, title, permissionKey } = item

            if (!categoryMap[category]) {
                categoryMap[category] = {
                    title: category,
                    key: `category-${category}`,
                    children: []
                }
            }

            const categoryNode = categoryMap[category]

            if (!categoryNode.subCategoryMap) {
                categoryNode.subCategoryMap = {}
            }

            if (!categoryNode.subCategoryMap[subCategory]) {
                const subCategoryNode = {
                    title: subCategory,
                    key: `sub-${category}-${subCategory}`,
                    children: []
                }
                categoryNode.subCategoryMap[subCategory] = subCategoryNode
                categoryNode.children.push(subCategoryNode)
            }

            const leafNode = {
                title: title || key,
                key: key,
                permissionKey: permissionKey,
                api: api,
                isLeaf: true
            }
            categoryNode.subCategoryMap[subCategory].children.push(leafNode)
        })

        const treeData = []

        categoryOrder.forEach(category => {
            if (categoryMap[category]) {
                const categoryNode = categoryMap[category]
                const subCategories = subCategoryOrder[category] || []

                const sortedChildren = []
                subCategories.forEach(subCategory => {
                    if (categoryNode.subCategoryMap[subCategory]) {
                        sortedChildren.push(categoryNode.subCategoryMap[subCategory])
                    }
                })

                Object.values(categoryNode.subCategoryMap).forEach(subNode => {
                    if (!subCategories.includes(subNode.title)) {
                        sortedChildren.push(subNode)
                    }
                })

                treeData.push({
                    title: categoryNode.title,
                    key: categoryNode.key,
                    children: sortedChildren
                })
            }
        })

        Object.keys(categoryMap).forEach(category => {
            if (!categoryOrder.includes(category)) {
                const categoryNode = categoryMap[category]
                const { subCategoryMap, ...rest } = categoryNode
                treeData.push(rest)
            }
        })

        setTreeData(treeData)

        const allKeys = []
        const collectKeys = (nodes) => {
            nodes.forEach(node => {
                allKeys.push(node.key)
                if (node.children) {
                    collectKeys(node.children)
                }
            })
        }
        collectKeys(treeData)
        setExpandedKeys(allKeys)
    }

    const handleCheck = (checkedKeysValue, info) => {
        const keys = checkedKeysValue.checked || checkedKeysValue
        const leafKeys = Array.isArray(keys) ? keys.filter(key => !key.startsWith('category-') && !key.startsWith('sub-')) : []
        setCheckedKeys(leafKeys)
        if (info.halfCheckedKeys) {
            setHalfCheckedKeys(info.halfCheckedKeys)
        }
    }

    const handleExpand = (expandedKeysValue) => {
        setExpandedKeys(expandedKeysValue)
    }

    const handleSelectAll = () => {
        setCheckedKeys(getAllLeafKeys)
    }

    const handleDeselectAll = () => {
        setCheckedKeys([])
    }

    const handleSearch = (value) => {
        setSearchValue(value)
        if (!value) {
            const allKeys = []
            const collectKeys = (nodes) => {
                nodes.forEach(node => {
                    allKeys.push(node.key)
                    if (node.children) {
                        collectKeys(node.children)
                    }
                })
            }
            collectKeys(treeData)
            setExpandedKeys(allKeys)
            return
        }

        const newExpandedKeys = []
        const searchInChildren = (nodes, parentKey) => {
            nodes.forEach(node => {
                const title = node.title?.toString().toLowerCase() || ''
                const matches = title.includes(value.toLowerCase())
                const childrenMatch = node.children ? searchInChildren(node.children, node.key) : false

                if (matches || childrenMatch) {
                    if (parentKey && !newExpandedKeys.includes(parentKey)) {
                        newExpandedKeys.push(parentKey)
                    }
                    if (!newExpandedKeys.includes(node.key)) {
                        newExpandedKeys.push(node.key)
                    }
                }
            })
            return false
        }

        searchInChildren(treeData, null)
        setExpandedKeys(newExpandedKeys)
    }

    const filterTreeNode = (node) => {
        if (!searchValue) return true
        const title = node.title?.toString().toLowerCase() || ''
        return title.includes(searchValue.toLowerCase())
    }

    const renderTreeNode = (node) => {
        if (node.isLeaf) {
            return (
                <span style={{ fontSize: '14px', color: '#595959' }}>
                    {node.title}
                </span>
            )
        }
        return (
            <span style={{ fontSize: '15px', fontWeight: 500, color: '#262626' }}>
                {node.title}
            </span>
        )
    }

    return (
        <Modal
            visible={visible}
            onCancel={onClose}
            footer={null}
            width={900}
            title={type === 'create' ? '创建角色' : '编辑角色'}
            destroyOnClose
        >
            <Form form={form} name="form_item_path" layout="vertical" onFinish={handleFormSubmit}>

                <MyFormItem name="name" label="角色名称"
                    rules={[
                        {
                            required: true,
                            message: '请输入角色名称!',
                        },
                    ]}>
                    <Input
                        value={spaceValue}
                        onChange={handleInputChange}
                        onKeyPress={handleKeyPress}
                        disabled={type === 'update'}
                        placeholder="请输入角色名称"
                    />
                </MyFormItem>

                <MyFormItem name="description" label="描述">
                    <Input placeholder="请输入角色描述" />
                </MyFormItem>

                <MyFormItem
                    label={
                        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
                            <span style={{ fontSize: '14px', fontWeight: 500 }}>选择权限</span>
                            <Space>
                                {selectedCount > 0 && (
                                    <Tag color="blue" style={{ margin: 0, borderRadius: '4px', padding: '2px 8px', fontSize: '13px' }}>
                                        已选 {selectedCount} / {totalCount} 项
                                    </Tag>
                                )}
                                {!disabledPermission && (
                                    <>
                                        <Tooltip title="全选所有权限">
                                            <Button
                                                type="text"
                                                size="small"
                                                icon={<CheckOutlined />}
                                                onClick={handleSelectAll}
                                                disabled={selectedCount === totalCount}
                                                style={{ color: '#1677ff' }}
                                            >
                                                全选
                                            </Button>
                                        </Tooltip>
                                        <Tooltip title="取消全选">
                                            <Button
                                                type="text"
                                                size="small"
                                                icon={<CloseOutlined />}
                                                onClick={handleDeselectAll}
                                                disabled={selectedCount === 0}
                                                style={{ color: '#ff4d4f' }}
                                            >
                                                取消
                                            </Button>
                                        </Tooltip>
                                    </>
                                )}
                            </Space>
                        </div>
                    }
                >
                    <div style={{ marginBottom: 12 }}>
                        <SearchInput
                            placeholder="搜索权限..."
                            prefix={<SearchOutlined />}
                            value={searchValue}
                            onChange={e => handleSearch(e.target.value)}
                            allowClear
                            style={{ width: '100%' }}
                        />
                    </div>
                    <div style={{
                        border: '1px solid #e8e8e8',
                        borderRadius: '8px',
                        padding: '16px',
                        maxHeight: 450,
                        overflowY: 'auto',
                        backgroundColor: '#fafafa',
                        boxShadow: '0 1px 2px rgba(0,0,0,0.03)'
                    }}>
                        <Tree
                            checkable
                            checkedKeys={checkedKeys}
                            onCheck={handleCheck}
                            expandedKeys={expandedKeys}
                            onExpand={handleExpand}
                            treeData={treeData}
                            filterTreeNode={filterTreeNode}
                            showLine={{ showLeafIcon: false }}
                            blockNode
                            disabled={disabledPermission}
                            titleRender={renderTreeNode}
                            checkStrictly={false}
                            style={{
                                background: 'transparent'
                            }}
                        />
                    </div>
                    {selectedCount === 0 && !disabledPermission && (
                        <div style={{ marginTop: 8, color: '#8c8c8c', fontSize: '13px' }}>
                            提示：请勾选上方权限列表中的权限项
                        </div>
                    )}
                </MyFormItem>

                <Divider style={{ margin: '24px 0 16px 0' }} />

                <div style={{ display: 'flex', justifyContent: 'flex-end', gap: 8 }}>
                    <Button onClick={onClose}>
                        取消
                    </Button>
                    <Button
                        type="primary"
                        htmlType="submit"
                        style={{
                            backgroundColor: '#000000'
                        }}
                    >
                        提交
                    </Button>
                </div>
            </Form>
        </Modal>
    )
}

export default UserRoleCreateModal