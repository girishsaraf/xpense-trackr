import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import InvestmentOverview from './InvestmentOverview';

jest.mock('axios');

describe('InvestmentOverview', () => {
    afterEach(() => {
        jest.clearAllMocks();
        jest.spyOn(console, 'error').mockRestore();
    });

    test('groups investments by type', async () => {
        axios.get.mockResolvedValue({
            data: [
                { investment_type: 'stocks', amount: '100' },
                { investment_type: 'stocks', amount: '50' },
                { investment_type: 'bonds', amount: '25.5' },
            ],
        });

        render(<InvestmentOverview />);

        expect(await screen.findByText('stocks:')).toBeInTheDocument();
        expect(screen.getByText('$150')).toBeInTheDocument();
        expect(screen.getByText('bonds:')).toBeInTheDocument();
        expect(screen.getByText('$25.5')).toBeInTheDocument();
    });

    test('reports fetch failures', async () => {
        const consoleError = jest.spyOn(console, 'error').mockImplementation(() => {});
        axios.get.mockRejectedValue(new Error('network'));

        render(<InvestmentOverview />);

        await waitFor(() => expect(consoleError).toHaveBeenCalledWith('Error fetching data:', expect.any(Error)));
    });
});
