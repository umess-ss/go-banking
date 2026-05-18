import { apiFetch } from "./api";
import type { Account } from "@/types/account";

type AccountsResponse = {
  accounts?: Account[];
  data?: Account[];
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
