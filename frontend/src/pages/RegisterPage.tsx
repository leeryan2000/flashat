import { useNavigate } from "react-router-dom";
import { RegisterForm } from "../components/RegisterForm";
import { useAuth } from "../context/AuthContext";
import { useEffect } from "react";
import { PATHS } from "../routes/paths";
import { MessageSquare } from "lucide-react";

export default function RegisterPage() {
  const { isAuthenticated, isLoading } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    if (!isLoading && isAuthenticated) navigate(PATHS.chat);
  }, [isLoading, isAuthenticated, navigate]);

  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4 relative overflow-hidden" style={{ background: "var(--sidebar-bg)" }}>

      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] rounded-full blur-[120px] pointer-events-none opacity-30" style={{ background: "var(--primary)" }} />

      <div className="w-full max-w-md rounded-3xl p-8 shadow-2xl relative z-10 animate-in fade-in zoom-in-95 duration-300 border" style={{ background: "var(--sidebar-item)", borderColor: "color-mix(in srgb, var(--primary) 20%, transparent)" }}>
        <div className="flex justify-center mb-8">
          <div className="h-16 w-16 rounded-2xl flex items-center justify-center rotate-3 border" style={{ background: "color-mix(in srgb, var(--primary) 15%, transparent)", borderColor: "color-mix(in srgb, var(--primary) 30%, transparent)" }}>
            <MessageSquare className="h-8 w-8" style={{ color: "var(--primary)" }} />
          </div>
        </div>
        <RegisterForm />
      </div>

      <p className="mt-8 text-sm relative z-10" style={{ color: "var(--text-faint)" }}>© 2026 Flashat</p>
    </div>
  );
}
