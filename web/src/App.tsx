import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { ProtectedRoute } from "./components/auth/protected-route";
import { PublicRoute } from "./components/auth/public-route";
import { ThemedToaster } from "./components/theme/themed-toaster";
import { ThemeToggle } from "./components/theme/theme-toggle";
import { DashboardPage } from "./pages/dashboard-page";
import { LoginPage } from "./pages/login-page";
import { RegisterPage } from "./pages/register-page";
import { AuthProvider } from "./providers/auth-provider";
import { ThemeProvider } from "./providers/theme-provider";

function App() {
  return (
    <ThemeProvider>
      <BrowserRouter>
        <AuthProvider>
          <ThemeToggle />
          <Routes>
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route element={<PublicRoute />}>
              <Route path="/login" element={<LoginPage />} />
              <Route path="/register" element={<RegisterPage />} />
            </Route>
            <Route element={<ProtectedRoute />}>
              <Route path="/dashboard" element={<DashboardPage />} />
            </Route>
            <Route path="*" element={<Navigate to="/dashboard" replace />} />
          </Routes>
          <ThemedToaster />
        </AuthProvider>
      </BrowserRouter>
    </ThemeProvider>
  );
}

export default App;
