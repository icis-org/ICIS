import './style.css';
import './app.css';
import {
    LoadICIContent,
    InstallApp,
    UninstallApp,
    ListInstalledApps,
    OpenFile,
    SelectDirectory,
} from '../wailsjs/go/main/App.js';
import { EventsOn } from '../wailsjs/runtime/runtime.js';

let currentScreen = 'home';
let currentICI = null;
let selectedDir = 'appdata';

function init() {
    document.querySelector('#app').innerHTML = layout();
    setupNav();
    setupHome();
    setupInstall();
    setupCreate();
    setupInstalled();
    listenEvents();
}

function layout() {
    return `
    <div class="sidebar">
        <div class="sidebar-header">
            <h1>ICIS</h1>
            <div class="subtitle">Configurable Install System</div>
        </div>
        <div class="nav-item active" data-screen="home">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M3 9l9-7 9 7v11a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/><polyline points="9 22 9 12 15 12 15 22"/></svg>
            Home
        </div>
        <div class="nav-item" data-screen="create">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="12" y1="18" x2="12" y2="12"/><line x1="9" y1="15" x2="15" y2="15"/></svg>
            Create .ici
        </div>
        <div class="nav-item" data-screen="installed">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
            Installed
        </div>
        <div class="sidebar-footer">
            <div class="version">ICIS v1.0.0</div>
        </div>
    </div>
    <div class="main-content">
        ${screenHome()}
        ${screenCreate()}
        ${screenInstall()}
        ${screenInstalled()}
    </div>`;
}

function screenHome() {
    return `
    <div id="screen-home" class="screen active">
        <div class="welcome">
            <div class="welcome-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
            </div>
            <h2>Welcome to ICIS</h2>
            <p>Open an <strong>.ici</strong> file to install an app, or create your own configuration file.</p>
            <div style="display:flex;gap:10px;margin-top:10px">
                <button class="btn btn-primary" id="btn-open-ici">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16"><path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/></svg>
                    Open .ici File
                </button>
                <button class="btn btn-secondary" id="btn-paste-ici">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16"><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><rect x="8" y="2" width="8" height="4" rx="1" ry="1"/></svg>
                    Paste .ici Content
                </button>
            </div>
            <div id="paste-modal" style="display:none;margin-top:15px;width:400px;text-align:left">
                <textarea id="paste-input" placeholder="Paste .ici content here..." style="width:100%;height:150px;padding:12px;background:var(--bg-tertiary);border:1px solid var(--border);border-radius:var(--radius-sm);color:var(--text-primary);font-size:12px;font-family:monospace;resize:vertical;outline:none"></textarea>
                <div style="display:flex;gap:8px;margin-top:8px;justify-content:flex-end">
                    <button class="btn btn-secondary btn-sm" id="btn-cancel-paste">Cancel</button>
                    <button class="btn btn-primary btn-sm" id="btn-submit-paste">Load</button>
                </div>
            </div>
        </div>
    </div>`;
}

function screenCreate() {
    return `
    <div id="screen-create" class="screen">
        <div class="ici-editor">
            <h2>Create .ici File</h2>
            <div class="form-row">
                <div class="form-group">
                    <label>App Name *</label>
                    <input type="text" id="ici-name" placeholder="e.g. Notepad++"/>
                </div>
                <div class="form-group">
                    <label>Version</label>
                    <input type="text" id="ici-version" placeholder="e.g. 8.7.1"/>
                </div>
            </div>
            <div class="form-group">
                <label>Description</label>
                <input type="text" id="ici-desc" placeholder="Free source code editor"/>
            </div>
            <div class="form-group">
                <label>Download URL *</label>
                <input type="text" id="ici-url" placeholder="https://github.com/..."/>
            </div>
            <div class="form-row">
                <div class="form-group">
                    <label>Archive Type</label>
                    <select id="ici-type">
                        <option value="zip">ZIP</option>
                        <option value="7z">7-Zip</option>
                        <option value="rar">RAR</option>
                        <option value="exe">EXE (Installer)</option>
                        <option value="msi">MSI</option>
                    </select>
                </div>
                <div class="form-group">
                    <label>Install Location</label>
                    <select id="ici-install-dir">
                        <option value="appdata">AppData</option>
                        <option value="programfiles">Program Files</option>
                        <option value="custom">Custom</option>
                    </select>
                </div>
            </div>
            <div class="form-group">
                <label>Shortcut Name</label>
                <input type="text" id="ici-shortcut" placeholder="Name for desktop shortcut"/>
            </div>
            <div class="form-group">
                <div class="checkbox-row">
                    <input type="checkbox" id="ici-startup"/>
                    <label for="ici-startup">Add to Startup</label>
                </div>
            </div>
            <div class="form-group">
                <label>Output Preview</label>
                <textarea id="ici-preview" readonly style="width:100%;height:160px;padding:12px;background:var(--bg-tertiary);border:1px solid var(--border);border-radius:var(--radius-sm);color:var(--text-secondary);font-size:12px;font-family:monospace;resize:vertical;outline:none"></textarea>
            </div>
            <div class="install-actions">
                <button class="btn btn-secondary" id="btn-download-ici">Download .ici</button>
                <button class="btn btn-primary" id="btn-install-from-create">Install</button>
            </div>
        </div>
    </div>`;
}

