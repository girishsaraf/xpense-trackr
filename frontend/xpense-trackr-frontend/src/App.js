// src/App.js

import React from 'react';
import {BrowserRouter as Router, Routes, Route, NavLink, Link} from 'react-router-dom';
import ExpenseList from './components/Expense/ExpenseList';
import './App.css';
import ExpenseForm from "./components/Expense/ExpenseForm";
import InvestmentList from "./components/Investment/InvestmentList";
import InvestmentForm from "./components/Investment/InvestmentForm"
import { faThumbtack} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import Homepage from "./components/Homepage/Homepage";

function App() {
    return (
        <Router>
            <div className="xpense-trackr-main">
                <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
                    <h1><FontAwesomeIcon icon={faThumbtack}/> Expense Tracker</h1>
                </Link>
                <nav>
                    <ul className="horizontal-tabs">
                        <li><NavLink to="/expenses" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Expense</NavLink></li>
                        <li><NavLink to="/investments" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Investments</NavLink></li>
                        <li><NavLink to="/budgets" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Plan your Budget</NavLink></li>
                        <li><NavLink to="/goals" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Investment Goals</NavLink></li>
                        <li><NavLink to="/trends" className={({ isActive }) => isActive ? 'active-tab' : undefined}>Trends</NavLink></li>
                    </ul>
                </nav>
                <Routes>
                    <Route path="/" element={<Homepage/>}/>
                    <Route path="/expenses" element={<ExpenseList/>}/>
                    <Route path="/expenses/new" element={<ExpenseForm/>}/>
                    <Route path="/investments" element={<InvestmentList/>}/>
                    <Route path="/investments/new" element={<InvestmentForm/>}/>
                    {/*<Route path="/categories" element={<CategoryList />} />*/}
                    {/*<Route path="/budgets" element={<BudgetList />} />*/}
                    {/*<Route path="/trends" element={<TrendsPage />} />*/}
                </Routes>
            </div>
        </Router>
    );
}

export default App;
