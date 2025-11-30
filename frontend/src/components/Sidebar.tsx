"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import {
  LayoutDashboard,
  Users,
  CreditCard,
  Activity,
  AlertTriangle,
  ShieldCheck,
  LogOut,
} from "lucide-react";
import { cn } from "../lib/utils";

const navigation = [
  { name: "Dashboard", href: "/", icon: LayoutDashboard },
  { name: "Clientes", href: "/clients", icon: Users },
  { name: "Cuentas", href: "/accounts", icon: CreditCard },
  { name: "Transacciones", href: "/transactions", icon: Activity },
  { name: "Alertas", href: "/alerts", icon: AlertTriangle },
];

export default function Sidebar() {
  const pathname = usePathname();
  const router = useRouter();

  const handleLogout = () => {
    localStorage.removeItem("token");
    router.push("/login");
  };

  return (
    <aside className="flex h-full w-64 flex-col bg-card border-r border-border fixed left-0 top-0 bottom-0 z-50 transition-colors duration-300">
      <div className="flex h-24 items-center px-6">
        <div className="flex items-center gap-3 group">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-primary text-primary-foreground shadow-zen transition-transform group-hover:scale-105">
            <ShieldCheck className="h-5 w-5" />
          </div>
          <span className="text-lg font-bold tracking-tight text-foreground">
            ZENTINEL
          </span>
        </div>
      </div>

      <nav className="flex-1 space-y-2 px-4 py-6 overflow-y-auto">
        {navigation.map((item) => {
          const isActive =
            pathname === item.href ||
            (item.href !== "/" && pathname.startsWith(item.href));

          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "group flex items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium transition-all duration-200 ease-out",
                isActive
                  ? "bg-primary text-primary-foreground shadow-zen translate-x-1"
                  : "text-muted-foreground hover:bg-secondary hover:text-foreground hover:translate-x-1",
              )}
            >
              <item.icon
                className={cn(
                  "h-5 w-5 shrink-0 transition-colors",
                  isActive
                    ? "text-primary-foreground"
                    : "text-muted-foreground group-hover:text-foreground",
                )}
                aria-hidden="true"
              />
              {item.name}
            </Link>
          );
        })}
      </nav>

      <div className="p-6 border-t border-border">
        <button
          onClick={handleLogout}
          className="flex w-full items-center gap-3 rounded-xl px-4 py-3 text-sm font-medium text-muted-foreground hover:bg-destructive/10 hover:text-destructive transition-all duration-200 ease-out mb-4"
        >
          <LogOut className="h-5 w-5 shrink-0" />
          Cerrar Sesión
        </button>
        <div className="flex items-center justify-center">
          <p className="text-[10px] font-medium text-muted-foreground/50 uppercase tracking-widest">
            v1.0.0
          </p>
        </div>
      </div>
    </aside>
  );
}
