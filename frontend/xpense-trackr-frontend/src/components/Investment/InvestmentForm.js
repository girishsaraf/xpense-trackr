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
    });

    const handleSubmit = async event => {
        event.preventDefault();

        let isRecurring = formData.recurring;

        console.log(formData);

        try {
            if(isRecurring){
                await axios.post(`${backendUrl}/api/recurring/investments`, formData);
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
        });
        navigate('/investments');
    }

    const handleInputChange = event => {
        const { name, value, type, checked } = event.target;
        console.log(name, value);
        const newValue = type === 'checkbox' ? checked : value;

        setFormData({
            ...formData,
            [name]: newValue,
        });
    };

    return (
        <div className="expense-form-container">
            <h2>Add New Investment</h2>
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
                        <option value="401_k">401(k)</option>
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
                    <label htmlFor="recurring">Is Recurring Investment?</label>
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
                    <button type="cancel" onClick={clearFormData}>Cancel</button>
                </div>
            </form>
            <div>
                <Notification message={notification.message} isError={notification.isError}/>
            </div>
        </div>
    )
        ;
};

export default InvestmentForm;
