const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8080/api/v1";

async function request(path, options = {}) {
  const response = await fetch(`${API_BASE}${path}`, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(options.headers ?? {})
    },
    ...options
  });

  if (response.status === 204) {
    return null;
  }

  const data = await response.json().catch(() => ({}));
  if (!response.ok) {
    const message = data.error ?? "Erro inesperado ao processar a requisição.";
    throw new Error(message);
  }

  return data;
}

export const api = {
  register(payload) {
    return request("/auth/register", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },
  login(payload) {
    return request("/auth/login", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },
  logout() {
    return request("/auth/logout", { method: "POST" });
  },
  me() {
    return request("/auth/me");
  },
  listTransactions() {
    return request("/transactions");
  },
  createTransaction(payload) {
    return request("/transactions", {
      method: "POST",
      body: JSON.stringify(payload)
    });
  },
  deleteTransaction(id) {
    return request(`/transactions/${id}`, {
      method: "DELETE"
    });
  },
  monthlyForecast(months = 12) {
    return request(`/forecast/monthly?months=${months}`);
  }
};
