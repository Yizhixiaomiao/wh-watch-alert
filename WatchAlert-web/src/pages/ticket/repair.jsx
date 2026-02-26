"use client"

import { useState } from "react"
import {
    Form,
    Input,
    Select,
    Button,
    Card,
    Row,
    Col,
    message,
    Upload,
    Space,
    Divider,
} from "antd"
import {
    UploadOutlined,
    SendOutlined,
} from "@ant-design/icons"
import { useNavigate } from "react-router-dom"
import { createTicket } from "../../api/ticket"
import { getUserList } from "../../api/user"

const { TextArea } = Input
const { Option } = Select

export const RepairForm = () => {
    const navigate = useNavigate()
    const [form] = Form.useForm()
    const [loading, setLoading] = useState(false)
    const [fileList, setFileList] = useState([])
    const [userList, setUserList] = useState([])

    const loadUserList = async () => {
        try {
            const res = await getUserList({ page: 1, size: 100 })
            if (res && res.list) {
                setUserList(res.list)
            }
        } catch (error) {
            console.error("加载用户列表失败:", error)
        }
    }

    const handleSubmit = async (values) => {
        setLoading(true)
        try {
            const params = {
                ...values,
                type: "Fault",
                source: "Manual",
                priority: values.priority || "P2",
                severity: values.severity || "Medium",
            }

            await createTicket(params)
            message.success("报修工单提交成功")
            navigate("/ticket")
        } catch (error) {
            console.error("提交失败:", error)
        } finally {
            setLoading(false)
        }
    }

    const uploadProps = {
        name: "file",
        multiple: true,
        fileList,
        onChange: (info) => {
            setFileList(info.fileList)
        },
        beforeUpload: (file) => {
            return false
        },
    }

    return (
        <div style={{ padding: "24px" }}>
            <Card title="人工报修" bordered={false}>
                <Form
                    form={form}
                    layout="vertical"
                    onFinish={handleSubmit}
                    initialValues={{
                        type: "Fault",
                        priority: "P2",
                        severity: "Medium",
                    }}
                >
                    <Row gutter={24}>
                        <Col span={12}>
                            <Form.Item
                                label="报修类型"
                                name="type"
                                rules={[{ required: true, message: "请选择报修类型" }]}
                            >
                                <Select placeholder="请选择报修类型">
                                    <Option value="Fault">故障报修</Option>
                                    <Option value="Query">咨询工单</Option>
                                    <Option value="Change">变更申请</Option>
                                </Select>
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item
                                label="优先级"
                                name="priority"
                                rules={[{ required: true, message: "请选择优先级" }]}
                            >
                                <Select placeholder="请选择优先级">
                                    <Option value="P0">P0 - 最高</Option>
                                    <Option value="P1">P1 - 高</Option>
                                    <Option value="P2">P2 - 中</Option>
                                    <Option value="P3">P3 - 低</Option>
                                    <Option value="P4">P4 - 最低</Option>
                                </Select>
                            </Form.Item>
                        </Col>
                    </Row>

                    <Form.Item
                        label="报修标题"
                        name="title"
                        rules={[{ required: true, message: "请输入报修标题" }]}
                    >
                        <Input placeholder="请简要描述故障" />
                    </Form.Item>

                    <Form.Item
                        label="故障描述"
                        name="description"
                        rules={[{ required: true, message: "请输入故障描述" }]}
                    >
                        <TextArea
                            rows={6}
                            placeholder="请详细描述故障现象、发生时间、影响范围等"
                        />
                    </Form.Item>

                    <Row gutter={24}>
                        <Col span={12}>
                            <Form.Item
                                label="影响范围"
                                name="faultCenterId"
                            >
                                <Select placeholder="请选择影响范围">
                                    <Option value="default">默认故障中心</Option>
                                    <Option value="server">服务器</Option>
                                    <Option value="network">网络</Option>
                                    <Option value="application">应用</Option>
                                    <Option value="database">数据库</Option>
                                </Select>
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item
                                label="指派给"
                                name="assignedTo"
                            >
                                <Select
                                    placeholder="请选择处理人"
                                    showSearch
                                    filterOption={(input, option) =>
                                        option.children.toLowerCase().indexOf(input.toLowerCase()) >= 0
                                    }
                                    onFocus={loadUserList}
                                >
                                    {userList.map((user) => (
                                        <Option key={user.id} value={user.id}>
                                            {user.username}
                                        </Option>
                                    ))}
                                </Select>
                            </Form.Item>
                        </Col>
                    </Row>

                    <Form.Item label="附件上传">
                        <Upload {...uploadProps}>
                            <Button icon={<UploadOutlined />}>上传附件</Button>
                        </Upload>
                    </Form.Item>

                    <Divider />

                    <Form.Item>
                        <Space>
                            <Button
                                type="primary"
                                htmlType="submit"
                                loading={loading}
                                icon={<SendOutlined />}
                            >
                                提交报修
                            </Button>
                            <Button onClick={() => navigate("/ticket")}>
                                取消
                            </Button>
                        </Space>
                    </Form.Item>
                </Form>
            </Card>
        </div>
    )
}

export default RepairForm