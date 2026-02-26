import http from '../utils/http';
import { message } from 'antd';
import {HandleApiError} from "../utils/lib";

async function getCurEventList(params, options = {}) {
    try {
        const queryString = Object.keys(params)
            .map(key => params[key] !== undefined ? `${key}=${params[key]}` : '')
            .filter(Boolean)
            .join('&');
        const res = await http('get', `/api/w8t/event/curEvent?${queryString}`, {}, options);
        return res;
    } catch (error) {
        HandleApiError(error)
        return error
    }
}

async function getHisEventList(params, options = {}) {
    try {
        const queryString = Object.keys(params)
            .map(key => params[key] !== undefined ? `${key}=${params[key]}` : '')
            .filter(Boolean)
            .join('&');

        const url = `/api/w8t/event/hisEvent?${queryString}`;
        const res = await http('get', url, {}, options);
        return res;
    } catch (error) {
        HandleApiError(error)
        return error
    }
}

async function ProcessAlertEvent(params) {
    return await http('post', `/api/w8t/event/processAlertEvent`,params);
}

async function archiveAlert(params) {
    try {
        const res = await http('post', '/api/w8t/event/archiveAlertEvent', params);
        message.success('归档成功');
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function ListEventComments(params) {
    return await http('get', '/api/w8t/event/listComments', params);
}

async function AddEventComment(params) {
    return await http('post', '/api/w8t/event/addComment', params);
}

async function DeleteEventComment(params) {
    return await http('post', '/api/w8t/event/deleteComment', params);
}

export {
    getCurEventList,
    getHisEventList,
    ProcessAlertEvent,
    archiveAlert,
    ListEventComments,
    AddEventComment,
    DeleteEventComment,
}