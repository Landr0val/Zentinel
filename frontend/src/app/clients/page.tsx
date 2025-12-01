"use client";

import Link from "next/link";
import {
  Plus,
  Search,
  Filter,
  MoreHorizontal,
  Loader2,
  User,
  Mail,
  Phone,
} from "lucide-react";
import PageHeader from "../../components/PageHeader";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../lib/api";
import { useState } from "react";
import { cn } from "../../lib/utils";
import { createSlug } from "../../lib/slug-manager";

export default function ClientsPage() {
  const [searchTerm, setSearchTerm] = useState("");

  const { data, isLoading, error } = useQuery({
    queryKey: ["clients", searchTerm],
    queryFn: () =>
      api.clients.list(searchTerm ? { search: searchTerm } : undefined),
  });

  const clients = data?.data || [];

  return (
    <div className="animate-in fade-in duration-500">
      <PageHeader
        title="Clientes"
        description="Gestión de perfiles y riesgo."
        action={
          <Link
            href="/clients/new"
            className="inline-flex items-center justify-center rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow-zen hover:bg-primary/90 transition-colors"
          >
            <Plus className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
            Nuevo Cliente
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
            placeholder="Buscar por nombre, email o documento..."
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
              <User className="h-5 w-5" />
            </div>
            <h3 className="text-lg font-medium text-foreground">
              Error al cargar
            </h3>
            <p className="text-muted-foreground mt-1">
              No se pudieron obtener los clientes.
            </p>
          </div>
        ) : (
          <>
            <div className="overflow-x-auto">
              <table className="min-w-full text-left text-sm">
                <thead className="bg-secondary/30 text-muted-foreground font-medium">
                  <tr>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Cliente
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Contacto
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Estado
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Riesgo
                    </th>
                    <th scope="col" className="px-6 py-4 font-medium">
                      Fecha Registro
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
                  {clients.length === 0 ? (
                    <tr>
                      <td
                        colSpan={6}
                        className="px-6 py-16 text-center text-muted-foreground"
                      >
                        No se encontraron clientes.
                      </td>
                    </tr>
                  ) : (
                    clients.map((client) => (
                      <tr
                        key={client.id}
                        className="group hover:bg-secondary/40 transition-colors"
                      >
                        <td className="px-6 py-4">
                          <div className="flex items-center gap-3">
                            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-secondary text-foreground font-medium border border-border">
                              {client.full_name.charAt(0)}
                            </div>
                            <div>
                              <div className="font-medium text-foreground">
                                {client.full_name}
                              </div>
                              <div className="text-xs text-muted-foreground font-mono mt-0.5">
                                {client.document_type}: {client.document_number}
                              </div>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <div className="flex flex-col gap-1">
                            <div className="flex items-center gap-2 text-xs text-muted-foreground">
                              <Mail className="h-3 w-3" />
                              <span>{client.email}</span>
                            </div>
                            <div className="flex items-center gap-2 text-xs text-muted-foreground">
                              <Phone className="h-3 w-3" />
                              <span>{client.phone}</span>
                            </div>
                          </div>
                        </td>
                        <td className="px-6 py-4">
                          <span
                            className={cn(
                              "inline-flex items-center rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset",
                              (client.status || "active") === "active"
                                ? "bg-secondary text-foreground ring-border"
                                : client.status === "inactive"
                                  ? "bg-secondary/50 text-muted-foreground ring-border/50"
                                  : "bg-destructive/10 text-destructive ring-destructive/20",
                            )}
                          >
                            {(client.status || "active") === "active"
                              ? "Activo"
                              : client.status === "inactive"
                                ? "Inactivo"
                                : "Bloqueado"}
                          </span>
                        </td>
                        <td className="px-6 py-4">
                          <span
                            className={cn(
                              "inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset",
                              client.risk_profile === "low"
                                ? "bg-secondary/50 text-muted-foreground ring-border/50"
                                : client.risk_profile === "standard"
                                  ? "bg-secondary text-foreground ring-border"
                                  : "bg-destructive/10 text-destructive ring-destructive/20",
                            )}
                          >
                            {client.risk_profile === "low"
                              ? "Bajo"
                              : client.risk_profile === "standard"
                                ? "Estándar"
                                : "Alto"}
                          </span>
                        </td>
                        <td className="px-6 py-4 text-muted-foreground text-xs">
                          {new Date(client.created_at).toLocaleDateString()}
                        </td>
                        <td className="px-6 py-4 text-right">
                          <Link
                            href={`/clients/${createSlug(client.id)}`}
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
            {clients.length > 0 && (
              <div className="flex items-center justify-between border-t border-border px-6 py-4">
                <div className="text-xs text-muted-foreground">
                  Mostrando{" "}
                  <span className="font-medium text-foreground">1</span> a{" "}
                  <span className="font-medium text-foreground">
                    {clients.length}
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
