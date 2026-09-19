<script lang="ts">
  import { onMount } from 'svelte';
  import { app, toasts } from '../lib/state.svelte';
  import {
    GetRegistryApps,
    GetRegistryURL,
    SetRegistryURL,
    LoadRegistryICI,
    EnterWizardMode,
  } from '../../wailsjs/go/main/App.js';
  import type { registry } from '../../wailsjs/go/models';

  let storeData = $state<registry.RegistryPackage[]>([]);
  let searchQuery = $state('');
  let statusText = $state('');
  let showSettings = $state(false);
  let registryUrl = $state('');

  onMount(async () => {
    await loadStore();
  });

  async function loadStore() {
    statusText = 'Loading...';
    try {
      const result = await GetRegistryApps();
      storeData = result.index?.packages || [];
      statusText = result.error || '';
    } catch {
      statusText = 'Failed to load store';
    }

    try {
      registryUrl = await GetRegistryURL();
    } catch {}
  }

  let filteredPackages = $derived(
    searchQuery
      ? storeData.filter(
          (p) =>
            p.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            (p.desc && p.desc.toLowerCase().includes(searchQuery.toLowerCase()))
        )
      : storeData
  );

  async function saveRegistryUrl() {
    if (!registryUrl.trim()) return;
    try {
      await SetRegistryURL(registryUrl.trim());
      toasts.add('Registry URL saved', 'success');
      showSettings = false;
      loadStore();
    } catch (e) {
      toasts.add('Failed to save URL: ' + e, 'error');
    }
  }

  async function installPackage(iciUrl: string) {
    if (!iciUrl) return;
    try {
      const ici = await LoadRegistryICI(iciUrl);
      if (ici) {
        if (!app.wizardMode) {
          app.wizardMode = true;
          await EnterWizardMode();
        }
        app.setICI(ici, iciUrl);
      }
    } catch (e) {
      toasts.add('Failed to load package: ' + e, 'error');
    }
  }
</script>

<div class="screen active">
  <div class="page-header">
    <h2>Store</h2>
    <div class="header-actions">
      <button class="btn btn-secondary btn-sm" onclick={() => showSettings = !showSettings}>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
          <circle cx="12" cy="12" r="3"/>
          <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.32 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/>
        </svg>
        Settings
      </button>
      <button class="btn btn-secondary btn-sm" onclick={loadStore}>
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
          <polyline points="23 4 23 10 17 10"/>
          <polyline points="1 20 1 14 7 14"/>
          <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
        </svg>
        Refresh
      </button>
    </div>
  </div>

  {#if statusText}
    <div class="store-status">{statusText}</div>
  {/if}

  <input
    type="text"
    class="store-search"
    placeholder="Search packages..."
    bind:value={searchQuery}
  />

  {#if showSettings}
    <div class="settings-panel glass-card">
      <span class="section-label">Registry URL</span>
      <div class="settings-row">
        <input type="text" class="settings-input" bind:value={registryUrl} />
        <button class="btn btn-primary btn-sm" onclick={saveRegistryUrl}>Save</button>
      </div>
    </div>
  {/if}

  <div class="store-list">
    {#if filteredPackages.length === 0}
      <div class="empty-state">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/>
          <line x1="3" y1="6" x2="21" y2="6"/>
          <path d="M16 10a4 4 0 0 1-8 0"/>
        </svg>
        <h3>No packages found</h3>
        <p>Check your registry URL or try a different search.</p>
      </div>
    {:else}
      {#each filteredPackages as pkg (pkg.name)}
        <div class="store-item">
          <div class="app-info">
            <span class="name">{pkg.name}{pkg.version ? ` v${pkg.version}` : ''}</span>
            <span class="meta">{pkg.desc || ''}{pkg.homepage ? ` \u00b7 ${pkg.homepage}` : ''}</span>
          </div>
          <div class="app-actions">
            <button class="btn btn-primary btn-sm" onclick={() => installPackage(pkg.ici)}>
              Install
            </button>
          </div>
        </div>
      {/each}
    {/if}
  </div>
</div>

<style>
  .screen {
    display: none;
    height: 100%;
  }

  .screen.active {
    display: block;
    animation: fadeIn 0.2s ease-in-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(6px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: var(--spacing-lg);
  }

  .page-header h2 {
    font-size: var(--font-size-xl);
    font-weight: 700;
  }

  .header-actions {
    display: flex;
    gap: 8px;
  }

  .store-status {
    margin-bottom: var(--spacing-sm);
    font-size: var(--font-size-xs);
    color: var(--text-secondary);
  }

  .store-search {
    width: 100%;
    padding: 10px 14px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    color: var(--text-primary);
    font-size: var(--font-size-sm);
    outline: none;
    box-sizing: border-box;
    margin-bottom: var(--spacing-md);
    transition: all 0.15s ease-in-out;
    font-family: inherit;
  }

  .store-search::placeholder {
    color: var(--text-muted);
  }

  .store-search:focus {
    border-color: var(--border-focus);
    box-shadow: 0 0 10px var(--accent-glow);
  }

  .settings-panel {
    padding: var(--spacing-md);
    margin-bottom: var(--spacing-md);
    border-radius: var(--radius-md);
  }

  .section-label {
    font-size: var(--font-size-sm);
    font-weight: 500;
    color: var(--text-secondary);
  }

  .settings-row {
    display: flex;
    gap: 8px;
    margin-top: 8px;
  }

  .settings-input {
    flex: 1;
    padding: 10px 14px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    color: var(--text-primary);
    font-size: var(--font-size-sm);
    outline: none;
    font-family: inherit;
    transition: all 0.15s ease-in-out;
  }

  .settings-input:focus {
    border-color: var(--border-focus);
    box-shadow: 0 0 10px var(--accent-glow);
  }

  .store-list {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .store-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: rgba(30, 32, 46, 0.4);
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: 10px;
    padding: 14px 18px;
    transition: all 0.15s ease-in-out;
  }

  .store-item:hover {
    background: rgba(30, 32, 46, 0.6);
    border-color: rgba(255, 255, 255, 0.12);
  }

  .app-info {
    display: flex;
    flex-direction: column;
    gap: 3px;
  }

  .app-info .name {
    font-size: var(--font-size-base);
    font-weight: 600;
    color: var(--text-primary);
  }

  .app-info .meta {
    font-size: var(--font-size-xs);
    color: var(--text-muted);
  }

  .app-actions {
    display: flex;
    gap: var(--spacing-sm);
  }

  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px var(--spacing-lg);
    gap: 12px;
  }

  .empty-state svg {
    width: 48px;
    height: 48px;
    color: var(--text-muted);
  }

  .empty-state h3 {
    font-size: var(--font-size-lg);
    font-weight: 600;
  }

  .empty-state p {
    color: var(--text-secondary);
    font-size: var(--font-size-sm);
    text-align: center;
  }
</style>
