// Ensure admin is authenticated; otherwise redirect to login (do not throw to keep logout handler usable)
if (!ensureToken("adminAccessToken", "./admin-login.html")) {
  window.location.href = "./admin-login.html";
}

let cachedUsers = [];
let userPage = 1;
const USERS_PAGE_SIZE = 10;

const adminInfo = JSON.parse(localStorage.getItem("adminInfo") || "{}");
if (adminInfo.email) {
  const el = document.getElementById("admin-info");
  if (el) el.innerText = `${adminInfo.name || "Admin"} (${adminInfo.email})`;
}

function setView(view) {
  ["users","logs","access","requests"].forEach(v => {
    const panel = document.getElementById("view-" + v);
    const nav = document.getElementById("nav-" + v);
    if (panel) panel.style.display = v === view ? "block" : "none";
    if (nav) nav.classList.toggle("active", v === view);
  });
  if (view === "users") loadUsers();
  if (view === "logs") loadLogs();
  if (view === "requests") loadRequests();
  if (view === "access") loadAdmins();
}

async function loadUsers() {
  try {
    const search = document.getElementById("search-users").value;
    renderUserDetail(null);
    const res = await api(`/v1/admin/users${search ? `?search=${encodeURIComponent(search)}` : ""}`, {
      headers: { Authorization: `Bearer ${getToken("adminAccessToken")}` }
    });
    cachedUsers = res.data || [];
    userPage = 1;
    renderUsersTable(cachedUsers);
  } catch (err) {
    handleAuthError(err);
    document.getElementById("users-table").innerText = "Error: " + err.message;
    toast("Failed to load users: " + err.message);
  }
}

async function toggleUserStatus(userId, status) {
  try {
    await api(`/v1/admin/users/${userId}/status`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${getToken("adminAccessToken")}`
      },
      body: JSON.stringify({ status })
    });
    toast(`User ${status}`);
    loadUsers();
  } catch (err) {
    toast("Failed: " + err.message);
  }
}

async function loadLogs() {
  try {
    const res = await api("/v1/admin/logs?limit=100", {
      headers: { Authorization: `Bearer ${getToken("adminAccessToken")}` }
    });
    renderLogsTable(res.data || []);
  } catch (err) {
    handleAuthError(err);
    document.getElementById("logs-table").innerText = "Error: " + err.message;
  }
}

async function createAdmin() {
  try {
    const body = {
      email: document.getElementById("new-admin-email").value,
      name: document.getElementById("new-admin-name").value,
      password: document.getElementById("new-admin-password").value
    };
    const res = await api("/v1/admin/admins", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${getToken("adminAccessToken")}`
      },
      body: JSON.stringify(body)
    });
    toast("Admin created: " + res.data.email);
    loadUsers();
    loadAdmins();
  } catch (err) {
    toast("Create admin failed: " + err.message);
  }
}

async function loadAdmins() {
  try {
    const search = document.getElementById("search-admins").value;
    const res = await api(`/v1/admin/admins${search ? `?search=${encodeURIComponent(search)}` : ""}`, {
      headers: { Authorization: `Bearer ${getToken("adminAccessToken")}` }
    });
    renderAdminsTable(res.data || []);
  } catch (err) {
    handleAuthError(err);
    document.getElementById("admins-table").innerText = "Error: " + err.message;
  }
}

async function loadRequests() {
  try {
    const res = await api("/v1/admin/requests", {
      headers: { Authorization: `Bearer ${getToken("adminAccessToken")}` }
    });
    renderRequestsTable(res.data || []);
  } catch (err) {
    handleAuthError(err);
    document.getElementById("requests-table").innerText = "Error: " + err.message;
  }
}

