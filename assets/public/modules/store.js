export const Store = {
  state: {
    token: localStorage.getItem("token"),
    currentConversationId: null,
    user: null,
    eventSource: null,
    isGenerating: false,
  },

  get token() {
    return this.state.token;
  },
  set token(val) {
    this.state.token = val;
    if (val) localStorage.setItem("token", val);
    else localStorage.removeItem("token");
  },

  get user() {
    return this.state.user;
  },
  set user(val) {
    this.state.user = val;
  },

  get currentConversationId() {
    return this.state.currentConversationId;
  },
  set currentConversationId(val) {
    this.state.currentConversationId = val;
  },

  get isGenerating() {
    return this.state.isGenerating;
  },
  set isGenerating(val) {
    this.state.isGenerating = val;
  },

  get eventSource() {
    return this.state.eventSource;
  },
  set eventSource(val) {
    this.state.eventSource = val;
  },

  // Helper to reset sensitive state on logout
  reset() {
    this.token = null;
    this.user = null;
    this.currentConversationId = null;
    this.isGenerating = false;
    if (this.eventSource) {
      this.eventSource.close();
      this.eventSource = null;
    }
  },
};
