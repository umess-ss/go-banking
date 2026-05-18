"use client";

import Link from "next/link";
import { useEffect, useMemo, useState } from "react";
import Container from "@/components/shared/Container";
import ProtectedRoute from "@/components/shared/ProtectedRoute";
import { getAccounts } from "@/services/account.service";
import { getTransactions } from "@/services/transaction.service";
import type { Account } from "@/types/account";
import type { Transaction } from "@/types/transaction";
import { getAuthUser, setAuthUser } from "@/lib/auth";
import type { AuthUser } from "@/types/auth";
import { getCurrentUser } from "@/services/auth.service";
import Alert from "@/components/ui/Alert";
import EmptyState from "@/components/ui/EmptyState";
import LoadingState from "@/components/ui/LoadingState";

export default function DashboardPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [user, setUser] = useState<AuthUser | null>(null);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function loadDashboard() {
      try {
        setError("");
        const cachedUser = getAuthUser();
        setUser(cachedUser);

        const [accountsData, transactionsData, currentUserData] = await Promise.all([
          getAccounts(),
          getTransactions(),
          cachedUser ? Promise.resolve(null) : getCurrentUser(),
        ]);

        setAccounts(accountsData);
        setTransactions(transactionsData);

        if (currentUserData?.data) {
          setAuthUser(currentUserData.data);
          setUser(currentUserData.data);
        }
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Failed to load dashboard"
        );
      } finally {
        setLoading(false);
      }
    }

    queueMicrotask(() => {
      void loadDashboard();
    });
  }, []);

  const totalBalance = useMemo(() => {
    return accounts.reduce(
      (sum, account) => sum + Number(account.balance || 0),
      0
    );
  }, [accounts]);

  const totalDeposits = useMemo(() => {
    return transactions
      .filter((transaction) => transaction.type === "deposit")
      .reduce((sum, transaction) => sum + Number(transaction.amount || 0), 0);
  }, [transactions]);

  const totalWithdrawals = useMemo(() => {
    return transactions
      .filter((transaction) => transaction.type === "withdraw")
      .reduce((sum, transaction) => sum + Number(transaction.amount || 0), 0);
  }, [transactions]);

  const totalTransfers = useMemo(() => {
    return transactions
      .filter((transaction) => transaction.type === "transfer")
      .reduce((sum, transaction) => sum + Number(transaction.amount || 0), 0);
  }, [transactions]);

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
  const recentTransactions = transactions.slice(0, 5);

  return (
    <ProtectedRoute>
      <Container>
        <div className="space-y-6">
          <div className="rounded-2xl border bg-white p-8 shadow-sm">
            <h1 className="text-3xl font-bold">
              {user?.name ? `Welcome, ${user.name}` : "Welcome"}
            </h1>
            <p className="mt-2 text-gray-600">
              Overview of your banking accounts, balances, and transactions.
            </p>
          </div>

          {loading && <LoadingState text="Loading dashboard..." />}

          {error && <Alert variant="error">{error}</Alert>}

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
                  <h2 className="mt-2 text-3xl font-bold">{accounts.length}</h2>
                </div>

                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Total Transactions</p>
                  <h2 className="mt-2 text-3xl font-bold">
                    {transactions.length}
                  </h2>
                </div>
              </div>

              <div className="grid gap-4 md:grid-cols-3">
                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Total Deposits</p>
                  <h2 className="mt-2 text-2xl font-bold">
                    NPR {totalDeposits.toLocaleString()}
                  </h2>
                </div>

                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Total Withdrawals</p>
                  <h2 className="mt-2 text-2xl font-bold">
                    NPR {totalWithdrawals.toLocaleString()}
                  </h2>
                </div>

                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Total Transfers</p>
                  <h2 className="mt-2 text-2xl font-bold">
                    NPR {totalTransfers.toLocaleString()}
                  </h2>
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

              <div className="grid gap-6 lg:grid-cols-2">
                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <div className="flex items-center justify-between">
                    <div>
                      <h2 className="text-xl font-semibold">Recent Accounts</h2>
                      <p className="mt-1 text-sm text-gray-600">
                        Latest accounts connected to your profile.
                      </p>
                    </div>

                    <Link href="/accounts" className="text-sm font-medium">
                      View all
                    </Link>
                  </div>

                  {recentAccounts.length === 0 ? (
                    <EmptyState variant="plain" className="mt-6">
                      No accounts found. Create your first account.
                    </EmptyState>
                  ) : (
                    <div className="mt-6 divide-y">
                      {recentAccounts.map((account) => (
                        <Link
                          key={account.id}
                          href={`/accounts/${account.id}`}
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
                              NPR{" "}
                              {Number(account.balance || 0).toLocaleString()}
                            </p>
                          </div>
                        </Link>
                      ))}
                    </div>
                  )}
                </div>

                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <div className="flex items-center justify-between">
                    <div>
                      <h2 className="text-xl font-semibold">
                        Recent Transactions
                      </h2>
                      <p className="mt-1 text-sm text-gray-600">
                        Latest deposit, withdrawal, and transfer activity.
                      </p>
                    </div>

                    <Link href="/transactions" className="text-sm font-medium">
                      View all
                    </Link>
                  </div>

                  {recentTransactions.length === 0 ? (
                    <EmptyState variant="plain" className="mt-6">
                      No transactions found yet.
                    </EmptyState>
                  ) : (
                    <div className="mt-6 divide-y">
                      {recentTransactions.map((transaction) => (
                        <div
                          key={transaction.id}
                          className="flex items-center justify-between py-4"
                        >
                          <div>
                            <p className="font-semibold capitalize">
                              {transaction.type}
                            </p>
                            <p className="mt-1 text-sm text-gray-500">
                              Ref: {transaction.reference_number}
                            </p>
                          </div>

                          <div className="text-right">
                            <p className="font-semibold">
                              NPR{" "}
                              {Number(transaction.amount || 0).toLocaleString()}
                            </p>
                            <p className="mt-1 text-sm capitalize text-gray-500">
                              {transaction.status}
                            </p>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              </div>
            </>
          )}
        </div>
      </Container>
    </ProtectedRoute>
  );
}
