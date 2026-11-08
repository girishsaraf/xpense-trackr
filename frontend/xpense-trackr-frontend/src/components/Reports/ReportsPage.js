import { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import { formatAmountMinor } from '../../utils/money';
import '../Ledger/LedgerPage.css';
import './ReportsPage.css';

const dateText = (date) => (
    `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
);

const ReportsPage = () => {
    const now = new Date();
    const [range, setRange] = useState({
        from: dateText(new Date(now.getFullYear(), now.getMonth(), 1)),
        to: dateText(now),
    });
    const [report, setReport] = useState(null);
    const [netWorth, setNetWorth] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');

    const refresh = useCallback(async () => {
        setError('');
        try {
            const query = `from=${encodeURIComponent(range.from)}&to=${encodeURIComponent(range.to)}`;
            const [cashflowResponse, netWorthResponse] = await Promise.all([
                axios.get(`${backendUrl}/api/reports/cashflow?${query}`),
                axios.get(`${backendUrl}/api/net-worth`),
            ]);
            setReport(cashflowResponse.data);
            setNetWorth(netWorthResponse.data);
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to load your reports.');
        } finally {
            setLoading(false);
        }
    }, [range.from, range.to]);

    useEffect(() => {
        refresh();
    }, [refresh]);

    const exportURL = `${backendUrl}/api/reports/cashflow.csv?from=${encodeURIComponent(range.from)}&to=${encodeURIComponent(range.to)}`;

    if (loading) {
        return <main className="ledger-page"><p role="status">Loading your reports…</p></main>;
    }

    return (
        <main className="ledger-page">
            <header className="ledger-heading">
                <div>
                    <p className="eyebrow">YOUR FINANCIAL PICTURE</p>
                    <h1>Reports &amp; net worth</h1>
                    <p>Income and spending stay separated by currency; no exchange rates are assumed.</p>
                </div>
            </header>

            {error && <p className="ledger-message ledger-error" role="alert">{error}</p>}

            <section aria-labelledby="cashflow-report-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div>
                        <p className="eyebrow">CASH FLOW</p>
                        <h2 id="cashflow-report-heading">Income and spending</h2>
                    </div>
                    <a className="report-export-link" href={exportURL}>Download CSV</a>
                </div>
                <form className="report-range-form" onSubmit={(event) => { event.preventDefault(); refresh(); }}>
                    <label>
                        From
                        <input type="date" value={range.from} onChange={(event) => setRange({ ...range, from: event.target.value })} required />
                    </label>
                    <label>
                        To
                        <input type="date" value={range.to} onChange={(event) => setRange({ ...range, to: event.target.value })} required />
                    </label>
                    <button className="ledger-primary-button" type="submit">Update report</button>
                </form>

                {!report?.currencies?.length ? (
                    <p className="ledger-empty">No income or expenses were recorded in this period.</p>
                ) : (
                    <div className="report-currency-list">
                        {report.currencies.map((currency) => (
                            <article className="ledger-panel" key={currency.currency}>
                                <h3>{currency.currency} cash flow</h3>
                                <dl className="report-totals">
                                    <div><dt>Income</dt><dd>{formatAmountMinor(currency.income_minor, currency.currency)}</dd></div>
                                    <div><dt>Spending</dt><dd>{formatAmountMinor(currency.expense_minor, currency.currency)}</dd></div>
                                    <div><dt>Net</dt><dd>{formatAmountMinor(currency.net_minor, currency.currency)}</dd></div>
                                </dl>
                                <h4>Spending by category</h4>
                                {currency.categories.length === 0 ? (
                                    <p className="ledger-empty">No expenses in this currency for the selected period.</p>
                                ) : (
                                    <table className="ledger-table">
                                        <thead><tr><th scope="col">Category</th><th scope="col">Amount</th></tr></thead>
                                        <tbody>
                                            {currency.categories.map((category, index) => (
                                                <tr key={`${category.category_id || 'uncategorized'}-${index}`}>
                                                    <td>{category.category_name || 'Uncategorized'}</td>
                                                    <td>{formatAmountMinor(category.amount_minor, currency.currency)}</td>
                                                </tr>
                                            ))}
                                        </tbody>
                                    </table>
                                )}
                            </article>
                        ))}
                    </div>
                )}
            </section>

            <section aria-labelledby="net-worth-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div>
                        <p className="eyebrow">AS OF {netWorth?.as_of}</p>
                        <h2 id="net-worth-heading">Net worth by currency</h2>
                    </div>
                </div>
                {!netWorth?.currencies?.length ? (
                    <p className="ledger-empty">Add an account to see your net worth.</p>
                ) : (
                    <div className="report-currency-list">
                        {netWorth.currencies.map((currency) => (
                            <article className="ledger-panel" key={currency.currency}>
                                <h3>{currency.currency}</h3>
                                <dl className="report-totals">
                                    <div><dt>Assets</dt><dd>{formatAmountMinor(currency.assets_minor, currency.currency)}</dd></div>
                                    <div><dt>Liabilities</dt><dd>{formatAmountMinor(currency.liabilities_minor, currency.currency)}</dd></div>
                                    <div><dt>Net worth</dt><dd>{formatAmountMinor(currency.net_worth_minor, currency.currency)}</dd></div>
                                </dl>
                                <ul className="net-worth-account-list">
                                    {currency.accounts.map((account) => (
                                        <li key={account.account_id}>
                                            <span>{account.name} <small>{account.type.replace('_', ' ')}</small></span>
                                            <strong>{formatAmountMinor(account.balance_minor, currency.currency)}</strong>
                                        </li>
                                    ))}
                                </ul>
                            </article>
                        ))}
                    </div>
                )}
            </section>
        </main>
    );
};

export default ReportsPage;
