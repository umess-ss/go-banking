import type { SelectHTMLAttributes } from "react";

type SelectProps = SelectHTMLAttributes<HTMLSelectElement> & {
  label?: string;
  children: React.ReactNode;
};

export default function Select({
  label,
  children,
  className = "",
  ...props
}: SelectProps) {
  return (
    <div>
      {label && <label className="text-sm font-medium">{label}</label>}
      <select
        className={`mt-1 w-full rounded-lg border px-3 py-2 outline-none transition focus:ring-2 focus:ring-black ${className}`}
        {...props}
      >
        {children}
      </select>
    </div>
  );
}
