import { useEffect, useRef } from "react";
import { Navigate, Outlet, useLocation } from "react-router-dom";
import { toast } from "sonner";
import { useAuth } from "../../providers/auth-provider";

export function ProtectedRoute() {
  const { user, isLoading } = useAuth();
  const location = useLocation();
  const hasShownToast = useRef(false);

  useEffect(() => {
    if (!isLoading && !user && !hasShownToast.current) {
      hasShownToast.current = true;
      const skipUnauthToast = sessionStorage.getItem("auth:skip-unauth-toast") === "1";
      if (skipUnauthToast) {
        sessionStorage.removeItem("auth:skip-unauth-toast");
        return;
      }
      toast.error("You need to login.");
    }
  }, [isLoading, user]);

  if (isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-zinc-50">
        <p className="text-sm text-zinc-600">Checking session...</p>
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  return <Outlet />;
}
