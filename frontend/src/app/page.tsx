import Link from "next/link";
import Container from "@/components/shared/Container";

export default function HomePage() {
  return (
    <Container>
      <section className="rounded-2xl border bg-white p-8 shadow-sm">
        <p className="text-sm font-medium text-gray-500">Go Banking</p>

        <h1 className="mt-3 text-4xl font-bold tracking-tight">
          Banking API Frontend
        </h1>

        <p className="mt-4 max-w-2xl text-gray-600">
          A Next.js frontend connected to your Go banking backend. This app will
          include authentication, accounts, transactions, and dashboard features.
        </p>

        <div className="mt-6 flex gap-3">
          <Link
            href="/login"
            className="rounded-lg bg-black px-4 py-2 text-white"
          >
            Login
          </Link>

          <Link
            href="/register"
            className="rounded-lg border px-4 py-2"
          >
            Create Account
          </Link>
        </div>
      </section>
    </Container>
  );
}