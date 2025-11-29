import {
  Account,
  Alert,
  ApiResponse,
  Client,
  CreateClientRequest,
  CreateTransactionRequest,
  DashboardStats,
  Transaction,
} from "../types";

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/api/v1";

async function fetcher<T>(endpoint: string, options?: RequestInit): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  const headers = {
    "Content-Type": "application/json",
    ...options?.headers,
  };

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(
      errorData.error?.message || `API Error: ${response.status} ${response.statusText}`
    );
  }

  return response.json();
}

export const api = {
  dashboard: {
    getStats: () => fetcher<ApiResponse<DashboardStats>>("/dashboard/stats"),
  },
  clients: {
    list: (params?: Record<string, string | number | boolean | null | undefined>) => {
      const searchParams = new URLSearchParams();
      if (params) {
        Object.entries(params).forEach(([key, value]) => {
          if (value !== undefined && value !== null && value !== "") {
            searchParams.append(key, String(value));
          }
        });
      }
      return fetcher<ApiResponse<Client[]>>(`/clients?${searchParams.toString()}`);
    },
    get: (id: string) => fetcher<ApiResponse<Client>>(`/clients/${id}`),
    create: (data: CreateClientRequest) =>
      fetcher<ApiResponse<Client>>("/clients", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    update: (id: string, data: Partial<Client>) =>
      fetcher<ApiResponse<Client>>(`/clients/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    getAccounts: (id: string) =>
      fetcher<ApiResponse<Account[]>>(`/clients/${id}/accounts`),
  },
  accounts: {
    list: (params?: Record<string, string | number | boolean | null | undefined>) => {
      const searchParams = new URLSearchParams();
      if (params) {
        Object.entries(params).forEach(([key, value]) => {
          if (value !== undefined && value !== null && value !== "") {
            searchParams.append(key, String(value));
          }
        });
      }
      return fetcher<ApiResponse<Account[]>>(`/accounts?${searchParams.toString()}`);
    },
    get: (id: string) => fetcher<ApiResponse<Account>>(`/accounts/${id}`),
    create: (data: Partial<Account>) =>
      fetcher<ApiResponse<Account>>("/accounts", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    update: (id: string, data: Partial<Account>) =>
      fetcher<ApiResponse<Account>>(`/accounts/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
  },
  transactions: {
    list: (params?: Record<string, string | number | boolean | null | undefined>) => {
      const searchParams = new URLSearchParams();
      if (params) {
        Object.entries(params).forEach(([key, value]) => {
          if (value !== undefined && value !== null && value !== "") {
            searchParams.append(key, String(value));
          }
        });
      }
      return fetcher<ApiResponse<Transaction[]>>(
        `/transactions?${searchParams.toString()}`
      );
    },
    get: (id: string) => fetcher<ApiResponse<Transaction>>(`/transactions/${id}`),
    create: (data: CreateTransactionRequest) =>
      fetcher<ApiResponse<Transaction>>("/transactions", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    analyze: (id: string) =>
      fetcher<ApiResponse<{ risk_score: number; analysis: string }>>(
        `/transactions/${id}/analyze`,
        {
          method: "POST",
        }
      ),
  },
  alerts: {
    list: (params?: Record<string, string | number | boolean | null | undefined>) => {
      const searchParams = new URLSearchParams();
      if (params) {
        Object.entries(params).forEach(([key, value]) => {
          if (value !== undefined && value !== null && value !== "") {
            searchParams.append(key, String(value));
          }
        });
      }
      return fetcher<ApiResponse<Alert[]>>(`/alerts?${searchParams.toString()}`);
    },
    get: (id: string) => fetcher<ApiResponse<Alert>>(`/alerts/${id}`),
    create: (data: Partial<Alert>) =>
      fetcher<ApiResponse<Alert>>("/alerts", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    updateStatus: (id: string, status: string) =>
      fetcher<ApiResponse<Alert>>(`/alerts/${id}/status`, {
        method: "PATCH",
        body: JSON.stringify({ status }),
      }),
  },
};