import Link from "next/link";
import { ArrowRight, MessageCircle } from "lucide-react";
import { Button } from "@/components/ui/button";

export function CTASection() {
  return (
    <section className="py-20 bg-primary text-primary-foreground">
      <div className="container mx-auto px-4">
        <div className="max-w-3xl mx-auto text-center">
          <h2 className="text-3xl md:text-4xl font-bold mb-4">
            Siap Mengelola Potensi GTK dengan Lebih Baik?
          </h2>
          <p className="text-lg opacity-90 mb-8">
            Bergabunglah dengan ratusan sekolah di wilayah Malang yang telah
            menggunakan SIPODI untuk memetakan dan mengembangkan potensi GTK.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button
              size="lg"
              variant="secondary"
              asChild
            >
              <Link href="/login">
                Masuk ke SIPODI
                <ArrowRight className="ml-2 h-4 w-4" />
              </Link>
            </Button>
            <Button
              size="lg"
              variant="outline"
              className="border-primary-foreground/30 text-primary-foreground hover:bg-primary-foreground/10"
              asChild
            >
              <Link href="https://wa.me/6281234567890" target="_blank">
                <MessageCircle className="mr-2 h-4 w-4" />
                Hubungi Admin
              </Link>
            </Button>
          </div>
        </div>
      </div>
    </section>
  );
}
