"use client";

import Link from "next/link";
import { FormEvent, useEffect, useState } from "react";
import Container from "@/components/shared/Container";
import ProtectedRoute from "@/components/shared/ProtectedRoute";
import { createAccount, getAccounts } from "@/services/account.service";
import type { Account } from "@/types/account";

export default function AccountsPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [accountName, setAccountName] = useState("");
  const [accountType, setAccountType] = useState("savings");
  const [currency, setCurrency] = useState("NPR");

  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  async function loadAccounts() {
    try {
      setError("");
      const data = await getAccounts();
      setAccounts(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load accounts");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    queueMicrotask(() => {
      void loadAccounts();
    });
  }, []);

  async function handleCreateAccount(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setError("");
    setSuccess("");
    setCreating(true);

    try {
      await createAccount({
        name: accountName,
        currency,
      });


      setSuccess("Account created successfully.");
      setAccountName("");
      await loadAccounts();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create account");
    } finally {
      setCreating(false);
    }
  }

  return (
    <ProtectedRoute>
      <Container>
        <div className="space-y-6">
          <div className="rounded-2xl border bg-white p-8 shadow-sm">
            <h1 className="text-3xl font-bold">Accounts</h1>
            <p className="mt-2 text-gray-600">
              Create and view your bank accounts.
            </p>
          </div>

          <div className="rounded-2xl border bg-white p-6 shadow-sm">
            <h2 className="text-xl font-semibold">Create New Account</h2>

            <form
              onSubmit={handleCreateAccount}
              className="mt-5 grid gap-4 md:grid-cols-3"
            >
              <div>
              <label className="text-sm font-medium">Account Name</label>
              <input
                type="text"
                value={accountName}
                onChange={(event) => setAccountName(event.target.value)}
                placeholder="Name"
                className="mt-1 w-full rounded-lg border px-3 py-2 outline-none focus:ring-2 focus:ring-black"
                required
              />
            </div>
              <div>
                <label className="text-sm font-medium">Account Type</label>
                <select
                  value={accountType}
                  onChange={(event) => setAccountType(event.target.value)}
                  className="mt-1 w-full rounded-lg border px-3 py-2 outline-none focus:ring-2 focus:ring-black"
                >
                  <option value="savings">Savings</option>
                  <option value="checking">Checking</option>
                  <option value="current">Current</option>
                </select>
              </div>

              <div>
                <label className="text-sm font-medium">Currency</label>
                <select
                  value={currency}
                  onChange={(event) => setCurrency(event.target.value)}
                  className="mt-1 w-full rounded-lg border px-3 py-2 outline-none focus:ring-2 focus:ring-black"
                >
                  <option value="NPR">NPR</option>
                  <option value="USD">USD</option>
                </select>
              </div>

              <div className="flex items-end">
                <button
                  type="submit"
                  disabled={creating}
                  className="w-full rounded-lg bg-black px-4 py-2 font-medium text-white disabled:opacity-60"
                >
                  {creating ? "Creating..." : "Create Account"}
                </button>
              </div>
            </form>

            {success && (
              <p className="mt-4 rounded-lg bg-green-50 px-3 py-2 text-sm text-green-700">
                {success}
              </p>
            )}
          </div>

          {loading && (
            <div className="rounded-2xl border bg-white p-6 shadow-sm">
              <p className="text-gray-600">Loading accounts...</p>
            </div>
          )}

          {error && (
            <div className="rounded-2xl border border-red-200 bg-red-50 p-6">
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}

          {!loading && !error && accounts.length === 0 && (
            <div className="rounded-2xl border bg-white p-6 shadow-sm">
              <p className="text-gray-600">
                No accounts found. Create your first account above.
              </p>
            </div>
          )}

          {!loading && !error && accounts.length > 0 && (
            <div className="grid gap-4 md:grid-cols-2">
              {accounts.map((account) => (
                <Link
                    href={`/accounts/${account.id}`}
                    key={account.id}
                    className="block rounded-2xl border bg-white p-6 shadow-sm transition hover:-translate-y-0.5 hover:shadow-md"
                  >
                  <div className="flex items-start justify-between gap-4">
                    <div>
                    <p className="text-sm capitalize text-gray-500">
                      {account.account_type || "Bank Account"}
                    </p>
                      <h2 className="mt-1 text-xl font-semibold">
                        {account.name || "Unnamed Account"}
                      </h2>
                      <h2 className="mt-1 text-xl font-semibold">
                        {account.account_number || account.id}
                      </h2>
                    </div>

                    <span className="rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-700">
                      {account.currency || "NPR"}
                    </span>
                  </div>

                  <p className="mt-6 text-3xl font-bold">
                    {account.currency} {Number(account.balance || 0).toLocaleString()}
                  </p>
                </Link>
              ))}
            </div>
          )}
        </div>
      </Container>
    </ProtectedRoute>
  );
}
