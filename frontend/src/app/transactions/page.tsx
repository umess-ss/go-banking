import Container from "@/components/shared/Container";
import ProtectedRoute from "@/components/shared/ProtectedRoute";

export default function TransactionsPage() {
  return (
    <ProtectedRoute>
      <Container>
        <div className="rounded-2xl border bg-white p-8 shadow-sm">
          <h1 className="text-3xl font-bold">Transactions</h1>
          <p className="mt-2 text-gray-600">
            Transaction history and transfer actions will be added here.
          </p>
        </div>
      </Container>
    </ProtectedRoute>
  );
}