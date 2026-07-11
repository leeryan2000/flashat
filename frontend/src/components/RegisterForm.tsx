import { useState } from "react";
import { useAuth } from "../context/AuthContext";
import { useNavigate, Link } from "react-router-dom";
import { PATHS } from "../routes/paths";
import { Mail, Lock, User, KeyRound, Loader2, ArrowRight } from "lucide-react";

const inputClass = "w-full border rounded-xl py-3 pl-10 pr-4 outline-none focus:ring-2 focus:border-transparent transition-all placeholder:text-[color:var(--text-faint)]";
const inputStyle = { background: "var(--input-bg)", borderColor: "var(--input-border)", color: "var(--text)", "--tw-ring-color": "var(--primary)" } as React.CSSProperties;

export function RegisterForm() {
  const { register, isLoading } = useAuth();
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [code, setCode] = useState("");
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    if (password.length < 6) { setError("Password must be at least 6 characters"); return; }
    try {
      await register(name, email, password, code);
      navigate(PATHS.login);
    } catch (err: any) {
      setError(err.message || "Registration failed. Please try again.");
    }
  };

  return (
    <div className="w-full max-w-md">
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold mb-2" style={{ color: "var(--text)" }}>Create an account</h2>
        <p style={{ color: "var(--text-soft)" }}>Join us to start chatting with friends</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-5">
        {error && (
          <div
            className="border rounded-lg p-3 text-sm text-center animate-in fade-in slide-in-from-top-2"
            style={{ background: "var(--danger-bg)", borderColor: "var(--danger-bg-hover)", color: "var(--danger-text)" }}
          >
            {error}
          </div>
        )}

        {[
          { icon: User, placeholder: "Username", type: "text", value: name, onChange: setName },
          { icon: Mail, placeholder: "Email address", type: "email", value: email, onChange: setEmail },
          { icon: Lock, placeholder: "Password", type: "password", value: password, onChange: setPassword },
          { icon: KeyRound, placeholder: "Invitation code", type: "text", value: code, onChange: setCode },
        ].map(({ icon: Icon, placeholder, type, value, onChange }) => (
          <div key={placeholder} className="relative">
            <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <Icon className="h-5 w-5" style={{ color: "var(--text-faint)" }} />
            </div>
            <input
              type={type}
              placeholder={placeholder}
              className={inputClass}
              style={inputStyle}
              value={value}
              onChange={(e) => onChange(e.target.value)}
              required
            />
          </div>
        ))}

        <button
          type="submit"
          disabled={isLoading}
          className="w-full text-white font-semibold py-3 rounded-xl transition-all shadow-lg active:scale-[0.98] disabled:opacity-70 disabled:cursor-not-allowed flex items-center justify-center gap-2 mt-2"
          style={{ background: "var(--primary)", boxShadow: "0 4px 20px var(--primary-shadow)" }}
        >
          {isLoading ? (
            <><Loader2 className="h-5 w-5 animate-spin" /><span>Creating account...</span></>
          ) : (
            <><span>Sign Up</span><ArrowRight className="h-5 w-5" /></>
          )}
        </button>

        <div className="text-center text-sm mt-6" style={{ color: "var(--text-faint)" }}>
          Already have an account?{" "}
          <Link to={PATHS.login} className="font-medium transition-colors hover:opacity-80" style={{ color: "var(--primary)" }}>
            Sign in
          </Link>
        </div>
      </form>
    </div>
  );
}
