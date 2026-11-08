import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import InvestmentList from './InvestmentList';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

const investments = [
    { id: 1, investment_type: '401k', amount: '100', date: '2026-01-15' },
    { id: 2, investment_type: 'mutual_funds', amount: '250', date: '2026-01-16' },
];

function renderInvestmentList() {
    return renderWithRouter(<InvestmentList />, { route: '/investments', path: '/investments' });
}

beforeEach(() => {
    jest.clearAllMocks();
    window.alert = jest.fn();
    window.confirm = jest.fn(() => true);

    axios.get.mockImplementation((url) => {
        if (url.endsWith('/api/investments')) {
            return Promise.resolve({ data: investments.map((investment) => ({ ...investment })) });
        }
        if (/\/api\/investments\/\d+$/.test(url)) {
            return Promise.resolve({ data: { ...investments[0] } });
        }
        return Promise.reject(new Error(`Unexpected GET ${url}`));
    });
    axios.put.mockResolvedValue({});
    axios.delete.mockResolvedValue({});
});

describe('InvestmentList', () => {
    test('renders investments with friendly type names', async () => {
        renderInvestmentList();

        expect(await screen.findByText('401(k)')).toBeInTheDocument();
        expect(screen.getByText('Mutual Funds')).toBeInTheDocument();
        expect(screen.getByText('$100')).toBeInTheDocument();
        expect(screen.getByText('$250')).toBeInTheDocument();
        expect(screen.getByRole('link', { name: /Add New/ })).toHaveAttribute('href', '/investments/new');
    });

    test('deletes an investment and refreshes the list', async () => {
        renderInvestmentList();
        await screen.findByText('401(k)');

        fireEvent.click(screen.getByLabelText('Delete investment 2'));

        await waitFor(() => {
            expect(axios.delete).toHaveBeenCalledWith(`${backendUrl}/api/investments/2`);
        });
        await waitFor(() => expect(window.alert).toHaveBeenCalledWith('Investment deleted successfully'));
        expect(axios.get).toHaveBeenCalledWith(`${backendUrl}/api/investments`);
    });

    test('edits an investment through the modal', async () => {
        renderInvestmentList();
        await screen.findByText('401(k)');

        fireEvent.click(screen.getByLabelText('Edit investment 1'));
        expect(screen.getByText('Edit Investment')).toBeInTheDocument();

        fireEvent.change(screen.getByLabelText('Amount:'), { target: { value: '175' } });
        fireEvent.change(screen.getByLabelText('Type:'), { target: { value: 'savings' } });
        fireEvent.click(screen.getByRole('button', { name: 'Submit' }));

        await waitFor(() => {
            expect(axios.put).toHaveBeenCalledWith(
                `${backendUrl}/api/investments/1`,
                expect.objectContaining({ amount: '175', investment_type: 'savings' })
            );
        });
        expect(window.alert).toHaveBeenCalledWith('Investment updated successfully');
        await waitFor(() => expect(screen.queryByText('Edit Investment')).not.toBeInTheDocument());
    });
});
