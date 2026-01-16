import type { Metadata } from "next";
import { Geist, Geist_Mono, Inter } from "next/font/google";
import { ThemeProvider } from "@/components/theme-provider";
import { AuthProvider } from "@/contexts/auth-context";
import "./globals.css";

const inter = Inter({ subsets: ["latin"], variable: "--font-sans" });

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "SIPODI - Sistem Informasi Potensi Diri GTK | Cabdin Malang",
  description:
    "SIPODI adalah sistem informasi terintegrasi untuk mengelola data, kompetensi, dan talenta Guru & Tenaga Kependidikan di wilayah Malang. Kelola potensi GTK dengan lebih cerdas.",
  openGraph: {
    title: "SIPODI - Sistem Informasi Potensi Diri GTK",
    description:
      "Kelola potensi GTK dengan lebih cerdas. Sistem terintegrasi untuk Cabang Dinas Pendidikan Wilayah Malang.",
    type: "website",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="id" className={inter.variable} suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <AuthProvider>
            {children}
          </AuthProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
