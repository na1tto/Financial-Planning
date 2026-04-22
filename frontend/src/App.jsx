import { useEffect, useMemo, useState } from "react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  Line,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis
} from "recharts";
import { api } from "./api";

const BRL = new Intl.NumberFormat("pt-BR", {
  style: "currency",
  currency: "BRL"
});

const today = new Date().toISOString().slice(0, 10);

const emptyTransaction = {
  kind: "expense",
  description: "",
  amount: "",
  dueDate: today
};

function App() {
  const [loading, setLoading] = useState(true);
  const [user, setUser] = useState(null);
  const [authMode, setAuthMode] = useState("login");
  const [authForm, setAuthForm] = useState({ name: "", email: "", password: "" });
  const [transactionForm, setTransactionForm] = useState(emptyTransaction);
  const [transactions, setTransactions] = useState([]);
  const [forecast, setForecast] = useState([]);
  const [error, setError] = useState("");

  useEffect(() => {
    initialize();
  }, []);

  async function initialize() {
    setLoading(true);
    setError("");

    try {
      const me = await api.me();
      setUser(me.user);
      await loadFinancialData();
    } catch {
      setUser(null);
    } finally {
      setLoading(false);
    }
  }

  async function loadFinancialData() {
    const [transactionsResponse, forecastResponse] = await Promise.all([
      api.listTransactions(),
      api.monthlyForecast(12)
    ]);

    setTransactions(transactionsResponse.transactions ?? []);
    setForecast(forecastResponse.data ?? []);
  }

  async function handleAuthSubmit(event) {
    event.preventDefault();
    setError("");

    try {
      if (authMode === "register") {
        const response = await api.register({
          name: authForm.name,
          email: authForm.email,
          password: authForm.password
        });
        setUser(response.user);
      } else {
        const response = await api.login({
          email: authForm.email,
          password: authForm.password
        });
        setUser(response.user);
      }

      setAuthForm({ name: "", email: "", password: "" });
      await loadFinancialData();
    } catch (requestError) {
      setError(requestError.message);
    }
  }

  async function handleLogout() {
    setError("");
    try {
      await api.logout();
    } finally {
      setUser(null);
      setTransactions([]);
      setForecast([]);
    }
  }

  async function handleTransactionSubmit(event) {
    event.preventDefault();
    setError("");

    try {
      await api.createTransaction({
        ...transactionForm,
        amount: Number(transactionForm.amount)
      });
      setTransactionForm(emptyTransaction);
      await loadFinancialData();
    } catch (requestError) {
      setError(requestError.message);
    }
  }

  async function handleDeleteTransaction(id) {
    setError("");
    try {
      await api.deleteTransaction(id);
      await loadFinancialData();
    } catch (requestError) {
      setError(requestError.message);
    }
  }

  const totals = useMemo(() => {
    return forecast.reduce(
      (accumulator, current) => {
        accumulator.income += current.income;
        accumulator.expense += current.expense;
        accumulator.net += current.net;
        return accumulator;
      },
      { income: 0, expense: 0, net: 0 }
    );
  }, [forecast]);

  if (loading) {
    return (
      <main className="page">
        <section className="shell shell-center">
          <h1>Planejamento Financeiro</h1>
          <p>Carregando dados...</p>
        </section>
      </main>
    );
  }

  return (
    <main className="page">
      <section className="shell">
        <header className="top">
          <div>
            <h1>Planejamento Financeiro</h1>
            <p>Controle de receitas e despesas futuras com projeções mensais.</p>
          </div>
          {user ? (
            <button className="button ghost" onClick={handleLogout} type="button">
              Sair
            </button>
          ) : null}
        </header>

        {error ? <p className="error">{error}</p> : null}

        {!user ? (
          <section className="card auth">
            <div className="auth-tabs">
              <button
                className={`tab ${authMode === "login" ? "active" : ""}`}
                onClick={() => setAuthMode("login")}
                type="button"
              >
                Entrar
              </button>
              <button
                className={`tab ${authMode === "register" ? "active" : ""}`}
                onClick={() => setAuthMode("register")}
                type="button"
              >
                Criar conta
              </button>
            </div>
            <form className="form" onSubmit={handleAuthSubmit}>
              {authMode === "register" ? (
                <label>
                  Nome
                  <input
                    required
                    value={authForm.name}
                    onChange={(event) => setAuthForm((state) => ({ ...state, name: event.target.value }))}
                  />
                </label>
              ) : null}

              <label>
                E-mail
                <input
                  required
                  type="email"
                  value={authForm.email}
                  onChange={(event) => setAuthForm((state) => ({ ...state, email: event.target.value }))}
                />
              </label>

              <label>
                Senha
                <input
                  required
                  minLength={8}
                  type="password"
                  value={authForm.password}
                  onChange={(event) => setAuthForm((state) => ({ ...state, password: event.target.value }))}
                />
              </label>

              <button className="button" type="submit">
                {authMode === "register" ? "Criar conta e entrar" : "Entrar"}
              </button>
            </form>
          </section>
        ) : (
          <section className="dashboard">
            <div className="grid">
              <article className="card">
                <h2>Nova movimentação futura</h2>
                <form className="form" onSubmit={handleTransactionSubmit}>
                  <label>
                    Tipo
                    <select
                      value={transactionForm.kind}
                      onChange={(event) => setTransactionForm((state) => ({ ...state, kind: event.target.value }))}
                    >
                      <option value="expense">Despesa</option>
                      <option value="income">Receita</option>
                    </select>
                  </label>

                  <label>
                    Descrição
                    <input
                      required
                      minLength={3}
                      value={transactionForm.description}
                      onChange={(event) =>
                        setTransactionForm((state) => ({ ...state, description: event.target.value }))
                      }
                    />
                  </label>

                  <label>
                    Valor (R$)
                    <input
                      required
                      min="0.01"
                      step="0.01"
                      type="number"
                      value={transactionForm.amount}
                      onChange={(event) => setTransactionForm((state) => ({ ...state, amount: event.target.value }))}
                    />
                  </label>

                  <label>
                    Data prevista
                    <input
                      required
                      type="date"
                      value={transactionForm.dueDate}
                      onChange={(event) => setTransactionForm((state) => ({ ...state, dueDate: event.target.value }))}
                    />
                  </label>

                  <button className="button" type="submit">
                    Salvar movimentação
                  </button>
                </form>
              </article>

              <article className="card">
                <h2>Resumo dos próximos 12 meses</h2>
                <div className="stats">
                  <div className="stat">
                    <span>Receitas</span>
                    <strong>{BRL.format(totals.income)}</strong>
                  </div>
                  <div className="stat">
                    <span>Despesas</span>
                    <strong>{BRL.format(totals.expense)}</strong>
                  </div>
                  <div className="stat">
                    <span>Saldo projetado</span>
                    <strong>{BRL.format(totals.net)}</strong>
                  </div>
                </div>
              </article>
            </div>

            <article className="card">
              <h2>Projeção mensal</h2>
              <div className="chart">
                <ResponsiveContainer width="100%" height={320}>
                  <BarChart data={forecast} margin={{ top: 16, right: 12, left: 0, bottom: 0 }}>
                    <CartesianGrid strokeDasharray="4 4" />
                    <XAxis dataKey="month" />
                    <YAxis tickFormatter={(value) => BRL.format(value)} />
                    <Tooltip formatter={(value) => BRL.format(value)} />
                    <Legend />
                    <Bar dataKey="income" name="Receita" fill="#14813b" />
                    <Bar dataKey="expense" name="Despesa" fill="#ba2d2d" />
                    <Line type="monotone" dataKey="net" name="Saldo" stroke="#1f4cbf" strokeWidth={3} />
                  </BarChart>
                </ResponsiveContainer>
              </div>
            </article>

            <article className="card">
              <h2>Lançamentos futuros</h2>
              {transactions.length === 0 ? (
                <p>Nenhum lançamento ainda. Comece cadastrando o primeiro.</p>
              ) : (
                <div className="table-wrapper">
                  <table>
                    <thead>
                      <tr>
                        <th>Data</th>
                        <th>Tipo</th>
                        <th>Descrição</th>
                        <th>Valor</th>
                        <th>Ação</th>
                      </tr>
                    </thead>
                    <tbody>
                      {transactions.map((transaction) => (
                        <tr key={transaction.id}>
                          <td>{transaction.dueDate}</td>
                          <td>{transaction.kind === "income" ? "Receita" : "Despesa"}</td>
                          <td>{transaction.description}</td>
                          <td>{BRL.format(transaction.amount)}</td>
                          <td>
                            <button
                              className="button danger"
                              onClick={() => handleDeleteTransaction(transaction.id)}
                              type="button"
                            >
                              Excluir
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </article>
          </section>
        )}
      </section>
    </main>
  );
}

export default App;
