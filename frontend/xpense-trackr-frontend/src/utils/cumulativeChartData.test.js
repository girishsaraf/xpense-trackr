import cumulativeChartData from './cumulativeChartData';

const rows = [
    { month: '2026-01', category: 'Food', amount: 100 },
    { month: '2026-01', category: 'Travel', amount: 50 },
    { month: '2026-02', category: 'Food', amount: 125 },
    { month: '2026-03', category: 'Food', amount: 175 },
    { month: '2026-03', category: 'Travel', amount: 80 },
];

describe('cumulativeChartData', () => {
    test('converts cumulative totals into monthly amounts with zeroes for missing categories', () => {
        const chart = cumulativeChartData(rows, 'category', (name) => name, 'monthly');

        expect(chart.type).toBe('bar');
        expect(chart.data.datasets).toEqual(expect.arrayContaining([
            expect.objectContaining({ label: 'Food', data: [100, 25, 50] }),
            expect.objectContaining({ label: 'Travel', data: [50, 0, 30] }),
        ]));
    });

    test('creates a category mix from the latest cumulative totals', () => {
        const chart = cumulativeChartData(rows, 'category', (name) => name, 'mix');

        expect(chart.type).toBe('doughnut');
        expect(chart.data.labels).toEqual(['Food', 'Travel']);
        expect(chart.data.datasets[0].data).toEqual([175, 80]);
    });
});
