import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import GoalList from './GoalList';
import GoalForm from './GoalForm';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

describe('Savings goal pages', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        axios.get.mockResolvedValue({
            data: [{ id: 1, name: 'Emergency fund', target_amount: '5000', current_amount: '1250', target_date: '2027-01-01' }],
        });
        axios.post.mockResolvedValue({});
    });

    test('shows goal progress and target date', async () => {
        renderWithRouter(<GoalList />);

        expect(await screen.findByRole('heading', { name: 'Emergency fund' })).toBeInTheDocument();
        expect(screen.getAllByText('$1,250.00')).toHaveLength(2);
        expect(screen.getByText('25% of the way there')).toBeInTheDocument();
        expect(screen.getByRole('progressbar', { name: 'Emergency fund progress' })).toHaveAttribute('aria-valuenow', '25');
    });

    test('creates a goal with current savings and a target date', async () => {
        renderWithRouter(<GoalForm />, {
            route: '/goals/new',
            path: '/goals/new',
            extraRoutes: [{ path: '/goals', element: <div>Goals list</div> }],
        });

        fireEvent.change(screen.getByLabelText('Goal name'), { target: { value: 'New car' } });
        fireEvent.change(screen.getByLabelText('Savings target'), { target: { value: '12000' } });
        fireEvent.change(screen.getByLabelText('Already saved'), { target: { value: '1500' } });
        fireEvent.change(screen.getByLabelText('Target date Optional'), { target: { value: '2028-08-01' } });
        fireEvent.click(screen.getByRole('button', { name: 'Create goal' }));

        await screen.findByText('Goals list');
        await waitFor(() => expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/goals`,
            { name: 'New car', target_amount: 12000, current_amount: 1500, target_date: '2028-08-01' }
        ));
    });
});
