import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import ExpenseList from './ExpenseList';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

const expenses = [
    { id: 1, category_id: '1', description: 'Groceries', amount: '12.5', date: '2026-01-15' },
    { id: 2, category_id: '2', description: 'Bus ticket', amount: '3', date: '2026-01-16' },
];
const categories = [
    { id: 1, name: 'Food' },
    { id: 2, name: 'Travel' },
];

function renderExpenseList() {
    return renderWithRouter(<ExpenseList />, { route: '/expenses', path: '/expenses' });
}

beforeEach(() => {
    jest.clearAllMocks();
    window.alert = jest.fn();

    axios.get.mockImplementation((url) => {
        if (url.endsWith('/api/expenses')) {
            return Promise.resolve({ data: expenses.map((expense) => ({ ...expense })) });
        }
        if (url.endsWith('/api/categories')) {
            return Promise.resolve({ data: categories });
        }
        if (/\/api\/expenses\/\d+$/.test(url)) {
            return Promise.resolve({ data: { ...expenses[0] } });
        }
        return Promise.reject(new Error(`Unexpected GET ${url}`));
    });
    axios.put.mockResolvedValue({});
    axios.delete.mockResolvedValue({});
});

describe('ExpenseList', () => {
    test('renders expenses with category names', async () => {
        renderExpenseList();

        expect(await screen.findByText('Groceries')).toBeInTheDocument();
        expect(screen.getByText('Food')).toBeInTheDocument();
        expect(screen.getByText('Bus ticket')).toBeInTheDocument();
        expect(screen.getByText('Travel')).toBeInTheDocument();
        expect(screen.getByText('$12.5')).toBeInTheDocument();
        expect(screen.getByRole('link', { name: /Add New/ })).toHaveAttribute('href', '/expenses/new');
    });

    test('deletes an expense and refreshes the list', async () => {
        renderExpenseList();
        await screen.findByText('Groceries');

        fireEvent.click(screen.getByLabelText('Delete expense 2'));

        await waitFor(() => {
            expect(axios.delete).toHaveBeenCalledWith(`${backendUrl}/api/expenses/2`);
        });
        await waitFor(() => expect(window.alert).toHaveBeenCalledWith('Expense deleted successfully'));
        expect(axios.get).toHaveBeenCalledWith(`${backendUrl}/api/expenses`);
    });

    test('edits an expense through the modal', async () => {
        renderExpenseList();
        await screen.findByText('Groceries');

        fireEvent.click(screen.getByLabelText('Edit expense 1'));
        expect(screen.getByText('Edit Expense')).toBeInTheDocument();

        fireEvent.change(screen.getByLabelText('Description:'), { target: { value: 'Updated groceries' } });
        fireEvent.change(screen.getByLabelText('Amount:'), { target: { value: '20' } });
        fireEvent.click(screen.getByRole('button', { name: 'Submit' }));

        await waitFor(() => {
            expect(axios.put).toHaveBeenCalledWith(
                `${backendUrl}/api/expenses/1`,
                expect.objectContaining({ description: 'Updated groceries', amount: '20' })
            );
        });
        expect(window.alert).toHaveBeenCalledWith('Expense updated successfully');
        await waitFor(() => expect(screen.queryByText('Edit Expense')).not.toBeInTheDocument());
    });
});
