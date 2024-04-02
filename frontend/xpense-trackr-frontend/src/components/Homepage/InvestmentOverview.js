import React, {useEffect, useState} from 'react';
import axios from "axios";
import backendUrl from "../../config";
import './Homepage.css';

const InvestmentOverview = () => {
    const [investmentData, setInvestmentData] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await axios.get(`${backendUrl}/api/investments`);
                const investments = response.data;
                const groupedInvestments = {};
                investments.forEach(investment => {
                    if(!groupedInvestments[investment.investment_type]) {
                        groupedInvestments[investment.investment_type] = 0
                    }
                    groupedInvestments[investment.investment_type] += parseFloat(investment.amount);
                });
                const investmentList = Object.entries(groupedInvestments).map(([type, totalAmount]) => ({
                    type: type,
                    totalAmount
                }));
                setInvestmentData(investmentList);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        fetchData();
    }, []);

    return (
        <div className="expense-overview">
            <h2>Investment Overview</h2>
            <ul className="overview-list">
                {investmentData.map((expense, index) => (
                    <li key={index} className="overview-item">
                        <span className="type">{expense.type}: </span>
                        <span className="amount">${expense.totalAmount}</span>
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default InvestmentOverview;
