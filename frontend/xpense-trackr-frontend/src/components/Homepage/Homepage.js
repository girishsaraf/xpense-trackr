import React from 'react';
import ExpenseOverview from './ExpenseOverview';
import InvestmentOverview from './InvestmentOverview';
import Graph from './Graph';
import './Homepage.css';

const Homepage = () => {
    return (
        <div className="homepage">
            <h1>Welcome to Your Financial Overview</h1>
            <div className="overview-section">
                <ExpenseOverview />
                <InvestmentOverview />
            </div>
            <div className="graph-section">
                <Graph />
            </div>
        </div>
    );
}

export default Homepage;
