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
import { WebSocketProvider } from "../context/WebSocketContext";
import Friendships from "../pages/Friendships";

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
              <WebSocketProvider url={import.meta.env.VITE_WS_URL}>
                <ChatProvider>
                  <Layout />
                </ChatProvider>
              </WebSocketProvider>
            </AuthCheck>
          }
        >
          <Route index element={<Chat />} />
          <Route path="friends" element={<Friendships />} />
          <Route path="profile" element={<Profile />} />
          <Route path="settings" element={<Settings />} />
        </Route>
        <Route path={PATHS.login} element={<LoginPage />} />
        {/* 404, etc. */}
      </Routes>
    </BrowserRouter>
  );
}
