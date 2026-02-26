/**
 * 网络请求配置
 */
import axios from 'axios';
import { message } from "antd";

const protocol = window.location.protocol;
const curUrl = window.location.hostname
const port = window.location.port;
axios.defaults.timeout = 100000;
axios.defaults.baseURL = `${protocol}//${curUrl}:${port}`;

// 请求缓存
const cache = new Map();
const CACHE_DURATION = 5 * 60 * 1000; // 5分钟缓存

// 请求队列（用于并发控制）
const pendingRequests = new Map();

/**
 * http request 拦截器
 */
axios.interceptors.request.use(
    (config) => {
        config.headers = {
            'Content-Type': 'application/json',
            'TenantID': localStorage.getItem('TenantID'),
        };
        if (localStorage.getItem('Authorization')) {
            config.headers.Authorization = `Bearer ${localStorage.getItem('Authorization')}`;
        }

        return config;
    },
    (error) => {
        return Promise.reject(error);
    }
);

/**
 * http response 拦截器
 */
axios.interceptors.response.use(
    (response) => {
        return response;
    },
    (error) => {
        switch (error?.response?.status){
            case 401:
                window.localStorage.removeItem('Authorization');
                window.history.replaceState(null, '', '/login');
            case 403:
                message.error("无权限访问!")
                window.history.replaceState(null, '', '/');
        }

        return Promise.reject(error);
    }
);

/**
 * 生成缓存键
 */
function getCacheKey(method, url, params) {
    return `${method}:${url}:${JSON.stringify(params)}`;
}

/**
 * 检查缓存
 */
function checkCache(method, url, params) {
    const key = getCacheKey(method, url, params);
    const cached = cache.get(key);
    if (cached && Date.now() - cached.timestamp < CACHE_DURATION) {
        return cached.data;
    }
    cache.delete(key);
    return null;
}

/**
 * 设置缓存
 */
function setCache(method, url, params, data) {
    const key = getCacheKey(method, url, params);
    cache.set(key, {
        data,
        timestamp: Date.now()
    });
}

/**
 * 检查是否有正在进行的相同请求
 */
function checkPendingRequest(url, params) {
    const key = getCacheKey('get', url, params);
    return pendingRequests.has(key);
}

/**
 * 设置正在进行的请求
 */
function setPendingRequest(url, params) {
    const key = getCacheKey('get', url, params);
    pendingRequests.set(key, true);
}

/**
 * 清除正在进行的请求
 */
function clearPendingRequest(url, params) {
    const key = getCacheKey('get', url, params);
    pendingRequests.delete(key);
}

/**
 * 重试机制
 */
async function requestWithRetry(method, url, data, options = {}) {
    const maxRetries = options.maxRetries || 3;
    const retryDelay = options.retryDelay || 1000;

    for (let i = 0; i < maxRetries; i++) {
        try {
            if (method === 'get') {
                return await axios.get(url, data);
            } else if (method === 'post') {
                return await axios.post(url, data);
            } else if (method === 'put') {
                return await axios.put(url, data);
            } else if (method === 'delete') {
                return await axios.delete(url, { data });
            } else if (method === 'patch') {
                return await axios.patch(url, data);
            }
        } catch (error) {
            const isRetryable = !error.response || 
                (error.response.status >= 500 && error.response.status < 600) ||
                error.code === 'ECONNABORTED' ||
                error.code === 'ETIMEDOUT';
            
            if (!isRetryable || i === maxRetries - 1) {
                throw error;
            }

            // 等待后重试
            await new Promise(resolve => setTimeout(resolve, retryDelay * (i + 1)));
        }
    }
}

/**
 * 封装get方法（带缓存）
 */
