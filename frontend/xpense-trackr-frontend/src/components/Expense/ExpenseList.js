// src/components/Expense/ExpenseList.js

import React, { useState, useEffect } from 'react';
import axios from 'axios';
import backendUrl from "../../config";
import {Link} from "react-router-dom";
import './ExpenseList.css';

const ExpenseList = () => {
    const [expenses, setExpenses] = useState([]);

    useEffect(() => {
        fetchExpenses();
    }, []);

    const fetchExpenses = async () => {
        try {
            const response = await axios.get(`${backendUrl}/api/expenses`);
            setExpenses(response.data);
        } catch (error) {
            console.error('Error fetching expenses:', error);
        }
    };

    const handleEditExpense = async (expenseId) => {
        try {
            // Fetch the expense data for the selected expenseId
            const response = await axios.get(`${process.env.REACT_APP_BACKEND_URL}/api/expenses/${expenseId}`);
            const expenseData = response.data;

            // Call the PUT endpoint to update the expense
            await axios.put(`${backendUrl}/api/expenses/${expenseId}`, expenseData);
            alert('Expense updated successfully');

            // Optionally, fetch the list of expenses to update the UI
            await fetchExpenses();
        } catch (error) {
            console.error('Error updating expense:', error);
            alert('Error updating expense. Please try again.');
        }
    };

    const handleDeleteExpense = async (expenseId) => {
        try {
            await axios.delete(`${backendUrl}/api/expenses/${expenseId}`);
            alert('Expense deleted successfully');
            await fetchExpenses(); // Fetch updated list of expenses after deletion
        } catch (error) {
            console.error('Error deleting expense:', error);
            alert('Error deleting expense. Please try again.');
        }
    };

    const formatDate = (dateString) => {
        const options = { year: 'numeric', month: 'short', day: 'numeric' };
        return new Date(dateString).toLocaleDateString('en-US', options);
    };

    return (
        <div className="expense-list-container">
            <div className="expense-list-header">
                <h2>Expense List</h2>
                {/* Add button to navigate to ExpenseForm page */}
                <Link to="/expenses/new" className="button-common">Add New Expense</Link>
            </div>
            <table className="expense-list-table">
                <thead>
                <tr>
                    <th>Description</th>
                    <th>Amount</th>
                    <th>Date</th>
                    <th colSpan="2">Action</th>
                </tr>
                </thead>
                <tbody>
                {expenses.map(expense => (
                    <tr key={expense.id}>
                        <td>{expense.description}</td>
                        <td>{expense.amount}</td>
                        <td>{formatDate(expense.date)}</td>
                        <td>
                            <button className="button-common" onClick={() => handleEditExpense(expense.id)}>Edit</button>
                        </td>
                        <td>
                            <button className="button-common" onClick={() => handleDeleteExpense(expense.id)}>Delete</button>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
        </div>
    );
};

export default ExpenseList;
