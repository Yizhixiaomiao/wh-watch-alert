import http from '../utils/http';
import { message } from 'antd';
import { HandleApiError } from "../utils/lib";

// 知识操作
async function createKnowledge(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/create', params);
        message.open({
            type: 'success',
            content: '知识创建成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function updateKnowledge(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/update', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function deleteKnowledge(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/delete', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getKnowledge(params) {
    try {
        const res = await http('get', '/api/w8t/knowledge/get', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getKnowledges(params) {
    try {
        const res = await http('get', '/api/w8t/knowledge/list', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function likeKnowledge(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/like', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function saveKnowledgeToTicket(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/save-to-ticket', params);
        message.open({
            type: 'success',
            content: '知识已添加到工单',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

// 分类操作
async function createKnowledgeCategory(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/category/create', params);
        message.open({
            type: 'success',
            content: '分类创建成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function updateKnowledgeCategory(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/category/update', params);
        message.open({
            type: 'success',
            content: '分类更新成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function deleteKnowledgeCategory(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/category/delete', params);
        message.open({
            type: 'success',
            content: '分类删除成功',
        });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getKnowledgeCategory(params) {
    try {
        const res = await http('get', '/api/w8t/knowledge/category/get', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

async function getKnowledgeCategories(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/knowledge/category/list', params, options);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

export {
    createKnowledge,
    updateKnowledge,
    deleteKnowledge,
    getKnowledge,
    getKnowledges,
    likeKnowledge,
    saveKnowledgeToTicket,
    createKnowledgeCategory,
    updateKnowledgeCategory,
    deleteKnowledgeCategory,
    getKnowledgeCategory,
    getKnowledgeCategories,
};