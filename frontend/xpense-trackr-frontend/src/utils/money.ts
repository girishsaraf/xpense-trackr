export function parseAmountMinor(value: string, currency = 'USD'): number | null {
    const digits = currencyFractionDigits(currency);
    const fraction = digits > 0 ? `(?:\\.(\\d{1,${digits}}))?` : '';
    const match = new RegExp(`^(-?)(0|[1-9]\\d*)${fraction}$`).exec(value.trim());
    if (!match) return null;

    const sign = match[1] === '-' ? -1 : 1;
    const scale = 10 ** digits;
    const fractionValue = digits > 0 ? Number((match[3] || '').padEnd(digits, '0')) : 0;
    const minor = sign * (Number(match[2]) * scale + fractionValue);
    return Number.isSafeInteger(minor) ? minor : null;
}

export function formatAmountMinor(amountMinor: number, currency: string): string {
    const digits = currencyFractionDigits(currency);
    return new Intl.NumberFormat(undefined, {
        style: 'currency',
        currency,
    }).format(amountMinor / (10 ** digits));
}

function currencyFractionDigits(currency: string): number {
    return new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency,
    }).resolvedOptions().maximumFractionDigits;
}
