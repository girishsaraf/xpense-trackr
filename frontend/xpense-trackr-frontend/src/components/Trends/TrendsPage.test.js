import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import TrendsPage from './TrendsPage';
import backendUrl from '../../config';

jest.mock('axios');

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
        expect(screen.getAllByText('Food')).toHaveLength(2);
        expect(screen.getByRole('list', { name: 'Monthly spending amounts' })).toBeInTheDocument();
    });

    test('reloads results when the selected period changes', async () => {
        render(<TrendsPage />);
        await screen.findByText('Monthly spending');

        fireEvent.change(screen.getByLabelText('Time period'), { target: { value: '12' } });

        await waitFor(() => expect(axios.get).toHaveBeenLastCalledWith(`${backendUrl}/api/trends/spending?months=12`));
    });
});
