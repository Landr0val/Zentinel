"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  LayoutDashboard,
  Users,
  CreditCard,
  Activity,
  AlertTriangle,
} from "lucide-react";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

const navigation = [
  { name: "Dashboard", href: "/", icon: LayoutDashboard },
  { name: "Clientes", href: "/clients", icon: Users },
  { name: "Cuentas", href: "/accounts", icon: CreditCard },
  { name: "Transacciones", href: "/transactions", icon: Activity },
  { name: "Alertas", href: "/alerts", icon: AlertTriangle },
];

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <div className="flex h-full w-64 flex-col bg-slate-900 text-white fixed left-0 top-0 bottom-0 z-50">
      <div className="flex h-16 items-center px-6 border-b border-slate-800">
        <div className="flex items-center gap-2 font-bold text-xl tracking-wider text-blue-500">
          <Activity className="h-6 w-6" />
          <span className="text-white">ZENTINEL</span>
        </div>
      </div>
      <nav className="flex-1 space-y-1 px-3 py-4 overflow-y-auto">
        {navigation.map((item) => {
          const isActive =
            pathname === item.href ||
            (item.href !== "/" && pathname.startsWith(item.href));
          return (
            <Link
              key={item.name}
              href={item.href}
              className={cn(
                "group flex items-center rounded-md px-3 py-2.5 text-sm font-medium transition-colors",
                isActive
                  ? "bg-blue-600 text-white shadow-md"
                  : "text-slate-400 hover:bg-slate-800 hover:text-white",
              )}
            >
              <item.icon
                className={cn(
                  "mr-3 h-5 w-5 shrink-0 transition-colors",
                  isActive
                    ? "text-white"
                    : "text-slate-400 group-hover:text-white",
                )}
                aria-hidden="true"
              />
              {item.name}
            </Link>
          );
        })}
      </nav>
      <div className="p-4 border-t border-slate-800">
        <p className="text-xs text-slate-500 text-center">v1.0.0</p>
      </div>
    </div>
  );
}
