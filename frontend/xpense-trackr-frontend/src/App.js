// src/App.js

import React from 'react';
import { BrowserRouter as Router, Routes, Route, NavLink } from 'react-router-dom';
import ExpenseList from './components/Expense/ExpenseList';
import './App.css';

function App() {
    return (
        <Router>
            <div>
                <h1>Xpense-Trackr</h1>
                <nav>
                    <ul className="horizontal-tabs">
                        <li><NavLink to="/expenses" activeClassName="active-tab">Expense</NavLink></li>
                        <li><NavLink to="/categories" activeClassName="active-tab">Category</NavLink></li>
                        <li><NavLink to="/budgets" activeClassName="active-tab">Budget</NavLink></li>
                        <li><NavLink to="/trends" activeClassName="active-tab">Trends</NavLink></li>
                    </ul>
                </nav>
                <Routes>
                    <Route path="/expenses" element={<ExpenseList />} />
                    {/*<Route path="/categories" element={<CategoryList />} />*/}
                    {/*<Route path="/budgets" element={<BudgetList />} />*/}
                    {/*<Route path="/trends" element={<TrendsPage />} />*/}
                </Routes>
            </div>
        </Router>
    );
}

export default App;
