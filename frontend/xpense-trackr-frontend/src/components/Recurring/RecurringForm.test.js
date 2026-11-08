import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import RecurringForm from './RecurringForm';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

describe('RecurringForm', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        axios.get.mockImplementation((url) => {
            if (url.endsWith('/api/categories')) return Promise.resolve({ data: [{ id: 2, name: 'Utilities' }] });
            return Promise.resolve({
                data: {
                    id: 4,
                    description: 'Internet',
                    amount: 40,
                    category_id: 2,
                    start_date: '2026-01-01T00:00:00Z',
                    end_date: '2026-12-01T00:00:00Z',
                    frequency: 'monthly',
                },
            });
        });
        axios.put.mockResolvedValue({});
    });

    test('updates an expense schedule using API-compatible dates and numeric values', async () => {
        renderWithRouter(<RecurringForm />, {
            route: '/recurring/expenses/4/edit',
            path: '/recurring/:kind/:id/edit',
            extraRoutes: [{ path: '/recurring', element: <div>Recurring activity</div> }],
        });

        fireEvent.change(await screen.findByLabelText('Amount per occurrence'), { target: { value: '55' } });
        fireEvent.change(screen.getByLabelText(/End date/), { target: { value: '2027-01-01' } });
        fireEvent.click(screen.getByRole('button', { name: 'Save changes' }));

        await screen.findByText('Recurring activity');
        await waitFor(() => expect(axios.put).toHaveBeenCalledWith(
            `${backendUrl}/api/recurring/expenses/4`,
            expect.objectContaining({
                amount: 55,
                category_id: 2,
                start_date: '2026-01-01T00:00:00Z',
                end_date: '2027-01-01T00:00:00Z',
            })
        ));
    });
});
