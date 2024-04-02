import React, { useEffect, useRef, useState } from 'react';
import Chart from 'chart.js/auto'; // Import Chart.js
import './Homepage.css';
import backendUrl from "../../config";
import 'chartjs-adapter-moment';

const InvestmentGraph = () => {
    const chartContainer = useRef(null);
    const chartInstance = useRef(null);
    const [investmentData, setInvestmentData] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                // Fetch investment data from API
                const response = await fetch(`${backendUrl}/api/investments`);
                const data = await response.json();
                setInvestmentData(data);
            } catch (error) {
                console.error('Error fetching investment data:', error);
            }
        };

        fetchData();
    }, []);

    useEffect(() => {
        if (chartContainer.current && investmentData.length > 0) {
            // Sort investment data by date
            const sortedData = investmentData.slice().sort((a, b) => new Date(a.date) - new Date(b.date));

            // Create the chart context
            const ctx = chartContainer.current.getContext('2d');

            // Fixed colors for the graphs
            const colors = ['#FFA280', '#80FFD6', '#80A4FF', '#FFF380', '#D980FF', '#80FF80'];

            // Create datasets for each investment type
            const datasets = {};
            sortedData.forEach((investment, index) => {
                if (!datasets[investment.investment_type]) {
                    datasets[investment.investment_type] = {
                        label: investment.investment_type,
                        data: [],
                        borderColor: colors[index % colors.length], // Use fixed colors for each dataset
                        borderWidth: 3, // Thicker lines
                    };
                }
                datasets[investment.investment_type].data.push({ x: investment.date, y: investment.amount });
            });

            // Define chart configuration
            const config = {
                type: 'line',
                data: {
                    datasets: Object.values(datasets),
                },
                options: {
                    scales: {
                        x: {
                            type: 'time',
                            time: {
                                unit: 'month', // Adjust time unit as needed
                            },
                        },
                        y: {
                            title: {
                                display: true,
                                text: 'Amount',
                            },
                        },
                    },
                    plugins: {
                        legend: {
                            display: true,
                            labels: {
                                boxWidth: 20, // Adjust legend box width
                                padding: 20, // Add padding between legend items
                            },
                        },
                    },
                },
            };

            // Create the chart instance
            chartInstance.current = new Chart(ctx, config);
        }

        // Cleanup function
        return () => {
            if (chartInstance.current) {
                chartInstance.current.destroy(); // Destroy the chart instance on unmount
            }
        };
    }, [investmentData]);

    return <canvas ref={chartContainer} className="graph-canvas" />;
};

export default InvestmentGraph;
