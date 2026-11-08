import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link, useNavigate, useParams } from 'react-router-dom';
import backendUrl from '../../config';

const asDateInput = (value) => value && !String(value).startsWith('0001-01-01') ? String(value).slice(0, 10) : '';

export default function RecurringForm() {
    const { kind, id } = useParams();
    const isExpense = kind === 'expenses';
    const navigate = useNavigate();
    const [categories, setCategories] = useState([]);
    const [formData, setFormData] = useState({
        description: '',
        amount: '',
        category_id: '',
        investment_type: 'stocks',
        start_date: '',
        end_date: '',
        frequency: 'monthly',
    });
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const loadForm = async () => {
            try {
                const requests = [];
                if (isExpense) requests.push(axios.get(`${backendUrl}/api/categories`));
                requests.push(axios.get(`${backendUrl}/api/recurring/${kind}/${id}`));
                const results = await Promise.all(requests);
                const recurring = results[results.length - 1].data;
                if (isExpense) setCategories(results[0].data || []);
                setFormData({
                    description: recurring.description || '',
                    amount: String(recurring.amount ?? ''),
                    category_id: String(recurring.category_id ?? ''),
                    investment_type: recurring.investment_type || 'stocks',
                    start_date: asDateInput(recurring.start_date),
                    end_date: asDateInput(recurring.end_date),
                    frequency: recurring.frequency || 'monthly',
                });
            } catch (requestError) {
                console.error('Error loading recurring schedule:', requestError);
                setError('We couldn’t load this schedule. Return to recurring activity and try again.');
            } finally {
                setLoading(false);
            }
        };
        loadForm();
    }, [id, isExpense, kind]);

    const updateField = (event) => {
        setFormData({ ...formData, [event.target.name]: event.target.value });
    };

    const submit = async (event) => {
        event.preventDefault();
        const amount = Number(formData.amount);
        if (!formData.description.trim() || !Number.isFinite(amount) || amount <= 0 ||
            !formData.start_date || !formData.frequency ||
            (isExpense && !formData.category_id)) {
            setError('Complete each required field and enter an amount greater than zero.');
            return;
        }
        setError('');
        const payload = {
            description: formData.description.trim(),
            amount,
            start_date: `${formData.start_date}T00:00:00Z`,
            end_date: formData.end_date ? `${formData.end_date}T00:00:00Z` : null,
            frequency: formData.frequency,
        };
        if (isExpense) payload.category_id = Number(formData.category_id);
        else payload.investment_type = formData.investment_type;

        try {
            await axios.put(`${backendUrl}/api/recurring/${kind}/${id}`, payload);
            navigate('/recurring');
        } catch (requestError) {
            console.error('Error updating recurring schedule:', requestError);
            setError('We couldn’t update this schedule. Please try again.');
        }
    };

    if (loading) return <p className="loading-message" role="status">Loading schedule…</p>;

    return (
        <main className="management-page form-page">
            <Link className="back-link" to="/recurring">← Back to recurring activity</Link>
            <section className="form-card">
                <p className="eyebrow">KEEP YOUR PLAN CURRENT</p>
                <h1>Edit recurring {isExpense ? 'expense' : 'investment'}</h1>
                <p className="form-description">Update the details or timing for this scheduled activity.</p>
                {error && <p className="inline-message" role="alert">{error}</p>}
                <form onSubmit={submit} noValidate>
                    <div className="field-group">
                        <label htmlFor="recurring-description">Description</label>
                        <input id="recurring-description" name="description" value={formData.description} onChange={updateField} maxLength="255" required />
                    </div>
                    {isExpense ? (
                        <div className="field-group">
                            <label htmlFor="recurring-category">Category</label>
                            <select id="recurring-category" name="category_id" value={formData.category_id} onChange={updateField} required>
                                {categories.map((category) => <option key={category.id} value={category.id}>{category.name}</option>)}
                            </select>
                        </div>
                    ) : (
                        <div className="field-group">
                            <label htmlFor="recurring-type">Investment type</label>
                            <select id="recurring-type" name="investment_type" value={formData.investment_type} onChange={updateField}>
                                <option value="stocks">Stocks</option>
                                <option value="401k">401(k)</option>
                                <option value="mutual_funds">Mutual Funds</option>
                                <option value="real_estate">Real Estate</option>
                                <option value="savings">Savings</option>
                                <option value="bonds">Bonds</option>
                            </select>
                        </div>
                    )}
                    <div className="field-group">
                        <label htmlFor="recurring-amount">Amount per occurrence</label>
                        <div className="currency-input">
                            <span aria-hidden="true">$</span>
                            <input id="recurring-amount" name="amount" type="number" min="0.01" step="0.01" inputMode="decimal" value={formData.amount} onChange={updateField} required />
                        </div>
                    </div>
                    <div className="form-field-row">
                        <div className="field-group">
                            <label htmlFor="recurring-start">Start date</label>
                            <input id="recurring-start" name="start_date" type="date" value={formData.start_date} onChange={updateField} required />
                        </div>
                        <div className="field-group">
                            <label htmlFor="recurring-end">End date <span className="optional-label">Optional</span></label>
                            <input id="recurring-end" name="end_date" type="date" value={formData.end_date} onChange={updateField} />
                        </div>
                    </div>
                    <div className="field-group">
                        <label htmlFor="recurring-frequency">Frequency</label>
                        <select id="recurring-frequency" name="frequency" value={formData.frequency} onChange={updateField} required>
                            <option value="daily">Daily</option>
                            <option value="weekly">Weekly</option>
                            <option value="biweekly">Every two weeks</option>
                            <option value="monthly">Monthly</option>
                        </select>
                    </div>
                    <div className="form-actions">
                        <Link className="secondary-action" to="/recurring">Cancel</Link>
                        <button className="primary-action" type="submit">Save changes</button>
                    </div>
                </form>
            </section>
        </main>
    );
}
