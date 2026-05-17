import Link from "next/link";

export default function Navbar() {
  return (
    <header className="border-b bg-white">
      <nav className="mx-auto flex max-w-6xl items-center justify-between px-4 py-4">
        <Link href="/" className="text-xl font-bold">
          Go Banking
        </Link>

        <div className="flex items-center gap-4 text-sm">
          <Link href="/login">Login</Link>
          <Link href="/register">Register</Link>
          <Link href="/dashboard">Dashboard</Link>
          <Link href="/accounts">Accounts</Link>
          <Link href="/transactions">Transactions</Link>
        </div>
      </nav>
    </header>
  );
}
