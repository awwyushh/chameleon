"use client";

import * as React from "react";
import { useState, useRef, useEffect } from "react";
import { useRouter } from "next/navigation";
import { API_URL } from "../../lib/api";
import {
    Card,
    CardHeader,
    CardTitle,
    CardDescription,
    CardContent,
    CardFooter,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Mail, ArrowRight, ArrowLeft } from "lucide-react";

export default function ForgotPasswordPage() {
    const router = useRouter();
    const [email, setEmail] = useState("");
    const [error, setError] = useState("");
    const [success, setSuccess] = useState("");
    const [loading, setLoading] = useState(false);

    const canvasRef = useRef(null);
    useEffect(() => {
        const canvas = canvasRef.current;
        const ctx = canvas?.getContext("2d");
        if (!canvas || !ctx) return;

        const setSize = () => {
            canvas.width = window.innerWidth;
            canvas.height = window.innerHeight;
        };
        setSize();

        let ps = [];
        let raf = 0;

        const make = () => ({
            x: Math.random() * canvas.width,
            y: Math.random() * canvas.height,
            v: Math.random() * 0.25 + 0.05,
            o: Math.random() * 0.35 + 0.15,
        });

        const init = () => {
            ps = [];
            const count = Math.floor((canvas.width * canvas.height) / 9000);
            for (let i = 0; i < count; i++) ps.push(make());
        };

        const draw = () => {
            ctx.clearRect(0, 0, canvas.width, canvas.height);
            ps.forEach((p) => {
                p.y -= p.v;
                if (p.y < 0) {
                    p.x = Math.random() * canvas.width;
                    p.y = canvas.height + Math.random() * 40;
                    p.v = Math.random() * 0.25 + 0.05;
                    p.o = Math.random() * 0.35 + 0.15;
                }
                ctx.fillStyle = `rgba(250,250,250,${p.o})`;
                ctx.fillRect(p.x, p.y, 0.7, 2.2);
            });
            raf = requestAnimationFrame(draw);
        };

        const onResize = () => {
            setSize();
            init();
        };

        window.addEventListener("resize", onResize);
        init();
        raf = requestAnimationFrame(draw);
        return () => {
            window.removeEventListener("resize", onResize);
            cancelAnimationFrame(raf);
        };
    }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();
        setError("");
        setSuccess("");
        setLoading(true);

        try {
            // TODO: Replace with actual forgot password API endpoint
            // const res = await fetch(`${API_URL}/auth/forgot-password`, {
            //   method: "POST",
            //   headers: { "Content-Type": "application/json" },
            //   body: JSON.stringify({ companyEmail: email }),
            // });

            // For now, simulate API call
            await new Promise((resolve) => setTimeout(resolve, 1000));

            // Simulate success response
            setSuccess("Password reset link has been sent to your email!");
            setLoading(false);

            // Optionally redirect after a delay
            // setTimeout(() => {
            //   router.push("/login");
            // }, 3000);
        } catch (err) {
            setError("Network error. Please try again.");
            setLoading(false);
        }
    };

    return (
        <section className="fixed inset-0 bg-zinc-950 text-zinc-50">
            <style>{`
        .accent-lines{position:absolute;inset:0;pointer-events:none;opacity:.7}
        .hline,.vline{position:absolute;background:#27272a;will-change:transform,opacity}
        .hline{left:0;right:0;height:1px;transform:scaleX(0);transform-origin:50% 50%;animation:drawX .8s cubic-bezier(.22,.61,.36,1) forwards}
        .vline{top:0;bottom:0;width:1px;transform:scaleY(0);transform-origin:50% 0%;animation:drawY .9s cubic-bezier(.22,.61,.36,1) forwards}
        .hline:nth-child(1){top:18%;animation-delay:.12s}
        .hline:nth-child(2){top:50%;animation-delay:.22s}
        .hline:nth-child(3){top:82%;animation-delay:.32s}
        .vline:nth-child(4){left:22%;animation-delay:.42s}
        .vline:nth-child(5){left:50%;animation-delay:.54s}
        .vline:nth-child(6){left:78%;animation-delay:.66s}
        .hline::after,.vline::after{content:"";position:absolute;inset:0;background:linear-gradient(90deg,transparent,rgba(250,250,250,.24),transparent);opacity:0;animation:shimmer .9s ease-out forwards}
        .hline:nth-child(1)::after{animation-delay:.12s}
        .hline:nth-child(2)::after{animation-delay:.22s}
        .hline:nth-child(3)::after{animation-delay:.32s}
        .vline:nth-child(4)::after{animation-delay:.42s}
        .vline:nth-child(5)::after{animation-delay:.54s}
        .vline:nth-child(6)::after{animation-delay:.66s}
        @keyframes drawX{0%{transform:scaleX(0);opacity:0}60%{opacity:.95}100%{transform:scaleX(1);opacity:.7}}
        @keyframes drawY{0%{transform:scaleY(0);opacity:0}60%{opacity:.95}100%{transform:scaleY(1);opacity:.7}}
        @keyframes shimmer{0%{opacity:0}35%{opacity:.25}100%{opacity:0}}

        /* === Card minimal fade-up animation === */
        .card-animate {
          opacity: 0;
          transform: translateY(20px);
          animation: fadeUp 0.8s cubic-bezier(.22,.61,.36,1) 0.4s forwards;
        }
        @keyframes fadeUp {
          to {
            opacity: 1;
            transform: translateY(0);
          }
        }
      `}</style>
            {/* Subtle vignette */}
            <div
                className="absolute inset-0 pointer-events-none [background:radial-gradient(80%_60%_at_50%_30%,rgba(255,255,255,0.06),transparent_60%)]" />
            {/* Animated accent lines */}
            <div className="accent-lines">
                <div className="hline" />
                <div className="hline" />
                <div className="hline" />
                <div className="vline" />
                <div className="vline" />
                <div className="vline" />
            </div>
            {/* Particles */}
            <canvas
                ref={canvasRef}
                className="absolute inset-0 w-full h-full opacity-50 mix-blend-screen pointer-events-none" />
            {/* Header */}
            <header
                className="absolute left-0 right-0 top-0 flex items-center justify-between px-6 py-4 border-b border-zinc-800/80">
                <button
                    onClick={() => router.push("/login")}
                    className="text-xs tracking-[0.14em] uppercase text-zinc-400 hover:text-zinc-200 transition-colors">
                    ← Back to Login
                </button>
                <Button
                    variant="outline"
                    onClick={() => router.push("/signup")}
                    className="h-9 rounded-lg border-zinc-800 bg-zinc-900 text-zinc-50 hover:bg-zinc-900/80">
                    <span className="mr-2">Sign Up</span>
                    <ArrowRight className="h-4 w-4" />
                </Button>
            </header>
            {/* Centered Forgot Password Card */}
            <div className="h-full w-full grid place-items-center px-4">
                <Card
                    className="card-animate w-full max-w-sm border-zinc-800 bg-zinc-900/70 backdrop-blur supports-[backdrop-filter]:bg-zinc-900/60">
                    <CardHeader className="space-y-1">
                        <CardTitle className="text-2xl">Forgot Password</CardTitle>
                        <CardDescription className="text-zinc-400">
                            Enter your email address and we'll send you a link to reset your password
                        </CardDescription>
                    </CardHeader>

                    <CardContent className="grid gap-5">
                        {!success ? (
                            <form onSubmit={handleSubmit} className="grid gap-5">
                                <div className="grid gap-2">
                                    <Label htmlFor="email" className="text-zinc-300">
                                        Company Email
                                    </Label>
                                    <div className="relative">
                                        <Mail
                                            className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-zinc-500" />
                                        <Input
                                            id="email"
                                            type="email"
                                            placeholder="you@company.com"
                                            value={email}
                                            onChange={(e) => setEmail(e.target.value)}
                                            required
                                            className="pl-10 bg-zinc-950 border-zinc-800 text-zinc-50 placeholder:text-zinc-600" />
                                    </div>
                                </div>

                                {error && (
                                    <div className="text-sm text-red-400 bg-red-950/20 border border-red-800/50 rounded-lg p-3">
                                        {error}
                                    </div>
                                )}

                                <Button
                                    type="submit"
                                    disabled={loading}
                                    className="w-full h-10 rounded-lg bg-zinc-50 text-zinc-900 hover:bg-zinc-200 disabled:opacity-50 disabled:cursor-not-allowed">
                                    {loading ? "Sending..." : "Send Reset Link"}
                                </Button>
                            </form>
                        ) : (
                            <div className="space-y-5">
                                <div className="text-sm text-green-400 bg-green-950/20 border border-green-800/50 rounded-lg p-4 text-center">
                                    {success}
                                </div>
                                <p className="text-sm text-zinc-400 text-center">
                                    Please check your email for the password reset link. If you don't see it, check your spam folder.
                                </p>
                                <Button
                                    onClick={() => router.push("/login")}
                                    className="w-full h-10 rounded-lg bg-zinc-50 text-zinc-900 hover:bg-zinc-200">
                                    <ArrowLeft className="h-4 w-4 mr-2" />
                                    Back to Login
                                </Button>
                            </div>
                        )}
                    </CardContent>

                    <CardFooter className="flex items-center justify-center text-sm text-zinc-400">
                        Remember your password?
                        <a
                            className="ml-1 text-zinc-200 hover:underline cursor-pointer"
                            onClick={() => router.push("/login")}>
                            Login
                        </a>
                    </CardFooter>
                </Card>
            </div>
        </section>
    );
}

