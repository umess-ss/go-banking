"use client";

import { useEffect, useState } from "react";
import Container from "@/components/shared/Container";
import ProtectedRoute from "@/components/shared/ProtectedRoute";
import { getAccounts } from "@/services/account.service";
import type { Account } from "@/types/account";

export default function AccountsPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
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

    loadAccounts();
  }, []);

  return (
    <ProtectedRoute>
      <Container>
        <div className="space-y-6">
          <div className="rounded-2xl border bg-white p-8 shadow-sm">
            <h1 className="text-3xl font-bold">Accounts</h1>
            <p className="mt-2 text-gray-600">
              View your bank accounts connected to your profile.
            </p>
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
                No accounts found. Account creation will be added in the next phase.
              </p>
            </div>
          )}

          {!loading && !error && accounts.length > 0 && (
            <div className="grid gap-4 md:grid-cols-2">
              {accounts.map((account) => (
                <div
                  key={account.id}
                  className="rounded-2xl border bg-white p-6 shadow-sm"
                >
                  <div className="flex items-start justify-between gap-4">
                    <div>
                      <p className="text-sm text-gray-500">
                        {account.account_type || "Bank Account"}
                      </p>
                      <h2 className="mt-1 text-xl font-semibold">
                        {account.account_number || account.id}
                      </h2>
                    </div>

                    <span className="rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-700">
                      {account.currency || "NPR"}
                    </span>
                  </div>

                  <p className="mt-6 text-3xl font-bold">
                    {(account.balance || 0).toLocaleString()}
                  </p>
                </div>
              ))}
            </div>
          )}
        </div>
      </Container>
    </ProtectedRoute>
  );
}