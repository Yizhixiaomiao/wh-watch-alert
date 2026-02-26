import { useState } from "react"
import {
    Form,
    Input,
    Button,
    Card,
    Upload,
    message,
    Result,
    Typography,
    Row,
    Col,
    Tag
} from "antd"
import { 
    CheckCircleOutlined, 
    PhoneOutlined,
    UploadOutlined
} from "@ant-design/icons"
import { mobileCreateTicket } from "../../../api/ticket"

const { TextArea } = Input
const { Title, Text } = Typography

export const MobileTicketCreate = () => {
    const [form] = Form.useForm()
    const [submitting, setSubmitting] = useState(false)
    const [submitResult, setSubmitResult] = useState(null)
    const [fileList, setFileList] = useState([])

    const handleSubmit = async () => {
        try {
            setSubmitting(true)
            
            const values = await form.validateFields()
            
            const submitData = {
                title: values.title,
                description: values.description,
                type: 'Fault',
                priority: 'P4',
                contactName: values.contactName,
                contactPhone: values.contactPhone,
                location: values.location,
                urgentLevel: '一般',
                images: fileList.map(file => file.url || file.response?.url).filter(Boolean),
                userAgent: navigator.userAgent,
                platform: 'mobile',
            }

            const res = await mobileCreateTicket(submitData)
            
            if (res && res.data) {
                setSubmitResult(res.data)
                message.success('报修提交成功！')
            }
        } catch (error) {
            message.error('提交失败，请重试')
        } finally {
            setSubmitting(false)
        }
    }

    const handleReset = () => {
        form.resetFields()
        setSubmitResult(null)
        setFileList([])
    }

    const uploadProps = {
        name: 'file',
        listType: "picture-card",
        fileList: fileList,
        onChange: (info) => {
            setFileList(info.fileList)
        },
        beforeUpload: (file) => {
            const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png'
            if (!isJpgOrPng) {
                message.error('只能上传 JPG/PNG 格式的图片!')
                return false
            }
            const isLt5M = file.size / 1024 / 1024 < 5
            if (!isLt5M) {
                message.error('图片大小不能超过 5MB!')
                return false
            }
            return false // 阻止自动上传，实际项目中需要配置上传接口
        },
    }

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
                    故障报修
                </Title>
            </div>

            {/* 报修表单 */}
            {!submitResult ? (
                <Card>
                    <Form
                        form={form}
                        layout="vertical"
                        onFinish={handleSubmit}
                    >
                        {/* 联系信息 */}
                        <Row gutter={16}>
                            <Col span={24}>
                                <Form.Item
                                    name="contactName"
                                    label="姓名"
                                    rules={[{ required: true, message: '请输入您的姓名' }]}
                                >
                                    <Input placeholder="请输入您的姓名" size="large" />
                                </Form.Item>
                            </Col>
                        </Row>

                        <Row gutter={16}>
                            <Col span={24}>
                                <Form.Item
                                    name="contactPhone"
                                    label="电话"
                                    rules={[{ required: true, message: '请输入电话号码' }]}
                                >
                                    <Input 
                                        placeholder="请输入手机号或短号（如：13800138000 或 6688）" 
                                        size="large"
                                        prefix={<PhoneOutlined />}
                                    />
                                </Form.Item>
                            </Col>
                        </Row>

                        <Row gutter={16}>
                            <Col span={24}>
                                <Form.Item
                                    name="location"
                                    label="位置"
                                    rules={[{ required: true, message: '请输入故障位置' }]}
                                >
                                    <Input 
                                        placeholder="请输入具体的故障位置" 
                                        size="large"
                                    />
                                </Form.Item>
                            </Col>
                        </Row>

                        <Row gutter={16}>
                            <Col span={24}>
                                <Form.Item
                                    name="title"
                                    label="故障描述"
                                    rules={[{ required: true, message: '请输入故障描述' }]}
                                >
                                    <TextArea 
                                        rows={4} 
                                        placeholder="请详细描述故障现象、发生时间、影响范围等信息" 
                                        size="large"
                                        style={{ fontSize: '16px' }}
                                    />
                                </Form.Item>
                            </Col>
                        </Row>

                        <Row gutter={16}>
                            <Col span={24}>
                                <Form.Item label="图片上传">
                                    <Upload {...uploadProps}>
                                        {fileList.length < 3 && (
                                            <div>
                                                <UploadOutlined />
                                                <div style={{ marginTop: 8 }}>上传图片</div>
                                            </div>
                                        )}
                                    </Upload>
                                    <div style={{ fontSize: '12px', color: '#999', marginTop: 4 }}>
                                        最多上传3张图片，每张不超过5MB，支持JPG/PNG格式
                                    </div>
                                </Form.Item>
                            </Col>
                        </Row>

                        {/* 提交按钮 */}
                        <Row gutter={16}>
                            <Col span={24}>
                                <Button 
                                    type="primary" 
                                    htmlType="submit" 
                                    loading={submitting}
                                    size="large" 
                                    block
                                    style={{ height: '50px', fontSize: '16px' }}
                                >
                                    {submitting ? '提交中...' : '提交报修'}
                                </Button>
                            </Col>
                        </Row>

                        <Row gutter={16} style={{ marginTop: '12px' }}>
                            <Col span={24}>
                                <Button 
                                    size="large" 
                                    block
                                    onClick={handleReset}
                                    style={{ height: '50px', fontSize: '16px' }}
                                >
                                    重置表单
                                </Button>
                            </Col>
                        </Row>
                    </Form>
                </Card>
            ) : (
                /* 提交成功结果 */
                <Result
                    icon={<CheckCircleOutlined style={{ color: '#52c41a', fontSize: '64px' }} />}
                    title="报修提交成功！"
                    subTitle={
                        <div style={{ textAlign: 'center' }}>
                            <div style={{ marginBottom: '16px' }}>
                                <Text strong style={{ fontSize: '18px' }}>
                                    工单号：{submitResult.ticketNo}
                                </Text>
                            </div>
                            <Text type="secondary" style={{ fontSize: '16px' }}>
                                预计{submitResult.waitTime}内处理
                            </Text>
                        </div>
                    }
                    extra={[
                        <Card key="details" style={{ textAlign: 'left', marginBottom: 16, fontSize: '16px' }}>
                            <Title level={4}>工单信息</Title>
                            <div style={{ fontSize: '16px', lineHeight: '1.8' }}>
                                <div><Text strong>工单号：</Text>{submitResult.ticketNo}</div>
                                <div><Text strong>预计处理时间：</Text>{submitResult.waitTime}</div>
                                <div><Text strong>联系电话：</Text>{form.getFieldValue('contactPhone')}</div>
                                <div><Text strong>故障位置：</Text>{form.getFieldValue('location')}</div>
                                <div style={{ marginTop: '8px' }}>
                                    <Text strong>状态：</Text>
                                    <Tag color="processing" style={{ fontSize: '14px' }}>待处理</Tag>
                                </div>
                            </div>
                        </Card>,
                        <Button 
                            type="primary" 
                            key="query" 
                            onClick={() => window.location.href = '/mobile/ticket/query?phone=' + form.getFieldValue('contactPhone')}
                            size="large"
                            block
                            style={{ height: '50px', fontSize: '16px', marginBottom: '12px' }}
                        >
                            查询处理进度
                        </Button>,
                        <Button 
                            key="new" 
                            onClick={handleReset}
                            size="large" 
                            block
                            style={{ height: '50px', fontSize: '16px' }}
                        >
                            再次报修
                        </Button>
                    ]}
                />
            )}
        </div>
    )
}