<script lang="ts">
  import { onMount } from 'svelte';
  import {
    WindowMinimise,
    WindowToggleMaximise,
    WindowIsMaximised,
    Quit,
  } from '../../wailsjs/runtime/runtime.js';

  let maximized = $state(false);

  onMount(async () => {
    maximized = await WindowIsMaximised();
  });

  async function handleMaximize() {
    await WindowToggleMaximise();
    maximized = await WindowIsMaximised();
  }
</script>

<header class="titlebar" class:maximized>
  <div class="titlebar-drag">
    <span class="titlebar-title">ICIS</span>
  </div>
  <div class="titlebar-controls">
    <button class="tb-btn" onclick={() => WindowMinimise()} aria-label="Minimize">
      <svg width="12" height="12" viewBox="0 0 12 12"><rect x="2" y="5.5" width="8" height="1" fill="currentColor" rx="0.5"/></svg>
    </button>
    <button class="tb-btn" onclick={handleMaximize} aria-label={maximized ? 'Restore' : 'Maximize'}>
      {#if maximized}
        <svg width="12" height="12" viewBox="0 0 12 12">
          <rect x="3.5" y="3.5" width="5.5" height="5.5" fill="none" stroke="currentColor" stroke-width="1" rx="0.5"/>
          <rect x="1" y="5.5" width="5.5" height="5.5" fill="var(--bg-panel)" stroke="currentColor" stroke-width="1" rx="0.5"/>
        </svg>
      {:else}
        <svg width="12" height="12" viewBox="0 0 12 12"><rect x="2" y="2" width="8" height="8" fill="none" stroke="currentColor" stroke-width="1" rx="0.5"/></svg>
      {/if}
    </button>
    <button class="tb-btn tb-close" onclick={() => Quit()} aria-label="Close">
      <svg width="12" height="12" viewBox="0 0 12 12"><line x1="2.5" y1="2.5" x2="9.5" y2="9.5" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"/><line x1="9.5" y1="2.5" x2="2.5" y2="9.5" stroke="currentColor" stroke-width="1.2" stroke-linecap="round"/></svg>
    </button>
  </div>
</header>

<style>
  .titlebar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 36px;
    background: rgba(22, 24, 34, 0.85);
    backdrop-filter: blur(16px) saturate(180%);
    -webkit-backdrop-filter: blur(16px) saturate(180%);
    border-radius: 10px 10px 0 0;
    flex-shrink: 0;
    user-select: none;
  }

  .titlebar.maximized {
    border-radius: 0;
  }

  .titlebar-drag {
    flex: 1;
    display: flex;
    align-items: center;
    padding-left: 14px;
    -webkit-app-region: drag;
    height: 100%;
  }

  .titlebar-title {
    font-size: 12px;
    font-weight: 600;
    color: var(--text-muted);
    letter-spacing: 0.3px;
  }

  .titlebar-controls {
    display: flex;
    height: 100%;
    -webkit-app-region: no-drag;
  }

  .tb-btn {
    width: 46px;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    color: var(--text-secondary);
    cursor: pointer;
    transition: all 0.15s ease-in-out;
  }

  .tb-btn:hover {
    background: rgba(255, 255, 255, 0.08);
    color: var(--text-primary);
  }

  .tb-close:hover {
    background: var(--danger);
    color: white;
  }
</style>
