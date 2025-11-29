"use client";

import {
  Activity,
  AlertTriangle,
  ArrowDownRight,
  ArrowUpRight,
  DollarSign,
  Users,
  Loader2,
} from "lucide-react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import Link from "next/link";

export default function Dashboard() {
  const {
    data: statsData,
    isLoading,
    error,
  } = useQuery({
    queryKey: ["dashboardStats"],
    queryFn: () => api.dashboard.getStats(),
  });

  if (isLoading) {
    return (
      <div className="flex h-[50vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-lg bg-red-50 p-4 text-red-800">
        <h3 className="text-lg font-medium">Error al cargar datos</h3>
        <p>No se pudo conectar con el servidor. Por favor intente más tarde.</p>
      </div>
    );
  }

  const stats = statsData?.data;

  const statCards = [
    {
      name: "Volumen Total",
      value: stats ? `$${stats.total_volume.toLocaleString()}` : "-",
      change: stats
        ? `${stats.volume_change_percentage > 0 ? "+" : ""}${stats.volume_change_percentage}%`
        : "-",
      trend: stats && stats.volume_change_percentage >= 0 ? "up" : "down",
      icon: DollarSign,
      color: "text-green-600",
      bg: "bg-green-100",
    },
    {
      name: "Transacciones",
      value: stats ? stats.transaction_count.toLocaleString() : "-",
      change: stats
        ? `${stats.transaction_change_percentage > 0 ? "+" : ""}${stats.transaction_change_percentage}%`
        : "-",
      trend: stats && stats.transaction_change_percentage >= 0 ? "up" : "down",
      icon: Activity,
      color: "text-blue-600",
      bg: "bg-blue-100",
    },
    {
      name: "Alertas de Fraude",
      value: stats ? stats.alert_count.toLocaleString() : "-",
      change: stats
        ? `${stats.alert_change_percentage > 0 ? "+" : ""}${stats.alert_change_percentage}%`
        : "-",
      trend: stats && stats.alert_change_percentage <= 0 ? "up" : "down", // Logic handled in render
      icon: AlertTriangle,
      color: "text-red-600",
      bg: "bg-red-100",
    },
    {
      name: "Clientes Activos",
      value: stats ? stats.active_clients.toLocaleString() : "-",
      change: stats
        ? `${stats.client_change_percentage > 0 ? "+" : ""}${stats.client_change_percentage}%`
        : "-",
      trend: stats && stats.client_change_percentage >= 0 ? "up" : "down",
      icon: Users,
      color: "text-purple-600",
      bg: "bg-purple-100",
    },
  ];

  // Format chart data
  const chartData =
    stats?.weekly_activity?.map((day) => ({
      name: day.date,
      transacciones: day.transactions,
      fraudes: day.alerts,
    })) || [];

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-slate-900">Dashboard</h1>
        <p className="text-slate-500">
          Resumen de actividad y monitoreo en tiempo real.
        </p>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {statCards.map((stat) => {
          const isAlerts = stat.name === "Alertas de Fraude";
          // Logic for color:
          // Normal: Up = Green, Down = Red
          // Alerts: Up = Red, Down = Green
          let trendColor = "text-slate-500";
          const TrendIcon = stat.trend === "up" ? ArrowUpRight : ArrowDownRight;

          if (isAlerts) {
            // For alerts, trend "up" (more alerts) is bad (red), "down" is good (green)
            // But the 'trend' property in stat object is calculated based on value change.
            // If change > 0, trend is 'up'.
            // If alerts increased (trend up), we want red.
            trendColor =
              stat.trend === "up" ? "text-red-500" : "text-green-500";
          } else {
            // For others, trend "up" is good (green)
            trendColor =
              stat.trend === "up" ? "text-green-500" : "text-red-500";
          }

          return (
            <div
              key={stat.name}
              className="rounded-xl bg-white p-6 shadow-sm border border-slate-100"
            >
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm font-medium text-slate-500">
                    {stat.name}
                  </p>
                  <p className="mt-2 text-3xl font-bold text-slate-900">
                    {stat.value}
                  </p>
                </div>
                <div className={`rounded-full p-3 ${stat.bg}`}>
                  <stat.icon className={`h-6 w-6 ${stat.color}`} />
                </div>
              </div>
              <div className="mt-4 flex items-center text-sm">
                <TrendIcon className={`mr-1 h-4 w-4 ${trendColor}`} />
                <span className={trendColor}>{stat.change}</span>
                <span className="ml-2 text-slate-400">vs mes anterior</span>
              </div>
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        {/* Chart */}
        <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
          <h3 className="text-lg font-semibold text-slate-900 mb-4">
            Actividad Semanal
          </h3>
          <div className="h-80 w-full">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                data={chartData}
                margin={{
                  top: 5,
                  right: 30,
                  left: 20,
                  bottom: 5,
                }}
              >
                <CartesianGrid strokeDasharray="3 3" vertical={false} />
                <XAxis
                  dataKey="name"
                  axisLine={false}
                  tickLine={false}
                  tick={{ fill: "#64748b" }}
                />
                <YAxis
                  axisLine={false}
                  tickLine={false}
                  tick={{ fill: "#64748b" }}
                />
                <Tooltip
                  contentStyle={{
                    backgroundColor: "#fff",
                    borderRadius: "8px",
                    border: "1px solid #e2e8f0",
                    boxShadow: "0 4px 6px -1px rgb(0 0 0 / 0.1)",
                  }}
                />
                <Legend />
                <Bar
                  dataKey="transacciones"
                  name="Transacciones"
                  fill="#3b82f6"
                  radius={[4, 4, 0, 0]}
                />
                <Bar
                  dataKey="fraudes"
                  name="Alertas"
                  fill="#ef4444"
                  radius={[4, 4, 0, 0]}
                />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Recent Alerts */}
        <div className="rounded-xl bg-white p-6 shadow-sm border border-slate-100">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-slate-900">
              Alertas Recientes
            </h3>
            <Link
              href="/alerts"
              className="text-sm font-medium text-blue-600 hover:text-blue-500"
            >
              Ver todas
            </Link>
          </div>
          <div className="space-y-4">
            {stats?.recent_alerts && stats.recent_alerts.length > 0 ? (
              stats.recent_alerts.map((alert) => (
                <div
                  key={alert.id}
                  className="flex items-start space-x-4 p-3 rounded-lg hover:bg-slate-50 transition-colors"
                >
                  <div
                    className={`rounded-full p-2 ${
                      alert.severity === "high"
                        ? "bg-red-100 text-red-600"
                        : alert.severity === "medium"
                          ? "bg-yellow-100 text-yellow-600"
                          : "bg-blue-100 text-blue-600"
                    }`}
                  >
                    <AlertTriangle className="h-4 w-4" />
                  </div>
                  <div className="flex-1">
                    <p className="text-sm font-medium text-slate-900">
                      {alert.description}
                    </p>
                    <p className="text-xs text-slate-500">
                      Tx: {alert.transaction_id} •{" "}
                      {new Date(alert.created_at).toLocaleDateString()}
                    </p>
                  </div>
                  <div className="text-right">
                    <Link
                      href={`/alerts/${alert.id}`}
                      className="text-xs font-medium text-blue-600 hover:underline"
                    >
                      Ver
                    </Link>
                  </div>
                </div>
              ))
            ) : (
              <p className="text-sm text-slate-500 text-center py-4">
                No hay alertas recientes.
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
