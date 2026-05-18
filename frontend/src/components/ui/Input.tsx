import type { InputHTMLAttributes } from "react";

type InputProps = InputHTMLAttributes<HTMLInputElement> & {
  label?: string;
};

export default function Input({ label, className = "", ...props }: InputProps) {
  return (
    <div>
      {label && <label className="text-sm font-medium">{label}</label>}
      <input
        className={`mt-1 w-full rounded-lg border px-3 py-2 outline-none transition focus:ring-2 focus:ring-black ${className}`}
        {...props}
      />
    </div>
  );
}
