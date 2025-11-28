const API_BASE = window.location.origin;

function toast(message) {
  const id = "toast";
  let node = document.getElementById(id);
  if (!node) {
    node = document.createElement("div");
    node.id = id;
    node.className = "toast";
    document.body.appendChild(node);
  }
  node.textContent = message;
  node.classList.add("show");
  setTimeout(() => node.classList.remove("show"), 2200);
}

async function api(path, options = {}) {
  const res = await fetch(`${API_BASE}${path}`, options);
  let body = {};
  try { body = await res.json(); } catch (_) {}
  const pickMessage = (b, statusText) => {
    if (!b) return statusText;
    if (typeof b === "string") return b;
    if (typeof b.message === "string") return b.message;
    if (typeof b.error === "string") return b.error;
    if (b.errors) return JSON.stringify(b.errors);
    return statusText;
  };
  if (!res.ok) {
    const msg = pickMessage(body, res.statusText);
    const err = new Error(msg);
    err.status = res.status;
    throw err;
  }
  return body;
}

function storeToken(key, value) { localStorage.setItem(key, value || ""); }
function getToken(key) { return localStorage.getItem(key) || ""; }

function ensureToken(key, redirectTo) {
  if (!getToken(key)) {
    if (redirectTo) window.location.href = redirectTo;
    return false;
  }
  return true;
}
