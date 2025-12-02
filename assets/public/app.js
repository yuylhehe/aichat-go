import { initAuth } from './modules/auth.js';
import { initChat } from './modules/chat.js';
import { ui } from './modules/ui.js';

initAuth();
initChat();

ui.toggleSidebar.onclick = () =>
  document.getElementById('sidebar').classList.toggle('hidden');

ui.openSettings.onclick = () => {
  ui.settingsPanel.classList.remove('hidden');
  ui.settingsPanel.classList.add('flex');
};

ui.closeSettings.onclick = () => {
  ui.settingsPanel.classList.add('hidden');
  ui.settingsPanel.classList.remove('flex');
};
