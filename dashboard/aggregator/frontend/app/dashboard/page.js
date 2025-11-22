"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { API_URL } from "../../lib/api";

export default function Dashboard() {
  const router = useRouter();

  const [user, setUser] = useState(null);
  const [kpis, setKpis] = useState(null);
  const [recentAttacks, setRecentAttacks] = useState([]);

  // AUTH + FETCH USER
  useEffect(() => {
    const token = localStorage.getItem("token");
    if (!token) return router.push("/login");

    const fetchUser = async () => {
      const res = await fetch(`${API_URL}/auth/me`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      const data = await res.json();
      if (!data.success) return router.push("/login");
      setUser(data.user);
    };

    fetchUser();
  }, []);

  // FETCH DASHBOARD DATA
  useEffect(() => {
    const fetchData = async () => {
      const res = await fetch(`${API_URL}/attacks/dashboard`);
      const data = await res.json();

      if (data.success) {
        setKpis(data.kpis);
        setRecentAttacks(data.recent);
      }
    };

    fetchData();
  }, []);

  if (!user || !kpis) {
    return (
      <div className="flex items-center justify-center h-screen text-xl">
        Loading...
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-[#0b0b0b] text-white px-6 py-8">

      <h1 className="text-4xl font-bold mb-6">Dashboard</h1>
      <p className="text-lg mb-10">Welcome, <b>{user.name}</b> 👋</p>

      {/* KPI CARDS */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-10">
        <KpiCard title="Total Attacks" value={kpis.totalAttacks} />
        <KpiCard title="SQL Injection Attempts" value={kpis.sqlAttacks} />
        <KpiCard title="XSS Attempts" value={kpis.xssAttacks} />
        <KpiCard title="Bruteforce Attempts" value={kpis.bruteAttacks} />
      </div>

      {/* ATTACK LIST */}
      <div className="bg-[#111] border border-black/30 p-6 rounded-xl">
        <h2 className="text-2xl font-semibold mb-4">Recent Attacks</h2>

        <div className="space-y-4">
          {recentAttacks.map((a) => (
            <div key={a._id} className="border-b border-neutral-800 pb-4">
              <p className="text-sm text-neutral-400">
                {new Date(a.time).toLocaleString()}
              </p>

              <p className="text-white mt-1">
                <b>{a.classification.label}</b> ({Math.round(a.classification.confidence)}%)
              </p>

              <p className="text-neutral-400 text-sm">
                IP: {a.srcIp} — {a.geo.country}, {a.geo.city}
              </p>

              <p className="mt-1 text-neutral-300">
                Input: <span className="text-red-400">{a.rawInput}</span>
              </p>

              <p className="mt-1 text-green-400">
                Deception: {a.deception}
              </p>

              <p className="text-yellow-300 text-sm">
                Tarpit Delay: {a.tarpitDelayMs}ms
              </p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function KpiCard({ title, value }) {
  return (
    <div className="bg-[#111] p-6 rounded-xl border border-black/30">
      <p className="text-neutral-400">{title}</p>
      <h2 className="text-3xl font-bold mt-2">{value}</h2>
    </div>
  );
}
