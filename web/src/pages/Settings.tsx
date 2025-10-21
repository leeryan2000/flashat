import { useAuth } from "../context/AuthContext";

export default function Settings() {
  const {logout} = useAuth();

  return (
    <div>
      <button onClick={logout} className="mt-4 w-full rounded-xl bg-emerald-500 text-white py-2 font-semibold shadow-lg hover:shadow-emerald-300/50 hover:brightness-110 active:scale-95 transition">Logout</button>
      <h1>This is Settings</h1>
    </div>
  );
}