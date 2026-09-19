<script lang="ts">
  interface Props {
    name: string;
    path: string;
    onfinish: () => void;
  }

  let { name, path, onfinish }: Props = $props();
</script>

<div class="wizard-page active">
  <div class="wizard-card glass-card centered">
    <div class="done-icon">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" width="48" height="48">
        <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/>
        <polyline points="22 4 12 14.01 9 11.01"/>
      </svg>
    </div>
    <h2>{name}</h2>
    <p class="done-path">{path}</p>
    <div class="wizard-actions">
      <button class="btn btn-primary" onclick={onfinish}>Finish</button>
    </div>
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

  .glass-card button {
    --wails-draggable: no-drag;
  }

  .centered {
    text-align: center;
  }

  h2 {
    font-size: var(--font-size-xl);
    font-weight: 600;
    color: var(--text-primary);
    margin: 0 0 var(--spacing-xs) 0;
  }

  .done-icon {
    color: var(--success);
    margin-bottom: var(--spacing-md);
  }

  .done-path {
    color: var(--text-secondary);
    font-size: var(--font-size-sm);
    word-break: break-all;
  }

  .wizard-actions {
    display: flex;
    justify-content: center;
    gap: var(--spacing-sm);
    margin-top: var(--spacing-xl);
  }
</style>
