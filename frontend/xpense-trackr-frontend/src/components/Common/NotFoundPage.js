import React from 'react';
import { Link } from 'react-router-dom';

export default function NotFoundPage() {
    return (
        <main className="management-page">
            <section className="empty-state">
                <span className="empty-state-icon" aria-hidden="true">404</span>
                <h1>We can’t find that page</h1>
                <p>The address may have changed, or the page may not exist.</p>
                <Link className="primary-action" to="/">Go to your overview</Link>
            </section>
        </main>
    );
}
