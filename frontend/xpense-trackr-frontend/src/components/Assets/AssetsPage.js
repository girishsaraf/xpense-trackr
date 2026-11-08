import { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import { formatAmountMinor, parseAmountMinor } from '../../utils/money';
import { supportedCurrencies } from '../../utils/currencies';
import '../Ledger/LedgerPage.css';
import './AssetsPage.css';

const localDate = () => {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
};

const minorToInput = (amountMinor, currency) => {
    const digits = new Intl.NumberFormat('en-US', { style: 'currency', currency }).resolvedOptions().maximumFractionDigits;
    return (amountMinor / (10 ** digits)).toFixed(digits);
};

const AssetsPage = () => {
    const [assets, setAssets] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [notice, setNotice] = useState('');
    const [drafts, setDrafts] = useState({});
    const [form, setForm] = useState({
        name: '',
        kind: 'asset',
        currency: 'USD',
        value: '',
        date: localDate(),
        source: 'manual',
    });

    const refresh = useCallback(async () => {
        setError('');
        try {
            const response = await axios.get(`${backendUrl}/api/assets`);
            setAssets(response.data || []);
        } catch (requestError) {
            console.error('Error loading assets:', requestError);
            setError('We couldn’t load your assets and liabilities. Check the local app server and try again.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        refresh();
    }, [refresh]);

    const createAsset = async (event) => {
        event.preventDefault();
        const valueMinor = parseAmountMinor(form.value, form.currency);
        if (!form.name.trim() || valueMinor === null || valueMinor < 0) {
            setError('Enter a name and a valid non-negative value for this currency.');
            return;
        }
        try {
            await axios.post(`${backendUrl}/api/assets`, {
                name: form.name.trim(),
                kind: form.kind,
                currency: form.currency,
                value_minor: valueMinor,
                date: form.date,
                source: form.source.trim() || 'manual',
            });
            setForm({
                name: '',
                kind: 'asset',
                currency: 'USD',
                value: '',
                date: localDate(),
                source: 'manual',
            });
            setNotice('Asset or liability added.');
            await refresh();
        } catch (requestError) {
            console.error('Error creating asset:', requestError);
            setError(requestError.response?.data || 'Unable to add this asset or liability.');
        }
    };

    const addValuation = async (event, asset) => {
        event.preventDefault();
        const draft = drafts[asset.id] || {};
        const valueMinor = parseAmountMinor(draft.value ?? '', asset.currency);
        if (valueMinor === null || valueMinor < 0 || !draft.date) {
            setError(`Enter a valid ${asset.currency} value and date for ${asset.name}.`);
            return;
        }
        try {
            await axios.post(`${backendUrl}/api/assets/${asset.id}/valuations`, {
                value_minor: valueMinor,
                date: draft.date,
                source: (draft.source || 'manual').trim(),
            });
            setNotice(`Valuation added for ${asset.name}.`);
            setDrafts((current) => ({ ...current, [asset.id]: { ...current[asset.id], value: '' } }));
            await refresh();
        } catch (requestError) {
            console.error('Error adding asset valuation:', requestError);
            setError(requestError.response?.data || 'Unable to add this valuation.');
        }
    };

    const deleteAsset = async (asset) => {
        if (!window.confirm(`Delete ${asset.name} and its valuation history?`)) return;
        try {
            await axios.delete(`${backendUrl}/api/assets/${asset.id}`);
            setNotice(`${asset.name} and its valuation history were deleted.`);
            await refresh();
        } catch (requestError) {
            console.error('Error deleting asset:', requestError);
            setError(requestError.response?.data || 'Unable to delete this asset.');
        }
    };

    return (
        <main className="ledger-page assets-page">
            <header className="ledger-heading">
                <div>
                    <p className="eyebrow">YOUR BALANCE SHEET</p>
                    <h1>Assets &amp; liabilities</h1>
                    <p>Track property, vehicles, loans, or other values with dated estimates and a clear source. No market prices or currency conversion are assumed.</p>
                </div>
            </header>
            {error && <p className="ledger-message ledger-error" role="alert">{error}</p>}
            {notice && <p className="ledger-message ledger-notice" role="status">{notice}</p>}

            <section aria-labelledby="add-asset-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div><p className="eyebrow">ADD TO YOUR BALANCE SHEET</p><h2 id="add-asset-heading">Add an asset or liability</h2></div>
                </div>
                <form className="asset-form" onSubmit={createAsset}>
                    <label>
                        Name
                        <input value={form.name} onChange={(event) => setForm({ ...form, name: event.target.value })} maxLength="120" required />
                    </label>
                    <label>
                        Type
                        <select value={form.kind} onChange={(event) => setForm({ ...form, kind: event.target.value })}>
                            <option value="asset">Asset</option>
                            <option value="liability">Liability</option>
                        </select>
                    </label>
                    <label>
                        Currency
                        <select value={form.currency} onChange={(event) => setForm({ ...form, currency: event.target.value })}>
                            {supportedCurrencies.map((currency) => <option key={currency} value={currency}>{currency}</option>)}
                        </select>
                    </label>
                    <label>
                        Current value ({form.currency})
                        <input type="number" min="0" step="any" value={form.value} onChange={(event) => setForm({ ...form, value: event.target.value })} required />
                    </label>
                    <label>
                        Valuation date
                        <input type="date" value={form.date} onChange={(event) => setForm({ ...form, date: event.target.value })} required />
                    </label>
                    <label>
                        Value source
                        <input value={form.source} onChange={(event) => setForm({ ...form, source: event.target.value })} maxLength="100" placeholder="Manual estimate, statement…" />
                    </label>
                    <button className="ledger-primary-button" type="submit">Add to balance sheet</button>
                </form>
            </section>

            <section aria-labelledby="assets-list-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div><p className="eyebrow">MANUAL VALUATIONS</p><h2 id="assets-list-heading">Your assets and liabilities</h2></div>
                </div>
                {loading ? (
                    <p role="status">Loading assets and liabilities…</p>
                ) : assets.length === 0 ? (
                    <p className="ledger-empty">No manual assets or liabilities yet. Add a current estimate to include it in net worth.</p>
                ) : (
                    <div className="asset-list">
                        {assets.map((asset) => {
                            const draft = drafts[asset.id] || {
                                value: asset.has_valuation ? minorToInput(asset.latest_value_minor, asset.currency) : '',
                                date: localDate(),
                                source: 'manual',
                            };
                            return (
                                <article className="ledger-panel asset-card" key={asset.id}>
                                    <div className="asset-card-heading">
                                        <div>
                                            <span className={`asset-kind asset-kind-${asset.kind}`}>{asset.kind}</span>
                                            <h3>{asset.name}</h3>
                                            {asset.has_valuation
                                                ? <p>As of {asset.valuation_date} · {asset.valuation_source}</p>
                                                : <p>No valuation recorded</p>}
                                        </div>
                                        <strong>{asset.has_valuation ? formatAmountMinor(asset.latest_value_minor, asset.currency) : '—'}</strong>
                                    </div>
                                    <form className="asset-valuation-form" aria-label={`Add valuation for ${asset.name}`} onSubmit={(event) => addValuation(event, asset)}>
                                        <label>
                                            New value ({asset.currency})
                                            <input type="number" min="0" step="any" value={draft.value} onChange={(event) => setDrafts((current) => ({ ...current, [asset.id]: { ...draft, value: event.target.value } }))} required />
                                        </label>
                                        <label>
                                            Date
                                            <input type="date" value={draft.date} onChange={(event) => setDrafts((current) => ({ ...current, [asset.id]: { ...draft, date: event.target.value } }))} required />
                                        </label>
                                        <label>
                                            Source
                                            <input value={draft.source} onChange={(event) => setDrafts((current) => ({ ...current, [asset.id]: { ...draft, source: event.target.value } }))} maxLength="100" />
                                        </label>
                                        <button className="ledger-primary-button" type="submit">Add valuation</button>
                                    </form>
                                    <button type="button" className="asset-delete-button" onClick={() => deleteAsset(asset)}>Delete item and history</button>
                                </article>
                            );
                        })}
                    </div>
                )}
            </section>
        </main>
    );
};

export default AssetsPage;
