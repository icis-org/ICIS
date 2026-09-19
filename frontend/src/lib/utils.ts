import type { iciparser } from '../../wailsjs/go/models';

export function formatBytes(bytes: number | undefined | null): string {
  if (!bytes || bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

export function buildICIContent(ici: iciparser.ICIFile, selectedDir: string, customDir?: string): string {
  const lines: string[] = [];
  lines.push(`name: ${ici.name}`);
  if (ici.version) lines.push(`version: ${ici.version}`);
  if (ici.desc) lines.push(`desc: ${ici.desc}`);
  lines.push(`url: ${ici.url}`);
  if (ici.type) lines.push(`type: ${ici.type}`);
  lines.push(`install_dir: ${selectedDir || 'appdata'}`);
  if (selectedDir === 'custom' && customDir) {
    lines.push(`custom_dir: ${customDir}`);
  }
  if (ici.shortcuts && ici.shortcuts.length > 0) {
    const entries = ici.shortcuts.map((s) => {
      const exe = s.exe || s.name;
      const name = s.name || s.exe;
      return exe === name ? exe : `${exe}=${name}`;
    });
    lines.push(`shortcuts: ${entries.join(', ')}`);
  } else if (ici.shortcut) {
    lines.push(`shortcut: ${ici.shortcut}`);
  }
  lines.push(`startup: ${ici.startup || false}`);
  return lines.join('\n');
}

export function buildICIPreview(
  name: string,
  version: string,
  desc: string,
  url: string,
  type: string,
  installDir: string,
  shortcutsText: string,
  startup: boolean
): string {
  const lines: string[] = [];
  if (name) lines.push(`name: ${name}`);
  if (version) lines.push(`version: ${version}`);
  if (desc) lines.push(`desc: ${desc}`);
  if (url) lines.push(`url: ${url}`);
  if (type) lines.push(`type: ${type}`);
  lines.push(`install_dir: ${installDir}`);
  if (shortcutsText) {
    const entries = shortcutsText.split('\n').map((s) => s.trim()).filter((s) => s);
    if (entries.length) lines.push(`shortcuts: ${entries.join(', ')}`);
  }
  lines.push(`startup: ${startup}`);
  return lines.join('\n');
}
