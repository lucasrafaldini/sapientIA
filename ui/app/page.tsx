export default function Home() {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-24">
      <div className="text-center">
        <h1 className="text-4xl font-bold mb-4">🌿 SapientIA</h1>
        <p className="text-xl text-gray-600 mb-8">
          Sabedoria é a inteligência que escuta.
        </p>
        <div className="bg-blue-50 border border-blue-200 rounded-lg p-6 max-w-2xl">
          <p className="text-sm text-blue-900">
            ⚠️ UI em desenvolvimento (v0.2)
          </p>
          <p className="text-sm text-gray-700 mt-2">
            A interface web será disponibilizada na versão 0.2 do roadmap.
            <br />
            Por enquanto, utilize a CLI:{" "}
            <code className="bg-gray-100 px-2 py-1 rounded">
              sapientia run pipeline.yaml
            </code>
          </p>
        </div>
      </div>
    </main>
  );
}
