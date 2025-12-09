import { CONFIG } from "./config.js";

// Cache for DOM elements to avoid repeated lookups
const elements = {};

const getEl = (id) => {
  if (!elements[id]) {
    elements[id] = document.getElementById(id);
  }
  return elements[id];
};

export const UI = {
  // Views
  get authView() {
    return getEl("auth-view");
  },
  get mainView() {
    return getEl("main-view");
  },

  // Auth Forms
  get loginForm() {
    return getEl("login-form");
  },
  get registerForm() {
    return getEl("register-form");
  },
  get toRegisterBtn() {
    return getEl("to-register-btn");
  },
  get toLoginBtn() {
    return getEl("to-login-btn");
  },
  get loginEmail() {
    return getEl("login-email");
  },
  get loginPassword() {
    return getEl("login-password");
  },
  get regName() {
    return getEl("reg-name");
  },
  get regEmail() {
    return getEl("reg-email");
  },
  get regPassword() {
    return getEl("reg-password");
  },

  // Sidebar
  get sidebar() {
    return getEl("sidebar");
  },
  get sidebarOverlay() {
    return getEl("sidebar-overlay");
  },
  get menuBtn() {
    return getEl("menu-btn");
  },
  get newChatBtn() {
    return getEl("new-chat-btn");
  },
  get searchConv() {
    return getEl("search-conv");
  },
  get convList() {
    return getEl("conv-list");
  },
  get userAvatar() {
    return getEl("user-avatar");
  },
  get userName() {
    return getEl("user-name");
  },
  get logoutBtn() {
    return getEl("logout-btn");
  },

  // Chat Area
  get chatContainer() {
    return getEl("chat-container");
  },
  get chatInput() {
    return getEl("chat-input");
  },
  get sendBtn() {
    return getEl("send-btn");
  },
  get emptyState() {
    return getEl("empty-state");
  },
  get thinkingToggle() {
    return getEl("thinking-toggle");
  },
  get statusIndicator() {
    return getEl("status-indicator");
  },

  // Modal
  get modalBackdrop() {
    return getEl("modal-backdrop");
  },
  get modalContent() {
    return getEl("modal-content");
  },
  get modalTitle() {
    return getEl("modal-title");
  },
  get modalInput() {
    return getEl("modal-input");
  },
  get modalCancel() {
    return getEl("modal-cancel");
  },
  get modalConfirm() {
    return getEl("modal-confirm");
  },

  // Toast
  get toastContainer() {
    return getEl("toast-container");
  },

  // Methods
  showToast(message, type = "info") {
    const toast = document.createElement("div");
    const bgClass = type === "error" ? "bg-red-500" : "bg-black";
    toast.className = `${bgClass} text-white px-4 py-3 rounded-lg shadow-lg text-sm flex items-center gap-2 transform transition-all duration-300 translate-y-2 opacity-0`;
    toast.innerHTML = `<span>${message}</span>`;

    this.toastContainer.appendChild(toast);

    // Animation
    requestAnimationFrame(() => {
      toast.classList.remove("translate-y-2", "opacity-0");
    });

    setTimeout(() => {
      toast.classList.add("opacity-0", "translate-x-full");
      setTimeout(() => toast.remove(), CONFIG.ANIMATION_DURATION);
    }, CONFIG.TOAST_DURATION);
  },

  toggleVisibility(element, show) {
    if (!element) return;
    if (show) element.classList.remove("hidden");
    else element.classList.add("hidden");
  },

  setLoading(isLoading) {
    if (isLoading) {
      this.statusIndicator.classList.remove("hidden");
    } else {
      this.statusIndicator.classList.add("hidden");
    }
  },

  escapeHtml(str) {
    if (!str) return "";
    return str
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;")
      .replace(/"/g, "&quot;")
      .replace(/'/g, "&#039;");
  },
};
