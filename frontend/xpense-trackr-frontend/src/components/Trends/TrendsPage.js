import React, { useEffect, useMemo, useState } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import ChartVisualization from '../Charts/ChartVisualization';
import formatCurrency from '../../utils/formatCurrency';

const periods = [3, 6, 12];
const monthLabel = (month) => {
    const [year, monthNumber] = month.split('-').map(Number);
    return new Date(year, monthNumber - 1, 1).toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
};

export default function TrendsPage() {
    const [months, setMonths] = useState(6);
    const [points, setPoints] = useState([]);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        let active = true;
        const loadTrends = async () => {
            setLoading(true);
            setError('');
            try {
                const response = await axios.get(`${backendUrl}/api/trends/spending?months=${months}`);
                if (active) setPoints(response.data || []);
            } catch (requestError) {
                console.error('Error fetching spending trends:', requestError);
                if (active) setError('We couldn’t load your spending trends. Check your connection and try again.');
            } finally {
                if (active) setLoading(false);
            }
        };
        loadTrends();
        return () => { active = false; };
    }, [months]);

    const totals = useMemo(() => {
        const byMonth = {};
        const byCategory = {};
        points.forEach(({ month, category, amount }) => {
            byMonth[month] = (byMonth[month] || 0) + Number(amount);
            byCategory[category] = (byCategory[category] || 0) + Number(amount);
        });
        return {
            months: Object.entries(byMonth).sort(([a], [b]) => a.localeCompare(b)),
            categories: Object.entries(byCategory).sort(([, a], [, b]) => b - a),
        };
    }, [points]);

    const totalSpent = totals.months.reduce((sum, [, amount]) => sum + amount, 0);
    const average = totals.months.length ? totalSpent / totals.months.length : 0;
    const highestCategory = totals.categories[0];
    const monthlyChartData = {
        labels: totals.months.map(([month]) => monthLabel(month)),
        datasets: [{
            label: 'Spending',
            data: totals.months.map(([, amount]) => amount),
            backgroundColor: '#4b83ed',
            borderRadius: 7,
            maxBarThickness: 42,
        }],
    };
    const categoryChartData = {
        labels: totals.categories.map(([category]) => category),
        datasets: [{
            label: 'Spending',
            data: totals.categories.map(([, amount]) => amount),
            backgroundColor: totals.categories.map((_, index) => ['#2563eb', '#06b6d4', '#8b5cf6', '#f59e0b', '#ec4899', '#0f766e'][index % 6]),
            borderColor: '#ffffff',
            borderWidth: 3,
        }],
    };

    return (
        <main className="management-page">
            <header className="page-intro">
                <div>
                    <p className="eyebrow">UNDERSTAND YOUR HABITS</p>
                    <h1>Spending trends</h1>
                    <p className="page-description">See where your money goes and how your spending changes over time.</p>
                </div>
                <div className="period-picker">
                    <span className="period-label">Time period</span>
                    <div className="period-control" role="group" aria-label="Time period">
                        {periods.map((period) => (
                            <button
                                key={period}
                                type="button"
                                aria-label={`${period} months`}
                                aria-pressed={months === period}
                                className={months === period ? 'selected' : ''}
                                onClick={() => setMonths(period)}
                            >
                                {period}M
                            </button>
                        ))}
                    </div>
                </div>
            </header>

            {error && <p className="inline-message" role="alert">{error}</p>}
            {loading ? <p className="loading-message" role="status">Loading spending trends…</p> : points.length === 0 ? (
                <section className="empty-state">
                    <span className="empty-state-icon" aria-hidden="true">↗</span>
                    <h2>Your trends will appear here</h2>
                    <p>Add a few expenses to see monthly patterns and category breakdowns.</p>
                </section>
            ) : (
                <>
                    <section className="summary-cards" aria-label="Spending summary">
                        <article className="summary-card">
                            <span>Total in selected period</span>
                            <strong>{formatCurrency(totalSpent)}</strong>
                        </article>
                        <article className="summary-card">
                            <span>Average per active month</span>
                            <strong>{formatCurrency(average)}</strong>
                        </article>
                        <article className="summary-card">
                            <span>Top category</span>
                            <strong>{highestCategory[0]}</strong>
                            <small>{formatCurrency(highestCategory[1])} in the selected period</small>
                        </article>
                    </section>
                    <div className="trend-grid">
                        <section className="data-card">
                            <div className="card-heading">
                                <div><p className="eyebrow">MONTH BY MONTH</p><h2>Monthly spending</h2></div>
                            </div>
                            <ChartVisualization
                                type="bar"
                                data={monthlyChartData}
                                ariaLabel="Bar chart showing spending by month"
                            />
                        </section>
                        <section className="data-card category-chart-card">
                            <div className="card-heading">
                                <div><p className="eyebrow">WHERE IT GOES</p><h2>By category</h2></div>
                            </div>
                            <ChartVisualization
                                type="doughnut"
                                data={categoryChartData}
                                ariaLabel="Doughnut chart showing spending distribution by category"
                            />
                        </section>
                    </div>
                </>
            )}
        </main>
    );
}
