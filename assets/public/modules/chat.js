import { api } from './api.js';
import { ui } from './ui.js';

let currentAssistantContent = null;
let assistantBuffer = '';
let es;

export function initChat() {
  ui.newConv.onclick = handleNewConv;
  ui.searchConv.oninput = () => loadConversations(ui.searchConv.value.trim());
  ui.send.onclick = handleSend;
}

export async function loadConversations(q) {
  const list = await api.get(
    `/conversations${q ? `?q=${encodeURIComponent(q)}` : ''}`,
  );
  ui.convList.innerHTML = '';
  list.items.forEach((c) => {
    const li = document.createElement('li');
    li.dataset.id = c.id;
    li.className =
      'px-3 py-2 rounded hover:bg-gray-100 flex items-center justify-between';
    li.innerHTML = `<span>${c.name}</span>
    <div class="space-x-2">
    <button data-id="${c.id}" class="pin px-2 py-1 text-xs border rounded">置顶</button>
    <button data-id="${c.id}" class="del px-2 py-1 text-xs border rounded">删除</button>
    </div>`;
    li.onclick = (e) => {
      if (
        e.target.classList.contains('pin') ||
        e.target.classList.contains('del')
      )
        return;
      const activeConv = document.querySelector('#convList li.active');
      if (activeConv) activeConv.classList.remove('active', 'bg-gray-200');
      li.classList.add('active', 'bg-gray-200');
      openConversation(c.id);
    };
    li.querySelector('.del').onclick = async (e) => {
      e.stopPropagation();
      if (!confirm('确定删除吗？')) return;
      await api.del(`/conversations/${c.id}`);
      loadConversations(q);
    };
    li.querySelector('.pin').onclick = (e) => {
      e.stopPropagation();
      alert('置顶功能后端字段新增后可用');
    };
    ui.convList.appendChild(li);
  });
}

async function handleNewConv() {
  const name = prompt('会话名称');
  if (!name) return;
  try {
    const user = await api.get('/users/profile');
    await api.post('/conversations', { name, userId: user.id });
    loadConversations();
  } catch (e) {
    console.error(e);
  }
}

async function openConversation(id) {
  console.log('openConversation id =', id);
  ui.chat.innerHTML = '';
  // 关闭已有的流
  if (es) {
      es.close();
      es = null;
  }
  const res = await api.get(`/messages/conversation/${id}`);
  res.items.forEach(renderMessage);
}

function renderMessage(m) {
  const div = document.createElement('div');
  div.className = `message ${m.type === 'user' ? 'user' : 'assistant'}`;
  const meta = document.createElement('div');
  meta.className = 'meta';
  meta.textContent = `${m.type} · ${new Date(m.createdAt || m.createAt || Date.now()).toLocaleString()}`;
  
  const contentWrapper = document.createElement('div');
  
  // 渲染思考过程详情（如果存在）
  if (m.reasoningContent) {
      const details = document.createElement('details');
      details.className = 'mb-2 p-2 bg-gray-50 rounded border text-sm text-gray-600';
      // 默认折叠状态，用于历史记录
      const summary = document.createElement('summary');
      summary.className = 'cursor-pointer font-medium text-gray-500 hover:text-gray-700 select-none';
      summary.textContent = '思考过程';
      const contentDiv = document.createElement('div');
      contentDiv.className = 'mt-2 prose prose-sm max-w-none';
      contentDiv.innerHTML = marked.parse(m.reasoningContent);
      
      details.appendChild(summary);
      details.appendChild(contentDiv);
      contentWrapper.appendChild(details);
  }

  const content = document.createElement('div');
  content.className = 'content prose prose-sm max-w-none'; // Added class for styling consistency
  content.innerHTML = marked.parse(m.content || '');
  
  contentWrapper.appendChild(content);
  
  div.appendChild(meta);
  div.appendChild(contentWrapper);
  ui.chat.appendChild(div);
  ui.chat.scrollTop = ui.chat.scrollHeight;
}

// 处理发送消息
async function handleSend() {
  const content = ui.input.value.trim();
  if (!content) return;

  setLoadingState(true);

  ui.sendState.textContent = '发送中...';
  const activeConv = document.querySelector('#convList li.active');
  const convId = activeConv?.dataset?.id ? Number(activeConv.dataset.id) : null;
  if (!convId) {
    ui.sendState.textContent = '请选择会话';
    setLoadingState(false);
    return;
  }
  
  try {
    await api.post('/messages', {
        conversationId: convId,
        content,
        type: 'user',
    });
    
    renderMessage({ type: 'user', content, createAt: new Date().toISOString() });
    ui.input.value = '';
    ui.sendState.textContent = '';
    
    createAssistantContainer();
    assistantBuffer = '';
    startStream({ conversationId: convId, message: content });
  } catch (e) {
    ui.sendState.textContent = '发送失败: ' + e.message;
    setLoadingState(false);
  }
}

