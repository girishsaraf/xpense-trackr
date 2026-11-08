import React, { useEffect, useRef } from 'react';
import Chart from 'chart.js/auto';
import './ChartVisualization.css';

const currency = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    maximumFractionDigits: 0,
});

export default function ChartVisualization({ type, data, ariaLabel }) {
    const canvasRef = useRef(null);
    const chartRef = useRef(null);

    useEffect(() => {
        const context = canvasRef.current?.getContext('2d');
        if (!context || !data?.labels?.length || !data?.datasets?.length) return undefined;

        const chart = new Chart(context, {
            type,
            data,
            options: {
                responsive: true,
                maintainAspectRatio: false,
                interaction: { intersect: false, mode: 'index' },
                plugins: {
                    legend: {
                        display: type === 'doughnut' || data.datasets.length > 1,
                        position: 'bottom',
                        labels: { usePointStyle: true, boxWidth: 8, padding: 18 },
                    },
                    tooltip: {
                        callbacks: {
                            label: (context) => {
                                const value = typeof context.parsed === 'number'
                                    ? context.parsed
                                    : context.parsed.y;
                                const label = type === 'doughnut' ? context.label : context.dataset.label;
                                return `${label}: ${currency.format(value || 0)}`;
                            },
                        },
                    },
                },
                scales: type === 'doughnut' ? {} : {
                    x: {
                        stacked: type === 'bar',
                        grid: { display: false },
                        ticks: { color: '#72819a', maxRotation: 0, autoSkip: true },
                    },
                    y: {
                        stacked: type === 'bar',
                        beginAtZero: true,
                        border: { display: false },
                        grid: { color: 'rgba(125, 146, 177, 0.16)' },
                        ticks: {
                            color: '#72819a',
                            callback: (value) => currency.format(value),
                        },
                    },
                },
                elements: {
                    line: { tension: 0.36, borderWidth: 2.5, fill: true },
                    point: { radius: 2, hoverRadius: 5, hitRadius: 10 },
                },
                ...(type === 'bar' && {
                    datasets: { bar: { borderRadius: 5, maxBarThickness: 36 } },
                }),
                ...(type === 'doughnut' && { cutout: '68%' }),
            },
        });
        chartRef.current = chart;

        return () => {
            if (typeof chart.destroy === 'function') chart.destroy();
            if (chartRef.current === chart) chartRef.current = null;
        };
    }, [data, type]);

    return (
        <div className={`chart-visualization chart-visualization-${type}`}>
            <canvas ref={canvasRef} role="img" aria-label={ariaLabel} />
        </div>
    );
}
