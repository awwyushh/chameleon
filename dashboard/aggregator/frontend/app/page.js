import Navbar from "../components/Navbar";

export default function Home() {
  return (
    <div>
      <Navbar />
      <div className="flex flex-col items-center justify-center h-[80vh] text-center p-5">
        <h1 className="text-5xl font-bold mb-4">
          Welcome to My App 🚀
        </h1>
        <p className="text-lg text-gray-600 max-w-2xl">
          A simple full-stack project with authentication using Next.js + MongoDB + Express.
        </p>
      </div>
    </div>
  );
}
