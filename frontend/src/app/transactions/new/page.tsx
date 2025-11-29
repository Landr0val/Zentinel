"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import { ArrowLeft, Save, Loader2 } from "lucide-react";
import PageHeader from "../../../components/PageHeader";
import { api } from "../../../lib/api";
import { CreateTransactionRequest, Account } from "../../../types";

export default function NewTransactionPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [accounts, setAccounts] = useState<Account[]>([]);

  const [formData, setFormData] = useState<CreateTransactionRequest>({
    account_id: "",
    amount: 0,
    currency: "USD",
    operation_type: "purchase",
    channel: "web",
    merchant: "",
    country: "",
    city: "",
  });

  useEffect(() => {
    const fetchAccounts = async () => {
      try {
        const response = await api.accounts.list({ page_size: 100 });
        setAccounts(response.data);
        if (response.data.length > 0) {
          setFormData((prev) => ({ ...prev, account_id: response.data[0].id }));
        }
      } catch (err) {
        console.error("Failed to fetch accounts", err);
      }
    };
    fetchAccounts();
  }, []);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: name === "amount" ? parseFloat(value) || 0 : value,
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    try {
      await api.transactions.create(formData);
      router.push("/transactions");
    } catch (err) {
      const message =
        err instanceof Error
          ? err.message
          : "Error al registrar la transacción";
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto animate-in fade-in duration-500">
      <div className="mb-6">
        <Link
          href="/transactions"
          className="inline-flex items-center text-sm text-muted-foreground hover:text-foreground mb-4 transition-colors"
        >
          <ArrowLeft className="mr-1 h-4 w-4" />
          Volver a Transacciones
        </Link>
        <PageHeader
          title="Registrar Nueva Transacción"
          description="Complete la información para registrar una nueva transacción en el sistema."
        />
      </div>

      {error && (
        <div className="mb-6 rounded-xl bg-destructive/10 p-4 border border-destructive/20">
          <div className="flex">
            <div className="ml-3">
              <h3 className="text-sm font-medium text-destructive">Error</h3>
              <div className="mt-2 text-sm text-destructive/90">
                <p>{error}</p>
              </div>
            </div>
          </div>
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="bg-card shadow-zen rounded-3xl border border-border overflow-hidden">
          <div className="px-6 py-8 sm:p-10">
            <div className="grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
              <div className="sm:col-span-4">
                <label
                  htmlFor="account_id"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Cuenta
                </label>
                <div className="mt-2">
                  <select
                    id="account_id"
                    name="account_id"
                    required
                    value={formData.account_id}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  >
                    <option value="" disabled>
                      Seleccione una cuenta
                    </option>
                    {accounts.map((account) => (
                      <option key={account.id} value={account.id}>
                        {account.account_number} - {account.client_name} (
                        {account.currency})
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="sm:col-span-2">
                <label
                  htmlFor="operation_type"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Tipo de Operación
                </label>
                <div className="mt-2">
                  <select
                    id="operation_type"
                    name="operation_type"
                    value={formData.operation_type}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  >
                    <option value="purchase">Compra</option>
                    <option value="deposit">Depósito</option>
                    <option value="withdrawal">Retiro</option>
                    <option value="transfer">Transferencia</option>
                  </select>
                </div>
              </div>

              <div className="sm:col-span-2">
                <label
                  htmlFor="amount"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Monto
                </label>
                <div className="mt-2">
                  <input
                    type="number"
                    name="amount"
                    id="amount"
                    required
                    min="0.01"
                    step="0.01"
                    value={formData.amount}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm placeholder:text-muted-foreground focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  />
                </div>
              </div>

              <div className="sm:col-span-2">
                <label
                  htmlFor="currency"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Moneda
                </label>
                <div className="mt-2">
                  <select
                    id="currency"
                    name="currency"
                    value={formData.currency}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  >
                    <option value="USD">USD</option>
                    <option value="EUR">EUR</option>
                    <option value="PEN">PEN</option>
                  </select>
                </div>
              </div>

              <div className="sm:col-span-2">
                <label
                  htmlFor="channel"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Canal
                </label>
                <div className="mt-2">
                  <select
                    id="channel"
                    name="channel"
                    value={formData.channel}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  >
                    <option value="mobile">Banca Móvil</option>
                    <option value="web">Banca por Internet</option>
                    <option value="atm">Cajero Automático</option>
                    <option value="branch">Agencia</option>
                  </select>
                </div>
              </div>

              <div className="sm:col-span-3">
                <label
                  htmlFor="merchant"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Comercio / Destinatario
                </label>
                <div className="mt-2">
                  <input
                    type="text"
                    name="merchant"
                    id="merchant"
                    value={formData.merchant}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm placeholder:text-muted-foreground focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  />
                </div>
              </div>

              <div className="sm:col-span-3">
                <label
                  htmlFor="country"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  País
                </label>
                <div className="mt-2">
                  <input
                    type="text"
                    name="country"
                    id="country"
                    value={formData.country}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm placeholder:text-muted-foreground focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  />
                </div>
              </div>

              <div className="sm:col-span-3">
                <label
                  htmlFor="city"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Ciudad
                </label>
                <div className="mt-2">
                  <input
                    type="text"
                    name="city"
                    id="city"
                    value={formData.city}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm placeholder:text-muted-foreground focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  />
                </div>
              </div>
            </div>
          </div>
          <div className="flex items-center justify-end gap-x-6 border-t border-border px-6 py-6 bg-secondary/30">
            <Link
              href="/transactions"
              className="text-sm font-semibold leading-6 text-foreground hover:text-foreground/80 transition-colors"
            >
              Cancelar
            </Link>
            <button
              type="submit"
              disabled={isLoading}
              className="inline-flex items-center rounded-xl bg-primary px-4 py-2.5 text-sm font-semibold text-primary-foreground shadow-sm hover:bg-primary/90 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary disabled:opacity-50 disabled:cursor-not-allowed transition-all"
            >
              {isLoading ? (
                <>
                  <Loader2 className="-ml-0.5 mr-2 h-4 w-4 animate-spin" />
                  Guardando...
                </>
              ) : (
                <>
                  <Save className="-ml-0.5 mr-2 h-4 w-4" aria-hidden="true" />
                  Registrar Transacción
                </>
              )}
            </button>
          </div>
        </div>
      </form>
    </div>
  );
}
