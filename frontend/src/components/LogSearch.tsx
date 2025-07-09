import React from 'react';
import { LogEntry } from '../types';

interface LogSearchProps {
  ip: string;
  setIp: (ip: string) => void;
  event: string;
  setEvent: (event: string) => void;
  results: LogEntry[];
  loading: boolean;
  error: string | null;
  onSearch: (ip: string, event: string) => void;
}

export const LogSearch: React.FC<LogSearchProps> = ({
  ip,
  setIp,
  event,
  setEvent,
  results,
  loading,
  error,
  onSearch,
}) => {
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSearch(ip, event);
  };

  return (
    <div className="bg-splunk-gray rounded-lg border border-splunk-light-gray mb-8">
      <div className="px-6 py-4 border-b border-splunk-light-gray">
        <h3 className="text-lg font-semibold text-white">Search Logs</h3>
      </div>
      <form className="px-6 py-4 flex flex-col md:flex-row gap-4 items-end" onSubmit={handleSubmit}>
        <div className="flex flex-col">
          <label className="text-gray-400 text-xs mb-1">Source IP</label>
          <input
            type="text"
            className="bg-splunk-darker border border-splunk-light-gray rounded px-3 py-2 text-white"
            value={ip}
            onChange={e => setIp(e.target.value)}
            placeholder="e.g. 192.168.1.100"
          />
        </div>
        <div className="flex flex-col">
          <label className="text-gray-400 text-xs mb-1">Event/Rule Name</label>
          <input
            type="text"
            className="bg-splunk-darker border border-splunk-light-gray rounded px-3 py-2 text-white"
            value={event}
            onChange={e => setEvent(e.target.value)}
            placeholder="e.g. Suspicious Login"
          />
        </div>
        <button
          type="submit"
          className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded font-semibold"
          disabled={loading}
        >
          {loading ? 'Searching...' : 'Search'}
        </button>
      </form>
      {error && <div className="px-6 pb-4 text-red-400">{error}</div>}
      <div className="overflow-x-auto px-6 pb-6">
        {results.length > 0 ? (
          <table className="w-full mt-4">
            <thead className="bg-splunk-darker">
              <tr>
                <th className="px-3 py-2 text-xs text-gray-400">Timestamp</th>
                <th className="px-3 py-2 text-xs text-gray-400">Level</th>
                <th className="px-3 py-2 text-xs text-gray-400">Rule Name</th>
                <th className="px-3 py-2 text-xs text-gray-400">Source IP</th>
                <th className="px-3 py-2 text-xs text-gray-400">Message</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-splunk-light-gray">
              {results.map((log, idx) => (
                <tr key={idx} className="hover:bg-splunk-darker">
                  <td className="px-3 py-2 text-sm text-white font-mono">{new Date(log.timestamp).toLocaleString()}</td>
                  <td className="px-3 py-2 text-sm text-white">{log.level}</td>
                  <td className="px-3 py-2 text-sm text-white">{log.ruleName}</td>
                  <td className="px-3 py-2 text-sm text-white font-mono">{log.sourceIP}</td>
                  <td className="px-3 py-2 text-sm text-white">{log.message}</td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <div className="text-gray-400 mt-4">No results found.</div>
        )}
      </div>
    </div>
  );
}; 