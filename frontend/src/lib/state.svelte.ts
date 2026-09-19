import type { iciparser, db } from '../../wailsjs/go/models';

export type Screen = 'home' | 'create' | 'install' | 'installed' | 'store';
export type WizardStep = 'confirm' | 'progress' | 'done';

export const app = $state({
  wizardMode: false,
  currentScreen: 'home' as Screen,
  currentICI: null as iciparser.ICIFile | null,
  pendingSource: '',
  selectedDir: 'appdata',
  customDir: '',
  wizardStep: 'confirm' as WizardStep,
  installing: false,

  setScreen(screen: Screen) {
    this.currentScreen = screen;
  },

  setWizardMode(on: boolean) {
    this.wizardMode = on;
  },

  setICI(ici: iciparser.ICIFile, source: string = '') {
    this.currentICI = ici;
    this.pendingSource = source;
  },

  setDir(dir: string) {
    this.selectedDir = dir;
  },

  setWizardStep(step: WizardStep) {
    this.wizardStep = step;
  },

  setInstalling(v: boolean) {
    this.installing = v;
  },

  reset() {
    this.currentICI = null;
    this.pendingSource = '';
    this.selectedDir = 'appdata';
    this.customDir = '';
    this.wizardStep = 'confirm';
    this.installing = false;
  },
});

export const toasts = $state({
  items: [] as Array<{ id: number; message: string; type: 'success' | 'error' | 'warning' }>,
  private_nextId: 0,

  add(message: string, type: 'success' | 'error' | 'warning' = 'success') {
    const id = this.private_nextId++;
    this.items.push({ id, message, type });
    setTimeout(() => this.remove(id), 4000);
  },

  remove(id: number) {
    const idx = this.items.findIndex((t) => t.id === id);
    if (idx !== -1) this.items.splice(idx, 1);
  },
});

export const installedApps = $state({
  apps: [] as db.InstalledApp[],
  loading: false,
});
