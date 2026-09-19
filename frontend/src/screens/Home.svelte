<script lang="ts">
  import { app, toasts } from '../lib/state.svelte';
  import { OpenFile, LoadICIFile, LoadICIContent } from '../../wailsjs/go/main/App.js';

  let showPasteModal = $state(false);
  let pasteContent = $state('');

  async function openFile() {
    try {
      const filePath = await OpenFile();
      if (filePath) {
        const ici = await LoadICIFile(filePath);
        if (ici) {
          app.setICI(ici, filePath);
          app.setScreen('install');
        }
      }
    } catch (err) {
      toasts.add(String(err), 'error');
    }
  }

  async function submitPaste() {
    const content = pasteContent.trim();
    if (!content) return;
    try {
      const ici = await LoadICIContent(content);
      if (ici) {
        showPasteModal = false;
        pasteContent = '';
        app.setICI(ici, '');
        app.setScreen('install');
      }
    } catch (err) {
      toasts.add(String(err), 'error');
    }
  }

  function cancelPaste() {
    showPasteModal = false;
    pasteContent = '';
  }
</script>

<div class="screen active">
  <div class="welcome glass-card">
    <div class="welcome-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
        <polyline points="7 10 12 15 17 10"/>
        <line x1="12" y1="15" x2="12" y2="3"/>
      </svg>
    </div>
    <h2>Welcome to ICIS</h2>
    <p class="welcome-desc">Open an <strong>.ici</strong> file to install an app, or create your own configuration file.</p>

    <div class="welcome-actions">
      <button class="btn btn-primary" onclick={openFile}>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
        </svg>
        Open .ici File
      </button>
      <button class="btn btn-secondary" onclick={() => showPasteModal = !showPasteModal}>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
          <path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/>
          <rect x="8" y="2" width="8" height="4" rx="1" ry="1"/>
        </svg>
        Paste .ici Content
      </button>
    </div>

    {#if showPasteModal}
      <div class="paste-modal">
        <textarea
          class="paste-input"
          placeholder="Paste .ici content here..."
          bind:value={pasteContent}
        ></textarea>
        <div class="paste-actions">
          <button class="btn btn-secondary btn-sm" onclick={cancelPaste}>Cancel</button>
          <button class="btn btn-primary btn-sm" onclick={submitPaste}>Load</button>
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .screen {
    display: none;
    height: 100%;
  }

  .screen.active {
    display: flex;
    animation: fadeIn 0.2s ease-in-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(6px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .glass-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    width: 100%;
    gap: var(--spacing-xl);
    background: rgba(30, 32, 46, 0.6);
    backdrop-filter: blur(16px) saturate(180%);
    -webkit-backdrop-filter: blur(16px) saturate(180%);
    border: 1px solid var(--border);
    border-top: 1px solid var(--border-highlight);
    border-radius: var(--radius-md);
    padding: var(--spacing-2xl);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  }

  .welcome-icon {
    width: 72px;
    height: 72px;
    background: var(--accent-glow);
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 24px var(--accent-glow);
  }

  .welcome-icon svg {
    width: 36px;
    height: 36px;
    color: var(--accent);
  }

  h2 {
    font-size: var(--font-size-2xl);
    font-weight: 700;
    letter-spacing: -0.5px;
    color: var(--text-primary);
    margin: 0;
  }

  .welcome-desc {
    color: var(--text-secondary);
    font-size: var(--font-size-base);
    max-width: 400px;
    text-align: center;
    line-height: 1.6;
    margin: 0;
  }

  .welcome-actions {
    display: flex;
    gap: var(--spacing-sm);
  }

  .paste-modal {
    width: 400px;
    text-align: left;
  }

  .paste-input {
    width: 100%;
    height: 150px;
    padding: 12px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    color: var(--text-primary);
    font-size: var(--font-size-xs);
    font-family: 'SFMono-Regular', Consolas, monospace;
    resize: vertical;
    outline: none;
    box-sizing: border-box;
  }

  .paste-input:focus {
    border-color: var(--border-focus);
    box-shadow: 0 0 10px var(--accent-glow);
  }

  .paste-actions {
    display: flex;
    gap: 8px;
    margin-top: 8px;
    justify-content: flex-end;
  }
</style>
