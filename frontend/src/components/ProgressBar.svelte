<script lang="ts">
  interface Props {
    percent?: number;
    text?: string;
    status?: string;
  }

  let { percent = 0, text = 'Preparing...', status = '' }: Props = $props();
</script>

<div class="progress-container">
  <div class="progress-bar-bg">
    <div class="progress-bar-fill" class:success={percent >= 100} style="width: {percent}%"></div>
  </div>
  <div class="progress-status">
    <span>{text}</span>
    <span>{Math.round(percent)}%</span>
  </div>
  {#if status}
    <div class="progress-extra">{status}</div>
  {/if}
</div>

<style>
  .progress-container {
    margin-top: var(--spacing-lg);
  }

  .progress-bar-bg {
    height: 6px;
    background: rgba(255, 255, 255, 0.06);
    border-radius: 3px;
    overflow: hidden;
    margin: var(--spacing-sm) 0;
  }

  .progress-bar-fill {
    height: 100%;
    background: var(--accent-gradient);
    border-radius: 3px;
    transition: width 0.3s ease-in-out;
    box-shadow: 0 0 8px var(--accent-glow);
  }

  .progress-bar-fill.success {
    background: var(--success);
    box-shadow: 0 0 8px rgba(52, 211, 153, 0.3);
  }

  .progress-status {
    font-size: var(--font-size-xs);
    color: var(--text-secondary);
    display: flex;
    justify-content: space-between;
  }

  .progress-extra {
    font-size: var(--font-size-xs);
    color: var(--warning);
    margin-top: var(--spacing-xs);
  }
</style>
