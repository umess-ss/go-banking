import { getToken, removeToken } from "@/lib/auth";

const API_BASE_URL =
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

type ApiFetchOptions = RequestInit & {
  auth?: boolean;
};

export async function apiFetch<T>(
  endpoint: string,
  options: ApiFetchOptions = {}
): Promise<T> {
  const { auth = true, headers, ...rest } = options;

  const token = getToken();

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...rest,
    headers: {
      "Content-Type": "application/json",
      ...(auth && token ? { Authorization: `Bearer ${token}` } : {}),
      ...headers,
    },
  });

  const data = await response.json().catch(() => null);

  // Only protected routes should trigger session-expired logout
  if (response.status === 401 && auth) {
    if (typeof window !== "undefined") {
      removeToken();
      window.location.href = "/login";
    }

    throw new Error("Session expired. Please login again.");
  }

  // Login/register 401 should show normal backend error
  if (!response.ok) {
    throw new Error(
      data?.error ||
        data?.message ||
        data?.data?.error ||
        `Request failed with status ${response.status}`
    );
  }

  return data as T;
}
