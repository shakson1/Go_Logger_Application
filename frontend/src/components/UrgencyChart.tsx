import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Bar } from 'react-chartjs-2';
import { UrgencyData } from '../types';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  Title,
  Tooltip,
  Legend
);

interface UrgencyChartProps {
  data: UrgencyData;
}

export const UrgencyChart: React.FC<UrgencyChartProps> = ({ data }) => {
  const chartData = {
    labels: ['Critical', 'High', 'Medium', 'Low'],
    datasets: [
      {
        label: 'Notable Events',
        data: [data.critical, data.high, data.medium, data.low],
        backgroundColor: [
          'rgba(239, 68, 68, 0.8)', // Red for critical
          'rgba(245, 158, 11, 0.8)', // Orange for high
          'rgba(59, 130, 246, 0.8)', // Blue for medium
          'rgba(34, 197, 94, 0.8)',  // Green for low
        ],
        borderColor: [
          'rgba(239, 68, 68, 1)',
          'rgba(245, 158, 11, 1)',
          'rgba(59, 130, 246, 1)',
          'rgba(34, 197, 94, 1)',
        ],
        borderWidth: 1,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
      title: {
        display: true,
        text: 'Notable Events by Urgency',
        color: '#ffffff',
        font: {
          size: 16,
          weight: 'bold' as const,
        },
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        grid: {
          color: 'rgba(255, 255, 255, 0.1)',
        },
        ticks: {
          color: '#ffffff',
        },
      },
      x: {
        grid: {
          color: 'rgba(255, 255, 255, 0.1)',
        },
        ticks: {
          color: '#ffffff',
        },
      },
    },
  };

  return (
    <div className="bg-splunk-gray rounded-lg p-6 border border-splunk-light-gray">
      <div className="h-64">
        <Bar data={chartData} options={options} />
      </div>
    </div>
  );
}; 