export type TransactionType = "deposit" | "withdraw" | "transfer";

export type TransactionStatus = "pending" | "success" | "failed";

export type Transaction = {
  id: number;
  type: TransactionType;
  from_account_id?: number | null;
  to_account_id?: number | null;
  amount: number;
  status: TransactionStatus;
  reference_number: string;
  created_at: string;
};
