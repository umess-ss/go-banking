import type { ButtonHTMLAttributes } from "react";

type ButtonVariant = "primary" | "secondary" | "danger";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: ButtonVariant;
};

export default function Button({
  children,
  className = "",
  variant = "primary",
  ...props
}: ButtonProps) {
  const variants: Record<ButtonVariant, string> = {
    primary: "bg-black text-white hover:bg-gray-800",
    secondary: "border bg-white text-black hover:bg-gray-50",
    danger: "bg-red-600 text-white hover:bg-red-700",
  };

  return (
    <button
      className={`rounded-lg px-4 py-2 font-medium transition disabled:cursor-not-allowed disabled:opacity-60 ${variants[variant]} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
}