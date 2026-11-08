import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import Homepage from './Homepage';

jest.mock('./ExpenseHome/ExpenseOverview', () => () => <div>ExpenseOverviewMock</div>);
jest.mock('./ExpenseHome/ExpenseGraph', () => () => <div>ExpenseGraphMock</div>);
jest.mock('./InvestmentHome/InvestmentOverview', () => () => <div>InvestmentOverviewMock</div>);
jest.mock('./InvestmentHome/InvestmentGraph', () => () => <div>InvestmentGraphMock</div>);

describe('Homepage', () => {
    const renderHomepage = () => render(<MemoryRouter><Homepage /></MemoryRouter>);

    test('shows the expense view by default', () => {
        renderHomepage();

        expect(screen.getByText('ExpenseOverviewMock')).toBeInTheDocument();
        expect(screen.getByText('ExpenseGraphMock')).toBeInTheDocument();
        expect(screen.queryByText('InvestmentOverviewMock')).not.toBeInTheDocument();
    });

    test('switches to the investment view', () => {
        renderHomepage();

        fireEvent.click(screen.getByRole('button', { name: 'Investments' }));

        expect(screen.getByText('InvestmentOverviewMock')).toBeInTheDocument();
        expect(screen.getByText('InvestmentGraphMock')).toBeInTheDocument();
        expect(screen.queryByText('ExpenseOverviewMock')).not.toBeInTheDocument();
    });

    test('switches back to the expense view', () => {
        renderHomepage();

        fireEvent.click(screen.getByRole('button', { name: 'Investments' }));
        fireEvent.click(screen.getByRole('button', { name: 'Expenses' }));

        expect(screen.getByText('ExpenseOverviewMock')).toBeInTheDocument();
        expect(screen.queryByText('InvestmentOverviewMock')).not.toBeInTheDocument();
    });
});
