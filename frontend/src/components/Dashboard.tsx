import React, { useState, useEffect } from 'react';
import { StatTile } from './StatTile';
import { UrgencyChart } from './UrgencyChart';
import { TimelineChart } from './TimelineChart';
import { TopEventsTable } from './TopEventsTable';
import { TopSourcesTable } from './TopSourcesTable';
import { LogSearch } from './LogSearch';
import { api } from '../services/api';
import { SummaryStats, UrgencyData, TimelineData, TopEvent, TopSource, LogEntry } from '../types';
import { Shield, Activity, AlertTriangle, Users } from 'lucide-react';

const REFRESH_OPTIONS = [5, 10, 15, 30];

export const Dashboard: React.FC = () => {
  const [summaryStats, setSummaryStats] = useState<SummaryStats | null>(null);
  const [urgencyData, setUrgencyData] = useState<UrgencyData | null>(null);
  const [timelineData, setTimelineData] = useState<TimelineData | null>(null);
  const [topEvents, setTopEvents] = useState<TopEvent[]>([]);
  const [topSources, setTopSources] = useState<TopSource[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshInterval, setRefreshInterval] = useState(30); // seconds

  // LogSearch state
  const [searchIp, setSearchIp] = useState('');
  const [searchEvent, setSearchEvent] = useState('');
  const [searchResults, setSearchResults] = useState<LogEntry[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchError, setSearchError] = useState<string | null>(null);

  // LogSearch handler
  const handleSearch = async (ip: string, event: string) => {
    setSearchLoading(true);
    setSearchError(null);
    try {
      const logs = await api.searchLogs(ip, event);
      setSearchResults(logs);
    } catch (err) {
      setSearchError('Failed to search logs');
    } finally {
      setSearchLoading(false);
    }
  };

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const [
          stats,
          urgency,
          timeline,
          events,
          sources
        ] = await Promise.all([
          api.getSummaryStats(),
          api.getUrgencyData(),
          api.getTimelineData(),
          api.getTopEvents(),
          api.getTopSources()
        ]);

        setSummaryStats(stats);
        setUrgencyData(urgency);
        setTimelineData(timeline);
        setTopEvents(events);
        setTopSources(sources);
        setError(null);
      } catch (err) {
        setError('Failed to load dashboard data. Please check if the backend server is running.');
        console.error('Dashboard data fetch error:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
    
    const interval = setInterval(fetchData, refreshInterval * 1000);
    return () => clearInterval(interval);
  }, [refreshInterval]);

  // Home button handler
  const handleHome = () => {
    window.scrollTo({ top: 0, behavior: 'smooth' });
    // Optionally clear search state:
    // setSearchIp(''); setSearchEvent(''); setSearchResults([]); setSearchError(null);
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-splunk-dark flex items-center justify-center">
        <div className="text-white text-xl">Loading dashboard...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-splunk-dark flex items-center justify-center">
        <div className="text-red-400 text-xl text-center max-w-md">
          {error}
          <br />
          <button 
            onClick={() => window.location.reload()} 
            className="mt-4 px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-splunk-dark">
      {/* Header */}
      <div className="bg-splunk-darker border-b border-splunk-light-gray">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center space-x-4">
              <Shield className="h-8 w-8 text-blue-400 mr-3" />
              <h1 className="text-xl font-bold text-white">Security Dashboard</h1>
              <button
                onClick={handleHome}
                className="ml-4 px-3 py-1 bg-splunk-gray text-white rounded hover:bg-splunk-light-gray border border-splunk-light-gray"
              >
                Home
              </button>
            </div>
            <div className="flex items-center space-x-4">
              <div className="text-sm text-gray-400">
                Last updated: {new Date().toLocaleTimeString()}
              </div>
              <div className="flex items-center">
                <span className="text-xs text-gray-400 mr-2">Refresh:</span>
                <select
                  className="bg-splunk-gray text-white border border-splunk-light-gray rounded px-2 py-1 text-xs"
                  value={refreshInterval}
                  onChange={e => setRefreshInterval(Number(e.target.value))}
                >
                  {REFRESH_OPTIONS.map(opt => (
                    <option key={opt} value={opt}>{opt}s</option>
                  ))}
                </select>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Log Search */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-8">
        <LogSearch
          ip={searchIp}
          setIp={setSearchIp}
          event={searchEvent}
          setEvent={setSearchEvent}
          results={searchResults}
          loading={searchLoading}
          error={searchError}
          onSearch={handleSearch}
        />
      </div>

      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Summary Stats */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <StatTile
            title="Access Notables"
            total={summaryStats?.accessNotables.total || 0}
            delta={summaryStats?.accessNotables.delta || 0}
            color="blue"
          />
          <StatTile
            title="Network Notables"
            total={summaryStats?.networkNotables.total || 0}
            delta={summaryStats?.networkNotables.delta || 0}
            color="green"
          />
          <StatTile
            title="Threat Notables"
            total={summaryStats?.threatNotables.total || 0}
            delta={summaryStats?.threatNotables.delta || 0}
            color="red"
          />
          <StatTile
            title="UBA Notables"
            total={summaryStats?.ubaNotables.total || 0}
            delta={summaryStats?.ubaNotables.delta || 0}
            color="purple"
          />
        </div>

        {/* Charts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          <UrgencyChart data={urgencyData!} />
          <TimelineChart data={timelineData!} />
        </div>

        {/* Tables */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <TopEventsTable events={topEvents} />
          <TopSourcesTable sources={topSources} />
        </div>
      </div>
    </div>
  );
}; 