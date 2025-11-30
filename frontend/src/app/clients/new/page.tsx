"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { ArrowLeft, Save, Loader2 } from "lucide-react";
import PageHeader from "../../../components/PageHeader";
import { api } from "../../../lib/api";
import { CreateClientRequest } from "../../../types";

export default function NewClientPage() {
  const router = useRouter();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState<CreateClientRequest>({
    document_type: "DNI",
    document_number: "",
    full_name: "",
    email: "",
    phone: "",
    initial_account: {
      account_number: "",
      account_type: "savings",
      currency: "USD",
    },
  });

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
  };

  const handleAccountChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>,
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      initial_account: {
        ...prev.initial_account!,
        [name]: value,
      },
    }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    try {
      await api.clients.create(formData);
      router.push("/clients");
    } catch (err) {
      const message =
        err instanceof Error ? err.message : "Error al crear el cliente";
      setError(message);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto animate-in fade-in duration-500">
      <div className="mb-6">
        <Link
          href="/clients"
          className="inline-flex items-center text-sm text-muted-foreground hover:text-foreground mb-4 transition-colors"
        >
          <ArrowLeft className="mr-1 h-4 w-4" />
          Volver a Clientes
        </Link>
        <PageHeader
          title="Registrar Nuevo Cliente"
          description="Complete la información para dar de alta un nuevo cliente y su cuenta inicial."
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
              <div className="sm:col-span-2">
                <label
                  htmlFor="document_type"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Tipo de Documento
                </label>
                <div className="mt-2">
                  <select
                    id="document_type"
                    name="document_type"
                    value={formData.document_type}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  >
                    <option value="DNI">DNI</option>
                    <option value="PASSPORT">Pasaporte</option>
                    <option value="RUC">RUC</option>
                  </select>
                </div>
              </div>

              <div className="sm:col-span-4">
                <label
                  htmlFor="document_number"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Número de Documento
                </label>
                <div className="mt-2">
                  <input
                    type="text"
                    name="document_number"
                    id="document_number"
                    required
                    value={formData.document_number}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm placeholder:text-muted-foreground focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  />
                </div>
              </div>

              <div className="sm:col-span-6">
                <label
                  htmlFor="full_name"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Nombre Completo
                </label>
                <div className="mt-2">
                  <input
                    type="text"
                    name="full_name"
                    id="full_name"
                    required
                    value={formData.full_name}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm placeholder:text-muted-foreground focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  />
                </div>
              </div>

              <div className="sm:col-span-3">
                <label
                  htmlFor="email"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Correo Electrónico
                </label>
                <div className="mt-2">
                  <input
                    id="email"
                    name="email"
                    type="email"
                    required
                    value={formData.email}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm placeholder:text-muted-foreground focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  />
                </div>
              </div>

              <div className="sm:col-span-3">
                <label
                  htmlFor="phone"
                  className="block text-sm font-medium leading-6 text-foreground"
                >
                  Teléfono
                </label>
                <div className="mt-2">
                  <input
                    type="text"
                    name="phone"
                    id="phone"
                    required
                    value={formData.phone}
                    onChange={handleChange}
                    className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm placeholder:text-muted-foreground focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                  />
                </div>
              </div>
            </div>

            <div className="mt-10 border-t border-border pt-10">
              <h3 className="text-base font-semibold leading-7 text-foreground mb-4">
                Cuenta Inicial
              </h3>
              <div className="grid grid-cols-1 gap-x-6 gap-y-8 sm:grid-cols-6">
                <div className="sm:col-span-3">
                  <label
                    htmlFor="account_number"
                    className="block text-sm font-medium leading-6 text-foreground"
                  >
                    Número de Cuenta
                  </label>
                  <div className="mt-2">
                    <input
                      type="text"
                      name="account_number"
                      id="account_number"
                      required
                      value={formData.initial_account?.account_number}
                      onChange={handleAccountChange}
                      className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm placeholder:text-muted-foreground focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                    />
                  </div>
                </div>

                <div className="sm:col-span-2">
                  <label
                    htmlFor="account_type"
                    className="block text-sm font-medium leading-6 text-foreground"
                  >
                    Tipo de Cuenta
                  </label>
                  <div className="mt-2">
                    <select
                      id="account_type"
                      name="account_type"
                      value={formData.initial_account?.account_type}
                      onChange={handleAccountChange}
                      className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                    >
                      <option value="savings">Ahorros</option>
                      <option value="checking">Corriente</option>
                    </select>
                  </div>
                </div>

                <div className="sm:col-span-1">
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
                      value={formData.initial_account?.currency}
                      onChange={handleAccountChange}
                      className="block w-full rounded-xl border border-input bg-background py-2.5 text-foreground shadow-sm focus:border-ring focus:ring-1 focus:ring-ring sm:text-sm sm:leading-6 outline-none transition-all"
                    >
                      <option value="USD">USD</option>
                      <option value="EUR">EUR</option>
                      <option value="COP">COP</option>
                    </select>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div className="flex items-center justify-end gap-x-6 border-t border-border px-6 py-6 bg-secondary/30">
            <Link
              href="/clients"
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
                  Guardar Cliente
                </>
              )}
            </button>
          </div>
        </div>
      </form>
    </div>
  );
}
