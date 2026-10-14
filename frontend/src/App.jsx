import { useEffect, useState } from "react";
import { listTodos, createTodo, updateTodo, deleteTodo } from "./api.js";

export default function App() {
  const [todos, setTodos] = useState([]);
  const [title, setTitle] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  // Run an API action; surface any failure in the UI.
  async function run(action) {
    setError("");
    try {
      await action();
    } catch (e) {
      setError(e.message);
    }
  }

  useEffect(() => {
    run(async () => setTodos(await listTodos())).finally(() => setLoading(false));
  }, []);

  const add = (e) => {
    e.preventDefault();
    const text = title.trim();
    if (!text) return;
    run(async () => {
      const created = await createTodo(text);
      setTodos((prev) => [...prev, created]);
      setTitle("");
    });
  };

  const toggle = (todo) =>
    run(async () => {
      const saved = await updateTodo({ ...todo, done: !todo.done });
      setTodos((prev) => prev.map((t) => (t.id === saved.id ? saved : t)));
    });

  const remove = (id) =>
    run(async () => {
      await deleteTodo(id);
      setTodos((prev) => prev.filter((t) => t.id !== id));
    });

  const remaining = todos.filter((t) => !t.done).length;

  return (
    <main className="app">
      <h1>Todo</h1>

      <form onSubmit={add} className="add">
        <input
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          placeholder="What needs to be done?"
          aria-label="New todo title"
        />
        <button type="submit">Add</button>
      </form>

      {error && <p className="error" role="alert">{error}</p>}

      {loading ? (
        <p className="muted">Loading…</p>
      ) : todos.length === 0 ? (
        <p className="muted">Nothing to do. Enjoy!</p>
      ) : (
        <ul>
          {todos.map((t) => (
            <li key={t.id} className={t.done ? "done" : ""}>
              <label>
                <input type="checkbox" checked={t.done} onChange={() => toggle(t)} />
                <span>{t.title}</span>
              </label>
              <button className="delete" onClick={() => remove(t.id)} aria-label={`Delete ${t.title}`}>
                ✕
              </button>
            </li>
          ))}
        </ul>
      )}

      <p className="muted">{remaining} item{remaining === 1 ? "" : "s"} left</p>
    </main>
  );
}
