import { CheckCircle2, Moon, SunDim } from "lucide-react";
import { useTheme } from "next-themes";
import { useEffect } from "react";
import { Button } from "~/components/ui/button";
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";
import { cn } from "~/lib/utils";

export function ThemeSwitch() {
    const { theme, setTheme } = useTheme();

    /* Update theme-color meta tag
     * when theme is updated */
    useEffect(() => {
        const themeColor = theme === "dark" ? "#020817" : "#fff";
        const metaThemeColor = document.querySelector(
            "meta[name='theme-color']"
        );
        if (metaThemeColor) metaThemeColor.setAttribute("content", themeColor);
    }, [theme]);

    return (
        <DropdownMenu modal={false}>
            <DropdownMenuTrigger asChild>
                <Button
                    variant="ghost"
                    size="icon"
                    className="scale-95 rounded-full"
                >
                    <SunDim className="size-[1.2rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90" />
                    <Moon className="absolute size-[1.2rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0" />
                    <span className="sr-only">Toggle theme</span>
                </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
                <DropdownMenuItem onClick={() => setTheme("light")}>
                    Light{" "}
                    <CheckCircle2
                        size={14}
                        className={cn("ml-auto", theme !== "light" && "hidden")}
                    />
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => setTheme("dark")}>
                    Dark
                    <CheckCircle2
                        size={14}
                        className={cn("ml-auto", theme !== "dark" && "hidden")}
                    />
                </DropdownMenuItem>
                <DropdownMenuItem onClick={() => setTheme("system")}>
                    System
                    <CheckCircle2
                        size={14}
                        className={cn(
                            "ml-auto",
                            theme !== "system" && "hidden"
                        )}
                    />
                </DropdownMenuItem>
            </DropdownMenuContent>
        </DropdownMenu>
    );
}
