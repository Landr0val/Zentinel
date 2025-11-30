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
  document_number: string;
  document_type: string;
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

export interface CreateClientRequest {
  document_type: string;
  document_number: string;
  full_name: string;
  email: string;
  phone: string;
  initial_account?: {
    account_number: string;
    account_type: AccountType;
    currency: string;
  };
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
  client_name?: string;
}

export interface CreateAccountRequest {
  client_id: string;
  account_number: string;
  account_type: AccountType;
  currency: string;
}

export type TransactionType = 'deposit' | 'withdrawal' | 'transfer' | 'purchase';
export type TransactionStatus = 'pending' | 'completed' | 'failed' | 'flagged';
export type TransactionChannel = 'mobile' | 'web' | 'atm' | 'branch';

export interface CreateTransactionRequest {
  account_id: string;
  amount: number;
  currency: string;
  operation_type: TransactionType;
  channel: TransactionChannel;
  merchant?: string;
  country?: string;
  city?: string;
}

export interface Transaction {
  id: string;
  code: string;
  account_id: string;
  operation_type: TransactionType;
  channel: TransactionChannel;
  amount: number;
  currency: string;
  merchant: string;
  country: string;
  city: string;
  status: TransactionStatus;
  risk_score: number;
  is_flagged: boolean;
  created_at: string;
  // Optional joined fields
  account_number?: string;
  client_name?: string;
}

export type AlertSeverity = 'low' | 'medium' | 'high' | 'critical';
export type AlertStatus = 'pending' | 'reviewing' | 'resolved' | 'false_positive';

export interface Alert {
  id: string;
  code: string;
  transaction_id?: string;
  transaction_code?: string;
  client_id: string;
  alert_type: string;
  severity: AlertSeverity;
  status: AlertStatus;
  description: string;
  ai_explanation?: string;
  reviewed_by?: string;
  reviewed_at?: string;
  blockchain_tx?: string;
  created_at: string;
  // Optional joined fields
  transaction_amount?: number;
  client_name?: string;
}

export interface DashboardStats {
  total_volume: number;
  total_transactions: number;
  flagged_count: number;
  flagged_volume: number;
  alerts_by_status: Record<string, number>;
  transaction_count?: number;
  alert_count?: number;
  active_clients?: number;
  volume_change_percentage?: number;
  transaction_change_percentage?: number;
  alert_change_percentage?: number;
  client_change_percentage?: number;
  recent_alerts?: Alert[];
  weekly_activity?: {
    date: string;
    transactions: number;
    alerts: number;
  }[];
}