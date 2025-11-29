"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import {
  ArrowLeft,
  CreditCard,
  DollarSign,
  Calendar,
  Activity,
  ArrowUpRight,
  ArrowDownRight,
} from "lucide-react";

// Mock data
const accountData = {
  id: "1",
  accountNumber: "1234-5678-9012",
  clientName: "Juan Pérez",
  clientId: "1",
  balance: 15420.50,
  currency: "MXN",
  status: "active",
  type: "savings",
  createdAt: "2023-05-20",
  lastActivity: "2024-03-12 14:30",
};

const recentTransactions = [
  {
    id: "txn_1",
    type: "payment",
    amount: 1250.00,
    merchant: "Amazon MX",
    date: "2024-03-12 14:30:00",
    status: "completed",
  },
  {
    id: "txn_3",
    type: "withdrawal",
    amount: 500.00,
    merchant: "Cajero ATM",
    date: "2024-03-11 09:15:00",
    status: "completed",
  },
  {
    id: "txn_5",
    type: "deposit",
    amount: 2000.00,
    merchant: "Transferencia SPEI",
    date: "2024-03-10 18:45:00",
    status: "completed",
  },
];

export default function AccountDetailPage() {
  const params = useParams();
  const id = params?.id as string;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <Link
          href="/accounts"
          className="p-2 rounded-full hover:bg-slate-100 text-slate-500 transition-colors"
        >
          <ArrowLeft className="h-5 w-5" />
        </Link>
        <div className="flex-1">
            <h1 className="text-2xl font-bold text-slate-900">Cuenta {accountData.accountNumber}</h1>
            <p className="text-slate-500">ID: {id}</p>
        </div>
        <div className="flex gap-2">
             <span
                className={`inline-flex items-center rounded-full px-3 py-1 text-sm font-medium ${
                accountData.status === "active"
                    ? "bg-green-100 text-green-800"
                    : "bg-slate-100 text-slate-800"
                }`}
            >
                {accountData.status === "active" ? "Activa" : "Inactiva"}
            </span>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        {/* Summary Card */}
        <div className="lg:col-span-2 space-y-6">
             <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
                <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
                    <div className="flex items-center justify-between mb-4">
                        <h3 className="text-sm font-medium text-slate-500">Saldo Disponible</h3>
                        <DollarSign className="h-5 w-5 text-slate-400" />
                    </div>
                    <div className="text-3xl font-bold text-slate-900">
                        ${accountData.balance.toLocaleString('es-MX', { minimumFractionDigits: 2 })}
                    </div>
                    <p className="text-sm text-slate-500 mt-1">{accountData.currency}</p>
                </div>
                <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
                    <div className="flex items-center justify-between mb-4">
                        <h3 className="text-sm font-medium text-slate-500">Tipo de Cuenta</h3>
                        <CreditCard className="h-5 w-5 text-slate-400" />
                    </div>
                    <div className="text-3xl font-bold text-slate-900 capitalize">
                        {accountData.type === 'savings' ? 'Ahorro' : 'Corriente'}
                    </div>
                    <p className="text-sm text-slate-500 mt-1">Banca Personal</p>
                </div>
             </div>

             {/* Transactions List */}
             <div className="rounded-xl bg-white shadow-sm border border-slate-100 overflow-hidden">
                <div className="px-6 py-5 border-b border-slate-100 flex items-center justify-between">
                    <h3 className="text-lg font-semibold text-slate-900">Últimos Movimientos</h3>
                    <Link href="/transactions" className="text-sm font-medium text-blue-600 hover:text-blue-500">
                        Ver historial completo
                    </Link>
                </div>
                <ul role="list" className="divide-y divide-slate-100">
                    {recentTransactions.map((txn) => (
                        <li key={txn.id} className="px-6 py-4 hover:bg-slate-50 transition-colors">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center">
                                    <div className={`rounded-full p-2 mr-4 ${
                                        txn.type === 'deposit' ? 'bg-green-100 text-green-600' : 'bg-slate-100 text-slate-600'
                                    }`}>
                                        {txn.type === 'deposit' ? <ArrowDownRight className="h-4 w-4" /> : <ArrowUpRight className="h-4 w-4" />}
                                    </div>
                                    <div>
                                        <p className="text-sm font-medium text-slate-900">{txn.merchant}</p>
                                        <p className="text-xs text-slate-500">{txn.date}</p>
                                    </div>
                                </div>
                                <div className="text-right">
                                    <p className={`text-sm font-medium ${
                                        txn.type === 'deposit' ? 'text-green-600' : 'text-slate-900'
                                    }`}>
                                        {txn.type === 'deposit' ? '+' : '-'}${txn.amount.toLocaleString('es-MX', { minimumFractionDigits: 2 })}
                                    </p>
                                    <p className="text-xs text-slate-500 capitalize">{txn.type}</p>
                                </div>
                            </div>
                        </li>
                    ))}
                </ul>
             </div>
        </div>

        {/* Sidebar Info */}
        <div className="space-y-6">
            <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
                <h3 className="text-lg font-semibold text-slate-900 mb-4">Detalles del Cliente</h3>
                <div className="flex items-center mb-6">
                    <div className="h-12 w-12 rounded-full bg-slate-100 flex items-center justify-center text-slate-500 font-bold text-lg">
                        {accountData.clientName.charAt(0)}
                    </div>
                    <div className="ml-4">
                        <div className="text-base font-medium text-slate-900">{accountData.clientName}</div>
                        <Link href={`/clients/${accountData.clientId}`} className="text-sm text-blue-600 hover:text-blue-500">
                            Ver perfil
                        </Link>
                    </div>
                </div>
                <div className="space-y-4 border-t border-slate-100 pt-4">
                    <div className="flex items-center justify-between">
                        <div className="flex items-center text-sm text-slate-500">
                            <Calendar className="mr-2 h-4 w-4" />
                            Apertura
                        </div>
                        <span className="text-sm font-medium text-slate-900">{accountData.createdAt}</span>
                    </div>
                    <div className="flex items-center justify-between">
                        <div className="flex items-center text-sm text-slate-500">
                            <Activity className="mr-2 h-4 w-4" />
                            Última Actividad
                        </div>
                        <span className="text-sm font-medium text-slate-900">{accountData.lastActivity.split(' ')[0]}</span>
                    </div>
                </div>
            </div>
        </div>
      </div>
    </div>
  );
}
