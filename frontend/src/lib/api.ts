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

  const token = typeof window !== "undefined" ? localStorage.getItem("token") : null;

  const headers = {
    "Content-Type": "application/json",
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    ...options?.headers,
  };

  const response = await fetch(url, {
    ...options,
    headers,
  });

  if (!response.ok) {
    if ((response.status === 401 || response.status === 403) && typeof window !== "undefined") {
      localStorage.removeItem("token");
      window.location.href = "/login";
    }

    const errorData = await response.json().catch(() => ({}));
    throw new Error(
      errorData.error?.message || `API Error: ${response.status} ${response.statusText}`
    );
  }

  return response.json();
}

export const api = {
  dashboard: {
    getStats: () => fetcher<ApiResponse<DashboardStats>>("/api/v1/dashboard/stats"),
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
      return fetcher<ApiResponse<Client[]>>(`/api/v1/clients?${searchParams.toString()}`);
    },
    get: (id: string) => fetcher<ApiResponse<Client>>(`/api/v1/clients/${id}`),
    create: (data: CreateClientRequest) =>
      fetcher<ApiResponse<Client>>("/api/v1/clients", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    update: (id: string, data: Partial<Client>) =>
      fetcher<ApiResponse<Client>>(`/api/v1/clients/${id}`, {
        method: "PUT",
        body: JSON.stringify(data),
      }),
    getAccounts: (id: string) =>
      fetcher<ApiResponse<Account[]>>(`/api/v1/clients/${id}/accounts`),
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
      return fetcher<ApiResponse<Account[]>>(`/api/v1/accounts?${searchParams.toString()}`);
    },
    get: (id: string) => fetcher<ApiResponse<Account>>(`/api/v1/accounts/${id}`),
    create: (data: Partial<Account>) =>
      fetcher<ApiResponse<Account>>("/api/v1/accounts", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    update: (id: string, data: Partial<Account>) =>
      fetcher<ApiResponse<Account>>(`/api/v1/accounts/${id}`, {
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
        `/api/v1/transactions?${searchParams.toString()}`
      );
    },
    get: (id: string) => fetcher<ApiResponse<Transaction>>(`/api/v1/transactions/${id}`),
    create: (data: CreateTransactionRequest) =>
      fetcher<ApiResponse<Transaction>>("/api/v1/transactions", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    analyze: (id: string) =>
      fetcher<ApiResponse<{ risk_score: number; analysis: string }>>(
        `/api/v1/transactions/${id}/analyze`,
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
      return fetcher<ApiResponse<Alert[]>>(`/api/v1/alerts?${searchParams.toString()}`);
    },
    get: (id: string) => fetcher<ApiResponse<Alert>>(`/api/v1/alerts/${id}`),
    create: (data: Partial<Alert>) =>
      fetcher<ApiResponse<Alert>>("/api/v1/alerts", {
        method: "POST",
        body: JSON.stringify(data),
      }),
    updateStatus: (id: string, status: string) =>
      fetcher<ApiResponse<Alert>>(`/api/v1/alerts/${id}/status`, {
        method: "PATCH",
        body: JSON.stringify({ status }),
      }),
  },
};
