type EmptyStateProps = {
  children: React.ReactNode;
  className?: string;
  variant?: "card" | "plain";
};

export default function EmptyState({
  children,
  className = "",
  variant = "card",
}: EmptyStateProps) {
  if (variant === "plain") {
    return <p className={`text-gray-600 ${className}`}>{children}</p>;
  }

  return (
    <div className={`rounded-2xl border bg-white p-6 shadow-sm ${className}`}>
      <p className="text-gray-600">{children}</p>
    </div>
  );
}
