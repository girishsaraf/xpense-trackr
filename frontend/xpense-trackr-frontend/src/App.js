// src/App.js

import { lazy, Suspense } from 'react';
import {BrowserRouter as Router, Routes, Route, NavLink, Link} from 'react-router-dom';
import './App.css';
import { faArrowTrendUp, faChartPie, faHouse, faThumbtack} from "@fortawesome/free-solid-svg-icons";
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
                    <nav aria-label="Main navigation">
                        <ul className="horizontal-tabs">
                            <li><NavLink to="/" end className={({ isActive }) => isActive ? 'active-tab' : undefined}><FontAwesomeIcon icon={faHouse}/> Overview</NavLink></li>
                            <li><NavLink to="/expenses" className={({ isActive }) => isActive ? 'active-tab' : undefined}><FontAwesomeIcon icon={faChartPie}/> Expense</NavLink></li>
                            <li><NavLink to="/transactions" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Transactions</NavLink></li>
                            <li><NavLink to="/forecast" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Cash flow</NavLink></li>
                            <li><NavLink to="/bills" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Bills</NavLink></li>
                            <li><NavLink to="/assets" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Assets</NavLink></li>
                            <li><NavLink to="/portfolio" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Portfolio</NavLink></li>
                            <li><NavLink to="/receipts" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Receipts</NavLink></li>
                            <li><NavLink to="/reports" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Reports</NavLink></li>
                            <li><NavLink to="/import-export" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Import / export</NavLink></li>
                            <li><NavLink to="/investments" className={({ isActive }) => isActive ? 'active-tab' : undefined}><FontAwesomeIcon icon={faArrowTrendUp}/> Investments</NavLink></li>
                            <li><NavLink to="/budgets" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Budgets</NavLink></li>
                            <li><NavLink to="/goals" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Goals</NavLink></li>
                            <li><NavLink to="/planning" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Plans</NavLink></li>
                            <li><NavLink to="/trends" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Trends</NavLink></li>
                            <li><NavLink to="/recurring" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Recurring</NavLink></li>
                            <li><NavLink to="/categories" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Categories</NavLink></li>
                        </ul>
                    </nav>
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
