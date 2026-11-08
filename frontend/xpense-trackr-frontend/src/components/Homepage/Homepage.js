import ExpenseOverview from "./ExpenseHome/ExpenseOverview";
import InvestmentOverview from "./InvestmentHome/InvestmentOverview";
import ExpenseGraph from "./ExpenseHome/ExpenseGraph";
import InvestmentGraph from "./InvestmentHome/InvestmentGraph";
import {useState} from "react"; // Import your CSS file
import { Link } from "react-router-dom";
import './Homepage.css';

const Homepage = () => {
    const [activeView, setActiveView] = useState('expense'); // Default to expense view

    const toggleView = (view) => {
        setActiveView(view);
    };

    return (
        <div className="homepage">
            <div className="page-intro">
                <div>
                    <p className="eyebrow">YOUR FINANCES AT A GLANCE</p>
                    <h1>Your money, in focus.</h1>
                    <p className="page-description">A clear view of your spending and investments, all in one place.</p>
                </div>
                <div className="dashboard-tools">
                    <div className="view-switch" aria-label="Financial overview">
                        <button
                            className={`toggle-button ${activeView === 'expense' ? 'active' : ''}`}
                            aria-pressed={activeView === 'expense'}
                            onClick={() => toggleView('expense')}
                        >
                            Expenses
                        </button>
                        <button
                            className={`toggle-button ${activeView === 'investment' ? 'active' : ''}`}
                            aria-pressed={activeView === 'investment'}
                            onClick={() => toggleView('investment')}
                        >
                            Investments
                        </button>
                    </div>
                    <div className="quick-actions" aria-label="Quick actions">
                        <Link className="secondary-action" to="/expenses/new">Add expense</Link>
                        <Link className="primary-action" to="/investments/new">Add investment</Link>
                    </div>
                </div>
            </div>

            <div className="overview-section">
                {activeView === 'expense' && <ExpenseOverview/>}
                {activeView === 'investment' && <InvestmentOverview/>}
            </div>
            <div className="graph-section">
                {activeView === 'expense' && (
                    <div className="graph-container">
                        <div className="section-heading">
                            <div>
                                <p className="eyebrow">MONTHLY BREAKDOWN</p>
                                <h2 className="graph-title">Spending over time</h2>
                            </div>
                        </div>
                        <ExpenseGraph/>
                    </div>
                )}
                {activeView === 'investment' && (
                    <div className="graph-container">
                        <div className="section-heading">
                            <div>
                                <p className="eyebrow">PORTFOLIO HISTORY</p>
                                <h2 className="graph-title">Investments over time</h2>
                            </div>
                        </div>
                        <InvestmentGraph/>
                    </div>
                )}
            </div>
        </div>
    );
};

export default Homepage;
