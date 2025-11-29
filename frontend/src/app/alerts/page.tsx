"use client";

import Link from "next/link";
import {
  Search,
  MoreHorizontal,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Clock,
  Loader2,
} from "lucide-react";
import PageHeader from "@/components/PageHeader";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useState } from "react";

export default function AlertsPage() {
  const [searchTerm, setSearchTerm] = useState("");

  const { data, isLoading, error } = useQuery({
    queryKey: ["alerts", searchTerm],
    queryFn: () =>
      api.alerts.list(searchTerm ? { search: searchTerm } : undefined),
  });

  const alerts = data?.data || [];

  return (
    <div>
      <PageHeader
        title="Alertas de Fraude"
        description="Gestión y resolución de alertas generadas por el sistema."
      />

      {/* Filters */}
      <div className="mb-6 flex flex-col gap-4 rounded-lg bg-white p-4 shadow-sm border border-slate-100 sm:flex-row sm:items-center sm:justify-between">
        <div className="relative flex-1 max-w-md">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
            <Search className="h-5 w-5 text-slate-400" aria-hidden="true" />
          </div>
          <input
            type="text"
            className="block w-full rounded-md border-0 py-1.5 pl-10 text-slate-900 ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6"
            placeholder="Buscar por ID de alerta o transacción..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <div className="flex gap-2">
          <select className="block rounded-md border-0 py-1.5 pl-3 pr-10 text-slate-900 ring-1 ring-inset ring-slate-300 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6">
            <option>Todos los estados</option>
            <option>Nuevas</option>
            <option>En Investigación</option>
            <option>Resueltas</option>
            <option>Falsos Positivos</option>
          </select>
          <select className="block rounded-md border-0 py-1.5 pl-3 pr-10 text-slate-900 ring-1 ring-inset ring-slate-300 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6">
            <option>Cualquier severidad</option>
            <option>Alta</option>
            <option>Media</option>
            <option>Baja</option>
          </select>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-hidden rounded-xl bg-white shadow-sm border border-slate-100">
        {isLoading ? (
          <div className="flex h-64 items-center justify-center">
            <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
          </div>
        ) : error ? (
          <div className="p-8 text-center text-red-500">
            Error al cargar las alertas. Por favor intente nuevamente.
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="min-w-full divide-y divide-slate-200">
                <thead className="bg-slate-50">
                  <tr>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500"
                    >
                      Alerta
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500"
                    >
                      Severidad
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500"
                    >
                      Descripción
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500"
                    >
                      Estado
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500"
                    >
                      Fecha
                    </th>
                    <th scope="col" className="relative px-6 py-3">
                      <span className="sr-only">Acciones</span>
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-200 bg-white">
                  {alerts.length === 0 ? (
                    <tr>
                      <td
                        colSpan={6}
                        className="px-6 py-12 text-center text-slate-500"
                      >
                        No se encontraron alertas.
                      </td>
                    </tr>
                  ) : (
                    alerts.map((alert) => (
                      <tr
                        key={alert.id}
                        className="hover:bg-slate-50 transition-colors"
                      >
                        <td className="whitespace-nowrap px-6 py-4">
                          <div className="flex items-center">
                            <div className="shrink-0">
                              <AlertTriangle
                                className={`h-5 w-5 ${
                                  alert.severity === "high"
                                    ? "text-red-500"
                                    : alert.severity === "medium"
                                      ? "text-yellow-500"
                                      : "text-blue-500"
                                }`}
                              />
                            </div>
                            <div className="ml-4">
                              <div className="text-sm font-medium text-slate-900">
                                {alert.id}
                              </div>
                              <div className="text-xs text-slate-500">
                                Tx: {alert.transaction_id}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="whitespace-nowrap px-6 py-4">
                          <span
                            className={`inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset ${
                              alert.severity === "high"
                                ? "bg-red-50 text-red-700 ring-red-600/20"
                                : alert.severity === "medium"
                                  ? "bg-yellow-50 text-yellow-800 ring-yellow-600/20"
                                  : "bg-blue-50 text-blue-700 ring-blue-600/20"
                            }`}
                          >
                            {alert.severity === "high"
                              ? "Alta"
                              : alert.severity === "medium"
                                ? "Media"
                                : "Baja"}
                          </span>
                        </td>
                        <td className="px-6 py-4">
                          <div
                            className="text-sm text-slate-900 max-w-xs truncate"
                            title={alert.description}
                          >
                            {alert.description}
                          </div>
                        </td>
                        <td className="whitespace-nowrap px-6 py-4">
                          <span
                            className={`inline-flex items-center gap-1 rounded-full px-2 py-1 text-xs font-medium ${
                              alert.status === "resolved"
                                ? "bg-green-50 text-green-700 ring-1 ring-inset ring-green-600/20"
                                : alert.status === "new"
                                  ? "bg-blue-50 text-blue-700 ring-1 ring-inset ring-blue-600/20"
                                  : alert.status === "investigating"
                                    ? "bg-yellow-50 text-yellow-800 ring-1 ring-inset ring-yellow-600/20"
                                    : "bg-slate-50 text-slate-700 ring-1 ring-inset ring-slate-600/20"
                            }`}
                          >
                            {alert.status === "resolved" && (
                              <CheckCircle className="h-3 w-3" />
                            )}
                            {alert.status === "new" && (
                              <AlertTriangle className="h-3 w-3" />
                            )}
                            {alert.status === "investigating" && (
                              <Clock className="h-3 w-3" />
                            )}
                            {alert.status === "false_positive" && (
                              <XCircle className="h-3 w-3" />
                            )}
                            <span className="capitalize">
                              {alert.status === "resolved"
                                ? "Resuelta"
                                : alert.status === "new"
                                  ? "Nueva"
                                  : alert.status === "investigating"
                                    ? "Investigando"
                                    : "Falso Positivo"}
                            </span>
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-6 py-4 text-sm text-slate-500">
                          {new Date(alert.created_at).toLocaleString()}
                        </td>
                        <td className="whitespace-nowrap px-6 py-4 text-right text-sm font-medium">
                          <Link
                            href={`/alerts/${alert.id}`}
                            className="text-slate-400 hover:text-blue-600"
                          >
                            <MoreHorizontal className="h-5 w-5" />
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
              <div className="flex items-center justify-between border-t border-slate-200 bg-white px-4 py-3 sm:px-6">
                <div className="flex flex-1 justify-between sm:hidden">
                  <button className="relative inline-flex items-center rounded-md border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50">
                    Anterior
                  </button>
                  <button className="relative ml-3 inline-flex items-center rounded-md border border-slate-300 bg-white px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50">
                    Siguiente
                  </button>
                </div>
                <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                  <div>
                    <p className="text-sm text-slate-700">
                      Mostrando <span className="font-medium">1</span> a{" "}
                      <span className="font-medium">{alerts.length}</span>{" "}
                      resultados
                    </p>
                  </div>
                  <div>
                    <nav
                      className="isolate inline-flex -space-x-px rounded-md shadow-sm"
                      aria-label="Pagination"
                    >
                      <button className="relative inline-flex items-center rounded-l-md px-2 py-2 text-slate-400 ring-1 ring-inset ring-slate-300 hover:bg-slate-50 focus:z-20 focus:outline-offset-0">
                        <span className="sr-only">Anterior</span>
                        <svg
                          className="h-5 w-5"
                          viewBox="0 0 20 20"
                          fill="currentColor"
                          aria-hidden="true"
                        >
                          <path
                            fillRule="evenodd"
                            d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z"
                            clipRule="evenodd"
                          />
                        </svg>
                      </button>
                      <button
                        aria-current="page"
                        className="relative z-10 inline-flex items-center bg-blue-600 px-4 py-2 text-sm font-semibold text-white focus:z-20 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
                      >
                        1
                      </button>
                      <button className="relative inline-flex items-center rounded-r-md px-2 py-2 text-slate-400 ring-1 ring-inset ring-slate-300 hover:bg-slate-50 focus:z-20 focus:outline-offset-0">
                        <span className="sr-only">Siguiente</span>
                        <svg
                          className="h-5 w-5"
                          viewBox="0 0 20 20"
                          fill="currentColor"
                          aria-hidden="true"
                        >
                          <path
                            fillRule="evenodd"
                            d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z"
                            clipRule="evenodd"
                          />
                        </svg>
                      </button>
                    </nav>
                  </div>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}
