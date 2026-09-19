<script lang="ts">
  import { onMount } from 'svelte';
  import { toasts, installedApps } from '../lib/state.svelte';
  import { ListInstalledApps, UninstallApp, LaunchApp, OpenHomepage } from '../../wailsjs/go/main/App.js';

  onMount(() => {
    refreshInstalled();
  });

  async function refreshInstalled() {
    installedApps.loading = true;
    try {
      const apps = await ListInstalledApps();
      installedApps.apps = apps || [];
    } catch (err) {
      toasts.add('Failed to load installed apps: ' + err, 'error');
    }
    installedApps.loading = false;
  }

  async function uninstall(name: string) {
    if (!confirm(`Are you sure you want to uninstall ${name}?`)) return;
    try {
      await UninstallApp(name);
      toasts.add(`${name} uninstalled`, 'success');
      refreshInstalled();
    } catch (err) {
      toasts.add(String(err), 'error');
    }
  }

  async function launch(name: string) {
    try {
      await LaunchApp(name);
    } catch (err) {
      toasts.add(String(err), 'error');
    }
  }

  async function homepage(name: string) {
    try {
      await OpenHomepage(name);
    } catch (err) {
      toasts.add(String(err), 'error');
    }
  }
</script>

<div class="screen active">
  <div class="page-header">
    <h2>Installed Apps</h2>
    <button class="btn btn-secondary" onclick={refreshInstalled}>
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
        <polyline points="23 4 23 10 17 10"/>
        <polyline points="1 20 1 14 7 14"/>
        <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
      </svg>
      Refresh
    </button>
  </div>

  <div class="installed-list">
    {#if installedApps.loading}
      <div class="empty-state">
        <p>Loading...</p>
      </div>
    {:else if installedApps.apps.length === 0}
      <div class="empty-state">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
          <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
          <polyline points="22 4 12 14.01 9 11.01"/>
        </svg>
        <h3>No apps installed</h3>
        <p>Open an .ici file or create one to install your first app.</p>
      </div>
    {:else}
      {#each installedApps.apps as appItem (appItem.name)}
        <div class="installed-item">
          <div class="app-info">
            <span class="name">{appItem.name}</span>
            <span class="meta">{appItem.version ? `v${appItem.version}` : ''}{appItem.installPath ? ` \u00b7 ${appItem.installPath}` : ''}</span>
          </div>
          <div class="app-actions">
            <button class="btn btn-secondary btn-sm" onclick={() => launch(appItem.name)}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                <polygon points="5 3 19 12 5 21 5 3"/>
              </svg>
              Launch
            </button>
            {#if appItem.homepage}
              <button class="btn btn-secondary btn-sm" onclick={() => homepage(appItem.name)}>
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                  <path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/>
                  <polyline points="15 3 21 3 21 9"/>
                  <line x1="10" y1="14" x2="21" y2="3"/>
                </svg>
                Homepage
              </button>
            {/if}
            <button class="btn btn-danger btn-sm" onclick={() => uninstall(appItem.name)}>
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14">
                <polyline points="3 6 5 6 21 6"/>
                <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/>
              </svg>
              Uninstall
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
    margin-bottom: var(--spacing-xl);
  }

  .page-header h2 {
    font-size: var(--font-size-xl);
    font-weight: 700;
  }

  .installed-list {
    display: flex;
    flex-direction: column;
    gap: var(--spacing-sm);
  }

  .installed-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: rgba(30, 32, 46, 0.4);
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: 10px;
    padding: 14px 18px;
    transition: all 0.15s ease-in-out;
  }

  .installed-item:hover {
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
