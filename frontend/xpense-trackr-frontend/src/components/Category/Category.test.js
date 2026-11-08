import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import CategoryList from './CategoryList';
import CategoryForm from './CategoryForm';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

describe('Category pages', () => {
    beforeEach(() => {
        jest.clearAllMocks();
        axios.get.mockResolvedValue({ data: [{ id: 1, name: 'Food' }] });
        axios.post.mockResolvedValue({});
    });

    test('shows categories and edit actions', async () => {
        renderWithRouter(<CategoryList />);

        expect(await screen.findByRole('heading', { name: 'Food' })).toBeInTheDocument();
        expect(screen.getByRole('link', { name: 'Edit Food category' })).toHaveAttribute('href', '/categories/1/edit');
    });

    test('creates a category', async () => {
        renderWithRouter(<CategoryForm />, {
            route: '/categories/new',
            path: '/categories/new',
            extraRoutes: [{ path: '/categories', element: <div>Categories list</div> }],
        });

        fireEvent.change(screen.getByLabelText('Category name'), { target: { value: 'Home' } });
        fireEvent.click(screen.getByRole('button', { name: 'Add category' }));

        await screen.findByText('Categories list');
        await waitFor(() => expect(axios.post).toHaveBeenCalledWith(`${backendUrl}/api/categories`, { name: 'Home' }));
    });
});
