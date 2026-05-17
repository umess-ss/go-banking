type ContainerProps = {
  children: React.ReactNode;
};

export default function Container({ children }: ContainerProps) {
  return (
    <main className="mx-auto w-full max-w-6xl px-4 py-6">
      {children}
    </main>
  );
}
