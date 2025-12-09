import { Auth } from "./modules/auth.js";
import { Chat } from "./modules/chat.js";

document.addEventListener("DOMContentLoaded", () => {
  Chat.init();
  Auth.init();
});
