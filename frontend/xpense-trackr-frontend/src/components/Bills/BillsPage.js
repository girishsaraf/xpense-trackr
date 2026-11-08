import React, { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import { formatAmountMinor, parseAmountMinor } from '../../utils/money';
import '../Ledger/LedgerPage.css';
import './BillsPage.css';

const minorToInput = (amountMinor, currency) => {
    const digits = new Intl.NumberFormat('en-US', { style: 'currency', currency }).resolvedOptions().maximumFractionDigits;
    return (amountMinor / (10 ** digits)).toFixed(digits);
};

function BillsPage() {
    const [bills, setBills] = useState([]);
    const [categories, setCategories] = useState([]);
    const [amounts, setAmounts] = useState({});
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [notice, setNotice] = useState('');

    const refresh = useCallback(async () => {
        setError('');
        try {
            const [billResponse, categoryResponse] = await Promise.all([
                axios.get(`${backendUrl}/api/bills`),
                axios.get(`${backendUrl}/api/categories`),
            ]);
            setBills(billResponse.data || []);
            setCategories(categoryResponse.data || []);
        } catch (requestError) {
            console.error('Error loading bills:', requestError);
            setError('We couldn’t load your upcoming bills. Check your connection and try again.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        refresh();
    }, [refresh]);

    const processOccurrence = async (bill, action) => {
        const amountMinor = action === 'record'
            ? parseAmountMinor(amounts[bill.id] ?? minorToInput(bill.amount_minor, bill.currency), bill.currency)
            : 0;
        if (action === 'record' && (amountMinor === null || amountMinor <= 0)) {
            setError('Enter the actual bill amount before recording payment.');
            return;
        }
        try {
            await axios.post(`${backendUrl}/api/scheduled-flows/${bill.id}/occurrences`, {
                scheduled_date: bill.next_date,
                action,
                amount_minor: amountMinor,
            });
            setAmounts((current) => {
                const next = { ...current };
                delete next[bill.id];
                return next;
            });
            setNotice(action === 'record' ? 'Payment recorded in your transaction ledger.' : 'This scheduled occurrence was skipped.');
            await refresh();
        } catch (requestError) {
            console.error('Error processing scheduled bill:', requestError);
            setError(requestError.response?.status === 409
                ? 'This bill occurrence was already processed with different details.'
                : 'We couldn’t process this bill. Refresh and try again.');
        }
    };

    if (loading) {
        return <main className="ledger-page"><p role="status">Loading bills and subscriptions…</p></main>;
    }

    return (
        <main className="ledger-page bills-page">
            <header className="ledger-heading">
                <div>
                    <p className="eyebrow">BILLS &amp; SUBSCRIPTIONS</p>
                    <h1>Upcoming bills</h1>
                    <p>Review reminders, compare recent amounts, then confirm each payment or skip it. Nothing is added to your ledger without your confirmation.</p>
                </div>
            </header>
            {error && <p className="ledger-message ledger-error" role="alert">{error}</p>}
            {notice && <p className="ledger-message ledger-notice" role="status">{notice}</p>}
            <p className="bills-explainer">Price-change signals compare the two latest expenses whose payee or description exactly matches the scheduled bill name. Verify the transactions before changing your plan.</p>
            {bills.length === 0 ? (
                <p className="ledger-empty">No active bills are scheduled. Add one from the cash-flow forecast.</p>
            ) : (
                <section aria-label="Scheduled bills" className="bills-grid">
                    {bills.map((bill) => {
                        const category = categories.find((item) => item.id === bill.category_id);
                        return (
                            <article className="bills-card" key={bill.id}>
                                <div className="bills-card-heading">
                                    <div><h2>{bill.label}</h2><p>{bill.account_name}{category ? ` · ${category.name}` : ''}</p></div>
                                    <span className={`bill-status bill-${bill.status.replace('_', '-')}`}>
                                        {bill.status === 'due_soon' ? 'Due soon' : bill.status === 'overdue' ? 'Overdue' : 'Upcoming'}
                                    </span>
                                </div>
                                <p className="bill-due-date">
                                    {bill.status === 'overdue'
                                        ? `${Math.abs(bill.days_until_due)} day${Math.abs(bill.days_until_due) === 1 ? '' : 's'} overdue`
                                        : bill.days_until_due === 0 ? 'Due today' : `Due ${bill.next_date}`}
                                    {' · '}{bill.frequency}
                                </p>
                                {bill.has_price_change && (
                                    <div className="bill-price-change" role="note">
                                        Recent amount {bill.change_percent > 0 ? 'increased' : 'decreased'} by {Math.abs(bill.change_percent)}%:
                                        {' '}{formatAmountMinor(bill.previous_amount_minor, bill.currency)} → {formatAmountMinor(bill.latest_amount_minor, bill.currency)}
                                    </div>
                                )}
                                <label className="bill-actual-amount">
                                    Actual amount paid ({bill.currency})
                                    <input
                                        type="number"
                                        min="0"
                                        step="any"
                                        value={amounts[bill.id] ?? minorToInput(bill.amount_minor, bill.currency)}
                                        onChange={(event) => setAmounts((current) => ({ ...current, [bill.id]: event.target.value }))}
                                    />
                                </label>
                                <div className="bill-actions">
                                    <button type="button" className="ledger-primary-button" onClick={() => processOccurrence(bill, 'record')}>
                                        Record payment
                                    </button>
                                    <button type="button" className="bill-skip-button" onClick={() => processOccurrence(bill, 'skip')}>
                                        Skip this occurrence
                                    </button>
                                </div>
                            </article>
                        );
                    })}
                </section>
            )}
        </main>
    );
}

export default BillsPage;
