import { createContext, useContext, useEffect, useState } from "react";

export type Theme = "violet" | "indigo" | "emerald" | "rose";

const themes: { id: Theme; label: string; color: string }[] = [
  { id: "violet", label: "Violet", color: "#7327f7" },
  { id: "indigo", label: "Indigo", color: "#4f46e5" },
  { id: "emerald", label: "Emerald", color: "#059669" },
];

type ThemeCtx = { theme: Theme; setTheme: (t: Theme) => void; themes: typeof themes };

const ThemeContext = createContext<ThemeCtx | null>(null);

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(
    () => (localStorage.getItem("theme") as Theme) ?? "violet"
  );

  useEffect(() => {
    const root = document.documentElement;
    if (theme === "violet") {
      root.removeAttribute("data-theme");
    } else {
      root.setAttribute("data-theme", theme);
    }
    localStorage.setItem("theme", theme);
  }, [theme]);

  const setTheme = (t: Theme) => setThemeState(t);

  return (
    <ThemeContext.Provider value={{ theme, setTheme, themes }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error("useTheme must be used inside ThemeProvider");
  return ctx;
}
