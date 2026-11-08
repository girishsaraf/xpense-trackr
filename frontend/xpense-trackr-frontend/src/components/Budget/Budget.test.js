import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import BudgetList from './BudgetList';
import BudgetForm from './BudgetForm';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

describe('Budget pages', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        axios.get.mockImplementation((url) => {
            if (url.endsWith('/api/budgets')) return Promise.resolve({ data: [{ id: 1, category_id: 1, amount: '100' }] });
            if (url.endsWith('/api/categories')) return Promise.resolve({ data: [{ id: 1, name: 'Food' }] });
            if (url.endsWith('/api/expenses')) return Promise.resolve({ data: [{ category_id: '1', amount: '40', date: `${new Date().getFullYear()}-${String(new Date().getMonth() + 1).padStart(2, '0')}-10` }] });
            return Promise.reject(new Error(`Unexpected GET ${url}`));
        });
        axios.post.mockResolvedValue({});
    });

    test('shows budget progress against this month’s category spending', async () => {
        renderWithRouter(<BudgetList />);

        expect(await screen.findByText('Food')).toBeInTheDocument();
        expect(screen.getByText('$40.00')).toBeInTheDocument();
        expect(screen.getByText('$60.00')).toBeInTheDocument();
        expect(screen.getByRole('progressbar', { name: 'Food budget used' })).toHaveAttribute('aria-valuenow', '40');
    });

    test('creates a budget with a category and monthly limit', async () => {
        axios.get.mockImplementation((url) => {
            if (url.endsWith('/api/budgets')) return Promise.resolve({ data: [] });
            if (url.endsWith('/api/categories')) return Promise.resolve({ data: [{ id: 1, name: 'Food' }] });
            return Promise.reject(new Error(`Unexpected GET ${url}`));
        });
        renderWithRouter(<BudgetForm />, {
            route: '/budgets/new',
            path: '/budgets/new',
            extraRoutes: [{ path: '/budgets', element: <div>Budget list</div> }],
        });

        await screen.findByRole('option', { name: 'Food' });
        fireEvent.change(screen.getByLabelText('Category'), { target: { value: '1' } });
        fireEvent.change(screen.getByLabelText('Monthly limit'), { target: { value: '250' } });
        fireEvent.click(screen.getByRole('button', { name: 'Create budget' }));

        await screen.findByText('Budget list');
        await waitFor(() => expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/budgets`,
            { category_id: 1, amount: 250 }
        ));
    });
});
