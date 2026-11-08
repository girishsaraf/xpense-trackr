import { useCallback, useEffect, useMemo, useState } from 'react';
import type { FormEvent } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import { formatAmountMinor, parseAmountMinor } from '../../utils/money';
import { supportedCurrencies } from '../../utils/currencies';
import '../Ledger/LedgerPage.css';
import './PortfolioPage.css';

type AssetClass = 'equity' | 'fixed_income' | 'cash' | 'real_estate' | 'commodity' | 'crypto' | 'other';

type PortfolioHolding = {
    id: number;
    name: string;
    ticker: string;
    asset_class: AssetClass;
    currency: string;
    quantity: string;
    cost_basis_minor: number;
    latest_value_minor: number;
    valuation_date: string;
    valuation_source: string;
    has_valuation: boolean;
};

type HoldingValuation = {
    id: number;
    date: string;
    source: string;
    value_minor: number;
};

type HoldingForm = {
    name: string;
    ticker: string;
    asset_class: AssetClass;
    currency: string;
    quantity: string;
    cost_basis: string;
    value: string;
    date: string;
    source: string;
};

type ValuationDraft = {
    value: string;
    date: string;
    source: string;
};

const assetClasses: Array<[AssetClass, string]> = [
    ['equity', 'Stocks / equity'],
    ['fixed_income', 'Bonds / fixed income'],
    ['cash', 'Cash equivalents'],
    ['real_estate', 'Real estate'],
    ['commodity', 'Commodities'],
    ['crypto', 'Crypto'],
    ['other', 'Other'],
];

const isAssetClass = (value: string): value is AssetClass => assetClasses.some(([assetClass]) => assetClass === value);

const localDate = () => {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
};

const minorToInput = (amountMinor: number, currency: string): string => {
    const digits = new Intl.NumberFormat('en-US', { style: 'currency', currency }).resolvedOptions().maximumFractionDigits;
    return (amountMinor / (10 ** digits)).toFixed(digits);
};

const emptyForm = (): HoldingForm => ({
    name: '',
    ticker: '',
    asset_class: 'equity',
    currency: 'USD',
    quantity: '',
    cost_basis: '',
    value: '',
    date: localDate(),
    source: 'manual',
});

