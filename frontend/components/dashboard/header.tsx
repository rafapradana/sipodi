"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { ChevronRight } from "lucide-react";
import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Separator } from "@/components/ui/separator";

// Breadcrumb path mappings
const pathLabels: Record<string, string> = {
    dashboard: "Dashboard",
    sekolah: "Sekolah",
    users: "Users",
    verifikasi: "Verifikasi Talenta",
    talenta: "Talenta",
};

export function DashboardHeader() {
    const pathname = usePathname();
    const pathSegments = pathname.split("/").filter(Boolean);

    // Generate breadcrumb items
    const breadcrumbs = pathSegments.map((segment, index) => {
        const href = "/" + pathSegments.slice(0, index + 1).join("/");
        const label = pathLabels[segment] || segment;
        const isLast = index === pathSegments.length - 1;

        return {
            href,
            label,
            isLast,
        };
    });

    return (
        <header className="flex h-16 shrink-0 items-center gap-2 border-b px-4">
            <SidebarTrigger className="-ml-1" />
            <Separator orientation="vertical" className="mr-2 h-4" />
            <Breadcrumb>
                <BreadcrumbList>
                    {breadcrumbs.map((crumb, index) => (
                        <BreadcrumbItem key={crumb.href}>
                            {index > 0 && <BreadcrumbSeparator />}
                            {crumb.isLast ? (
                                <BreadcrumbPage>{crumb.label}</BreadcrumbPage>
                            ) : (
                                <BreadcrumbLink asChild>
                                    <Link href={crumb.href}>{crumb.label}</Link>
                                </BreadcrumbLink>
                            )}
                        </BreadcrumbItem>
                    ))}
                </BreadcrumbList>
            </Breadcrumb>
        </header>
    );
}
