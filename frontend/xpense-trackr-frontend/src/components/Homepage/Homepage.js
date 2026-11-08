import { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import { Link } from 'react-router-dom';
import backendUrl from '../../config';
import { formatAmountMinor } from '../../utils/money';
import './Dashboard.css';

const Homepage = () => {
    const [accounts, setAccounts] = useState([]);
    const [transactions, setTransactions] = useState([]);
    const [cashflow, setCashflow] = useState(null);
    const [netWorth, setNetWorth] = useState(null);
    const [forecast, setForecast] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    const loadDashboard = useCallback(async () => {
        setError('');
        try {
            const [accountResponse, transactionResponse, cashflowResponse, netWorthResponse, forecastResponse] = await Promise.all([
                axios.get(`${backendUrl}/api/accounts`),
                axios.get(`${backendUrl}/api/transactions?limit=8`),
                axios.get(`${backendUrl}/api/reports/cashflow`),
                axios.get(`${backendUrl}/api/net-worth`),
                axios.get(`${backendUrl}/api/forecast/cashflow?days=30`),
            ]);
            setAccounts(accountResponse.data || []);
            setTransactions(transactionResponse.data || []);
            setCashflow(cashflowResponse.data);
            setNetWorth(netWorthResponse.data);
            setForecast(forecastResponse.data);
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to load your financial overview.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        loadDashboard();
    }, [loadDashboard]);

    if (loading) {
        return <main className="dashboard-page"><p role="status">Loading your financial overview…</p></main>;
    }

    return (
        <main className="dashboard-page">
            <header className="dashboard-heading">
                <div>
                    <p className="eyebrow">YOUR FINANCES AT A GLANCE</p>
                    <h1>Your money, in focus.</h1>
                    <p>One private, local view of the accounts and activity you choose to track.</p>
                </div>
                <Link className="primary-action" to="/transactions">Add transaction</Link>
            </header>

            {error && <p className="ledger-message ledger-error" role="alert">{error}</p>}

            {accounts.length === 0 ? (
                <section className="dashboard-onboarding">
                    <span className="dashboard-onboarding-mark" aria-hidden="true">$</span>
                    <p className="eyebrow">START WITH THE BASICS</p>
                    <h2>Your finances stay on this device.</h2>
                    <p>Add an account and your first transaction. No bank connection or cloud account is required.</p>
                    <Link className="primary-action" to="/transactions">Set up an account</Link>
                    <Link className="secondary-action" to="/import-export">Import existing transactions</Link>
                </section>
            ) : (
                <>
                    <section aria-labelledby="dashboard-cashflow-heading" className="dashboard-section">
                        <div className="dashboard-section-heading">
                            <div>
                                <p className="eyebrow">THIS MONTH · BY CURRENCY</p>
                                <h2 id="dashboard-cashflow-heading">Income and spending</h2>
                            </div>
                            <Link to="/reports">View reports</Link>
                        </div>
                        {!cashflow?.currencies?.length ? (
                            <p className="dashboard-empty">No income or spending recorded yet this month.</p>
                        ) : (
                            <div className="dashboard-cashflow-grid">
                                {cashflow.currencies.map((item) => (
                                    <article className="dashboard-currency-card" key={item.currency}>
                                        <span>{item.currency}</span>
                                        <dl>
                                            <div><dt>Income</dt><dd>{formatAmountMinor(item.income_minor, item.currency)}</dd></div>
                                            <div><dt>Spending</dt><dd>{formatAmountMinor(item.expense_minor, item.currency)}</dd></div>
                                            <div><dt>Net</dt><dd>{formatAmountMinor(item.net_minor, item.currency)}</dd></div>
                                        </dl>
                                    </article>
                                ))}
                            </div>
                        )}
                    </section>

                    <section aria-labelledby="dashboard-accounts-heading" className="dashboard-section">
                        <div className="dashboard-section-heading">
                            <div>
                                <p className="eyebrow">BALANCES</p>
                                <h2 id="dashboard-accounts-heading">Your accounts</h2>
                            </div>
                            <Link to="/reports">Net worth</Link>
                        </div>
                        <ul className="dashboard-account-list">
                            {accounts.map((account) => (
                                <li key={account.id}>
                                    <span><strong>{account.name}</strong><small>{account.type.replace('_', ' ')}</small></span>
                                    <strong>{formatAmountMinor(account.balance_minor, account.currency)}</strong>
                                </li>
                            ))}
                        </ul>
                        {netWorth?.currencies?.length > 0 && (
                            <p className="dashboard-note">
                                Net worth is shown separately for each currency; no exchange rate is assumed.
                            </p>
                        )}
                    </section>

                    <section aria-labelledby="dashboard-upcoming-heading" className="dashboard-section">
                        <div className="dashboard-section-heading">
                            <div>
                                <p className="eyebrow">NEXT 30 DAYS</p>
                                <h2 id="dashboard-upcoming-heading">Expected income and bills</h2>
                            </div>
                            <Link to="/forecast">Plan cash flow</Link>
                        </div>
                        {!forecast?.events?.length ? (
                            <p className="dashboard-empty">Add expected income or bills to see them here.</p>
                        ) : (
                            <ul className="dashboard-account-list">
                                {forecast.events.slice(0, 5).map((event, index) => (
                                    <li key={`${event.scheduled_flow_id}-${event.date}-${index}`}>
                                        <span><strong>{event.label}</strong><small>{event.date}</small></span>
                                        <strong>{event.type === 'expense' ? '−' : '+'}{formatAmountMinor(event.amount_minor, event.currency)}</strong>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </section>

                    <section aria-labelledby="dashboard-activity-heading" className="dashboard-section">
                        <div className="dashboard-section-heading">
                            <div>
                                <p className="eyebrow">RECENT ACTIVITY</p>
                                <h2 id="dashboard-activity-heading">Latest transactions</h2>
                            </div>
                            <Link to="/transactions">All transactions</Link>
                        </div>
                        {transactions.length === 0 ? (
                            <p className="dashboard-empty">Your transactions will show up here.</p>
                        ) : (
                            <div className="ledger-table-wrap">
                                <table className="ledger-table">
                                    <thead><tr><th scope="col">Date</th><th scope="col">Payee</th><th scope="col">Type</th><th scope="col">Amount</th></tr></thead>
                                    <tbody>
                                        {transactions.map((transaction) => (
                                            <tr key={transaction.id}>
                                                <td>{transaction.date}</td>
                                                <td>{transaction.payee || transaction.description || 'Transaction'}</td>
                                                <td>{transaction.type}</td>
                                                <td>{formatAmountMinor(transaction.amount_minor, transaction.currency)}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}
                    </section>
                </>
            )}

            <nav aria-label="Finance shortcuts" className="dashboard-shortcuts">
                <Link to="/import-export">Import or export data</Link>
                <Link to="/receipts">Receipts</Link>
                <Link to="/goals">Savings goals</Link>
                <Link to="/budgets">Budgets</Link>
            </nav>
        </main>
    );
};

export default Homepage;
