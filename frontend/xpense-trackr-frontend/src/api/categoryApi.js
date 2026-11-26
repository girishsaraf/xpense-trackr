import apiClient from './client';

export const listCategories = () => apiClient.get('/api/categories');