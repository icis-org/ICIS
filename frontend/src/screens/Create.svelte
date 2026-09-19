<script lang="ts">
  import { app, toasts } from '../lib/state.svelte';
  import { LoadICIContent } from '../../wailsjs/go/main/App.js';
  import { buildICIPreview } from '../lib/utils';

  let name = $state('');
  let version = $state('');
  let desc = $state('');
  let url = $state('');
  let type = $state('zip');
  let installDir = $state('appdata');
  let shortcutsText = $state('');
  let startup = $state(false);

  let preview = $derived(
    buildICIPreview(name, version, desc, url, type, installDir, shortcutsText, startup)
  );

  const typeOptions = ['zip', '7z', 'rar', 'exe', 'msi'];
  const typeLabels: Record<string, string> = { zip: 'ZIP', '7z': '7Z', rar: 'RAR', exe: 'EXE', msi: 'MSI' };
  const dirOptions = ['appdata', 'programfiles', 'custom'];
  const dirLabels: Record<string, string> = { appdata: 'AppData', programfiles: 'Program Files', custom: 'Custom' };

  function downloadICi() {
    if (!preview) return;
    const blob = new Blob([preview], { type: 'text/plain' });
    const a = document.createElement('a');
    a.href = URL.createObjectURL(blob);
    a.download = (name || 'app') + '.ici';
    a.click();
  }

  async function installFromCreate() {
    if (!preview) {
      toasts.add('Fill in the form first', 'error');
      return;
    }
    try {
      const ici = await LoadICIContent(preview);
      if (ici) {
        app.setICI(ici, '');
        app.setScreen('install');
      }
    } catch (err) {
      toasts.add(String(err), 'error');
    }
  }
</script>

<div class="screen active">
  <div class="ici-editor glass-card">
    <h2>Create .ici File</h2>

    <div class="form-row">
      <div class="form-group">
        <label for="ici-name">App Name</label>
        <input type="text" id="ici-name" placeholder="e.g. Notepad++" bind:value={name} />
      </div>
      <div class="form-group">
        <label for="ici-version">Version</label>
        <input type="text" id="ici-version" placeholder="e.g. 8.7.1" bind:value={version} />
      </div>
    </div>

    <div class="form-group">
      <label for="ici-desc">Description</label>
      <input type="text" id="ici-desc" placeholder="Free source code editor" bind:value={desc} />
    </div>

    <div class="form-group">
      <label for="ici-url">Download URL</label>
      <input type="text" id="ici-url" placeholder="https://github.com/..." bind:value={url} />
    </div>

    <div class="form-group">
      <label>Archive Type</label>
      <div class="segment-group">
        {#each typeOptions as opt}
          <button
            class="segment-item"
            class:selected={type === opt}
            onclick={() => type = opt}
          >{typeLabels[opt]}</button>
        {/each}
      </div>
    </div>

    <div class="form-group">
      <label>Install Location</label>
      <div class="segment-group">
        {#each dirOptions as opt}
          <button
            class="segment-item"
            class:selected={installDir === opt}
            onclick={() => installDir = opt}
          >{dirLabels[opt]}</button>
        {/each}
      </div>
    </div>

    <div class="form-group">
      <label for="ici-shortcuts">Shortcuts (one per line: exe.exe=Display Name)</label>
      <textarea
        id="ici-shortcuts"
        placeholder={"mpv.exe=MPV Player\ntools/recorder.exe=Recorder"}
        bind:value={shortcutsText}
      ></textarea>
    </div>

    <div class="form-group">
      <div class="checkbox-row">
        <input type="checkbox" id="ici-startup" bind:checked={startup} />
        <label for="ici-startup">Add to Startup</label>
      </div>
    </div>

    <div class="form-group">
      <label for="ici-preview">Output Preview</label>
      <textarea id="ici-preview" readonly value={preview}></textarea>
    </div>

    <div class="install-actions">
      <button class="btn btn-secondary" onclick={downloadICi}>Download .ici</button>
      <button class="btn btn-primary" onclick={installFromCreate}>Install</button>
    </div>
  </div>
</div>

<style>
  .screen {
    display: none;
    height: 100%;
  }

  .screen.active {
    display: flex;
    flex-direction: column;
    overflow: hidden;
    animation: fadeIn 0.2s ease-in-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(6px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .glass-card {
    flex: 1;
    min-height: 0;
    overflow-y: auto;
    background: rgba(30, 32, 46, 0.6);
    backdrop-filter: blur(16px) saturate(180%);
    -webkit-backdrop-filter: blur(16px) saturate(180%);
    border: 1px solid var(--border);
    border-top: 1px solid var(--border-highlight);
    border-radius: var(--radius-md);
    padding: var(--spacing-xl);
    max-width: 600px;
    margin: 0 auto;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    width: 100%;
  }

  h2 {
    font-size: var(--font-size-lg);
    font-weight: 700;
    margin: 0 0 var(--spacing-lg) 0;
    color: var(--text-primary);
  }

  textarea {
    width: 100%;
    height: 80px;
    padding: 10px 14px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid var(--border);
    border-radius: var(--radius-sm);
    color: var(--text-primary);
    font-size: var(--font-size-xs);
    font-family: 'SFMono-Regular', Consolas, monospace;
    resize: none;
    outline: none;
    box-sizing: border-box;
    transition: all 0.15s ease-in-out;
  }

  textarea:focus {
    border-color: var(--border-focus);
    box-shadow: 0 0 10px var(--accent-glow);
  }

  textarea[readonly] {
    height: 160px;
    color: var(--text-secondary);
  }

  .install-actions {
    display: flex;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-lg);
    justify-content: flex-end;
  }
</style>
