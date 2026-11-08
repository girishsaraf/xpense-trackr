import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import backendUrl from '../../config';
import AssetsPage from './AssetsPage';

jest.mock('axios');

beforeEach(() => {
    axios.get.mockResolvedValue({
        data: [{
            id: 8,
            name: 'Home',
            kind: 'asset',
            currency: 'USD',
            latest_value_minor: 25000000,
            valuation_date: '2026-10-01',
            valuation_source: 'Appraisal',
            has_valuation: true,
        }],
    });
    axios.post.mockResolvedValue({ data: {} });
});

afterEach(() => {
    jest.clearAllMocks();
});

test('records dated asset valuations with their source', async () => {
    render(<AssetsPage />);
    expect(await screen.findByRole('heading', { name: 'Home' })).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('New value (USD)'), { target: { value: '260000' } });
    fireEvent.change(screen.getByLabelText('Source'), { target: { value: 'Annual appraisal' } });
    fireEvent.click(screen.getByRole('button', { name: 'Add valuation' }));

    await waitFor(() => {
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/assets/8/valuations`,
            expect.objectContaining({ value_minor: 26000000, source: 'Annual appraisal' }),
        );
    });
    expect(await screen.findByRole('status')).toHaveTextContent('Valuation added for Home.');
});
