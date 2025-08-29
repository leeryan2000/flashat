import { createContext, useEffect, useState } from "react";
import { api } from "../api/Api";




type User = { 
    uid: string;
    email: string;
    name: string;
}

interface AuthContextType {
    user: User | null;
    login: (user: User) => void;
    logout: () => void;
    isLoading: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        (async () => {
            try {
                const user = await api<User>("/auth/login");
                setUser(user);
            } catch (error) {
                console.error("Failed to fetch user:", error);
            } finally {
                setIsLoading(false);
            }
        })();
    }, []);

    const login = async(email: string, password: string) => {
        setIsLoading(true);
        try {
            const user = await api<User>("/auth/login", {
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
    }

    const logout = async() => {
        setIsLoading(true);
        try {
            await api("/auth/logout", { method: "POST" });
            setUser(null);
        } catch (error) {
            console.error("Logout failed:", error);
            throw error;
        } finally {
            setIsLoading(false);
        }
    }

}