export function get(url, params = {}, options = {}) {
    return new Promise((resolve, reject) => {
        // 检查缓存
        if (!options.skipCache) {
            const cached = checkCache('get', url, params);
            if (cached) {
                return resolve(cached);
            }
        }

        requestWithRetry('get', url, { params })
            .then((response) => {
                // 设置缓存
                if (!options.skipCache) {
                    setCache('get', url, params, response.data);
                }
                resolve(response.data);
            })
            .catch((error) => {
                reject(error);
            });
    });
}

/**
 * 封装post方法
 */
export function post(url, data, options = {}) {
    return new Promise((resolve, reject) => {
        requestWithRetry('post', url, data)
            .then((response) => {
                resolve(response.data);
            })
            .catch((error) => {
                reject(error);
            });
    });
}

/**
 * 封装patch方法
 */
export function patch(url, data = {}, options = {}) {
    return new Promise((resolve, reject) => {
        requestWithRetry('patch', url, data, options)
            .then((response) => {
                resolve(response.data);
            })
            .catch((err) => {
                msag(err);
                reject(err);
            });
    });
}

/**
 * 封装put方法
 */
export function put(url, data = {}, options = {}) {
    return new Promise((resolve, reject) => {
        requestWithRetry('put', url, data, options)
            .then((response) => {
                resolve(response.data);
            })
            .catch((err) => {
                msag(err);
                reject(err);
            });
    });
}

/**
 * 封装delete方法
 */
export function del(url, data = {}, options = {}) {
    return new Promise((resolve, reject) => {
        requestWithRetry('delete', url, { data }, options)
            .then((response) => {
                resolve(response.data);
            })
            .catch((error) => {
                reject(error);
            });
    });
}

/**
 * 统一接口处理，返回数据
 */
export default function (method, url, param, options = {}) {
    return new Promise((resolve, reject) => {
        switch (method) {
            case 'get':
                get(url, param, options)
                    .then(function (response) {
                        resolve(response);
                    })
                    .catch(function (error) {
                        reject(error);
                    });
                break;
            case 'post':
                post(url, param, options)
                    .then(function (response) {
                        resolve(response);
                    })
                    .catch(function (error) {
                        console.error('get request POST failed.', error);
                        reject(error);
                    });
                break;
            case 'put':
                put(url, param, options)
                    .then(function (response) {
                        resolve(response);
                    })
                    .catch(function (error) {
                        reject(error);
                    });
                break;
            case 'delete':
                del(url, param, options)
                    .then(function (response) {
                        resolve(response);
                    })
                    .catch(function (error) {
                        reject(error);
                    });
                break;
            case 'patch':
                patch(url, param, options)
                    .then(function (response) {
                        resolve(response);
                    })
                    .catch(function (error) {
                        reject(error);
                    });
                break;
            default:
                break;
        }
    });
}

// 清除所有缓存
export function clearCache() {
    cache.clear();
}

// 清除特定缓存
export function clearCacheByUrl(url) {
    for (const [key] of cache.keys()) {
        if (key.includes(url)) {
            console.log('[clearCacheByUrl] Deleting cache key:', key)
            cache.delete(key);
        }
    }
    // 同时清除所有正在进行的请求
    for (const [key] of pendingRequests.keys()) {
        if (key.includes(url)) {
            console.log('[clearCacheByUrl] Clearing pending request:', key)
            pendingRequests.delete(key);
        }
    }
}

// 失败提示
function msag(err) {
    if (err && err.response) {
        switch (err.response.status) {
            case 400:
                alert(err.response.data.error.details);
                break;
            case 401:
                alert('未授权，请登录');
                break;

            case 403:
                alert('拒绝访问');
                break;

            case 404:
                alert('请求地址出错');
                break;

            case 408:
                alert('请求超时');
                break;

            case 500:
                alert('服务器内部错误');
                break;

            case 501:
                alert('服务未实现');
                break;

            case 502:
                alert('网关错误');
                break;

            case 503:
                alert('服务不可用');
                break;

            case 504:
                alert('网关超时');
                break;

            case 505:
                alert('HTTP版本不受支持');
                break;
            default:
        }
    }
}