"use client";

import Link from "next/link";
import { Plus, Search, Filter, MoreHorizontal, Loader2 } from "lucide-react";
import PageHeader from "@/components/PageHeader";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { useState } from "react";

export default function ClientsPage() {
  const [searchTerm, setSearchTerm] = useState("");

  const { data, isLoading, error } = useQuery({
    queryKey: ["clients", searchTerm],
    queryFn: () =>
      api.clients.list(searchTerm ? { search: searchTerm } : undefined),
  });

  const clients = data?.data || [];

  return (
    <div>
      <PageHeader
        title="Clientes"
        description="Gestión de clientes y perfiles de riesgo."
        action={
          <Link
            href="/clients/new"
            className="inline-flex items-center justify-center rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600"
          >
            <Plus className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
            Nuevo Cliente
          </Link>
        }
      />

      {/* Filters */}
      <div className="mb-6 flex items-center justify-between gap-4 rounded-lg bg-white p-4 shadow-sm border border-slate-100">
        <div className="relative flex-1 max-w-md">
          <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
            <Search className="h-5 w-5 text-slate-400" aria-hidden="true" />
          </div>
          <input
            type="text"
            className="block w-full rounded-md border-0 py-1.5 pl-10 text-slate-900 ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6"
            placeholder="Buscar por nombre, email o teléfono..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <button className="inline-flex items-center gap-2 rounded-md bg-white px-3 py-2 text-sm font-semibold text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 hover:bg-slate-50">
          <Filter className="h-4 w-4 text-slate-500" />
          Filtros
        </button>
      </div>

      {/* Table */}
      <div className="overflow-hidden rounded-xl bg-white shadow-sm border border-slate-100">
        {isLoading ? (
          <div className="flex h-64 items-center justify-center">
            <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
          </div>
        ) : error ? (
          <div className="p-8 text-center text-red-500">
            Error al cargar los clientes. Por favor intente nuevamente.
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
                      Cliente
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500"
                    >
                      Contacto
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
                      Riesgo
                    </th>
                    <th
                      scope="col"
                      className="px-6 py-3 text-left text-xs font-medium uppercase tracking-wider text-slate-500"
                    >
                      Fecha Registro
                    </th>
                    <th scope="col" className="relative px-6 py-3">
                      <span className="sr-only">Acciones</span>
                    </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-slate-200 bg-white">
                  {clients.length === 0 ? (
                    <tr>
                      <td
                        colSpan={6}
                        className="px-6 py-12 text-center text-slate-500"
                      >
                        No se encontraron clientes.
                      </td>
                    </tr>
                  ) : (
                    clients.map((client) => (
                      <tr
                        key={client.id}
                        className="hover:bg-slate-50 transition-colors"
                      >
                        <td className="whitespace-nowrap px-6 py-4">
                          <div className="flex items-center">
                            <div className="h-10 w-10 shrink-0 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 font-bold">
                              {client.full_name.charAt(0)}
                            </div>
                            <div className="ml-4">
                              <div className="text-sm font-medium text-slate-900">
                                {client.full_name}
                              </div>
                              <div className="text-sm text-slate-500">
                                ID: {client.id}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="whitespace-nowrap px-6 py-4">
                          <div className="text-sm text-slate-900">
                            {client.email}
                          </div>
                          <div className="text-sm text-slate-500">
                            {client.phone}
                          </div>
                        </td>
                        <td className="whitespace-nowrap px-6 py-4">
                          <span
                            className={`inline-flex rounded-full px-2 text-xs font-semibold leading-5 ${
                              client.status === "active"
                                ? "bg-green-100 text-green-800"
                                : client.status === "inactive"
                                  ? "bg-slate-100 text-slate-800"
                                  : "bg-red-100 text-red-800"
                            }`}
                          >
                            {client.status === "active"
                              ? "Activo"
                              : client.status === "inactive"
                                ? "Inactivo"
                                : "Bloqueado"}
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-6 py-4">
                          <span
                            className={`inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset ${
                              client.risk_level === "low"
                                ? "bg-green-50 text-green-700 ring-green-600/20"
                                : client.risk_level === "medium"
                                  ? "bg-yellow-50 text-yellow-800 ring-yellow-600/20"
                                  : "bg-red-50 text-red-700 ring-red-600/20"
                            }`}
                          >
                            {client.risk_level === "low"
                              ? "Bajo"
                              : client.risk_level === "medium"
                                ? "Medio"
                                : "Alto"}
                          </span>
                        </td>
                        <td className="whitespace-nowrap px-6 py-4 text-sm text-slate-500">
                          {new Date(client.created_at).toLocaleDateString()}
                        </td>
                        <td className="whitespace-nowrap px-6 py-4 text-right text-sm font-medium">
                          <Link
                            href={`/clients/${client.id}`}
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

            {/* Pagination (Placeholder - API needs to support pagination metadata return) */}
            {clients.length > 0 && (
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
                      <span className="font-medium">{clients.length}</span>{" "}
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
