<script lang="ts">
  import { app, toasts } from '../lib/state.svelte';
  import { InstallApp } from '../../wailsjs/go/main/App.js';
  import { buildICIContent } from '../lib/utils';
  import DirSelector from '../components/DirSelector.svelte';
  import ProgressBar from '../components/ProgressBar.svelte';
  import { onMount } from 'svelte';

  let customDirValue = $state('');
  let showProgress = $state(false);
  let progressPercent = $state(0);
  let progressText = $state('Preparing...');
  let installing = $state(false);
  let installComplete = $state(false);
  let completeName = $state('');
  let completeMessage = $state('');

  onMount(() => {
    const onDownload = (e: Event) => {
      const data = (e as CustomEvent).detail;
      progressPercent = data.percent || 0;
      progressText = `Downloading... ${formatBytes(data.downloaded)} / ${formatBytes(data.total)}`;
    };
    const onInstall = (e: Event) => {
      const data = (e as CustomEvent).detail;
      progressText = data.message || data.status;
    };
    const onExtract = (e: Event) => {
      const data = (e as CustomEvent).detail;
      progressText = `Extracting: ${data.file} (${data.current}/${data.total})`;
    };
    const onComplete = (e: Event) => {
      const data = (e as CustomEvent).detail;
      installComplete = true;
      completeName = data.name;
      completeMessage = data.message;
      progressPercent = 100;
      progressText = data.message;
      installing = false;
    };

    window.addEventListener('icis:download-progress', onDownload);
    window.addEventListener('icis:install-progress', onInstall);
    window.addEventListener('icis:extract-progress', onExtract);
    window.addEventListener('icis:install-complete', onComplete);

    return () => {
      window.removeEventListener('icis:download-progress', onDownload);
      window.removeEventListener('icis:install-progress', onInstall);
      window.removeEventListener('icis:extract-progress', onExtract);
      window.removeEventListener('icis:install-complete', onComplete);
    };
  });

  function formatBytes(bytes: number | undefined | null): string {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }

  async function doInstall() {
    if (!app.currentICI) return;

    let installDir = '';
    if (app.selectedDir === 'custom') {
      installDir = customDirValue;
      if (!installDir) {
        toasts.add('Please select a custom directory', 'error');
        return;
      }
    }

    installing = true;
    showProgress = true;
    progressPercent = 0;
    progressText = 'Preparing installation...';
    installComplete = false;

    const iciContent = buildICIContent(app.currentICI, app.selectedDir, installDir);

    try {
      await InstallApp(iciContent, installDir);
    } catch (err) {
      toasts.add(String(err), 'error');
      installing = false;
      showProgress = false;
    }
  }

  function goBack() {
    showProgress = false;
    installComplete = false;
    app.setScreen('home');
  }
</script>

<div class="screen active">
  <div class="install-card glass-card">
    {#if installComplete}
      <div class="install-done">
        <div class="done-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="48" height="48">
            <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
            <polyline points="22 4 12 14.01 9 11.01"/>
          </svg>
        </div>
        <h2>{completeName}</h2>
        <p class="done-message">{completeMessage}</p>
        <div class="install-actions" style="justify-content:center">
          <button class="btn btn-primary" onclick={goBack}>Done</button>
        </div>
      </div>
    {:else}
      <h2>{app.currentICI?.name || ''}</h2>
      {#if app.currentICI?.version}
        <div class="version">{app.currentICI.version}</div>
      {/if}
      {#if app.currentICI?.desc}
        <div class="desc">{app.currentICI.desc}</div>
      {/if}

      <div class="detail-section">
        <div class="detail-row">
          <span class="detail-label">Download URL</span>
          <span class="detail-value detail-url">{app.currentICI?.url || ''}</span>
        </div>
        <div class="detail-row">
          <span class="detail-label">Archive Type</span>
          <span class="detail-value">{(app.currentICI?.type || 'zip').toUpperCase()}</span>
        </div>
      </div>

      {#if app.currentICI?.shortcuts && app.currentICI.shortcuts.length > 0}
        <div class="shortcuts-section">
          <span class="section-label">Shortcuts</span>
          <div class="shortcuts-list">
            {#each app.currentICI.shortcuts as s}
              <div class="shortcut-item">{s.name || s.exe}</div>
            {/each}
          </div>
        </div>
      {:else if app.currentICI?.shortcut}
        <div class="shortcuts-section">
          <span class="section-label">Shortcut</span>
          <div class="shortcuts-list">
            <div class="shortcut-item">{app.currentICI.shortcut}</div>
          </div>
        </div>
      {/if}

      <div class="dir-section">
        <span class="section-label">Install Location</span>
        <DirSelector />
      </div>

      {#if showProgress}
        <ProgressBar percent={progressPercent} text={progressText} />
      {/if}

      <div class="install-actions">
        <button class="btn btn-secondary" onclick={goBack} disabled={installing}>Back</button>
        {#if !showProgress}
          <button class="btn btn-primary" onclick={doInstall} disabled={installing}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
              <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
              <polyline points="7 10 12 15 17 10"/>
              <line x1="12" y1="15" x2="12" y2="3"/>
            </svg>
            Install
          </button>
        {/if}
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
    align-items: flex-start;
    justify-content: center;
    animation: fadeIn 0.2s ease-in-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(6px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .glass-card {
    background: rgba(30, 32, 46, 0.6);
    backdrop-filter: blur(16px) saturate(180%);
    -webkit-backdrop-filter: blur(16px) saturate(180%);
    border: 1px solid var(--border);
    border-top: 1px solid var(--border-highlight);
    border-radius: var(--radius-md);
    padding: var(--spacing-xl);
    max-width: 600px;
    width: 100%;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  }

  h2 {
    font-size: var(--font-size-xl);
    font-weight: 700;
    margin: 0 0 var(--spacing-xs) 0;
    color: var(--text-primary);
  }

  .version {
    color: var(--accent);
    font-size: var(--font-size-sm);
    font-weight: 600;
  }

  .desc {
    color: var(--text-secondary);
    margin: var(--spacing-md) 0;
    font-size: var(--font-size-base);
    line-height: 1.6;
  }

  .detail-section {
    background: rgba(255, 255, 255, 0.03);
    border-radius: var(--radius-sm);
    border: 1px solid var(--border);
    overflow: hidden;
  }

  .detail-row {
    display: flex;
    justify-content: space-between;
    padding: 10px 14px;
    font-size: var(--font-size-sm);
  }

  .detail-row + .detail-row {
    border-top: 1px solid var(--border);
  }

  .detail-label {
    color: var(--text-muted);
  }

  .detail-value {
    color: var(--text-primary);
    font-weight: 500;
  }

  .detail-url {
    max-width: 300px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .section-label {
    font-size: var(--font-size-sm);
    font-weight: 500;
    color: var(--text-secondary);
  }

  .shortcuts-section {
    margin-top: var(--spacing-md);
  }

  .shortcuts-list {
    margin-top: 6px;
    font-size: var(--font-size-sm);
    color: var(--text-primary);
  }

  .shortcut-item {
    padding: 4px 0;
  }

  .dir-section {
    margin-top: var(--spacing-md);
  }

  .install-actions {
    display: flex;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-lg);
    justify-content: flex-end;
  }

  .install-done {
    text-align: center;
  }

  .done-icon {
    color: var(--success);
    margin-bottom: var(--spacing-md);
  }

  .done-message {
    color: var(--text-secondary);
    font-size: var(--font-size-sm);
  }
</style>
