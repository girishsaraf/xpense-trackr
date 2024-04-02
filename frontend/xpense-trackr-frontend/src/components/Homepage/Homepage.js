import React from 'react';
import ExpenseOverview from './ExpenseOverview';
import InvestmentOverview from './InvestmentOverview';
import InvestmentGraph from './InvestmentGraph';
import './Homepage.css';
import ExpenseGraph from "./ExpenseGraph";

const Homepage = () => {
    return (
        <div className="homepage">
            <h1>Welcome to Your Financial Overview</h1>
            <div className="overview-section">
                <ExpenseOverview />
                <InvestmentOverview />
            </div>
            <div className="graph-section">
                <div className="graph-container">
                    <h2 className="graph-title">Expense Graph</h2>
                    <ExpenseGraph />
                    <div className="legend">
                        {/* Legend for Expense Graph */}
                    </div>
                </div>
                <div className="graph-container">
                    <h2 className="graph-title">Investment Graph</h2>
                    <InvestmentGraph />
                    <div className="legend">
                        {/* Legend for Investment Graph */}
                    </div>
                </div>
            </div>
        </div>
    );
}

export default Homepage;
