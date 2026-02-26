// 微信相关API和服务

import http from '../utils/http';
import { HandleApiError } from "../utils/lib";

/**
 * 获取微信JS-SDK配置
 */
export async function getWechatConfig(url) {
    try {
        const res = await http('get', '/api/w8t/wechat/config', { url });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

/**
 * 获取微信用户信息（如果需要）
 */
export async function getWechatUserInfo(code) {
    try {
        const res = await http('get', '/api/w8t/wechat/user', { code });
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

/**
 * 发送微信模板消息
 */
export async function sendWechatTemplateMessage(params) {
    try {
        const res = await http('post', '/api/w8t/wechat/template-message', params);
        return res;
    } catch (error) {
        HandleApiError(error);
        return error;
    }
}

export {
    getWechatConfig,
    getWechatUserInfo,
    sendWechatTemplateMessage
};