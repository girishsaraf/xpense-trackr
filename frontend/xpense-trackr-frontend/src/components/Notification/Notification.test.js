import React from 'react';
import { act, render, screen } from '@testing-library/react';
import Notification from './Notification';

describe('Notification', () => {
    afterEach(() => {
        jest.useRealTimers();
    });

    test('shows a success message', () => {
        render(<Notification message="Saved" isError={false} />);

        const notification = screen.getByText('Saved');
        expect(notification).toHaveClass('notification', 'success', 'show');
    });

    test('shows an error message', () => {
        render(<Notification message="Failed" isError />);

        expect(screen.getByText('Failed')).toHaveClass('notification', 'error', 'show');
    });

    test('hides itself after three seconds', () => {
        jest.useFakeTimers();
        render(<Notification message="Saved" isError={false} />);

        const notification = screen.getByText('Saved');
        expect(notification).toHaveClass('show');

        act(() => {
            jest.advanceTimersByTime(3000);
        });

        expect(notification).not.toHaveClass('show');
    });
});
