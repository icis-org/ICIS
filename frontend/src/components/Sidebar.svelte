<script lang="ts">
  import { app } from '../lib/state.svelte';
  import type { Screen } from '../lib/state.svelte';
  import favicon from '../assets/favicon.png';

  const navItems: Array<{ screen: Screen; label: string; icon: string }> = [
    { screen: 'home', label: 'Home', icon: 'home' },
    { screen: 'create', label: 'Create .ici', icon: 'create' },
    { screen: 'installed', label: 'Installed', icon: 'installed' },
    { screen: 'store', label: 'Store', icon: 'store' },
  ];

  function navigate(screen: Screen) {
    app.setScreen(screen);
  }
</script>

<aside class="sidebar">
  <div class="sidebar-header">
    <img src={favicon} alt="ICIS" class="sidebar-logo" />
    <div class="sidebar-brand">
      <h1 class="sidebar-title">ICIS</h1>
      <div class="sidebar-subtitle">Configurable Install System</div>
    </div>
  </div>

  <nav class="sidebar-nav">
    {#each navItems as item}
      <button
        class="nav-item"
        class:active={app.currentScreen === item.screen}
        onclick={() => navigate(item.screen)}
      >
        {#if item.icon === 'home'}
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
        {:else if item.icon === 'create'}
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="12" y1="18" x2="12" y2="12"/><line x1="9" y1="15" x2="15" y2="15"/></svg>
        {:else if item.icon === 'installed'}
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
        {:else if item.icon === 'store'}
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>
        {/if}
        <span>{item.label}</span>
      </button>
    {/each}
  </nav>

  <div class="sidebar-footer">
    <span class="version">ICIS v1.0.0</span>
  </div>
</aside>

<style>
  .sidebar {
    width: 220px;
    min-width: 220px;
    background: rgba(22, 24, 34, 0.75);
    backdrop-filter: blur(16px) saturate(180%);
    -webkit-backdrop-filter: blur(16px) saturate(180%);
    border-right: 1px solid var(--border);
    display: flex;
    flex-direction: column;
    overflow-y: auto;
  }

  .sidebar-header {
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    padding: var(--spacing-md) var(--spacing-lg);
    border-bottom: 1px solid var(--border);
  }

  .sidebar-logo {
    width: 28px;
    height: 28px;
    border-radius: 6px;
    flex-shrink: 0;
  }

  .sidebar-brand {
    display: flex;
    flex-direction: column;
  }

  .sidebar-title {
    font-size: var(--font-size-lg);
    font-weight: 700;
    letter-spacing: -0.5px;
    color: var(--text-primary);
    margin: 0;
    line-height: 1.2;
  }

  .sidebar-subtitle {
    font-size: 10px;
    color: var(--text-muted);
    font-weight: 500;
    line-height: 1.2;
  }

  .sidebar-nav {
    flex: 1;
    padding: var(--spacing-sm);
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 10px 14px;
    cursor: pointer;
    transition: all 0.15s ease-in-out;
    color: var(--text-secondary);
    font-size: var(--font-size-sm);
    font-weight: 500;
    border-radius: var(--radius-sm);
    border: none;
    background: none;
    width: 100%;
    text-align: left;
    font-family: inherit;
    position: relative;
  }

  .nav-item:hover {
    background: rgba(255, 255, 255, 0.06);
    color: var(--text-primary);
  }

  .nav-item.active {
    background: var(--bg-active);
    color: var(--accent);
  }

  .nav-item.active::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 2px;
    height: 20px;
    background: var(--accent);
    border-radius: 0 2px 2px 0;
  }

  .nav-item svg {
    width: 18px;
    height: 18px;
    flex-shrink: 0;
  }

  .sidebar-footer {
    padding: var(--spacing-md) var(--spacing-lg);
    border-top: 1px solid var(--border);
  }

  .version {
    font-size: var(--font-size-xs);
    color: var(--text-muted);
  }
</style>
