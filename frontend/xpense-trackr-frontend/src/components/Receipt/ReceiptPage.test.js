import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import ReceiptPage from './ReceiptPage';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

beforeEach(() => {
    jest.clearAllMocks();
    axios.get.mockImplementation((url) => {
        if (url.endsWith('/api/accounts')) {
            return Promise.resolve({ data: [{ id: 1, name: 'Checking', currency: 'USD' }] });
        }
        if (url.includes('/api/transactions')) {
            return Promise.resolve({
                data: [{
                    id: 2, account_id: 1, type: 'expense', amount_minor: 1250,
                    currency: 'USD', date: '2026-10-12', payee: 'Market',
                }],
            });
        }
        if (url.endsWith('/api/receipts')) return Promise.resolve({ data: [] });
        return Promise.reject(new Error(`Unexpected GET ${url}`));
    });
    axios.post.mockResolvedValue({});
});

test('uploads a supported receipt for a selected transaction', async () => {
    renderWithRouter(<ReceiptPage />, { route: '/receipts', path: '/receipts' });

    await screen.findByRole('heading', { name: 'Attach a receipt' });
    fireEvent.change(screen.getByLabelText('Transaction'), { target: { value: '2' } });
    const file = new File(['%PDF-1.4 receipt'], 'receipt.pdf', { type: 'application/pdf' });
    fireEvent.change(screen.getByLabelText(/PDF, PNG, or JPEG/), { target: { files: [file] } });
    fireEvent.click(screen.getByRole('button', { name: 'Attach receipt' }));

    await waitFor(() => {
        expect(axios.post).toHaveBeenCalledWith(
            `${backendUrl}/api/receipts`,
            expect.any(FormData),
        );
    });
    const payload = axios.post.mock.calls[0][1];
    expect(payload.get('transaction_id')).toBe('2');
    expect(payload.get('file')).toBe(file);
    expect(await screen.findByText('Receipt attached to the transaction.')).toBeInTheDocument();
});