async function updateRequest(id, status) {
  try {
    const feedback = prompt(`Optional feedback for ${status}:`) || null;
    await api(`/v1/admin/requests/${id}/status`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${getToken("adminAccessToken")}`
      },
      body: JSON.stringify({ status, feedback })
    });
    toast(`Request ${status}`);
    loadRequests();
  } catch (err) {
    toast("Failed: " + err.message);
  }
}

async function viewUserInfo(userId) {
  const cached = cachedUsers.find(u => u.id === userId);
  if (cached) {
    renderUserDetail(cached);
  } else {
    renderUserDetail(null, null, true);
  }

  try {
    const res = await api(`/v1/admin/users/${userId}`, {
      headers: { Authorization: `Bearer ${getToken("adminAccessToken")}` }
    });
    renderUserDetail(res.data);
  } catch (err) {
    handleAuthError(err);
    renderUserDetail(null, err.message);
    toast("Failed to load user: " + err.message);
  }
}

function logout() {
  storeToken("adminAccessToken", "");
  localStorage.removeItem("adminInfo");
  window.location.href = "./index.html";
}

function handleAuthError(err) {
  if (err && (err.status === 401 || err.status === 403)) {
    toast("Session expired. Please log in again.");
    logout();
  }
}

function renderUsersTable(list) {
  const total = (list || []).length;
  const totalPages = Math.max(1, Math.ceil(total / USERS_PAGE_SIZE));
  if (userPage > totalPages) userPage = totalPages;
  const start = (userPage - 1) * USERS_PAGE_SIZE;
  const paged = (list || []).slice(start, start + USERS_PAGE_SIZE);

  const rows = paged.map(u => {
    const nextStatus = u.status === "disabled" ? "active" : "disabled";
    const btn = `<div style="display:flex; gap:6px; flex-wrap:wrap;">
      <button class="btn btn-ghost" style="padding:6px 10px;" onclick="viewUserInfo('${u.id}')">View info</button>
      <button class="btn" style="padding:6px 10px;" onclick="toggleUserStatus('${u.id}','${nextStatus}')">${nextStatus === "active" ? "Enable" : "Disable"}</button>
    </div>`;
    return `<tr>
      <td>${u.name || ""}</td>
      <td>${u.email}</td>
      <td><span class="status ${u.status}">${u.status}</span></td>
      <td>${u.role || "user"}</td>
      <td>${btn}</td>
    </tr>`;
  }).join("");
  const pager = `
    <div style="margin:6px 0;color:var(--muted); display:flex; justify-content:space-between; align-items:center; flex-wrap:wrap; gap:8px;">
      <div>Total: ${total}</div>
      <div style="display:flex; gap:6px; align-items:center;">
        <button class="btn btn-ghost" style="padding:6px 10px;" onclick="changeUserPage(-1)" ${userPage <= 1 ? "disabled" : ""}>Prev</button>
        <span>Page ${userPage} of ${totalPages}</span>
        <button class="btn btn-ghost" style="padding:6px 10px;" onclick="changeUserPage(1)" ${userPage >= totalPages ? "disabled" : ""}>Next</button>
      </div>
    </div>`;

  document.getElementById("users-table").innerHTML = `
    ${pager}
    <table>
      <tr><th>Name</th><th>Email</th><th>Status</th><th>Role</th><th>Action</th></tr>
      ${rows || "<tr><td colspan='5'>No users found</td></tr>"}
    </table>
    ${pager}`;
}

function changeUserPage(delta) {
  const totalPages = Math.max(1, Math.ceil((cachedUsers.length || 0) / USERS_PAGE_SIZE));
  userPage = Math.min(totalPages, Math.max(1, userPage + delta));
  renderUsersTable(cachedUsers);
}

function renderUserDetail(user, errorMessage, loading) {
  const target = document.getElementById("user-detail-body");
  if (!target) return;

  if (!user) {
    target.innerHTML = errorMessage
      ? `<div style="color:var(--danger);">Error: ${errorMessage}</div>`
      : loading
        ? "Loading user..."
        : "Select a user to view details.";
    return;
  }

  const info = [
    ["Name", user.name || ""],
    ["Email", user.email],
    ["Role", user.role || "user"],
    ["Status", user.status],
    ["First name", user.first_name || ""],
    ["Last name", user.last_name || ""],
    ["Phone", user.phone || ""],
    ["Signup source", user.signup_source || ""],
    ["Created at", user.created_at || ""],
    ["Updated at", user.updated_at || ""],
    ["Last login", user.last_login_at || ""],
    ["Email verified", user.email_verified ? "Yes" : "No"],
    ["Phone verified", user.phone_verified ? "Yes" : "No"]
  ];

  target.innerHTML = `<div class="grid-2">
    ${info.map(([k, v]) => `<div><div style="font-weight:600;">${k}</div><div style="color:var(--muted); word-break:break-all;">${v || "-"}</div></div>`).join("")}
  </div>`;
}

function renderLogsTable(list) {
  const rows = (list || []).map(log => `
    <tr>
      <td>${log.action}</td>
      <td>${log.actor_role || ""}</td>
      <td>${log.actor_id || ""}</td>
      <td>${log.target_type || ""}</td>
      <td>${log.target_id || ""}</td>
      <td>${log.created_at || ""}</td>
    </tr>`).join("");
  document.getElementById("logs-table").innerHTML = `<table>
    <tr><th>Action</th><th>Role</th><th>Actor</th><th>Target Type</th><th>Target ID</th><th>When</th></tr>
    ${rows || "<tr><td colspan='6'>No logs</td></tr>"}
  </table>`;
}

function renderAdminsTable(admins) {
  const rows = (admins || []).map(a => `
    <tr>
      <td>${a.name || ""}</td>
      <td>${a.email}</td>
      <td><span class="status ${a.status}">${a.status}</span></td>
    </tr>`).join("");
  document.getElementById("admins-table").innerHTML = `<table>
    <tr><th>Name</th><th>Email</th><th>Status</th></tr>
    ${rows || "<tr><td colspan='3'>No admins found</td></tr>"}
  </table>`;
}

function renderRequestsTable(list) {
  const rows = (list || []).map(r => `
    <tr>
      <td>${r.title || ""}</td>
      <td>${r.details || ""}</td>
      <td>${r.user_id}</td>
      <td><span class="status ${r.status}">${r.status}</span></td>
      <td>${r.feedback || ""}</td>
      <td style="display:flex; gap:6px; flex-wrap:wrap;">
        <button class="btn" style="padding:6px 10px;" onclick="updateRequest('${r.id}','approved')">Approve</button>
        <button class="btn" style="padding:6px 10px;" onclick="updateRequest('${r.id}','pending')">Pending</button>
        <button class="btn" style="padding:6px 10px;" onclick="updateRequest('${r.id}','rejected')">Reject</button>
      </td>
    </tr>`).join("");
  document.getElementById("requests-table").innerHTML = `<table>
    <tr><th>Title</th><th>Details</th><th>User</th><th>Status</th><th>Feedback</th><th>Actions</th></tr>
    ${rows || "<tr><td colspan='6'>No requests</td></tr>"}
  </table>`;
}

// initialize
setView("users");
