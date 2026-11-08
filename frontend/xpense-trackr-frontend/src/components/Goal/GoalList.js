import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import backendUrl from '../../config';
import formatCurrency from '../../utils/formatCurrency';

export default function GoalList() {
    const [goals, setGoals] = useState([]);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(true);

    const loadGoals = useCallback(async () => {
        setLoading(true);
        setError('');
        try {
            const response = await axios.get(`${backendUrl}/api/goals`);
            setGoals(response.data || []);
        } catch (requestError) {
            console.error('Error fetching savings goals:', requestError);
            setError('We couldn’t load your goals. Check your connection and try again.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadGoals();
    }, [loadGoals]);

    const deleteGoal = async (goal) => {
        if (!window.confirm(`Delete the “${goal.name}” goal?`)) return;
        try {
            await axios.delete(`${backendUrl}/api/goals/${goal.id}`);
            await loadGoals();
        } catch (requestError) {
            console.error('Error deleting savings goal:', requestError);
            setError('We couldn’t delete this goal. Please try again.');
        }
    };

    const totalSaved = goals.reduce((sum, goal) => sum + Number(goal.current_amount), 0);
    const totalTarget = goals.reduce((sum, goal) => sum + Number(goal.target_amount), 0);
    const formatTargetDate = (date) => {
        if (!date) return 'No target date';
        const [year, month, day] = date.split('-').map(Number);
        return new Date(year, month - 1, day).toLocaleDateString('en-US', { month: 'long', year: 'numeric' });
    };

    return (
        <main className="management-page">
            <header className="page-intro">
                <div>
                    <p className="eyebrow">MAKE ROOM FOR WHAT MATTERS</p>
                    <h1>Savings goals</h1>
                    <p className="page-description">Track the things you’re saving for and celebrate each step forward.</p>
                </div>
                <Link className="primary-action" to="/goals/new">Add a goal</Link>
            </header>

            <section className="summary-cards" aria-label="Savings goal summary">
                <article className="summary-card">
                    <span>Saved toward your goals</span>
                    <strong>{formatCurrency(totalSaved)}</strong>
                </article>
                <article className="summary-card">
                    <span>Total goal targets</span>
                    <strong>{formatCurrency(totalTarget)}</strong>
                </article>
                <article className="summary-card">
                    <span>Goals in progress</span>
                    <strong>{goals.filter((goal) => Number(goal.current_amount) < Number(goal.target_amount)).length}</strong>
                </article>
            </section>

            {error && <p className="inline-message" role="alert">{error}</p>}
            {loading ? <p className="loading-message" role="status">Loading your goals…</p> : (
                goals.length === 0 ? (
                    <section className="empty-state">
                        <span className="empty-state-icon" aria-hidden="true">◎</span>
                        <h2>Give your savings a destination</h2>
                        <p>Create a goal for an emergency fund, a trip, or anything you’re looking forward to.</p>
                        <Link className="primary-action" to="/goals/new">Create your first goal</Link>
                    </section>
                ) : (
                    <section className="goal-grid" aria-label="Your savings goals">
                        {goals.map((goal) => {
                            const target = Number(goal.target_amount);
                            const saved = Number(goal.current_amount);
                            const percent = target > 0 ? (saved / target) * 100 : 0;
                            const complete = percent >= 100;
                            return (
                                <article className="goal-card" key={goal.id}>
                                    <div className="goal-card-top">
                                        <span className={`goal-status ${complete ? 'complete' : ''}`}>{complete ? 'Goal reached' : 'In progress'}</span>
                                        <div className="row-actions">
                                            <Link to={`/goals/${goal.id}/edit`} aria-label={`Edit ${goal.name} goal`}>Edit</Link>
                                            <button type="button" className="text-button danger-text" onClick={() => deleteGoal(goal)} aria-label={`Delete ${goal.name} goal`}>Delete</button>
                                        </div>
                                    </div>
                                    <h2>{goal.name}</h2>
                                    <p className="goal-target-date">{formatTargetDate(goal.target_date)}</p>
                                    <div className="goal-amounts">
                                        <strong>{formatCurrency(saved)}</strong>
                                        <span>of {formatCurrency(target)}</span>
                                    </div>
                                    <div
                                        className="progress-track"
                                        role="progressbar"
                                        aria-label={`${goal.name} progress`}
                                        aria-valuemin="0"
                                        aria-valuemax="100"
                                        aria-valuenow={Math.min(100, Math.round(percent))}
                                    >
                                        <span style={{ width: `${Math.min(100, percent)}%` }} />
                                    </div>
                                    <p className="goal-progress-copy">{complete ? 'You’ve reached this goal. Congratulations!' : `${Math.min(100, Math.round(percent))}% of the way there`}</p>
                                </article>
                            );
                        })}
                    </section>
                )
            )}
        </main>
    );
}
