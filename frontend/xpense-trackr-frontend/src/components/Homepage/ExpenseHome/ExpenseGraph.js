import React, { useEffect, useMemo, useState } from 'react';
import backendUrl from '../../../config';
import ChartVisualization from '../../Charts/ChartVisualization';
import cumulativeChartData from '../../../utils/cumulativeChartData';
import './ExpenseHome.css';

const periods = [3, 6, 12];
const views = [
    { id: 'trend', label: 'Trend' },
    { id: 'monthly', label: 'Monthly' },
    { id: 'mix', label: 'Category mix' },
];

const ExpenseGraph = () => {
    const [rows, setRows] = useState([]);
    const [lastMonths, setLastMonths] = useState(6);
    const [view, setView] = useState('trend');
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [retryKey, setRetryKey] = useState(0);

    useEffect(() => {
        let active = true;
        const fetchData = async () => {
            setLoading(true);
            setError('');
            try {
                const response = await fetch(`${backendUrl}/api/expenses/cumulative?lastMonths=${lastMonths}`);
                if (!response.ok) throw new Error(`Request failed with status ${response.status}`);
                const data = await response.json();
                if (active) setRows(data || []);
            } catch (requestError) {
                console.error('Error fetching data:', requestError);
                if (active) setError('We couldn’t load your spending chart. Check your connection and try again.');
            } finally {
                if (active) setLoading(false);
            }
        };
        fetchData();
        return () => { active = false; };
    }, [lastMonths, retryKey]);

    const chart = useMemo(
        () => cumulativeChartData(rows, 'category', (category) => category, view),
        [rows, view]
    );

    return (
        <section className="finance-graph" aria-label="Spending chart">
            <div className="graph-controls">
                <div className="chart-mode-control" role="group" aria-label="Chart display">
                    {views.map((item) => (
                        <button
                            key={item.id}
                            type="button"
                            aria-pressed={view === item.id}
                            className={view === item.id ? 'selected' : ''}
                            onClick={() => setView(item.id)}
                        >
                            {item.label}
                        </button>
                    ))}
                </div>
                <div className="chart-period-control" role="group" aria-label="Chart time period">
                    {periods.map((months) => (
                        <button
                            key={months}
                            type="button"
                            aria-label={`${months} months`}
                            aria-pressed={lastMonths === months}
                            className={lastMonths === months ? 'selected' : ''}
                            onClick={() => setLastMonths(months)}
                        >
                            {months}M
                        </button>
                    ))}
                </div>
            </div>
            {loading ? (
                <p className="chart-message" role="status">Loading spending history…</p>
            ) : error ? (
                <div className="chart-message chart-error" role="alert">
                    <span>{error}</span>
                    <button type="button" onClick={() => setRetryKey((key) => key + 1)}>Try again</button>
                </div>
            ) : rows.length === 0 ? (
                <p className="chart-message">Add expenses to see your spending history here.</p>
            ) : (
                <ChartVisualization
                    type={chart.type}
                    data={chart.data}
                    ariaLabel={`${views.find((item) => item.id === view).label} of spending over the last ${lastMonths} months`}
                />
            )}
        </section>
    );
};

export default ExpenseGraph;
