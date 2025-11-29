"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import {
  ArrowLeft,
  Edit,
  CreditCard,
  Activity,
  Phone,
  Mail,
  Calendar,
  Shield,
  MapPin,
  Loader2,
} from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import { cn } from "../../../lib/utils";

export default function ClientDetailPage() {
  const params = useParams();
  const id = params?.id as string;

  const {
    data: clientData,
    isLoading: isLoadingClient,
    error: errorClient,
  } = useQuery({
    queryKey: ["client", id],
    queryFn: () => api.clients.get(id),
    enabled: !!id,
  });

  const { data: accountsData, isLoading: isLoadingAccounts } = useQuery({
    queryKey: ["clientAccounts", id],
    queryFn: () => api.clients.getAccounts(id),
    enabled: !!id,
  });

  const client = clientData?.data;
  const accounts = accountsData?.data || [];

  if (isLoadingClient) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (errorClient || !client) {
    return (
      <div className="rounded-xl border border-destructive/20 bg-destructive/5 p-6 text-destructive">
        <h3 className="text-lg font-medium">Error al cargar cliente</h3>
        <p>No se pudo encontrar la información del cliente.</p>
        <Link
          href="/clients"
          className="mt-4 inline-block text-sm font-medium underline hover:text-destructive/80"
        >
          Volver a la lista
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6 animate-in fade-in duration-500">
      <div className="flex items-center gap-4 mb-6">
        <Link
          href="/clients"
          className="p-2 rounded-full hover:bg-secondary text-muted-foreground transition-colors"
        >
          <ArrowLeft className="h-5 w-5" />
        </Link>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-foreground">
            {client.full_name}
          </h1>
          <p className="text-muted-foreground">ID: {id}</p>
        </div>
        <Link
          href={`/clients/${id}/edit`}
          className="inline-flex items-center justify-center rounded-xl border border-border bg-card px-4 py-2 text-sm font-medium text-foreground shadow-sm hover:bg-secondary transition-colors"
        >
          <Edit className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
          Editar
        </Link>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Client Info Card */}
        <div className="lg:col-span-1 space-y-6">
          <div className="rounded-xl bg-card p-6 shadow-zen border border-border">
            <h3 className="text-lg font-semibold text-foreground mb-4">
              Información Personal
            </h3>
            <div className="space-y-4">
              <div className="flex items-start gap-3">
                <Mail className="h-5 w-5 text-muted-foreground mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-foreground">Email</p>
                  <p className="text-sm text-muted-foreground">
                    {client.email}
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <Phone className="h-5 w-5 text-muted-foreground mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-foreground">
                    Teléfono
                  </p>
                  <p className="text-sm text-muted-foreground">
                    {client.phone}
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <MapPin className="h-5 w-5 text-muted-foreground mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-foreground">
                    Dirección
                  </p>
                  <p className="text-sm text-muted-foreground">
                    {client.address}, {client.city}, {client.state}{" "}
                    {client.zip_code}
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <Calendar className="h-5 w-5 text-muted-foreground mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-foreground">
                    Fecha de Registro
                  </p>
                  <p className="text-sm text-muted-foreground">
                    {new Date(client.created_at).toLocaleDateString()}
                  </p>
                </div>
              </div>
            </div>

            <div className="mt-6 pt-6 border-t border-border">
              <h4 className="text-sm font-medium text-foreground mb-3">
                Estado del Cliente
              </h4>
              <div className="flex flex-col gap-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">Estado</span>
                  <span
                    className={cn(
                      "inline-flex rounded-full px-2 py-1 text-xs font-semibold capitalize ring-1 ring-inset",
                      client.status === "active"
                        ? "bg-secondary text-foreground ring-border"
                        : client.status === "inactive"
                          ? "bg-secondary/50 text-muted-foreground ring-border/50"
                          : "bg-destructive/10 text-destructive ring-destructive/20",
                    )}
                  >
                    {client.status === "active"
                      ? "Activo"
                      : client.status === "inactive"
                        ? "Inactivo"
                        : "Bloqueado"}
                  </span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-muted-foreground">
                    Nivel de Riesgo
                  </span>
                  <span
                    className={cn(
                      "inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset",
                      client.risk_level === "low"
                        ? "bg-secondary/50 text-muted-foreground ring-border/50"
                        : client.risk_level === "medium"
                          ? "bg-secondary text-foreground ring-border"
                          : "bg-destructive/10 text-destructive ring-destructive/20",
                    )}
                  >
                    <Shield className="h-3 w-3" />
                    {client.risk_level === "low"
                      ? "Bajo"
                      : client.risk_level === "medium"
                        ? "Medio"
                        : "Alto"}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Main Content Area */}
        <div className="lg:col-span-2 space-y-6">
          {/* Accounts Section */}
          <div className="rounded-xl bg-card p-6 shadow-zen border border-border">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-foreground">
                Cuentas Asociadas
              </h3>
              <Link
                href="/accounts"
                className="text-sm font-medium text-primary hover:text-primary/80"
              >
                Ver todas
              </Link>
            </div>

            {isLoadingAccounts ? (
              <div className="flex justify-center py-8">
                <Loader2 className="h-6 w-6 animate-spin text-primary" />
              </div>
            ) : accounts.length === 0 ? (
              <div className="text-center py-8 text-muted-foreground">
                No hay cuentas asociadas a este cliente.
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                {accounts.map((account) => (
                  <div
                    key={account.id}
                    className="relative flex items-center space-x-3 rounded-xl border border-border bg-card px-6 py-5 shadow-sm transition-colors hover:bg-secondary/30"
                  >
                    <div className="shrink-0">
                      <div className="h-10 w-10 rounded-full bg-secondary flex items-center justify-center text-foreground">
                        <CreditCard className="h-5 w-5" />
                      </div>
                    </div>
                    <div className="min-w-0 flex-1">
                      <Link
                        href={`/accounts/${account.id}`}
                        className="focus:outline-none"
                      >
                        <span className="absolute inset-0" aria-hidden="true" />
                        <p className="text-sm font-medium text-foreground capitalize">
                          {account.type === "savings" ? "Ahorro" : "Corriente"}
                        </p>
                        <p className="truncate text-sm text-muted-foreground font-mono">
                          {account.account_number}
                        </p>
                      </Link>
                    </div>
                    <div className="flex flex-col items-end">
                      <span className="text-sm font-bold text-foreground">
                        $
                        {account.balance.toLocaleString("es-MX", {
                          minimumFractionDigits: 2,
                        })}
                      </span>
                      <span className="text-xs text-muted-foreground">
                        {account.currency}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Recent Activity Placeholder */}
          <div className="rounded-xl bg-card p-6 shadow-zen border border-border">
            <h3 className="text-lg font-semibold text-foreground mb-4">
              Actividad Reciente
            </h3>
            <div className="text-center py-8 text-muted-foreground">
              <Activity className="h-12 w-12 mx-auto text-muted-foreground/30 mb-3" />
              <p>No hay actividad reciente para mostrar.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
