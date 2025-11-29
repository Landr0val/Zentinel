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
} from "lucide-react";

// Mock data
const transactionData = {
  id: "txn_3",
  amount: 8000.0,
  currency: "MXN",
  date: "2024-03-12 13:15:00",
  merchant: "Cajero ATM #402",
  location: "Ciudad de México, Centro",
  status: "flagged",
  type: "withdrawal",
  accountNumber: "1234-5678-9012",
  clientName: "Juan Pérez",
  clientId: "1",
  accountId: "1",
  riskScore: 85,
  riskLevel: "high",
  aiAnalysis:
    "La transacción presenta un comportamiento anómalo. El monto retirado excede el promedio histórico de retiros en cajeros automáticos para este cliente ($2,000 MXN). Además, la ubicación del cajero difiere del patrón habitual de geolocalización del cliente en los últimos 30 días.",
  riskFactors: [
    "Monto inusual para el canal (ATM)",
    "Ubicación geográfica atípica",
    "Hora no habitual para este tipo de operación",
  ],
};

export default function TransactionDetailPage() {
  const params = useParams();
  const id = params?.id as string;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <Link
          href="/transactions"
          className="p-2 rounded-full hover:bg-slate-100 text-slate-500 transition-colors"
        >
          <ArrowLeft className="h-5 w-5" />
        </Link>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-slate-900">
            Transacción {id}
          </h1>
          <p className="text-slate-500">Detalles y Análisis de Riesgo</p>
        </div>
        <div className="flex gap-2">
          <span
            className={`inline-flex items-center gap-1 rounded-full px-3 py-1 text-sm font-medium ${
              transactionData.status === "completed"
                ? "bg-green-100 text-green-800"
                : transactionData.status === "flagged"
                  ? "bg-red-100 text-red-800"
                  : "bg-yellow-100 text-yellow-800"
            }`}
          >
            {transactionData.status === "flagged" && (
              <AlertTriangle className="h-4 w-4" />
            )}
            {transactionData.status === "completed" && (
              <CheckCircle className="h-4 w-4" />
            )}
            <span className="capitalize">
              {transactionData.status === "flagged"
                ? "Sospechosa"
                : transactionData.status === "completed"
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
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100 overflow-hidden relative">
            <div className="absolute top-0 right-0 p-4 opacity-10">
              <BrainCircuit className="h-32 w-32 text-blue-600" />
            </div>
            <div className="relative z-10">
              <div className="flex items-center gap-2 mb-4">
                <BrainCircuit className="h-6 w-6 text-blue-600" />
                <h3 className="text-lg font-semibold text-slate-900">
                  Análisis de IA (Sentinel)
                </h3>
              </div>

              <div className="flex items-center gap-6 mb-6">
                <div className="text-center">
                  <div
                    className={`text-4xl font-bold ${
                      transactionData.riskScore > 70
                        ? "text-red-600"
                        : transactionData.riskScore > 30
                          ? "text-yellow-600"
                          : "text-green-600"
                    }`}
                  >
                    {transactionData.riskScore}/100
                  </div>
                  <div className="text-xs font-medium text-slate-500 uppercase tracking-wide mt-1">
                    Score de Riesgo
                  </div>
                </div>
                <div className="h-12 w-px bg-slate-200"></div>
                <div className="flex-1">
                  <p className="text-slate-700 leading-relaxed">
                    {transactionData.aiAnalysis}
                  </p>
                </div>
              </div>

              <div className="bg-slate-50 rounded-lg p-4">
                <h4 className="text-sm font-medium text-slate-900 mb-3">
                  Factores de Riesgo Detectados:
                </h4>
                <ul className="space-y-2">
                  {transactionData.riskFactors.map((factor, index) => (
                    <li
                      key={index}
                      className="flex items-start gap-2 text-sm text-slate-600"
                    >
                      <AlertTriangle className="h-4 w-4 text-orange-500 mt-0.5 shrink-0" />
                      {factor}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>

          {/* Transaction Details */}
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <h3 className="text-lg font-semibold text-slate-900 mb-4">
              Detalles de la Operación
            </h3>
            <dl className="grid grid-cols-1 gap-x-4 gap-y-6 sm:grid-cols-2">
              <div>
                <dt className="text-sm font-medium text-slate-500">Monto</dt>
                <dd className="mt-1 text-2xl font-bold text-slate-900">
                  $
                  {transactionData.amount.toLocaleString("es-MX", {
                    minimumFractionDigits: 2,
                  })}{" "}
                  <span className="text-sm font-normal text-slate-500">
                    {transactionData.currency}
                  </span>
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-slate-500">
                  Comercio / Destino
                </dt>
                <dd className="mt-1 text-lg font-medium text-slate-900">
                  {transactionData.merchant}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-slate-500">
                  Fecha y Hora
                </dt>
                <dd className="mt-1 text-sm text-slate-900 flex items-center gap-2">
                  <Clock className="h-4 w-4 text-slate-400" />
                  {transactionData.date}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-slate-500">
                  Ubicación
                </dt>
                <dd className="mt-1 text-sm text-slate-900 flex items-center gap-2">
                  <MapPin className="h-4 w-4 text-slate-400" />
                  {transactionData.location}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-slate-500">Tipo</dt>
                <dd className="mt-1 text-sm text-slate-900 capitalize">
                  {transactionData.type}
                </dd>
              </div>
            </dl>
          </div>
        </div>

        {/* Sidebar Info */}
        <div className="space-y-6">
          {/* Actions */}
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <h3 className="text-sm font-medium text-slate-900 mb-4">
              Acciones Disponibles
            </h3>
            <div className="space-y-3">
              <button className="w-full flex justify-center items-center gap-2 rounded-md bg-green-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-green-500">
                <CheckCircle className="h-4 w-4" />
                Aprobar Transacción
              </button>
              <button className="w-full flex justify-center items-center gap-2 rounded-md bg-red-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-red-500">
                <XCircle className="h-4 w-4" />
                Bloquear y Reportar
              </button>
              <button className="w-full flex justify-center items-center gap-2 rounded-md bg-white px-3 py-2 text-sm font-semibold text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 hover:bg-slate-50">
                <Shield className="h-4 w-4" />
                Investigar Más
              </button>
            </div>
          </div>

          {/* Related Entities */}
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <h3 className="text-sm font-medium text-slate-900 mb-4">
              Entidades Relacionadas
            </h3>

            <div className="space-y-6">
              <div>
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs font-medium text-slate-500 uppercase">
                    Cliente
                  </span>
                  <Link
                    href={`/clients/${transactionData.clientId}`}
                    className="text-xs text-blue-600 hover:underline"
                  >
                    Ver perfil
                  </Link>
                </div>
                <div className="flex items-center gap-3">
                  <div className="h-8 w-8 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 font-bold text-xs">
                    {transactionData.clientName.charAt(0)}
                  </div>
                  <div>
                    <p className="text-sm font-medium text-slate-900">
                      {transactionData.clientName}
                    </p>
                    <p className="text-xs text-slate-500">
                      ID: {transactionData.clientId}
                    </p>
                  </div>
                </div>
              </div>

              <div className="border-t border-slate-100 pt-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-xs font-medium text-slate-500 uppercase">
                    Cuenta Origen
                  </span>
                  <Link
                    href={`/accounts/${transactionData.accountId}`}
                    className="text-xs text-blue-600 hover:underline"
                  >
                    Ver cuenta
                  </Link>
                </div>
                <div className="flex items-center gap-3">
                  <div className="h-8 w-8 rounded-full bg-blue-50 flex items-center justify-center text-blue-600">
                    <CreditCard className="h-4 w-4" />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-slate-900 font-mono">
                      {transactionData.accountNumber}
                    </p>
                    <p className="text-xs text-slate-500">Ahorro • MXN</p>
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
