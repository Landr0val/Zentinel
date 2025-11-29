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
} from "lucide-react";

// Mock data
const alertData = {
  id: "alt_1",
  transactionId: "txn_3",
  severity: "high",
  status: "new",
  description: "Transacción inusual en ubicación geográfica distante",
  createdAt: "2024-03-12 13:15:05",
  updatedAt: "2024-03-12 13:15:05",
  assignedTo: "Sin asignar",
  ruleId: "GEO-VEL-001",
  ruleName: "Velocidad Geográfica Imposible",
  notes: [
    {
      id: 1,
      author: "Sistema Sentinel",
      text: "Alerta generada automáticamente. La distancia entre la transacción anterior y esta es de 500km en menos de 30 minutos.",
      date: "2024-03-12 13:15:05",
      type: "system",
    },
  ],
  transaction: {
    id: "txn_3",
    amount: 8000.0,
    currency: "MXN",
    merchant: "Cajero ATM #402",
    date: "2024-03-12 13:15:00",
    clientName: "Juan Pérez",
    clientId: "1",
  },
};

export default function AlertDetailPage() {
  const params = useParams();
  const id = params?.id as string;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <Link
          href="/alerts"
          className="p-2 rounded-full hover:bg-slate-100 text-slate-500 transition-colors"
        >
          <ArrowLeft className="h-5 w-5" />
        </Link>
        <div className="flex-1">
          <h1 className="text-2xl font-bold text-slate-900">Alerta {id}</h1>
          <p className="text-slate-500">Gestión de Incidencia</p>
        </div>
        <div className="flex gap-2">
          <span
            className={`inline-flex items-center gap-1 rounded-full px-3 py-1 text-sm font-medium ${
              alertData.status === "resolved"
                ? "bg-green-100 text-green-800"
                : alertData.status === "new"
                  ? "bg-blue-100 text-blue-800"
                  : alertData.status === "investigating"
                    ? "bg-yellow-100 text-yellow-800"
                    : "bg-slate-100 text-slate-800"
            }`}
          >
            {alertData.status === "resolved" && (
              <CheckCircle className="h-4 w-4" />
            )}
            {alertData.status === "new" && (
              <AlertTriangle className="h-4 w-4" />
            )}
            {alertData.status === "investigating" && (
              <Clock className="h-4 w-4" />
            )}
            {alertData.status === "false_positive" && (
              <XCircle className="h-4 w-4" />
            )}
            <span className="capitalize">
              {alertData.status === "resolved"
                ? "Resuelta"
                : alertData.status === "new"
                  ? "Nueva"
                  : alertData.status === "investigating"
                    ? "Investigando"
                    : "Falso Positivo"}
            </span>
          </span>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Main Content */}
        <div className="lg:col-span-2 space-y-6">
          {/* Alert Details */}
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <div className="flex items-start gap-4">
              <div
                className={`p-3 rounded-full ${
                  alertData.severity === "high"
                    ? "bg-red-100 text-red-600"
                    : alertData.severity === "medium"
                      ? "bg-yellow-100 text-yellow-600"
                      : "bg-blue-100 text-blue-600"
                }`}
              >
                <AlertTriangle className="h-6 w-6" />
              </div>
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-slate-900">
                  {alertData.description}
                </h3>
                <p className="text-slate-500 mt-1">
                  Regla activada:{" "}
                  <span className="font-mono text-xs bg-slate-100 px-1 py-0.5 rounded">
                    {alertData.ruleId}
                  </span>{" "}
                  {alertData.ruleName}
                </p>
              </div>
            </div>
          </div>

          {/* Timeline / Notes */}
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <h3 className="text-lg font-semibold text-slate-900 mb-6">
              Historial y Notas
            </h3>
            <div className="space-y-6">
              {alertData.notes.map((note) => (
                <div key={note.id} className="flex gap-4">
                  <div className="flex flex-col items-center">
                    <div className="h-8 w-8 rounded-full bg-slate-100 flex items-center justify-center">
                      {note.type === "system" ? (
                        <Shield className="h-4 w-4 text-slate-500" />
                      ) : (
                        <User className="h-4 w-4 text-slate-500" />
                      )}
                    </div>
                    <div className="w-px h-full bg-slate-200 my-2"></div>
                  </div>
                  <div className="flex-1 pb-4">
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-sm font-medium text-slate-900">
                        {note.author}
                      </span>
                      <span className="text-xs text-slate-500">
                        {note.date}
                      </span>
                    </div>
                    <div className="text-sm text-slate-600 bg-slate-50 p-3 rounded-lg">
                      {note.text}
                    </div>
                  </div>
                </div>
              ))}

              {/* Add Note Input */}
              <div className="flex gap-4 pt-2">
                <div className="h-8 w-8 rounded-full bg-blue-100 flex items-center justify-center shrink-0">
                  <User className="h-4 w-4 text-blue-600" />
                </div>
                <div className="flex-1">
                  <div className="relative">
                    <textarea
                      rows={3}
                      className="block w-full rounded-md border-0 py-1.5 text-slate-900 shadow-sm ring-1 ring-inset ring-slate-300 placeholder:text-slate-400 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6"
                      placeholder="Agregar una nota de investigación..."
                    />
                    <div className="absolute bottom-2 right-2">
                      <button className="inline-flex items-center rounded-md bg-blue-600 px-2.5 py-1.5 text-xs font-semibold text-white shadow-sm hover:bg-blue-500">
                        <Send className="h-3 w-3 mr-1" />
                        Agregar
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Sidebar Info */}
        <div className="space-y-6">
          {/* Status Actions */}
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <h3 className="text-sm font-medium text-slate-900 mb-4">
              Gestionar Estado
            </h3>
            <div className="space-y-3">
              <select className="block w-full rounded-md border-0 py-1.5 pl-3 pr-10 text-slate-900 ring-1 ring-inset ring-slate-300 focus:ring-2 focus:ring-inset focus:ring-blue-600 sm:text-sm sm:leading-6">
                <option value="new">Nueva</option>
                <option value="investigating">En Investigación</option>
                <option value="resolved">Resuelta (Confirmado Fraude)</option>
                <option value="false_positive">Falso Positivo</option>
              </select>
              <button className="w-full rounded-md bg-blue-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500">
                Actualizar Estado
              </button>
            </div>
            <div className="mt-6 pt-6 border-t border-slate-100">
              <div className="flex justify-between items-center mb-2">
                <span className="text-sm text-slate-500">Asignado a</span>
                <button className="text-xs text-blue-600 hover:underline">
                  Cambiar
                </button>
              </div>
              <div className="flex items-center gap-2">
                <div className="h-6 w-6 rounded-full bg-slate-200 flex items-center justify-center text-xs font-medium text-slate-600">
                  ?
                </div>
                <span className="text-sm font-medium text-slate-900">
                  {alertData.assignedTo}
                </span>
              </div>
            </div>
          </div>

          {/* Related Transaction */}
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-slate-900">
                Transacción Relacionada
              </h3>
              <Link
                href={`/transactions/${alertData.transactionId}`}
                className="text-xs text-blue-600 hover:underline"
              >
                Ver detalles
              </Link>
            </div>
            <div className="bg-slate-50 rounded-lg p-4 space-y-3">
              <div className="flex justify-between items-start">
                <div>
                  <p className="text-sm font-bold text-slate-900">
                    $
                    {alertData.transaction.amount.toLocaleString("es-MX", {
                      minimumFractionDigits: 2,
                    })}
                  </p>
                  <p className="text-xs text-slate-500">
                    {alertData.transaction.currency}
                  </p>
                </div>
                <span className="text-xs font-mono text-slate-400">
                  {alertData.transactionId}
                </span>
              </div>
              <div className="flex items-center gap-2 text-sm text-slate-700">
                <CreditCard className="h-4 w-4 text-slate-400" />
                {alertData.transaction.merchant}
              </div>
              <div className="flex items-center gap-2 text-sm text-slate-700">
                <Clock className="h-4 w-4 text-slate-400" />
                {alertData.transaction.date}
              </div>
            </div>
          </div>

          {/* Related Client */}
          <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-sm font-medium text-slate-900">
                Cliente Afectado
              </h3>
              <Link
                href={`/clients/${alertData.transaction.clientId}`}
                className="text-xs text-blue-600 hover:underline"
              >
                Ver perfil
              </Link>
            </div>
            <div className="flex items-center gap-3">
              <div className="h-10 w-10 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 font-bold">
                {alertData.transaction.clientName.charAt(0)}
              </div>
              <div>
                <p className="text-sm font-medium text-slate-900">
                  {alertData.transaction.clientName}
                </p>
                <p className="text-xs text-slate-500">
                  ID: {alertData.transaction.clientId}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
