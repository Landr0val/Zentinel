"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import {
  ArrowLeft,
  AlertTriangle,
  CheckCircle,
  XCircle,
  Clock,
  Shield,
  User,
  CreditCard,
  Send,
  MoreHorizontal,
  Loader2,
  AlertCircle,
} from "lucide-react";
import { cn } from "../../../lib/utils";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../../lib/api";

export default function AlertDetailPage() {
  const params = useParams();
  const id = params?.id as string;

  const {
    data: alertRes,
    isLoading: isLoadingAlert,
    error,
  } = useQuery({
    queryKey: ["alert", id],
    queryFn: () => api.alerts.get(id),
  });

  const alert = alertRes?.data;

  const { data: txnRes } = useQuery({
    queryKey: ["transaction", alert?.transaction_id],
    queryFn: () => api.transactions.get(alert!.transaction_id as string),
    enabled: !!alert?.transaction_id,
  });
  const transaction = txnRes?.data;

  const { data: clientRes } = useQuery({
    queryKey: ["client", alert?.client_id],
    queryFn: () => api.clients.get(alert!.client_id),
    enabled: !!alert?.client_id,
  });
  const client = clientRes?.data;

  if (isLoadingAlert) {
    return (
      <div className="flex h-96 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error || !alert) {
    return (
      <div className="flex h-96 flex-col items-center justify-center text-center">
        <AlertCircle className="h-12 w-12 text-destructive mb-4" />
        <h2 className="text-lg font-semibold">Error al cargar la alerta</h2>
        <p className="text-muted-foreground">
          No se pudo encontrar la información solicitada.
        </p>
      </div>
    );
  }

  const notes: {
    id: string;
    author: string;
    text: string;
    date: string;
    type: string;
  }[] = [];
  if (alert.ai_explanation) {
    notes.push({
      id: "ai-expl",
      author: "Sistema Sentinel (IA)",
      text: alert.ai_explanation,
      date: new Date(alert.created_at).toLocaleString(),
      type: "system",
    });
  }

  return (
    <div className="animate-in fade-in duration-500">
      {/* Header */}
      <div className="flex items-center gap-4 mb-8">
        <Link
          href="/alerts"
          className="p-2 rounded-full hover:bg-secondary text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeft className="h-5 w-5" />
        </Link>
        <div className="flex-1">
          <h1 className="text-3xl font-light tracking-tight text-foreground">
            Alerta {alert.code}
          </h1>
          <p className="mt-1 text-muted-foreground">Gestión de Incidencia</p>
        </div>
        <div className="flex gap-2">
          <span
            className={cn(
              "inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-sm font-medium border",
              alert.status === "resolved"
                ? "bg-green-500/10 text-green-600 border-green-500/20"
                : alert.status === "pending"
                  ? "bg-blue-500/10 text-blue-600 border-blue-500/20"
                  : alert.status === "reviewing"
                    ? "bg-yellow-500/10 text-yellow-600 border-yellow-500/20"
                    : "bg-secondary text-muted-foreground border-border",
            )}
          >
            {alert.status === "resolved" && (
              <CheckCircle className="h-3.5 w-3.5" />
            )}
            {alert.status === "pending" && (
              <AlertTriangle className="h-3.5 w-3.5" />
            )}
            {alert.status === "reviewing" && <Clock className="h-3.5 w-3.5" />}
            {alert.status === "false_positive" && (
              <XCircle className="h-3.5 w-3.5" />
            )}
            <span className="capitalize">
              {alert.status === "resolved"
                ? "Resuelta"
                : alert.status === "pending"
                  ? "Pendiente"
                  : alert.status === "reviewing"
                    ? "En Revisión"
                    : "Falso Positivo"}
            </span>
          </span>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Alert Details */}
          <div className="rounded-3xl border border-border bg-card p-6 shadow-zen">
            <div className="flex items-start gap-4">
              <div
                className={cn(
                  "p-3 rounded-2xl",
                  alert.severity === "critical" || alert.severity === "high"
                    ? "bg-destructive/10 text-destructive"
                    : alert.severity === "medium"
                      ? "bg-yellow-500/10 text-yellow-600"
                      : "bg-blue-500/10 text-blue-600",
                )}
              >
                <AlertTriangle className="h-6 w-6" />
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-medium text-foreground">
                  {alert.description}
                </h3>
                <p className="text-muted-foreground mt-2 text-sm">
                  Tipo de Alerta:{" "}
                  <span className="font-mono text-xs bg-secondary px-1.5 py-0.5 rounded text-foreground ml-1 uppercase">
                    {alert.alert_type}
                  </span>
                </p>
              </div>
            </div>
          </div>

          {/* Timeline / Notes */}
          <div className="rounded-3xl border border-border bg-card p-6 shadow-zen">
            <div className="flex items-center justify-between mb-6">
              <h3 className="text-lg font-medium text-foreground">
                Historial y Notas
              </h3>
              <button className="text-muted-foreground hover:text-foreground transition-colors">
                <MoreHorizontal className="h-5 w-5" />
              </button>
            </div>

            <div className="space-y-8 relative before:absolute before:inset-y-0 before:left-[19px] before:w-px before:bg-border">
              {notes.map((note) => (
                <div key={note.id} className="relative flex gap-4">
                  <div className="flex flex-col items-center shrink-0 z-10">
                    <div className="h-10 w-10 rounded-full bg-card border border-border flex items-center justify-center shadow-sm">
                      {note.type === "system" ? (
                        <Shield className="h-4 w-4 text-muted-foreground" />
                      ) : (
                        <User className="h-4 w-4 text-muted-foreground" />
                      )}
                    </div>
                  </div>
                  <div className="flex-1 pt-1">
                    <div className="flex items-center justify-between mb-2">
                      <span className="text-sm font-medium text-foreground">
                        {note.author}
                      </span>
                      <span className="text-xs text-muted-foreground">
                        {note.date}
                      </span>
                    </div>
                    <div className="text-sm text-muted-foreground bg-secondary/30 p-4 rounded-2xl border border-border/50">
                      {note.text}
                    </div>
                  </div>
                </div>
              ))}
            </div>

            {/* Add Note Input */}
            <div className="flex gap-4 pt-8 mt-2 relative z-10">
              <div className="h-10 w-10 rounded-full bg-foreground flex items-center justify-center shrink-0 shadow-sm">
                <User className="h-4 w-4 text-background" />
              </div>
              <div className="flex-1">
                <div className="relative">
                  <textarea
                    rows={3}
                    className="block w-full rounded-2xl border border-border bg-secondary/30 py-3 px-4 text-sm text-foreground placeholder:text-muted-foreground focus:border-foreground focus:ring-0 transition-all shadow-sm hover:border-foreground/50 outline-none resize-none"
                    placeholder="Agregar una nota de investigación..."
                  />
                  <div className="absolute bottom-3 right-3">
                    <button className="inline-flex items-center rounded-xl bg-foreground px-3 py-1.5 text-xs font-medium text-background shadow-sm hover:bg-foreground/90 transition-colors">
                      <Send className="h-3 w-3 mr-1.5" />
                      Agregar
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Sidebar Info */}
        <div className="space-y-6">
          {/* Status Actions */}
          <div className="rounded-3xl border border-border bg-card p-6 shadow-zen">
            <h3 className="text-sm font-medium text-foreground mb-4">
              Gestionar Estado
            </h3>
            <div className="space-y-3">
              <div className="relative">
                <select className="appearance-none block w-full rounded-xl border border-border bg-card py-2.5 pl-4 pr-10 text-sm text-foreground focus:border-foreground focus:ring-0 transition-all shadow-sm hover:border-foreground/50 cursor-pointer outline-none">
                  <option value="pending">Pendiente</option>
                  <option value="reviewing">En Revisión</option>
                  <option value="resolved">Resuelta (Confirmado Fraude)</option>
                  <option value="false_positive">Falso Positivo</option>
                </select>
                <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-3 text-muted-foreground">
                  <svg
                    className="h-4 w-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M19 9l-7 7-7-7"
                    ></path>
                  </svg>
                </div>
              </div>

              <button className="w-full rounded-xl bg-foreground px-4 py-2.5 text-sm font-medium text-background shadow-sm hover:bg-foreground/90 transition-colors">
                Actualizar Estado
              </button>
            </div>
            <div className="mt-6 pt-6 border-t border-border">
              <div className="flex justify-between items-center mb-3">
                <span className="text-sm text-muted-foreground">
                  Asignado a
                </span>
                <button className="text-xs text-foreground font-medium hover:underline">
                  Cambiar
                </button>
              </div>
              <div className="flex items-center gap-3 p-2 rounded-xl hover:bg-secondary/50 transition-colors cursor-pointer">
                <div className="h-8 w-8 rounded-full bg-secondary flex items-center justify-center text-xs font-medium text-muted-foreground border border-border">
                  ?
                </div>
                <span className="text-sm font-medium text-foreground">
                  {alert.reviewed_by || "Sin asignar"}
                </span>
              </div>
            </div>
          </div>

          {/* Related Transaction */}
          <div className="rounded-3xl border border-border bg-card p-6 shadow-zen">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-foreground">
                Transacción Relacionada
              </h3>
              {alert.transaction_id && (
                <Link
                  href={`/transactions/${alert.transaction_id}`}
                  className="text-xs text-foreground font-medium hover:underline"
                >
                  Ver detalles
                </Link>
              )}
            </div>
            {transaction ? (
              <div className="bg-secondary/30 rounded-2xl p-4 space-y-4 border border-border/50">
                <div className="flex justify-between items-start">
                  <div>
                    <p className="text-lg font-medium text-foreground">
                      $
                      {transaction.amount.toLocaleString("es-MX", {
                        minimumFractionDigits: 2,
                      })}
                    </p>
                    <p className="text-xs text-muted-foreground font-medium">
                      {transaction.currency}
                    </p>
                  </div>
                  <span className="text-xs font-mono text-muted-foreground bg-card px-1.5 py-0.5 rounded border border-border">
                    {transaction.code}
                  </span>
                </div>
                <div className="space-y-2 pt-2 border-t border-border/50">
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <CreditCard className="h-3.5 w-3.5" />
                    <span className="truncate">{transaction.merchant}</span>
                  </div>
                  <div className="flex items-center gap-2 text-sm text-muted-foreground">
                    <Clock className="h-3.5 w-3.5" />
                    <span>
                      {new Date(transaction.created_at).toLocaleString()}
                    </span>
                  </div>
                </div>
              </div>
            ) : (
              <div className="p-4 text-center text-sm text-muted-foreground bg-secondary/30 rounded-2xl border border-border/50">
                {isLoadingAlert || !alert.transaction_id
                  ? "Cargando..."
                  : "Información no disponible"}
              </div>
            )}
          </div>

          {/* Related Client */}
          <div className="rounded-3xl border border-border bg-card p-6 shadow-zen">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-foreground">
                Cliente Afectado
              </h3>
              {alert.client_id && (
                <Link
                  href={`/clients/${alert.client_id}`}
                  className="text-xs text-foreground font-medium hover:underline"
                >
                  Ver perfil
                </Link>
              )}
            </div>
            {client ? (
              <div className="flex items-center gap-4 p-2 rounded-xl hover:bg-secondary/50 transition-colors cursor-pointer">
                <div className="h-10 w-10 rounded-full bg-secondary border border-border flex items-center justify-center text-muted-foreground font-medium">
                  {client.full_name.charAt(0)}
                </div>
                <div>
                  <p className="text-sm font-medium text-foreground">
                    {client.full_name}
                  </p>
                  <p className="text-xs text-muted-foreground">
                    {client.document_type}: {client.document_number}
                  </p>
                </div>
              </div>
            ) : (
              <div className="p-4 text-center text-sm text-muted-foreground bg-secondary/30 rounded-2xl border border-border/50">
                {isLoadingAlert || !alert.client_id
                  ? "Cargando..."
                  : "Información no disponible"}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
