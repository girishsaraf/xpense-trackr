import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import Pagination from './Pagination';

describe('Pagination', () => {
    test('renders one button per page', () => {
        render(<Pagination currentPage={1} itemsPerPage={15} totalItems={35} paginate={jest.fn()} />);

        expect(screen.getByRole('button', { name: '1' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: '2' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: '3' })).toBeInTheDocument();
    });

    test('calls paginate with the selected page', () => {
        const paginate = jest.fn();
        render(<Pagination currentPage={1} itemsPerPage={15} totalItems={35} paginate={paginate} />);

        fireEvent.click(screen.getByRole('button', { name: '2' }));

        expect(paginate).toHaveBeenCalledWith(2);
    });

    test('highlights the current page', () => {
        render(<Pagination currentPage={2} itemsPerPage={15} totalItems={35} paginate={jest.fn()} />);

        expect(screen.getByRole('button', { name: '2' }).parentElement).toHaveClass('active');
        expect(screen.getByRole('button', { name: '1' }).parentElement).not.toHaveClass('active');
    });

    test('renders no buttons when there are no items', () => {
        render(<Pagination currentPage={1} itemsPerPage={15} totalItems={0} paginate={jest.fn()} />);

        expect(screen.queryAllByRole('button')).toHaveLength(0);
    });
});
