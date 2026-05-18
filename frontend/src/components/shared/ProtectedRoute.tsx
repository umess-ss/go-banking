"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { getToken } from "@/lib/auth";

type ProtectedRouteProps = {
  children: React.ReactNode;
};

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const router = useRouter();
  const [checking, setChecking] = useState(true);

  useEffect(() => {
    queueMicrotask(() => {
      const token = getToken();

      if (!token) {
        router.replace("/login");
        return;
      }

      setChecking(false);
    });
  }, [router]);

  if (checking) {
    return (
      <main className="mx-auto w-full max-w-6xl px-4 py-6">
        <div className="rounded-2xl border bg-white p-8 shadow-sm">
          <p className="text-gray-600">Checking authentication...</p>
        </div>
      </main>
    );
  }

  return <>{children}</>;
}
