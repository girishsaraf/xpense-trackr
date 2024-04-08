// src/components/Expense/ExpenseForm.js

import React, { useState } from 'react';
import axios from 'axios';
import backendUrl from "../../config";
import './ExpenseForm.css';
import { useNavigate } from 'react-router-dom';


const ExpenseForm = () => {

    const navigate = useNavigate();
    const [formData, setFormData] = useState({
        description: '',
        amount: '',
        date: '',
        category_id: '',
        recurring: false,
        start_date: '',
        end_date: '',
        frequency: '',
    });

    const handleSubmit = async event => {
        event.preventDefault();

        let isRecurring = formData.recurring;

        try
        {
            if(isRecurring){
                await axios.post(`${backendUrl}/api/recurring/expenses`, formData);
                alert('Expense added successfully');
            }
            else {
                await axios.post(`${backendUrl}/api/expenses`, formData);
                alert('Expense added successfully');
            }
            clearFormData();

            navigate('/expenses');
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
            recurring: false,
            start_date: '',
            end_date: '',
            frequency: '',
        });
        navigate('/expenses');
    }

    const handleInputChange = event => {
        const { name, value, type, checked } = event.target;
        const newValue = type === 'checkbox' ? checked : value;

        setFormData({
            ...formData,
            [name]: newValue,
        });
    };

    return (
        <div className="expense-form-container">
            <h2>Add Expense</h2>
            <form onSubmit={handleSubmit}>
                <div className="input-container">
                    <label htmlFor="description">Description:</label>
                    <input
                        type="text"
                        id="description"
                        name="description"
                        value={formData.description}
                        onChange={handleInputChange}
                    />
                </div>
                <div className="input-container">
                    <label htmlFor="amount">Amount:</label>
                    <input
                        type="number"
                        id="amount"
                        name="amount"
                        value={formData.amount}
                        onChange={handleInputChange}
                    />
                </div>
                <div className="input-container">
                    <label htmlFor="category">Category:</label>
                    <input
                        type="text"
                        id="category_id"
                        name="category_id"
                        value={formData.category_id}
                        onChange={handleInputChange}
                    />
                </div>
                <div className="input-container">
                    <label htmlFor="date">Date:</label>
                    <input
                        type="date"
                        id="date"
                        name="date"
                        value={formData.date}
                        onChange={handleInputChange}
                    />
                </div>
                <div className="input-container">
                    <label htmlFor="recurring">Is Recurring Expense?</label>
                    <input
                        type="checkbox"
                        id="recurring"
                        name="recurring"
                        value={formData.recurring}
                        onChange={handleInputChange}
                    />
                </div>
                {formData.recurring && (
                    <>
                        <div className="input-container">
                            <label htmlFor="startDate">Start Date:</label>
                            <input
                                type="date"
                                id="startDate"
                                name="startDate"
                                value={formData.start_date}
                                defaultValue={formData.date}
                                onChange={handleInputChange}
                            />
                        </div>
                        <div className="input-container">
                            <label htmlFor="endDate">End Date:</label>
                            <input
                                type="date"
                                id="endDate"
                                name="endDate"
                                value={formData.end_date}
                                defaultValue={formData.date}
                                onChange={handleInputChange}
                            />
                        </div>
                        <div className="input-container">
                            <label htmlFor="frequency">Frequency:</label>
                            <select
                                id="frequency"
                                name="frequency"
                                value={formData.frequency}
                                onChange={handleInputChange}
                            >
                                <option value="daily">Daily</option>
                                <option value="weekly">Weekly</option>
                                <option value="biweekly">Bi-Weekly</option>
                                <option value="monthly">Monthly</option>
                            </select>
                        </div>
                    </>
                )}
                <div className="buttons-container">
                    <button type="submit">Add</button>
                    <button type="button" onClick={clearFormData}>Cancel</button>
                </div>
            </form>
        </div>
    );

};

export default ExpenseForm;
