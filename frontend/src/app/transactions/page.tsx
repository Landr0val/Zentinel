"use client";

import Link from "next/link";
import {
  Search,
  Filter,
  MoreHorizontal,
  ArrowUpRight,
  ArrowDownRight,
  AlertCircle,
  Check,
  Clock,
  Loader2,
  ArrowRightLeft,
  Plus,
} from "lucide-react";
import PageHeader from "../../components/PageHeader";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../lib/api";
import { useState } from "react";
import { cn } from "../../lib/utils";
import { createSlug } from "../../lib/slug-manager";
import { useCurrency } from "../../context/CurrencyContext";

export default function TransactionsPage() {
  const [searchTerm, setSearchTerm] = useState("");
  const { formatAmount } = useCurrency();

  const { data, isLoading, error } = useQuery({
    queryKey: ["transactions", searchTerm],
    queryFn: () =>
      api.transactions.list(searchTerm ? { search: searchTerm } : undefined),
  });

  const transactions = data?.data || [];

  return (
    <div className="animate-in fade-in duration-500">
      <PageHeader
        title="Transacciones"
        description="Historial y monitoreo de movimientos."
        action={
          <Link
            href="/transactions/new"
            className="inline-flex items-center justify-center rounded-xl bg-foreground px-4 py-2 text-sm font-medium text-background shadow-sm hover:bg-foreground/90 transition-colors"
          >
            <Plus className="mr-2 h-4 w-4" />
            Nueva Transacción
          </Link>
        }
      />

      {/* Filters */}
      <div className="mb-8 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="relative flex-1 max-w-md group">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
            <Search
              className="h-4 w-4 text-muted-foreground group-focus-within:text-foreground transition-colors"
              aria-hidden="true"
            />
          </div>
          <input
            type="text"
            className="block w-full rounded-xl border border-border bg-card py-2.5 pl-10 text-sm text-foreground placeholder:text-muted-foreground focus:border-foreground focus:ring-0 transition-all shadow-sm hover:border-foreground/50 outline-none"
            placeholder="Buscar por ID, cuenta o comercio..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <div className="flex gap-3">
          <div className="relative">
            <select className="appearance-none block w-full rounded-xl border border-border bg-card py-2.5 pl-4 pr-10 text-sm text-foreground focus:border-foreground focus:ring-0 transition-all shadow-sm hover:border-foreground/50 cursor-pointer outline-none">
              <option>Todos los tipos</option>
              <option>Pagos</option>
              <option>Transferencias</option>
              <option>Retiros</option>
              <option>Depósitos</option>
            </select>
            <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-3 text-muted-foreground">
              <Filter className="h-3 w-3" />
            </div>
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-hidden rounded-3xl border border-border bg-card shadow-zen">
        {isLoading ? (
          <div className="flex h-64 items-center justify-center">
            <Loader2 className="h-6 w-6 animate-spin text-primary" />
          </div>
        ) : error ? (
          <div className="p-12 text-center">
            <div className="inline-flex h-10 w-10 items-center justify-center rounded-full bg-destructive/10 text-destructive mb-4">
              <AlertCircle className="h-5 w-5" />
            </div>
            <h3 className="text-lg font-medium text-foreground">
              Error al cargar
            </h3>
            <p className="text-muted-foreground mt-1">
              No se pudieron obtener las transacciones.
            </p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="min-w-full text-left text-sm">
                <thead className="bg-secondary/30 text-muted-foreground font-medium">
                  <tr>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Transacción
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Cuenta
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Monto
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Estado
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Riesgo
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Fecha
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-4 font-medium text-right"
                    >
                      Acciones
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-border">
                  {transactions.length === 0 ? (
                    <tr>
                      <td
                        colSpan={7}
                        className="px-6 py-16 text-center text-muted-foreground"
                      >
                        No se encontraron transacciones.
                      </td>
                    </tr>
                  ) : (
                    transactions.map((txn) => (
                      <tr
                        key={txn.id}
                        className="group hover:bg-secondary/40 transition-colors"
                      >
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-3">
                            <div
                              className={cn(
                                "flex h-8 w-8 items-center justify-center rounded-full border transition-colors",
                                txn.operation_type === "deposit"
                                  ? "border-border bg-secondary/50 text-foreground"
                                  : "border-border bg-card text-muted-foreground",
                              )}
                            >
                              {txn.operation_type === "deposit" ? (
                                <ArrowDownRight className="h-4 w-4" />
                              ) : txn.operation_type === "withdrawal" ? (
                                <ArrowUpRight className="h-4 w-4" />
                              ) : (
                                <ArrowRightLeft className="h-4 w-4" />
                              )}
                            </div>
                            <div>
                              <div className="font-medium text-foreground">
                                {txn.merchant}
                              </div>
                              <div className="text-xs text-muted-foreground">
                                <span className="font-mono mr-2 opacity-70">
                                  {txn.code}
                                </span>
                                <span className="capitalize">
                                  {txn.operation_type}
                                </span>
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="font-mono text-xs text-muted-foreground bg-secondary/50 px-2 py-1 rounded-md inline-block">
                            {txn.account_number || "N/A"}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="font-medium text-foreground">
                            {txn.operation_type === "deposit" ? "+" : "-"}
                            {formatAmount(txn.amount, txn.currency)}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-2">
                            {txn.status === "completed" && (
                              <Check className="h-3 w-3 text-foreground" />
                            )}
                            {txn.status === "pending" && (
                              <Clock className="h-3 w-3 text-muted-foreground" />
                            )}
                            {(txn.status === "flagged" ||
                              txn.status === "failed") && (
                              <AlertCircle className="h-3 w-3 text-destructive" />
                            )}
                            <span
                              className={cn(
                                "text-xs font-medium",
                                txn.status === "flagged" ||
                                  txn.status === "failed"
                                  ? "text-destructive"
                                  : "text-muted-foreground",
                              )}
                            >
                              {txn.status === "completed"
                                ? "Completada"
                                : txn.status === "pending"
                                  ? "Pendiente"
                                  : txn.status === "flagged"
                                    ? "Sospechosa"
                                    : "Fallida"}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-3">
                            <div className="h-1.5 w-16 rounded-full bg-secondary overflow-hidden">
                              <div
                                className={cn(
                                  "h-full rounded-full transition-all",
                                  txn.risk_score < 20
                                    ? "bg-muted-foreground/30"
                                    : txn.risk_score < 70
                                      ? "bg-foreground"
                                      : "bg-destructive",
                                )}
                                style={{ width: `${txn.risk_score}%` }}
                              />
                            </div>
                            <span className="text-xs font-medium text-muted-foreground">
                              {txn.risk_score}%
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4 text-muted-foreground text-xs">
                          {new Date(txn.created_at).toLocaleString()}
                        </td>
                        <td className="px-6 py-4 text-right">
                          <Link
                            href={`/transactions/${createSlug(txn.id)}`}
                            className="inline-flex h-8 w-8 items-center justify-center rounded-full text-muted-foreground hover:bg-secondary hover:text-foreground transition-colors"
                          >
                            <MoreHorizontal className="h-4 w-4" />
                          </Link>
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>

            {/* Pagination */}
            {transactions.length > 0 && (
              <div className="flex items-center justify-between border-t border-border px-6 py-4">
                <div className="text-xs text-muted-foreground">
                  Mostrando{" "}
                  <span className="font-medium text-foreground">1</span> a{" "}
                  <span className="font-medium text-foreground">
                    {transactions.length}
                  </span>{" "}
                  resultados
                </div>
                <div className="flex gap-2">
                  <button className="rounded-lg border border-border px-3 py-1 text-xs font-medium text-muted-foreground hover:bg-secondary hover:text-foreground transition-colors disabled:opacity-50">
                    Anterior
                  </button>
                  <button className="rounded-lg border border-border px-3 py-1 text-xs font-medium text-muted-foreground hover:bg-secondary hover:text-foreground transition-colors disabled:opacity-50">
                    Siguiente
                  </button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
