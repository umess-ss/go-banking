type AlertVariant = "success" | "error" | "info";

type AlertProps = {
  children: React.ReactNode;
  variant?: AlertVariant;
};

export default function Alert({ children, variant = "info" }: AlertProps) {
  const variants: Record<AlertVariant, string> = {
    success: "bg-green-50 text-green-700 border-green-200",
    error: "bg-red-50 text-red-600 border-red-200",
    info: "bg-gray-50 text-gray-600 border-gray-200",
  };

  return (
    <div className={`rounded-lg border px-4 py-3 text-sm ${variants[variant]}`}>
      {children}
    </div>
  );
}
