// src/components/Expense/ExpenseForm.js

import React, { useState } from 'react';
import axios from 'axios';
import backendUrl from "../../config";
import './InvestmentForm.css';
import { useNavigate } from 'react-router-dom';


const InvestmentForm = () => {

    const navigate = useNavigate();

    const [formData, setFormData] = useState({
        type: '',
        amount: '',
        date: '',
    });

    const handleSubmit = async event => {
        event.preventDefault();

        try {
            await axios.post(`${backendUrl}/api/investments`, formData);
            alert('Investment added successfully');
            clearFormData();

            navigate('/investments');
        } catch (error) {
            console.error('Error adding expense:', error);
            alert('Error adding investment. Please try again.');
        }
    };

    const clearFormData = () => {
        setFormData({
            type: '',
            amount: '',
            date: '',
        });
        navigate('/investments');
    }

    const handleInputChange = event => {
        setFormData({
            ...formData,
            [event.target.name]: event.target.value,
        });
    };

    return (
        <div className="expense-form-container">
            <h2>Add New Investment</h2>
            <form onSubmit={handleSubmit}>
                <div>
                    <label htmlFor="type">Investment Type:</label>
                    <input
                        type="text"
                        id="type"
                        name="type"
                        value={formData.type}
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

export default InvestmentForm;
