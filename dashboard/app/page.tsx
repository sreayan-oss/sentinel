import CpuChart from "../components/CpuChart"; // <--- Import logic

export default function Home() {
  return (
    <main className="min-h-screen p-8 bg-neutral-950 text-white">
      <div className="max-w-5xl mx-auto">
        
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold tracking-tight text-blue-500">
            Sentinel
          </h1>
          <p className="text-neutral-400">
            Real-time Distributed System Monitor
          </p>
        </div>

        {/* Dashboard Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          
          {/* Card 1: CPU Usage */}
          <div className="p-6 rounded-xl bg-neutral-900 border border-neutral-800 h-96">
            <h2 className="text-lg font-semibold mb-4 text-neutral-200">
              CPU Usage (%)
            </h2>
            {/* THE CHART IS HERE NOW */}
            <CpuChart />
          </div>

          {/* Card 2: Status */}
          <div className="p-6 rounded-xl bg-neutral-900 border border-neutral-800 h-96">
             <h2 className="text-lg font-semibold mb-4 text-neutral-200">
              Active Agents
            </h2>
            <div className="text-4xl font-bold text-green-500">
              1
            </div>
            <p className="text-sm text-neutral-500 mt-2">Agent-007 is online</p>
          </div>

        </div>
      </div>
    </main>
  );
}