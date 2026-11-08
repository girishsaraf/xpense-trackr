import React from 'react';
import { fireEvent, screen, waitFor } from '@testing-library/react';
import axios from 'axios';
import ImportExportPage from './ImportExportPage';
import backendUrl from '../../config';
import { renderWithRouter } from '../../testUtils';

jest.mock('axios');

const preview = {
    total: 1,
    ready: 1,
    duplicates: 0,
    invalid: 0,
    imported: 0,
    rows: [{
        line: 2,
        status: 'ready',
        transaction: { date: '2026-10-12', payee: 'Market', amount_minor: 1234, currency: 'USD' },
    }],
};

beforeEach(() => {
    jest.clearAllMocks();
    axios.post
        .mockResolvedValueOnce({ data: preview })
        .mockResolvedValueOnce({ data: { ...preview, ready: 0, imported: 1, rows: [{ ...preview.rows[0], status: 'imported' }] } });
});

test('previews a transaction CSV before committing the import', async () => {
    renderWithRouter(<ImportExportPage />, { route: '/import-export', path: '/import-export' });

    const file = new File([
        'account_id,type,amount,currency,date,payee\n1,expense,12.34,USD,2026-10-12,Market',
    ], 'transactions.csv', { type: 'text/csv' });
    fireEvent.change(screen.getByLabelText('Choose CSV file'), { target: { files: [file] } });
    fireEvent.click(screen.getByRole('button', { name: 'Preview import' }));

    expect(await screen.findByText(/1 rows · 1 ready · 0 duplicates · 0 invalid/)).toBeInTheDocument();
    fireEvent.click(screen.getByRole('button', { name: 'Import 1 ready' }));

    await waitFor(() => {
        expect(axios.post).toHaveBeenNthCalledWith(
            1,
            `${backendUrl}/api/transactions/import.csv`,
            expect.any(FormData),
        );
        expect(axios.post).toHaveBeenNthCalledWith(
            2,
            `${backendUrl}/api/transactions/import.csv?commit=true`,
            expect.any(FormData),
        );
    });
    expect(await screen.findByText('1 transaction imported.')).toBeInTheDocument();
    expect(screen.getByRole('link', { name: 'Download transaction CSV' })).toHaveAttribute('href', `${backendUrl}/api/transactions/export.csv`);
});
