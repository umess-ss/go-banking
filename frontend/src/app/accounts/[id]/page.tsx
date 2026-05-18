"use client";

import Link from "next/link";
import { useParams } from "next/navigation";
import { useEffect, useState } from "react";
import Container from "@/components/shared/Container";
import ProtectedRoute from "@/components/shared/ProtectedRoute";
import { getAccountById } from "@/services/account.service";
import { getTransactionsByAccountId } from "@/services/transaction.service";
import type { Account } from "@/types/account";
import type { Transaction } from "@/types/transaction";

export default function AccountDetailPage() {
  const params = useParams();
  const accountId = Number(params.id);

  const [account, setAccount] = useState<Account | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function loadAccountDetails() {
      try {
        setError("");

        if (!accountId || Number.isNaN(accountId)) {
          setError("Invalid account ID.");
          return;
        }

        const [accountData, transactionData] = await Promise.all([
          getAccountById(accountId),
          getTransactionsByAccountId(accountId),
        ]);

        setAccount(accountData);
        setTransactions(transactionData);
      } catch (err) {
        setError(
          err instanceof Error ? err.message : "Failed to load account details"
        );
      } finally {
        setLoading(false);
      }
    }

    loadAccountDetails();
  }, [accountId]);

  return (
    <ProtectedRoute>
      <Container>
        <div className="space-y-6">
          <div className="rounded-2xl border bg-white p-8 shadow-sm">
            <Link href="/accounts" className="text-sm text-gray-500">
              ← Back to accounts
            </Link>

            <h1 className="mt-4 text-3xl font-bold">Account Details</h1>
            <p className="mt-2 text-gray-600">
              View account information and account-specific transactions.
            </p>
          </div>

          {loading && (
            <div className="rounded-2xl border bg-white p-6 shadow-sm">
              <p className="text-gray-600">Loading account details...</p>
            </div>
          )}

          {error && (
            <div className="rounded-2xl border border-red-200 bg-red-50 p-6">
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}

          {!loading && !error && account && (
            <>
              <div className="grid gap-4 md:grid-cols-3">
                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Account Name</p>
                  <h2 className="mt-2 text-2xl font-bold">{account.name}</h2>
                </div>

                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Account Type</p>
                  <h2 className="mt-2 text-2xl font-bold capitalize">
                    {account.account_type}
                  </h2>
                </div>

                <div className="rounded-2xl border bg-white p-6 shadow-sm">
                  <p className="text-sm text-gray-500">Balance</p>
                  <h2 className="mt-2 text-2xl font-bold">
                    {account.currency}{" "}
                    {Number(account.balance || 0).toLocaleString()}
                  </h2>
                </div>
              </div>

              <div className="rounded-2xl border bg-white p-6 shadow-sm">
                <p className="text-sm text-gray-500">Account Number</p>
                <h2 className="mt-2 text-xl font-semibold">
                  {account.account_number}
                </h2>
              </div>

              <div className="rounded-2xl border bg-white shadow-sm">
                <div className="border-b p-6">
                  <h2 className="text-xl font-semibold">
                    Account Transactions
                  </h2>
                  <p className="mt-1 text-sm text-gray-600">
                    Transactions related to this account.
                  </p>
                </div>

                {transactions.length === 0 ? (
                  <p className="p-6 text-gray-600">
                    No transactions found for this account.
                  </p>
                ) : (
                  <div className="divide-y">
                    {transactions.map((transaction) => (
                      <div
                        key={transaction.id}
                        className="flex items-center justify-between gap-4 p-6"
                      >
                        <div>
                          <p className="font-semibold capitalize">
                            {transaction.type}
                          </p>
                          <p className="mt-1 text-sm text-gray-500">
                            Ref: {transaction.reference_number}
                          </p>
                          <p className="mt-1 text-xs text-gray-400">
                            {transaction.created_at
                              ? new Date(transaction.created_at).toLocaleString()
                              : "No date"}
                          </p>
                        </div>

                        <div className="text-right">
                          <p className="text-lg font-bold">
                            NPR{" "}
                            {Number(transaction.amount || 0).toLocaleString()}
                          </p>
                          <span className="mt-1 inline-block rounded-full bg-gray-100 px-3 py-1 text-xs capitalize text-gray-700">
                            {transaction.status}
                          </span>
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
