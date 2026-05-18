type PageHeaderProps = {
  label?: string;
  title: string;
  description?: string;
};

export default function PageHeader({
  label,
  title,
  description,
}: PageHeaderProps) {
  return (
    <div className="mb-8">
      {label && (
        <p className="text-sm font-semibold uppercase tracking-wide text-blue-400">
          {label}
        </p>
      )}

      <h1 className="mt-2 text-3xl font-bold text-white">{title}</h1>

      {description && (
        <p className="mt-2 max-w-2xl text-sm text-slate-400">{description}</p>
      )}
    </div>
  );
}
