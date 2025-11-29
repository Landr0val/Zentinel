export interface Pagination {
  page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
}

export interface ApiResponse<T> {
  data: T;
  pagination?: Pagination;
  error?: {
    code: string;
    message: string;
    details?: unknown;
  };
}

export type ClientStatus = 'active' | 'inactive' | 'blocked';
export type RiskLevel = 'low' | 'medium' | 'high';

export interface Client {
  id: string;
  full_name: string;
  email: string;
  phone: string;
  address: string;
  city: string;
  state: string;
  zip_code: string;
  status: ClientStatus;
  risk_level: RiskLevel;
  created_at: string;
  updated_at: string;
}

export type AccountType = 'savings' | 'checking';
export type AccountStatus = 'active' | 'frozen' | 'closed';

export interface Account {
  id: string;
  client_id: string;
  account_number: string;
  type: AccountType;
  currency: string;
  balance: number;
  status: AccountStatus;
  created_at: string;
  updated_at: string;
  // Optional joined fields often returned by list endpoints
  client_name?: string;
}

export type TransactionType = 'deposit' | 'withdrawal' | 'transfer' | 'payment';
export type TransactionStatus = 'pending' | 'completed' | 'failed' | 'flagged';

export interface Transaction {
  id: string;
  account_id: string;
  type: TransactionType;
  amount: number;
  currency: string;
  merchant: string;
  location: string;
  status: TransactionStatus;
  risk_score: number;
  analysis_result?: string;
  created_at: string;
  // Optional joined fields
  account_number?: string;
  client_name?: string;
}

export type AlertSeverity = 'low' | 'medium' | 'high' | 'critical';
export type AlertStatus = 'new' | 'investigating' | 'resolved' | 'false_positive';

export interface Alert {
  id: string;
  transaction_id: string;
  type: string;
  severity: AlertSeverity;
  status: AlertStatus;
  description: string;
  created_at: string;
  updated_at: string;
  // Optional joined fields
  transaction_amount?: number;
  client_name?: string;
}

export interface DashboardStats {
  total_volume: number;
  transaction_count: number;
  alert_count: number;
  active_clients: number;
  volume_change_percentage: number;
  transaction_change_percentage: number;
  alert_change_percentage: number;
  client_change_percentage: number;
  recent_alerts: Alert[];
  weekly_activity: {
    date: string;
    transactions: number;
    alerts: number;
  }[];
}