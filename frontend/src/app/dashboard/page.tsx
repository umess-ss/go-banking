"use client";

import { useEffect, useMemo, useState } from "react";
import Container from "@/components/shared/Container";
import ProtectedRoute from "@/components/shared/ProtectedRoute";
import { getAccounts } from "@/services/account.service";
import type { Account } from "@/types/account";

export default function DashboardPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function loadDashboard() {
      try {
        setError("");
        const data = await getAccounts();
        setAccounts(data);
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Failed to load dashboard"
        );
      } finally {
        setLoading(false);
      }
    }

    loadDashboard();
  }, []);

  const totalBalance = useMemo(() => {
    return accounts.reduce((sum, account) => sum + Number(account.balance || 0), 0);
  }, [accounts]);

  const savingsCount = accounts.filter(
    (account) => account.account_type === "savings"
  ).length;

  const checkingCount = accounts.filter(
    (account) => account.account_type === "checking"
  ).length;

  const currentCount = accounts.filter(
    (account) => account.account_type === "current"
  ).length;

  const recentAccounts = accounts.slice(0, 4);

  return (
    <ProtectedRoute>
      <Container>
        <div className="space-y-6">
          <div className="rounded-2xl border bg-white p-8 shadow-sm">
            <h1 className="text-3xl font-bold">Dashboard</h1>
            <p className="mt-2 text-gray-600">
              Overview of your banking accounts and balances.
            </p>
          </div>

          {loading && (
            <div className="rounded-2xl border bg-white p-6 shadow-sm">
              <p className="text-gray-600">Loading dashboard...</p>
            </div>
          )}

          {error && (
            <div className="rounded-2xl border border-red-200 bg-red-50 p-6">
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}

          {!loading && !error && (
            <>
              <div className="grid gap-4 md:grid-cols-3">
                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Total Balance</p>
                  <h2 className="mt-2 text-3xl font-bold">
                    NPR {totalBalance.toLocaleString()}
                  </h2>
                </div>

                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Total Accounts</p>
                  <h2 className="mt-2 text-3xl font-bold">
                    {accounts.length}
                  </h2>
                </div>

                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Currency</p>
                  <h2 className="mt-2 text-3xl font-bold">NPR</h2>
                </div>
              </div>

              <div className="grid gap-4 md:grid-cols-3">
                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Savings Accounts</p>
                  <h2 className="mt-2 text-2xl font-bold">{savingsCount}</h2>
                </div>

                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Checking Accounts</p>
                  <h2 className="mt-2 text-2xl font-bold">{checkingCount}</h2>
                </div>

                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Current Accounts</p>
                  <h2 className="mt-2 text-2xl font-bold">{currentCount}</h2>
                </div>
              </div>

              <div className="rounded-2xl border bg-white p-6 shadow-sm">
                <div className="flex items-center justify-between">
                  <div>
                    <h2 className="text-xl font-semibold">Recent Accounts</h2>
                    <p className="mt-1 text-sm text-gray-600">
                      Latest accounts connected to your profile.
                    </p>
                  </div>
                </div>

                {recentAccounts.length === 0 ? (
                  <p className="mt-6 text-gray-600">
                    No accounts found. Create your first account from the
                    Accounts page.
                  </p>
                ) : (
                  <div className="mt-6 divide-y">
                    {recentAccounts.map((account) => (
                      <div
                        key={account.id}
                        className="flex items-center justify-between py-4"
                      >
                        <div>
                          <p className="font-semibold">{account.name}</p>
                          <p className="mt-1 text-sm text-gray-500">
                            {account.account_number}
                          </p>
                        </div>

                        <div className="text-right">
                          <p className="text-sm capitalize text-gray-500">
                            {account.account_type}
                          </p>
                          <p className="font-semibold">
                            NPR {Number(account.balance || 0).toLocaleString()}
                          </p>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </>
          )}
        </div>
      </Container>
    </ProtectedRoute>
  );
}