const PortfolioPage = () => {
    const [holdings, setHoldings] = useState<PortfolioHolding[]>([]);
    const [form, setForm] = useState<HoldingForm>(emptyForm);
    const [valuationDrafts, setValuationDrafts] = useState<Record<number, ValuationDraft>>({});
    const [history, setHistory] = useState<Record<number, HoldingValuation[] | undefined>>({});
    const [expanded, setExpanded] = useState<Record<number, boolean>>({});
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [notice, setNotice] = useState('');

    const refresh = useCallback(async () => {
        setError('');
        try {
            const response = await axios.get<PortfolioHolding[]>(`${backendUrl}/api/portfolio/holdings`);
            setHoldings(response.data || []);
        } catch (requestError) {
            console.error('Error loading portfolio holdings:', requestError);
            setError('We couldn’t load your portfolio. Check the local app server and try again.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        refresh();
    }, [refresh]);

    const totals = useMemo(() => {
        const byCurrency = new Map<string, number>();
        const byClass = new Map<string, number>();
        holdings.filter((holding) => holding.has_valuation).forEach((holding) => {
            const currencyTotal = byCurrency.get(holding.currency) || 0;
            byCurrency.set(holding.currency, currencyTotal + holding.latest_value_minor);
            const key = `${holding.currency}:${holding.asset_class}`;
            const classTotal = byClass.get(key) || 0;
            byClass.set(key, classTotal + holding.latest_value_minor);
        });
        return { byCurrency, byClass };
    }, [holdings]);

    const addHolding = async (event: FormEvent<HTMLFormElement>) => {
        event.preventDefault();
        const costBasisMinor = parseAmountMinor(form.cost_basis, form.currency);
        const valueMinor = parseAmountMinor(form.value, form.currency);
        if (!form.name.trim() || !/^[0-9]+(\.[0-9]{1,8})?$/.test(form.quantity) ||
            Number(form.quantity) <= 0 || costBasisMinor === null || costBasisMinor < 0 ||
            valueMinor === null || valueMinor < 0) {
            setError('Enter a name, positive quantity, and valid non-negative cost basis and value.');
            return;
        }
        try {
            await axios.post(`${backendUrl}/api/portfolio/holdings`, {
                name: form.name.trim(),
                ticker: form.ticker.trim(),
                asset_class: form.asset_class,
                currency: form.currency,
                quantity: form.quantity,
                cost_basis_minor: costBasisMinor,
                value_minor: valueMinor,
                date: form.date,
                source: form.source.trim() || 'manual',
            });
            setForm(emptyForm());
            setNotice('Holding added to your portfolio.');
            await refresh();
        } catch (requestError) {
            console.error('Error creating portfolio holding:', requestError);
            setError(errorMessage(requestError, 'Unable to add this holding.'));
        }
    };

    const addValuation = async (event: FormEvent<HTMLFormElement>, holding: PortfolioHolding) => {
        event.preventDefault();
        const draft = valuationDrafts[holding.id] || {};
        const valueMinor = parseAmountMinor(draft.value ?? '', holding.currency);
        if (valueMinor === null || valueMinor < 0 || !draft.date) {
            setError(`Enter a valid ${holding.currency} value and date for ${holding.name}.`);
            return;
        }
        try {
            await axios.post(`${backendUrl}/api/portfolio/holdings/${holding.id}/valuations`, {
                value_minor: valueMinor,
                date: draft.date,
                source: (draft.source || 'manual').trim(),
            });
            setNotice(`Valuation added for ${holding.name}.`);
            setValuationDrafts((current) => ({ ...current, [holding.id]: { ...current[holding.id], value: '' } }));
            setHistory((current) => ({ ...current, [holding.id]: undefined }));
            await refresh();
        } catch (requestError) {
            console.error('Error adding portfolio valuation:', requestError);
            setError(errorMessage(requestError, 'Unable to add this valuation.'));
        }
    };

    const toggleHistory = async (holding: PortfolioHolding) => {
        if (expanded[holding.id]) {
            setExpanded((current) => ({ ...current, [holding.id]: false }));
            return;
        }
        setExpanded((current) => ({ ...current, [holding.id]: true }));
        if (history[holding.id]) return;
        try {
            const response = await axios.get<HoldingValuation[]>(`${backendUrl}/api/portfolio/holdings/${holding.id}/valuations`);
            setHistory((current) => ({ ...current, [holding.id]: response.data || [] }));
        } catch (requestError) {
            console.error('Error loading portfolio valuation history:', requestError);
            setError('Unable to load valuation history.');
            setExpanded((current) => ({ ...current, [holding.id]: false }));
        }
    };

    const deleteHolding = async (holding: PortfolioHolding) => {
        if (!window.confirm(`Delete ${holding.name} and its valuation history?`)) return;
        try {
            await axios.delete(`${backendUrl}/api/portfolio/holdings/${holding.id}`);
            setNotice(`${holding.name} and its valuation history were deleted.`);
            await refresh();
        } catch (requestError) {
            console.error('Error deleting portfolio holding:', requestError);
            setError(errorMessage(requestError, 'Unable to delete this holding.'));
        }
    };

    return (
        <main className="ledger-page portfolio-page">
            <header className="ledger-heading">
                <div>
                    <p className="eyebrow">YOUR INVESTMENTS</p>
                    <h1>Portfolio</h1>
                    <p>Track quantities, cost basis, and dated manual valuations. No live market prices or currency conversion are assumed.</p>
                </div>
            </header>
            <p className="portfolio-double-count-note">
                Portfolio values are included in net worth. To avoid double-counting, don’t also enter the same investments as manual assets or in an account balance.
            </p>
            {error && <p className="ledger-message ledger-error" role="alert">{error}</p>}
            {notice && <p className="ledger-message ledger-notice" role="status">{notice}</p>}

            <section aria-labelledby="portfolio-summary-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div><p className="eyebrow">CURRENT VALUATIONS</p><h2 id="portfolio-summary-heading">Portfolio by currency</h2></div>
                </div>
                {!totals.byCurrency.size ? (
                    <p className="ledger-empty">Add a holding to start tracking your investments.</p>
                ) : (
                    <div className="portfolio-summary-grid">
                        {[...totals.byCurrency.entries()].sort(([left], [right]) => left.localeCompare(right)).map(([currency, total]) => {
                            const classes = assetClasses
                                .map(([assetClass, label]) => ({
                                    assetClass,
                                    label,
                                    value: totals.byClass.get(`${currency}:${assetClass}`) || 0,
                                }))
                                .filter((item) => item.value > 0);
                            return (
                                <article className="ledger-panel" key={currency}>
                                    <h3>{currency}</h3>
                                    <strong className="portfolio-total">{formatAmountMinor(total, currency)}</strong>
                                    <ul>
                                        {classes.map((item) => (
                                            <li key={item.assetClass}>
                                                <span>{item.label}</span>
                                                <strong>{formatAmountMinor(item.value, currency)} · {((item.value / total) * 100).toFixed(1)}%</strong>
                                            </li>
                                        ))}
                                    </ul>
                                </article>
                            );
                        })}
                    </div>
                )}
            </section>

            <section aria-labelledby="add-holding-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div><p className="eyebrow">MANUAL PORTFOLIO ENTRY</p><h2 id="add-holding-heading">Add a holding</h2></div>
                </div>
                <form className="portfolio-form" onSubmit={addHolding}>
                    <label>Name<input value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} maxLength={120} required /></label>
                    <label>Symbol<input value={form.ticker} onChange={(event) => setForm({ ...form, ticker: event.target.value })} maxLength={16} placeholder="Optional" /></label>
                    <label>
                        Asset class
                        <select value={form.asset_class} onChange={(event) => {
                            const assetClass = event.target.value;
                            if (isAssetClass(assetClass)) setForm({ ...form, asset_class: assetClass });
                        }}>
                            {assetClasses.map(([value, label]) => <option key={value} value={value}>{label}</option>)}
                        </select>
                    </label>
                    <label>Currency
                        <select value={form.currency} onChange={(event) => setForm({ ...form, currency: event.target.value })}>
                            {supportedCurrencies.map((currency) => <option key={currency} value={currency}>{currency}</option>)}
                        </select>
                    </label>
                    <label>Quantity<input type="number" min="0.00000001" step="0.00000001" value={form.quantity} onChange={(event) => setForm({ ...form, quantity: event.target.value })} required /></label>
                    <label>Cost basis ({form.currency})<input type="number" min="0" step="any" value={form.cost_basis} onChange={(event) => setForm({ ...form, cost_basis: event.target.value })} required /></label>
                    <label>Current value ({form.currency})<input type="number" min="0" step="any" value={form.value} onChange={(event) => setForm({ ...form, value: event.target.value })} required /></label>
                    <label>Valuation date<input type="date" value={form.date} onChange={(event) => setForm({ ...form, date: event.target.value })} required /></label>
                    <label>Value source<input value={form.source} onChange={(event) => setForm({ ...form, source: event.target.value })} maxLength={100} placeholder="Manual estimate, statement…" /></label>
                    <button className="ledger-primary-button" type="submit">Add holding</button>
                </form>
            </section>

            <section aria-labelledby="holdings-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div><p className="eyebrow">POSITIONS &amp; VALUE HISTORY</p><h2 id="holdings-heading">Your holdings</h2></div>
                </div>
                {loading ? (
                    <p role="status">Loading portfolio…</p>
                ) : holdings.length === 0 ? (
                    <p className="ledger-empty">No holdings yet.</p>
                ) : (
                    <div className="portfolio-holding-list">
                        {holdings.map((holding) => {
                            const draft = valuationDrafts[holding.id] || {
                                value: holding.has_valuation ? minorToInput(holding.latest_value_minor, holding.currency) : '',
                                date: localDate(),
                                source: 'manual',
                            };
                            const gainLoss = holding.latest_value_minor - holding.cost_basis_minor;
                            return (
                                <article className="ledger-panel portfolio-holding-card" key={holding.id}>
                                    <div className="portfolio-holding-heading">
                                        <div>
                                            <p className="portfolio-holding-meta">{holding.asset_class.replace('_', ' ')}{holding.ticker ? ` · ${holding.ticker}` : ''} · {holding.quantity} units</p>
                                            <h3>{holding.name}</h3>
                                            <p>{holding.has_valuation ? `Valued ${holding.valuation_date} · ${holding.valuation_source}` : 'No valuation recorded yet'}</p>
                                        </div>
                                        <div className="portfolio-holding-values">
                                            <strong>{holding.has_valuation ? formatAmountMinor(holding.latest_value_minor, holding.currency) : '—'}</strong>
                                            {holding.has_valuation && (
                                                <span>{gainLoss >= 0 ? 'Gain' : 'Loss'} vs. cost basis: {formatAmountMinor(gainLoss, holding.currency)}</span>
                                            )}
                                        </div>
                                    </div>
                                    <form className="portfolio-valuation-form" aria-label={`Add valuation for ${holding.name}`} onSubmit={(event) => addValuation(event, holding)}>
                                        <label>New value ({holding.currency})
                                            <input type="number" min="0" step="any" value={draft.value} onChange={(event) => setValuationDrafts((current) => ({ ...current, [holding.id]: { ...draft, value: event.target.value } }))} required />
                                        </label>
                                        <label>Date<input type="date" value={draft.date} onChange={(event) => setValuationDrafts((current) => ({ ...current, [holding.id]: { ...draft, date: event.target.value } }))} required /></label>
                                        <label>Source<input value={draft.source} onChange={(event) => setValuationDrafts((current) => ({ ...current, [holding.id]: { ...draft, source: event.target.value } }))} maxLength={100} /></label>
                                        <button className="ledger-primary-button" type="submit">Add valuation</button>
                                    </form>
                                    <div className="portfolio-holding-actions">
                                        <button type="button" className="portfolio-history-button" onClick={() => toggleHistory(holding)}>{expanded[holding.id] ? 'Hide valuation history' : 'Show valuation history'}</button>
                                        <button type="button" className="asset-delete-button" onClick={() => deleteHolding(holding)}>Delete holding and history</button>
                                    </div>
                                    {expanded[holding.id] && (
                                        <ul className="portfolio-valuation-history">
                                            {(history[holding.id] || []).map((valuation) => (
                                                <li key={valuation.id}>
                                                    <span>{valuation.date} · {valuation.source}</span>
                                                    <strong>{formatAmountMinor(valuation.value_minor, holding.currency)}</strong>
                                                </li>
                                            ))}
                                            {!history[holding.id]?.length && <li>No valuation history yet.</li>}
                                        </ul>
                                    )}
                                </article>
                            );
                        })}
                    </div>
                )}
            </section>
        </main>
    );
};

const errorMessage = (error: unknown, fallback: string): string => {
    if (axios.isAxiosError(error) && typeof error.response?.data === 'string') {
        return error.response.data;
    }
    return fallback;
};

export default PortfolioPage;
