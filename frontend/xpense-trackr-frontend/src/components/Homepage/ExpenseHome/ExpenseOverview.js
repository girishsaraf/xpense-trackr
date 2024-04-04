import React, {useEffect, useState} from 'react';
import backendUrl from "../../config";
import axios from "axios";
import './Homepage.css';

const ExpenseOverview = () => {
    const [expenseData, setExpenseData] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [expensesResponse, categoriesResponse] = await Promise.all([
                    axios.get(`${backendUrl}/api/expenses`),
                    axios.get(`${backendUrl}/api/categories`)
                ]);
                const expenses = expensesResponse.data;
                const categories = categoriesResponse.data.reduce((acc, category) => {
                    acc[category.id] = category.name;
                    return acc;
                }, {});
                const groupedExpenses = {};
                expenses.forEach(expense => {
                    const categoryName = categories[expense.category_id];
                    if (!groupedExpenses[categoryName]) {
                        groupedExpenses[categoryName] = 0;
                    }
                    groupedExpenses[categoryName] += parseFloat(expense.amount);
                });
                const expenseList = Object.entries(groupedExpenses).map(([categoryName, totalAmount]) => ({
                    category: categoryName,
                    totalAmount
                }));
                setExpenseData(expenseList);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        fetchData();
    }, []);

    return (
        <div className="expense-overview">
            <h2>Expense Overview</h2>
            <div className="overview-container">
                {expenseData.map((expense, index) => (
                    <div key={index} className="overview-box">
                        <span className="category">{expense.category}: </span>
                        <span className="amount">${expense.totalAmount}</span>
                    </div>
                ))}
            </div>
        </div>
    );
}

export default ExpenseOverview;
