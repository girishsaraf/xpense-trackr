import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import { formatAmountMinor, parseAmountMinor } from '../../utils/money';
import './PlanningPage.css';

const today = () => new Date().toISOString().slice(0, 10);
const currentPeriod = () => new Date().toISOString().slice(0, 7);

function PlanningPage() {
    const [goals, setGoals] = useState([]);
    const [envelopes, setEnvelopes] = useState([]);
    const [categories, setCategories] = useState([]);
    const [period, setPeriod] = useState(currentPeriod);
    const [error, setError] = useState('');
    const [goalForm, setGoalForm] = useState({ name: '', currency: 'USD', target: '', target_date: '' });
    const [contributions, setContributions] = useState({});
    const [contributionHistory, setContributionHistory] = useState({});
    const [expandedGoal, setExpandedGoal] = useState(null);
    const [envelopeForm, setEnvelopeForm] = useState({ category_id: '', currency: 'USD', amount: '', rollover: false });

    const loadData = useCallback(async () => {
        setError('');
        try {
            const [goalResponse, envelopeResponse, categoryResponse] = await Promise.all([
                axios.get(`${backendUrl}/api/planning/goals`),
                axios.get(`${backendUrl}/api/planning/envelopes`, { params: { period } }),
                axios.get(`${backendUrl}/api/categories`),
            ]);
            setGoals(goalResponse.data || []);
            setEnvelopes(envelopeResponse.data || []);
            setCategories(categoryResponse.data || []);
        } catch (requestError) {
            console.error('Error loading financial plans:', requestError);
            setError('We couldn’t load your goals and envelopes. Check the connection and try again.');
        }
    }, [period]);

    useEffect(() => {
        loadData();
    }, [loadData]);

    const createGoal = async (event) => {
        event.preventDefault();
        const targetAmountMinor = parseAmountMinor(goalForm.target, goalForm.currency);
        if (!goalForm.name.trim() || targetAmountMinor === null || targetAmountMinor <= 0) {
            setError('Enter a goal name and a valid target amount.');
            return;
        }
        try {
            await axios.post(`${backendUrl}/api/planning/goals`, {
                name: goalForm.name,
                currency: goalForm.currency,
                target_amount_minor: targetAmountMinor,
                target_date: goalForm.target_date,
            });
            setGoalForm({ name: '', currency: 'USD', target: '', target_date: '' });
            await loadData();
        } catch (requestError) {
            console.error('Error creating financial goal:', requestError);
            setError('We couldn’t save this goal. Check the details and try again.');
        }
    };

    const createContribution = async (event, goalId) => {
        event.preventDefault();
        const form = contributions[goalId] || { amount: '', direction: 'contribution', date: today() };
        const amountMinor = parseAmountMinor(form.amount, goals.find((goal) => goal.id === goalId)?.currency || 'USD');
        if (amountMinor === null || amountMinor <= 0) {
            setError('Enter a valid contribution or withdrawal amount.');
            return;
        }
        try {
            await axios.post(`${backendUrl}/api/planning/goals/${goalId}/contributions`, {
                amount_minor: amountMinor,
                direction: form.direction,
                date: form.date,
                note: form.note || '',
            });
            setContributions((current) => ({ ...current, [goalId]: { amount: '', direction: 'contribution', date: today() } }));
            await loadData();
        } catch (requestError) {
            console.error('Error recording goal contribution:', requestError);
            setError(requestError.response?.status === 409
                ? 'That withdrawal is larger than the goal’s current balance.'
                : 'We couldn’t record this contribution. Check the details and try again.');
        }
    };

    const createEnvelope = async (event) => {
        event.preventDefault();
        const amountMinor = parseAmountMinor(envelopeForm.amount, envelopeForm.currency);
        if (!envelopeForm.category_id || amountMinor === null || amountMinor <= 0) {
            setError('Choose a category and enter a valid monthly amount.');
            return;
        }
        try {
            await axios.post(`${backendUrl}/api/planning/envelopes`, {
                category_id: Number(envelopeForm.category_id),
                currency: envelopeForm.currency,
                monthly_limit_minor: amountMinor,
                rollover: envelopeForm.rollover,
            });
            setEnvelopeForm({ category_id: '', currency: 'USD', amount: '', rollover: false });
            await loadData();
        } catch (requestError) {
            console.error('Error creating budget envelope:', requestError);
            setError('We couldn’t save this envelope. Check its category and try again.');
        }
    };

    const deleteEnvelope = async (id) => {
        try {
            await axios.delete(`${backendUrl}/api/planning/envelopes/${id}`);
            await loadData();
        } catch (requestError) {
            console.error('Error deleting budget envelope:', requestError);
            setError('We couldn’t remove this envelope. Try again.');
        }
    };

    const toggleContributionHistory = async (goalId) => {
        if (expandedGoal === goalId) {
            setExpandedGoal(null);
            return;
        }
        try {
            const response = await axios.get(`${backendUrl}/api/planning/goals/${goalId}/contributions`);
            setContributionHistory((current) => ({ ...current, [goalId]: response.data || [] }));
            setExpandedGoal(goalId);
        } catch (requestError) {
            console.error('Error loading goal contribution history:', requestError);
            setError('We couldn’t load this goal’s contribution history.');
        }
    };

    return (
        <main className="planning-page" aria-labelledby="planning-title">
            <header className="planning-heading">
                <div>
                    <p className="planning-eyebrow">PLAN WITH PURPOSE</p>
                    <h1 id="planning-title">Goals &amp; envelopes</h1>
                    <p>Track progress toward savings goals and give monthly spending a clear limit.</p>
                </div>
                <label>
                    Budget period
                    <input type="month" value={period} onChange={(event) => setPeriod(event.target.value)} />
                </label>
            </header>
            {error && <p className="planning-error" role="alert">{error}</p>}

            <section className="planning-section" aria-labelledby="goals-title">
                <div className="planning-section-heading">
                    <div><p className="planning-eyebrow">SAVING MILESTONES</p><h3 id="goals-title">Your goals</h3></div>
                </div>
                <form className="planning-form" onSubmit={createGoal}>
                    <label>Goal name<input value={goalForm.name} onChange={(event) => setGoalForm({ ...goalForm, name: event.target.value })} required /></label>
                    <label>Target amount<input type="number" min="0" step="0.01" value={goalForm.target} onChange={(event) => setGoalForm({ ...goalForm, target: event.target.value })} required /></label>
                    <label>Currency<input maxLength="3" value={goalForm.currency} onChange={(event) => setGoalForm({ ...goalForm, currency: event.target.value.toUpperCase() })} required /></label>
                    <label>Target date<input type="date" value={goalForm.target_date} onChange={(event) => setGoalForm({ ...goalForm, target_date: event.target.value })} /></label>
                    <button type="submit">Add goal</button>
                </form>
                {goals.length === 0 ? <p className="planning-empty">Add a savings goal to start tracking progress.</p> : (
                    <div className="planning-cards">
                        {goals.map((goal) => {
                            const progress = Math.min(100, Math.round((goal.current_amount_minor / goal.target_amount_minor) * 100));
                            const form = contributions[goal.id] || { amount: '', direction: 'contribution', date: today() };
                            return (
                                <article className="planning-card" key={goal.id}>
                                    <div className="planning-card-top">
                                        <div><h4>{goal.name}</h4><p>{formatAmountMinor(goal.current_amount_minor, goal.currency)} of {formatAmountMinor(goal.target_amount_minor, goal.currency)}</p></div>
                                        <strong>{progress}%</strong>
                                    </div>
                                    <progress value={progress} max="100" aria-label={`${goal.name} progress`} />
                                    {goal.required_monthly_minor > 0 && <p className="planning-note">About {formatAmountMinor(goal.required_monthly_minor, goal.currency)} per month to meet your target date.</p>}
                                    <button className="planning-secondary" type="button" onClick={() => toggleContributionHistory(goal.id)}>
                                        {expandedGoal === goal.id ? 'Hide history' : 'View history'}
                                    </button>
                                    {expandedGoal === goal.id && (
                                        <ul className="planning-history">
                                            {(contributionHistory[goal.id] || []).length === 0
                                                ? <li>No contributions recorded yet.</li>
                                                : contributionHistory[goal.id].map((entry) => (
                                                    <li key={entry.id}>
                                                        <span>{entry.date} · {entry.direction}</span>
                                                        <strong>{formatAmountMinor(entry.amount_minor, goal.currency)}</strong>
                                                    </li>
                                                ))}
                                        </ul>
                                    )}
                                    <form className="planning-inline-form" onSubmit={(event) => createContribution(event, goal.id)}>
                                        <label>Entry type<select value={form.direction} onChange={(event) => setContributions({ ...contributions, [goal.id]: { ...form, direction: event.target.value } })}><option value="contribution">Contribution</option><option value="withdrawal">Withdrawal</option></select></label>
                                        <label>Amount<input type="number" min="0" step="0.01" value={form.amount} onChange={(event) => setContributions({ ...contributions, [goal.id]: { ...form, amount: event.target.value } })} required /></label>
                                        <label>Date<input type="date" value={form.date} onChange={(event) => setContributions({ ...contributions, [goal.id]: { ...form, date: event.target.value } })} required /></label>
                                        <button type="submit">Record</button>
                                    </form>
                                </article>
                            );
                        })}
                    </div>
                )}
            </section>

            <section className="planning-section" aria-labelledby="envelopes-title">
                <div className="planning-section-heading">
                    <div><p className="planning-eyebrow">MONTHLY SPENDING PLAN</p><h3 id="envelopes-title">Budget envelopes</h3></div>
                </div>
                <form className="planning-form" onSubmit={createEnvelope}>
                    <label>Category<select value={envelopeForm.category_id} onChange={(event) => setEnvelopeForm({ ...envelopeForm, category_id: event.target.value })} required><option value="">Choose category</option>{categories.map((category) => <option value={category.id} key={category.id}>{category.name}</option>)}</select></label>
                    <label>Monthly limit<input type="number" min="0" step="0.01" value={envelopeForm.amount} onChange={(event) => setEnvelopeForm({ ...envelopeForm, amount: event.target.value })} required /></label>
                    <label>Currency<input maxLength="3" value={envelopeForm.currency} onChange={(event) => setEnvelopeForm({ ...envelopeForm, currency: event.target.value.toUpperCase() })} required /></label>
                    <label className="planning-checkbox"><input type="checkbox" checked={envelopeForm.rollover} onChange={(event) => setEnvelopeForm({ ...envelopeForm, rollover: event.target.checked })} /> Roll unused balance forward</label>
                    <button type="submit">Add envelope</button>
                </form>
                <p className="planning-note">Rollover carries each prior month’s unused or overspent amount forward. Actual spending comes from categorized ledger transactions.</p>
                {envelopes.length === 0 ? <p className="planning-empty">No envelopes for this period yet.</p> : (
                    <div className="planning-cards">
                        {envelopes.map((envelope) => {
                            const category = categories.find((item) => item.id === envelope.category_id);
                            return (
                                <article className="planning-card" key={envelope.id}>
                                    <div className="planning-card-top">
                                        <div><h4>{category?.name || 'Category'}</h4><p>{formatAmountMinor(envelope.spent_minor, envelope.currency)} spent of {formatAmountMinor(envelope.monthly_limit_minor + envelope.rollover_minor, envelope.currency)} available</p></div>
                                        <strong className={envelope.available_minor < 0 ? 'planning-over' : ''}>{formatAmountMinor(envelope.available_minor, envelope.currency)} left</strong>
                                    </div>
                                    <button className="planning-secondary" type="button" onClick={() => deleteEnvelope(envelope.id)} aria-label={`Remove ${category?.name || 'category'} envelope`}>Remove envelope</button>
                                </article>
                            );
                        })}
                    </div>
                )}
            </section>
        </main>
    );
}

export default PlanningPage;
