import { createContext, useContext, useEffect, useMemo, useState } from "react";

type Theme = "light" | "dark";

type ThemeContextValue = {
  theme: Theme;
  toggleTheme: () => void;
};

const THEME_COOKIE = "theme";
const ONE_YEAR_SECONDS = 60 * 60 * 24 * 365;

const ThemeContext = createContext<ThemeContextValue | null>(null);

function readThemeFromCookie(): Theme {
  const cookiePart = document.cookie
    .split("; ")
    .find((part) => part.startsWith(`${THEME_COOKIE}=`));
  const value = cookiePart?.split("=")[1];
  return value === "dark" ? "dark" : "light";
}

function writeThemeCookie(theme: Theme) {
  document.cookie = `${THEME_COOKIE}=${theme}; Max-Age=${ONE_YEAR_SECONDS}; Path=/; SameSite=Lax`;
}

function applyTheme(theme: Theme) {
  const root = document.documentElement;
  if (theme === "dark") {
    root.classList.add("dark");
  } else {
    root.classList.remove("dark");
  }
}

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setTheme] = useState<Theme>(() => readThemeFromCookie());

  useEffect(() => {
    applyTheme(theme);
    writeThemeCookie(theme);
  }, [theme]);

  const value = useMemo(
    () => ({
      theme,
      toggleTheme: () => setTheme((current) => (current === "light" ? "dark" : "light")),
    }),
    [theme],
  );

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used inside ThemeProvider");
  }
  return context;
}
