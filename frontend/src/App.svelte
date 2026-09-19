<script lang="ts">
  import { onMount } from 'svelte';
  import { app, toasts } from './lib/state.svelte';
  import {
    IsWizardMode,
    GetPendingFile,
    GetPendingProtocolURL,
    IsAutoInstall,
    LoadICIFile,
    LoadRegistryICI,
    EnterWizardMode,
  } from '../wailsjs/go/main/App.js';
  import { EventsOn } from '../wailsjs/runtime/runtime.js';
  import TitleBar from './components/TitleBar.svelte';
  import Sidebar from './components/Sidebar.svelte';
  import Toast from './components/Toast.svelte';
  import Home from './screens/Home.svelte';
  import Create from './screens/Create.svelte';
  import Install from './screens/Install.svelte';
  import Installed from './screens/Installed.svelte';
  import Store from './screens/Store.svelte';
  import WizardLayout from './wizard/WizardLayout.svelte';

  let ready = $state(false);

  onMount(async () => {
    try {
      app.wizardMode = await IsWizardMode();
    } catch {
      app.wizardMode = false;
    }

    listenEvents();
    await checkPendingFile();

    if (app.wizardMode) {
      await checkPendingProtocol();
    }

    ready = true;
  });

  async function checkPendingFile() {
    try {
      const pending = await GetPendingFile();
      if (pending && pending.path) {
        const ici = await LoadICIFile(pending.path);
        if (ici) {
          app.setICI(ici, pending.path);
          if (pending.auto) {
            app.wizardMode = true;
          }
        }
      }
    } catch (e) {
      toasts.add('Failed to load .ici file: ' + (e instanceof Error ? e.message : e), 'error');
    }
  }

  async function checkPendingProtocol() {
    try {
      const protocolURL = await GetPendingProtocolURL();
      if (protocolURL) {
        toasts.add('Loading package from: ' + protocolURL, 'warning');
        const ici = await LoadRegistryICI(protocolURL);
        if (ici) {
          app.setICI(ici, protocolURL);
        } else {
          toasts.add('Failed to load package: no data returned', 'error');
        }
      }
    } catch (e) {
      toasts.add('Failed to load package: ' + (e instanceof Error ? e.message : e), 'error');
    }
  }

  function listenEvents() {
    EventsOn('ici-loaded', async (ici: any) => {
      app.setICI(ici, app.pendingSource);
      try {
        const auto = await IsAutoInstall();
        if (auto) {
          setTimeout(() => {
            app.installing = true;
          }, 500);
        }
      } catch {}
    });

    EventsOn('protocol-ici', async (iciURL: string) => {
      try {
        const ici = await LoadRegistryICI(iciURL);
        if (ici) {
          if (!app.wizardMode) {
            app.wizardMode = true;
            await EnterWizardMode();
          }
          app.setICI(ici, iciURL);
        }
      } catch (e) {
        toasts.add('Failed to load package: ' + e, 'error');
      }
    });

    EventsOn('download-progress', (data: any) => {
      const event = new CustomEvent('icis:download-progress', { detail: data });
      window.dispatchEvent(event);
    });

    EventsOn('install-progress', (data: any) => {
      const event = new CustomEvent('icis:install-progress', { detail: data });
      window.dispatchEvent(event);

      if (data.status === 'shortcut-warning') {
        toasts.add(data.message || data.status, 'warning');
      }
    });

    EventsOn('extract-progress', (data: any) => {
      const event = new CustomEvent('icis:extract-progress', { detail: data });
      window.dispatchEvent(event);
    });

    EventsOn('install-complete', (data: any) => {
      const event = new CustomEvent('icis:install-complete', { detail: data });
      window.dispatchEvent(event);
    });

    EventsOn('uninstall-complete', (data: any) => {
      toasts.add(data.message, 'success');
    });
  }
</script>

{#if !ready}
  <div class="loading-screen"></div>
{:else if app.wizardMode}
  <WizardLayout />
{:else}
  <div class="app-root">
    <TitleBar />
    <div class="app-body">
      <Sidebar />
      <main class="main-content">
        {#if app.currentScreen === 'home'}
          <Home />
        {:else if app.currentScreen === 'create'}
          <Create />
        {:else if app.currentScreen === 'install'}
          <Install />
        {:else if app.currentScreen === 'installed'}
          <Installed />
        {:else if app.currentScreen === 'store'}
          <Store />
        {/if}
      </main>
    </div>
  </div>
{/if}

{#each toasts.items as t (t.id)}
  <Toast message={t.message} type={t.type} ondismiss={() => toasts.remove(t.id)} />
{/each}

<style>
  .loading-screen {
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;
    background: transparent;
  }

  .app-root {
    height: 100vh;
    display: flex;
    flex-direction: column;
    background: var(--bg-base);
    border: 1px solid var(--border);
    border-radius: var(--radius);
    box-shadow: 0 20px 40px rgba(0, 0, 0, 0.5);
    overflow: hidden;
    box-sizing: border-box;
  }

  .app-body {
    flex: 1;
    display: flex;
    overflow: hidden;
  }

  .main-content {
    flex: 1;
    overflow-y: auto;
    padding: var(--spacing-xl);
    background: var(--bg-base);
  }
</style>
