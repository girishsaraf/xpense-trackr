import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link, useNavigate, useParams } from 'react-router-dom';
import backendUrl from '../../config';

const emptyGoal = { name: '', target_amount: '', current_amount: '0', target_date: '' };

export default function GoalForm() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [formData, setFormData] = useState(emptyGoal);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(Boolean(id));

    useEffect(() => {
        if (!id) return;
        axios.get(`${backendUrl}/api/goals/${id}`)
            .then(({ data }) => setFormData({
                name: data.name || '',
                target_amount: String(data.target_amount ?? ''),
                current_amount: String(data.current_amount ?? 0),
                target_date: data.target_date || '',
            }))
            .catch((requestError) => {
                console.error('Error loading savings goal:', requestError);
                setError('We couldn’t load this goal. Return to your goals and try again.');
            })
            .finally(() => setLoading(false));
    }, [id]);

    const updateField = (event) => {
        setFormData({ ...formData, [event.target.name]: event.target.value });
    };

    const submit = async (event) => {
        event.preventDefault();
        setError('');
        const target = Number(formData.target_amount);
        const current = Number(formData.current_amount);
        if (!formData.name.trim() || !Number.isFinite(target) || target <= 0 || !Number.isFinite(current) || current < 0) {
            setError('Add a name, a target greater than zero, and a current balance of zero or more.');
            return;
        }
        const payload = {
            name: formData.name.trim(),
            target_amount: target,
            current_amount: current,
            target_date: formData.target_date,
        };
        try {
            if (id) await axios.put(`${backendUrl}/api/goals/${id}`, payload);
            else await axios.post(`${backendUrl}/api/goals`, payload);
            navigate('/goals');
        } catch (requestError) {
            console.error('Error saving savings goal:', requestError);
            setError('We couldn’t save this goal. Please check the details and try again.');
        }
    };

    if (loading) return <p className="loading-message" role="status">Loading goal…</p>;

    return (
        <main className="management-page form-page">
            <Link className="back-link" to="/goals">← Back to goals</Link>
            <section className="form-card">
                <p className="eyebrow">A LITTLE CLOSER EVERY DAY</p>
                <h1>{id ? 'Edit savings goal' : 'Create a savings goal'}</h1>
                <p className="form-description">Name what you’re saving for, set your target, and track your progress.</p>
                {error && <p className="inline-message" role="alert">{error}</p>}
                <form onSubmit={submit} noValidate>
                    <div className="field-group">
                        <label htmlFor="goal-name">Goal name</label>
                        <input id="goal-name" name="name" value={formData.name} onChange={updateField} maxLength="120" placeholder="For example, emergency fund" required />
                    </div>
                    <div className="form-field-row">
                        <div className="field-group">
                            <label htmlFor="goal-target">Savings target</label>
                            <div className="currency-input">
                                <span aria-hidden="true">$</span>
                                <input id="goal-target" name="target_amount" type="number" min="0.01" step="0.01" inputMode="decimal" value={formData.target_amount} onChange={updateField} placeholder="0.00" required />
                            </div>
                        </div>
                        <div className="field-group">
                            <label htmlFor="goal-current">Already saved</label>
                            <div className="currency-input">
                                <span aria-hidden="true">$</span>
                                <input id="goal-current" name="current_amount" type="number" min="0" step="0.01" inputMode="decimal" value={formData.current_amount} onChange={updateField} placeholder="0.00" required />
                            </div>
                        </div>
                    </div>
                    <div className="field-group">
                        <label htmlFor="goal-date">Target date <span className="optional-label">Optional</span></label>
                        <input id="goal-date" name="target_date" type="date" value={formData.target_date} onChange={updateField} />
                    </div>
                    <div className="form-actions">
                        <Link className="secondary-action" to="/goals">Cancel</Link>
                        <button className="primary-action" type="submit">{id ? 'Save changes' : 'Create goal'}</button>
                    </div>
                </form>
            </section>
        </main>
    );
}
