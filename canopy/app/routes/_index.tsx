import type { Route } from "./+types/home";
import { Welcome } from "../welcome/welcome";
import { Button } from "~/components/ui/button";
import { ThemeSwitch } from "~/components/theme-switch";
import {
    Sidebar,
    SidebarContent,
    SidebarFooter,
    SidebarGroup,
    SidebarHeader,
    SidebarProvider,
    SidebarTrigger,
} from "~/components/ui/sidebar";

export function meta({}: Route.MetaArgs) {
    return [
        { title: "New React Router App" },
        { name: "description", content: "Welcome to React Router!" },
    ];
}

export default function Home(children: { children: React.ReactNode }) {
    return (
        <>
            <p>Hello</p>
            <Button>Hey</Button>
            <Button variant="ghost">Ghost?</Button>
            <ThemeSwitch />
            <Welcome />
        </>
    );
}
