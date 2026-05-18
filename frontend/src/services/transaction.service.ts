import { apiFetch } from "./api";
import type { Transaction } from "@/types/transaction";

type TransactionsResponse = {
  success?: boolean;
  message?: string;
  data?: Transaction[];
  transactions?: Transaction[];
};

type TransactionResponse = {
  success?: boolean;
  message?: string;
  data?: Transaction;
  transaction?: Transaction;
};

export async function getTransactions() {
  const response = await apiFetch<TransactionsResponse | Transaction[]>(
    "/transactions",
    {
      method: "GET",
      auth: true,
    }
  );

  if (Array.isArray(response)) {
    return response;
  }

  return response.data || response.transactions || [];
}

export async function deposit(accountId: number, amount: number) {
  const response = await apiFetch<TransactionResponse>(
    `/accounts/${accountId}/deposit`,
    {
      method: "POST",
      body: JSON.stringify({ amount }),
      auth: true,
    }
  );

  return response.data || response.transaction;
}

export async function withdraw(accountId: number, amount: number) {
  const response = await apiFetch<TransactionResponse>(
    `/accounts/${accountId}/withdraw`,
    {
      method: "POST",
      body: JSON.stringify({ amount }),
      auth: true,
    }
  );

  return response.data || response.transaction;
}


export async function transfer(
  fromAccountId: number,
  toAccountId: number,
  amount: number
) {
  const response = await apiFetch<TransactionResponse>("/transfer", {
    method: "POST",
    body: JSON.stringify({
      from_account_id: fromAccountId,
      to_account_id: toAccountId,
      amount,
    }),
    auth: true,
  });

  return response.data || response.transaction;
}