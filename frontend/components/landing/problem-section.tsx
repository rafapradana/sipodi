import { FolderOpen, Search, FileText, BarChart2 } from "lucide-react";
import { Card, CardContent } from "@/components/ui/card";

const problems = [
  {
    icon: FolderOpen,
    title: "Data Tidak Terpusat",
    description:
      "Informasi GTK tersebar di berbagai dokumen dan sistem yang tidak saling terhubung.",
  },
  {
    icon: Search,
    title: "Potensi Tersembunyi",
    description:
      "Kompetensi, prestasi, dan bakat GTK tidak terdokumentasi dengan baik sehingga sulit dipetakan.",
  },
  {
    icon: FileText,
    title: "Proses Masih Manual",
    description:
      "Perencanaan dan pembinaan GTK masih mengandalkan cara konvensional yang memakan waktu.",
  },
  {
    icon: BarChart2,
    title: "Minim Data Analitik",
    description:
      "Pengambilan keputusan strategis tidak didukung oleh dashboard dan laporan yang memadai.",
  },
];

export function ProblemSection() {
  return (
    <section className="py-20 bg-muted/30">
      <div className="container mx-auto px-4">
        <div className="text-center mb-12">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Tantangan Pengelolaan Data GTK Saat Ini
          </h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-6xl mx-auto">
          {problems.map((problem) => (
            <Card
              key={problem.title}
              className="border-destructive/20 bg-destructive/5 dark:bg-destructive/10"
            >
              <CardContent className="pt-6">
                <div className="flex flex-col items-center text-center">
                  <div className="mb-4 rounded-full bg-destructive/10 p-3">
                    <problem.icon className="h-6 w-6 text-destructive" />
                  </div>
                  <h3 className="font-semibold mb-2">{problem.title}</h3>
                  <p className="text-sm text-muted-foreground">
                    {problem.description}
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
