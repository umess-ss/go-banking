"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { getToken, removeToken } from "@/lib/auth";

const navLinks = [
  { href: "/dashboard", label: "Dashboard" },
  { href: "/accounts", label: "Accounts" },
  { href: "/transactions", label: "Transactions" },
];

export default function Navbar() {
  const pathname = usePathname();
  const router = useRouter();
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    setLoggedIn(Boolean(getToken()));
  }, [pathname]);

  function handleLogout() {
    removeToken();
    setLoggedIn(false);
    router.push("/login");
  }

  return (
    <header className="sticky top-0 z-50 border-b bg-white/95 backdrop-blur">
      <nav className="mx-auto flex max-w-6xl items-center justify-between px-4 py-5">
        <Link href="/dashboard" className="text-lg font-bold tracking-tight">
          Go Banking
        </Link>

        <div className="flex items-center gap-8 text-sm">
          {loggedIn ? (
            <>
              {navLinks.map((link) => {
                const active =
                  pathname === link.href ||
                  pathname.startsWith(`${link.href}/`);

                return (
                  <Link
                    key={link.href}
                    href={link.href}
                    className={`relative font-medium transition ${
                      active
                        ? "text-black"
                        : "text-gray-500 hover:text-black"
                    }`}
                  >
                    {link.label}

                    {active && (
                      <span className="absolute -bottom-2 left-0 h-0.5 w-full rounded-full bg-black" />
                    )}
                  </Link>
                );
              })}

              <button
                onClick={handleLogout}
                className="font-medium text-gray-500 transition hover:text-red-600"
              >
                Logout
              </button>
            </>
          ) : (
            <>
              <Link
                href="/login"
                className={`relative font-medium transition ${
                  pathname === "/login"
                    ? "text-black"
                    : "text-gray-500 hover:text-black"
                }`}
              >
                Login
                {pathname === "/login" && (
                  <span className="absolute -bottom-2 left-0 h-0.5 w-full rounded-full bg-black" />
                )}
              </Link>

              <Link
                href="/register"
                className={`relative font-medium transition ${
                  pathname === "/register"
                    ? "text-black"
                    : "text-gray-500 hover:text-black"
                }`}
              >
                Register
                {pathname === "/register" && (
                  <span className="absolute -bottom-2 left-0 h-0.5 w-full rounded-full bg-black" />
                )}
              </Link>
            </>
          )}
        </div>
      </nav>
    </header>
  );
}