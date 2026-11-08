import React, { useEffect, useMemo, useState } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import formatCurrency from '../../utils/formatCurrency';

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
    const maxCategory = highestCategory?.[1] || 1;
    const maxMonth = Math.max(1, ...totals.months.map(([, amount]) => amount));

    return (
        <main className="management-page">
            <header className="page-intro">
                <div>
                    <p className="eyebrow">UNDERSTAND YOUR HABITS</p>
                    <h1>Spending trends</h1>
                    <p className="page-description">See where your money goes and how your spending changes over time.</p>
                </div>
                <div className="period-picker">
                    <label htmlFor="trend-period">Time period</label>
                    <select id="trend-period" value={months} onChange={(event) => setMonths(Number(event.target.value))}>
                        <option value="3">Last 3 months</option>
                        <option value="6">Last 6 months</option>
                        <option value="12">Last 12 months</option>
                    </select>
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
                            <div className="trend-bars" role="list" aria-label="Monthly spending amounts">
                                {totals.months.map(([month, amount]) => (
                                    <div className="trend-bar-row" role="listitem" key={month}>
                                        <span>{monthLabel(month)}</span>
                                        <div className="trend-bar-track" aria-hidden="true"><span style={{ width: `${(amount / maxMonth) * 100}%` }} /></div>
                                        <strong>{formatCurrency(amount)}</strong>
                                    </div>
                                ))}
                            </div>
                        </section>
                        <section className="data-card">
                            <div className="card-heading">
                                <div><p className="eyebrow">WHERE IT GOES</p><h2>By category</h2></div>
                            </div>
                            <div className="category-breakdown" role="list" aria-label="Spending by category">
                                {totals.categories.map(([category, amount]) => (
                                    <div className="category-breakdown-row" role="listitem" key={category}>
                                        <div className="breakdown-heading"><span>{category}</span><strong>{formatCurrency(amount)}</strong></div>
                                        <div className="trend-bar-track" aria-hidden="true"><span style={{ width: `${(amount / maxCategory) * 100}%` }} /></div>
                                    </div>
                                ))}
                            </div>
                        </section>
                    </div>
                </>
            )}
        </main>
    );
}
