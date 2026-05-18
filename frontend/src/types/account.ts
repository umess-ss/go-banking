export type Account = {
  id: number;
  user_id?: number;
  name: string;
  account_number: string;
  account_type: string;
  balance: number;
  currency: string;
  created_at?: string;
  updated_at?: string;
};

export type CreateAccountPayload = {
  name: string;
  account_type?: string;
  currency: string;
};