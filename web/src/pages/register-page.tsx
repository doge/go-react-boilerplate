import { useState } from "react";
import type { SubmitEvent } from "react";
import { Loader2 } from "lucide-react";
import { Link, useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { Button } from "../components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "../components/ui/card";
import { Input } from "../components/ui/input";
import { Label } from "../components/ui/label";
import { ApiError, registerUser } from "../lib/api";

export function RegisterPage() {
  const navigate = useNavigate();
  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  async function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!email || !username || !password) {
      toast.error("Email, username, and password are required.");
      return;
    }

    try {
      setIsSubmitting(true);
      setFieldErrors({});
      await registerUser(email, username, password);
      toast.success("Account created. You can now login.");
      navigate("/login", { replace: true });
    } catch (error) {
      if (error instanceof ApiError && error.fields) {
        setFieldErrors(error.fields);
      }
      const message = error instanceof Error ? error.message : "Registration failed";
      toast.error(message);
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <main className="relative min-h-screen overflow-hidden bg-zinc-100 dark:bg-zinc-950">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_right,#e4e4e7,transparent_45%),radial-gradient(circle_at_bottom_left,#d4d4d8,transparent_45%)] dark:bg-[radial-gradient(circle_at_top_right,#27272a,transparent_45%),radial-gradient(circle_at_bottom_left,#18181b,transparent_45%)]" />
      <div className="relative mx-auto flex min-h-screen w-full max-w-6xl items-center justify-center px-4 py-10">
        <Card className="w-full max-w-md dark:border-zinc-800 dark:bg-zinc-900">
          <CardHeader>
            <CardTitle>Register</CardTitle>
            <CardDescription className="dark:text-zinc-400">Create an account to access your dashboard.</CardDescription>
          </CardHeader>
          <CardContent>
            <form className="space-y-4" onSubmit={handleSubmit}>
              <div className="space-y-2">
                <Label htmlFor="email">Email</Label>
                <Input
                  id="email"
                  type="email"
                  autoComplete="email"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  className={fieldErrors.email ? "border-red-500 focus-visible:ring-red-500 dark:border-red-500" : ""}
                />
                {fieldErrors.email ? (
                  <p className="text-sm text-red-600 dark:text-red-400">{fieldErrors.email}</p>
                ) : null}
              </div>
              <div className="space-y-2">
                <Label htmlFor="username">Username</Label>
                <Input
                  id="username"
                  autoComplete="username"
                  value={username}
                  onChange={(event) => setUsername(event.target.value)}
                  className={fieldErrors.username ? "border-red-500 focus-visible:ring-red-500 dark:border-red-500" : ""}
                />
                {fieldErrors.username ? (
                  <p className="text-sm text-red-600 dark:text-red-400">{fieldErrors.username}</p>
                ) : null}
              </div>
              <div className="space-y-2">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  type="password"
                  autoComplete="new-password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  className={fieldErrors.password ? "border-red-500 focus-visible:ring-red-500 dark:border-red-500" : ""}
                />
                {fieldErrors.password ? (
                  <p className="text-sm text-red-600 dark:text-red-400">{fieldErrors.password}</p>
                ) : null}
              </div>
              <Button type="submit" className="w-full disabled:opacity-100" disabled={isSubmitting}>
                {isSubmitting ? (
                  <>
                    <Loader2 className="size-4 animate-spin" />
                    Creating account...
                  </>
                ) : (
                  "Create account"
                )}
              </Button>
            </form>

            <p className="mt-4 text-sm text-zinc-600 dark:text-zinc-400">
              Already registered?{" "}
              <Link className="font-medium text-zinc-900 underline underline-offset-4 dark:text-zinc-100" to="/login">
                Go to login
              </Link>
            </p>
          </CardContent>
        </Card>
      </div>
    </main>
  );
}
