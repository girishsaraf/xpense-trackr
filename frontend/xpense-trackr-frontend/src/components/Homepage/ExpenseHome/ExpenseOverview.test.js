import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import ExpenseOverview from './ExpenseOverview';

jest.mock('axios');

describe('ExpenseOverview', () => {
    afterEach(() => {
        jest.clearAllMocks();
        jest.spyOn(console, 'error').mockRestore();
    });

    test('groups expenses by category', async () => {
        axios.get.mockImplementation((url) => {
            if (url.endsWith('/api/expenses')) {
                return Promise.resolve({
                    data: [
                        { category_id: '1', amount: '10.5' },
                        { category_id: '1', amount: '4.5' },
                        { category_id: '2', amount: '5' },
                    ],
                });
            }
            return Promise.resolve({
                data: [
                    { id: 1, name: 'Food' },
                    { id: 2, name: 'Travel' },
                ],
            });
        });

        render(<ExpenseOverview />);

        expect(await screen.findByText('Food:')).toBeInTheDocument();
        expect(screen.getByText('$15')).toBeInTheDocument();
        expect(screen.getByText('Travel:')).toBeInTheDocument();
        expect(screen.getByText('$5')).toBeInTheDocument();
    });

    test('reports fetch failures', async () => {
        const consoleError = jest.spyOn(console, 'error').mockImplementation(() => {});
        axios.get.mockRejectedValue(new Error('network'));

        render(<ExpenseOverview />);

        await waitFor(() => expect(consoleError).toHaveBeenCalledWith('Error fetching data:', expect.any(Error)));
    });
});
