"use client";

import * as React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import {
    LayoutDashboard,
    School,
    Users,
    CheckCircle,
    Star,
    LogOut,
    ChevronUp,
    GraduationCap,
    Bell,
} from "lucide-react";

import { useAuth } from "@/hooks/use-auth";
import { cn } from "@/lib/utils";
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarRail,
} from "@/components/ui/sidebar";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import { ThemeToggle } from "@/components/theme-toggle";
import { USER_ROLE_LABELS } from "@/types";

// Navigation item type
interface NavItem {
    title: string;
    href: string;
    icon: React.ComponentType<{ className?: string }>;
    badge?: boolean;
}

// Navigation items for Super Admin
const superAdminNav: NavItem[] = [
    {
        title: "Dashboard",
        href: "/dashboard",
        icon: LayoutDashboard,
    },
    {
        title: "Sekolah",
        href: "/dashboard/sekolah",
        icon: School,
    },
    {
        title: "Users",
        href: "/dashboard/users",
        icon: Users,
    },
    {
        title: "Verifikasi Talenta",
        href: "/dashboard/verifikasi",
        icon: CheckCircle,
        badge: true,
    },
];

const commonNav: NavItem[] = [
    {
        title: "Talenta",
        href: "/dashboard/talenta",
        icon: Star,
    },
    {
        title: "Notifikasi",
        href: "/dashboard/notifications",
        icon: Bell,
    },
];

export function AppSidebar({ pendingCount = 0 }: { pendingCount?: number }) {
    const pathname = usePathname();
    const { user, logout } = useAuth();

    const handleLogout = async () => {
        await logout();
        window.location.href = "/login";
    };

    const getInitials = (name: string) => {
        return name
            .split(" ")
            .map((n) => n[0])
            .join("")
            .toUpperCase()
            .slice(0, 2);
    };

    // Determine which nav items to show based on role
    const navItems =
        user?.role === "super_admin" ? [...superAdminNav, ...commonNav] : commonNav;

    return (
        <Sidebar collapsible="icon">
            <SidebarHeader>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <SidebarMenuButton size="lg" asChild>
                            <Link href="/dashboard">
                                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                                    <GraduationCap className="size-4" />
                                </div>
                                <div className="grid flex-1 text-left text-sm leading-tight">
                                    <span className="truncate font-semibold">SIPODI</span>
                                    <span className="truncate text-xs text-muted-foreground">
                                        Cadbin Malang
                                    </span>
                                </div>
                            </Link>
                        </SidebarMenuButton>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarHeader>

            <SidebarContent>
                <SidebarGroup>
                    <SidebarGroupLabel>Menu</SidebarGroupLabel>
                    <SidebarGroupContent>
                        <SidebarMenu>
                            {navItems.map((item) => (
                                <SidebarMenuItem key={item.href}>
                                    <SidebarMenuButton
                                        asChild
                                        isActive={
                                            pathname === item.href ||
                                            (item.href !== "/dashboard" &&
                                                pathname.startsWith(item.href))
                                        }
                                        tooltip={item.title}
                                    >
                                        <Link href={item.href}>
                                            <item.icon />
                                            <span>{item.title}</span>
                                            {item.badge && pendingCount > 0 && (
                                                <Badge
                                                    variant="destructive"
                                                    className="ml-auto size-5 justify-center p-0 text-xs"
                                                >
                                                    {pendingCount > 99 ? "99+" : pendingCount}
                                                </Badge>
                                            )}
                                        </Link>
                                    </SidebarMenuButton>
                                </SidebarMenuItem>
                            ))}
                        </SidebarMenu>
                    </SidebarGroupContent>
                </SidebarGroup>
            </SidebarContent>

            <SidebarFooter>
                <SidebarMenu>
                    <SidebarMenuItem>
                        <DropdownMenu>
                            <DropdownMenuTrigger asChild>
                                <SidebarMenuButton
                                    size="lg"
                                    className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                                >
                                    <Avatar className="h-8 w-8 rounded-lg">
                                        <AvatarImage
                                            src={user?.photo_url || undefined}
                                            alt={user?.full_name}
                                        />
                                        <AvatarFallback className="rounded-lg">
                                            {user?.full_name ? getInitials(user.full_name) : "U"}
                                        </AvatarFallback>
                                    </Avatar>
                                    <div className="grid flex-1 text-left text-sm leading-tight">
                                        <span className="truncate font-semibold">
                                            {user?.full_name || "User"}
                                        </span>
                                        <span className="truncate text-xs text-muted-foreground">
                                            {user?.role ? USER_ROLE_LABELS[user.role] : ""}
                                        </span>
                                    </div>
                                    <ChevronUp className="ml-auto size-4" />
                                </SidebarMenuButton>
                            </DropdownMenuTrigger>
                            <DropdownMenuContent
                                className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                                side="top"
                                align="end"
                                sideOffset={4}
                            >
                                <div className="flex items-center justify-between px-2 py-1.5">
                                    <span className="text-sm font-medium">Theme</span>
                                    <ThemeToggle />
                                </div>
                                <DropdownMenuSeparator />
                                <DropdownMenuItem onClick={handleLogout}>
                                    <LogOut className="mr-2 size-4" />
                                    Logout
                                </DropdownMenuItem>
                            </DropdownMenuContent>
                        </DropdownMenu>
                    </SidebarMenuItem>
                </SidebarMenu>
            </SidebarFooter>

            <SidebarRail />
        </Sidebar>
    );
}
