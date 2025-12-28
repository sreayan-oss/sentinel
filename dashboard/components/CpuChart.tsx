"use client"; // <--- Crucial: Tells Next.js this runs in the browser

import { useEffect, useState } from "react";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

// Data shape matches your Go Struct
type Metric = {
  id: number;
  agent_id: string;
  cpu: number;
  timestamp: number;
};

export default function CpuChart() {
  const [data, setData] = useState<Metric[]>([]);

  // The Fetch Loop
  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await fetch("http://localhost:8080/metrics");
        const json = await res.json();
        
        // Reverse array so graph moves Left -> Right (Oldest -> Newest)
        // Convert timestamp to readable time string
        const formatted = json.reverse().map((item: Metric) => ({
          ...item,
          time: new Date(item.timestamp * 1000).toLocaleTimeString([], { 
              hour12: false, 
              hour: '2-digit', 
              minute: '2-digit', 
              second: '2-digit' 
          }),
        }));
        
        setData(formatted);
      } catch (err) {
        console.error("Failed to fetch metrics:", err);
      }
    };

    // Fetch immediately, then every 2 seconds
    fetchData();
    const interval = setInterval(fetchData, 2000);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="w-full h-full min-h-[250px]">
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data} margin={{ top: 10, right: 30, left: 0, bottom: 20 }}>
          <CartesianGrid strokeDasharray="3 3" stroke="#333" />
          <XAxis 
            dataKey="time" 
            stroke="#888" 
            fontSize={12} 
            tick={{fill: '#6b7280'}}
          />
          <YAxis 
            stroke="#888" 
            fontSize={12} 
            domain={[0, 100]} // Fix Y-axis from 0% to 100%
            tick={{fill: '#6b7280'}}
          />
          <Tooltip 
            contentStyle={{ backgroundColor: '#111', border: '1px solid #333' }}
            itemStyle={{ color: '#fff' }}
          />
          <Line
            type="monotone" // Makes the line curvy
            dataKey="cpu"
            stroke="#3b82f6" // Tailwind Blue-500
            strokeWidth={2}
            dot={false} // Hide dots for a cleaner look
            activeDot={{ r: 8 }}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  );
}