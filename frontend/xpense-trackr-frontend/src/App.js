// src/App.js

import { lazy, Suspense, useEffect, useState } from 'react';
import {BrowserRouter as Router, Routes, Route, NavLink, Link, useLocation} from 'react-router-dom';
import './App.css';
import { faHouse, faThumbtack} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import OfflineNotice from "./components/Common/OfflineNotice";

const ExpenseList = lazy(() => import('./components/Expense/ExpenseList'));
const ExpenseForm = lazy(() => import('./components/Expense/ExpenseForm'));
const InvestmentList = lazy(() => import('./components/Investment/InvestmentList'));
const InvestmentForm = lazy(() => import('./components/Investment/InvestmentForm'));
const Homepage = lazy(() => import('./components/Homepage/Homepage'));
const BudgetList = lazy(() => import('./components/Budget/BudgetList'));
const BudgetForm = lazy(() => import('./components/Budget/BudgetForm'));
const GoalList = lazy(() => import('./components/Goal/GoalList'));
const GoalForm = lazy(() => import('./components/Goal/GoalForm'));
const TrendsPage = lazy(() => import('./components/Trends/TrendsPage'));
const CategoryList = lazy(() => import('./components/Category/CategoryList'));
const CategoryForm = lazy(() => import('./components/Category/CategoryForm'));
const RecurringPage = lazy(() => import('./components/Recurring/RecurringPage'));
const RecurringForm = lazy(() => import('./components/Recurring/RecurringForm'));
const NotFoundPage = lazy(() => import('./components/Common/NotFoundPage'));
const LedgerPage = lazy(() => import('./components/Ledger/LedgerPage'));
const ForecastPage = lazy(() => import('./components/Forecast/ForecastPage'));
const ReceiptPage = lazy(() => import('./components/Receipt/ReceiptPage'));
const ReportsPage = lazy(() => import('./components/Reports/ReportsPage'));
const ImportExportPage = lazy(() => import('./components/ImportExport/ImportExportPage'));
const PlanningPage = lazy(() => import('./components/Planning/PlanningPage'));
const BillsPage = lazy(() => import('./components/Bills/BillsPage'));
const AssetsPage = lazy(() => import('./components/Assets/AssetsPage'));
const PortfolioPage = lazy(() => import('./components/Portfolio/PortfolioPage'));

const navigationGroups = [
    {
        label: 'Activity',
        links: [
            ['Expenses', '/expenses'],
            ['Transactions', '/transactions'],
            ['Bills', '/bills'],
            ['Recurring', '/recurring'],
            ['Receipts', '/receipts'],
            ['Categories', '/categories'],
        ],
    },
    {
        label: 'Planning',
        links: [
            ['Cash flow', '/forecast'],
            ['Budgets', '/budgets'],
            ['Goals', '/goals'],
            ['Plans', '/planning'],
        ],
    },
    {
        label: 'Wealth',
        links: [
            ['Investments', '/investments'],
            ['Portfolio', '/portfolio'],
            ['Assets', '/assets'],
        ],
    },
    {
        label: 'More',
        links: [
            ['Reports', '/reports'],
            ['Trends', '/trends'],
            ['Import / export', '/import-export'],
        ],
    },
];

function MainNavigation() {
    const { pathname } = useLocation();
    const [openGroup, setOpenGroup] = useState(null);

    useEffect(() => {
        setOpenGroup(null);
    }, [pathname]);

    return (
        <nav className="app-navigation" aria-label="Main navigation">
            <NavLink to="/" end className={({ isActive }) => `overview-link${isActive ? ' active-tab' : ''}`}>
                <FontAwesomeIcon icon={faHouse}/> Overview
            </NavLink>
            {navigationGroups.map((group) => {
                const isGroupActive = group.links.some(([, path]) => pathname.startsWith(path));
                return (
                    <details
                        className={`navigation-group${isGroupActive ? ' group-active' : ''}`}
                        key={group.label}
                        open={openGroup === group.label}
                    >
                        <summary onClick={(event) => {
                            event.preventDefault();
                            setOpenGroup((current) => current === group.label ? null : group.label);
                        }}>
                            {group.label}
                        </summary>
                        <div className="navigation-menu">
                            {group.links.map(([label, path]) => (
                                <NavLink
                                    key={path}
                                    to={path}
                                    className={({ isActive }) => isActive ? 'active-tab' : undefined}
                                >
                                    {label}
                                </NavLink>
                            ))}
                        </div>
                    </details>
                );
            })}
        </nav>
    );
}

function App() {
    return (
        <Router>
            <div className="xpense-trackr-main">
                <header className="app-header">
                    <Link to="/" className="app-brand">
                        <span className="brand-mark"><FontAwesomeIcon icon={faThumbtack}/></span>
                        <span>
                            <span className="brand-eyebrow">PERSONAL FINANCE</span>
                            <h1>Expense Tracker</h1>
                        </span>
                    </Link>
                    <MainNavigation/>
                </header>
                <OfflineNotice />
                <Suspense fallback={<main className="ledger-page"><p role="status">Loading page…</p></main>}>
                    <Routes>
                        <Route path="/" element={<Homepage/>}/>
                        <Route path="/expenses" element={<ExpenseList/>}/>
                        <Route path="/expenses/new" element={<ExpenseForm/>}/>
                        <Route path="/transactions" element={<LedgerPage/>}/>
                        <Route path="/forecast" element={<ForecastPage/>}/>
                        <Route path="/bills" element={<BillsPage/>}/>
                        <Route path="/assets" element={<AssetsPage/>}/>
                        <Route path="/portfolio" element={<PortfolioPage/>}/>
                        <Route path="/receipts" element={<ReceiptPage/>}/>
                        <Route path="/reports" element={<ReportsPage/>}/>
                        <Route path="/import-export" element={<ImportExportPage/>}/>
                        <Route path="/investments" element={<InvestmentList/>}/>
                        <Route path="/investments/new" element={<InvestmentForm/>}/>
                        <Route path="/budgets" element={<BudgetList/>}/>
                        <Route path="/budgets/new" element={<BudgetForm/>}/>
                        <Route path="/budgets/:id/edit" element={<BudgetForm/>}/>
                        <Route path="/goals" element={<GoalList/>}/>
                        <Route path="/goals/new" element={<GoalForm/>}/>
                        <Route path="/goals/:id/edit" element={<GoalForm/>}/>
                        <Route path="/planning" element={<PlanningPage/>}/>
                        <Route path="/trends" element={<TrendsPage/>}/>
                        <Route path="/recurring" element={<RecurringPage/>}/>
                        <Route path="/recurring/:kind/:id/edit" element={<RecurringForm/>}/>
                        <Route path="/categories" element={<CategoryList/>}/>
                        <Route path="/categories/new" element={<CategoryForm/>}/>
                        <Route path="/categories/:id/edit" element={<CategoryForm/>}/>
                        <Route path="*" element={<NotFoundPage/>}/>
                    </Routes>
                </Suspense>
            </div>
        </Router>
    );
}

export default App;
