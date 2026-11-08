const seriesColors = ['#2563eb', '#06b6d4', '#8b5cf6', '#f59e0b', '#ec4899', '#0f766e'];

const monthLabel = (month) => {
    const [year, monthNumber] = month.split('-').map(Number);
    return new Date(year, monthNumber - 1, 1).toLocaleDateString('en-US', {
        month: 'short',
        year: '2-digit',
    });
};

export default function cumulativeChartData(rows, seriesKey, labelForSeries, view) {
    const months = [...new Set(rows.map((row) => row.month))].sort();
    const series = [...new Set(rows.map((row) => row[seriesKey]))];
    const totalsBySeries = new Map(series.map((name) => [name, new Map()]));

    rows.forEach((row) => {
        totalsBySeries.get(row[seriesKey]).set(row.month, Number(row.amount));
    });

    const trendDatasets = series.map((name, index) => {
        const totals = totalsBySeries.get(name);
        let currentTotal = 0;
        let previousTotal = 0;
        const values = months.map((month) => {
            if (totals.has(month)) currentTotal = totals.get(month);
            const value = view === 'monthly' ? currentTotal - previousTotal : currentTotal;
            previousTotal = currentTotal;
            return value;
        });
        const color = seriesColors[index % seriesColors.length];
        return {
            label: labelForSeries(name),
            data: values,
            borderColor: color,
            backgroundColor: `${color}24`,
            ...(view === 'monthly' && { stack: 'monthly' }),
        };
    });

    if (view === 'mix') {
        const values = trendDatasets.map((dataset) => dataset.data[dataset.data.length - 1] || 0);
        return {
            type: 'doughnut',
            data: {
                labels: trendDatasets.map((dataset) => dataset.label),
                datasets: [{
                    label: 'Total',
                    data: values,
                    backgroundColor: values.map((_, index) => seriesColors[index % seriesColors.length]),
                    borderColor: '#ffffff',
                    borderWidth: 3,
                }],
            },
        };
    }

    return {
        type: view === 'monthly' ? 'bar' : 'line',
        data: {
            labels: months.map(monthLabel),
            datasets: trendDatasets,
        },
    };
}
