export type Account = {
  id: number;
  user_id?: number;
  name: string;
  account_number: string;
  balance: number;
  currency: string;
  created_at?: string;
  updated_at?: string;
};

export type CreateAccountPayload = {
  name: string;
  currency: string;
};