export interface QueryDigest {
  classes: Class[];
  global: Global;
}

export interface Class {
  attribute: string;
  checksum: string;
  distillate: string;
  example: Example;
  fingerprint: string;
  histograms: Histograms;
  metrics: Metrics;
  query_count: number;
  tables?: Table[];
  ts_max: string;
  ts_min: string;
}

export interface Example {
  Query_time: string;
  query: string;
  ts: string;
  as_select?: string;
}

export interface Histograms {
  Query_time: number[];
}

export interface Metrics {
  Lock_time: MetricsRow;
  Query_length: MetricsRow;
  Query_time: MetricsRow;
  Rows_examined: MetricsRow;
  Rows_sent: MetricsRow;
  // 使わない & 型合わせのため
  // host: Host;
  // user: User;
}

export interface MetricsRow {
  avg: string;
  max: string;
  median: string;
  min: string;
  pct: string;
  pct_95: string;
  stddev: string;
  sum: string;
}

export interface Host {
  value: string;
}

export interface User {
  value: string;
}

export interface Table {
  create: string;
  status: string;
}

export interface Global {
  files: File[];
  metrics: Metrics2;
  query_count: number;
  unique_query_count: number;
}

export interface File {
  name: string;
  size: number;
}

export interface Metrics2 {
  Lock_time: MetricsRow2;
  Query_length: MetricsRow2;
  Query_time: MetricsRow2;
  Rows_examined: MetricsRow2;
  Rows_sent: MetricsRow2;
}

export interface MetricsRow2 {
  avg: string;
  max: string;
  median: string;
  min: string;
  pct_95: string;
  stddev: string;
  sum: string;
}
