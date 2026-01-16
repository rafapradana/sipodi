import Link from "next/link";
import { ArrowRight, BarChart3, Award, GraduationCap } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";

export function HeroSection() {
  return (
    <section id="beranda" className="relative overflow-hidden py-20 md:py-32">
      {/* Background gradient */}
      <div className="absolute inset-0 -z-10 bg-gradient-to-b from-primary/5 via-background to-background" />
      
      {/* Floating elements */}
      <div className="absolute top-20 left-10 -z-10 opacity-20 dark:opacity-10">
        <GraduationCap className="h-24 w-24 text-primary" />
      </div>
      <div className="absolute bottom-20 right-10 -z-10 opacity-20 dark:opacity-10">
        <Award className="h-20 w-20 text-primary" />
      </div>
      <div className="absolute top-40 right-20 -z-10 opacity-20 dark:opacity-10">
        <BarChart3 className="h-16 w-16 text-primary" />
      </div>

      <div className="container mx-auto px-4">
        <div className="flex flex-col items-center text-center max-w-4xl mx-auto">
          {/* Badge */}
          <Badge variant="secondary" className="mb-6">
            🏛️ Cabang Dinas Pendidikan Wilayah Malang
          </Badge>

          {/* Headline */}
          <h1 className="text-4xl md:text-5xl lg:text-6xl font-bold tracking-tight mb-6">
            Kelola Potensi GTK dengan{" "}
            <span className="text-primary">Lebih Cerdas</span>
          </h1>

          {/* Subheadline */}
          <p className="text-lg md:text-xl text-muted-foreground mb-8 max-w-2xl">
            Sistem informasi terintegrasi untuk memetakan kompetensi, prestasi,
            dan talenta Guru & Tenaga Kependidikan di wilayah Malang.
          </p>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row gap-4">
            <Button size="lg" asChild>
              <Link href="/login">
                Mulai Sekarang
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
            <Button size="lg" variant="outline" asChild>
              <Link href="#fitur">Pelajari Lebih Lanjut</Link>
            </Button>
          </div>

          {/* Dashboard Preview */}
          <div className="mt-16 w-full max-w-5xl">
            <div className="relative rounded-xl border bg-card shadow-2xl overflow-hidden">
              <div className="absolute inset-0 bg-gradient-to-t from-background/80 to-transparent z-10" />
              <div className="p-4 md:p-8">
                {/* Mock Dashboard */}
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
                  {[
                    { label: "Total GTK", value: "5,234" },
                    { label: "Sekolah", value: "156" },
                    { label: "Talenta", value: "12,847" },
                    { label: "Terverifikasi", value: "89%" },
                  ].map((stat) => (
                    <div
                      key={stat.label}
                      className="rounded-lg bg-muted/50 p-4 text-center"
                    >
                      <p className="text-2xl md:text-3xl font-bold">
                        {stat.value}
                      </p>
                      <p className="text-sm text-muted-foreground">
                        {stat.label}
                      </p>
                    </div>
                  ))}
                </div>
                <div className="h-32 md:h-48 rounded-lg bg-muted/30 flex items-center justify-center">
                  <BarChart3 className="h-16 w-16 text-muted-foreground/50" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
