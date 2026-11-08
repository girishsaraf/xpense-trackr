import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import axios from 'axios';
import App from './App';

jest.mock('./components/Homepage/Homepage', () => () => <div>Homepage</div>);
jest.mock('axios');

beforeEach(() => {
    window.history.pushState({}, '', '/');
    jest.clearAllMocks();
    axios.get.mockResolvedValue({ data: [] });
});

describe('App', () => {
    test('renders the header and navigation', () => {
        render(<App />);

        expect(screen.getByRole('heading', { name: /Expense Tracker/ })).toBeInTheDocument();
        expect(screen.getByRole('link', { name: 'Expense' })).toHaveAttribute('href', '/expenses');
        expect(screen.getByRole('link', { name: 'Investments' })).toHaveAttribute('href', '/investments');
        expect(screen.getByText('Homepage')).toBeInTheDocument();
    });

    test('navigates to the expense list', async () => {
        render(<App />);

        fireEvent.click(screen.getByRole('link', { name: 'Expense' }));

        expect(await screen.findByText('Recent Expense Overview')).toBeInTheDocument();
    });

    test('navigates to the investment list', async () => {
        render(<App />);

        fireEvent.click(screen.getByRole('link', { name: 'Investments' }));

        expect(await screen.findByText('Recent Investment Overview')).toBeInTheDocument();
    });
});
