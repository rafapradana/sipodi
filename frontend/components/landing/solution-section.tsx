import { Target, CheckCircle2, TrendingUp } from "lucide-react";

const benefits = [
  {
    icon: Target,
    title: "Satu Database Terpusat",
    description: "Semua data GTK tersimpan aman dan dapat diakses kapan saja",
  },
  {
    icon: CheckCircle2,
    title: "Verifikasi Terstruktur",
    description: "Sistem approval memastikan validitas setiap data talenta",
  },
  {
    icon: TrendingUp,
    title: "Dashboard Analitik",
    description: "Visualisasi data untuk pengambilan keputusan yang lebih baik",
  },
];

export function SolutionSection() {
  return (
    <section className="py-20">
      <div className="container mx-auto px-4">
        <div className="max-w-4xl mx-auto">
          <div className="text-center mb-12">
            <h2 className="text-3xl md:text-4xl font-bold mb-4">
              SIPODI: Solusi Terintegrasi untuk Pengelolaan GTK
            </h2>
            <p className="text-lg text-muted-foreground max-w-2xl mx-auto">
              SIPODI hadir sebagai sistem informasi terpadu yang menghubungkan
              data GTK, sekolah, dan talenta dalam satu platform. Dengan SIPODI,
              Anda dapat memetakan potensi, memverifikasi prestasi, dan mengambil
              keputusan berbasis data — semua dalam satu tempat.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {benefits.map((benefit) => (
              <div
                key={benefit.title}
                className="flex flex-col items-center text-center"
              >
                <div className="mb-4 rounded-full bg-primary/10 p-4">
                  <benefit.icon className="h-8 w-8 text-primary" />
                </div>
                <h3 className="font-semibold mb-2">{benefit.title}</h3>
                <p className="text-sm text-muted-foreground">
                  {benefit.description}
                </p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
