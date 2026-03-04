import http from '../utils/http';
import { message } from 'antd';
import { HandleApiError } from "../utils/lib";

// 评审操作
async function assignReviewers(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/review/assign', params);
        message.open({
            type: 'success',
            content: '评委分配成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function submitReview(params) {
    try {
        const res = await http('post', '/api/w8t/ticket/review/submit', params);
        message.open({
            type: 'success',
            content: '评审提交成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getReview(params) {
    try {
        const res = await http('get', '/api/w8t/ticket/review/get', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getReviews(params) {
    try {
        const res = await http('get', '/api/w8t/ticket/review/list', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

export {
    assignReviewers,
    submitReview,
    getReview,
    getReviews,
};