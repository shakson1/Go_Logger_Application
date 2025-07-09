import React from 'react';
import { TrendingUp, TrendingDown } from 'lucide-react';

interface StatTileProps {
  title: string;
  total: number;
  delta: number;
  color: string;
}

export const StatTile: React.FC<StatTileProps> = ({ title, total, delta, color }) => {
  const isPositive = delta >= 0;
  
  return (
    <div className="bg-splunk-gray rounded-lg p-6 border border-splunk-light-gray">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-gray-400 text-sm font-medium">{title}</p>
          <p className="text-2xl font-bold text-white mt-2">{total.toLocaleString()}</p>
        </div>
        <div className={`flex items-center space-x-1 ${isPositive ? 'text-green-400' : 'text-red-400'}`}>
          {isPositive ? (
            <TrendingUp className="w-4 h-4" />
          ) : (
            <TrendingDown className="w-4 h-4" />
          )}
          <span className="text-sm font-medium">
            {isPositive ? '+' : ''}{delta}
          </span>
        </div>
      </div>
    </div>
  );
}; 