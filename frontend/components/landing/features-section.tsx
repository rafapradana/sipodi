import {
  Users,
  Trophy,
  CheckSquare,
  School,
  BarChart3,
  FileDown,
} from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

const features = [
  {
    icon: Users,
    title: "Manajemen Data GTK",
    description:
      "Kelola profil lengkap GTK termasuk NUPTK, NIP, jabatan, dan informasi personal dalam satu sistem.",
  },
  {
    icon: Trophy,
    title: "Pencatatan Talenta",
    description:
      "Dokumentasikan pelatihan, prestasi lomba, dan minat/bakat GTK dengan form yang terstruktur.",
  },
  {
    icon: CheckSquare,
    title: "Verifikasi Berjenjang",
    description:
      "Admin sekolah memverifikasi data talenta GTK sebelum masuk ke database resmi.",
  },
  {
    icon: School,
    title: "Database Sekolah",
    description:
      "Data lengkap sekolah termasuk NPSN, status, alamat, dan daftar GTK yang terdaftar.",
  },
  {
    icon: BarChart3,
    title: "Dashboard Analitik",
    description:
      "Visualisasi statistik GTK, talenta, dan performa per sekolah maupun wilayah.",
  },
  {
    icon: FileDown,
    title: "Export Laporan",
    description:
      "Unduh laporan dalam format PDF atau Excel untuk kebutuhan dokumentasi dan pelaporan.",
  },
];

export function FeaturesSection() {
  return (
    <section id="fitur" className="py-20 bg-muted/30">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Fitur Lengkap untuk Setiap Kebutuhan
          </h2>
          <p className="text-lg text-muted-foreground">
            Dirancang khusus untuk ekosistem pendidikan di wilayah Malang
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 max-w-6xl mx-auto">
          {features.map((feature) => (
            <Card
              key={feature.title}
              className="group hover:shadow-lg transition-all duration-300 hover:-translate-y-1"
            >
              <CardContent className="pt-6">
                <div className="flex flex-col items-center text-center">
                  <div className="mb-4 rounded-full bg-primary/10 p-3 group-hover:bg-primary/20 transition-colors">
                    <feature.icon className="h-6 w-6 text-primary" />
                  </div>
                  <h3 className="font-semibold mb-2">{feature.title}</h3>
                  <p className="text-sm text-muted-foreground">
                    {feature.description}
                  </p>
                </div>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
