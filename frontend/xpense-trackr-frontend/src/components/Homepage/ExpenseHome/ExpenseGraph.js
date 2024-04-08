import React, {useEffect, useRef, useState} from 'react';
import Chart from 'chart.js/auto'; // Import Chart.js
import backendUrl from "../../../config";
import 'chartjs-adapter-moment';
import './ExpenseHome.css';

const ExpenseGraph = () => {
    const [graphData, setGraphData] = useState(null);
    const [lastMonths, setLastMonths] = useState(6); // Default is last 6 months
    const chartRef = useRef(null);
    const [chartInstance, setChartInstance] = useState(null);

    useEffect(() => {
        const fetchData = async () => {
            try {
                const [expensesResponse] = await Promise.all([
                    fetch(`${backendUrl}/api/expenses/cumulative?lastMonths=${lastMonths}`)
                ]);

                const expensesData = await expensesResponse.json();

                const datasets = {};
                const colors = ['#FFA280', '#80FFD6', '#80A4FF', '#FFF380', '#D980FF', '#33FF33'];

                expensesData.forEach((item, index) => {
                    datasets[item.category] = datasets[item.category] || { label: item.category, data: [], borderColor: colors[index % colors.length] };
                    datasets[item.category].data.push({ x: item.month, y: item.amount });
                });

                const chartData = {
                    labels: [...new Set(expensesData.map(item => item.month))],
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

export default ExpenseGraph;
