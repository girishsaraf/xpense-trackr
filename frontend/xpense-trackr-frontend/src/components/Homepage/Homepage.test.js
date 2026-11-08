import React from 'react';
import { screen } from '@testing-library/react';
import axios from 'axios';
import Homepage from './Homepage';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

beforeEach(() => {
    jest.clearAllMocks();
    axios.get.mockImplementation((url) => {
        if (url.endsWith('/api/accounts')) {
            return Promise.resolve({ data: [{ id: 1, name: 'Checking', type: 'checking', currency: 'USD', balance_minor: 12500 }] });
        }
        if (url.includes('/api/transactions')) {
            return Promise.resolve({ data: [{ id: 2, account_id: 1, date: '2026-10-12', type: 'expense', payee: 'Market', amount_minor: 1250, currency: 'USD' }] });
        }
        if (url.includes('/api/reports/cashflow')) {
            return Promise.resolve({ data: { currencies: [{ currency: 'USD', income_minor: 5000, expense_minor: 1250, net_minor: 3750 }] } });
        }
        if (url.endsWith('/api/net-worth')) {
            return Promise.resolve({ data: { currencies: [{ currency: 'USD', net_worth_minor: 12500 }] } });
        }
        if (url.includes('/api/forecast/cashflow')) {
            return Promise.resolve({ data: { events: [{ scheduled_flow_id: 3, date: '2026-10-15', label: 'Rent', type: 'expense', amount_minor: 85000, currency: 'USD' }] } });
        }
        return Promise.reject(new Error(`Unexpected GET ${url}`));
    });
});

test('shows local account, cashflow, forecast, and activity summaries', async () => {
    renderWithRouter(<Homepage />, { route: '/', path: '/' });

    expect(await screen.findByRole('heading', { name: 'Your money, in focus.' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Income and spending' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Your accounts' })).toBeInTheDocument();
    expect(screen.getByText('Checking')).toBeInTheDocument();
    expect(screen.getByText('Rent')).toBeInTheDocument();
    expect(screen.getByText('Market')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Add transaction' })).toHaveAttribute('href', '/transactions');
});

test('guides a first-run user to add an account or import data', async () => {
    axios.get.mockImplementation((url) => {
        if (url.endsWith('/api/accounts')) return Promise.resolve({ data: [] });
        if (url.includes('/api/transactions')) return Promise.resolve({ data: [] });
        if (url.includes('/api/reports/cashflow')) return Promise.resolve({ data: { currencies: [] } });
        if (url.endsWith('/api/net-worth')) return Promise.resolve({ data: { currencies: [] } });
        if (url.includes('/api/forecast/cashflow')) return Promise.resolve({ data: { events: [] } });
        return Promise.reject(new Error(`Unexpected GET ${url}`));
    });
    renderWithRouter(<Homepage />, { route: '/', path: '/' });

    expect(await screen.findByRole('heading', { name: 'Your finances stay on this device.' })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Set up an account' })).toHaveAttribute('href', '/transactions');
    expect(screen.getByRole('link', { name: 'Import existing transactions' })).toHaveAttribute('href', '/import-export');
});
