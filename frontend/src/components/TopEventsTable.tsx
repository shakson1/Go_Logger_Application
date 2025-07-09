import React, { useState } from 'react';
import { TopEvent } from '../types';
import { SparklineChart } from './SparklineChart';
import { api } from '../services/api';
import { LogEntry } from '../types';

interface TopEventsTableProps {
  events: TopEvent[];
}

const DrilldownModal: React.FC<{
  ruleName: string;
  onClose: () => void;
}> = ({ ruleName, onClose }) => {
  const [logs, setLogs] = useState<LogEntry[] | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  React.useEffect(() => {
    setLoading(true);
    setError(null);
    api.searchLogs(undefined, ruleName)
      .then(setLogs)
      .catch(() => setError('Failed to fetch logs'))
      .finally(() => setLoading(false));
  }, [ruleName]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-60">
      <div className="bg-splunk-gray rounded-lg shadow-lg w-full max-w-3xl max-h-[80vh] overflow-y-auto">
        <div className="flex justify-between items-center px-6 py-4 border-b border-splunk-light-gray">
          <h2 className="text-lg font-bold text-white">Logs for: {ruleName}</h2>
          <button onClick={onClose} className="text-gray-400 hover:text-white text-2xl">&times;</button>
        </div>
        <div className="px-6 py-4">
          {loading && <div className="text-white">Loading...</div>}
          {error && <div className="text-red-400">{error}</div>}
          {logs && logs.length === 0 && <div className="text-gray-400">No logs found for this event.</div>}
          {logs && logs.length > 0 && (
            <table className="w-full text-sm mt-2">
              <thead className="bg-splunk-darker">
                <tr>
                  <th className="px-2 py-2 text-xs text-gray-400">Timestamp</th>
                  <th className="px-2 py-2 text-xs text-gray-400">Level</th>
                  <th className="px-2 py-2 text-xs text-gray-400">Source IP</th>
                  <th className="px-2 py-2 text-xs text-gray-400">Message</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-splunk-light-gray">
                {logs.map((log, idx) => (
                  <tr key={idx} className="hover:bg-splunk-darker">
                    <td className="px-2 py-2 text-white font-mono">{new Date(log.timestamp).toLocaleString()}</td>
                    <td className="px-2 py-2 text-white">{log.level}</td>
                    <td className="px-2 py-2 text-white font-mono">{log.sourceIP}</td>
                    <td className="px-2 py-2 text-white">{log.message}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      </div>
    </div>
  );
};

export const TopEventsTable: React.FC<TopEventsTableProps> = ({ events }) => {
  const [currentPage, setCurrentPage] = useState(1);
  const itemsPerPage = 5;
  const totalPages = Math.ceil(events.length / itemsPerPage);
  const [drilldown, setDrilldown] = useState<string | null>(null);

  const startIndex = (currentPage - 1) * itemsPerPage;
  const endIndex = startIndex + itemsPerPage;
  const currentEvents = events.slice(startIndex, endIndex);

  const getUrgencyColor = (urgency: string) => {
    switch (urgency) {
      case 'critical':
        return 'text-red-400';
      case 'high':
        return 'text-orange-400';
      case 'medium':
        return 'text-blue-400';
      case 'low':
        return 'text-green-400';
      default:
        return 'text-gray-400';
    }
  };

  const getUrgencyBadge = (urgency: string) => {
    const colorClass = getUrgencyColor(urgency);
    return (
      <span className={`px-2 py-1 rounded-full text-xs font-medium bg-opacity-20 ${colorClass} bg-current`}>
        {urgency.toUpperCase()}
      </span>
    );
  };

  return (
    <div className="bg-splunk-gray rounded-lg border border-splunk-light-gray">
      <div className="px-6 py-4 border-b border-splunk-light-gray">
        <h3 className="text-lg font-semibold text-white">Top Notable Events</h3>
      </div>
      
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead className="bg-splunk-darker">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Rule Name
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Trend
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Count
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-400 uppercase tracking-wider">
                Urgency
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-splunk-light-gray">
            {currentEvents.map((event, index) => (
              <tr
                key={index}
                className="hover:bg-splunk-darker cursor-pointer"
                onClick={() => setDrilldown(event.ruleName)}
              >
                <td className="px-6 py-4 whitespace-nowrap text-sm text-white">
                  {event.ruleName}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <SparklineChart data={event.sparkline} />
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-white">
                  {event.count.toLocaleString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {getUrgencyBadge(event.urgency)}
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
              Showing {startIndex + 1} to {Math.min(endIndex, events.length)} of {events.length} results
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
      {drilldown && (
        <DrilldownModal ruleName={drilldown} onClose={() => setDrilldown(null)} />
      )}
    </div>
  );
}; 