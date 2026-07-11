import { createContext, useContext, useEffect, useState } from "react";

export type Theme = "violet" | "indigo" | "emerald" | "rose";
export type Mode = "dark" | "light";

const themes: { id: Theme; label: string; color: string }[] = [
  { id: "violet", label: "Violet", color: "#7327f7" },
  { id: "indigo", label: "Indigo", color: "#4f46e5" },
  { id: "emerald", label: "Emerald", color: "#059669" },
];

type ThemeCtx = {
  theme: Theme;
  setTheme: (t: Theme) => void;
  themes: typeof themes;
  mode: Mode;
  setMode: (m: Mode) => void;
};

const ThemeContext = createContext<ThemeCtx | null>(null);

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const [theme, setThemeState] = useState<Theme>(
    () => (localStorage.getItem("theme") as Theme) ?? "violet"
  );
  const [mode, setModeState] = useState<Mode>(
    () => (localStorage.getItem("mode") as Mode) ?? "light"
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

  useEffect(() => {
    const root = document.documentElement;
    if (mode === "light") {
      root.removeAttribute("data-mode");
    } else {
      root.setAttribute("data-mode", mode);
    }
    localStorage.setItem("mode", mode);
  }, [mode]);

  const setTheme = (t: Theme) => setThemeState(t);
  const setMode = (m: Mode) => setModeState(m);

  return (
    <ThemeContext.Provider value={{ theme, setTheme, themes, mode, setMode }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error("useTheme must be used inside ThemeProvider");
  return ctx;
}
