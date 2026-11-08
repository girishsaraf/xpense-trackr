import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import backendUrl from '../../config';
import formatCurrency from '../../utils/formatCurrency';

const formatDate = (date) => {
    if (!date || String(date).startsWith('0001-01-01')) return 'No end date';
    const parsedDate = new Date(date);
    return Number.isNaN(parsedDate.getTime()) ? 'No end date' : parsedDate.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
};

export default function RecurringPage() {
    const [view, setView] = useState('expenses');
    const [expenses, setExpenses] = useState([]);
    const [investments, setInvestments] = useState([]);
    const [categories, setCategories] = useState([]);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(true);

    const loadSchedules = useCallback(async () => {
        setLoading(true);
        setError('');
        try {
            const [expenseResponse, investmentResponse, categoryResponse] = await Promise.all([
                axios.get(`${backendUrl}/api/recurring/expenses`),
                axios.get(`${backendUrl}/api/recurring/investments`),
                axios.get(`${backendUrl}/api/categories`),
            ]);
            setExpenses(expenseResponse.data || []);
            setInvestments(investmentResponse.data || []);
            setCategories(categoryResponse.data || []);
        } catch (requestError) {
            console.error('Error fetching recurring schedules:', requestError);
            setError('We couldn’t load your recurring activity. Check your connection and try again.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadSchedules();
    }, [loadSchedules]);

    const deleteSchedule = async (item, kind) => {
        const label = kind === 'expenses' ? 'recurring expense' : 'recurring investment';
        if (!window.confirm(`Stop this ${label}?`)) return;
        try {
            await axios.delete(`${backendUrl}/api/recurring/${kind}/${item.id}`);
            await loadSchedules();
        } catch (requestError) {
            console.error(`Error deleting ${label}:`, requestError);
            setError(`We couldn’t stop this ${label}. Please try again.`);
        }
    };

    const entries = view === 'expenses' ? expenses : investments;
    const monthlyEstimate = (amount, frequency) => {
        const monthlyFactor = { daily: 30, weekly: 4.33, biweekly: 2.17, monthly: 1 }[frequency] || 1;
        return Number(amount) * monthlyFactor;
    };
    const nameFor = (item) => view === 'expenses'
        ? categories.find((category) => Number(category.id) === Number(item.category_id))?.name || 'Expense'
        : item.investment_type;

    return (
        <main className="management-page">
            <header className="page-intro">
                <div>
                    <p className="eyebrow">AUTOMATIC, NOT FORGOTTEN</p>
                    <h1>Recurring activity</h1>
                    <p className="page-description">Review scheduled expenses and investment contributions in one place.</p>
                </div>
            </header>
            <div className="view-switch recurring-switch" aria-label="Recurring activity type">
                <button type="button" aria-pressed={view === 'expenses'} className={`toggle-button ${view === 'expenses' ? 'active' : ''}`} onClick={() => setView('expenses')}>Expenses ({expenses.length})</button>
                <button type="button" aria-pressed={view === 'investments'} className={`toggle-button ${view === 'investments' ? 'active' : ''}`} onClick={() => setView('investments')}>Investments ({investments.length})</button>
            </div>
            {error && <p className="inline-message" role="alert">{error}</p>}
            {loading ? <p className="loading-message" role="status">Loading recurring activity…</p> : entries.length === 0 ? (
                <section className="empty-state">
                    <span className="empty-state-icon" aria-hidden="true">↻</span>
                    <h2>No recurring {view} yet</h2>
                    <p>When you add a recurring {view === 'expenses' ? 'expense' : 'investment'}, you can review or stop it here.</p>
                </section>
            ) : (
                <section className="data-card">
                    <div className="card-heading">
                        <div><p className="eyebrow">SCHEDULED</p><h2>{view === 'expenses' ? 'Recurring expenses' : 'Recurring investments'}</h2></div>
                        <span className="muted-copy">Approx. {formatCurrency(entries.reduce((sum, item) => sum + monthlyEstimate(item.amount, item.frequency), 0))} / month</span>
                    </div>
                    <div className="recurring-list">
                        {entries.map((item) => (
                            <article className="recurring-row" key={item.id}>
                                <div className="recurring-details">
                                    <strong>{item.description || nameFor(item)}</strong>
                                    <span>{nameFor(item)} <span aria-hidden="true">·</span> {item.frequency} <span aria-hidden="true">·</span> Through {formatDate(item.end_date)}</span>
                                </div>
                                <div className="row-actions">
                                    <strong>{formatCurrency(item.amount)}</strong>
                                    <Link to={`/recurring/${view}/${item.id}/edit`} aria-label={`Edit recurring ${view === 'expenses' ? 'expense' : 'investment'} ${item.description || nameFor(item)}`}>Edit</Link>
                                    <button className="text-button danger-text" type="button" onClick={() => deleteSchedule(item, view)} aria-label={`Stop recurring ${view === 'expenses' ? 'expense' : 'investment'} ${item.description || nameFor(item)}`}>Stop</button>
                                </div>
                            </article>
                        ))}
                    </div>
                </section>
            )}
        </main>
    );
}
