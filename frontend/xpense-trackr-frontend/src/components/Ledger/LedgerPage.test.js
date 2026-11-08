import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import LedgerPage from './LedgerPage';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

const accounts = [
    { id: 1, name: 'Checking', type: 'checking', currency: 'USD', balance_minor: 12500 },
    { id: 2, name: 'Savings', type: 'savings', currency: 'USD', balance_minor: 2000 },
    { id: 3, name: 'Euro account', type: 'checking', currency: 'EUR', balance_minor: 2000 },
];

beforeEach(() => {
    jest.clearAllMocks();
    axios.get.mockImplementation((url) => {
        if (url.endsWith('/api/accounts')) return Promise.resolve({ data: accounts });
        if (url.endsWith('/api/categories')) return Promise.resolve({ data: [{ id: 3, name: 'Food' }] });
        if (url.endsWith('/api/category-rules')) return Promise.resolve({ data: [] });
        if (url.includes('/api/transactions')) return Promise.resolve({ data: [] });
        return Promise.reject(new Error(`Unexpected GET ${url}`));
    });
    axios.post.mockResolvedValue({});
});

test('creates an exact-minor-unit transaction from the local accounts view', async () => {
    renderWithRouter(<LedgerPage />, { route: '/transactions', path: '/transactions' });

    expect(await screen.findByRole('heading', { name: 'Checking' })).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText('Account'), { target: { value: '1' } });
    fireEvent.change(screen.getByLabelText('Category'), { target: { value: '3' } });
    fireEvent.change(screen.getByLabelText('Amount (USD)'), { target: { value: '12.34' } });
    fireEvent.change(screen.getByLabelText('Date'), { target: { value: '2026-11-08' } });
    fireEvent.change(screen.getByLabelText('Payee'), { target: { value: 'Market' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save transaction' }));

    await waitFor(() => {
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/transactions`,
            expect.objectContaining({
                account_id: 1,
                category_id: 3,
                type: 'expense',
                amount_minor: 1234,
                currency: 'USD',
                date: '2026-11-08',
                payee: 'Market',
            }),
        );
    });
    expect(await screen.findByText('Transaction added.')).toBeInTheDocument();
});

test('transfer options only include accounts using the selected currency', async () => {
    renderWithRouter(<LedgerPage />, { route: '/transactions', path: '/transactions' });

    await screen.findByRole('heading', { name: 'Checking' });
    fireEvent.change(screen.getByLabelText('Type'), { target: { value: 'transfer' } });
    fireEvent.change(screen.getByLabelText('From account'), { target: { value: '1' } });

    const destinations = screen.getByLabelText('To account');
    expect(destinations).toHaveTextContent('Savings');
    expect(destinations).not.toHaveTextContent('Euro account');
});

test('creates a merchant categorization rule', async () => {
    renderWithRouter(<LedgerPage />, { route: '/transactions', path: '/transactions' });

    await screen.findByRole('heading', { name: 'Checking' });
    fireEvent.change(screen.getByLabelText('Rule name'), { target: { value: 'Groceries' } });
    fireEvent.change(screen.getByLabelText('Payee contains'), { target: { value: 'market' } });
    fireEvent.change(screen.getByLabelText('Apply category'), { target: { value: '3' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save rule' }));

    await waitFor(() => {
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/category-rules`,
            expect.objectContaining({
                name: 'Groceries',
                match_value: 'market',
                category_id: 3,
                enabled: true,
            }),
        );
    });
    expect(await screen.findByText('Categorization rule added.')).toBeInTheDocument();
});
