import { useState } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import { formatAmountMinor } from '../../utils/money';
import '../Ledger/LedgerPage.css';
import './ImportExportPage.css';

const ImportExportPage = () => {
    const [file, setFile] = useState(null);
    const [report, setReport] = useState(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const [notice, setNotice] = useState('');

    const sendFile = async (commit) => {
        if (!file) {
            setError('Choose a CSV file first.');
            return;
        }
        setLoading(true);
        setError('');
        setNotice('');
        const formData = new FormData();
        formData.append('file', file);
        try {
            const response = await axios.post(
                `${backendUrl}/api/transactions/import.csv${commit ? '?commit=true' : ''}`,
                formData,
            );
            setReport(response.data);
            if (commit) setNotice(`${response.data.imported} transaction${response.data.imported === 1 ? '' : 's'} imported.`);
        } catch (requestError) {
            if (requestError.response?.data?.rows) {
                setReport(requestError.response.data);
            } else {
                setError(requestError.response?.data || 'Unable to process this CSV file.');
            }
        } finally {
            setLoading(false);
        }
    };

    const transaction = (row) => row.transaction;

    return (
        <main className="ledger-page">
            <header className="ledger-heading">
                <div>
                    <p className="eyebrow">PORTABLE BY DESIGN</p>
                    <h1>Import &amp; export</h1>
                    <p>Preview and validate transaction CSVs before importing. Re-imports are detected and skipped.</p>
                </div>
            </header>

            {error && <p className="ledger-message ledger-error" role="alert">{error}</p>}
            {notice && <p className="ledger-message ledger-notice" role="status">{notice}</p>}

            <section className="ledger-panel import-export-panel" aria-labelledby="transaction-import-heading">
                <h2 id="transaction-import-heading">Import transactions</h2>
                <p>CSV columns: <code>account_id, type, amount, currency, date, payee, description, category_id, transfer_account_id</code>. Export a CSV below to get a compatible template. Maximum 5 MB and 5,000 rows.</p>
                <label className="import-file-label">
                    Choose CSV file
                    <input type="file" accept=".csv,text/csv" onChange={(event) => {
                        setFile(event.target.files?.[0] || null);
                        setReport(null);
                        setNotice('');
                    }} />
                </label>
                <button type="button" className="ledger-primary-button" disabled={!file || loading} onClick={() => sendFile(false)}>
                    {loading ? 'Processing…' : 'Preview import'}
                </button>
            </section>

            {report && (
                <section aria-labelledby="import-preview-heading" className="ledger-section">
                    <div className="ledger-section-heading">
                        <div>
                            <p className="eyebrow">REVIEW BEFORE SAVING</p>
                            <h2 id="import-preview-heading">Import preview</h2>
                        </div>
                        <button
                            type="button"
                            className="ledger-primary-button"
                            disabled={loading || report.invalid > 0 || report.ready === 0 || Boolean(notice)}
                            onClick={() => sendFile(true)}
                        >
                            {loading ? 'Importing…' : `Import ${report.ready} ready`}
                        </button>
                    </div>
                    <p className="import-summary" role="status">
                        {report.total} rows · {report.ready} ready · {report.duplicates} duplicates · {report.invalid} invalid
                    </p>
                    <div className="ledger-table-wrap">
                        <table className="ledger-table">
                            <thead>
                                <tr><th scope="col">Line</th><th scope="col">Status</th><th scope="col">Date</th><th scope="col">Payee</th><th scope="col">Amount</th><th scope="col">Detail</th></tr>
                            </thead>
                            <tbody>
                                {report.rows.map((row, index) => {
                                    const item = transaction(row);
                                    return (
                                        <tr key={`${row.line}-${index}`}>
                                            <td>{row.line}</td>
                                            <td>{row.status}</td>
                                            <td>{item?.date || '—'}</td>
                                            <td>{item?.payee || '—'}</td>
                                            <td>{item ? formatAmountMinor(item.amount_minor, item.currency) : '—'}</td>
                                            <td>{row.message || 'Validated'}</td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>
                </section>
            )}

            <section aria-labelledby="transaction-export-heading" className="ledger-section import-export-section">
                <div>
                    <p className="eyebrow">YOUR DATA</p>
                    <h2 id="transaction-export-heading">Export transactions</h2>
                    <p>Download every ledger transaction in the same CSV format accepted by the importer.</p>
                </div>
                <a className="ledger-primary-button import-export-download" href={`${backendUrl}/api/transactions/export.csv`}>
                    Download transaction CSV
                </a>
            </section>
        </main>
    );
};

export default ImportExportPage;
