<script lang="ts">
  import { app } from '../lib/state.svelte';
  import DirSelector from '../components/DirSelector.svelte';

  interface Props {
    oninstall: (customDir: string) => void;
    oncancel: () => void;
  }

  let { oninstall, oncancel }: Props = $props();
  let customDirValue = $state('');

  function handleInstall() {
    const dir = app.selectedDir === 'custom' ? customDirValue : '';
    if (app.selectedDir === 'custom' && !dir) return;
    oninstall(dir);
  }
</script>

<div class="wizard-page active">
  <div class="wizard-card glass-card">
    <h2>{app.currentICI?.name || ''}</h2>
    {#if app.currentICI?.version}
      <div class="version">{app.currentICI.version}</div>
    {/if}
    {#if app.currentICI?.desc}
      <div class="desc">{app.currentICI.desc}</div>
    {/if}

    <div class="detail-section">
      <div class="detail-row">
        <span class="detail-label">Source</span>
        <span class="detail-value detail-url">{app.pendingSource || ''}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Download</span>
        <span class="detail-value detail-url">{app.currentICI?.url || ''}</span>
      </div>
      <div class="detail-row">
        <span class="detail-label">Type</span>
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

    <div class="wizard-actions">
      <button class="btn btn-secondary" onclick={oncancel}>Cancel</button>
      <button class="btn btn-primary" onclick={handleInstall}>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16">
          <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/>
          <polyline points="7 10 12 15 17 10"/>
          <line x1="12" y1="15" x2="12" y2="3"/>
        </svg>
        Install
      </button>
    </div>
  </div>
</div>

<style>
  .wizard-page {
    display: none;
    padding: var(--spacing-xl);
    justify-content: center;
    align-items: flex-start;
  }

  .wizard-page.active {
    display: flex;
  }

  .glass-card {
    width: 100%;
    max-width: 520px;
    background: rgba(30, 32, 46, 0.6);
    backdrop-filter: blur(16px) saturate(180%);
    -webkit-backdrop-filter: blur(16px) saturate(180%);
    border: 1px solid var(--border);
    border-top: 1px solid var(--border-highlight);
    border-radius: var(--radius-md);
    padding: var(--spacing-xl);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
  }

  h2 {
    font-size: var(--font-size-xl);
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 var(--spacing-xs) 0;
  }

  .version {
    color: var(--accent);
    font-size: var(--font-size-sm);
    font-weight: 600;
  }

  .desc {
    color: var(--text-secondary);
    margin: var(--spacing-sm) 0;
    font-size: var(--font-size-base);
    line-height: 1.6;
  }

  .detail-section {
    background: rgba(255, 255, 255, 0.03);
    border-radius: var(--radius-sm);
    border: 1px solid var(--border);
    overflow: hidden;
    margin-top: var(--spacing-md);
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
    max-width: 350px;
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

  .wizard-actions {
    display: flex;
    justify-content: flex-end;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-xl);
  }
</style>