function screenInstall() {
    return `
    <div id="screen-install" class="screen">
        <div class="install-card">
            <h2 id="install-name"></h2>
            <div class="version" id="install-version"></div>
            <div class="desc" id="install-desc"></div>
            <div class="detail-row">
                <span class="detail-label">Download URL</span>
                <span class="detail-value" id="install-url" style="max-width:300px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap"></span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Archive Type</span>
                <span class="detail-value" id="install-type"></span>
            </div>
            <div style="margin-top:15px">
                <label style="font-size:12px;font-weight:600;color:var(--text-secondary);text-transform:uppercase;letter-spacing:0.5px">Install Location</label>
                <div class="install-dir-select" id="dir-select">
                    <div class="dir-option selected" data-dir="appdata">AppData</div>
                    <div class="dir-option" data-dir="programfiles">Program Files</div>
                    <div class="dir-option" data-dir="custom">Custom...</div>
                </div>
                <div class="custom-dir-row" id="custom-dir-row" style="display:none">
                    <input class="custom-dir-input" id="custom-dir-input" placeholder="Select custom directory..."/>
                    <button class="btn btn-secondary" id="btn-browse-dir" style="font-size:12px">Browse</button>
                </div>
            </div>
            <div class="progress-container" id="progress-section" style="display:none">
                <div class="progress-bar-bg">
                    <div class="progress-bar-fill" id="progress-bar"></div>
                </div>
                <div class="progress-status">
                    <span id="progress-text">Preparing...</span>
                    <span id="progress-percent">0%</span>
                </div>
            </div>
            <div class="install-actions" id="install-actions">
                <button class="btn btn-secondary" id="btn-back-install">Back</button>
                <button class="btn btn-primary" id="btn-do-install">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                    Install
                </button>
            </div>
        </div>
    </div>`;
}

function screenInstalled() {
    return `
    <div id="screen-installed" class="screen">
        <div class="page-header">
            <h2>Installed Apps</h2>
            <button class="btn btn-secondary" id="btn-refresh-installed">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14"><polyline points="23 4 23 10 17 10"/><polyline points="1 20 1 14 7 14"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
                Refresh
            </button>
        </div>
        <div id="installed-list" class="installed-list"></div>
    </div>`;
}

function setupNav() {
    document.querySelectorAll('.nav-item').forEach(item => {
        item.addEventListener('click', () => {
            const screen = item.dataset.screen;
            switchScreen(screen);
        });
    });
}

function switchScreen(name) {
    currentScreen = name;
    document.querySelectorAll('.nav-item').forEach(n => n.classList.remove('active'));
    const navItem = document.querySelector(`.nav-item[data-screen="${name}"]`);
    if (navItem) navItem.classList.add('active');
    document.querySelectorAll('.screen').forEach(s => s.classList.remove('active'));
    const screen = document.getElementById(`screen-${name}`);
    if (screen) screen.classList.add('active');

    if (name === 'installed') {
        refreshInstalled();
    }
}

function setupHome() {
    document.getElementById('btn-open-ici').addEventListener('click', async () => {
        try {
            const filePath = await OpenFile();
            if (filePath) {
                const { LoadICIFile } = await import('../wailsjs/go/main/App.js');
                const ici = await LoadICIFile(filePath);
                if (ici) {
                    showInstallScreen(ici);
                }
            }
        } catch (err) {
            showToast(err, 'error');
        }
    });

    document.getElementById('btn-paste-ici').addEventListener('click', () => {
        const modal = document.getElementById('paste-modal');
        modal.style.display = modal.style.display === 'none' ? 'block' : 'none';
    });

    document.getElementById('btn-cancel-paste').addEventListener('click', () => {
        document.getElementById('paste-modal').style.display = 'none';
        document.getElementById('paste-input').value = '';
    });

    document.getElementById('btn-submit-paste').addEventListener('click', async () => {
        const content = document.getElementById('paste-input').value.trim();
        if (!content) return;
        try {
            const ici = await LoadICIContent(content);
            if (ici) {
                document.getElementById('paste-modal').style.display = 'none';
                document.getElementById('paste-input').value = '';
                showInstallScreen(ici);
            }
        } catch (err) {
            showToast(err, 'error');
        }
    });
}

