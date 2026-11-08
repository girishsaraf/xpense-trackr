import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import InvestmentForm from './InvestmentForm';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

function renderInvestmentForm() {
    return renderWithRouter(<InvestmentForm />, {
        route: '/investments/new',
        path: '/investments/new',
        extraRoutes: [{ path: '/investments', element: <div>Investments Page</div> }],
    });
}

function fillInvestmentFields() {
    fireEvent.change(screen.getByLabelText('Investment Type:'), { target: { value: 'bonds' } });
    fireEvent.change(screen.getByLabelText('Amount:'), { target: { value: '300' } });
    fireEvent.change(screen.getByLabelText('Date:'), { target: { value: '2026-02-01' } });
}

beforeEach(() => {
    jest.clearAllMocks();
    axios.post.mockResolvedValue({});
});

describe('InvestmentForm', () => {
    test('submits a one-off investment', async () => {
        renderInvestmentForm();
        fillInvestmentFields();

        fireEvent.click(screen.getByRole('button', { name: 'Add investment' }));

        await screen.findByText('Investments Page');
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/investments`,
            expect.objectContaining({
                investment_type: 'bonds',
                amount: '300',
                date: '2026-02-01',
                recurring: false,
            })
        );
    });

    test('submits a recurring investment', async () => {
        renderInvestmentForm();
        fillInvestmentFields();

        fireEvent.click(screen.getByLabelText('Is Recurring Investment?'));
        fireEvent.change(screen.getByLabelText('Description:'), { target: { value: 'Monthly contribution' } });
        fireEvent.change(screen.getByLabelText('Start Date:'), { target: { value: '2026-02-01' } });
        fireEvent.change(screen.getByLabelText('End Date:'), { target: { value: '2026-08-01' } });
        fireEvent.change(screen.getByLabelText('Frequency:'), { target: { value: 'weekly' } });

        fireEvent.click(screen.getByRole('button', { name: 'Add investment' }));

        await screen.findByText('Investments Page');
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/recurring/investments`,
            expect.objectContaining({
                recurring: true,
                frequency: 'weekly',
                amount: 300,
                description: 'Monthly contribution',
                start_date: '2026-02-01T00:00:00Z',
                end_date: '2026-08-01T00:00:00Z',
            })
        );
    });

    test('navigates back when cancelled', async () => {
        renderInvestmentForm();

        fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));

        await screen.findByText('Investments Page');
        expect(axios.post).not.toHaveBeenCalled();
    });

    test('shows a notification when the request fails', async () => {
        axios.post.mockRejectedValue(new Error('network'));
        renderInvestmentForm();
        fillInvestmentFields();

        fireEvent.click(screen.getByRole('button', { name: 'Add investment' }));

        expect(await screen.findByText('Error adding investment')).toBeInTheDocument();
        await waitFor(() => expect(axios.post).toHaveBeenCalled());
    });
});
