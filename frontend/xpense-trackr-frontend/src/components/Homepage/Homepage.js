import ExpenseOverview from "./ExpenseHome/ExpenseOverview";
import InvestmentOverview from "./InvestmentHome/InvestmentOverview";
import ExpenseGraph from "./ExpenseHome/ExpenseGraph";
import InvestmentGraph from "./InvestmentHome/InvestmentGraph";
import {useState} from "react"; // Import your CSS file
import './Homepage.css';

const Homepage = () => {
    const [activeView, setActiveView] = useState('expense'); // Default to expense view

    const toggleView = (view) => {
        setActiveView(view);
    };

    return (
        <div className="homepage">
            <h1>Welcome to Your Financial Overview</h1>
            <div className="toggle-container" onClick={() => toggleView(activeView === 'expense' ? 'investment' : 'expense')}>
                <button
                    className={`toggle-button ${activeView === 'expense' ? 'active' : ''}`}
                    onClick={() => toggleView('expense')}
                >
                    Expenses
                </button>
                <div className="toggle-switch"
                     onClick={() => toggleView(activeView === 'expense' ? 'investment' : 'expense')}></div>
                <button
                    className={`toggle-button ${activeView === 'investment' ? 'active' : ''}`}
                    onClick={() => toggleView('investment')}
                >
                    Investments
                </button>
            </div>

            <div className="overview-section">
                {activeView === 'expense' && <ExpenseOverview/>}
                {activeView === 'investment' && <InvestmentOverview/>}
            </div>
            <div className="graph-section">
                {activeView === 'expense' && (
                    <div className="graph-container">
                        <h2 className="graph-title">Expense Graph</h2>
                        <ExpenseGraph/>
                        <div className="legend">
                            {/* Legend for Expense Graph */}
                        </div>
                    </div>
                )}
                {activeView === 'investment' && (
                    <div className="graph-container">
                        <h2 className="graph-title">Investment Graph</h2>
                        <InvestmentGraph/>
                        <div className="legend">
                            {/* Legend for Investment Graph */}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Homepage;
