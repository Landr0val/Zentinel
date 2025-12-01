"use client";

import Link from "next/link";
import {
  Search,
  MoreHorizontal,
  AlertCircle,
  Check,
  X,
  Clock,
  Loader2,
  Filter,
  ShieldAlert,
} from "lucide-react";
import PageHeader from "../../components/PageHeader";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../lib/api";
import { useState } from "react";
import { cn } from "../../lib/utils";
import { createSlug } from "../../lib/slug-manager";

export default function AlertsPage() {
  const [searchTerm, setSearchTerm] = useState("");
  const [status, setStatus] = useState("");
  const [severity, setSeverity] = useState("");

  const { data, isLoading, error } = useQuery({
    queryKey: ["alerts", searchTerm, status, severity],
    queryFn: () =>
      api.alerts.list({
        search: searchTerm || undefined,
        status: status || undefined,
        severity: severity || undefined,
      }),
  });

  const alerts = data?.data || [];

  return (
    <div className="animate-in fade-in duration-500">
      <PageHeader
        title="Alertas de Fraude"
        description="Gestión y resolución de incidentes."
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
            placeholder="Buscar por ID de alerta o transacción..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <div className="flex gap-3">
          <div className="relative">
            <select
              value={status}
              onChange={(e) => setStatus(e.target.value)}
              className="appearance-none block w-full rounded-xl border border-border bg-card py-2.5 pl-4 pr-10 text-sm text-foreground focus:border-foreground focus:ring-0 transition-all shadow-sm hover:border-foreground/50 cursor-pointer outline-none"
            >
              <option value="">Todos los estados</option>
              <option value="pending">Pendientes</option>
              <option value="reviewing">En Revisión</option>
              <option value="resolved">Resueltas</option>
              <option value="false_positive">Falsos Positivos</option>
            </select>
            <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-3 text-muted-foreground">
              <Filter className="h-3 w-3" />
            </div>
          </div>
          <div className="relative">
            <select
              value={severity}
              onChange={(e) => setSeverity(e.target.value)}
              className="appearance-none block w-full rounded-xl border border-border bg-card py-2.5 pl-4 pr-10 text-sm text-foreground focus:border-foreground focus:ring-0 transition-all shadow-sm hover:border-foreground/50 cursor-pointer outline-none"
            >
              <option value="">Cualquier severidad</option>
              <option value="critical">Crítica</option>
              <option value="high">Alta</option>
              <option value="medium">Media</option>
              <option value="low">Baja</option>
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
              No se pudieron obtener las alertas.
            </p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="min-w-full text-left text-sm">
                <thead className="bg-secondary/30 text-muted-foreground font-medium">
                  <tr>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Alerta
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Severidad
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Descripción
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Estado
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
                  {alerts.length === 0 ? (
                    <tr>
                      <td
                        colSpan={6}
                        className="px-6 py-16 text-center text-muted-foreground"
                      >
                        No se encontraron alertas.
                      </td>
                    </tr>
                  ) : (
                    alerts.map((alert) => (
                      <tr
                        key={alert.id}
                        className="group hover:bg-secondary/40 transition-colors"
                      >
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-3">
                            <div
                              className={cn(
                                "flex h-8 w-8 items-center justify-center rounded-full border transition-colors",
                                alert.severity === "critical" ||
                                  alert.severity === "high"
                                  ? "border-destructive/20 bg-destructive/10 text-destructive"
                                  : alert.severity === "medium"
                                    ? "border-border bg-secondary text-foreground"
                                    : "border-border bg-card text-muted-foreground",
                              )}
                            >
                              <ShieldAlert className="h-4 w-4" />
                            </div>
                            <div>
                              <div className="font-medium text-foreground">
                                {alert.code}
                              </div>
                              <div className="text-xs text-muted-foreground font-mono">
                                Tx: {alert.transaction_code}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <span
                            className={cn(
                              "inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset",
                              alert.severity === "critical" ||
                                alert.severity === "high"
                                ? "bg-destructive/10 text-destructive ring-destructive/20"
                                : alert.severity === "medium"
                                  ? "bg-secondary text-foreground ring-border"
                                  : "bg-secondary/50 text-muted-foreground ring-border/50",
                            )}
                          >
                            {alert.severity === "critical"
                              ? "Crítica"
                              : alert.severity === "high"
                                ? "Alta"
                                : alert.severity === "medium"
                                  ? "Media"
                                  : "Baja"}
                          </span>
                        </td>
                        <td className="px-6 py-4">
                          <div
                            className="text-foreground max-w-xs truncate"
                            title={alert.description}
                          >
                            {alert.description}
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-2">
                            {alert.status === "resolved" && (
                              <Check className="h-3 w-3 text-foreground" />
                            )}
                            {alert.status === "pending" && (
                              <AlertCircle className="h-3 w-3 text-destructive" />
                            )}
                            {alert.status === "reviewing" && (
                              <Clock className="h-3 w-3 text-muted-foreground" />
                            )}
                            {alert.status === "false_positive" && (
                              <X className="h-3 w-3 text-muted-foreground" />
                            )}
                            <span
                              className={cn(
                                "text-xs font-medium",
                                alert.status === "pending"
                                  ? "text-destructive"
                                  : alert.status === "resolved"
                                    ? "text-foreground"
                                    : "text-muted-foreground",
                              )}
                            >
                              {alert.status === "resolved"
                                ? "Resuelta"
                                : alert.status === "pending"
                                  ? "Pendiente"
                                  : alert.status === "reviewing"
                                    ? "En Revisión"
                                    : "Falso Positivo"}
                            </span>
                          </div>
                        </td>
                        <td className="px-6 py-4 text-muted-foreground text-xs">
                          {new Date(alert.created_at).toLocaleString()}
                        </td>
                        <td className="px-6 py-4 text-right">
                          <Link
                            href={`/alerts/${createSlug(alert.id)}`}
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
            {alerts.length > 0 && (
              <div className="flex items-center justify-between border-t border-border px-6 py-4">
                <div className="text-xs text-muted-foreground">
                  Mostrando{" "}
                  <span className="font-medium text-foreground">1</span> a{" "}
                  <span className="font-medium text-foreground">
                    {alerts.length}
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
