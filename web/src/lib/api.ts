const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ??
  `${window.location.protocol}//${window.location.hostname}:8000`;

type RequestInitWithRetry = RequestInit & {
  retryOnUnauthorized?: boolean;
};

export class ApiError extends Error {
  code?: string;
  fields?: Record<string, string>;

  constructor(message: string, options?: { code?: string; fields?: Record<string, string> }) {
    super(message);
    this.name = "ApiError";
    this.code = options?.code;
    this.fields = options?.fields;
  }
}

let accessToken: string | null = null;
let refreshPromise: Promise<boolean> | null = null;

function withAuthHeader(init: RequestInit = {}): RequestInit {
  const headers = new Headers(init.headers ?? {});

  if (accessToken) {
    headers.set("Authorization", `Bearer ${accessToken}`);
  }

  return {
    ...init,
    headers,
    credentials: "include",
  };
}

async function request(path: string, init: RequestInitWithRetry = {}): Promise<Response> {
  const { retryOnUnauthorized = true, ...rest } = init;

  const response = await fetch(`${API_BASE_URL}${path}`, withAuthHeader(rest));
  if (response.status !== 401 || !retryOnUnauthorized) {
    return response;
  }

  const refreshed = await refreshAccessToken();
  if (!refreshed) {
    return response;
  }

  return fetch(`${API_BASE_URL}${path}`, withAuthHeader(rest));
}

function clearAccessToken(): void {
  accessToken = null;
}

async function parseError(response: Response): Promise<ApiError> {
  try {
    const data = (await response.json()) as {
      message?: string;
      code?: string;
      fields?: Record<string, string>;
    };
    return new ApiError(data.message?.trim() || "Request failed.", {
      code: data.code?.trim(),
      fields: data.fields,
    });
  } catch {
    return new ApiError("Request failed.");
  }
}

export async function registerUser(email: string, username: string, password: string): Promise<void> {
  const response = await request("/auth/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, username, password }),
    retryOnUnauthorized: false,
  });

  if (!response.ok) {
    throw await parseError(response);
  }
}

export async function login(username: string, password: string): Promise<{ email: string }> {
  const response = await request("/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
    retryOnUnauthorized: false,
  });

  if (!response.ok) {
    throw await parseError(response);
  }

  const data = (await response.json()) as {
    email?: string;
    Email?: string;
    accessToken?: string;
    AccessToken?: string;
  };

  const resolvedAccessToken = data.accessToken ?? data.AccessToken;
  const resolvedEmail = data.email ?? data.Email;

  if (!resolvedAccessToken) {
    throw new Error("Login succeeded but access token was missing");
  }

  accessToken = resolvedAccessToken;
  return { email: resolvedEmail ?? "" };
}

export async function refreshAccessToken(): Promise<boolean> {
  if (refreshPromise) {
    return refreshPromise;
  }

  refreshPromise = (async () => {
    const response = await fetch(`${API_BASE_URL}/auth/refresh`, {
      method: "POST",
      credentials: "include",
    });

    if (!response.ok) {
      clearAccessToken();
      return false;
    }

    const data = (await response.json()) as {
      accessToken?: string;
      AccessToken?: string;
    };
    const resolvedAccessToken = data.accessToken ?? data.AccessToken;
    if (!resolvedAccessToken) {
      clearAccessToken();
      return false;
    }
    accessToken = resolvedAccessToken;
    return true;
  })();

  try {
    return await refreshPromise;
  } finally {
    refreshPromise = null;
  }
}

export async function logout(): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/auth/logout`, {
    method: "POST",
    credentials: "include",
  });
  if (!response.ok) {
    throw await parseError(response);
  }
  clearAccessToken();
}

export async function getSession(): Promise<{ uid: string; email: string }> {
  const response = await request("/auth/session", {
    method: "GET",
  });

  if (!response.ok) {
    throw new Error("Unauthorized");
  }

  const data = (await response.json()) as {
    uid?: string;
    email?: string;
  };

  return {
    uid: data.uid ?? "",
    email: data.email ?? "",
  };
}
