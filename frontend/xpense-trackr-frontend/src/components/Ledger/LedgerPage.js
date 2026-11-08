import { useCallback, useEffect, useMemo, useState } from 'react';
import axios from 'axios';
import backendUrl from '../../config';
import { formatAmountMinor, parseAmountMinor } from '../../utils/money';
import './LedgerPage.css';

const accountTypes = [
    ['checking', 'Checking'],
    ['savings', 'Savings'],
    ['cash', 'Cash'],
    ['credit_card', 'Credit card'],
    ['investment', 'Investment'],
    ['loan', 'Loan'],
];

const currencies = [
    'AED', 'AUD', 'BHD', 'BRL', 'CAD', 'CHF', 'CLP', 'CNY', 'DKK', 'EUR',
    'GBP', 'HKD', 'IDR', 'INR', 'IQD', 'ISK', 'JPY', 'JOD', 'KRW', 'KWD',
    'MXN', 'MYR', 'NOK', 'NZD', 'OMR', 'PHP', 'PLN', 'SAR', 'SEK', 'SGD',
    'THB', 'TND', 'TRY', 'USD', 'VND', 'ZAR',
];

const localDate = () => {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
};

const LedgerPage = () => {
    const [accounts, setAccounts] = useState([]);
    const [categories, setCategories] = useState([]);
    const [transactions, setTransactions] = useState([]);
    const [categoryRules, setCategoryRules] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [notice, setNotice] = useState('');
    const [accountForm, setAccountForm] = useState({
        name: '',
        type: 'checking',
        currency: 'USD',
        openingBalance: '0',
    });
    const [transactionForm, setTransactionForm] = useState({
        accountId: '',
        transferAccountId: '',
        type: 'expense',
        categoryId: '',
        amount: '',
        payee: '',
        description: '',
        date: localDate(),
    });
    const [ruleForm, setRuleForm] = useState({ name: '', matchValue: '', categoryId: '' });

    const refresh = useCallback(async () => {
        setError('');
        try {
            const [accountResponse, categoryResponse, transactionResponse, rulesResponse] = await Promise.all([
                axios.get(`${backendUrl}/api/accounts`),
                axios.get(`${backendUrl}/api/categories`),
                axios.get(`${backendUrl}/api/transactions?limit=50`),
                axios.get(`${backendUrl}/api/category-rules`),
            ]);
            setAccounts(accountResponse.data || []);
            setCategories(categoryResponse.data || []);
            setTransactions(transactionResponse.data || []);
            setCategoryRules(rulesResponse.data || []);
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to load accounts and transactions.');
        } finally {
            setLoading(false);
        }
    }, []);

    useEffect(() => {
        refresh();
    }, [refresh]);

    const selectedAccount = useMemo(
        () => accounts.find((account) => String(account.id) === transactionForm.accountId),
        [accounts, transactionForm.accountId],
    );
    const transferAccounts = useMemo(
        () => accounts.filter(
            (account) => account.id !== selectedAccount?.id && account.currency === selectedAccount?.currency,
        ),
        [accounts, selectedAccount],
    );

    const updateTransactionField = (event) => {
        const { name, value } = event.target;
        setTransactionForm((current) => ({
            ...current,
            [name]: value,
            ...(name === 'accountId' ? { transferAccountId: '' } : {}),
            ...(name === 'type' && value === 'transfer' ? { categoryId: '' } : {}),
        }));
    };

    const submitAccount = async (event) => {
        event.preventDefault();
        const openingBalanceMinor = parseAmountMinor(accountForm.openingBalance, accountForm.currency.toUpperCase());
        if (openingBalanceMinor === null || !accountForm.name.trim()) {
            setError('Enter an account name and a valid opening balance for the selected currency.');
            return;
        }

        try {
            await axios.post(`${backendUrl}/api/accounts`, {
                name: accountForm.name,
                type: accountForm.type,
                currency: accountForm.currency.toUpperCase(),
                opening_balance_minor: openingBalanceMinor,
            });
            setAccountForm({ name: '', type: 'checking', currency: 'USD', openingBalance: '0' });
            setNotice('Account added.');
            await refresh();
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to add account.');
        }
    };

    const submitTransaction = async (event) => {
        event.preventDefault();
        const amountMinor = selectedAccount
            ? parseAmountMinor(transactionForm.amount, selectedAccount.currency)
            : null;
        if (!selectedAccount || amountMinor === null || amountMinor <= 0) {
            setError('Choose an account and enter a positive amount using its currency precision.');
            return;
        }
        if (transactionForm.type === 'transfer' && !transactionForm.transferAccountId) {
            setError('Choose a destination account in the same currency.');
            return;
        }

        const payload = {
            account_id: selectedAccount.id,
            type: transactionForm.type,
            amount_minor: amountMinor,
            currency: selectedAccount.currency,
            date: transactionForm.date,
            payee: transactionForm.payee,
            description: transactionForm.description,
        };
        if (transactionForm.type === 'transfer') {
            payload.transfer_account_id = Number(transactionForm.transferAccountId);
        } else if (transactionForm.categoryId) {
            payload.category_id = Number(transactionForm.categoryId);
        }

        try {
            await axios.post(`${backendUrl}/api/transactions`, payload);
            setTransactionForm((current) => ({
                ...current,
                amount: '',
                payee: '',
                description: '',
                date: localDate(),
            }));
            setNotice('Transaction added.');
            await refresh();
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to add transaction.');
        }
    };

    const submitCategoryRule = async (event) => {
        event.preventDefault();
        if (!ruleForm.name.trim() || !ruleForm.matchValue.trim() || !ruleForm.categoryId) {
            setError('Enter a rule name, payee text, and category.');
            return;
        }
        try {
            await axios.post(`${backendUrl}/api/category-rules`, {
                name: ruleForm.name.trim(),
                match_value: ruleForm.matchValue.trim(),
                category_id: Number(ruleForm.categoryId),
                priority: 0,
                enabled: true,
            });
            setRuleForm({ name: '', matchValue: '', categoryId: '' });
            setNotice('Categorization rule added.');
            await refresh();
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to add categorization rule.');
        }
    };

    const deleteCategoryRule = async (id) => {
        try {
            await axios.delete(`${backendUrl}/api/category-rules/${id}`);
            setNotice('Categorization rule deleted.');
            await refresh();
        } catch (requestError) {
            setError(requestError.response?.data || 'Unable to delete categorization rule.');
        }
    };

    if (loading) {
        return <main className="ledger-page"><p role="status">Loading your accounts…</p></main>;
    }

    return (
        <main className="ledger-page">
            <header className="ledger-heading">
                <div>
                    <p className="eyebrow">YOUR MONEY</p>
                    <h1>Accounts &amp; transactions</h1>
                    <p>Keep income, spending, and transfers together in one local ledger.</p>
                </div>
            </header>

            {error && <p className="ledger-message ledger-error" role="alert">{error}</p>}
            {notice && <p className="ledger-message ledger-notice" role="status">{notice}</p>}

            <section aria-labelledby="accounts-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div>
                        <p className="eyebrow">BALANCES</p>
                        <h2 id="accounts-heading">Your accounts</h2>
                    </div>
                </div>
                {accounts.length === 0 ? (
                    <p className="ledger-empty">Add an account to start recording transactions.</p>
                ) : (
                    <ul className="ledger-account-grid">
                        {accounts.map((account) => (
                            <li key={account.id} className="ledger-account-card">
                                <span className="ledger-account-type">{account.type.replace('_', ' ')}</span>
                                <h3>{account.name}</h3>
                                <p>{formatAmountMinor(account.balance_minor, account.currency)}</p>
                            </li>
                        ))}
                    </ul>
                )}
            </section>

            <div className="ledger-forms">
                <section aria-labelledby="add-account-heading" className="ledger-panel">
                    <h2 id="add-account-heading">Add an account</h2>
                    <form onSubmit={submitAccount} className="ledger-form">
                        <label>
                            Account name
                            <input
                                name="name"
                                value={accountForm.name}
                                onChange={(event) => setAccountForm({ ...accountForm, name: event.target.value })}
                                maxLength="120"
                                required
                            />
                        </label>
                        <label>
                            Account type
                            <select
                                value={accountForm.type}
                                onChange={(event) => setAccountForm({ ...accountForm, type: event.target.value })}
                            >
                                {accountTypes.map(([value, label]) => <option key={value} value={value}>{label}</option>)}
                            </select>
                        </label>
                        <label>
                            Currency
                            <select
                                value={accountForm.currency}
                                onChange={(event) => setAccountForm({ ...accountForm, currency: event.target.value })}
                            >
                                {currencies.map((currency) => <option key={currency} value={currency}>{currency}</option>)}
                            </select>
                        </label>
                        <label>
                            Opening balance
                            <input
                                inputMode="decimal"
                                value={accountForm.openingBalance}
                                onChange={(event) => setAccountForm({ ...accountForm, openingBalance: event.target.value })}
                                required
                            />
                        </label>
                        <button type="submit" className="ledger-primary-button">Add account</button>
                    </form>
                </section>

                <section aria-labelledby="add-transaction-heading" className="ledger-panel">
                    <h2 id="add-transaction-heading">Add a transaction</h2>
                    {accounts.length === 0 ? (
                        <p className="ledger-empty">Add an account before recording a transaction.</p>
                    ) : (
                        <form onSubmit={submitTransaction} className="ledger-form">
                            <label>
                                Type
                                <select name="type" value={transactionForm.type} onChange={updateTransactionField}>
                                    <option value="expense">Expense</option>
                                    <option value="income">Income</option>
                                    <option value="transfer">Transfer</option>
                                </select>
                            </label>
                            <label>
                                {transactionForm.type === 'transfer' ? 'From account' : 'Account'}
                                <select name="accountId" value={transactionForm.accountId} onChange={updateTransactionField} required>
                                    <option value="">Choose an account</option>
                                    {accounts.map((account) => (
                                        <option key={account.id} value={account.id}>
                                            {account.name} ({account.currency})
                                        </option>
                                    ))}
                                </select>
                            </label>
                            {transactionForm.type === 'transfer' && (
                                <label>
                                    To account
                                    <select
                                        name="transferAccountId"
                                        value={transactionForm.transferAccountId}
                                        onChange={updateTransactionField}
                                        required
                                    >
                                        <option value="">Choose a destination</option>
                                        {transferAccounts.map((account) => (
                                            <option key={account.id} value={account.id}>{account.name}</option>
                                        ))}
                                    </select>
                                </label>
                            )}
                            {transactionForm.type !== 'transfer' && (
                                <label>
                                    Category
                                    <select name="categoryId" value={transactionForm.categoryId} onChange={updateTransactionField}>
                                        <option value="">Uncategorized</option>
                                        {categories.map((category) => (
                                            <option key={category.id} value={category.id}>{category.name}</option>
                                        ))}
                                    </select>
                                </label>
                            )}
                            <label>
                                Amount {selectedAccount ? `(${selectedAccount.currency})` : ''}
                                <input
                                    name="amount"
                                    inputMode="decimal"
                                    value={transactionForm.amount}
                                    onChange={updateTransactionField}
                                    placeholder="0.00"
                                    required
                                />
                            </label>
                            <label>
                                Date
                                <input name="date" type="date" value={transactionForm.date} onChange={updateTransactionField} required />
                            </label>
                            <label>
                                Payee
                                <input name="payee" value={transactionForm.payee} onChange={updateTransactionField} maxLength="180" />
                            </label>
                            <label>
                                Note
                                <input name="description" value={transactionForm.description} onChange={updateTransactionField} maxLength="500" />
                            </label>
                            <button type="submit" className="ledger-primary-button">Save transaction</button>
                        </form>
                    )}
                </section>
            </div>

            <section aria-labelledby="recent-transactions-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div>
                        <p className="eyebrow">ACTIVITY</p>
                        <h2 id="recent-transactions-heading">Recent transactions</h2>
                    </div>
                </div>
                {transactions.length === 0 ? (
                    <p className="ledger-empty">Your transactions will appear here as you add them.</p>
                ) : (
                    <div className="ledger-table-wrap">
                        <table className="ledger-table">
                            <thead>
                                <tr><th scope="col">Date</th><th scope="col">Description</th><th scope="col">Account</th><th scope="col">Type</th><th scope="col">Amount</th></tr>
                            </thead>
                            <tbody>
                                {transactions.map((transaction) => {
                                    const account = accounts.find((item) => item.id === transaction.account_id);
                                    const destination = accounts.find((item) => item.id === transaction.transfer_account_id);
                                    const label = transaction.payee || transaction.description || 'Transaction';
                                    return (
                                        <tr key={transaction.id}>
                                            <td>{transaction.date}</td>
                                            <td>{transaction.type === 'transfer' ? `${account?.name || 'Account'} → ${destination?.name || 'Account'}` : label}</td>
                                            <td>{account?.name || '—'}</td>
                                            <td>{transaction.type}</td>
                                            <td>{formatAmountMinor(transaction.amount_minor, transaction.currency)}</td>
                                        </tr>
                                    );
                                })}
                            </tbody>
                        </table>
                    </div>
                )}
            </section>

            <section aria-labelledby="category-rules-heading" className="ledger-section">
                <div className="ledger-section-heading">
                    <div>
                        <p className="eyebrow">AUTOMATION</p>
                        <h2 id="category-rules-heading">Merchant rules</h2>
                    </div>
                </div>
                <div className="ledger-forms">
                    <form onSubmit={submitCategoryRule} className="ledger-panel ledger-form">
                        <h3>Create a rule</h3>
                        <label>
                            Rule name
                            <input
                                value={ruleForm.name}
                                onChange={(event) => setRuleForm({ ...ruleForm, name: event.target.value })}
                                maxLength="120"
                                required
                            />
                        </label>
                        <label>
                            Payee contains
                            <input
                                value={ruleForm.matchValue}
                                onChange={(event) => setRuleForm({ ...ruleForm, matchValue: event.target.value })}
                                maxLength="180"
                                required
                            />
                        </label>
                        <label>
                            Apply category
                            <select
                                value={ruleForm.categoryId}
                                onChange={(event) => setRuleForm({ ...ruleForm, categoryId: event.target.value })}
                                required
                            >
                                <option value="">Choose a category</option>
                                {categories.map((category) => (
                                    <option key={category.id} value={category.id}>{category.name}</option>
                                ))}
                            </select>
                        </label>
                        <button type="submit" className="ledger-primary-button" disabled={categories.length === 0}>
                            Save rule
                        </button>
                    </form>
                    <div className="ledger-panel">
                        <h3>Active rules</h3>
                        {categoryRules.length === 0 ? (
                            <p className="ledger-empty">No merchant rules yet. Rules apply automatically to new uncategorized transactions.</p>
                        ) : (
                            <ul className="ledger-rule-list">
                                {categoryRules.map((rule) => (
                                    <li key={rule.id}>
                                        <span>
                                            <strong>{rule.name}</strong>
                                            <small>Payee contains “{rule.match_value}” · {categories.find((category) => category.id === rule.category_id)?.name || 'Category'}</small>
                                        </span>
                                        <button type="button" onClick={() => deleteCategoryRule(rule.id)} aria-label={`Delete ${rule.name}`}>
                                            Delete
                                        </button>
                                    </li>
                                ))}
                            </ul>
                        )}
                    </div>
                </div>
            </section>
        </main>
    );
};

export default LedgerPage;
