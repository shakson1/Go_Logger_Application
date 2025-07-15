import { SummaryStats, UrgencyData, TimelineData, TopEvent, TopSource, LogEntry } from '../types';

const API_BASE_URL = '/api';

export const api = {
  async getSummaryStats(): Promise<SummaryStats> {
    const response = await fetch(`${API_BASE_URL}/summary`);
    if (!response.ok) {
      throw new Error('Failed to fetch summary stats');
    }
    return response.json();
  },

  async getUrgencyData(): Promise<UrgencyData> {
    const response = await fetch(`${API_BASE_URL}/urgency`);
    if (!response.ok) {
      throw new Error('Failed to fetch urgency data');
    }
    return response.json();
  },

  async getTimelineData(): Promise<TimelineData> {
    const response = await fetch(`${API_BASE_URL}/timeline`);
    if (!response.ok) {
      throw new Error('Failed to fetch timeline data');
    }
    return response.json();
  },

  async getTopEvents(): Promise<TopEvent[]> {
    const response = await fetch(`${API_BASE_URL}/top-events`);
    if (!response.ok) {
      throw new Error('Failed to fetch top events');
    }
    return response.json();
  },

  async getTopSources(): Promise<TopSource[]> {
    const response = await fetch(`${API_BASE_URL}/top-sources`);
    if (!response.ok) {
      throw new Error('Failed to fetch top sources');
    }
    return response.json();
  },

  async searchLogs(ip?: string, event?: string): Promise<LogEntry[]> {
    const params = new URLSearchParams();
    if (ip) params.append('ip', ip);
    if (event) params.append('event', event);
    const response = await fetch(`${API_BASE_URL}/logs?${params.toString()}`);
    if (!response.ok) {
      throw new Error('Failed to search logs');
    }
    return response.json();
  },

  async ingestLog(entry: LogEntry): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/logs`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(entry),
    });
    if (!response.ok) {
      throw new Error('Failed to ingest log');
    }
  },
}; 