import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { Link, useNavigate, useParams } from 'react-router-dom';
import backendUrl from '../../config';

export default function CategoryForm() {
    const { id } = useParams();
    const navigate = useNavigate();
    const [name, setName] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(Boolean(id));

    useEffect(() => {
        if (!id) return;
        axios.get(`${backendUrl}/api/categories/${id}`)
            .then(({ data }) => setName(data.name || ''))
            .catch((requestError) => {
                console.error('Error loading category:', requestError);
                setError('We couldn’t load this category. Return to categories and try again.');
            })
            .finally(() => setLoading(false));
    }, [id]);

    const submit = async (event) => {
        event.preventDefault();
        const categoryName = name.trim();
        if (!categoryName) {
            setError('Enter a category name to continue.');
            return;
        }
        setError('');
        try {
            if (id) await axios.put(`${backendUrl}/api/categories/${id}`, { name: categoryName });
            else await axios.post(`${backendUrl}/api/categories`, { name: categoryName });
            navigate('/categories');
        } catch (requestError) {
            console.error('Error saving category:', requestError);
            setError('We couldn’t save this category. The name may already be in use.');
        }
    };

    if (loading) return <p className="loading-message" role="status">Loading category…</p>;

    return (
        <main className="management-page form-page">
            <Link className="back-link" to="/categories">← Back to categories</Link>
            <section className="form-card">
                <p className="eyebrow">MAKE IT YOURS</p>
                <h1>{id ? 'Edit category' : 'Add a category'}</h1>
                <p className="form-description">Use a short name that makes your transactions easy to recognize.</p>
                {error && <p className="inline-message" role="alert">{error}</p>}
                <form onSubmit={submit} noValidate>
                    <div className="field-group">
                        <label htmlFor="category-name">Category name</label>
                        <input id="category-name" value={name} onChange={(event) => setName(event.target.value)} maxLength="50" placeholder="For example, Home" required autoFocus />
                    </div>
                    <div className="form-actions">
                        <Link className="secondary-action" to="/categories">Cancel</Link>
                        <button className="primary-action" type="submit">{id ? 'Save changes' : 'Add category'}</button>
                    </div>
                </form>
            </section>
        </main>
    );
}