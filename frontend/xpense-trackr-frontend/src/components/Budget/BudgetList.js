import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import backendUrl from '../../config';
import formatCurrency from '../../utils/formatCurrency';

const now = new Date();
const currentMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;

export default function BudgetList() {
    const [budgets, setBudgets] = useState([]);
    const [categories, setCategories] = useState([]);
    const [expenses, setExpenses] = useState([]);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(true);

    const loadBudgets = useCallback(async () => {
        setLoading(true);
        setError('');
        try {
            const [budgetResponse, categoryResponse, expenseResponse] = await Promise.all([
                axios.get(`${backendUrl}/api/budgets`),
                axios.get(`${backendUrl}/api/categories`),
                axios.get(`${backendUrl}/api/expenses`),
            ]);
            setBudgets(budgetResponse.data || []);
            setCategories(categoryResponse.data || []);
            setExpenses(expenseResponse.data || []);
        } catch (requestError) {
            console.error('Error fetching budgets:', requestError);
            setError('We couldn’t load your budgets. Check your connection and try again.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadBudgets();
    }, [loadBudgets]);

    const categoryName = (id) => categories.find((category) => Number(category.id) === Number(id))?.name || 'Unknown category';
    const spentByCategory = expenses.reduce((totals, expense) => {
        if (String(expense.date).slice(0, 7) === currentMonth) {
            const key = Number(expense.category_id);
            totals[key] = (totals[key] || 0) + Number(expense.amount);
        }
        return totals;
    }, {});
    const totalBudget = budgets.reduce((sum, budget) => sum + Number(budget.amount), 0);
    const totalSpent = budgets.reduce((sum, budget) => sum + (spentByCategory[Number(budget.category_id)] || 0), 0);

    const deleteBudget = async (id) => {
        if (!window.confirm('Delete this monthly budget?')) return;
        try {
            await axios.delete(`${backendUrl}/api/budgets/${id}`);
            await loadBudgets();
        } catch (requestError) {
            console.error('Error deleting budget:', requestError);
            setError('We couldn’t delete that budget. Please try again.');
        }
    };

    return (
        <main className="management-page">
            <header className="page-intro">
                <div>
                    <p className="eyebrow">PLAN WITH CONFIDENCE</p>
                    <h1>Monthly budgets</h1>
                    <p className="page-description">Set a limit for each category and keep an eye on this month’s spending.</p>
                </div>
                <Link className="primary-action" to="/budgets/new">Create a budget</Link>
            </header>

            <section className="summary-cards" aria-label="Budget summary">
                <article className="summary-card">
                    <span>Planned this month</span>
                    <strong>{formatCurrency(totalBudget)}</strong>
                </article>
                <article className="summary-card">
                    <span>Spent in budgeted categories</span>
                    <strong>{formatCurrency(totalSpent)}</strong>
                </article>
                <article className="summary-card">
                    <span>Remaining</span>
                    <strong>{formatCurrency(Math.max(0, totalBudget - totalSpent))}</strong>
                </article>
            </section>

            {error && <p className="inline-message" role="alert">{error}</p>}
            {loading ? <p className="loading-message" role="status">Loading budgets…</p> : (
                budgets.length === 0 ? (
                    <section className="empty-state">
                        <span className="empty-state-icon" aria-hidden="true">$</span>
                        <h2>Start with a monthly plan</h2>
                        <p>Create a budget for a category to see how your spending tracks against it.</p>
                        <Link className="primary-action" to="/budgets/new">Create your first budget</Link>
                    </section>
                ) : (
                    <section className="data-card">
                        <div className="card-heading">
                            <div>
                                <p className="eyebrow">YOUR PLAN</p>
                                <h2>Category budgets</h2>
                            </div>
                            <span className="muted-copy">{budgets.length} {budgets.length === 1 ? 'category' : 'categories'}</span>
                        </div>
                        <div className="budget-list">
                            {budgets.map((budget) => {
                                const spent = spentByCategory[Number(budget.category_id)] || 0;
                                const limit = Number(budget.amount);
                                const percent = limit > 0 ? (spent / limit) * 100 : 0;
                                const overBudget = percent >= 100;
                                return (
                                    <article className="budget-row" key={budget.id}>
                                        <div className="budget-row-heading">
                                            <div>
                                                <h3>{categoryName(budget.category_id)}</h3>
                                                <p>{formatCurrency(spent)} spent <span aria-hidden="true">·</span> {formatCurrency(limit)} monthly limit</p>
                                            </div>
                                            <div className="row-actions">
                                                <strong className={overBudget ? 'over-budget' : ''}>{formatCurrency(Math.max(0, limit - spent))} left</strong>
                                                <Link to={`/budgets/${budget.id}/edit`} aria-label={`Edit ${categoryName(budget.category_id)} budget`}>Edit</Link>
                                                <button type="button" className="text-button danger-text" onClick={() => deleteBudget(budget.id)} aria-label={`Delete ${categoryName(budget.category_id)} budget`}>Delete</button>
                                            </div>
                                        </div>
                                        <div
                                            className={`progress-track ${overBudget ? 'is-over' : ''}`}
                                            role="progressbar"
                                            aria-label={`${categoryName(budget.category_id)} budget used`}
                                            aria-valuemin="0"
                                            aria-valuemax="100"
                                            aria-valuenow={Math.min(100, Math.round(percent))}
                                        >
                                            <span style={{ width: `${Math.min(100, percent)}%` }} />
                                        </div>
                                    </article>
                                );
                            })}
                        </div>
                    </section>
                )
            )}
        </main>
    );
}