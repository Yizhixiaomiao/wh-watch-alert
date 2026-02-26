import http from '../utils/http';
import { message } from 'antd';
import { HandleApiError } from "../utils/lib";

// 工单基础操作
async function getTicketList(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/ticket/list', params, options);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getTicket(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/ticket/get', params, options);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function createTicket(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/create', params);
        message.open({
            type: 'success',
            content: '工单创建成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function updateTicket(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/update', params);
        message.open({
            type: 'success',
            content: '工单更新成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function deleteTicket(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/delete', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

// 工单状态操作
async function assignTicket(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/assign', params);
        message.open({
            type: 'success',
            content: '工单分配成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function claimTicket(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/claim', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function transferTicket(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/transfer', params);
        message.open({
            type: 'success',
            content: '工单转派成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function escalateTicket(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/escalate', params);
        message.open({
            type: 'success',
            content: '工单升级成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function resolveTicket(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/resolve', params);
        message.open({
            type: 'success',
            content: '工单已标记为解决',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function closeTicket(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/close', params);
        message.open({
            type: 'success',
            content: '工单已关闭',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function reopenTicket(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/reopen', params);
        message.open({
            type: 'success',
            content: '工单已重新打开',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

// 评论和日志
async function addTicketComment(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/comment', params);
        message.open({
            type: 'success',
            content: '评论添加成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getTicketComments(params) {
    try {
        const res = await http('get', '/api/w8t/ticket/comments', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getTicketWorkLogs(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/ticket/worklog', params, options);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

// 统计
async function getTicketStatistics(params) {
    try {
        const res = await http('get', '/api/w8t/ticket/statistics', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

// 工单模板
async function getTicketTemplateList(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/ticket/template/list', params, options);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function createTicketTemplate(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/template/create', params);
        message.open({
            type: 'success',
            content: '模板创建成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function updateTicketTemplate(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/template/update', params);
        message.open({
            type: 'success',
            content: '模板更新成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function deleteTicketTemplate(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/template/delete', params);
        message.open({
            type: 'success',
            content: '模板删除成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

// SLA策略
async function getTicketSLAPolicyList(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/ticket/sla/list', params, options);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function createTicketSLAPolicy(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/sla/create', params);
        message.open({
            type: 'success',
            content: 'SLA策略创建成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function updateTicketSLAPolicy(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/sla/update', params);
        message.open({
            type: 'success',
            content: 'SLA策略更新成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function deleteTicketSLAPolicy(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/sla/delete', params);
        message.open({
            type: 'success',
            content: 'SLA策略删除成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

// 移动端专用接口
async function mobileCreateTicket(params) {
    try {
        const res = await http('post', '/api/w8t/mobile/ticket/create', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function mobileQueryTicket(params) {
    try {
        const res = await http('get', '/api/w8t/mobile/ticket/query', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

// 处理步骤操作
async function addTicketStep(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/step/add', params);
        message.open({
            type: 'success',
            content: '步骤添加成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function updateTicketStep(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/step/update', params);
        message.open({
            type: 'success',
            content: '步骤更新成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function deleteTicketStep(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/step/delete', params);
        message.open({
            type: 'success',
            content: '步骤删除成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getTicketSteps(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/ticket/steps', params, options);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function reorderTicketSteps(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/step/reorder', params);
        message.open({
            type: 'success',
            content: '步骤排序成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

export {
    // 工单基础操作
    getTicketList,
    getTicket,
    createTicket,
    updateTicket,
    deleteTicket,
    // 工单状态操作
    assignTicket,
    claimTicket,
    transferTicket,
    escalateTicket,
    resolveTicket,
    closeTicket,
    reopenTicket,
    // 处理步骤操作
    addTicketStep,
    updateTicketStep,
    deleteTicketStep,
    getTicketSteps,
    reorderTicketSteps,
    // 评论和日志
    addTicketComment,
    getTicketComments,
    getTicketWorkLogs,
    // 统计
    getTicketStatistics,
    // 工单模板
    getTicketTemplateList,
    createTicketTemplate,
    updateTicketTemplate,
    deleteTicketTemplate,
    // SLA策略
    getTicketSLAPolicyList,
    createTicketSLAPolicy,
    updateTicketSLAPolicy,
    deleteTicketSLAPolicy,
    // 移动端接口
    mobileCreateTicket,
    mobileQueryTicket,
};