function setupInstall() {
    document.querySelectorAll('.dir-option').forEach(opt => {
        opt.addEventListener('click', () => {
            document.querySelectorAll('.dir-option').forEach(o => o.classList.remove('selected'));
            opt.classList.add('selected');
            selectedDir = opt.dataset.dir;
            document.getElementById('custom-dir-row').style.display = selectedDir === 'custom' ? 'flex' : 'none';
        });
    });

    document.getElementById('btn-browse-dir').addEventListener('click', async () => {
        try {
            const dir = await SelectDirectory();
            if (dir) {
                document.getElementById('custom-dir-input').value = dir;
            }
        } catch (err) {
            console.error(err);
        }
    });

    document.getElementById('btn-back-install').addEventListener('click', () => {
        switchScreen('home');
    });

    document.getElementById('btn-do-install').addEventListener('click', doInstall);
}

function setupCreate() {
    const fields = ['ici-name', 'ici-version', 'ici-desc', 'ici-url', 'ici-type', 'ici-install-dir', 'ici-shortcut', 'ici-startup'];
    fields.forEach(id => {
        const el = document.getElementById(id);
        if (el) {
            el.addEventListener('input', updateCreatePreview);
            el.addEventListener('change', updateCreatePreview);
        }
    });

    document.getElementById('btn-download-ici').addEventListener('click', () => {
        const preview = document.getElementById('ici-preview').value;
        if (!preview) return;
        const blob = new Blob([preview], { type: 'text/plain' });
        const a = document.createElement('a');
        a.href = URL.createObjectURL(blob);
        a.download = (document.getElementById('ici-name').value || 'app') + '.ici';
        a.click();
    });

    document.getElementById('btn-install-from-create').addEventListener('click', async () => {
        const content = document.getElementById('ici-preview').value;
        if (!content) {
            showToast('Fill in the form first', 'error');
            return;
        }
        try {
            const ici = await LoadICIContent(content);
            if (ici) {
                showInstallScreen(ici);
            }
        } catch (err) {
            showToast(err, 'error');
        }
    });
}

function updateCreatePreview() {
    const name = document.getElementById('ici-name').value;
    const version = document.getElementById('ici-version').value;
    const desc = document.getElementById('ici-desc').value;
    const url = document.getElementById('ici-url').value;
    const type = document.getElementById('ici-type').value;
    const installDir = document.getElementById('ici-install-dir').value;
    const shortcut = document.getElementById('ici-shortcut').value;
    const startup = document.getElementById('ici-startup').checked;

    let lines = [];
    if (name) lines.push(`name: ${name}`);
    if (version) lines.push(`version: ${version}`);
    if (desc) lines.push(`desc: ${desc}`);
    if (url) lines.push(`url: ${url}`);
    if (type) lines.push(`type: ${type}`);
    lines.push(`install_dir: ${installDir}`);
    if (shortcut) lines.push(`shortcut: ${shortcut}`);
    lines.push(`startup: ${startup}`);

    document.getElementById('ici-preview').value = lines.join('\n');
}

function setupInstalled() {
    document.getElementById('btn-refresh-installed').addEventListener('click', refreshInstalled);
}

function showInstallScreen(ici) {
    currentICI = ici;
    document.getElementById('install-name').textContent = ici.name;
    document.getElementById('install-version').textContent = ici.version || '';
    document.getElementById('install-desc').textContent = ici.desc || '';
    document.getElementById('install-url').textContent = ici.url;
    document.getElementById('install-type').textContent = (ici.type || 'zip').toUpperCase();

    document.getElementById('progress-section').style.display = 'none';
    document.getElementById('progress-bar').style.width = '0%';
    document.getElementById('progress-bar').classList.remove('success');
    document.getElementById('install-actions').style.display = 'flex';

    switchScreen('install');
}

