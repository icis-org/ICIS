<script lang="ts">
  import ProgressBar from '../components/ProgressBar.svelte';

  interface Props {
    percent: number;
    text: string;
    warnings: string[];
  }

  let { percent, text, warnings }: Props = $props();
</script>

<div class="wizard-page active">
  <div class="wizard-card glass-card">
    <h2>Installing...</h2>
    <ProgressBar {percent} {text} />
    {#if warnings.length > 0}
      <div class="warnings">
        {#each warnings as w}
          <div class="warning-item">{w}</div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .wizard-page {
    display: none;
    height: 100%;
    padding: 0;
    justify-content: center;
    align-items: flex-start;
    box-sizing: border-box;
  }

  .wizard-page.active {
    display: flex;
  }

  .glass-card {
    width: 100%;
    max-width: 560px;
    background: rgba(22, 24, 34, 0.85);
    backdrop-filter: blur(20px) saturate(180%);
    -webkit-backdrop-filter: blur(20px) saturate(180%);
    border: 1px solid var(--border);
    border-top: 1px solid var(--border-highlight);
    border-radius: var(--radius-md);
    padding: var(--spacing-lg);
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    --wails-draggable: drag;
    box-sizing: border-box;
  }

  h2 {
    font-size: var(--font-size-xl);
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 var(--spacing-md) 0;
  }

  .warnings {
    margin-top: var(--spacing-md);
  }

  .warning-item {
    color: var(--warning);
    font-size: var(--font-size-xs);
    margin-top: 4px;
  }
</style>
