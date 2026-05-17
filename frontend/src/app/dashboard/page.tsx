import Container from "@/components/shared/Container";

export default function DashboardPage() {
  return (
    <Container>
      <div className="rounded-2xl border bg-white p-8 shadow-sm">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <p className="mt-2 text-gray-600">
          Account overview and banking summary will be shown here.
        </p>
      </div>
    </Container>
  );
}
