import React from 'react';
import { act, fireEvent, render, screen } from '@testing-library/react';
import axios from 'axios';
import backendUrl from '../../config';
import BillsPage from './BillsPage';

jest.mock('axios');

describe('BillsPage', () => {
    beforeEach(() => {
        axios.get.mockImplementation((url) => {
            if (url.endsWith('/api/bills')) {
                return Promise.resolve({
                    data: [{
                        id: 4,
                        label: 'Internet',
                        account_name: 'Checking',
                        category_id: 2,
                        currency: 'USD',
                        amount_minor: 1250,
                        next_date: '2026-11-09',
                        frequency: 'monthly',
                        status: 'due_soon',
                        days_until_due: 1,
                        has_price_change: true,
                        previous_amount_minor: 1000,
                        latest_amount_minor: 1250,
                        change_percent: 25,
                    }],
                });
            }
            return Promise.resolve({ data: [{ id: 2, name: 'Utilities' }] });
        });
        axios.post.mockResolvedValue({ data: {} });
    });

    afterEach(() => {
        jest.clearAllMocks();
    });

    it('shows price changes and records a confirmed actual payment', async () => {
        await act(async () => {
            render(<BillsPage />);
            await Promise.resolve();
            await Promise.resolve();
        });
        expect(await screen.findByRole('heading', { name: 'Internet' })).toBeInTheDocument();
        expect(screen.getByText(/increased by 25%/)).toBeInTheDocument();
        fireEvent.change(screen.getByLabelText('Actual amount paid (USD)'), { target: { value: '13.50' } });
        await act(async () => {
            fireEvent.click(screen.getByRole('button', { name: 'Record payment' }));
            await new Promise((resolve) => setTimeout(resolve, 0));
        });
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/scheduled-flows/4/occurrences`,
            { scheduled_date: '2026-11-09', action: 'record', amount_minor: 1350 },
        );
        expect(screen.getByRole('status')).toHaveTextContent('Payment recorded in your transaction ledger.');
    });

    it('records a skipped occurrence without creating a transaction', async () => {
        await act(async () => {
            render(<BillsPage />);
            await Promise.resolve();
            await Promise.resolve();
        });
        await screen.findByRole('heading', { name: 'Internet' });
        await act(async () => {
            fireEvent.click(screen.getByRole('button', { name: 'Skip this occurrence' }));
            await new Promise((resolve) => setTimeout(resolve, 0));
        });
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/scheduled-flows/4/occurrences`,
            { scheduled_date: '2026-11-09', action: 'skip', amount_minor: 0 },
        );
        expect(screen.getByRole('status')).toHaveTextContent('This scheduled occurrence was skipped.');
    });
});
