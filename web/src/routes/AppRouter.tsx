import type React from "react";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import LoginPage from "../pages/LoginPage";
import { useAuth } from "../context/AuthContext";
import DashBoard from "../pages/DashBoard";
import Layout from "../layouts/layout";
import Profile from "../pages/Profile";
import Settings from "../pages/Settings";

function AuthCheck({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, isLoading } = useAuth();
  if (isLoading) {
    return null;
  } // could show a loading spinner here
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />;
}

export default function AppRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/dashboard"
          element={
            <AuthCheck>
              <Layout />
            </AuthCheck>
          }
        >
          <Route index element={<DashBoard />} />
          <Route path="profile" element={<Profile />} />
          <Route path="settings" element={<Settings />} />
        </Route>
        <Route path="/login" element={<LoginPage />} />
        {/* 404, etc. */}
      </Routes>
    </BrowserRouter>
  );
}
