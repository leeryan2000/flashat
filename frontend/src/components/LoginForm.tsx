import { useState } from "react";
import { useAuth } from "../context/AuthContext";
import { Link, useNavigate } from "react-router-dom";
import { PATHS } from "../routes/paths";
import { ArrowRight, Loader2, Lock, Mail } from "lucide-react";

export function LoginForm() {
  const { login, isLoading } = useAuth();
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    try {
      await login(email, password);
      navigate(PATHS.chat);
    } catch (error) {
      setError(error instanceof Error ? error.message : "Login failed");
    }
  };

  return (
    <div className="w-full max-w-md">
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold mb-2" style={{ color: "var(--text)" }}>Flashat</h2>
        <p className="text-sm" style={{ color: "var(--text-soft)" }}>Welcome back</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">

        {error && (
          <div
            className="border rounded-lg p-3 text-sm text-center animate-in fade-in slide-in-from-top-2"
            style={{ background: "var(--danger-bg)", borderColor: "var(--danger-bg-hover)", color: "var(--danger-text)" }}
          >
            {error}
          </div>
        )}

        <div className="space-y-4">
          {/* Email */}
          <div className="relative group">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Mail className="h-5 w-5 transition-colors" style={{ color: "var(--text-faint)" }} />
            </div>
            <input
              type="text"
              placeholder="Email address"
              className="w-full border rounded-xl py-3 pl-10 pr-4 outline-none focus:ring-2 focus:border-transparent transition-all placeholder:text-[color:var(--text-faint)]"
              style={{ background: "var(--input-bg)", borderColor: "var(--input-border)", color: "var(--text)", "--tw-ring-color": "var(--primary)" } as React.CSSProperties}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          {/* Password */}
          <div className="relative group">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Lock className="h-5 w-5 transition-colors" style={{ color: "var(--text-faint)" }} />
            </div>
            <input
              type="password"
              placeholder="Password"
              className="w-full border rounded-xl py-3 pl-10 pr-4 outline-none focus:ring-2 focus:border-transparent transition-all placeholder:text-[color:var(--text-faint)]"
              style={{ background: "var(--input-bg)", borderColor: "var(--input-border)", color: "var(--text)", "--tw-ring-color": "var(--primary)" } as React.CSSProperties}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>
        </div>

        {/* Submit */}
        <button
          type="submit"
          disabled={isLoading}
          className="w-full text-white font-semibold py-3 rounded-xl transition-all shadow-lg active:scale-[0.98] disabled:opacity-70 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          style={{ background: "var(--primary)", boxShadow: "0 4px 20px var(--primary-shadow)" }}
        >
          {isLoading ? (
            <><Loader2 className="h-5 w-5 animate-spin" /><span>Signing in...</span></>
          ) : (
            <><span>Sign In</span><ArrowRight className="h-5 w-5" /></>
          )}
        </button>

        <div className="text-center text-sm" style={{ color: "var(--text-faint)" }}>
          Don't have an account?{" "}
          <Link
            to={PATHS.register}
            className="font-medium transition-colors hover:opacity-80"
            style={{ color: "var(--primary)" }}
          >
            Sign up
          </Link>
        </div>
      </form>
    </div>
  );
}
