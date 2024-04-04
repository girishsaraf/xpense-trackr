// src/components/Expense/ExpenseList.js

import React, {useEffect, useState} from 'react';
import axios from 'axios';
import backendUrl from "../../config";
import {Link} from "react-router-dom";
import './ExpenseList.css';
import Pagination from "../Pagination/Pagination";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import {faEdit, faPlus, faTag, faTrashAlt} from '@fortawesome/free-solid-svg-icons';


const ExpenseList = () => {
    const [expenses, setExpenses] = useState([]);
    const [categories, setCategories] = useState([]);
    const [showModal, setShowModal] = useState(false);
    const [currentExpense, setCurrentExpense] = useState(null);
    const [editedAmount, setEditedAmount] = useState(0);
    const [editedDate, setEditedDate] = useState('');
    const [editedDescription, setEditedDescription] = useState('');
    const [editedCategory, setEditedCategory] = useState('');

    const [currentPage, setCurrentPage] = useState(1);
    const [expensesPerPage] = useState(15); // Number of expenses per page

    // Calculate the indexes of the expenses to display on the current page
    const indexOfLastExpense = currentPage * expensesPerPage;
    const indexOfFirstExpense = indexOfLastExpense - expensesPerPage;
    const currentExpenses = expenses.slice(indexOfFirstExpense, indexOfLastExpense);

    // Logic for changing the current page
    const paginate = (pageNumber) => setCurrentPage(pageNumber);

    useEffect(() => {
        fetchExpenses();
    }, []);

    const fetchExpenses = async () => {
        try {
            const response = await axios.get(`${backendUrl}/api/expenses`);

            // Sort expenses by date parameter in descending order
            const sortedExpenses = response.data.sort((a, b) => new Date(b.date) - new Date(a.date));
            setExpenses(sortedExpenses);

            const responseCategories = await axios.get(`${backendUrl}/api/categories`);
            setCategories(responseCategories.data);
        } catch (error) {
            console.error('Error fetching expenses:', error);
        }
    };

    const getCategoryNameById = (categoryId) => {
        const category = categories.find(cat => cat.id === +categoryId);
        return category ? category.name : 'Unknown Category';
    };


    const handleSubmitEdit = async () => {
        try {
            // Fetch the expense data for the selected expenseId
            const expenseId = currentExpense.id;
            const response = await axios.get(`${backendUrl}/api/expenses/${expenseId}`);
            const expenseData = response.data;

            // Update the expense data with edited amount and date
            expenseData.date = editedDate.split("T")[0];
            expenseData.amount = editedAmount.toString();
            expenseData.description = editedDescription.toString();
            expenseData.category_id = editedCategory.toString();

            // Call the PUT endpoint to update the expense
            await axios.put(`${backendUrl}/api/expenses/${expenseId}`, expenseData);
            alert('Expense updated successfully');

            // Optionally, fetch the list of expenses to update the UI
            fetchExpenses();
        } catch (error) {
            console.error('Error updating expense:', error);
            alert('Error updating expense. Please try again.');
        }
        setShowModal(false);
    };


    const handleOpenModal = (expense) => {
        setCurrentExpense(expense);
        setEditedAmount(expense.amount);
        setEditedDate(expense.date);
        setEditedDescription(expense.description);
        setEditedCategory(expense.category_id);
        setShowModal(true);
    };

    const handleCloseModal = () => {
        setShowModal(false);
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
                <h2>Recent Expense Overview</h2>
                {/* Add button to navigate to ExpenseForm page */}
                <Link to="/expenses/new" className="button-common-add">
                    <FontAwesomeIcon icon={faPlus} style={{ marginRight: '5px' }} /> Add New
                </Link>
            </div>
            <table className="expense-list-table">
                <thead>
                <tr>
                    <th>Category</th>
                    <th>Description</th>
                    <th>Amount</th>
                    <th>Date</th>
                    <th colSpan="2">Action</th>
                </tr>
                </thead>
                <tbody>
                {showModal && (
                    <div className="modal">
                        <div className="modal-content">
                            <h2>Edit Expense</h2>
                            <div>
                                <label htmlFor="edited-description">Description:</label>
                                <input
                                    type="text"
                                    id="edited-description"
                                    value={editedDescription}
                                    onChange={(e) => setEditedDescription(e.target.value)}
                                />
                            </div>
                            <div>
                                <label htmlFor="edited-amount">Amount:</label>
                                <input
                                    type="number"
                                    id="edited-amount"
                                    value={editedAmount}
                                    onChange={(e) => setEditedAmount(e.target.value)}
                                />
                            </div>
                            <div>
                                <label htmlFor="edited-date">Date:</label>
                                <input
                                    type="date"
                                    id="edited-date"
                                    value={editedDate}
                                    onChange={(e) => setEditedDate(e.target.value)}
                                />
                            </div>
                            <div>
                                <label htmlFor="edited-category">Category:</label>
                                <select
                                    id="edited-category"
                                    value={editedCategory}
                                    onChange={(e) => setEditedCategory(e.target.value)}
                                >
                                    {categories.map(category => (
                                        <option key={category.id} value={category.id}>
                                            {category.name}
                                        </option>
                                    ))}
                                </select>
                            </div>
                            <button onClick={handleSubmitEdit}>Submit</button>
                            <button onClick={handleCloseModal}>Cancel</button>
                        </div>
                    </div>
                )}
                {currentExpenses.map(expense => (
                    <tr key={expense.id}>
                        <td>
                            <div
                                className={`category-tag category-${getCategoryNameById(expense.category_id).toLowerCase()}`}>
                                <FontAwesomeIcon icon={faTag}/> {getCategoryNameById(expense.category_id)}
                            </div>
                        </td>
                        <td>{expense.description}</td>
                        <td>${expense.amount}</td>
                        <td>{formatDate(expense.date)}</td>
                        <td>
                            <button className="button-common" onClick={() => handleOpenModal(expense)}><FontAwesomeIcon
                                icon={faEdit}/></button>
                        </td>
                        <td>
                            <button className="button-common" onClick={() => handleDeleteExpense(expense.id)}>
                                <FontAwesomeIcon icon={faTrashAlt}/></button>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
            <Pagination
                currentPage={currentPage}
                itemsPerPage={expensesPerPage}
                totalItems={expenses.length}
                paginate={paginate}
            />
        </div>

    );
};

export default ExpenseList;
