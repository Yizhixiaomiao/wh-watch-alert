import http from '../utils/http';
import { message } from 'antd';
import { HandleApiError } from "../utils/lib";

// 工时标准操作
async function createWorkHoursStandard(params) {
    try {
        const res = await http('post', '/api/w8t/work-hours/standard/create', params);
        message.open({
            type: 'success',
            content: '工时标准创建成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function updateWorkHoursStandard(params) {
    try {
        const res = await http('post', '/api/w8t/work-hours/standard/update', params);
        message.open({
            type: 'success',
            content: '工时标准更新成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function deleteWorkHoursStandard(params) {
    try {
        const res = await http('post', '/api/w8t/work-hours/standard/delete', params);
        message.open({
            type: 'success',
            content: '工时标准删除成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getWorkHoursStandard(params) {
    try {
        const res = await http('get', '/api/w8t/work-hours/standard/get', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getWorkHoursStandards(params) {
    try {
        const res = await http('get', '/api/w8t/work-hours/standard/list', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function calculateWorkHours(params) {
    try {
        const res = await http('post', '/api/w8t/work-hours/calculate', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

export {
    createWorkHoursStandard,
    updateWorkHoursStandard,
    deleteWorkHoursStandard,
    getWorkHoursStandard,
    getWorkHoursStandards,
    calculateWorkHours,
};