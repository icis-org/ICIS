<script lang="ts">
  import { app, toasts } from '../lib/state.svelte';
  import { InstallApp } from '../../wailsjs/go/main/App.js';
  import { Quit } from '../../wailsjs/runtime/runtime.js';
  import { buildICIContent } from '../lib/utils';
  import WizardConfirm from './WizardConfirm.svelte';
  import WizardProgress from './WizardProgress.svelte';
  import WizardDone from './WizardDone.svelte';

  let progressPercent = $state(0);
  let progressText = $state('Preparing...');
  let warnings = $state<string[]>([]);
  let doneName = $state('');
  let donePath = $state('');

  async function doInstall(customDir: string) {
    if (!app.currentICI) return;

    app.setInstalling(true);
    app.setWizardStep('progress');
    progressPercent = 0;
    progressText = 'Preparing installation...';
    warnings = [];

    const iciContent = buildICIContent(app.currentICI, app.selectedDir, customDir);

    try {
      await InstallApp(iciContent, customDir);
    } catch (err) {
      toasts.add(String(err), 'error');
      app.setWizardStep('confirm');
      app.setInstalling(false);
    }
  }

  function handleCancel() {
    try { Quit(); } catch { window.close(); }
  }

  function handleFinish() {
    try { Quit(); } catch { window.close(); }
  }

  $effect(() => {
    const onDownload = (e: Event) => {
      const data = (e as CustomEvent).detail;
      progressPercent = data.percent || 0;
      progressText = `Downloading... ${fmtBytes(data.downloaded)} / ${fmtBytes(data.total)}`;
    };
    const onInstall = (e: Event) => {
      const data = (e as CustomEvent).detail;
      progressText = data.message || data.status;
      if (data.status === 'shortcut-warning') {
        warnings = [...warnings, data.message];
      }
    };
    const onExtract = (e: Event) => {
      const data = (e as CustomEvent).detail;
      progressText = `Extracting: ${data.file} (${data.current}/${data.total})`;
    };
    const onComplete = (e: Event) => {
      const data = (e as CustomEvent).detail;
      doneName = data.name;
      donePath = data.path;
      progressPercent = 100;
      app.setWizardStep('done');
      app.setInstalling(false);
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

  function fmtBytes(bytes: number | undefined | null): string {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
  }
</script>

<div class="wizard-root">
  <div class="wizard-header">
    <h1>ICIS Installer</h1>
  </div>
  <div class="wizard-body">
    {#if app.wizardStep === 'confirm'}
      <WizardConfirm oninstall={doInstall} oncancel={handleCancel} />
    {:else if app.wizardStep === 'progress'}
      <WizardProgress percent={progressPercent} text={progressText} {warnings} />
    {:else if app.wizardStep === 'done'}
      <WizardDone name={doneName} path={donePath} onfinish={handleFinish} />
    {/if}
  </div>
</div>

<style>
  .wizard-root {
    width: 100vw;
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: var(--bg-base);
    overflow: hidden;
  }

  .wizard-header {
    padding: var(--spacing-md) var(--spacing-lg);
    background: rgba(22, 24, 34, 0.85);
    backdrop-filter: blur(16px) saturate(180%);
    -webkit-backdrop-filter: blur(16px) saturate(180%);
    border-bottom: 1px solid var(--border);
    flex-shrink: 0;
  }

  .wizard-header h1 {
    font-size: var(--font-size-lg);
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
  }

  .wizard-body {
    flex: 1;
    overflow-y: auto;
    position: relative;
  }
</style>
