import { useCallback, useEffect, useMemo, useState } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import { formatAmountMinor, parseAmountMinor } from '../../utils/money';
import '../Ledger/LedgerPage.css';
import './ForecastPage.css';

const localDate = (date = new Date()) => (
    `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
);

const ForecastPage = () => {
    const [accounts, setAccounts] = useState([]);
    const [flows, setFlows] = useState([]);
    const [forecast, setForecast] = useState(null);
    const [days, setDays] = useState(90);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [notice, setNotice] = useState('');
    const [form, setForm] = useState({
        accountId: '',
        label: '',
        type: 'expense',
        amount: '',
        frequency: 'monthly',
        nextDate: localDate(new Date(Date.now() + 86400000)),
    });

    const refresh = useCallback(async () => {
        setError('');
        try {
            const [accountResponse, flowResponse, forecastResponse] = await Promise.all([
                axios.get(`${backendUrl}/api/accounts`),
                axios.get(`${backendUrl}/api/scheduled-flows`),
                axios.get(`${backendUrl}/api/forecast/cashflow?days=${days}`),
            ]);
            setAccounts(accountResponse.data || []);
            setFlows(flowResponse.data || []);
            setForecast(forecastResponse.data);
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to load your cash-flow forecast.');
        } finally {
            setLoading(false);
        }
    }, [days]);

    useEffect(() => {
        refresh();
    }, [refresh]);

    const selectedAccount = useMemo(
        () => accounts.find((account) => String(account.id) === form.accountId),
        [accounts, form.accountId],
    );

    const submitFlow = async (event) => {
        event.preventDefault();
        const amountMinor = selectedAccount ? parseAmountMinor(form.amount, selectedAccount.currency) : null;
        if (!selectedAccount || amountMinor === null || amountMinor <= 0 || !form.label.trim()) {
            setError('Choose an account, provide a label, and enter a positive amount.');
            return;
        }
        try {
            await axios.post(`${backendUrl}/api/scheduled-flows`, {
                account_id: selectedAccount.id,
                label: form.label.trim(),
                type: form.type,
                amount_minor: amountMinor,
                currency: selectedAccount.currency,
                frequency: form.frequency,
                next_date: form.nextDate,
                enabled: true,
            });
            setForm((current) => ({ ...current, label: '', amount: '' }));
            setNotice('Scheduled item added to your forecast.');
            await refresh();
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to add this scheduled item.');
        }
    };

    const deleteFlow = async (id) => {
        try {
            await axios.delete(`${backendUrl}/api/scheduled-flows/${id}`);
            setNotice('Scheduled item removed.');
            await refresh();
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to remove this scheduled item.');
        }
    };

    if (loading) {
        return <main className="ledger-page"><p role="status">Loading your cash-flow forecast…</p></main>;
    }

    return (
        <main className="ledger-page forecast-page">
            <header className="ledger-heading">
                <div>
                    <p className="eyebrow">PLAN AHEAD</p>
                    <h1>Cash-flow forecast</h1>
                    <p>Projected balances use only the future income and bills you schedule here. They never create transactions.</p>
                </div>
                <label className="forecast-range">
                    Forecast range
                    <select value={days} onChange={(event) => setDays(Number(event.target.value))}>
                        <option value={30}>30 days</option>
                        <option value={90}>90 days</option>
                        <option value={180}>180 days</option>
                        <option value={365}>1 year</option>
                    </select>
                </label>
            </header>

            {error && <p className="ledger-message ledger-error" role="alert">{error}</p>}
            {notice && <p className="ledger-message ledger-notice" role="status">{notice}</p>}

            {!accounts.length ? (
                <p className="ledger-empty">Add an account before planning future income or bills.</p>
            ) : (
                <section aria-labelledby="forecast-balances-heading" className="ledger-section">
                    <div className="ledger-section-heading">
                        <div>
                            <p className="eyebrow">ESTIMATES · {forecast?.start_date} TO {forecast?.end_date}</p>
                            <h2 id="forecast-balances-heading">Projected account balances</h2>
                        </div>
                    </div>
                    <ul className="ledger-account-grid">
                        {(forecast?.accounts || []).map((account) => (
                            <li key={account.account_id} className="ledger-account-card">
                                <span className="ledger-account-type">{account.currency}</span>
                                <h3>{account.name}</h3>
                                <p>{formatAmountMinor(account.projected_balance_minor, account.currency)}</p>
                                <small>Starting at {formatAmountMinor(account.current_balance_minor, account.currency)}</small>
                            </li>
                        ))}
                    </ul>
                </section>
            )}

            <section aria-labelledby="scheduled-items-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div>
                        <p className="eyebrow">EXPECTED INCOME &amp; BILLS</p>
                        <h2 id="scheduled-items-heading">Scheduled items</h2>
                    </div>
                </div>
                <div className="ledger-forms">
                    <form onSubmit={submitFlow} className="ledger-panel ledger-form">
                        <h3>Add expected income or bill</h3>
                        <label>
                            Account
                            <select
                                value={form.accountId}
                                onChange={(event) => setForm({ ...form, accountId: event.target.value })}
                                required
                            >
                                <option value="">Choose an account</option>
                                {accounts.map((account) => (
                                    <option key={account.id} value={account.id}>{account.name} ({account.currency})</option>
                                ))}
                            </select>
                        </label>
                        <label>
                            Type
                            <select value={form.type} onChange={(event) => setForm({ ...form, type: event.target.value })}>
                                <option value="expense">Bill or expense</option>
                                <option value="income">Expected income</option>
                            </select>
                        </label>
                        <label>
                            Name
                            <input value={form.label} onChange={(event) => setForm({ ...form, label: event.target.value })} maxLength="180" required />
                        </label>
                        <label>
                            Amount {selectedAccount ? `(${selectedAccount.currency})` : ''}
                            <input inputMode="decimal" value={form.amount} onChange={(event) => setForm({ ...form, amount: event.target.value })} required />
                        </label>
                        <label>
                            Repeats
                            <select value={form.frequency} onChange={(event) => setForm({ ...form, frequency: event.target.value })}>
                                <option value="once">Once</option>
                                <option value="weekly">Weekly</option>
                                <option value="biweekly">Every two weeks</option>
                                <option value="monthly">Monthly</option>
                                <option value="yearly">Yearly</option>
                            </select>
                        </label>
                        <label>
                            Next date
                            <input type="date" value={form.nextDate} onChange={(event) => setForm({ ...form, nextDate: event.target.value })} required />
                        </label>
                        <button type="submit" className="ledger-primary-button" disabled={accounts.length === 0}>Add to forecast</button>
                    </form>

                    <div className="ledger-panel">
                        <h3>Upcoming occurrences</h3>
                        {!forecast?.events?.length ? (
                            <p className="ledger-empty">No scheduled items fall within this forecast range.</p>
                        ) : (
                            <ul className="forecast-event-list">
                                {forecast.events.map((event, index) => {
                                    const account = accounts.find((item) => item.id === event.account_id);
                                    return (
                                        <li key={`${event.scheduled_flow_id}-${event.date}-${index}`}>
                                            <span>
                                                <strong>{event.label}</strong>
                                                <small>{event.date} · {account?.name || 'Account'}</small>
                                            </span>
                                            <strong className={event.type === 'expense' ? 'forecast-expense' : 'forecast-income'}>
                                                {event.type === 'expense' ? '−' : '+'}{formatAmountMinor(event.amount_minor, event.currency)}
                                            </strong>
                                        </li>
                                    );
                                })}
                            </ul>
                        )}
                    </div>
                </div>
            </section>

            <section aria-labelledby="manage-scheduled-items-heading" className="ledger-section">
                <h2 id="manage-scheduled-items-heading">Manage schedules</h2>
                {flows.length === 0 ? (
                    <p className="ledger-empty">No schedules saved yet.</p>
                ) : (
                    <ul className="forecast-schedule-list">
                        {flows.map((flow) => (
                            <li key={flow.id}>
                                <span>
                                    <strong>{flow.label}</strong>
                                    <small>{flow.frequency} · next on {flow.next_date} · {accounts.find((account) => account.id === flow.account_id)?.name || 'Account'}</small>
                                </span>
                                <button type="button" onClick={() => deleteFlow(flow.id)} aria-label={`Remove ${flow.label}`}>
                                    Remove
                                </button>
                            </li>
                        ))}
                    </ul>
                )}
            </section>
        </main>
    );
};

export default ForecastPage;
