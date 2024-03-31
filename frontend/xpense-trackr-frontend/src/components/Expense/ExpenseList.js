// src/components/Expense/ExpenseList.js

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import backendUrl from "../../config";

const ExpenseList = () => {
    const [expenses, setExpenses] = useState([]);

    useEffect(() => {
        fetchExpenses();
    }, []);

    const fetchExpenses = async () => {
        try {
            const response = await axios.get(`$${backendUrl}/api/expenses`);
            setExpenses(response.data);
        } catch (error) {
            console.error('Error fetching expenses:', error);
        }
    };

    return (
        <div>
            <h2>Expense List</h2>
            <ul>
                {expenses.map(expense => (
                    <li key={expense.id}>
                        <div>{expense.title}</div>
                        <div>{expense.amount}</div>
                    </li>
                ))}
            </ul>
        </div>
    );
};

export default ExpenseList;
