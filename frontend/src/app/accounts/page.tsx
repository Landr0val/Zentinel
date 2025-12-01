"use client";

import Link from "next/link";
import {
  Search,
  Filter,
  MoreHorizontal,
  CreditCard,
  Loader2,
  Wallet,
  Plus,
} from "lucide-react";
import PageHeader from "../../components/PageHeader";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../lib/api";
import { useState } from "react";
import { cn } from "../../lib/utils";
import { createSlug } from "../../lib/slug-manager";

export default function AccountsPage() {
  const [searchTerm, setSearchTerm] = useState("");

  const { data, isLoading, error } = useQuery({
    queryKey: ["accounts", searchTerm],
    queryFn: () =>
      api.accounts.list(searchTerm ? { search: searchTerm } : undefined),
  });

  const accounts = data?.data || [];

  return (
    <div className="animate-in fade-in duration-500">
      <PageHeader
        title="Cuentas"
        description="Gestión de productos y saldos."
        action={
          <Link
            href="/accounts/new"
            className="inline-flex items-center justify-center rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow-zen hover:bg-primary/90 transition-colors"
          >
            <Plus className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
            Nueva Cuenta
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
            placeholder="Buscar por número de cuenta o cliente..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <div className="flex gap-3">
          <button className="inline-flex items-center gap-2 rounded-xl border border-border bg-card px-4 py-2.5 text-sm font-medium text-foreground shadow-sm hover:bg-secondary transition-colors">
            <Filter className="h-4 w-4 text-muted-foreground" />
            Filtros
          </button>
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
              <CreditCard className="h-5 w-5" />
            </div>
            <h3 className="text-lg font-medium text-foreground">
              Error al cargar
            </h3>
            <p className="text-muted-foreground mt-1">
              No se pudieron obtener las cuentas.
            </p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="min-w-full text-left text-sm">
                <thead className="bg-secondary/30 text-muted-foreground font-medium">
                  <tr>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Cuenta
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Cliente
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Tipo
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Saldo
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Estado
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Fecha Creación
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
                  {accounts.length === 0 ? (
                    <tr>
                      <td
                        colSpan={7}
                        className="px-6 py-16 text-center text-muted-foreground"
                      >
                        No se encontraron cuentas.
                      </td>
                    </tr>
                  ) : (
                    accounts.map((account) => (
                      <tr
                        key={account.id}
                        className="group hover:bg-secondary/40 transition-colors"
                      >
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-3">
                            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-secondary text-foreground border border-border">
                              <Wallet className="h-5 w-5" />
                            </div>
                            <div>
                              <div className="font-medium text-foreground font-mono">
                                {account.account_number}
                              </div>
                              <div className="text-xs text-muted-foreground">
                                {account.currency}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-2">
                            <div className="h-6 w-6 rounded-full bg-secondary flex items-center justify-center text-[10px] font-medium text-muted-foreground">
                              {account.client_name?.charAt(0) || "?"}
                            </div>
                            <span className="font-medium text-foreground">
                              {account.client_name || "Desconocido"}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <span className="inline-flex items-center rounded-md bg-secondary/50 px-2 py-1 text-xs font-medium text-muted-foreground ring-1 ring-inset ring-border/50 capitalize">
                            {account.type === "savings"
                              ? "Ahorro"
                              : "Corriente"}
                          </span>
                        </td>
                        <td className="px-6 py-4">
                          <div className="font-medium text-foreground">
                            $
                            {account.balance.toLocaleString("es-MX", {
                              minimumFractionDigits: 2,
                            })}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <span
                            className={cn(
                              "inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset",
                              account.status === "active"
                                ? "bg-secondary text-foreground ring-border"
                                : account.status === "frozen"
                                  ? "bg-secondary/50 text-muted-foreground ring-border/50"
                                  : "bg-destructive/10 text-destructive ring-destructive/20",
                            )}
                          >
                            {account.status === "active"
                              ? "Activa"
                              : account.status === "frozen"
                                ? "Congelada"
                                : "Cerrada"}
                          </span>
                        </td>
                        <td className="px-6 py-4 text-muted-foreground text-xs">
                          {new Date(account.created_at).toLocaleDateString()}
                        </td>
                        <td className="px-6 py-4 text-right">
                          <Link
                            href={`/accounts/${createSlug(account.id)}`}
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
            {accounts.length > 0 && (
              <div className="flex items-center justify-between border-t border-border px-6 py-4">
                <div className="text-xs text-muted-foreground">
                  Mostrando{" "}
                  <span className="font-medium text-foreground">1</span> a{" "}
                  <span className="font-medium text-foreground">
                    {accounts.length}
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