/**
 * 创建助手消息容器
 */
function createAssistantContainer() {
  const div = document.createElement('div');
  div.className = 'message assistant';
  const meta = document.createElement('div');
  meta.className = 'meta';
  meta.textContent = `answer · ${new Date().toLocaleString()}`;
  const content = document.createElement('div');
  content.className = 'content prose prose-sm max-w-none';
  div.appendChild(meta);
  div.appendChild(content);
  ui.chat.appendChild(div);
  currentAssistantContent = content;
  ui.chat.scrollTop = ui.chat.scrollHeight;
}

// 启动事件流
function startStream({ conversationId, message }) {
  if (es) es.close();
  ui.status.textContent = '生成中…';

  const params = new URLSearchParams();
  if (message) params.append('prompt', message);
  
  const token = api.getToken();
  if (token) params.append('token', token);

  const url = `${api.base}/ai/stream/${conversationId}?${params.toString()}`;
  
  // 开启思考过程详情（如果勾选）
  if (ui.thinkingToggle && ui.thinkingToggle.checked) {
      params.append('thinking', 'enabled');
  } else {
      params.append('thinking', 'disabled');
  }
  
  // 构建最终的URL
  const finalUrl = `${api.base}/ai/stream/${conversationId}?${params.toString()}`;
  es = new EventSource(finalUrl);

  let reasoningBuffer = '';
  let currentReasoningContent = null;

  // 处理事件流数据
  es.onmessage = (e) => {
    try {
      const data = JSON.parse(e.data);
      if (data.type === 'heartbeat') return;

      if (data.type === 'token' || data.type === 'message') {
        if (!currentAssistantContent) createAssistantContainer();
        
        // 当正式回答开始，折叠推理过程
        if (currentReasoningContent && currentReasoningContent.parentElement.open) {
            currentReasoningContent.parentElement.open = false;
        }
        
        assistantBuffer += data.content || '';
        currentAssistantContent.innerHTML = marked.parse(assistantBuffer);
        ui.chat.scrollTop = ui.chat.scrollHeight;
      }

      if (data.type === 'reasoning') {
        if (!currentAssistantContent) createAssistantContainer();
        
        // 创建推理容器
        if (!currentReasoningContent) {
            const details = document.createElement('details');
            details.className = 'mb-2 p-2 bg-gray-50 rounded border text-sm text-gray-600';
            details.open = true; // Default open to show thinking
            const summary = document.createElement('summary');
            summary.className = 'cursor-pointer font-medium text-gray-500 hover:text-gray-700 select-none';
            summary.textContent = '思考过程';
            const contentDiv = document.createElement('div');
            contentDiv.className = 'mt-2 prose prose-sm max-w-none';
            
            details.appendChild(summary);
            details.appendChild(contentDiv);
            
            // 插入思考过程详情到消息容器中
            currentAssistantContent.parentElement.insertBefore(details, currentAssistantContent);
            // 初始化思考过程内容容器
            currentReasoningContent = contentDiv;
        }
        
        reasoningBuffer += data.content || '';
        currentReasoningContent.innerHTML = marked.parse(reasoningBuffer);
        ui.chat.scrollTop = ui.chat.scrollHeight;
      }

      if (data.type === 'error') {
        ui.status.textContent = `生成失败：${data.message || '连接异常'}`;
        currentAssistantContent = null;
        assistantBuffer = '';
        es.close();
        setLoadingState(false);
        return;
      }

      if (data.type === 'finish') {
        ui.status.textContent = '对话完成';
        if (currentAssistantContent) {
          currentAssistantContent.innerHTML = marked.parse(assistantBuffer || '');
        }
        currentAssistantContent = null;
        assistantBuffer = '';
        currentReasoningContent = null;
        reasoningBuffer = '';
        es.close();
        setLoadingState(false);
      }
    } catch (error) {
      console.error('SSE Parse Error', error);
    }
  };

  es.onerror = () => {
    if (ui.status.textContent.startsWith('生成失败')) return;
    ui.status.textContent = '连接异常';
    currentAssistantContent = null;
    assistantBuffer = '';
    es.close();
    setLoadingState(false);
  };
}

// 设置加载状态（禁用输入和发送按钮）
function setLoadingState(loading) {
  ui.input.disabled = loading;
  ui.send.disabled = loading;
  if (loading) {
    ui.send.classList.add('opacity-50', 'cursor-not-allowed');
  } else {
    ui.send.classList.remove('opacity-50', 'cursor-not-allowed');
    ui.input.focus();
  }
}
