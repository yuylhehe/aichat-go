import { API } from "./api.js";
import { Store } from "./store.js";
import { UI } from "./ui.js";
import { API_BASE } from "./config.js";

export const Chat = {
  renamingId: null,

  init() {
    this.bindEvents();

    // Listen to Auth events
    window.addEventListener("auth:login", () => this.loadConversations());
    window.addEventListener("auth:logout", () => this.clearChat());
  },

  bindEvents() {
    UI.newChatBtn.onclick = () => this.startNewChat();
    UI.sendBtn.onclick = () => this.sendMessage();

    // Close dropdowns when clicking outside
    document.addEventListener("click", (e) => {
      if (!e.target.closest(".menu-trigger")) {
        document
          .querySelectorAll(".menu-dropdown")
          .forEach((el) => el.classList.add("hidden"));
      }
    });

    // Auto-resize textarea
    UI.chatInput.addEventListener("input", function () {
      this.style.height = "auto";
      this.style.height = this.scrollHeight + "px";
      if (this.value === "") this.style.height = "auto";

      // Enable/disable send button
      UI.sendBtn.disabled = !this.value.trim() || Store.isGenerating;
      if (this.value.trim() && !Store.isGenerating) {
        UI.sendBtn.classList.remove("opacity-30");
      } else {
        UI.sendBtn.classList.add("opacity-30");
      }
    });

    UI.chatInput.addEventListener("keydown", (e) => {
      if (e.key === "Enter" && !e.shiftKey) {
        e.preventDefault();
        this.sendMessage();
      }
    });

    UI.searchConv.oninput = (e) =>
      this.loadConversations(e.target.value.trim());

    // Sidebar Mobile Toggle
    UI.menuBtn.onclick = () => {
      UI.sidebar.classList.remove("-translate-x-full");
      UI.toggleVisibility(UI.sidebarOverlay, true);
    };

    UI.sidebarOverlay.onclick = () => {
      UI.sidebar.classList.add("-translate-x-full");
      UI.toggleVisibility(UI.sidebarOverlay, false);
    };

    // Modal Logic
    const closeModal = () => {
      UI.modalBackdrop.classList.remove("opacity-100");
      UI.modalBackdrop.classList.add("opacity-0");
      setTimeout(() => UI.toggleVisibility(UI.modalBackdrop, false), 300);
      this.renamingId = null;
    };

    UI.modalCancel.onclick = closeModal;

    UI.modalConfirm.onclick = async () => {
      if (this.renamingId) {
        const newName = UI.modalInput.value.trim();
        if (newName) {
          await this.renameConversation(this.renamingId, newName);
        }
      }
      closeModal();
    };
  },

  async renameConversation(id, name) {
    try {
      await API.put(`/conversations/${id}`, { name });
      this.loadConversations();
    } catch (e) {
      UI.showToast(e.message, "error");
    }
  },

  async loadConversations(q = "") {
    try {
      const res = await API.get(
        `/conversations${q ? `?q=${encodeURIComponent(q)}` : ""}`
      );
      const list = res.items || [];

      UI.convList.innerHTML = "";
      list.forEach((conv) => {
        const item = this.createConversationItem(conv);
        UI.convList.appendChild(item);
      });
    } catch (e) {
      console.error(e);
    }
  },

  createConversationItem(conv) {
    const item = document.createElement("div");
    const isActive = conv.id == Store.currentConversationId;
    item.className = `group flex items-center gap-3 p-2 rounded-lg cursor-pointer transition-colors relative ${
      isActive ? "bg-gray-200" : "hover:bg-gray-100"
    }`;

    item.innerHTML = `
            <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-900 truncate">${UI.escapeHtml(
                  conv.name
                )}</p>
            </div>
            <div class="relative">
                <button class="menu-trigger opacity-0 group-hover:opacity-100 p-1 text-gray-500 hover:text-black transition-all rounded-md hover:bg-gray-200/50" title="Options">
                    <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24"><path d="M12 8c1.1 0 2-.9 2-2s-.9-2-2-2-2 .9-2 2 .9 2 2 2zm0 2c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"></path></svg>
                </button>
                <div class="menu-dropdown hidden absolute right-0 top-8 w-32 bg-white rounded-lg shadow-xl border border-gray-100 z-50 overflow-hidden py-1">
                    <button class="rename-btn w-full text-left px-4 py-2 text-xs text-gray-700 hover:bg-gray-50 flex items-center gap-2">
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"></path></svg>
                        Rename
                    </button>
                    <button class="delete-btn w-full text-left px-4 py-2 text-xs text-red-600 hover:bg-gray-50 flex items-center gap-2">
                        <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path></svg>
                        Delete
                    </button>
                </div>
            </div>
        `;

    item.onclick = (e) => {
      // Prevent opening if clicking menu elements
      if (
        e.target.closest(".menu-trigger") ||
        e.target.closest(".menu-dropdown")
      )
        return;

      this.openConversation(conv.id);
      // Mobile: close sidebar
      if (window.innerWidth < 768) {
        UI.sidebar.classList.add("-translate-x-full");
        UI.toggleVisibility(UI.sidebarOverlay, false);
      }
    };

    const trigger = item.querySelector(".menu-trigger");
    const dropdown = item.querySelector(".menu-dropdown");
    const renameBtn = item.querySelector(".rename-btn");
    const deleteBtn = item.querySelector(".delete-btn");

    trigger.onclick = (e) => {
      e.stopPropagation();
      // Close others
      document
        .querySelectorAll(".menu-dropdown")
        .forEach((el) => el !== dropdown && el.classList.add("hidden"));
      dropdown.classList.toggle("hidden");
    };

    renameBtn.onclick = (e) => {
      e.stopPropagation();
      dropdown.classList.add("hidden");
      trigger.classList.remove("opacity-100");
      
      this.renamingId = conv.id;
      UI.modalTitle.innerText = "Rename Chat";
      UI.modalInput.value = conv.name;
      
      UI.toggleVisibility(UI.modalBackdrop, true);
      // Small delay for transition
      requestAnimationFrame(() => {
        UI.modalBackdrop.classList.remove("opacity-0");
        UI.modalBackdrop.classList.add("opacity-100");
      });
      
      setTimeout(() => UI.modalInput.focus(), 100);
    };

    deleteBtn.onclick = async (e) => {
      e.stopPropagation();
      dropdown.classList.add("hidden");
      if (confirm("Delete this conversation?")) {
        await API.del(`/conversations/${conv.id}`);
        if (Store.currentConversationId == conv.id) {
          this.startNewChat();
        } else {
          this.loadConversations();
        }
      }
    };

    return item;
  },

  startNewChat() {
    Store.currentConversationId = null;
    this.clearChat();
    this.loadConversations(); // Refresh highlight
    UI.chatInput.focus();
  },

  clearChat() {
    UI.chatContainer.innerHTML = "";
    UI.toggleVisibility(UI.emptyState, true);
    Store.currentConversationId = null;
    // Reset active state in list
    Array.from(UI.convList.children).forEach((child) =>
      child.classList.remove("bg-gray-200")
    );
  },

  async openConversation(id) {
    if (Store.currentConversationId === id) return;
    Store.currentConversationId = id;

    // Highlight in list
    this.loadConversations();

    UI.toggleVisibility(UI.emptyState, false);
    UI.chatContainer.innerHTML =
      '<div class="flex justify-center p-4"><span class="loading-dot"></span></div>';

    try {
      const res = await API.get(`/messages/conversation/${id}`);
      UI.chatContainer.innerHTML = "";

      const messages = res.items || [];
      if (messages.length === 0) {
        UI.toggleVisibility(UI.emptyState, true);
      } else {
        messages.forEach((msg) => this.appendMessage(msg));
        this.scrollToBottom();
      }
    } catch (e) {
      UI.chatContainer.innerHTML = `<div class="text-red-500 text-center p-4">Failed to load messages: ${e.message}</div>`;
    }
  },

  async sendMessage() {
    const content = UI.chatInput.value.trim();
    if (!content || Store.isGenerating) return;

    UI.chatInput.value = "";
    UI.chatInput.style.height = "auto";
    UI.sendBtn.disabled = true;
    UI.sendBtn.classList.add("opacity-30");
    UI.toggleVisibility(UI.emptyState, false);

    // Optimistic UI for User Message
    this.appendMessage({ type: "user", content, createdAt: new Date() });
    this.scrollToBottom();

    Store.isGenerating = true;
    UI.setLoading(true);

    try {
      // 1. Ensure Conversation
      if (!Store.currentConversationId) {
        const name = content.slice(0, 30) || "New Chat";
        const conv = await API.post("/conversations", { name });
        Store.currentConversationId = conv.id;
        this.loadConversations();
      }

      // 2. Send Message
      await API.post("/messages", {
        conversationId: Store.currentConversationId,
        content,
        type: "user",
      });

      // 3. Start SSE
      this.startSSE(Store.currentConversationId, content);
    } catch (e) {
      UI.showToast(e.message, "error");
      Store.isGenerating = false;
      UI.sendBtn.disabled = false;
      UI.sendBtn.classList.remove("opacity-30");
      UI.setLoading(false);
    }
  },

  appendMessage(msg) {
    const isUser = msg.type === "user";
    const div = document.createElement("div");
    div.className = `flex w-full ${
      isUser ? "justify-end" : "justify-start"
    } animate-fade-in`;

    const bubble = document.createElement("div");
    bubble.className = isUser ? "message-user" : "message-assistant";

    let htmlContent = "";

    // Reasoning
    if (msg.reasoningContent) {
      htmlContent += this.buildReasoningHTML(msg.reasoningContent);
    }

    htmlContent += `<div class="prose prose-sm max-w-none ${
      isUser ? "prose-invert" : ""
    }">${marked.parse(msg.content || "")}</div>`;

    bubble.innerHTML = htmlContent;
    div.appendChild(bubble);
    UI.chatContainer.appendChild(div);
  },

  buildReasoningHTML(content) {
    return `
            <details class="mb-3 group">
                <summary class="cursor-pointer list-none text-xs font-medium text-gray-500 flex items-center gap-1 select-none hover:text-gray-800 transition-colors">
                    <svg class="w-3 h-3 transition-transform group-open:rotate-90" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
                    Thinking Process
                </summary>
                <div class="mt-2 pl-4 border-l-2 border-gray-100 text-gray-500 text-sm prose prose-sm max-w-none">
                    ${marked.parse(content)}
                </div>
            </details>
        `;
  },

  startSSE(conversationId, prompt) {
    if (Store.eventSource) Store.eventSource.close();

    const params = new URLSearchParams();
    if (prompt) params.append("prompt", prompt);
    if (Store.token) params.append("token", Store.token);
    params.append(
      "thinking",
      UI.thinkingToggle.checked ? "enabled" : "disabled"
    );

    const url = `${API_BASE}/ai/stream/${conversationId}?${params.toString()}`;
    Store.eventSource = new EventSource(url);

    // Create a placeholder for assistant response
    const div = document.createElement("div");
    div.className = "flex w-full justify-start animate-fade-in";
    div.innerHTML = `
            <div class="message-assistant">
                <div class="assistant-content prose prose-sm max-w-none"></div>
            </div>
        `;
    UI.chatContainer.appendChild(div);
    const contentContainer = div.querySelector(".assistant-content");

    let fullContent = "";
    let fullReasoning = "";
    let reasoningContainer = null;

    Store.eventSource.onmessage = (e) => {
      try {
        const data = JSON.parse(e.data);
        if (data.type === "heartbeat") return;

        if (data.type === "reasoning") {
          if (!reasoningContainer) {
            const details = document.createElement("details");
            details.open = true;
            details.className = "mb-3 group";
            details.innerHTML = `
                            <summary class="cursor-pointer list-none text-xs font-medium text-gray-500 flex items-center gap-1 select-none hover:text-gray-800 transition-colors">
                                <svg class="w-3 h-3 transition-transform group-open:rotate-90" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path></svg>
                                Thinking Process
                            </summary>
                            <div class="reasoning-body mt-2 pl-4 border-l-2 border-gray-100 text-gray-500 text-sm prose prose-sm max-w-none"></div>
                        `;
            contentContainer.parentElement.insertBefore(
              details,
              contentContainer
            );
            reasoningContainer = details.querySelector(".reasoning-body");
          }
          fullReasoning += data.content;
          reasoningContainer.innerHTML = marked.parse(fullReasoning);
          this.scrollToBottom();
        }

        if (data.type === "token" || data.type === "message") {
          // If we switch from reasoning to content, close reasoning details
          if (reasoningContainer && reasoningContainer.parentElement.open) {
            reasoningContainer.parentElement.open = false;
          }

          fullContent += data.content;
          contentContainer.innerHTML = marked.parse(fullContent);
          this.scrollToBottom();
        }

        if (data.type === "finish" || data.type === "error") {
          this.endSSE();
          if (data.type === "error") {
            UI.showToast(data.message || "Generation failed", "error");
          }
        }
      } catch (err) {
        console.error(err);
      }
    };

    Store.eventSource.onerror = () => {
      this.endSSE();
    };
  },

  endSSE() {
    if (Store.eventSource) {
      Store.eventSource.close();
      Store.eventSource = null;
    }
    Store.isGenerating = false;
    UI.setLoading(false);
    UI.sendBtn.disabled = false;
    UI.sendBtn.classList.remove("opacity-30");
  },

  scrollToBottom() {
    UI.chatContainer.scrollTop = UI.chatContainer.scrollHeight;
  },
};
