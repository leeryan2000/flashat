import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { api } from "../api/api";

type User = {
  uid: string;
  email: string;
  name: string;
  user_avatar_url?: string;
};

type AvatarUploadUrlResponse = {
  upload_url: string;
  object_url: string;
};

interface AuthContextType {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => void;
  logout: () => void;
  register: (username: string, email: string, password: string, code: string) => void;
  updateProfile: (updatedData: {name: string}) => void;
  getAvatarUploadUrl: () => Promise<AvatarUploadUrlResponse>;
  updateAvatarUrl: (url: string) => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isInitialLoading, setIsInitialLoading] = useState(true);
  const [isLoading, setIsLoading] = useState(false);

  useEffect(() => {
    (async () => {
      try {
        const user = await api<User>("/auth/me");
        setUser(user);
      } catch (error) {
        console.error("Failed to fetch user:", error);
        setUser(null);
      } finally {
        setIsInitialLoading(false);
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

  const register = async (name: string, email: string, password: string, code: string) => {
    try {
      await api<User>("/user/register",
        {
          method: "POST",
          body: JSON.stringify({ name, email, password, code }),
        }
      );
    } catch (error) {
      console.error("Registration failed:", error);
      throw error;
    }
  }

  const updateProfile = async (updatedData: {name: string}) => {
    try {
      const updatedUser = await api<User>("/user/auth/name", {
        method: "PUT",
        body: JSON.stringify(updatedData),
      });
      setUser(updatedUser);
    } catch (error) {
      console.error("Failed to update profile:", error);
      throw error;
    }
  };

  const getAvatarUploadUrl = async () => {
    return api<AvatarUploadUrlResponse>("/user/auth/avatar/upload-url");
  };

  const updateAvatarUrl = async (url: string) => {
    try {
      const updatedUser = await api<User>("/user/auth/avatar", {
        method: "PUT",
        body: JSON.stringify({ avatar_url: url }),
      });
      setUser(updatedUser);
    } catch (error) {
      console.error("Failed to update avatar:", error);
      throw error;
    }
  };

  // Memoize the context value
  const value = useMemo(
    () => ({
      user,
      // see if user is truthy, if not then its not authenticated
      isAuthenticated: !isInitialLoading && !!user,
      login,
      logout,
      register,
      isLoading,
      updateProfile,
      getAvatarUploadUrl,
      updateAvatarUrl
    }),
    [user, isInitialLoading, isLoading, login, logout, register]
  );

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuth must be used within an AuthProvider");
    }
    return context;
}
