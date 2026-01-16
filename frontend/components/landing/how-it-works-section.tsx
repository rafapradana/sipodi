import { Lock, PenLine, CheckCircle } from "lucide-react";

const steps = [
  {
    number: "01",
    icon: Lock,
    title: "Login ke Sistem",
    description:
      "Masuk menggunakan akun yang telah didaftarkan oleh admin. Setiap role memiliki akses sesuai kewenangannya.",
  },
  {
    number: "02",
    icon: PenLine,
    title: "Lengkapi Data",
    description:
      "GTK mengisi data diri dan talenta. Admin sekolah mengelola data GTK di sekolahnya.",
  },
  {
    number: "03",
    icon: CheckCircle,
    title: "Verifikasi & Pantau",
    description:
      "Admin memverifikasi data talenta. Dashboard menampilkan statistik dan laporan secara real-time.",
  },
];

export function HowItWorksSection() {
  return (
    <section id="cara-kerja" className="py-20">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Cara Kerja SIPODI
          </h2>
          <p className="text-lg text-muted-foreground">
            Tiga langkah sederhana untuk mulai mengelola data GTK
          </p>
        </div>

        <div className="max-w-5xl mx-auto">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 relative">
            {/* Connector line - desktop only */}
            <div className="hidden md:block absolute top-16 left-1/6 right-1/6 h-0.5 bg-border" />

            {steps.map((step, index) => (
              <div key={step.number} className="relative">
                {/* Vertical connector - mobile only */}
                {index < steps.length - 1 && (
                  <div className="md:hidden absolute left-1/2 top-32 h-16 w-0.5 bg-border -translate-x-1/2" />
                )}

                <div className="flex flex-col items-center text-center">
                  {/* Number badge */}
                  <div className="relative z-10 mb-4">
                    <div className="flex h-16 w-16 items-center justify-center rounded-full bg-primary text-primary-foreground text-2xl font-bold">
                      {step.number}
                    </div>
                  </div>

                  {/* Icon */}
                  <div className="mb-4 rounded-full bg-muted p-3">
                    <step.icon className="h-6 w-6 text-muted-foreground" />
                  </div>

                  {/* Content */}
                  <h3 className="font-semibold mb-2">{step.title}</h3>
                  <p className="text-sm text-muted-foreground max-w-xs">
                    {step.description}
                  </p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
