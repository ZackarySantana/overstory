import type { Route } from "./+types/home";
import { Welcome } from "../welcome/welcome";
import { Button } from "~/components/ui/button";
import { ThemeSwitch } from "~/components/theme-switch";

export function meta({}: Route.MetaArgs) {
    return [
        { title: "New React Router App" },
        { name: "description", content: "Welcome to React Router!" },
    ];
}

export default function Home() {
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
