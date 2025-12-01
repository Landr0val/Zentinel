"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import {
  ArrowLeft,
  CreditCard,
  DollarSign,
  Calendar,
  Activity,
  ArrowUpRight,
  ArrowDownRight,
  Loader2,
} from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import { createSlug, getIdFromSlug } from "../../../lib/slug-manager";
import { useCurrency } from "../../../context/CurrencyContext";

export default function AccountDetailPage() {
  const params = useParams();
  const slug = params?.id as string;
  const id = getIdFromSlug(slug) || slug;
  const { formatAmount } = useCurrency();

  const {
    data: accountData,
    isLoading: isLoadingAccount,
    error: errorAccount,
  } = useQuery({
    queryKey: ["account", id],
    queryFn: () => api.accounts.get(id),
  });

  const account = accountData?.data;

  const { data: clientData } = useQuery({
    queryKey: ["client", account?.client_id],
    queryFn: () => api.clients.get(account!.client_id),
    enabled: !!account?.client_id,
  });
  const client = clientData?.data;

  const { data: transactionsData, isLoading: isLoadingTransactions } = useQuery(
    {
      queryKey: ["accountTransactions", account?.id],
      queryFn: () => api.transactions.list({ account_id: account?.id }),
      enabled: !!account?.id,
    },
  );

  const transactions = transactionsData?.data || [];

  if (isLoadingAccount) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (errorAccount || !account) {
    return (
      <div className="rounded-xl border border-destructive/20 bg-destructive/5 p-6 text-destructive">
        <h3 className="text-lg font-medium">Error al cargar cuenta</h3>
        <p>No se pudo encontrar la información de la cuenta.</p>
        <Link
          href="/accounts"
          className="mt-4 inline-block text-sm font-medium underline hover:text-destructive/80"
        >
          Volver a la lista
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <Link
          href="/accounts"
          className="p-2 rounded-full hover:bg-secondary text-muted-foreground transition-colors"
        >
          <ArrowLeft className="h-5 w-5" />
        </Link>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-foreground">
            Cuenta {account.account_number}
          </h1>
          <p className="text-muted-foreground">
            {client
              ? `${client.document_type}: ${client.document_number}`
              : `ID: ${id}`}
          </p>
        </div>
        <div className="flex gap-2">
          <span
            className={`inline-flex items-center rounded-full px-3 py-1 text-sm font-medium ${
              account.status === "active"
                ? "bg-secondary text-foreground ring-1 ring-inset ring-border"
                : "bg-secondary/50 text-muted-foreground ring-1 ring-inset ring-border/50"
            }`}
          >
            {account.status === "active" ? "Activa" : "Inactiva"}
          </span>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Summary Card */}
        <div className="lg:col-span-2 space-y-6">
          <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
            <div className="rounded-xl bg-card p-6 shadow-zen border border-border">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-sm font-medium text-muted-foreground">
                  Saldo Disponible
                </h3>
                <DollarSign className="h-5 w-5 text-muted-foreground" />
              </div>
              <div className="text-3xl font-bold text-foreground">
                {formatAmount(account.balance, account.currency)}
              </div>
            </div>
            <div className="rounded-xl bg-card p-6 shadow-zen border border-border">
              <div className="flex items-center justify-between mb-4">
                <h3 className="text-sm font-medium text-muted-foreground">
                  Tipo de Cuenta
                </h3>
                <CreditCard className="h-5 w-5 text-muted-foreground" />
              </div>
              <div className="text-3xl font-bold text-foreground capitalize">
                {account.type === "savings" ? "Ahorro" : "Corriente"}
              </div>
              <p className="text-sm text-muted-foreground mt-1">
                Banca Personal
              </p>
            </div>
          </div>

          {/* Transactions List */}
          <div className="rounded-xl bg-card shadow-zen border border-border overflow-hidden">
            <div className="px-6 py-5 border-b border-border flex items-center justify-between">
              <h3 className="text-lg font-semibold text-foreground">
                Últimos Movimientos
              </h3>
              <Link
                href="/transactions"
                className="text-sm font-medium text-primary hover:text-primary/80"
              >
                Ver historial completo
              </Link>
            </div>
            {isLoadingTransactions ? (
              <div className="flex justify-center py-8">
                <Loader2 className="h-6 w-6 animate-spin text-primary" />
              </div>
            ) : transactions.length === 0 ? (
              <div className="px-6 py-8 text-center text-muted-foreground">
                No hay movimientos recientes.
              </div>
            ) : (
              <ul role="list" className="divide-y divide-border">
                {transactions.slice(0, 5).map((txn) => (
                  <li
                    key={txn.id}
                    className="px-6 py-4 hover:bg-secondary/50 transition-colors"
                  >
                    <div className="flex items-center justify-between">
                      <div className="flex items-center">
                        <div
                          className={`rounded-full p-2 mr-4 ${
                            txn.operation_type === "deposit"
                              ? "bg-secondary text-foreground"
                              : "bg-secondary/50 text-muted-foreground"
                          }`}
                        >
                          {txn.operation_type === "deposit" ? (
                            <ArrowDownRight className="h-4 w-4" />
                          ) : (
                            <ArrowUpRight className="h-4 w-4" />
                          )}
                        </div>
                        <div>
                          <p className="text-sm font-medium text-foreground">
                            {txn.merchant}
                          </p>
                          <p className="text-xs text-muted-foreground">
                            {new Date(txn.created_at).toLocaleString()}
                          </p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p
                          className={`text-sm font-medium ${
                            txn.operation_type === "deposit"
                              ? "text-foreground"
                              : "text-foreground"
                          }`}
                        >
                          {txn.operation_type === "deposit" ? "+" : "-"}
                          {formatAmount(txn.amount, txn.currency)}
                        </p>
                        <p className="text-xs text-muted-foreground capitalize">
                          {txn.operation_type}
                        </p>
                      </div>
                    </div>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>

        {/* Sidebar Info */}
        <div className="space-y-6">
          <div className="rounded-xl bg-card p-6 shadow-zen border border-border">
            <h3 className="text-lg font-semibold text-foreground mb-4">
              Detalles del Cliente
            </h3>
            <div className="flex items-center mb-6">
              <div className="h-12 w-12 rounded-full bg-secondary flex items-center justify-center text-foreground font-bold text-lg">
                {account.client_name?.charAt(0) || "?"}
              </div>
              <div className="ml-4">
                <div className="text-base font-medium text-foreground">
                  {account.client_name || "Desconocido"}
                </div>
                <Link
                  href={`/clients/${createSlug(account.client_id)}`}
                  className="text-sm text-primary hover:text-primary/80"
                >
                  Ver perfil
                </Link>
              </div>
            </div>
            <div className="space-y-4 border-t border-border pt-4">
              <div className="flex items-center justify-between">
                <div className="flex items-center text-sm text-muted-foreground">
                  <Calendar className="mr-2 h-4 w-4" />
                  Apertura
                </div>
                <span className="text-sm font-medium text-foreground">
                  {new Date(account.created_at).toLocaleDateString()}
                </span>
              </div>
              <div className="flex items-center justify-between">
                <div className="flex items-center text-sm text-muted-foreground">
                  <Activity className="mr-2 h-4 w-4" />
                  Última Actividad
                </div>
                <span className="text-sm font-medium text-foreground">
                  {new Date(account.updated_at).toLocaleDateString()}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
