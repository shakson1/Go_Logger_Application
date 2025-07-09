import React, { useState } from 'react';
import { TopSource } from '../types';
import { SparklineChart } from './SparklineChart';

interface TopSourcesTableProps {
  sources: TopSource[];
}

export const TopSourcesTable: React.FC<TopSourcesTableProps> = ({ sources }) => {
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 5;
  const totalPages = Math.ceil(sources.length / itemsPerPage);
  
  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentSources = sources.slice(startIndex, endIndex);

  const getCategoryColor = (category: string) => {
    switch (category) {
      case 'access':
        return 'text-blue-400';
      case 'network':
        return 'text-green-400';
      case 'threat':
        return 'text-red-400';
      case 'uba':
        return 'text-purple-400';
      default:
        return 'text-gray-400';
    }
  };

  const getCategoryBadge = (category: string) => {
    const colorClass = getCategoryColor(category);
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium bg-opacity-20 ${colorClass} bg-current`}>
        {category.toUpperCase()}
      </span>
    );
  };

  return (
    <div className="bg-splunk-gray rounded-lg border border-splunk-light-gray">
      <div className="px-6 py-4 border-b border-splunk-light-gray">
        <h3 className="text-lg font-semibold text-white">Top Event Sources</h3>
      </div>
      
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-splunk-darker">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Source IP
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Trend
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Count
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Category
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-splunk-light-gray">
            {currentSources.map((source, index) => (
              <tr key={index} className="hover:bg-splunk-darker">
                <td className="px-6 py-4 whitespace-nowrap text-sm text-white font-mono">
                  {source.sourceIP}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <SparklineChart data={source.sparkline} />
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-white">
                  {source.count.toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {getCategoryBadge(source.category)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      
      {totalPages > 1 && (
        <div className="px-6 py-4 border-t border-splunk-light-gray">
          <div className="flex items-center justify-between">
            <div className="text-sm text-gray-400">
              Showing {startIndex + 1} to {Math.min(endIndex, sources.length)} of {sources.length} results
            </div>
            <div className="flex space-x-2">
              <button
                onClick={() => setCurrentPage(Math.max(1, currentPage - 1))}
                disabled={currentPage === 1}
                className="px-3 py-1 text-sm bg-splunk-darker text-white rounded border border-splunk-light-gray disabled:opacity-50 disabled:cursor-not-allowed hover:bg-splunk-light-gray"
              >
                Previous
              </button>
              <span className="px-3 py-1 text-sm text-white">
                {currentPage} of {totalPages}
              </span>
              <button
                onClick={() => setCurrentPage(Math.min(totalPages, currentPage + 1))}
                disabled={currentPage === totalPages}
                className="px-3 py-1 text-sm bg-splunk-darker text-white rounded border border-splunk-light-gray disabled:opacity-50 disabled:cursor-not-allowed hover:bg-splunk-light-gray"
              >
                Next
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}; 