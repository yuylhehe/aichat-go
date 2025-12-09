import { API } from "./api.js";
import { Store } from "./store.js";
import { UI } from "./ui.js";

export const Auth = {
  init() {
    // Listen for 401 events
    window.addEventListener("auth:unauthorized", () =>
      this.logout("Session expired")
    );

    if (Store.token) {
      this.showMain();
      this.fetchProfile();
    } else {
      this.showAuth();
    }

    this.bindEvents();
  },

  bindEvents() {
    UI.toRegisterBtn.onclick = () => {
      UI.toggleVisibility(UI.loginForm, false);
      UI.toggleVisibility(UI.registerForm, true);
    };

    UI.toLoginBtn.onclick = () => {
      UI.toggleVisibility(UI.registerForm, false);
      UI.toggleVisibility(UI.loginForm, true);
    };

    UI.loginForm.onsubmit = async (e) => {
      e.preventDefault();
      const email = UI.loginEmail.value;
      const password = UI.loginPassword.value;

      try {
        const res = await API.post("/auth/login", { email, password });
        this.login(res.accessToken, res.user);
      } catch (err) {
        UI.showToast(err.message, "error");
      }
    };

    UI.registerForm.onsubmit = async (e) => {
      e.preventDefault();
      const username = UI.regName.value;
      const email = UI.regEmail.value;
      const password = UI.regPassword.value;

      try {
        const res = await API.post("/auth/register", {
          username,
          email,
          password,
        });
        this.login(res.accessToken, res.user);
      } catch (err) {
        UI.showToast(err.message, "error");
      }
    };

    UI.logoutBtn.onclick = () => this.logout();
  },

  login(token, user) {
    Store.token = token;
    Store.user = user;
    this.updateProfileUI(user);
    this.showMain();
  },

  logout(msg) {
    Store.reset();
    if (msg) UI.showToast(msg);
    this.showAuth();
    // Dispatch logout event so Chat module can clear state
    window.dispatchEvent(new CustomEvent("auth:logout"));
  },

  async fetchProfile() {
    try {
      const user = await API.get("/users/profile");
      Store.user = user;
      this.updateProfileUI(user);
    } catch (e) {
      // Token might be invalid, handled by interceptor
    }
  },

  updateProfileUI(user) {
    if (user) {
      UI.userName.textContent = user.name || user.username || "User";
      UI.userAvatar.textContent = (user.name ||
        user.username ||
        "U")[0].toUpperCase();
    }
  },

  showAuth() {
    UI.toggleVisibility(UI.authView, true);
    UI.toggleVisibility(UI.mainView, false);
    UI.mainView.classList.remove("opacity-100");
  },

  showMain() {
    UI.toggleVisibility(UI.authView, false);
    UI.toggleVisibility(UI.mainView, true);
    // Small delay for fade in
    setTimeout(() => UI.mainView.classList.add("opacity-100"), 50);
    // Dispatch login event
    window.dispatchEvent(new CustomEvent("auth:login"));
  },
};
