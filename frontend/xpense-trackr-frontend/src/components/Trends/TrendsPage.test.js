import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import TrendsPage from './TrendsPage';
import backendUrl from '../../config';

jest.mock('axios');
jest.mock('../Charts/ChartVisualization', () => ({ ariaLabel }) => (
    <div role="img" aria-label={ariaLabel} />
));

describe('TrendsPage', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        axios.get.mockResolvedValue({
            data: [
                { month: '2026-01', category: 'Food', amount: 75 },
                { month: '2026-01', category: 'Travel', amount: 25 },
                { month: '2026-02', category: 'Food', amount: 100 },
            ],
        });
    });

    test('summarizes monthly spending and categories', async () => {
        render(<TrendsPage />);

        expect(await screen.findByText('$200.00')).toBeInTheDocument();
        expect(screen.getByText('Food')).toBeInTheDocument();
        expect(screen.getByRole('img', { name: 'Bar chart showing spending by month' })).toBeInTheDocument();
        expect(screen.getByRole('img', { name: 'Doughnut chart showing spending distribution by category' })).toBeInTheDocument();
    });

    test('reloads results when the selected period changes', async () => {
        render(<TrendsPage />);
        await screen.findByText('Monthly spending');

        fireEvent.click(screen.getByRole('button', { name: '12 months' }));

        await waitFor(() => expect(axios.get).toHaveBeenLastCalledWith(`${backendUrl}/api/trends/spending?months=12`));
        expect(screen.getByRole('button', { name: '12 months' })).toHaveAttribute('aria-pressed', 'true');
    });
});
