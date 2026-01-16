"use client";

import { useEffect, useState, useRef } from "react";
import { School, Users, Trophy, Zap } from "lucide-react";

const stats = [
  { icon: School, value: 150, suffix: "+", label: "Sekolah Terdaftar" },
  { icon: Users, value: 5000, suffix: "+", label: "GTK Terdata" },
  { icon: Trophy, value: 10000, suffix: "+", label: "Talenta Tercatat" },
  { icon: Zap, value: 99, suffix: "%", label: "Uptime Sistem" },
];

function AnimatedCounter({
  value,
  suffix,
}: {
  value: number;
  suffix: string;
}) {
  const [count, setCount] = useState(0);
  const [isVisible, setIsVisible] = useState(false);
  const ref = useRef<HTMLSpanElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true);
        }
      },
      { threshold: 0.1 }
    );

    if (ref.current) {
      observer.observe(ref.current);
    }

    return () => observer.disconnect();
  }, []);

  useEffect(() => {
    if (!isVisible) return;

    const duration = 2000;
    const steps = 60;
    const increment = value / steps;
    let current = 0;

    const timer = setInterval(() => {
      current += increment;
      if (current >= value) {
        setCount(value);
        clearInterval(timer);
      } else {
        setCount(Math.floor(current));
      }
    }, duration / steps);

    return () => clearInterval(timer);
  }, [isVisible, value]);

  const formatNumber = (num: number) => {
    return num.toLocaleString("id-ID");
  };

  return (
    <span ref={ref}>
      {formatNumber(count)}
      {suffix}
    </span>
  );
}

export function StatisticsSection() {
  return (
    <section className="py-20">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Dipercaya untuk Mengelola Data GTK Wilayah Malang
          </h2>
        </div>

        <div className="max-w-4xl mx-auto">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-8">
            {stats.map((stat) => (
              <div key={stat.label} className="text-center">
                <div className="flex justify-center mb-3">
                  <div className="rounded-full bg-primary/10 p-3">
                    <stat.icon className="h-6 w-6 text-primary" />
                  </div>
                </div>
                <p className="text-3xl md:text-4xl font-bold mb-1">
                  <AnimatedCounter value={stat.value} suffix={stat.suffix} />
                </p>
                <p className="text-sm text-muted-foreground">{stat.label}</p>
              </div>
            ))}
          </div>

          {/* Trust badges */}
          <div className="mt-12 flex flex-wrap items-center justify-center gap-6">
            <div className="flex items-center gap-2 rounded-full bg-muted px-4 py-2">
              <Building2Icon className="h-5 w-5 text-muted-foreground" />
              <span className="text-sm font-medium">
                Sistem Resmi Cabdin Malang
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}

function Building2Icon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
    >
      <path d="M6 22V4a2 2 0 0 1 2-2h8a2 2 0 0 1 2 2v18Z" />
      <path d="M6 12H4a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2h2" />
      <path d="M18 9h2a2 2 0 0 1 2 2v9a2 2 0 0 1-2 2h-2" />
      <path d="M10 6h4" />
      <path d="M10 10h4" />
      <path d="M10 14h4" />
      <path d="M10 18h4" />
    </svg>
  );
}