async function doInstall() {
    if (!currentICI) return;

    let installDir = '';
    if (selectedDir === 'custom') {
        installDir = document.getElementById('custom-dir-input').value;
        if (!installDir) {
            showToast('Please select a custom directory', 'error');
            return;
        }
    }

    const btn = document.getElementById('btn-do-install');
    btn.disabled = true;
    btn.innerHTML = 'Installing...';

    document.getElementById('progress-section').style.display = 'block';
    document.getElementById('progress-text').textContent = 'Preparing installation...';
    document.getElementById('progress-percent').textContent = '0%';

    const iciContent = document.getElementById('ici-preview')?.value || buildICIContent(currentICI);

    try {
        await InstallApp(iciContent, installDir);
    } catch (err) {
        showToast(err, 'error');
        btn.disabled = false;
        btn.innerHTML = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg> Install`;
        document.getElementById('progress-section').style.display = 'none';
    }
}

function buildICIContent(ici) {
    let lines = [];
    lines.push(`name: ${ici.name}`);
    if (ici.version) lines.push(`version: ${ici.version}`);
    if (ici.desc) lines.push(`desc: ${ici.desc}`);
    lines.push(`url: ${ici.url}`);
    if (ici.type) lines.push(`type: ${ici.type}`);
    lines.push(`install_dir: ${selectedDir || 'appdata'}`);
    if (selectedDir === 'custom') {
        lines.push(`custom_dir: ${document.getElementById('custom-dir-input')?.value || ''}`);
    }
    if (ici.shortcut) lines.push(`shortcut: ${ici.shortcut}`);
    lines.push(`startup: ${ici.startup || false}`);
    return lines.join('\n');
}

async function refreshInstalled() {
    try {
        const apps = await ListInstalledApps();
        const list = document.getElementById('installed-list');

        if (!apps || apps.length === 0) {
            list.innerHTML = `
            <div class="empty-state">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
                <h3>No apps installed</h3>
                <p>Open an .ici file or create one to install your first app.</p>
            </div>`;
            return;
        }

        list.innerHTML = apps.map(app => `
            <div class="installed-item">
                <div class="app-info">
                    <span class="name">${app.name}</span>
                    <span class="meta">${app.version ? 'v' + app.version + ' · ' : ''}${app.installPath}</span>
                </div>
                <div class="app-actions">
                    <button class="btn btn-danger" onclick="window.uninstallApp('${app.name}')">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14"><polyline points="3 6 5 6 21 6"/><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"/></svg>
                        Uninstall
                    </button>
                </div>
            </div>
        `).join('');
    } catch (err) {
        showToast('Failed to load installed apps: ' + err, 'error');
    }
}

window.uninstallApp = async function(name) {
    if (!confirm(`Are you sure you want to uninstall ${name}?`)) return;
    try {
        await UninstallApp(name);
        showToast(`${name} uninstalled`, 'success');
        refreshInstalled();
    } catch (err) {
        showToast(err, 'error');
    }
};

function listenEvents() {
    EventsOn('ici-loaded', (ici) => {
        showInstallScreen(ici);
    });

    EventsOn('download-progress', (data) => {
        const percent = Math.round(data.percent || 0);
        document.getElementById('progress-bar').style.width = percent + '%';
        document.getElementById('progress-text').textContent = `Downloading... ${formatBytes(data.downloaded)} / ${formatBytes(data.total)}`;
        document.getElementById('progress-percent').textContent = percent + '%';
    });

    EventsOn('install-progress', (data) => {
        const statusMessages = {
            downloading: 'Downloading...',
            extracting: 'Extracting files...',
            shortcut: 'Creating shortcuts...',
            saving: 'Saving installation record...',
        };
        document.getElementById('progress-text').textContent = data.message || statusMessages[data.status] || data.status;
    });

    EventsOn('extract-progress', (data) => {
        document.getElementById('progress-text').textContent = `Extracting: ${data.file} (${data.current}/${data.total})`;
    });

    EventsOn('install-complete', (data) => {
        const bar = document.getElementById('progress-bar');
        bar.style.width = '100%';
        bar.classList.add('success');
        document.getElementById('progress-text').textContent = data.message;
        document.getElementById('progress-percent').textContent = '100%';

        const btn = document.getElementById('btn-do-install');
        btn.disabled = false;
        btn.innerHTML = `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg> Install`;

        showToast(data.message, 'success');

        setTimeout(() => {
            document.getElementById('progress-section').style.display = 'none';
            bar.style.width = '0%';
            bar.classList.remove('success');
        }, 3000);
    });

    EventsOn('uninstall-complete', (data) => {
        showToast(data.message, 'success');
    });
}

function formatBytes(bytes) {
    if (!bytes || bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

function showToast(message, type) {
    const existing = document.querySelector('.toast');
    if (existing) existing.remove();

    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    document.body.appendChild(toast);

    setTimeout(() => {
        toast.remove();
    }, 4000);
}

init();
