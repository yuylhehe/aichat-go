import { api } from './api.js';
import { ui } from './ui.js';
import { loadConversations } from './chat.js';

export function initAuth() {
  ui.logout.addEventListener('click', handleLogout);
  ui.loginBtn.addEventListener('click', handleLogin);
  ui.registerBtn.addEventListener('click', handleRegister);

  const t = api.getToken();
  if (t) {
    ui.authPanel.classList.add('hidden');
    loadConversations();

    fetchUserProfile();
  }
}

function handleLogout() {
  api.setToken(null);
  location.reload();
}

function showRegisterPanel() {
  ui.authPanel.classList.add('hidden');
  ui.registerPanel.classList.remove('hidden');
}

async function handleLogin() {
  const email = ui.loginEmail.value.trim();
  const password = ui.loginPassword.value.trim();
  try {
    const res = await api.post('/auth/login', { email, password });
    api.setToken(res.accessToken);
    ui.authPanel.classList.add('hidden');
    ui.userName.textContent = res.user.name;
    loadConversations();
  } catch (e) {
    ui.authError.textContent = e.message;
  }
}

async function handleRegister() {
  const username = ui.regName.value.trim();
  const email = ui.regEmail.value.trim();
  const password = ui.regPassword.value.trim();
  
  if (password.length < 6) {
    ui.registerError.textContent = '密码至少6位';
    return;
  }
  
  try {
    const res = await api.post('/auth/register', { username, email, password });
    api.setToken(res.accessToken);
    ui.registerPanel.classList.add('hidden');
    ui.userName.textContent = res.user.username;
    loadConversations();
  } catch (e) {
    ui.registerError.textContent = e.message;
  }
}

async function fetchUserProfile() {
  try {
    const user = await api.get('/users/profile');
    if (user && user.name) {
        ui.userName.textContent = user.name;
    }
  } catch (e) {
    console.error('Failed to fetch profile', e);
  }
}
