export default function LoadingState({ text = "Loading..." }: { text?: string }) {
  return (
    <div className="rounded-2xl border bg-white p-6 shadow-sm">
      <p className="text-gray-600">{text}</p>
    </div>
  );
}
