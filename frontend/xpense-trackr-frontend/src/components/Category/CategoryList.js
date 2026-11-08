import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import backendUrl from '../../config';

export default function CategoryList() {
    const [categories, setCategories] = useState([]);
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(true);

    const loadCategories = useCallback(async () => {
        setLoading(true);
        setError('');
        try {
            const response = await axios.get(`${backendUrl}/api/categories`);
            setCategories(response.data || []);
        } catch (requestError) {
            console.error('Error fetching categories:', requestError);
            setError('We couldn’t load your categories. Check your connection and try again.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadCategories();
    }, [loadCategories]);

    const deleteCategory = async (category) => {
        if (!window.confirm(`Delete the “${category.name}” category?`)) return;
        try {
            await axios.delete(`${backendUrl}/api/categories/${category.id}`);
            await loadCategories();
        } catch (requestError) {
            console.error('Error deleting category:', requestError);
            setError('This category could not be deleted. It may still be linked to transactions.');
        }
    };

    return (
        <main className="management-page">
            <header className="page-intro">
                <div>
                    <p className="eyebrow">KEEP THINGS ORGANIZED</p>
                    <h1>Categories</h1>
                    <p className="page-description">Manage the labels you use to organize your spending.</p>
                </div>
                <Link className="primary-action" to="/categories/new">Add category</Link>
            </header>
            {error && <p className="inline-message" role="alert">{error}</p>}
            {loading ? <p className="loading-message" role="status">Loading categories…</p> : categories.length === 0 ? (
                <section className="empty-state">
                    <span className="empty-state-icon" aria-hidden="true">#</span>
                    <h2>No categories yet</h2>
                    <p>Add categories to make it easier to understand your spending.</p>
                    <Link className="primary-action" to="/categories/new">Create a category</Link>
                </section>
            ) : (
                <section className="data-card">
                    <div className="card-heading">
                        <div><p className="eyebrow">YOUR LABELS</p><h2>Spending categories</h2></div>
                        <span className="muted-copy">{categories.length} total</span>
                    </div>
                    <div className="category-list">
                        {categories.map((category) => (
                            <article className="category-row" key={category.id}>
                                <span className="category-dot" aria-hidden="true">{category.name.slice(0, 1).toUpperCase()}</span>
                                <h3>{category.name}</h3>
                                <div className="row-actions">
                                    <Link to={`/categories/${category.id}/edit`} aria-label={`Edit ${category.name} category`}>Edit</Link>
                                    <button type="button" className="text-button danger-text" onClick={() => deleteCategory(category)} aria-label={`Delete ${category.name} category`}>Delete</button>
                                </div>
                            </article>
                        ))}
                    </div>
                </section>
            )}
        </main>
    );
}