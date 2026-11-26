import apiClient from './client';

export const listExpenses = () => apiClient.get('/api/expenses');

export const getExpense = (id) => apiClient.get(`/api/expenses/${id}`);

export const createExpense = (expense) => apiClient.post('/api/expenses', expense);

export const updateExpense = (id, expense) => apiClient.put(`/api/expenses/${id}`, expense);

export const deleteExpense = (id) => apiClient.delete(`/api/expenses/${id}`);

export const createRecurringExpense = (expense) =>
    apiClient.post('/api/recurring/expenses', expense);