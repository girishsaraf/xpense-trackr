// src/components/Expense/ExpenseForm.js

import React, { useState } from 'react';
import axios from 'axios';
import backendUrl from "../../config";
import './ExpenseForm.css';

const ExpenseForm = () => {
    const [formData, setFormData] = useState({
        description: '',
        amount: '',
        date: '',
        category_id: '',
        user_id: '',
    });

    const handleSubmit = async event => {
        event.preventDefault();

        try {
            await axios.post(`${backendUrl}/api/expenses`, formData);
            alert('Expense added successfully');
            clearFormData();
        } catch (error) {
            console.error('Error adding expense:', error);
            alert('Error adding expense. Please try again.');
        }
    };

    const clearFormData = () => {
        setFormData({
            description: '',
            amount: '',
            date: '',
            category_id: '',
            user_id: '',
        });
    }

    const handleInputChange = event => {
        setFormData({
            ...formData,
            [event.target.name]: event.target.value,
        });
    };

    return (
        <div className="expense-form-container">
            <h2>Add Expense</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <label htmlFor="user_id">User:</label>
                    <input
                        type="number"
                        id="user_id"
                        name="user_id"
                        value={formData.user_id}
                        onChange={handleInputChange}
                    />
                </div>
                <div>
                    <label htmlFor="description">Description:</label>
                    <input
                        type="text"
                        id="description"
                        name="description"
                        value={formData.description}
                        onChange={handleInputChange}
                    />
                </div>
                <div>
                    <label htmlFor="amount">Amount:</label>
                    <input
                        type="number"
                        id="amount"
                        name="amount"
                        value={formData.amount}
                        onChange={handleInputChange}
                    />
                </div>
                <div>
                    <label htmlFor="category">Category:</label>
                    <input
                        type="text"
                        id="category_id"
                        name="category_id"
                        value={formData.category_id}
                        onChange={handleInputChange}
                    />
                </div>
                <div>
                    <label htmlFor="date">Date:</label>
                    <input
                        type="date"
                        id="date"
                        name="date"
                        value={formData.date}
                        onChange={handleInputChange}
                    />
                </div>
                <button type="submit">Add Expense</button>
                <button type="cancel" onClick={clearFormData}>Cancel</button>
            </form>
        </div>
    );
};

export default ExpenseForm;
