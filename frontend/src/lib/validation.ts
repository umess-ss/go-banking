export function parsePositiveAmount(value: string) {
  const amount = Number(value);

  if (!value.trim()) {
    return {
      valid: false,
      message: "Amount is required.",
      amount: 0,
    };
  }

  if (Number.isNaN(amount)) {
    return {
      valid: false,
      message: "Amount must be a valid number.",
      amount: 0,
    };
  }

  if (amount <= 0) {
    return {
      valid: false,
      message: "Amount must be greater than 0.",
      amount: 0,
    };
  }

  if (amount > 1_000_000) {
    return {
      valid: false,
      message: "Amount cannot exceed NPR 1,000,000.",
      amount: 0,
    };
  }

  return {
    valid: true,
    message: "",
    amount,
  };
}
