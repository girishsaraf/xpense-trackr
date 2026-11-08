import React from 'react';
import { act, fireEvent, render, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import backendUrl from '../../config';
import PortfolioPage from './PortfolioPage';

jest.mock('axios');
global.IS_REACT_ACT_ENVIRONMENT = true;

beforeEach(() => {
    axios.get.mockImplementation((url) => {
        if (url.endsWith('/api/portfolio/holdings/1/valuations')) {
            return Promise.resolve({
                data: [{
                    id: 2,
                    date: '2026-10-01',
                    source: 'Broker statement',
                    value_minor: 25000,
                }],
            });
        }
        return Promise.resolve({
            data: [{
                id: 1,
                name: 'Index fund',
                ticker: 'VTI',
                asset_class: 'equity',
                currency: 'USD',
                quantity: '2',
                cost_basis_minor: 20000,
                latest_value_minor: 25000,
                valuation_date: '2026-10-01',
                valuation_source: 'Broker statement',
                has_valuation: true,
            }],
        });
    });
    axios.post.mockResolvedValue({ data: {} });
});

afterEach(() => {
    jest.clearAllMocks();
});

test('shows portfolio allocation and valuation history', async () => {
    render(<PortfolioPage />);

    expect(await screen.findByRole('heading', { name: 'Index fund' })).toBeInTheDocument();
    expect(screen.getAllByText('Stocks / equity')).toHaveLength(2);
    expect(screen.getByText('Gain vs. cost basis: $50.00')).toBeInTheDocument();

    await act(async () => {
        fireEvent.click(screen.getByRole('button', { name: 'Show valuation history' }));
        await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(screen.getByText('2026-10-01 · Broker statement')).toBeInTheDocument();
    expect(axios.get).toHaveBeenCalledWith(`${backendUrl}/api/portfolio/holdings/1/valuations`);
});

test('creates a holding with exact minor-unit cost basis and valuation', async () => {
    render(<PortfolioPage />);
    await screen.findByRole('heading', { name: 'Index fund' });

    fireEvent.change(screen.getByLabelText('Name'), { target: { value: 'Treasury fund' } });
    fireEvent.change(screen.getByLabelText('Quantity'), { target: { value: '1.25' } });
    fireEvent.change(screen.getByLabelText('Cost basis (USD)'), { target: { value: '100.00' } });
    fireEvent.change(screen.getByLabelText('Current value (USD)'), { target: { value: '105.50' } });
    await act(async () => {
        fireEvent.click(screen.getByRole('button', { name: 'Add holding' }));
        await new Promise((resolve) => setTimeout(resolve, 0));
    });
    await waitFor(() => {
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/portfolio/holdings`,
            expect.objectContaining({
                name: 'Treasury fund',
                asset_class: 'equity',
                currency: 'USD',
                quantity: '1.25',
                cost_basis_minor: 10000,
                value_minor: 10550,
            }),
        );
    });
});
