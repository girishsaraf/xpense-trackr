import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import ExpenseForm from './ExpenseForm';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

function renderExpenseForm() {
    return renderWithRouter(<ExpenseForm />, {
        route: '/expenses/new',
        path: '/expenses/new',
        extraRoutes: [{ path: '/expenses', element: <div>Expenses Page</div> }],
    });
}

async function fillExpenseFields() {
    await screen.findByRole('option', { name: 'Food' });
    fireEvent.change(screen.getByLabelText('Description:'), { target: { value: 'Lunch' } });
    fireEvent.change(screen.getByLabelText('Amount:'), { target: { value: '15.5' } });
    fireEvent.change(screen.getByLabelText('Category:'), { target: { value: '1' } });
    fireEvent.change(screen.getByLabelText('Date:'), { target: { value: '2026-01-20' } });
}

beforeEach(() => {
    jest.clearAllMocks();
    window.alert = jest.fn();
    axios.post.mockResolvedValue({});
    axios.get.mockResolvedValue({ data: [{ id: 1, name: 'Food' }] });
});

describe('ExpenseForm', () => {
    test('submits a one-off expense', async () => {
        renderExpenseForm();
        await fillExpenseFields();

        fireEvent.click(screen.getByRole('button', { name: 'Add' }));

        await screen.findByText('Expenses Page');
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/expenses`,
            expect.objectContaining({
                description: 'Lunch',
                amount: '15.5',
                category_id: '1',
                date: '2026-01-20',
                recurring: false,
            })
        );
        expect(window.alert).toHaveBeenCalledWith('Expense added successfully');
    });

    test('submits a recurring expense', async () => {
        renderExpenseForm();
        await fillExpenseFields();

        fireEvent.click(screen.getByLabelText('Is Recurring Expense?'));
        fireEvent.change(screen.getByLabelText('Start Date:'), { target: { value: '2026-01-20' } });
        fireEvent.change(screen.getByLabelText('End Date:'), { target: { value: '2026-06-20' } });
        fireEvent.change(screen.getByLabelText('Frequency:'), { target: { value: 'monthly' } });

        fireEvent.click(screen.getByRole('button', { name: 'Add' }));

        await screen.findByText('Expenses Page');
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/recurring/expenses`,
            expect.objectContaining({
                recurring: true,
                frequency: 'monthly',
                category_id: 1,
                amount: 15.5,
                start_date: '2026-01-20T00:00:00Z',
                end_date: '2026-06-20T00:00:00Z',
            })
        );
    });

    test('navigates back when cancelled', async () => {
        renderExpenseForm();

        fireEvent.click(screen.getByRole('button', { name: 'Cancel' }));

        await screen.findByText('Expenses Page');
        expect(axios.post).not.toHaveBeenCalled();
    });

    test('shows an alert when the request fails', async () => {
        const consoleError = jest.spyOn(console, 'error').mockImplementation(() => {});
        axios.post.mockRejectedValue(new Error('network'));
        renderExpenseForm();
        await fillExpenseFields();

        fireEvent.click(screen.getByRole('button', { name: 'Add' }));

        await waitFor(() => {
            expect(window.alert).toHaveBeenCalledWith('Error adding expense. Please try again.');
        });
        expect(consoleError).toHaveBeenCalled();
        consoleError.mockRestore();
    });
});
