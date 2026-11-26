import axios from 'axios';
import backendUrl from '../config';

const apiClient = {
    get: (path, config) => config === undefined
        ? axios.get(`${backendUrl}${path}`)
        : axios.get(`${backendUrl}${path}`, config),
    post: (path, data, config) => config === undefined
        ? axios.post(`${backendUrl}${path}`, data)
        : axios.post(`${backendUrl}${path}`, data, config),
    put: (path, data, config) => config === undefined
        ? axios.put(`${backendUrl}${path}`, data)
        : axios.put(`${backendUrl}${path}`, data, config),
    delete: (path, config) => config === undefined
        ? axios.delete(`${backendUrl}${path}`)
        : axios.delete(`${backendUrl}${path}`, config),
};

export default apiClient;
