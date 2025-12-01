"use client";

import React, { createContext, useContext, useState, useEffect } from "react";

type Currency = "USD" | "EUR" | "COP";

interface CurrencyContextType {
  currency: Currency;
  setCurrency: (currency: Currency) => void;
  formatAmount: (amount: number, fromCurrency?: string) => string;
  convertAmount: (amount: number, fromCurrency?: string) => number;
}

const CurrencyContext = createContext<CurrencyContextType | undefined>(
  undefined,
);

// Tasas de cambio fijas para el ejemplo (base USD)
const EXCHANGE_RATES: Record<string, number> = {
  USD: 1,
  EUR: 0.92,
  COP: 3900,
};

export function CurrencyProvider({ children }: { children: React.ReactNode }) {
  const [currency, setCurrency] = useState<Currency>(() => {
    if (typeof window !== "undefined") {
      const saved = localStorage.getItem("preferredCurrency") as Currency;
      if (saved && EXCHANGE_RATES[saved]) {
        return saved;
      }
    }
    return "USD";
  });

  useEffect(() => {
    localStorage.setItem("preferredCurrency", currency);
  }, [currency]);

  const convertAmount = (amount: number, fromCurrency: string = "USD") => {
    // Normalizar moneda origen
    const from = fromCurrency.toUpperCase();

    // Si la moneda origen es la misma que la seleccionada, no convertir
    if (from === currency) return amount;

    // Convertir a USD primero (moneda base)
    const amountInUSD = amount / (EXCHANGE_RATES[from] || 1);

    // Convertir de USD a moneda destino
    return amountInUSD * EXCHANGE_RATES[currency];
  };

  const formatAmount = (amount: number, fromCurrency: string = "USD") => {
    const convertedAmount = convertAmount(amount, fromCurrency);

    return new Intl.NumberFormat("es-MX", {
      style: "currency",
      currency: currency,
      minimumFractionDigits: 2,
      maximumFractionDigits: 2,
    }).format(convertedAmount);
  };

  return (
    <CurrencyContext.Provider
      value={{ currency, setCurrency, formatAmount, convertAmount }}
    >
      {children}
    </CurrencyContext.Provider>
  );
}

export function useCurrency() {
  const context = useContext(CurrencyContext);
  if (context === undefined) {
    throw new Error("useCurrency must be used within a CurrencyProvider");
  }
  return context;
}
