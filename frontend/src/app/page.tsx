import Link from "next/link";

export default function HomePage() {
  return (
    <main className="min-h-screen bg-slate-950 text-white">
      <section className="mx-auto flex min-h-screen max-w-6xl flex-col items-center justify-center px-6 text-center">
        <p className="mb-4 rounded-full border border-blue-500/30 bg-blue-500/10 px-4 py-2 text-sm font-semibold text-blue-300">
          Go Banking Platform
        </p>

        <h1 className="max-w-4xl text-4xl font-bold tracking-tight md:text-6xl">
          Secure banking dashboard built with Go, PostgreSQL, and Next.js.
        </h1>

        <p className="mt-6 max-w-2xl text-lg leading-8 text-slate-300">
          Manage users, bank accounts, deposits, withdrawals, transfers, and
          transaction history through a clean modern interface.
        </p>

        <div className="mt-10 flex flex-col gap-4 sm:flex-row">
          <Link
            href="/login"
            className="rounded-xl bg-blue-600 px-6 py-3 font-semibold transition hover:bg-blue-700"
          >
            Login
          </Link>

          <Link
            href="/register"
            className="rounded-xl border border-slate-700 px-6 py-3 font-semibold transition hover:bg-slate-900"
          >
            Create Account
          </Link>
        </div>
      </section>
    </main>
  );
}