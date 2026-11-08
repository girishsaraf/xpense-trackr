import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link, useNavigate, useParams } from 'react-router-dom';
import backendUrl from '../../config';

export default function BudgetForm() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [categories, setCategories] = useState([]);
    const [existingBudgets, setExistingBudgets] = useState([]);
    const [formData, setFormData] = useState({ category_id: '', amount: '' });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(Boolean(id));

    useEffect(() => {
        const loadForm = async () => {
            try {
                const requests = [
                    axios.get(`${backendUrl}/api/categories`),
                    axios.get(`${backendUrl}/api/budgets`),
                ];
                if (id) requests.push(axios.get(`${backendUrl}/api/budgets/${id}`));
                const [categoryResponse, budgetsResponse, budgetResponse] = await Promise.all(requests);
                setCategories(categoryResponse.data || []);
                setExistingBudgets(budgetsResponse.data || []);
                if (budgetResponse) {
                    setFormData({
                        category_id: String(budgetResponse.data.category_id),
                        amount: String(budgetResponse.data.amount),
                    });
                }
            } catch (requestError) {
                console.error('Error loading budget form:', requestError);
                setError('We couldn’t load the budget details. Please return to budgets and try again.');
            } finally {
                setLoading(false);
            }
        };
        loadForm();
    }, [id]);

    const submit = async (event) => {
        event.preventDefault();
        setError('');
        const amount = Number(formData.amount);
        if (!formData.category_id || !Number.isFinite(amount) || amount <= 0) {
            setError('Choose a category and enter an amount greater than zero.');
            return;
        }
        const payload = { category_id: Number(formData.category_id), amount };
        try {
            if (id) await axios.put(`${backendUrl}/api/budgets/${id}`, payload);
            else await axios.post(`${backendUrl}/api/budgets`, payload);
            navigate('/budgets');
        } catch (requestError) {
            console.error('Error saving budget:', requestError);
            setError('We couldn’t save this budget. Please try again.');
        }
    };

    const availableCategories = categories.filter((category) =>
        id
            ? Number(category.id) === Number(formData.category_id) ||
                !existingBudgets.some((budget) => Number(budget.category_id) === Number(category.id))
            : !existingBudgets.some((budget) => Number(budget.category_id) === Number(category.id))
    );

    if (loading) return <p className="loading-message" role="status">Loading budget…</p>;

    return (
        <main className="management-page form-page">
            <Link className="back-link" to="/budgets">← Back to budgets</Link>
            <section className="form-card">
                <p className="eyebrow">MAKE A PLAN</p>
                <h1>{id ? 'Edit budget' : 'Create a budget'}</h1>
                <p className="form-description">Choose a spending category and set the monthly amount you want to stay within.</p>
                {error && <p className="inline-message" role="alert">{error}</p>}
                {!id && categories.length > 0 && availableCategories.length === 0 && (
                    <p className="inline-message" role="status">Every category already has a budget. Add another category to create a new budget.</p>
                )}
                <form onSubmit={submit} noValidate>
                    <div className="field-group">
                        <label htmlFor="budget-category">Category</label>
                        <select
                            id="budget-category"
                            value={formData.category_id}
                            onChange={(event) => setFormData({ ...formData, category_id: event.target.value })}
                            required
                        >
                            <option value="">Select a category</option>
                            {availableCategories.map((category) => <option key={category.id} value={category.id}>{category.name}</option>)}
                        </select>
                    </div>
                    <div className="field-group">
                        <label htmlFor="budget-amount">Monthly limit</label>
                        <div className="currency-input">
                            <span aria-hidden="true">$</span>
                            <input
                                id="budget-amount"
                                type="number"
                                min="0.01"
                                step="0.01"
                                inputMode="decimal"
                                value={formData.amount}
                                onChange={(event) => setFormData({ ...formData, amount: event.target.value })}
                                placeholder="0.00"
                                required
                            />
                        </div>
                    </div>
                    <div className="form-actions">
                        <Link className="secondary-action" to="/budgets">Cancel</Link>
                        <button className="primary-action" type="submit" disabled={!id && availableCategories.length === 0}>{id ? 'Save changes' : 'Create budget'}</button>
                    </div>
                </form>
            </section>
        </main>
    );
}