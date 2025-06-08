// components/app-sidebar.tsx

"use client"; // Required for the hooks in the footer

import type { LucideIcon } from "lucide-react";
import { useEffect } from "react";
import { useTheme } from "next-themes";
import {
    Sidebar as ShadSidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarHeader,
    SidebarMenu,
    SidebarMenuBadge,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarProvider,
    SidebarTrigger,
} from "../ui/sidebar";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuSeparator,
    DropdownMenuSub,
    DropdownMenuSubContent,
    DropdownMenuSubTrigger,
    DropdownMenuTrigger,
} from "../ui/dropdown-menu";
import {
    Atom,
    Bell,
    Bookmark,
    CheckCircle2,
    ChevronRight,
    CircleUser,
    Flame,
    HelpCircle,
    LayoutGrid,
    LogOut,
    Settings,
    Sprout,
    SunDim,
    Trees,
    ChevronUp,
    Sun,
} from "lucide-react";
import { cn } from "~/lib/utils"; // Assuming this is your path to cn

type OverviewItem = {
    title: string;
    href: string;
    icon: LucideIcon;
    count: number;
};

type ProjectItem = {
    name: string;
    href: string;
    icon: LucideIcon;
    updateCount: number;
};

const overviewItems: OverviewItem[] = [
    { title: "New Updates", href: "/feed", icon: Bell, count: 12 },
    {
        title: "Subscribed",
        href: "/projects/subscribed",
        icon: Bookmark,
        count: 8,
    },
    { title: "All Projects", href: "/projects", icon: LayoutGrid, count: 42 },
];

const projectItems: ProjectItem[] = [
    {
        name: "Phoenix Initiative",
        href: "/projects/phoenix",
        icon: Flame,
        updateCount: 5,
    },
    {
        name: "QuantumLeap",
        href: "/projects/quantum",
        icon: Atom,
        updateCount: 2,
    },
    {
        name: "EcoBuilders",
        href: "/projects/eco",
        icon: Sprout,
        updateCount: 0,
    },
];

const SidebarHeaderComponent = () => (
    <SidebarHeader>
        <SidebarMenu>
            <SidebarMenuItem>
                <SidebarMenuButton asChild>
                    <a href="/" className="cursor-pointer">
                        <Trees className="h-5 w-5" />
                        <span className="font-semibold">Overstory</span>
                    </a>
                </SidebarMenuButton>
            </SidebarMenuItem>
        </SidebarMenu>
    </SidebarHeader>
);

const OverviewSection = ({ items }: { items: OverviewItem[] }) => (
    <SidebarGroup>
        <SidebarGroupLabel>Overview</SidebarGroupLabel>
        <SidebarGroupContent>
            <SidebarMenu>
                {items.map((item) => (
                    <SidebarMenuItem key={item.title}>
                        <SidebarMenuButton asChild>
                            <a href={item.href}>
                                <item.icon className="h-4 w-4" />
                                <span>{item.title}</span>
                            </a>
                        </SidebarMenuButton>
                        <SidebarMenuBadge>{item.count}</SidebarMenuBadge>
                    </SidebarMenuItem>
                ))}
            </SidebarMenu>
        </SidebarGroupContent>
    </SidebarGroup>
);

const ProjectsSection = ({ items }: { items: ProjectItem[] }) => (
    <SidebarGroup>
        <SidebarGroupLabel>Your Projects</SidebarGroupLabel>
        <SidebarGroupContent>
            <SidebarMenu>
                {items.map((project) => (
                    <SidebarMenuItem key={project.name}>
                        <SidebarMenuButton asChild>
                            <a href={project.href}>
                                <project.icon className="h-4 w-4" />
                                <span>{project.name}</span>
                            </a>
                        </SidebarMenuButton>
                        {project.updateCount > 0 && (
                            <SidebarMenuBadge>
                                {project.updateCount}
                            </SidebarMenuBadge>
                        )}
                    </SidebarMenuItem>
                ))}
            </SidebarMenu>
        </SidebarGroupContent>
    </SidebarGroup>
);

// Converted to a full component to use hooks
function SidebarFooterComponent() {
    const { theme, setTheme } = useTheme();

    useEffect(() => {
        const themeColor = theme === "dark" ? "#020817" : "#fff";
        const metaThemeColor = document.querySelector(
            "meta[name='theme-color']"
        );
        if (metaThemeColor) metaThemeColor.setAttribute("content", themeColor);
    }, [theme]);

    return (
        <SidebarFooter>
            <SidebarMenu>
                <SidebarMenuItem>
                    <DropdownMenu modal={false}>
                        <DropdownMenuTrigger asChild>
                            <SidebarMenuButton className="w-full">
                                <CircleUser className="h-5 w-5" />
                                <span>Jane Doe</span>
                                <ChevronUp className="ml-auto h-4 w-4" />
                            </SidebarMenuButton>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent
                            side="top"
                            className="w-[--radix-popper-anchor-width]"
                        >
                            <DropdownMenuItem>
                                <HelpCircle className="mr-2 h-4 w-4" />
                                <span className="pr-2">Help & Feedback</span>
                            </DropdownMenuItem>

                            <DropdownMenuSub>
                                <DropdownMenuSubTrigger>
                                    <Sun className="mr-4 h-4 w-4" />
                                    <span>Theme</span>
                                </DropdownMenuSubTrigger>
                                <DropdownMenuSubContent>
                                    <DropdownMenuItem
                                        onClick={() => setTheme("light")}
                                    >
                                        Light
                                        <CheckCircle2
                                            size={14}
                                            className={cn(
                                                "ml-auto",
                                                theme !== "light" && "hidden"
                                            )}
                                        />
                                    </DropdownMenuItem>
                                    <DropdownMenuItem
                                        onClick={() => setTheme("dark")}
                                    >
                                        Dark
                                        <CheckCircle2
                                            size={14}
                                            className={cn(
                                                "ml-auto",
                                                theme !== "dark" && "hidden"
                                            )}
                                        />
                                    </DropdownMenuItem>
                                    <DropdownMenuItem
                                        onClick={() => setTheme("system")}
                                    >
                                        System
                                        <CheckCircle2
                                            size={14}
                                            className={cn(
                                                "ml-auto",
                                                theme !== "system" && "hidden"
                                            )}
                                        />
                                    </DropdownMenuItem>
                                </DropdownMenuSubContent>
                            </DropdownMenuSub>

                            <DropdownMenuItem>
                                <Settings className="mr-2 h-4 w-4" />
                                <span>Settings</span>
                            </DropdownMenuItem>

                            <DropdownMenuSeparator />
                            <DropdownMenuItem>
                                <LogOut className="mr-2 h-4 w-4" />
                                <span>Log out</span>
                            </DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </SidebarMenuItem>
            </SidebarMenu>
        </SidebarFooter>
    );
}

export function Sidebar({ children }: { children: React.ReactNode }) {
    return (
        <SidebarProvider>
            <ShadSidebar>
                <SidebarHeaderComponent />
                <SidebarContent>
                    <OverviewSection items={overviewItems} />
                    <ProjectsSection items={projectItems} />
                </SidebarContent>
                <SidebarFooterComponent />
            </ShadSidebar>
            <main>
                <SidebarTrigger />
                {children}
            </main>
        </SidebarProvider>
    );
}
