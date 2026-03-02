import { Moon, Sun } from "lucide-react";
import { Button } from "../ui/button";
import { useTheme } from "../../providers/theme-provider";

export function ThemeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <Button
      type="button"
      variant="outline"
      size="icon"
      className="fixed right-4 top-4 z-50 bg-white/80 backdrop-blur dark:border-zinc-700 dark:bg-zinc-900/80 dark:hover:bg-zinc-800"
      onClick={toggleTheme}
      aria-label="Toggle dark mode"
      title={theme === "dark" ? "Switch to light mode" : "Switch to dark mode"}
    >
      {theme === "dark" ? <Sun className="size-4" /> : <Moon className="size-4" />}
    </Button>
  );
}
