// src/App.js

import React from 'react';
import { BrowserRouter as Router, Routes, Route, NavLink } from 'react-router-dom';
import ExpenseList from './components/Expense/ExpenseList';
import './App.css';
import ExpenseForm from "./components/Expense/ExpenseForm";
import InvestmentList from "./components/Investment/InvestmentList";
import InvestmentForm from "./components/Investment/InvestmentForm"

function App() {
    return (
        <Router>
            <div className="xpense-trackr-main">
                <h1>Xpense-Trackr</h1>
                <nav>
                    <ul className="horizontal-tabs">
                        <li><NavLink to="/expenses" activeClassName="active-tab">Expense</NavLink></li>
                        <li><NavLink to="/investments" activeClassName="active-tab">Investments</NavLink></li>
                        <li><NavLink to="/budgets" activeClassName="active-tab">Plan your Budget</NavLink></li>
                        <li><NavLink to="/goals" activeClassName="active-tab">Investment Goals</NavLink></li>
                        <li><NavLink to="/trends" activeClassName="active-tab">Trends</NavLink></li>
                    </ul>
                </nav>
                <Routes>
                    <Route path="/expenses" element={<ExpenseList />} />
                    <Route path="/expenses/new" element={<ExpenseForm />} />
                    <Route path="/investments" element={<InvestmentList />} />
                    <Route path="/investments/new" element={<InvestmentForm />} />
                    {/*<Route path="/categories" element={<CategoryList />} />*/}
                    {/*<Route path="/budgets" element={<BudgetList />} />*/}
                    {/*<Route path="/trends" element={<TrendsPage />} />*/}
                </Routes>
            </div>
        </Router>
    );
}

export default App;
