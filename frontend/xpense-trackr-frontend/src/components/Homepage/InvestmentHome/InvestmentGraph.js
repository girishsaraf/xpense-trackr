import React, { useState, useEffect, useRef } from 'react';
import Chart from 'chart.js/auto'; // Import Chart.js
import backendUrl from "../../config";
import 'chartjs-adapter-moment';

const InvestmentGraph = () => {
    const [graphData, setGraphData] = useState(null);
    const [lastMonths, setLastMonths] = useState(6); // Default is last 6 months
    const chartRef = useRef(null);
    const [chartInstance, setChartInstance] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch(`${backendUrl}/api/investments/cumulative?lastMonths=${lastMonths}`);
                const data = await response.json();
                const datasets = {};

                const colors = ['#FFA280', '#80FFD6', '#80A4FF', '#FFF380', '#D980FF', '#33FF33'];

                // Group data by investment type
                data.forEach((item, index) => {
                    if (!datasets[item.type]) {
                        datasets[item.type] = { label: item.type, data: [], borderColor: colors[index % colors.length] };
                    }
                    datasets[item.type].data.push({ x: item.month, y: item.amount });
                });

                // Prepare data for Chart.js
                const chartData = {
                    labels: [...new Set(data.map(item => item.month))],
                    datasets: Object.values(datasets),
                };

                setGraphData(chartData);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        fetchData();
    }, [lastMonths]);

    useEffect(() => {
        // Destroy the existing chart instance before rendering a new one
        if (chartInstance) {
            chartInstance.destroy();
        }
        if (graphData && chartRef.current) {
            const ctx = chartRef.current.getContext('2d');
            const newChartInstance = new Chart(ctx, {
                type: 'line',
                data: graphData,
                options: {
                    scales: {
                        x: {
                            type: 'time',
                            time: {
                                unit: 'month',
                            },
                        },
                        y: {
                            title: {
                                display: true,
                                text: 'Amount',
                            },
                        },
                    },
                },
            });
            setChartInstance(newChartInstance);
        }
    }, [graphData]);

    const handleButtonClick = (months) => {
        setLastMonths(months+1);
    };

    return (
        <div className="investment-graph">
            <div className="buttons-container">
                <button onClick={() => handleButtonClick(2)}>2 Months</button>
                <button onClick={() => handleButtonClick(4)}>4 Months</button>
                <button onClick={() => handleButtonClick(6)}>6 Months</button>
            </div>
            <canvas ref={chartRef}></canvas>
        </div>
    );
};

export default InvestmentGraph;
