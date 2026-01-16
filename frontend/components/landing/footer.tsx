import Link from "next/link";
import { Users, Mail, Phone, MapPin } from "lucide-react";
import { Separator } from "@/components/ui/separator";

const navigation = [
  { href: "#beranda", label: "Beranda" },
  { href: "#fitur", label: "Fitur" },
  { href: "#cara-kerja", label: "Cara Kerja" },
  { href: "/login", label: "Masuk" },
];

const legal = [
  { href: "/privacy", label: "Kebijakan Privasi" },
  { href: "/terms", label: "Syarat & Ketentuan" },
];

export function Footer() {
  return (
    <footer className="bg-card border-t">
      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          {/* Brand */}
          <div className="space-y-4">
            <Link href="/" className="flex items-center gap-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
                <Users className="h-5 w-5 text-primary-foreground" />
              </div>
              <span className="text-xl font-bold">SIPODI</span>
            </Link>
            <p className="text-sm text-muted-foreground">
              Sistem Informasi Potensi Diri GTK
            </p>
            <p className="text-xs text-muted-foreground">
              © {new Date().getFullYear()} Cabang Dinas Pendidikan Wilayah Malang
            </p>
          </div>

          {/* Navigation */}
          <div>
            <h3 className="font-semibold mb-4">Navigasi</h3>
            <ul className="space-y-2">
              {navigation.map((item) => (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {item.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          {/* Contact */}
          <div>
            <h3 className="font-semibold mb-4">Kontak</h3>
            <ul className="space-y-3">
              <li className="flex items-start gap-2 text-sm text-muted-foreground">
                <MapPin className="h-4 w-4 shrink-0 mt-0.5" />
                <span>
                  Jl. Veteran No. 1, Kota Malang,
                  <br />
                  Jawa Timur 65145
                </span>
              </li>
              <li className="flex items-center gap-2 text-sm text-muted-foreground">
                <Phone className="h-4 w-4 shrink-0" />
                <span>(0341) 123-456</span>
              </li>
              <li className="flex items-center gap-2 text-sm text-muted-foreground">
                <Mail className="h-4 w-4 shrink-0" />
                <span>sipodi@cabdinmalang.go.id</span>
              </li>
            </ul>
          </div>

          {/* Legal */}
          <div>
            <h3 className="font-semibold mb-4">Legal</h3>
            <ul className="space-y-2">
              {legal.map((item) => (
                <li key={item.href}>
                  <Link
                    href={item.href}
                    className="text-sm text-muted-foreground hover:text-foreground transition-colors"
                  >
                    {item.label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>

        <Separator className="my-8" />

        <div className="text-center text-sm text-muted-foreground">
          Dikembangkan untuk Cabang Dinas Pendidikan Wilayah Malang
        </div>
      </div>
    </footer>
  );
}
