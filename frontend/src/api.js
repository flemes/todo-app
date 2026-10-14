const BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

async function request(path, options = {}) {
  const res = await fetch(`${BASE}${path}`, {
    headers: { "Content-Type": "application/json" },
    ...options,
  });
  if (!res.ok) {
    let message = `Request failed (${res.status})`;
    try {
      message = (await res.json()).error ?? message;
    } catch {
      /* body was not JSON */
    }
    throw new Error(message);
  }
  return res.status === 204 ? null : res.json();
}

export const listTodos = () => request("/api/todos");

export const createTodo = (title) =>
  request("/api/todos", { method: "POST", body: JSON.stringify({ title }) });

export const updateTodo = (todo) =>
  request(`/api/todos/${todo.id}`, {
    method: "PUT",
    body: JSON.stringify({ title: todo.title, done: todo.done }),
  });

export const deleteTodo = (id) => request(`/api/todos/${id}`, { method: "DELETE" });
