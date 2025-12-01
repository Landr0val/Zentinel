"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import {
  ArrowLeft,
  Shield,
  AlertTriangle,
  CheckCircle,
  XCircle,
  MapPin,
  Clock,
  CreditCard,
  BrainCircuit,
  Loader2,
  AlertCircle,
} from "lucide-react";
import { cn } from "../../../lib/utils";
import { useQuery } from "@tanstack/react-query";
import { api } from "../../../lib/api";
import { createSlug, getIdFromSlug } from "../../../lib/slug-manager";
import { useCurrency } from "../../../context/CurrencyContext";

export default function TransactionDetailPage() {
  const params = useParams();
  const slug = params?.id as string;
  const id = getIdFromSlug(slug) || slug;
  const { formatAmount } = useCurrency();

  const {
    data: transactionRes,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["transaction", id],
    queryFn: () => api.transactions.get(id),
  });

  const transaction = transactionRes?.data;

  // Fetch account details to get client info if available
  const { data: accountRes } = useQuery({
    queryKey: ["account", transaction?.account_id],
    queryFn: () => api.accounts.get(transaction!.account_id),
    enabled: !!transaction?.account_id,
  });
  const account = accountRes?.data;

  if (isLoading) {
    return (
      <div className="flex h-96 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error || !transaction) {
    return (
      <div className="flex h-96 flex-col items-center justify-center text-center">
        <AlertCircle className="h-12 w-12 text-destructive mb-4" />
        <h2 className="text-lg font-semibold">
          Error al cargar la transacción
        </h2>
        <p className="text-muted-foreground">
          No se pudo encontrar la información solicitada.
        </p>
      </div>
    );
  }

  // Mock AI Analysis (since backend doesn't persist full explanation on transaction yet, only on alert)
  const aiAnalysis = {
    explanation:
      transaction.risk_score > 50
        ? "La transacción presenta un comportamiento anómalo. El monto o la ubicación difieren significativamente del patrón habitual del cliente."
        : "La transacción se ajusta a los patrones de comportamiento habituales del cliente. No se detectaron anomalías significativas.",
    riskFactors:
      transaction.risk_score > 50
        ? [
            "Monto inusual para el canal",
            "Ubicación geográfica atípica",
            "Hora no habitual",
          ]
        : ["Comportamiento normal"],
  };

  return (
    <div className="animate-in fade-in duration-500">
      {/* Header */}
      <div className="flex items-center gap-4 mb-8">
        <Link
          href="/transactions"
          className="p-2 rounded-full hover:bg-secondary text-muted-foreground hover:text-foreground transition-colors"
        >
          <ArrowLeft className="h-5 w-5" />
        </Link>
        <div className="flex-1">
          <h1 className="text-3xl font-light tracking-tight text-foreground">
            Transacción {transaction.code}
          </h1>
          <p className="mt-1 text-muted-foreground">
            Detalles y Análisis de Riesgo
          </p>
        </div>
        <div className="flex gap-2">
          <span
            className={cn(
              "inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-sm font-medium border",
              transaction.status === "completed"
                ? "bg-green-500/10 text-green-600 border-green-500/20"
                : transaction.status === "flagged" ||
                    transaction.status === "failed"
                  ? "bg-destructive/10 text-destructive border-destructive/20"
                  : "bg-secondary text-muted-foreground border-border",
            )}
          >
            {transaction.status === "flagged" ||
            transaction.status === "failed" ? (
              <AlertTriangle className="h-3.5 w-3.5" />
            ) : transaction.status === "completed" ? (
              <CheckCircle className="h-3.5 w-3.5" />
            ) : (
              <Clock className="h-3.5 w-3.5" />
            )}
            <span className="capitalize">
              {transaction.status === "flagged"
                ? "Sospechosa"
                : transaction.status === "failed"
                  ? "Fallida"
                  : transaction.status === "completed"
                    ? "Completada"
                    : "Pendiente"}
            </span>
          </span>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Main Details */}
        <div className="lg:col-span-2 space-y-6">
          {/* AI Analysis Card */}
          <div className="rounded-3xl border border-border bg-card p-6 shadow-zen overflow-hidden relative">
            <div className="absolute top-0 right-0 p-4 opacity-5">
              <BrainCircuit className="h-32 w-32 text-foreground" />
            </div>
            <div className="relative z-10">
              <div className="flex items-center gap-2 mb-6">
                <div className="p-2 rounded-xl bg-primary/10 text-primary">
                  <BrainCircuit className="h-5 w-5" />
                </div>
                <h3 className="text-lg font-medium text-foreground">
                  Análisis de IA (Sentinel)
                </h3>
              </div>

              <div className="flex flex-col sm:flex-row items-center gap-8 mb-6">
                <div className="text-center shrink-0">
                  <div
                    className={cn(
                      "text-5xl font-light tracking-tighter",
                      transaction.risk_score > 70
                        ? "text-destructive"
                        : transaction.risk_score > 30
                          ? "text-yellow-600"
                          : "text-green-600",
                    )}
                  >
                    {transaction.risk_score}
                    <span className="text-2xl text-muted-foreground">/100</span>
                  </div>
                  <div className="text-xs font-medium text-muted-foreground uppercase tracking-widest mt-2">
                    Score de Riesgo
                  </div>
                </div>
                <div className="hidden sm:block h-16 w-px bg-border"></div>
                <div className="flex-1 text-center sm:text-left">
                  <p className="text-muted-foreground leading-relaxed">
                    {aiAnalysis.explanation}
                  </p>
                </div>
              </div>

              <div className="bg-secondary/30 rounded-2xl p-5 border border-border/50">
                <h4 className="text-sm font-medium text-foreground mb-3">
                  Factores de Riesgo Detectados:
                </h4>
                <ul className="space-y-2">
                  {aiAnalysis.riskFactors.map((factor, index) => (
                    <li
                      key={index}
                      className="flex items-center gap-2.5 text-sm text-muted-foreground"
                    >
                      <div className="h-1.5 w-1.5 rounded-full bg-orange-500 shrink-0" />
                      {factor}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>

          {/* Transaction Details */}
          <div className="rounded-3xl border border-border bg-card p-6 shadow-zen">
            <h3 className="text-lg font-medium text-foreground mb-6">
              Detalles de la Operación
            </h3>
            <dl className="grid grid-cols-1 gap-x-4 gap-y-8 sm:grid-cols-2">
              <div>
                <dt className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">
                  Monto
                </dt>
                <dd className="text-3xl font-light text-foreground tracking-tight">
                  {formatAmount(transaction.amount, transaction.currency)}
                </dd>
              </div>
              <div>
                <dt className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">
                  Comercio / Destino
                </dt>
                <dd className="text-lg font-medium text-foreground">
                  {transaction.merchant}
                </dd>
              </div>
              <div>
                <dt className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">
                  Fecha y Hora
                </dt>
                <dd className="text-sm text-foreground flex items-center gap-2">
                  <Clock className="h-4 w-4 text-muted-foreground" />
                  {new Date(transaction.created_at).toLocaleString()}
                </dd>
              </div>
              <div>
                <dt className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">
                  Ubicación
                </dt>
                <dd className="text-sm text-foreground flex items-center gap-2">
                  <MapPin className="h-4 w-4 text-muted-foreground" />
                  {transaction.city}, {transaction.country}
                </dd>
              </div>
              <div>
                <dt className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">
                  Tipo de Operación
                </dt>
                <dd className="text-sm text-foreground capitalize bg-secondary inline-block px-2 py-1 rounded-md">
                  {transaction.operation_type}
                </dd>
              </div>
              <div>
                <dt className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">
                  Canal
                </dt>
                <dd className="text-sm text-foreground capitalize">
                  {transaction.channel}
                </dd>
              </div>
            </dl>
          </div>
        </div>

        {/* Sidebar Info */}
        <div className="space-y-6">
          {/* Actions */}
          <div className="rounded-3xl border border-border bg-card p-6 shadow-zen">
            <h3 className="text-sm font-medium text-foreground mb-4">
              Acciones Disponibles
            </h3>
            <div className="space-y-3">
              <button className="w-full flex justify-center items-center gap-2 rounded-xl bg-foreground px-4 py-2.5 text-sm font-medium text-background shadow-sm hover:bg-foreground/90 transition-colors">
                <Shield className="h-4 w-4" />
                Investigar Más
              </button>
              <button className="w-full flex justify-center items-center gap-2 rounded-xl bg-secondary px-4 py-2.5 text-sm font-medium text-foreground shadow-sm hover:bg-secondary/80 transition-colors border border-border">
                <AlertTriangle className="h-4 w-4" />
                Reportar Anomalía
              </button>
            </div>
          </div>

          {/* Related Entities */}
          <div className="rounded-3xl border border-border bg-card p-6 shadow-zen">
            <h3 className="text-sm font-medium text-foreground mb-4">
              Entidades Relacionadas
            </h3>

            <div className="space-y-6">
              {account && (
                <div>
                  <div className="flex items-center justify-between mb-3">
                    <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">
                      Cliente
                    </span>
                    <Link
                      href={`/clients/${createSlug(account.client_id)}`}
                      className="text-xs text-foreground font-medium hover:underline"
                    >
                      Ver perfil
                    </Link>
                  </div>
                  <div className="flex items-center gap-3 p-2 rounded-xl hover:bg-secondary/50 transition-colors cursor-pointer">
                    <div className="h-10 w-10 rounded-full bg-secondary border border-border flex items-center justify-center text-muted-foreground font-medium">
                      {account.client_name?.charAt(0) || "?"}
                    </div>
                    <div>
                      <p className="text-sm font-medium text-foreground">
                        {account.client_name}
                      </p>
                      <p className="text-xs text-muted-foreground">
                        ID: {account.client_id.substring(0, 8)}...
                      </p>
                    </div>
                  </div>
                </div>
              )}

              <div className="border-t border-border pt-4">
                <div className="flex items-center justify-between mb-3">
                  <span className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">
                    Cuenta Origen
                  </span>
                  <Link
                    href={`/accounts/${createSlug(transaction.account_id)}`}
                    className="text-xs text-foreground font-medium hover:underline"
                  >
                    Ver cuenta
                  </Link>
                </div>
                <div className="flex items-center gap-3 p-2 rounded-xl bg-secondary/30 border border-border/50">
                  <div className="h-10 w-10 rounded-full bg-card border border-border flex items-center justify-center text-foreground">
                    <CreditCard className="h-4 w-4" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-foreground font-mono">
                      {account?.account_number || "Cargando..."}
                    </p>
                    <p className="text-xs text-muted-foreground capitalize">
                      {account?.type || "Cuenta"} • {account?.currency || "MXN"}
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
