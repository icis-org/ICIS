<script lang="ts">
  interface Props {
    message: string;
    type: 'success' | 'error' | 'warning';
    ondismiss: () => void;
  }

  let { message, type, ondismiss }: Props = $props();
</script>

<div class="toast {type}">
  <span>{message}</span>
  <button class="toast-dismiss" onclick={ondismiss} aria-label="Dismiss">
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
  </button>
</div>

<style>
  .toast {
    position: fixed;
    bottom: 20px;
    right: 20px;
    padding: 12px 16px;
    border-radius: var(--radius-sm);
    font-size: var(--font-size-sm);
    font-weight: 500;
    z-index: 1000;
    animation: slideUp 0.25s ease-in-out;
    display: flex;
    align-items: center;
    gap: var(--spacing-sm);
    max-width: 400px;
    background: rgba(30, 32, 46, 0.9);
    backdrop-filter: blur(16px);
    -webkit-backdrop-filter: blur(16px);
    border: 1px solid var(--border);
    box-shadow: 0 4px 24px rgba(0, 0, 0, 0.4);
  }

  .toast.success {
    color: var(--success);
    border-color: rgba(52, 211, 153, 0.3);
  }

  .toast.error {
    color: var(--danger);
    border-color: rgba(248, 113, 113, 0.3);
  }

  .toast.warning {
    color: var(--warning);
    border-color: rgba(251, 191, 36, 0.3);
  }

  .toast-dismiss {
    background: none;
    border: none;
    color: inherit;
    cursor: pointer;
    padding: 2px;
    display: flex;
    opacity: 0.6;
    transition: opacity 0.15s ease-in-out;
  }

  .toast-dismiss:hover {
    opacity: 1;
  }

  @keyframes slideUp {
    from { transform: translateY(12px); opacity: 0; }
    to { transform: translateY(0); opacity: 1; }
  }
</style>
