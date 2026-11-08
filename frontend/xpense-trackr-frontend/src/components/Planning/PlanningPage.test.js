import React from 'react';
import { act, fireEvent, render, screen, waitFor, within } from '@testing-library/react';
import axios from 'axios';
import PlanningPage from './PlanningPage';

jest.mock('axios');

describe('PlanningPage', () => {
    beforeEach(() => {
        axios.get.mockImplementation((url) => {
            if (url.includes('/contributions')) {
                return Promise.resolve({ data: [{ id: 8, date: '2026-10-01', direction: 'contribution', amount_minor: 2500 }] });
            }
            if (url.endsWith('/goals')) {
                return Promise.resolve({
                    data: [{
                        id: 1,
                        name: 'Emergency fund',
                        currency: 'USD',
                        target_amount_minor: 10000,
                        current_amount_minor: 2500,
                        remaining_minor: 7500,
                        required_monthly_minor: 1250,
                    }],
                });
            }
            if (url.endsWith('/envelopes')) {
                return Promise.resolve({
                    data: [{
                        id: 2,
                        category_id: 3,
                        currency: 'USD',
                        monthly_limit_minor: 10000,
                        spent_minor: 3000,
                        rollover_minor: 1000,
                        available_minor: 8000,
                    }],
                });
            }
            return Promise.resolve({ data: [{ id: 3, name: 'Food' }] });
        });
        axios.post.mockResolvedValue({ data: {} });
        axios.delete.mockResolvedValue({ data: {} });
    });

    afterEach(() => {
        jest.clearAllMocks();
    });

    it('records a goal contribution and creates a monthly envelope', async () => {
        await act(async () => {
            render(<PlanningPage />);
            await Promise.resolve();
            await Promise.resolve();
        });

        const goalHeading = await screen.findByRole('heading', { name: 'Emergency fund' });
        const goalCard = goalHeading.closest('article');
        expect(within(goalCard).getByText(/25\.00 of \$100\.00/)).toBeInTheDocument();
        fireEvent.click(within(goalCard).getByRole('button', { name: 'View history' }));
        expect(await within(goalCard).findByText('2026-10-01 · contribution')).toBeInTheDocument();
        fireEvent.change(within(goalCard).getByLabelText('Amount'), { target: { value: '10.00' } });
        await act(async () => {
            fireEvent.click(within(goalCard).getByRole('button', { name: 'Record' }));
            await waitFor(() => {
                expect(axios.post).toHaveBeenCalledWith(
                    expect.stringContaining('/api/planning/goals/1/contributions'),
                    expect.objectContaining({ amount_minor: 1000, direction: 'contribution' }),
                );
            });
            await Promise.resolve();
        });

        fireEvent.change(screen.getByLabelText('Category'), { target: { value: '3' } });
        fireEvent.change(screen.getByLabelText('Monthly limit'), { target: { value: '50.00' } });
        await act(async () => {
            fireEvent.click(screen.getByRole('button', { name: 'Add envelope' }));
            await waitFor(() => {
                expect(axios.post).toHaveBeenCalledWith(
                    expect.stringContaining('/api/planning/envelopes'),
                    expect.objectContaining({ category_id: 3, monthly_limit_minor: 5000, rollover: false }),
                );
            });
            await Promise.resolve();
        });
    });
});
