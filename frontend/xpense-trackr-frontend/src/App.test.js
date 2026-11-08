import React from 'react';
import { fireEvent, render, screen, within } from '@testing-library/react';
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
    test('renders the header and navigation', async () => {
        render(<App />);

        expect(screen.getByRole('heading', { name: /Expense Tracker/ })).toBeInTheDocument();
        const navigation = screen.getByRole('navigation', { name: 'Main navigation' });
        expect(navigation.querySelectorAll(':scope > a')).toHaveLength(1);
        expect(navigation.querySelectorAll('details:not([open]) .navigation-menu a')).toHaveLength(16);
        expect(within(navigation).getByText('Activity')).toBeInTheDocument();
        expect(within(navigation).getByText('Planning')).toBeInTheDocument();
        expect(within(navigation).getByText('Wealth')).toBeInTheDocument();
        expect(within(navigation).getByText('More')).toBeInTheDocument();
        expect(await screen.findByText('Homepage')).toBeInTheDocument();
    });

    test('navigates to the expense list', async () => {
        render(<App />);

        fireEvent.click(screen.getByText('Activity'));
        expect(screen.getByRole('link', { name: 'Expenses' })).toHaveAttribute('href', '/expenses');
        fireEvent.click(screen.getByRole('link', { name: 'Expenses' }));

        expect(await screen.findByText('Recent Expense Overview')).toBeInTheDocument();
        expect(screen.getByRole('navigation', { name: 'Main navigation' }).querySelectorAll('details[open]')).toHaveLength(0);
    });

    test('keeps only the selected navigation menu open', () => {
        render(<App />);
        const navigation = screen.getByRole('navigation', { name: 'Main navigation' });

        fireEvent.click(screen.getByText('Activity'));
        expect(navigation.querySelectorAll('details[open]')).toHaveLength(1);

        fireEvent.click(screen.getByText('Planning'));
        expect(navigation.querySelectorAll('details[open]')).toHaveLength(1);
        expect(navigation.querySelector('details[open] summary')).toHaveTextContent('Planning');
    });

    test('navigates to the investment list', async () => {
        render(<App />);

        fireEvent.click(screen.getByText('Wealth'));
        fireEvent.click(screen.getByRole('link', { name: 'Investments' }));

        expect(await screen.findByText('Recent Investment Overview')).toBeInTheDocument();
    });
});
