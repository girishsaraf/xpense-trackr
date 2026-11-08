import { formatAmountMinor, parseAmountMinor } from './money';

describe('parseAmountMinor', () => {
    test.each([
        ['12', 1200],
        ['12.3', 1230],
        ['12.34', 1234],
        ['-0.01', -1],
    ])('parses %s without floating-point rounding', (value, expected) => {
        expect(parseAmountMinor(value)).toBe(expected);
    });

    test.each(['', '1.234', '1,000.00', '1e2', '90071992547409.92'])(
        'rejects invalid or unsafe amount %s',
        (value) => {
            expect(parseAmountMinor(value)).toBeNull();
        },
    );
});

test('formats minor units using the selected currency', () => {
    expect(formatAmountMinor(1234, 'USD')).toContain('12.34');
});

test('respects currencies with zero and three fractional digits', () => {
    expect(parseAmountMinor('1234', 'JPY')).toBe(1234);
    expect(parseAmountMinor('12.34', 'JPY')).toBeNull();
    expect(parseAmountMinor('1.234', 'KWD')).toBe(1234);
    expect(formatAmountMinor(1234, 'KWD')).toContain('1.234');
});
