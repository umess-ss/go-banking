import Navbar from "./Navbar";

type AppShellProps = {
  children: React.ReactNode;
};

export default function AppShell({ children }: AppShellProps) {
  return (
    <div className="min-h-screen bg-gray-50 text-gray-950">
      <Navbar />
      {children}
    </div>
  );
}