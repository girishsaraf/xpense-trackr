export type CashflowCategory = {
    category_id?: number | null;
    category_name: string;
    amount_minor: number;
};

export type CashflowCurrency = {
    currency: string;
    income_minor: number;
    expense_minor: number;
    net_minor: number;
    categories: CashflowCategory[];
};

export type CashflowReport = {
    from: string;
    to: string;
    currencies: CashflowCurrency[];
};

export type PortfolioHoldingItem = {
    holding_id: number;
    name: string;
    ticker: string;
    asset_class: string;
    quantity: string;
    cost_basis_minor: number;
    value_minor: number;
    as_of: string;
    source: string;
};

export type NetWorthAccount = {
    account_id: number;
    name: string;
    type: string;
    balance_minor: number;
};

export type NetWorthItem = {
    asset_id: number;
    name: string;
    kind: string;
    value_minor: number;
    as_of: string;
    source: string;
};

export type NetWorthCurrency = {
    currency: string;
    assets_minor: number;
    liabilities_minor: number;
    net_worth_minor: number;
    portfolio_minor?: number;
    accounts: NetWorthAccount[];
    manual_items?: NetWorthItem[];
    holdings?: PortfolioHoldingItem[];
};

export type NetWorthReport = {
    as_of: string;
    currencies: NetWorthCurrency[];
};

export type NetWorthHistoryReport = {
    months: number;
    entries: Array<{
        as_of: string;
        currencies: NetWorthCurrency[];
    }>;
};

export type FinancialInsight = {
    kind: string;
    title: string;
    explanation: string;
    currency: string;
    category_name?: string;
    current_amount_minor: number;
    previous_amount_minor: number;
    change_percent?: number;
};

export type FinancialInsightsReport = {
    from: string;
    to: string;
    previous_from: string;
    previous_to: string;
    insights: FinancialInsight[];
};
