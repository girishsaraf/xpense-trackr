// src/components/Expense/ExpenseForm.js

import React, { useState } from 'react';
import axios from 'axios';
import backendUrl from "../../config";
import './InvestmentForm.css';
import { useNavigate } from 'react-router-dom';
import Notification from "../Notification/Notification";


const InvestmentForm = () => {

    const navigate = useNavigate();
    const [notification, setNotification] = useState({ message: '', isError: false });

    const [formData, setFormData] = useState({
        investment_type: 'stocks',
        amount: '',
        date: '',
        recurring: false,
        start_date: '',
        end_date: '',
        frequency: '',
        description: '',
    });

    const handleSubmit = async event => {
        event.preventDefault();

        let isRecurring = formData.recurring;

        try {
            if(isRecurring){
                await axios.post(`${backendUrl}/api/recurring/investments`, {
                    ...formData,
                    amount: Number(formData.amount),
                    start_date: `${formData.start_date}T00:00:00Z`,
                    end_date: formData.end_date ? `${formData.end_date}T00:00:00Z` : null,
                });
                setNotification({ message: 'Recurring Investment added successfully', isError: false });
            }
            else{
                await axios.post(`${backendUrl}/api/investments`, formData);
                setNotification({ message: 'Investment added successfully', isError: false });
            }
            clearFormData();

            navigate('/investments');
        } catch (error) {
            setNotification({ message: 'Error adding investment', isError: true });
        }
    };

    const clearFormData = () => {
        setFormData({
            investment_type: 'stocks',
            amount: '',
            date: '',
            recurring: false,
            start_date: '',
            end_date: '',
            frequency: '',
            description: '',
        });
        navigate('/investments');
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
        <main className="management-page form-page">
            <section className="expense-form-container">
            <p className="eyebrow">GROW YOUR FUTURE</p>
            <h1>Add an investment</h1>
            <p className="form-description">Record a contribution and keep your portfolio history up to date.</p>
            <form onSubmit={handleSubmit}>
                <div className="input-container">
                    <label htmlFor="investment_type">Investment Type:</label>
                    <select
                        id="investment_type"
                        name="investment_type"
                        value={formData.investment_type}
                        onChange={handleInputChange}
                    >
                        <option value="stocks">Stocks</option>
                        <option value="401k">401(k)</option>
                        <option value="mutual_funds">Mutual Funds</option>
                        <option value="real_estate">Real Estate</option>
                        <option value="savings">Savings</option>
                        <option value="bonds">Bonds</option>
                    </select>
                </div>
                <div className="input-container">
                    <label htmlFor="amount">Amount:</label>
                    <input
                        type="number"
                        id="amount"
                        name="amount"
                        value={formData.amount}
                        onChange={handleInputChange}
                        min="0.01"
                        step="0.01"
                        inputMode="decimal"
                        required
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
                        required
                    />
                </div>
                <div className="input-container">
                    <label htmlFor="recurring">Is Recurring Investment?</label>
                    <input
                        type="checkbox"
                        id="recurring"
                        name="recurring"
                        checked={formData.recurring}
                        onChange={handleInputChange}
                    />
                </div>
                {formData.recurring && (
                    <>
                        <div className="input-container">
                            <label htmlFor="recurring-description">Description:</label>
                            <input
                                type="text"
                                id="recurring-description"
                                name="description"
                                value={formData.description}
                                onChange={handleInputChange}
                                maxLength="255"
                                required
                            />
                        </div>
                        <div className="input-container">
                            <label htmlFor="startDate">Start Date:</label>
                            <input
                                type="date"
                                id="startDate"
                                name="start_date"
                                value={formData.start_date}
                                onChange={handleInputChange}
                                required
                            />
                        </div>
                        <div className="input-container">
                            <label htmlFor="endDate">End Date:</label>
                            <input
                                type="date"
                                id="endDate"
                                name="end_date"
                                value={formData.end_date}
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
                                required
                            >
                                <option value="">Choose a frequency</option>
                                <option value="daily">Daily</option>
                                <option value="weekly">Weekly</option>
                                <option value="biweekly">Bi-Weekly</option>
                                <option value="monthly">Monthly</option>
                            </select>
                        </div>
                    </>
                )}
                <div className="buttons-container">
                    <button className="primary-action" type="submit">Add investment</button>
                    <button className="secondary-action" type="button" onClick={clearFormData}>Cancel</button>
                </div>
            </form>
            <Notification message={notification.message} isError={notification.isError}/>
            </section>
        </main>
    )
        ;
};

export default InvestmentForm;
