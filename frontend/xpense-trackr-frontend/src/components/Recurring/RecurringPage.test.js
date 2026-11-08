import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import axios from 'axios';
import RecurringPage from './RecurringPage';

jest.mock('axios');

describe('RecurringPage', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        axios.get.mockImplementation((url) => {
            if (url.endsWith('/api/recurring/expenses')) {
                return Promise.resolve({ data: [{ id: 1, description: 'Internet', category_id: 1, amount: 50, frequency: 'monthly', end_date: '2027-01-01' }] });
            }
            if (url.endsWith('/api/recurring/investments')) {
                return Promise.resolve({ data: [{ id: 2, description: '401(k) contribution', investment_type: '401k', amount: 250, frequency: 'monthly' }] });
            }
            return Promise.resolve({ data: [{ id: 1, name: 'Utilities' }] });
        });
    });

    test('lists recurring expenses and switches to recurring investments', async () => {
        render(<MemoryRouter><RecurringPage /></MemoryRouter>);

        expect(await screen.findByText('Internet')).toBeInTheDocument();
        expect(screen.getByText(/Utilities/)).toBeInTheDocument();

        fireEvent.click(screen.getByRole('button', { name: 'Investments (1)' }));
        expect(await screen.findByText('401(k) contribution')).toBeInTheDocument();
    });
});
