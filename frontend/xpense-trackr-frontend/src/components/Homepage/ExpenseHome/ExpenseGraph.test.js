import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import Chart from 'chart.js/auto';
import ExpenseGraph from './ExpenseGraph';
import backendUrl from '../../../config';

jest.mock('chart.js/auto', () => ({
    __esModule: true,
    default: jest.fn(function () {
        return { destroy: jest.fn() };
    }),
}));

describe('ExpenseGraph', () => {
    const originalFetch = global.fetch;
    const originalGetContext = HTMLCanvasElement.prototype.getContext;
    const graphData = [
        { month: '2026-01', category: 'Food', amount: 100 },
        { month: '2026-02', category: 'Food', amount: 120 },
        { month: '2026-02', category: 'Travel', amount: 50 },
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
        render(<ExpenseGraph />);

        await waitFor(() => {
            expect(global.fetch).toHaveBeenCalledWith(`${backendUrl}/api/expenses/cumulative?lastMonths=6`);
        });
        await waitFor(() => expect(Chart).toHaveBeenCalled());

        const chartConfig = Chart.mock.calls[0][1];
        expect(chartConfig.data.datasets).toHaveLength(2);
        expect(chartConfig.data.datasets.map((dataset) => dataset.label)).toEqual(expect.arrayContaining(['Food', 'Travel']));
    });

    test('re-fetches when the selected period changes', async () => {
        render(<ExpenseGraph />);

        await waitFor(() => expect(global.fetch).toHaveBeenCalledTimes(1));

        fireEvent.click(screen.getByRole('button', { name: '3 months' }));

        await waitFor(() => {
            expect(global.fetch).toHaveBeenLastCalledWith(`${backendUrl}/api/expenses/cumulative?lastMonths=3`);
        });
    });

    test('switches to a stacked monthly comparison view', async () => {
        render(<ExpenseGraph />);
        await waitFor(() => expect(Chart).toHaveBeenCalled());

        fireEvent.click(screen.getByRole('button', { name: 'Monthly' }));

        await waitFor(() => expect(Chart.mock.calls[Chart.mock.calls.length - 1][1].type).toBe('bar'));
        expect(screen.getByRole('button', { name: 'Monthly' })).toHaveAttribute('aria-pressed', 'true');
    });

    test('reports fetch failures', async () => {
        const consoleError = jest.spyOn(console, 'error').mockImplementation(() => {});
        global.fetch.mockRejectedValue(new Error('network'));

        render(<ExpenseGraph />);

        await waitFor(() => expect(consoleError).toHaveBeenCalledWith('Error fetching data:', expect.any(Error)));
        consoleError.mockRestore();
    });
});
