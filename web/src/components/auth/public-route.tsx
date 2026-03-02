import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../../providers/auth-provider";

export function PublicRoute() {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-zinc-50 dark:bg-zinc-950">
        <p className="text-sm text-zinc-600 dark:text-zinc-400">Checking session...</p>
      </div>
    );
  }

  if (user) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}
