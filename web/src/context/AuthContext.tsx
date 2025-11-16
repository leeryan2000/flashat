import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { api } from "../api/api";

type User = {
  uid: string;
  email: string;
  name: string;
};

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => void;
  logout: () => void;
  isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({
  children,
  blockUntilReady = true,          // <-- optional opt-out
  Fallback = DefaultAuthFallback,  // <-- optional custom loading UI
}: {
  children: React.ReactNode;
  blockUntilReady?: boolean;
  Fallback?: React.ComponentType;
}) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    (async () => {
      try {
        const user = await api<User>("/auth/me");
        setUser(user);
      } catch (error) {
        console.error("Failed to fetch user:", error);
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);
  
  const login = async (email: string, password: string) => {
    setIsLoading(true);
    try {
      const user = await api<User>("/login", {
        method: "POST",
        body: JSON.stringify({ email, password }),
      });
      setUser(user);
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const logout = async () => {
    setIsLoading(true);
    try {
      await api("/auth/logout", { method: "DELETE" });
      setUser(null);
    } catch (error) {
      console.error("Logout failed:", error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  // Memoize the context value
  const value = useMemo(
    () => ({
      user,
      // see if user is truthy, if not then its not authenticated
      isAuthenticated: !isLoading && !!user,
      login,
      logout,
      isLoading,
    }),
    [user, isLoading]
  );

  if (blockUntilReady && isLoading) {
    return <Fallback />;
  } 

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

function DefaultAuthFallback() {
  return (
    <div className="grid min-h-screen place-items-center text-slate-500">
      Loading session…
    </div>
  );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}