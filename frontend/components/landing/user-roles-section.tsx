import { Building2, School, GraduationCap, Check } from "lucide-react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

const roles = [
  {
    icon: Building2,
    title: "Super Admin",
    subtitle: "Cabang Dinas Pendidikan",
    color: "bg-blue-500",
    capabilities: [
      "Kelola semua data sekolah",
      "Kelola semua user (admin & GTK)",
      "Akses dashboard wilayah",
      "Export laporan lengkap",
    ],
  },
  {
    icon: School,
    title: "Admin Sekolah",
    subtitle: "Operator Sekolah",
    color: "bg-emerald-500",
    capabilities: [
      "Kelola data GTK di sekolahnya",
      "Verifikasi talenta GTK",
      "Akses dashboard sekolah",
      "Export laporan sekolah",
    ],
  },
  {
    icon: GraduationCap,
    title: "GTK",
    subtitle: "Guru & Tenaga Kependidikan",
    color: "bg-amber-500",
    capabilities: [
      "Kelola data diri",
      "Input dan update talenta",
      "Upload dokumen pendukung",
      "Terima notifikasi verifikasi",
    ],
  },
];

export function UserRolesSection() {
  return (
    <section id="tentang" className="py-20 bg-muted/30">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Satu Sistem, Tiga Peran
          </h2>
          <p className="text-lg text-muted-foreground">
            Setiap pengguna memiliki akses dan kemampuan sesuai perannya
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-5xl mx-auto">
          {roles.map((role) => (
            <Card
              key={role.title}
              className="group hover:shadow-lg transition-all duration-300 hover:-translate-y-1 overflow-hidden"
            >
              <CardHeader className={`${role.color} text-white`}>
                <div className="flex items-center gap-3">
                  <role.icon className="h-8 w-8" />
                  <div>
                    <CardTitle className="text-lg">{role.title}</CardTitle>
                    <p className="text-sm opacity-90">{role.subtitle}</p>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="pt-6">
                <ul className="space-y-3">
                  {role.capabilities.map((capability) => (
                    <li key={capability} className="flex items-start gap-2">
                      <Check className="h-5 w-5 text-primary shrink-0 mt-0.5" />
                      <span className="text-sm">{capability}</span>
                    </li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          ))}
        </div>
      </div>
    </section>
  );
}
