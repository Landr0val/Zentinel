"use client";

import {
  Activity,
  DollarSign,
  Users,
  Loader2,
  ShieldAlert,
  TrendingUp,
  TrendingDown,
} from "lucide-react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { useQuery } from "@tanstack/react-query";
import { api } from "../lib/api";
import Link from "next/link";
import { cn } from "../lib/utils";

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
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="rounded-xl border border-destructive/20 bg-destructive/5 p-6 text-destructive">
        <h3 className="text-lg font-medium">Error al cargar datos</h3>
        <p className="text-sm opacity-80">
          No se pudo conectar con el servidor.
        </p>
      </div>
    );
  }

  const stats = statsData?.data;

  const statCards = [
    {
      name: "Volumen Total",
      value:
        stats?.total_volume !== undefined
          ? `$${stats.total_volume.toLocaleString()}`
          : "-",
      change: stats?.volume_change_percentage
        ? `${stats.volume_change_percentage > 0 ? "+" : ""}${stats.volume_change_percentage}%`
        : "-",
      trend:
        stats?.volume_change_percentage !== undefined &&
        stats.volume_change_percentage >= 0
          ? "up"
          : "down",
      icon: DollarSign,
    },
    {
      name: "Transacciones",
      value:
        stats?.total_transactions !== undefined
          ? stats.total_transactions.toLocaleString()
          : "-",
      change: stats?.transaction_change_percentage
        ? `${stats.transaction_change_percentage > 0 ? "+" : ""}${stats.transaction_change_percentage}%`
        : "-",
      trend:
        stats?.transaction_change_percentage !== undefined &&
        stats.transaction_change_percentage >= 0
          ? "up"
          : "down",
      icon: Activity,
    },
    {
      name: "Alertas de Fraude",
      value:
        stats?.flagged_count !== undefined
          ? stats.flagged_count.toLocaleString()
          : "-",
      change: stats?.alert_change_percentage
        ? `${stats.alert_change_percentage > 0 ? "+" : ""}${stats.alert_change_percentage}%`
        : "-",
      trend:
        stats?.alert_change_percentage !== undefined &&
        stats.alert_change_percentage <= 0
          ? "up"
          : "down",
      icon: ShieldAlert,
    },
    {
      name: "Volumen de Riesgo",
      value:
        stats?.flagged_volume !== undefined
          ? `$${stats.flagged_volume.toLocaleString()}`
          : "-",
      change: "-",
      trend: "down",
      icon: ShieldAlert,
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
    <div className="space-y-8 animate-in fade-in duration-500">
      <div className="flex items-end justify-between">
        <div>
          <h1 className="text-3xl font-light tracking-tight text-foreground">
            Dashboard
          </h1>
          <p className="mt-2 text-muted-foreground">Monitoreo en tiempo real</p>
        </div>
        <div className="text-sm text-muted-foreground border border-border px-3 py-1 rounded-full hidden sm:block">
          {new Date().toLocaleDateString("es-ES", {
            weekday: "long",
            year: "numeric",
            month: "long",
            day: "numeric",
          })}
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4">
        {statCards.map((stat) => {
          const TrendIcon = stat.trend === "up" ? TrendingUp : TrendingDown;
          return (
            <div
              key={stat.name}
              className="group relative overflow-hidden rounded-2xl border border-border bg-card p-6 shadow-zen transition-all hover:shadow-lg"
            >
              <div className="flex items-center justify-between">
                <div className="rounded-xl bg-secondary/50 p-2.5 text-foreground transition-colors group-hover:bg-primary group-hover:text-primary-foreground">
                  <stat.icon className="h-5 w-5" />
                </div>
                <div
                  className={cn(
                    "flex items-center gap-1 text-xs font-medium px-2 py-1 rounded-full border",
                    stat.trend === "up"
                      ? "border-border bg-secondary/30 text-foreground"
                      : "border-border bg-secondary/30 text-muted-foreground",
                  )}
                >
                  <TrendIcon className="h-3 w-3" />
                  <span>{stat.change}</span>
                </div>
              </div>
              <div className="mt-4">
                <p className="text-sm font-medium text-muted-foreground">
                  {stat.name}
                </p>
                <p className="mt-1 text-2xl font-semibold tracking-tight text-foreground">
                  {stat.value}
                </p>
              </div>
            </div>
          );
        })}
      </div>

      <div className="grid grid-cols-1 gap-8 lg:grid-cols-3">
        {/* Chart */}
        <div className="lg:col-span-2 rounded-3xl border border-border bg-card p-8 shadow-zen">
          <div className="mb-6 flex items-center justify-between">
            <h3 className="text-lg font-medium text-foreground">
              Actividad Semanal
            </h3>
          </div>
          <div className="h-[350px] w-full">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart
                data={chartData}
                margin={{
                  top: 10,
                  right: 10,
                  left: -20,
                  bottom: 0,
                }}
                barGap={8}
              >
                <CartesianGrid
                  strokeDasharray="3 3"
                  vertical={false}
                  stroke="var(--border)"
                  opacity={0.4}
                />
                <XAxis
                  dataKey="name"
                  axisLine={false}
                  tickLine={false}
                  tick={{ fill: "var(--muted-foreground)", fontSize: 12 }}
                  dy={10}
                />
                <YAxis
                  axisLine={false}
                  tickLine={false}
                  tick={{ fill: "var(--muted-foreground)", fontSize: 12 }}
                />
                <Tooltip
                  cursor={{ fill: "var(--secondary)", opacity: 0.4 }}
                  contentStyle={{
                    backgroundColor: "var(--popover)",
                    borderColor: "var(--border)",
                    borderRadius: "12px",
                    boxShadow: "var(--shadow-zen)",
                    color: "var(--popover-foreground)",
                  }}
                  itemStyle={{ color: "var(--foreground)" }}
                />
                <Bar
                  dataKey="transacciones"
                  name="Transacciones"
                  fill="var(--foreground)"
                  radius={[6, 6, 6, 6]}
                  barSize={32}
                  animationDuration={1000}
                />
                <Bar
                  dataKey="fraudes"
                  name="Alertas"
                  fill="var(--muted-foreground)"
                  radius={[6, 6, 6, 6]}
                  barSize={32}
                  animationDuration={1000}
                  opacity={0.3}
                />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Recent Alerts */}
        <div className="rounded-3xl border border-border bg-card p-8 shadow-zen flex flex-col">
          <div className="flex items-center justify-between mb-6">
            <h3 className="text-lg font-medium text-foreground">
              Alertas Recientes
            </h3>
            <Link
              href="/alerts"
              className="text-xs font-medium text-muted-foreground hover:text-foreground transition-colors border border-border px-3 py-1 rounded-full hover:bg-secondary"
            >
              Ver todas
            </Link>
          </div>
          <div className="flex-1 overflow-y-auto pr-2 space-y-4">
            {stats?.recent_alerts && stats.recent_alerts.length > 0 ? (
              stats.recent_alerts.map((alert) => (
                <div
                  key={alert.id}
                  className="group flex flex-col gap-3 rounded-2xl border border-border bg-secondary/10 p-4 transition-all hover:bg-secondary/40 hover:border-secondary-foreground/10"
                >
                  <div className="flex items-start justify-between">
                    <div className="flex items-center gap-2">
                      <div
                        className={cn(
                          "h-2 w-2 rounded-full",
                          alert.severity === "high"
                            ? "bg-foreground"
                            : alert.severity === "medium"
                              ? "bg-muted-foreground"
                              : "bg-border",
                        )}
                      />
                      <span className="text-xs font-medium text-muted-foreground uppercase tracking-wider">
                        {alert.severity === "high"
                          ? "Crítico"
                          : alert.severity === "medium"
                            ? "Medio"
                            : "Bajo"}
                      </span>
                    </div>
                    <span className="text-[10px] text-muted-foreground">
                      {new Date(alert.created_at).toLocaleDateString()}
                    </span>
                  </div>

                  <div>
                    <p className="text-sm font-medium text-foreground leading-snug">
                      {alert.description}
                    </p>
                    <p className="mt-1 text-xs text-muted-foreground font-mono">
                      Ref: {alert.transaction_code}
                    </p>
                  </div>

                  <div className="pt-2 mt-1 border-t border-border/50 flex justify-end">
                    <Link
                      href={`/alerts/${alert.id}`}
                      className="text-xs font-medium text-foreground hover:underline decoration-1 underline-offset-4"
                    >
                      Revisar detalles &rarr;
                    </Link>
                  </div>
                </div>
              ))
            ) : (
              <div className="flex flex-col items-center justify-center h-40 text-muted-foreground">
                <ShieldAlert className="h-8 w-8 mb-2 opacity-20" />
                <p className="text-sm">Sin alertas recientes</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
