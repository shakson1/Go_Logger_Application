import React from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip,
} from 'chart.js';
import { Line } from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Tooltip
);

interface SparklineChartProps {
  data: number[];
  color?: string;
}

export const SparklineChart: React.FC<SparklineChartProps> = ({ data, color = '#3B82F6' }) => {
  const chartData = {
    labels: data.map((_, index) => index.toString()),
    datasets: [
      {
        data: data,
        borderColor: color,
        backgroundColor: color + '20',
        borderWidth: 1,
        fill: false,
        tension: 0.4,
        pointRadius: 0,
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
      tooltip: {
        enabled: false,
      },
    },
    scales: {
      y: {
        display: false,
      },
      x: {
        display: false,
      },
    },
    elements: {
      point: {
        radius: 0,
      },
    },
  };

  return (
    <div className="w-20 h-8">
      <Line data={chartData} options={options} />
    </div>
  );
}; 