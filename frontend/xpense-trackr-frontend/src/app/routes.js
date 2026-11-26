import { lazy } from 'react';

const ExpenseList = lazy(() => import('../components/Expense/ExpenseList'));
const ExpenseForm = lazy(() => import('../components/Expense/ExpenseForm'));
const InvestmentList = lazy(() => import('../components/Investment/InvestmentList'));
const InvestmentForm = lazy(() => import('../components/Investment/InvestmentForm'));
const Homepage = lazy(() => import('../components/Homepage/Homepage'));
const BudgetList = lazy(() => import('../components/Budget/BudgetList'));
const BudgetForm = lazy(() => import('../components/Budget/BudgetForm'));
const GoalList = lazy(() => import('../components/Goal/GoalList'));
const GoalForm = lazy(() => import('../components/Goal/GoalForm'));
const TrendsPage = lazy(() => import('../components/Trends/TrendsPage'));
const CategoryList = lazy(() => import('../components/Category/CategoryList'));
const CategoryForm = lazy(() => import('../components/Category/CategoryForm'));
const RecurringPage = lazy(() => import('../components/Recurring/RecurringPage'));
const RecurringForm = lazy(() => import('../components/Recurring/RecurringForm'));
const NotFoundPage = lazy(() => import('../components/Common/NotFoundPage'));
const LedgerPage = lazy(() => import('../components/Ledger/LedgerPage'));
const ForecastPage = lazy(() => import('../components/Forecast/ForecastPage'));
const ReceiptPage = lazy(() => import('../components/Receipt/ReceiptPage'));
const ReportsPage = lazy(() => import('../components/Reports/ReportsPage'));
const ImportExportPage = lazy(() => import('../components/ImportExport/ImportExportPage'));
const PlanningPage = lazy(() => import('../components/Planning/PlanningPage'));
const BillsPage = lazy(() => import('../components/Bills/BillsPage'));
const AssetsPage = lazy(() => import('../components/Assets/AssetsPage'));
const PortfolioPage = lazy(() => import('../components/Portfolio/PortfolioPage'));

export const appRoutes = [
    ['/', Homepage],
    ['/expenses', ExpenseList],
    ['/expenses/new', ExpenseForm],
    ['/transactions', LedgerPage],
    ['/forecast', ForecastPage],
    ['/bills', BillsPage],
    ['/assets', AssetsPage],
    ['/portfolio', PortfolioPage],
    ['/receipts', ReceiptPage],
    ['/reports', ReportsPage],
    ['/import-export', ImportExportPage],
    ['/investments', InvestmentList],
    ['/investments/new', InvestmentForm],
    ['/budgets', BudgetList],
    ['/budgets/new', BudgetForm],
    ['/budgets/:id/edit', BudgetForm],
    ['/goals', GoalList],
    ['/goals/new', GoalForm],
    ['/goals/:id/edit', GoalForm],
    ['/planning', PlanningPage],
    ['/trends', TrendsPage],
    ['/recurring', RecurringPage],
    ['/recurring/:kind/:id/edit', RecurringForm],
    ['/categories', CategoryList],
    ['/categories/new', CategoryForm],
    ['/categories/:id/edit', CategoryForm],
    ['*', NotFoundPage],
];
