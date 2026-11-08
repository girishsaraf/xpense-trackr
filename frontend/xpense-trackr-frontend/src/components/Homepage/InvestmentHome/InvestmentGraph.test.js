import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import Chart from 'chart.js/auto';
import InvestmentGraph from './InvestmentGraph';
import backendUrl from '../../../config';

jest.mock('chart.js/auto', () => ({
    __esModule: true,
    default: jest.fn(function () {
        return { destroy: jest.fn() };
    }),
}));

describe('InvestmentGraph', () => {
    const originalFetch = global.fetch;
    const originalGetContext = HTMLCanvasElement.prototype.getContext;
    const graphData = [
        { month: '2026-01', type: '401k', amount: 500 },
        { month: '2026-02', type: 'stocks', amount: 300 },
    ];

    beforeEach(() => {
        Chart.mockClear();
        global.fetch = jest.fn().mockResolvedValue({
            ok: true,
            json: () => Promise.resolve(graphData),
        });
        HTMLCanvasElement.prototype.getContext = jest.fn(() => ({}));
    });

    afterEach(() => {
        global.fetch = originalFetch;
        HTMLCanvasElement.prototype.getContext = originalGetContext;
    });

    test('fetches data and renders a chart', async () => {
        render(<InvestmentGraph />);

        await waitFor(() => {
            expect(global.fetch).toHaveBeenCalledWith(`${backendUrl}/api/investments/cumulative?lastMonths=6`);
        });
        await waitFor(() => expect(Chart).toHaveBeenCalled());

        const chartConfig = Chart.mock.calls[0][1];
        expect(chartConfig.data.datasets.map((dataset) => dataset.label)).toEqual(
            expect.arrayContaining(['401(k)', 'Stocks'])
        );
    });

    test('re-fetches when the selected period changes', async () => {
        render(<InvestmentGraph />);

        await waitFor(() => expect(global.fetch).toHaveBeenCalledTimes(1));

        fireEvent.click(screen.getByRole('button', { name: '12 months' }));

        await waitFor(() => {
            expect(global.fetch).toHaveBeenLastCalledWith(`${backendUrl}/api/investments/cumulative?lastMonths=12`);
        });
    });

    test('switches to a portfolio mix chart', async () => {
        render(<InvestmentGraph />);
        await waitFor(() => expect(Chart).toHaveBeenCalled());

        fireEvent.click(screen.getByRole('button', { name: 'Portfolio mix' }));

        await waitFor(() => expect(Chart.mock.calls[Chart.mock.calls.length - 1][1].type).toBe('doughnut'));
        expect(screen.getByRole('button', { name: 'Portfolio mix' })).toHaveAttribute('aria-pressed', 'true');
    });

    test('reports fetch failures', async () => {
        const consoleError = jest.spyOn(console, 'error').mockImplementation(() => {});
        global.fetch.mockRejectedValue(new Error('network'));

        render(<InvestmentGraph />);

        await waitFor(() => expect(consoleError).toHaveBeenCalledWith('Error fetching data:', expect.any(Error)));
        consoleError.mockRestore();
    });
});
