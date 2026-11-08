import { useCallback, useEffect, useState } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import { formatAmountMinor } from '../../utils/money';
import '../Ledger/LedgerPage.css';
import './ReceiptPage.css';

const maxReceiptSize = 10 * 1024 * 1024;

const ReceiptPage = () => {
    const [accounts, setAccounts] = useState([]);
    const [transactions, setTransactions] = useState([]);
    const [receipts, setReceipts] = useState([]);
    const [transactionId, setTransactionId] = useState('');
    const [file, setFile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [notice, setNotice] = useState('');

    const refresh = useCallback(async () => {
        setError('');
        try {
            const [accountResponse, transactionResponse, receiptResponse] = await Promise.all([
                axios.get(`${backendUrl}/api/accounts`),
                axios.get(`${backendUrl}/api/transactions?limit=100`),
                axios.get(`${backendUrl}/api/receipts`),
            ]);
            setAccounts(accountResponse.data || []);
            setTransactions(transactionResponse.data || []);
            setReceipts(receiptResponse.data || []);
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to load receipts and transactions.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        refresh();
    }, [refresh]);

    const submitReceipt = async (event) => {
        event.preventDefault();
        if (!transactionId || !file) {
            setError('Choose a transaction and a receipt file.');
            return;
        }
        if (file.size === 0 || file.size > maxReceiptSize) {
            setError('Receipt files must be between 1 byte and 10 MB.');
            return;
        }
        const formData = new FormData();
        formData.append('transaction_id', transactionId);
        formData.append('file', file);
        try {
            await axios.post(`${backendUrl}/api/receipts`, formData);
            setFile(null);
            event.target.reset();
            setTransactionId('');
            setNotice('Receipt attached to the transaction.');
            await refresh();
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to attach this receipt.');
        }
    };

    const deleteReceipt = async (id) => {
        try {
            await axios.delete(`${backendUrl}/api/receipts/${id}`);
            setNotice('Receipt deleted.');
            await refresh();
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to delete this receipt.');
        }
    };

    const transactionLabel = (transaction) => {
        const account = accounts.find((item) => item.id === transaction.account_id);
        const label = transaction.payee || transaction.description || transaction.type;
        return `${transaction.date} · ${label} · ${account?.name || 'Account'} · ${formatAmountMinor(transaction.amount_minor, transaction.currency)}`;
    };

    if (loading) {
        return <main className="ledger-page"><p role="status">Loading receipts…</p></main>;
    }

    return (
        <main className="ledger-page">
            <header className="ledger-heading">
                <div>
                    <p className="eyebrow">LOCAL DOCUMENTS</p>
                    <h1>Receipts</h1>
                    <p>Attach PDF or image receipts to transactions. Files stay in your local app data folder.</p>
                </div>
            </header>

            {error && <p className="ledger-message ledger-error" role="alert">{error}</p>}
            {notice && <p className="ledger-message ledger-notice" role="status">{notice}</p>}

            <section className="ledger-panel receipt-upload-panel" aria-labelledby="receipt-upload-heading">
                <h2 id="receipt-upload-heading">Attach a receipt</h2>
                {transactions.length === 0 ? (
                    <p className="ledger-empty">Add a transaction before attaching a receipt.</p>
                ) : (
                    <form onSubmit={submitReceipt} className="ledger-form">
                        <label>
                            Transaction
                            <select value={transactionId} onChange={(event) => setTransactionId(event.target.value)} required>
                                <option value="">Choose a transaction</option>
                                {transactions.map((transaction) => (
                                    <option key={transaction.id} value={transaction.id}>{transactionLabel(transaction)}</option>
                                ))}
                            </select>
                        </label>
                        <label>
                            PDF, PNG, or JPEG (up to 10 MB)
                            <input
                                type="file"
                                accept="application/pdf,image/png,image/jpeg,.pdf,.png,.jpg,.jpeg"
                                onChange={(event) => setFile(event.target.files?.[0] || null)}
                                required
                            />
                        </label>
                        <button type="submit" className="ledger-primary-button">Attach receipt</button>
                    </form>
                )}
            </section>

            <section aria-labelledby="saved-receipts-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div>
                        <p className="eyebrow">DOCUMENTS</p>
                        <h2 id="saved-receipts-heading">Saved receipts</h2>
                    </div>
                </div>
                {receipts.length === 0 ? (
                    <p className="ledger-empty">No receipts attached yet.</p>
                ) : (
                    <ul className="receipt-list">
                        {receipts.map((receipt) => {
                            const transaction = transactions.find((item) => item.id === receipt.transaction_id);
                            return (
                                <li key={receipt.id}>
                                    <span>
                                        <strong>{receipt.file_name}</strong>
                                        <small>{transaction ? transactionLabel(transaction) : 'Transaction unavailable'} · {(receipt.size_bytes / 1024).toFixed(0)} KB</small>
                                    </span>
                                    <span className="receipt-actions">
                                        <a href={`${backendUrl}/api/receipts/${receipt.id}/file`}>Download</a>
                                        <button type="button" onClick={() => deleteReceipt(receipt.id)} aria-label={`Delete ${receipt.file_name}`}>
                                            Delete
                                        </button>
                                    </span>
                                </li>
                            );
                        })}
                    </ul>
                )}
            </section>
        </main>
    );
};

export default ReceiptPage;
