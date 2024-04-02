import React, {useEffect, useRef, useState} from 'react';
import Chart from 'chart.js/auto'; // Import Chart.js
import './Homepage.css';
import backendUrl from "../../config";
import 'chartjs-adapter-moment';

const ExpenseGraph = () => {
    const chartContainer = useRef(null);
    const chartInstance = useRef(null);
    const [expenseData, setExpenseData] = useState([]);
    const [categoryData, setCategoryData] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            try {
                // Fetch expense data from API
                const expenseResponse = await fetch(`${backendUrl}/api/expenses`);
                const expenseData = await expenseResponse.json();
                setExpenseData(expenseData);

                // Fetch category data from API
                const categoryResponse = await fetch(`${backendUrl}/api/categories`);
                const categoryData = await categoryResponse.json();
                const categories = categoryData.reduce((acc, category) => {
                    acc[category.id] = category.name;
                    return acc;
                }, {});
                setCategoryData(categories);

                // Update expense data with category names
                const updatedExpenseData = expenseData.map(expense => ({
                    ...expense,
                    category: categories[expense.category_id]
                }));
                setExpenseData(updatedExpenseData);
            } catch (error) {
                console.error('Error fetching data:', error);
            }
        };

        fetchData();
    }, []);

    useEffect(() => {

        if (chartContainer.current && expenseData.length > 0) {
            // Sort investment data by date
            const sortedData = expenseData.slice().sort((a, b) => new Date(a.date) - new Date(b.date));

            // Create the chart context
            const ctx = chartContainer.current.getContext('2d');

            // Fixed colors for the graphs
            const colors = ['#FFA280', '#80FFD6', '#80A4FF', '#FFF380', '#D980FF'];


            // Create datasets for each investment type
            const datasets = {};
            sortedData.forEach((expense, index) => {
                if (!datasets[expense.category]) {
                    datasets[expense.category] = {
                        label: expense.category,
                        data: [],
                        borderColor: colors[index % colors.length], // Use fixed colors for each dataset
                        borderWidth: 3, // Thicker lines
                    };
                }
                datasets[expense.category].data.push({ x: expense.date, y: expense.amount });
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
    }, [expenseData]);

    return <canvas ref={chartContainer} className="graph-canvas" />;
};

export default ExpenseGraph;
