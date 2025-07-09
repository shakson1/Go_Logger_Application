export interface StatTile {
  total: number;
  delta: number;
}

export interface SummaryStats {
  accessNotables: StatTile;
  networkNotables: StatTile;
  threatNotables: StatTile;
  ubaNotables: StatTile;
}

export interface UrgencyData {
  critical: number;
  high: number;
  medium: number;
  low: number;
}

export interface TimelineSeries {
  name: string;
  data: number[];
  color: string;
}

export interface TimelineData {
  labels: string[];
  series: TimelineSeries[];
}

export interface TopEvent {
  ruleName: string;
  sparkline: number[];
  count: number;
  urgency: string;
}

export interface TopSource {
  sourceIP: string;
  sparkline: number[];
  count: number;
  category: string;
}

export interface LogEntry {
  timestamp: string;
  level: string;
  message: string;
  ruleName: string;
  sourceIP: string;
  metadata: Record<string, string>;
} 