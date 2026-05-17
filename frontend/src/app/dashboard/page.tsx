import Container from "@/components/shared/Container";
import ProtectedRoute from "@/components/shared/ProtectedRoute";

export default function DashboardPage() {
  return (
    <ProtectedRoute>
      <Container>
        <div className="rounded-2xl border bg-white p-8 shadow-sm">
          <h1 className="text-3xl font-bold">Dashboard</h1>
          <p className="mt-2 text-gray-600">
            You are logged in. Banking overview will be shown here.
          </p>
        </div>
      </Container>
    </ProtectedRoute>
  );
}