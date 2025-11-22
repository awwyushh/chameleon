"use client";
import Link from "next/link";

export default function Navbar() {
  return (
    <nav className="w-full flex justify-between items-center p-5 bg-white shadow-md">
      <h1 className="text-2xl font-bold">My App</h1>

      <div className="space-x-4">
        <Link href="/login" className="px-4 py-2 bg-blue-500 text-white rounded">
          Login
        </Link>
        <Link href="/signup" className="px-4 py-2 bg-green-500 text-white rounded">
          Signup
        </Link>
      </div>
    </nav>
  );
}
