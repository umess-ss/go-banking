import Container from "@/components/shared/Container";

export default function AccountsPage() {
  return (
    <Container>
      <div className="rounded-2xl border bg-white p-8 shadow-sm">
        <h1 className="text-3xl font-bold">Accounts</h1>
        <p className="mt-2 text-gray-600">
          User bank accounts will be listed here.
        </p>
      </div>
    </Container>
  );
}
