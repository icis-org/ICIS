<script lang="ts">
  import { SelectDirectory } from '../../wailsjs/go/main/App.js';
  import { app } from '../lib/state.svelte';

  interface Props {
    customDirInputId?: string;
  }

  let { customDirInputId = 'custom-dir-input' }: Props = $props();

  let customDirValue = $state('');

  async function browse() {
    try {
      const dir = await SelectDirectory();
      if (dir) {
        customDirValue = dir;
      }
    } catch {}
  }
</script>

<div class="dir-selector">
  <div class="segment-group">
    <button
      class="segment-item"
      class:selected={app.selectedDir === 'appdata'}
      onclick={() => app.setDir('appdata')}
    >AppData</button>
    <button
      class="segment-item"
      class:selected={app.selectedDir === 'programfiles'}
      onclick={() => app.setDir('programfiles')}
    >Program Files</button>
    <button
      class="segment-item"
      class:selected={app.selectedDir === 'custom'}
      onclick={() => app.setDir('custom')}
    >Custom...</button>
  </div>

  {#if app.selectedDir === 'custom'}
    <div class="custom-dir-row">
      <input
        class="custom-dir-input"
        id={customDirInputId}
        placeholder="Select custom directory..."
        bind:value={customDirValue}
      />
      <button class="btn btn-secondary btn-sm" onclick={browse}>Browse</button>
    </div>
  {/if}
</div>

<style>
  .dir-selector {
    margin-top: var(--spacing-sm);
  }

  .custom-dir-row {
    display: flex;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-sm);
  }

  .custom-dir-input {
    flex: 1;
    padding: 10px 14px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    color: var(--text-primary);
    font-size: var(--font-size-sm);
    outline: none;
    transition: all 0.15s ease-in-out;
    font-family: inherit;
  }

  .custom-dir-input::placeholder {
    color: var(--text-muted);
  }

  .custom-dir-input:focus {
    border-color: var(--border-focus);
    box-shadow: 0 0 10px var(--accent-glow);
  }
</style>
