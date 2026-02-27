import http from '../utils/http';
import { message } from 'antd';
import { HandleApiError } from "../utils/lib";

// 知识操作
async function createKnowledge(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/create', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function updateKnowledge(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/update', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function deleteKnowledge(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/delete', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function getKnowledge(params) {
    try {
        const res = await http('get', '/api/w8t/knowledge/get', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function getKnowledges(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/knowledge/list', params, options);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function likeKnowledge(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/like', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function saveKnowledgeToTicket(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/save-to-ticket', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

// 分类操作
async function createKnowledgeCategory(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/category/create', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function updateKnowledgeCategory(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/category/update', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function deleteKnowledgeCategory(params) {
    try {
        const res = await http('post', '/api/w8t/knowledge/category/delete', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function getKnowledgeCategory(params) {
    try {
        const res = await http('get', '/api/w8t/knowledge/category/get', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
    }
}

async function getKnowledgeCategories(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/knowledge/category/list', params, options);
        return res;
    } catch (error) {
        HandleApiError(error);
        throw error;
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