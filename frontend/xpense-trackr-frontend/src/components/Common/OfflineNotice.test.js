import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import OfflineNotice from './OfflineNotice';

test('explains offline limits and updates when connectivity returns', () => {
    Object.defineProperty(navigator, 'onLine', { configurable: true, value: false });
    render(<OfflineNotice />);

    expect(screen.getByRole('status')).toHaveTextContent(
        'The app shell is available, but financial data isn’t stored for offline access.',
    );

    Object.defineProperty(navigator, 'onLine', { configurable: true, value: true });
    fireEvent(window, new Event('online'));
    expect(screen.queryByRole('status')).not.toBeInTheDocument();
});
