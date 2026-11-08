import React from 'react';
import { screen } from '@testing-library/react';
import axios from 'axios';
import ReportsPage from './ReportsPage';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

beforeEach(() => {
    jest.clearAllMocks();
    axios.get.mockImplementation((url) => {
        if (url.includes('/api/reports/cashflow')) {
            return Promise.resolve({
                data: {
                    from: '2026-10-01',
                    to: '2026-10-31',
                    currencies: [{
                        currency: 'USD',
                        income_minor: 5000,
                        expense_minor: 1250,
                        net_minor: 3750,
                        categories: [{ category_id: null, category_name: '', amount_minor: 1250 }],
                    }],
                },
            });
        }
        if (url.endsWith('/api/net-worth')) {
            return Promise.resolve({
                data: {
                    as_of: '2026-10-31',
                    currencies: [{
                        currency: 'USD',
                        assets_minor: 10000,
                        liabilities_minor: 0,
                        net_worth_minor: 10000,
                        accounts: [{ account_id: 1, name: 'Checking', type: 'checking', balance_minor: 10000 }],
                    }],
                },
            });
        }
        return Promise.reject(new Error(`Unexpected GET ${url}`));
    });
});

test('shows cashflow totals and currency-separated net worth', async () => {
    renderWithRouter(<ReportsPage />, { route: '/reports', path: '/reports' });

    expect(await screen.findByRole('heading', { name: 'USD cash flow' })).toBeInTheDocument();
    expect(screen.getByText('Uncategorized')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Net worth by currency' })).toBeInTheDocument();
    expect(screen.getByText('Checking')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Download CSV' })).toHaveAttribute('href', expect.stringContaining('/api/reports/cashflow.csv'));
});
