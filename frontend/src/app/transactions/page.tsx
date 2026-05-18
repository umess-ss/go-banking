"use client";

import { FormEvent, useEffect, useState } from "react";
import Container from "@/components/shared/Container";
import ProtectedRoute from "@/components/shared/ProtectedRoute";
import { getAccounts } from "@/services/account.service";
import { parsePositiveAmount } from "@/lib/validation";
import {
  deposit,
  getTransactions,
  transfer,
  withdraw,
} from "@/services/transaction.service";
import type { Account } from "@/types/account";
import type { Transaction } from "@/types/transaction";

export default function TransactionsPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [transactions, setTransactions] = useState<Transaction[]>([]);

  const [selectedAccountId, setSelectedAccountId] = useState("");
  const [amount, setAmount] = useState("");
  const [actionType, setActionType] = useState<"deposit" | "withdraw">(
    "deposit"
  );

  const [fromAccountId, setFromAccountId] = useState("");
  const [toAccountId, setToAccountId] = useState("");
  const [transferAmount, setTransferAmount] = useState("");

  const [loading, setLoading] = useState(true);
  const [submittingAction, setSubmittingAction] = useState<
    "deposit" | "withdraw" | "transfer" | null
  >(null);

  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  async function loadData() {
    try {
      setError("");

      const [accountsData, transactionsData] = await Promise.all([
        getAccounts(),
        getTransactions(),
      ]);

      setAccounts(accountsData);
      setTransactions(transactionsData);

      if (!selectedAccountId && accountsData.length > 0) {
        setSelectedAccountId(String(accountsData[0].id));
      }

      if (!fromAccountId && accountsData.length > 0) {
        setFromAccountId(String(accountsData[0].id));
      }

      if (!toAccountId && accountsData.length > 1) {
        setToAccountId(String(accountsData[1].id));
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load data");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    queueMicrotask(() => {
      void loadData();
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function handleMoneyAction(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError("");
    setSuccess("");

    const amountResult = parsePositiveAmount(amount);

    if (!selectedAccountId) {
      setError("Please select an account.");
      return;
    }


    if (!amountResult.valid) {
  setError(amountResult.message);
  return;
}

const numericAmount = amountResult.amount;
const selectedAccount = accounts.find(
  (account) => String(account.id) === selectedAccountId
);

if (actionType === "withdraw" && selectedAccount) {
  const balance = Number(selectedAccount.balance || 0);

  if (numericAmount > balance) {
    setError("Withdrawal amount cannot exceed account balance.");
    return;
  }
}

    try {
      setSubmittingAction(actionType);

      if (actionType === "deposit") {
        await deposit(Number(selectedAccountId), numericAmount);
        setSuccess("Deposit completed successfully.");
      } else {
        await withdraw(Number(selectedAccountId), numericAmount);
        setSuccess("Withdrawal completed successfully.");
      }

      setAmount("");
      await loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : `${actionType} failed`);
    } finally {
      setSubmittingAction(null);
    }
  }

  async function handleTransfer(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    setError("");
    setSuccess("");

    const amountResult = parsePositiveAmount(transferAmount);

    if (!fromAccountId || !toAccountId) {
      setError("Please select both source and destination accounts.");
      return;
    }

    if (fromAccountId === toAccountId) {
      setError("Source and destination accounts cannot be the same.");
      return;
    }

    if (!amountResult.valid) {
  setError(amountResult.message);
  return;
}

const numericAmount = amountResult.amount;
const sourceAccount = accounts.find(
  (account) => String(account.id) === fromAccountId
);

if (sourceAccount) {
  const balance = Number(sourceAccount.balance || 0);

  if (numericAmount > balance) {
    setError("Transfer amount cannot exceed source account balance.");
    return;
  }
}

    try {
      setSubmittingAction("transfer");

      await transfer(Number(fromAccountId), Number(toAccountId), numericAmount);

      setSuccess("Transfer completed successfully.");
      setTransferAmount("");
      await loadData();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Transfer failed");
    } finally {
      setSubmittingAction(null);
    }
  }

  return (
    <ProtectedRoute>
      <Container>
        <div className="space-y-6">
          <div className="rounded-2xl border bg-white p-8 shadow-sm">
            <h1 className="text-3xl font-bold">Transactions</h1>
            <p className="mt-2 text-gray-600">
              Deposit, withdraw, transfer money, and view your transaction
              history.
            </p>
          </div>

          <div className="rounded-2xl border bg-white p-6 shadow-sm">
            <h2 className="text-xl font-semibold">Deposit / Withdraw</h2>
            <p className="mt-1 text-sm text-gray-600">
              Choose an account and enter an amount.
            </p>

            <form
              onSubmit={handleMoneyAction}
              className="mt-5 grid gap-4 md:grid-cols-4"
            >
              <div className="md:col-span-2">
                <label className="text-sm font-medium">Account</label>
                <select
                  value={selectedAccountId}
                  onChange={(event) => setSelectedAccountId(event.target.value)}
                  className="mt-1 w-full rounded-lg border px-3 py-2 outline-none focus:ring-2 focus:ring-black"
                  required
                >
                  {accounts.length === 0 ? (
                    <option value="">No accounts found</option>
                  ) : (
                    accounts.map((account) => (
                      <option key={account.id} value={account.id}>
                        {account.name} — {account.account_number} — NPR{" "}
                        {Number(account.balance || 0).toLocaleString()}
                      </option>
                    ))
                  )}
                </select>
              </div>

              <div>
                <label className="text-sm font-medium">Amount</label>
                <input
                  type="number"
                  min="1"
                  step="0.01"
                  value={amount}
                  onChange={(event) => setAmount(event.target.value)}
                  placeholder="1000"
                  className="mt-1 w-full rounded-lg border px-3 py-2 outline-none focus:ring-2 focus:ring-black"
                  required
                />
              </div>

              <div className="flex gap-2 md:items-end">
                <button
                  type="submit"
                  disabled={submittingAction !== null || accounts.length === 0}
                  
                  onClick={() => setActionType("deposit")}
                  className="w-full rounded-lg bg-black px-4 py-2 font-medium text-white disabled:opacity-60"
                >
                  {submittingAction === "deposit" ? "Depositing..." : "Deposit"}
                </button>

                <button
                  type="submit"
                  disabled={submittingAction !== null || accounts.length === 0}
                  onClick={() => setActionType("withdraw")}
                  className="w-full rounded-lg border px-4 py-2 font-medium disabled:opacity-60"
                >
                  {submittingAction === "withdraw"
                    ? "Withdrawing..."
                    : "Withdraw"}
                </button>
              </div>
            </form>
          </div>

          <div className="rounded-2xl border bg-white p-6 shadow-sm">
            <h2 className="text-xl font-semibold">Transfer Money</h2>
            <p className="mt-1 text-sm text-gray-600">
              Move money between your own accounts.
            </p>

            <form
              onSubmit={handleTransfer}
              className="mt-5 grid gap-4 md:grid-cols-4"
            >
              <div>
                <label className="text-sm font-medium">From Account</label>
                <select
                  value={fromAccountId}
                  onChange={(event) => setFromAccountId(event.target.value)}
                  className="mt-1 w-full rounded-lg border px-3 py-2 outline-none focus:ring-2 focus:ring-black"
                  required
                >
                  {accounts.length === 0 ? (
                    <option value="">No accounts found</option>
                  ) : (
                    accounts.map((account) => (
                      <option key={account.id} value={account.id}>
                        {account.name} — NPR{" "}
                        {Number(account.balance || 0).toLocaleString()}
                      </option>
                    ))
                  )}
                </select>
              </div>

              <div>
                <label className="text-sm font-medium">To Account</label>
                <select
                  value={toAccountId}
                  onChange={(event) => setToAccountId(event.target.value)}
                  className="mt-1 w-full rounded-lg border px-3 py-2 outline-none focus:ring-2 focus:ring-black"
                  required
                >
                  {accounts.length === 0 ? (
                    <option value="">No accounts found</option>
                  ) : (
                    accounts.map((account) => (
                      <option key={account.id} value={account.id}>
                        {account.name} — {account.account_number}
                      </option>
                    ))
                  )}
                </select>
              </div>

              <div>
                <label className="text-sm font-medium">Amount</label>
                <input
                  type="number"
                  min="1"
                  step="0.01"
                  value={transferAmount}
                  onChange={(event) => setTransferAmount(event.target.value)}
                  placeholder="500"
                  className="mt-1 w-full rounded-lg border px-3 py-2 outline-none focus:ring-2 focus:ring-black"
                  required
                />
              </div>

              <div className="flex md:items-end">
                <button
                  type="submit"
                  disabled={submittingAction !== null || accounts.length < 2}
                  className="w-full rounded-lg bg-black px-4 py-2 font-medium text-white disabled:opacity-60"
                >
                  {submittingAction === "transfer"
                    ? "Transferring..."
                    : "Transfer"}
                </button>
              </div>
            </form>
          </div>

          {success && (
            <p className="rounded-lg bg-green-50 px-3 py-2 text-sm text-green-700">
              {success}
            </p>
          )}

          {loading && (
            <div className="rounded-2xl border bg-white p-6 shadow-sm">
              <p className="text-gray-600">Loading transactions...</p>
            </div>
          )}

          {error && (
            <div className="rounded-2xl border border-red-200 bg-red-50 p-6">
              <p className="text-sm text-red-600">{error}</p>
            </div>
          )}

          {!loading && !error && transactions.length === 0 && (
            <div className="rounded-2xl border bg-white p-6 shadow-sm">
              <p className="text-gray-600">
                No transactions found. Try making a deposit first.
              </p>
            </div>
          )}

          {!loading && transactions.length > 0 && (
            <div className="rounded-2xl border bg-white shadow-sm">
              <div className="border-b p-6">
                <h2 className="text-xl font-semibold">Transaction History</h2>
                <p className="mt-1 text-sm text-gray-600">
                  Latest banking activity.
                </p>
              </div>

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
                        NPR {Number(transaction.amount || 0).toLocaleString()}
                      </p>

                      <span className="mt-1 inline-block rounded-full bg-gray-100 px-3 py-1 text-xs capitalize text-gray-700">
                        {transaction.status}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </Container>
    </ProtectedRoute>
  );
}
