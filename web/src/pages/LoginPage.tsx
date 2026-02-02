import { useNavigate } from "react-router-dom";
import { LoginForm } from "../components/LoginForm";
import { useAuth } from "../context/AuthContext";
import { useEffect } from "react";
import { PATHS } from "../routes/paths";
import { MessageSquare } from "lucide-react";

export default function LoginPage() {
  const { isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
        // redirect the user if logged in already
        if (!isLoading && isAuthenticated) {
            navigate(PATHS.chat);
        }
    }, [isLoading, isAuthenticated, navigate]);

  return (
    // Full screen container with dark background
    <div className="min-h-screen bg-slate-950 flex flex-col items-center justify-center p-4">
      
      {/* Decorative Background Glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] bg-indigo-500/20 rounded-full blur-[100px] pointer-events-none" />

      {/* Main Card */}
      <div className="w-full max-w-md bg-slate-900 border border-slate-800 rounded-3xl p-8 shadow-2xl relative z-10 animate-in fade-in zoom-in-95 duration-300">
        
        {/* App Logo Section */}
        <div className="flex justify-center mb-8">
          <div className="h-16 w-16 bg-indigo-500/10 rounded-2xl flex items-center justify-center rotate-3 border border-indigo-500/20">
            <MessageSquare className="h-8 w-8 text-indigo-500" />
          </div>
        </div>

        {/* Form Component */}
        <LoginForm />
        
      </div>

      {/* Footer / Copyright */}
      <p className="mt-8 text-slate-600 text-sm relative z-10">
        © 2026 Flashat
      </p>
    </div>
  );
}
