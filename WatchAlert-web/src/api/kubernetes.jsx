import http from '../utils/http';
import {HandleApiError} from "../utils/lib";

async function getKubernetesResourceList(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/kubernetes/getResourceList', params, options);
        return res;
    } catch (error) {
        HandleApiError(error)
        return error
    }
}

async function getKubernetesReasonList(params, options = {}) {
    try {
        const res = await http('get', '/api/w8t/kubernetes/getReasonList', params, options);
        return res;
    } catch (error) {
        HandleApiError(error)
        return error
    }
}

export {
    getKubernetesResourceList,
    getKubernetesReasonList
}