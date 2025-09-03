import type React from "react";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import LoginPage from "../pages/LoginPage";
import { useAuth } from "../context/AuthContext";
import Chat from "../pages/Chat";
import Profile from "../pages/Profile";
import Settings from "../pages/Settings";
import Layout from "../layouts/Layout";
import { PATHS } from "./paths";
import { ChatProvider } from "../context/ChatContext";

function AuthCheck({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();
  if (isLoading) {
    return null;
  } // could show a loading spinner here
  return isAuthenticated ? <>{children}</> : <Navigate to={PATHS.login} replace />;
}

export default function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path={PATHS.chat}
          element={
            <AuthCheck>
                <ChatProvider>
                  <Layout />
                </ChatProvider>
            </AuthCheck>
          }
        >
          <Route index element={<Chat />} />
          <Route path="profile" element={<Profile />} />
          <Route path="settings" element={<Settings />} />
        </Route>
        <Route path={PATHS.login} element={<LoginPage />} />
        {/* 404, etc. */}
      </Routes>
    </BrowserRouter>
  );
}
