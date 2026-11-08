// src/components/Expense/ExpenseForm.js

import React, { useEffect, useState } from 'react';
import axios from 'axios';
import backendUrl from "../../config";
import './ExpenseForm.css';
import { useNavigate } from 'react-router-dom';


const ExpenseForm = () => {

    const navigate = useNavigate();
    const [categories, setCategories] = useState([]);
    const [categoryError, setCategoryError] = useState('');
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

    useEffect(() => {
        const loadCategories = async () => {
            try {
                const response = await axios.get(`${backendUrl}/api/categories`);
                setCategories(response.data || []);
            } catch (error) {
                console.error('Error fetching categories:', error);
                setCategoryError('Categories could not be loaded. Please try again before adding an expense.');
            }
        };
        loadCategories();
    }, []);

    const handleSubmit = async event => {
        event.preventDefault();

        let isRecurring = formData.recurring;

        try
        {
            if(isRecurring){
                await axios.post(`${backendUrl}/api/recurring/expenses`, {
                    ...formData,
                    category_id: Number(formData.category_id),
                    amount: Number(formData.amount),
                    start_date: `${formData.start_date}T00:00:00Z`,
                    end_date: formData.end_date ? `${formData.end_date}T00:00:00Z` : null,
                });
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
        <main className="management-page form-page">
        <section className="expense-form-container">
            <p className="eyebrow">KEEP YOUR PICTURE CLEAR</p>
            <h1>Add an expense</h1>
            <p className="form-description">Capture a purchase and keep your monthly spending up to date.</p>
            <form onSubmit={handleSubmit}>
                {categoryError && <p className="inline-message" role="alert">{categoryError}</p>}
                <div className="input-container">
                    <label htmlFor="description">Description:</label>
                    <input
                        type="text"
                        id="description"
                        name="description"
                        value={formData.description}
                        onChange={handleInputChange}
                        maxLength="255"
                        required
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
                        min="0.01"
                        step="0.01"
                        inputMode="decimal"
                        required
                    />
                </div>
                <div className="input-container">
                    <label htmlFor="category_id">Category:</label>
                    <select
                        id="category_id"
                        name="category_id"
                        value={formData.category_id}
                        onChange={handleInputChange}
                        required
                    >
                        <option value="">Select a category</option>
                        {categories.map((category) => <option key={category.id} value={category.id}>{category.name}</option>)}
                    </select>
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
                    <label htmlFor="recurring">Is Recurring Expense?</label>
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
        </section>
        </main>
    );

};

export default ExpenseForm;
