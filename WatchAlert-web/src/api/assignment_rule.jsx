import http from '../utils/http';
import { message } from 'antd';
import { HandleApiError } from "../utils/lib";

// 规则操作
async function createAssignmentRule(params) {
    try {
        const res = await http('post', '/api/w8t/assignment-rule/create', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function updateAssignmentRule(params) {
    try {
        const res = await http('post', '/api/w8t/assignment-rule/update', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function deleteAssignmentRule(params) {
    try {
        const res = await http('post', '/api/w8t/assignment-rule/delete', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function getAssignmentRule(params) {
    try {
        const res = await http('get', '/api/w8t/assignment-rule/get', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function getAssignmentRules(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/assignment-rule/list', params, options);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function matchAssignmentRule(params) {
    try {
        const res = await http('post', '/api/w8t/assignment-rule/match', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function autoAssignTicket(params) {
    try {
        const res = await http('post', '/api/w8t/assignment-rule/auto-assign', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
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