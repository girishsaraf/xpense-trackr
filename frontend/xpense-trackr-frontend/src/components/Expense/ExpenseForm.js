// src/components/Expense/ExpenseForm.js

import React, { useState } from 'react';
import axios from 'axios';
import backendUrl from "../../config";

const ExpenseForm = () => {
    const [formData, setFormData] = useState({
        title: '',
        amount: '',
    });

    const handleSubmit = async event => {
        event.preventDefault();

        try {
            await axios.post(`${backendUrl}/api/expenses`, formData);
            alert('Expense added successfully');
            setFormData({
                title: '',
                amount: '',
            });
        } catch (error) {
            console.error('Error adding expense:', error);
            alert('Error adding expense. Please try again.');
        }
    };

    const handleInputChange = event => {
        setFormData({
            ...formData,
            [event.target.name]: event.target.value,
        });
    };

    return (
        <div>
            <h2>Add Expense</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <label htmlFor="title">Title:</label>
                    <input
                        type="text"
                        id="title"
                        name="title"
                        value={formData.title}
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
                <button type="submit">Add Expense</button>
            </form>
        </div>
    );
};

export default ExpenseForm;
