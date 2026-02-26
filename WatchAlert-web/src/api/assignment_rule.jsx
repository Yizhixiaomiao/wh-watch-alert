import http from '../utils/http';
import { message } from 'antd';
import { HandleApiError } from "../utils/lib";

// 规则操作
async function createAssignmentRule(params) {
    try {
        const res = await http('post', '/api/w8t/assignment-rule/create', params);
        message.open({
            type: 'success',
            content: '分配规则创建成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function updateAssignmentRule(params) {
    try {
        const res = await http('post', '/api/w8t/assignment-rule/update', params);
        message.open({
            type: 'success',
            content: '分配规则更新成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function deleteAssignmentRule(params) {
    try {
        const res = await http('post', '/api/w8t/assignment-rule/delete', params);
        message.open({
            type: 'success',
            content: '分配规则删除成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getAssignmentRule(params) {
    try {
        const res = await http('get', '/api/w8t/assignment-rule/get', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getAssignmentRules(params) {
    try {
        const res = await http('get', '/api/w8t/assignment-rule/list', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function matchAssignmentRule(params) {
    try {
        const res = await http('post', '/api/w8t/assignment-rule/match', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function autoAssignTicket(params) {
    try {
        const res = await http('post', '/api/w8t/assignment-rule/auto-assign', params);
        message.open({
            type: 'success',
            content: '工单自动分配成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

export {
    createAssignmentRule,
    updateAssignmentRule,
    deleteAssignmentRule,
    getAssignmentRule,
    getAssignmentRules,
    matchAssignmentRule,
    autoAssignTicket,
};