import { useNavigate } from "react-router-dom";
import { RegisterForm } from "../components/RegisterForm"; // Import the form above
import { useAuth } from "../context/AuthContext";
import { useEffect } from "react";
import { PATHS } from "../routes/paths";
import { MessageSquare } from "lucide-react";

export default function RegisterPage() {
  const { isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    // Redirect if already logged in
    if (!isLoading && isAuthenticated) {
      navigate(PATHS.chat);
    }
  }, [isLoading, isAuthenticated, navigate]);

  return (
    <div className="min-h-screen bg-slate-950 flex flex-col items-center justify-center p-4 relative overflow-hidden">
      
      {/* Background Glow Effect (Same as Login) */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] bg-indigo-500/20 rounded-full blur-[100px] pointer-events-none" />

      {/* Main Card */}
      <div className="w-full max-w-md bg-slate-900 border border-slate-800 rounded-3xl p-8 shadow-2xl relative z-10 animate-in fade-in zoom-in-95 duration-300">
        
        {/* App Logo */}
        <div className="flex justify-center mb-8">
          <div className="h-16 w-16 bg-indigo-500/10 rounded-2xl flex items-center justify-center rotate-3 border border-indigo-500/20 shadow-inner">
            <MessageSquare className="h-8 w-8 text-indigo-500" />
          </div>
        </div>

        {/* The Form */}
        <RegisterForm />
        
      </div>

      {/* Footer */}
      <p className="mt-8 text-slate-600 text-sm relative z-10">
        © 2026 Flashat
      </p>
    </div>
  );
}