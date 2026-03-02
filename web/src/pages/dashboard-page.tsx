import { LogOut } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { Button } from "../components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../components/ui/card";
import { useAuth } from "../providers/auth-provider";

export function DashboardPage() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  async function handleLogout() {
    try {
      sessionStorage.setItem("auth:skip-unauth-toast", "1");
      await logout();
      toast.success("Logged out.");
      navigate("/login", { replace: true });
    } catch (error) {
      const message = error instanceof Error ? error.message : "Logout failed";
      toast.error(message);
    }
  }

  return (
    <main className="min-h-screen bg-zinc-100 px-4 py-10 dark:bg-zinc-950">
      <div className="mx-auto flex w-full max-w-4xl flex-col gap-6">
        <header className="rounded-2xl border border-zinc-200 bg-white p-6 shadow-sm dark:border-zinc-800 dark:bg-zinc-900">
          <h1 className="text-3xl font-semibold tracking-tight text-zinc-950 dark:text-zinc-100">Dashboard</h1>
          <p className="mt-2 text-zinc-600 dark:text-zinc-400">Your home away from home.</p>
        </header>

        <div className="grid gap-4 md:grid-cols-2">
          <Card className="dark:border-zinc-800 dark:bg-zinc-900">
            <CardHeader>
              <CardTitle>Account</CardTitle>
              <CardDescription className="dark:text-zinc-400">Your account details.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-2 text-sm">
              <p>
                <span className="font-medium text-zinc-950 dark:text-zinc-100">Email:</span> {user?.email}
              </p>
              <p>
                <span className="font-medium text-zinc-950 dark:text-zinc-100">UID:</span> {user?.uid}
              </p>
            </CardContent>
          </Card>

          <Card className="dark:border-zinc-800 dark:bg-zinc-900">
            <CardHeader>
              <CardTitle>Session</CardTitle>
            </CardHeader>
            <CardContent>
              <Button type="button" variant="outline" onClick={handleLogout}>
                <LogOut className="size-4" />
                Logout
              </Button>
            </CardContent>
          </Card>
        </div>
      </div>
    </main>
  );
}
