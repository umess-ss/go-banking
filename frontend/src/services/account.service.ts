import { apiFetch } from "./api";
import type { Account, CreateAccountPayload } from "@/types/account";

type AccountsResponse = {
  accounts?: Account[];
  data?: Account[];
  message?: string;
};

type AccountResponse = {
  account?: Account;
  data?: Account;
  message?: string;
};

export async function getAccounts() {
  const response = await apiFetch<AccountsResponse | Account[]>("/accounts", {
    method: "GET",
    auth: true,
  });

  if (Array.isArray(response)) {
    return response;
  }

  return response.accounts || response.data || [];
}

export async function createAccount(payload: CreateAccountPayload) {
  const response = await apiFetch<AccountResponse | Account>("/accounts", {
    method: "POST",
    body: JSON.stringify(payload),
    auth: true,
  });

  if ("account" in response && response.account) {
    return response.account;
  }

  if ("data" in response && response.data) {
    return response.data;
  }

  return response as Account;
}