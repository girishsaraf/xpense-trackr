// src/components/Expense/ExpenseList.js

import React, {useEffect, useState} from 'react';
import axios from 'axios';
import backendUrl from "../../config";
import {Link} from "react-router-dom";
import './InvestmentList.css';
import Pagination from "../Pagination/Pagination";
import {FaPlus} from 'react-icons/fa';
import {faEdit, faPlus, faTag, faTrashAlt} from "@fortawesome/free-solid-svg-icons";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";

const InvestmentList = () => {
    const [investments, setInvestments] = useState([]);
    const [showInvModal, setShowInvModal] = useState(false);
    const [currentInvestment, setCurrentInvestment] = useState(null);
    const [editedAmount, setEditedAmount] = useState(0);
    const [editedDate, setEditedDate] = useState('');
    const [editedType, setEditedType] = useState('');

    const [currentPage, setCurrentPage] = useState(1);
    const [investmentsPerPage] = useState(15); // Number of expenses per page

    // Calculate the indexes of the expenses to display on the current page
    const indexOfLastInvestment = currentPage * investmentsPerPage;
    const indexOfFirstInvestment = indexOfLastInvestment - investmentsPerPage;
    const currentInvestments = investments.slice(indexOfFirstInvestment, indexOfLastInvestment);

    // Logic for changing the current page
    const paginate = (pageNumber) => setCurrentPage(pageNumber);

    useEffect(() => {
        fetchInvestments();
    }, []);

    const fetchInvestments = async () => {
        try {
            const response = await axios.get(`${backendUrl}/api/investments`);

            // Sort expenses by date parameter in descending order
            const sortedInvestments = response.data.sort((a, b) => new Date(b.date) - new Date(a.date));
            setInvestments(sortedInvestments);
        } catch (error) {
            console.error('Error fetching investments:', error);
        }
    };


    const handleSubmitEdit = async () => {
        try {
            // Fetch the expense data for the selected expenseId
            const investmentId = currentInvestment.id;
            const response = await axios.get(`${backendUrl}/api/investments/${investmentId}`);
            const investmentData = response.data;

            // Update the expense data with edited amount and date
            investmentData.date = editedDate.split("T")[0];
            investmentData.amount = editedAmount.toString();
            investmentData.investment_type = editedType.toString();

            // Call the PUT endpoint to update the expense
            await axios.put(`${backendUrl}/api/investments/${investmentId}`, investmentData);
            alert('Investment updated successfully');

            // Optionally, fetch the list of expenses to update the UI
            await fetchInvestments();
        } catch (error) {
            console.error('Error updating investment:', error);
            alert('Error updating investment. Please try again.');
        }
        setShowInvModal(false);
    };


    const handleOpenModal = (expense) => {
        setCurrentInvestment(expense);
        setEditedAmount(expense.amount);
        setEditedDate(expense.date);
        setEditedType(expense.investment_type);
        setShowInvModal(true);
    };

    const handleCloseModal = () => {
        setShowInvModal(false);
    };

    const handleDeleteInvestment = async (investmentId) => {
        try {
            await axios.delete(`${backendUrl}/api/investments/${investmentId}`);
            alert('Investment deleted successfully');
            await fetchInvestments(); // Fetch updated list of expenses after deletion
        } catch (error) {
            console.error('Error deleting investment:', error);
            alert('Error deleting investment. Please try again.');
        }
    };

    const getTypeName = (type) => {
        switch (type.toLowerCase()) {
            case 'stocks':
                return 'Stocks';
            case '401k':
                return '401(k)';
            case 'mutual_funds':
                return 'Mutual Funds';
            case 'real_estate':
                return 'Real Estate';
            case 'savings':
                return 'Savings';
            case 'bonds':
                return 'Bonds';
            default:
                return type;
        }
    };

    const formatDate = (dateString) => {
        const options = { year: 'numeric', month: 'short', day: 'numeric' };
        return new Date(dateString).toLocaleDateString('en-US', options);
    };

    return (
        <div className="investment-list-container">
            <div className="investment-list-header">
                <h2>Recent Investment Overview</h2>
                {/* Add button to navigate to InvestmentForm page */}
                <Link to="/investments/new" className="button-common-add">
                    <FontAwesomeIcon icon={faPlus} style={{ marginRight: '5px' }} /> Add New
                </Link>
            </div>
            <table className="investment-list-table">
                <thead>
                <tr>
                    <th>Investment Type</th>
                    <th>Amount</th>
                    <th>Date</th>
                    <th colSpan="2">Action</th>
                </tr>
                </thead>
                <tbody>
                {showInvModal && (
                    <div className="modal">
                        <div className="modal-content">
                            <h2>Edit Investment</h2>
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
                                <label htmlFor="edited-amount">Amount:</label>
                                <input
                                    type="number"
                                    id="edited-amount"
                                    value={editedAmount}
                                    onChange={(e) => setEditedAmount(e.target.value)}
                                />
                            </div>
                            <div>
                                <label htmlFor="edited-type">Type:</label>
                                <select
                                    id="edited-type"
                                    value={editedType}
                                    onChange={(e) => setEditedType(e.target.value)}
                                >
                                    <option value="stocks">Stocks</option>
                                    <option value="401k">401(k)</option>
                                    <option value="mutual_funds">Mutual Funds</option>
                                    <option value="real_estate">Real Estate</option>
                                    <option value="savings">Savings</option>
                                    <option value="bonds">Bonds</option>
                                </select>
                            </div>
                            <button onClick={handleSubmitEdit}>Submit</button>
                            <button onClick={handleCloseModal}>Cancel</button>
                        </div>
                    </div>
                )}
                {currentInvestments.map(investment => (
                    <tr key={investment.id}>
                        <td>
                            <div className={`type-tag type-${investment.investment_type.toLowerCase()}`}>
                                <FontAwesomeIcon icon={faTag}/> {getTypeName(investment.investment_type)}
                            </div>
                        </td>
                        <td>${investment.amount}</td>
                        <td>{formatDate(investment.date)}</td>
                        <td>
                            <button className="button-common" onClick={() => handleOpenModal(investment)}>
                                <FontAwesomeIcon icon={faEdit}/></button>
                        </td>
                        <td>
                            <button className="button-common" onClick={() => handleDeleteInvestment(investment.id)}>
                                <FontAwesomeIcon icon={faTrashAlt}/></button>
                        </td>
                    </tr>
                ))}
                </tbody>
            </table>
            <Pagination
                currentPage={currentPage}
                itemsPerPage={investmentsPerPage}
                totalItems={investments.length}
                paginate={paginate}
            />
        </div>

    );
};

export default InvestmentList;
