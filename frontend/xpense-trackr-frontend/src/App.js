// src/App.js

import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ExpenseList from './components/Expense/ExpenseList';
import ExpenseForm from './components/Expense/ExpenseForm';

function App() {
    return (
        <Router>
            <div>
                <h1>Expense Tracking App</h1>
                <Routes>
                    <Route path="/expenses" element={<ExpenseList />} />
                    <Route path="/expenses/new" element={<ExpenseForm />} />
                    {/* Other routes for different components */}
                </Routes>
            </div>
        </Router>
    );
}

export default App;
