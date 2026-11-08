import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import ForecastPage from './ForecastPage';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

const accounts = [
    { id: 1, name: 'Checking', type: 'checking', currency: 'USD', balance_minor: 12500 },
];
const forecast = {
    start_date: '2026-01-01',
    end_date: '2026-04-01',
    days: 90,
    accounts: [{
        account_id: 1,
        name: 'Checking',
        currency: 'USD',
        current_balance_minor: 12500,
        projected_balance_minor: 10000,
        points: [],
    }],
    events: [],
};

beforeEach(() => {
    jest.clearAllMocks();
    axios.get.mockImplementation((url) => {
        if (url.endsWith('/api/accounts')) return Promise.resolve({ data: accounts });
        if (url.endsWith('/api/scheduled-flows')) return Promise.resolve({ data: [] });
        if (url.includes('/api/forecast/cashflow')) return Promise.resolve({ data: forecast });
        return Promise.reject(new Error(`Unexpected GET ${url}`));
    });
    axios.post.mockResolvedValue({});
});

test('creates an expected bill in the cashflow forecast', async () => {
    renderWithRouter(<ForecastPage />, { route: '/forecast', path: '/forecast' });

    expect(await screen.findByRole('heading', { name: 'Checking' })).toBeInTheDocument();
    fireEvent.change(screen.getByLabelText('Account'), { target: { value: '1' } });
    fireEvent.change(screen.getByLabelText('Name'), { target: { value: 'Rent' } });
    fireEvent.change(screen.getByLabelText('Amount (USD)'), { target: { value: '850.25' } });
    fireEvent.change(screen.getByLabelText('Next date'), { target: { value: '2026-01-15' } });
    fireEvent.click(screen.getByRole('button', { name: 'Add to forecast' }));

    await waitFor(() => {
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/scheduled-flows`,
            expect.objectContaining({
                account_id: 1,
                label: 'Rent',
                type: 'expense',
                amount_minor: 85025,
                currency: 'USD',
                frequency: 'monthly',
                next_date: '2026-01-15',
            }),
        );
    });
    expect(await screen.findByText('Scheduled item added to your forecast.')).toBeInTheDocument();
});
