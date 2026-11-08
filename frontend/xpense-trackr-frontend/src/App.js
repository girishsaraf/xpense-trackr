// src/App.js

import React from 'react';
import {BrowserRouter as Router, Routes, Route, NavLink, Link} from 'react-router-dom';
import ExpenseList from './components/Expense/ExpenseList';
import './App.css';
import ExpenseForm from "./components/Expense/ExpenseForm";
import InvestmentList from "./components/Investment/InvestmentList";
import InvestmentForm from "./components/Investment/InvestmentForm"
import { faArrowTrendUp, faChartPie, faHouse, faThumbtack} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import Homepage from "./components/Homepage/Homepage";
import BudgetList from "./components/Budget/BudgetList";
import BudgetForm from "./components/Budget/BudgetForm";
import GoalList from "./components/Goal/GoalList";
import GoalForm from "./components/Goal/GoalForm";
import TrendsPage from "./components/Trends/TrendsPage";
import CategoryList from "./components/Category/CategoryList";
import CategoryForm from "./components/Category/CategoryForm";
import RecurringPage from "./components/Recurring/RecurringPage";
import RecurringForm from "./components/Recurring/RecurringForm";
import NotFoundPage from "./components/Common/NotFoundPage";
import LedgerPage from "./components/Ledger/LedgerPage";
import ForecastPage from "./components/Forecast/ForecastPage";
import ReceiptPage from "./components/Receipt/ReceiptPage";
import ReportsPage from "./components/Reports/ReportsPage";
import ImportExportPage from "./components/ImportExport/ImportExportPage";

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
                            <li><NavLink to="/receipts" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Receipts</NavLink></li>
                            <li><NavLink to="/reports" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Reports</NavLink></li>
                            <li><NavLink to="/import-export" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Import / export</NavLink></li>
                            <li><NavLink to="/investments" className={({ isActive }) => isActive ? 'active-tab' : undefined}><FontAwesomeIcon icon={faArrowTrendUp}/> Investments</NavLink></li>
                            <li><NavLink to="/budgets" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Budgets</NavLink></li>
                            <li><NavLink to="/goals" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Goals</NavLink></li>
                            <li><NavLink to="/trends" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Trends</NavLink></li>
                            <li><NavLink to="/recurring" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Recurring</NavLink></li>
                            <li><NavLink to="/categories" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Categories</NavLink></li>
                        </ul>
                    </nav>
                </header>
                <Routes>
                    <Route path="/" element={<Homepage/>}/>
                    <Route path="/expenses" element={<ExpenseList/>}/>
                    <Route path="/expenses/new" element={<ExpenseForm/>}/>
                    <Route path="/transactions" element={<LedgerPage/>}/>
                    <Route path="/forecast" element={<ForecastPage/>}/>
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
                    <Route path="/trends" element={<TrendsPage/>}/>
                    <Route path="/recurring" element={<RecurringPage/>}/>
                    <Route path="/recurring/:kind/:id/edit" element={<RecurringForm/>}/>
                    <Route path="/categories" element={<CategoryList/>}/>
                    <Route path="/categories/new" element={<CategoryForm/>}/>
                    <Route path="/categories/:id/edit" element={<CategoryForm/>}/>
                    <Route path="*" element={<NotFoundPage/>}/>
                </Routes>
            </div>
        </Router>
    );
}

export default App;
