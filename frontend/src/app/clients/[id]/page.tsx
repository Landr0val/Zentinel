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
import { api } from "@/lib/api";

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
        <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
      </div>
    );
  }

  if (errorClient || !client) {
    return (
      <div className="rounded-lg bg-red-50 p-4 text-red-800">
        <h3 className="text-lg font-medium">Error al cargar cliente</h3>
        <p>No se pudo encontrar la información del cliente.</p>
        <Link
          href="/clients"
          className="mt-4 inline-block text-sm font-medium underline"
        >
          Volver a la lista
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-4 mb-6">
        <Link
          href="/clients"
          className="p-2 rounded-full hover:bg-slate-100 text-slate-500 transition-colors"
        >
          <ArrowLeft className="h-5 w-5" />
        </Link>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-slate-900">
            {client.full_name}
          </h1>
          <p className="text-slate-500">ID: {id}</p>
        </div>
        <Link
          href={`/clients/${id}/edit`}
          className="inline-flex items-center justify-center rounded-md bg-white px-4 py-2 text-sm font-medium text-slate-700 shadow-sm ring-1 ring-inset ring-slate-300 hover:bg-slate-50"
        >
          <Edit className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
          Editar
        </Link>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Client Info Card */}
        <div className="lg:col-span-1 space-y-6">
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <h3 className="text-lg font-semibold text-slate-900 mb-4">
              Información Personal
            </h3>
            <div className="space-y-4">
              <div className="flex items-start gap-3">
                <Mail className="h-5 w-5 text-slate-400 mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-slate-900">Email</p>
                  <p className="text-sm text-slate-500">{client.email}</p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <Phone className="h-5 w-5 text-slate-400 mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-slate-900">Teléfono</p>
                  <p className="text-sm text-slate-500">{client.phone}</p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <MapPin className="h-5 w-5 text-slate-400 mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-slate-900">
                    Dirección
                  </p>
                  <p className="text-sm text-slate-500">
                    {client.address}, {client.city}, {client.state}{" "}
                    {client.zip_code}
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <Calendar className="h-5 w-5 text-slate-400 mt-0.5" />
                <div>
                  <p className="text-sm font-medium text-slate-900">
                    Fecha de Registro
                  </p>
                  <p className="text-sm text-slate-500">
                    {new Date(client.created_at).toLocaleDateString()}
                  </p>
                </div>
              </div>
            </div>

            <div className="mt-6 pt-6 border-t border-slate-100">
              <h4 className="text-sm font-medium text-slate-900 mb-3">
                Estado del Cliente
              </h4>
              <div className="flex flex-col gap-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm text-slate-500">Estado</span>
                  <span
                    className={`inline-flex rounded-full px-2 py-1 text-xs font-semibold capitalize ${
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
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm text-slate-500">
                    Nivel de Riesgo
                  </span>
                  <span
                    className={`inline-flex items-center gap-1 rounded-md px-2 py-1 text-xs font-medium ring-1 ring-inset ${
                      client.risk_level === "low"
                        ? "bg-green-50 text-green-700 ring-green-600/20"
                        : client.risk_level === "medium"
                          ? "bg-yellow-50 text-yellow-800 ring-yellow-600/20"
                          : "bg-red-50 text-red-700 ring-red-600/20"
                    }`}
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
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-slate-900">
                Cuentas Asociadas
              </h3>
              <Link
                href="/accounts"
                className="text-sm font-medium text-blue-600 hover:text-blue-500"
              >
                Ver todas
              </Link>
            </div>

            {isLoadingAccounts ? (
              <div className="flex justify-center py-8">
                <Loader2 className="h-6 w-6 animate-spin text-blue-600" />
              </div>
            ) : accounts.length === 0 ? (
              <div className="text-center py-8 text-slate-500">
                No hay cuentas asociadas a este cliente.
              </div>
            ) : (
              <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                {accounts.map((account) => (
                  <div
                    key={account.id}
                    className="relative flex items-center space-x-3 rounded-lg border border-slate-300 bg-white px-6 py-5 shadow-sm focus-within:ring-2 focus-within:ring-blue-500 focus-within:ring-offset-2 hover:border-slate-400"
                  >
                    <div className="shrink-0">
                      <div className="h-10 w-10 rounded-full bg-blue-50 flex items-center justify-center text-blue-600">
                        <CreditCard className="h-6 w-6" />
                      </div>
                    </div>
                    <div className="min-w-0 flex-1">
                      <Link
                        href={`/accounts/${account.id}`}
                        className="focus:outline-none"
                      >
                        <span className="absolute inset-0" aria-hidden="true" />
                        <p className="text-sm font-medium text-slate-900 capitalize">
                          {account.type === "savings" ? "Ahorro" : "Corriente"}
                        </p>
                        <p className="truncate text-sm text-slate-500">
                          {account.account_number}
                        </p>
                      </Link>
                    </div>
                    <div className="flex flex-col items-end">
                      <span className="text-sm font-bold text-slate-900">
                        $
                        {account.balance.toLocaleString("es-MX", {
                          minimumFractionDigits: 2,
                        })}
                      </span>
                      <span className="text-xs text-slate-500">
                        {account.currency}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>

          {/* Recent Activity Placeholder */}
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <h3 className="text-lg font-semibold text-slate-900 mb-4">
              Actividad Reciente
            </h3>
            <div className="text-center py-8 text-slate-500">
              <Activity className="h-12 w-12 mx-auto text-slate-300 mb-3" />
              <p>No hay actividad reciente para mostrar.</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
