import './style.css';
import './app.css';
import {
    LoadICIContent,
    LoadICIFile,
    InstallApp,
    UninstallApp,
    ListInstalledApps,
    OpenFile,
    SelectDirectory,
    IsAutoInstall,
    GetPendingFile,
    IsWizardMode,
    EnterWizardMode,
    IsInstalling,
    GetPendingProtocolURL,
    GetRegistryApps,
    LoadRegistryICI,
    GetRegistryURL,
    SetRegistryURL,
} from '../wailsjs/go/main/App.js';
import { EventsOn, Quit } from '../wailsjs/runtime/runtime.js';

let currentScreen = 'home';
let currentICI = null;
let pendingSource = '';
let selectedDir = 'appdata';
let wizardMode = false;
let wizardStep = 'confirm';
let storeData = [];
let eventsBound = false;

function init() {
    detectWizardMode();
}

async function detectWizardMode() {
    try {
        wizardMode = await IsWizardMode();
    } catch (e) {
        wizardMode = false;
    }

    if (wizardMode) {
        document.querySelector('#app').innerHTML = wizardLayout();
        setupWizard();
        listenEvents();
        checkPendingProtocol();
        checkPendingFile();
    } else {
        document.querySelector('#app').innerHTML = layout();
        setupNav();
        setupHome();
        setupInstall();
        setupCreate();
        setupInstalled();
        setupStore();
        listenEvents();
        checkPendingFile();
    }
}

function layout() {
    return `
    <div class="sidebar">
        <div class="sidebar-header">
            <img src="./public/favicon.png" alt="ICIS" style="width:28px;height:28px;border-radius:4px;margin-bottom:4px"/>
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
        <div class="nav-item" data-screen="store">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>
            Store
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
        ${screenStore()}
    </div>`;
}

function wizardLayout() {
    return `
    <div class="wizard-root">
        <div class="wizard-header">
            <h1>ICIS Installer</h1>
        </div>
        <div class="wizard-body" id="wizard-body">
            ${wizardPageConfirm()}
            ${wizardPageProgress()}
            ${wizardPageDone()}
        </div>
    </div>`;
}

function wizardPageConfirm() {
    return `
    <div id="wizard-confirm" class="wizard-page active">
        <div class="wizard-card">
            <h2 id="wiz-name"></h2>
            <div class="version" id="wiz-version"></div>
            <div class="desc" id="wiz-desc"></div>
            <div class="detail-row">
                <span class="detail-label">Source</span>
                <span class="detail-value" id="wiz-source" style="max-width:350px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap"></span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Download</span>
                <span class="detail-value" id="wiz-url" style="max-width:350px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap"></span>
            </div>
            <div class="detail-row">
                <span class="detail-label">Type</span>
                <span class="detail-value" id="wiz-type"></span>
            </div>
            <div id="wiz-shortcuts-section" style="display:none;margin-top:12px">
                <label style="font-size:12px;font-weight:600;color:var(--text-secondary);text-transform:uppercase;letter-spacing:0.5px">Shortcuts</label>
                <div id="wiz-shortcuts-list" style="margin-top:6px;font-size:13px;color:var(--text-primary)"></div>
            </div>
            <div style="margin-top:15px">
                <label style="font-size:12px;font-weight:600;color:var(--text-secondary);text-transform:uppercase;letter-spacing:0.5px">Install Location</label>
                <div class="install-dir-select" id="wiz-dir-select">
                    <div class="dir-option selected" data-dir="appdata">AppData</div>
                    <div class="dir-option" data-dir="programfiles">Program Files</div>
                    <div class="dir-option" data-dir="custom">Custom...</div>
                </div>
                <div class="custom-dir-row" id="wiz-custom-dir-row" style="display:none">
                    <input class="custom-dir-input" id="wiz-custom-dir-input" placeholder="Select custom directory..."/>
                    <button class="btn btn-secondary" id="wiz-browse-dir" style="font-size:12px">Browse</button>
                </div>
            </div>
            <div class="wizard-actions">
                <button class="btn btn-secondary" id="wiz-cancel">Cancel</button>
                <button class="btn btn-primary" id="wiz-install">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="16" height="16"><path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>
                    Install
                </button>
            </div>
        </div>
    </div>`;
}

function wizardPageProgress() {
    return `
    <div id="wizard-progress" class="wizard-page">
        <div class="wizard-card">
            <h2>Installing <span id="wiz-prog-name"></span></h2>
            <div class="progress-container" style="display:block">
                <div class="progress-bar-bg">
                    <div class="progress-bar-fill" id="wiz-progress-bar"></div>
                </div>
                <div class="progress-status">
                    <span id="wiz-progress-text">Preparing...</span>
                    <span id="wiz-progress-percent">0%</span>
                </div>
            </div>
            <div id="wiz-progress-warnings" style="margin-top:12px"></div>
        </div>
    </div>`;
}

function wizardPageDone() {
    return `
    <div id="wizard-done" class="wizard-page">
        <div class="wizard-card" style="text-align:center">
            <div class="wizard-done-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="48" height="48"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
            </div>
            <h2 id="wiz-done-name"></h2>
            <p id="wiz-done-path" style="color:var(--text-secondary);font-size:13px;word-break:break-all"></p>
            <div class="wizard-actions" style="justify-content:center;margin-top:20px">
                <button class="btn btn-primary" id="wiz-done-finish">Finish</button>
            </div>
        </div>
    </div>`;
}

function wizardShowStep(step) {
    wizardStep = step;
    document.querySelectorAll('.wizard-page').forEach(p => p.classList.remove('active'));
    const page = document.getElementById('wizard-' + step);
    if (page) page.classList.add('active');
}

function setupWizard() {
    document.querySelectorAll('#wiz-dir-select .dir-option').forEach(opt => {
        opt.addEventListener('click', () => {
            document.querySelectorAll('#wiz-dir-select .dir-option').forEach(o => o.classList.remove('selected'));
            opt.classList.add('selected');
            selectedDir = opt.dataset.dir;
            document.getElementById('wiz-custom-dir-row').style.display = selectedDir === 'custom' ? 'flex' : 'none';
        });
    });

    const browseBtn = document.getElementById('wiz-browse-dir');
    if (browseBtn) {
        browseBtn.addEventListener('click', async () => {
            try {
                const dir = await SelectDirectory();
                if (dir) {
                    document.getElementById('wiz-custom-dir-input').value = dir;
                }
            } catch (e) {}
        });
    }

    document.getElementById('wiz-install').addEventListener('click', doWizardInstall);

    document.getElementById('wiz-cancel').addEventListener('click', async () => {
        if (wizardMode) {
            try { Quit(); } catch (e) { window.close(); }
        } else {
            switchScreen('home');
        }
    });

    document.getElementById('wiz-done-finish').addEventListener('click', async () => {
        if (wizardMode) {
            try { Quit(); } catch (e) { window.close(); }
        } else {
            switchScreen('home');
        }
    });
}

function showWizardConfirm(ici, source) {
    currentICI = ici;
    document.getElementById('wiz-name').textContent = ici.name;
    document.getElementById('wiz-version').textContent = ici.version || '';
    document.getElementById('wiz-desc').textContent = ici.desc || '';
    document.getElementById('wiz-source').textContent = source || '';
    document.getElementById('wiz-url').textContent = ici.url;
    document.getElementById('wiz-type').textContent = (ici.type || 'zip').toUpperCase();

    const sc = document.getElementById('wiz-shortcuts-section');
    const sl = document.getElementById('wiz-shortcuts-list');
    if (ici.shortcuts && ici.shortcuts.length > 0) {
        sc.style.display = 'block';
        sl.innerHTML = ici.shortcuts.map(s => `<div style="padding:4px 0">  ${s.name || s.Name || s}</div>`).join('');
    } else if (ici.shortcut) {
        sc.style.display = 'block';
        sl.innerHTML = `<div style="padding:4px 0">  ${ici.shortcut}</div>`;
    } else {
        sc.style.display = 'none';
    }

    selectedDir = 'appdata';
    document.querySelectorAll('#wiz-dir-select .dir-option').forEach(o => o.classList.remove('selected'));
    document.querySelector('#wiz-dir-select .dir-option[data-dir="appdata"]').classList.add('selected');
    document.getElementById('wiz-custom-dir-row').style.display = 'none';

    document.getElementById('wiz-progress-bar').style.width = '0%';
    document.getElementById('wiz-progress-bar').classList.remove('success');
    document.getElementById('wiz-progress-warnings').innerHTML = '';

    document.getElementById('wiz-prog-name').textContent = ici.name;
    wizardShowStep('confirm');
}

async function doWizardInstall() {
    if (!currentICI) return;

    let installDir = '';
    if (selectedDir === 'custom') {
        installDir = document.getElementById('wiz-custom-dir-input').value;
        if (!installDir) {
            showToast('Please select a custom directory', 'error');
            return;
        }
    }

    document.getElementById('wiz-progress-bar').style.width = '0%';
    document.getElementById('wiz-progress-text').textContent = 'Preparing installation...';
    document.getElementById('wiz-progress-percent').textContent = '0%';
    document.getElementById('wiz-progress-warnings').innerHTML = '';

    wizardShowStep('progress');

    const iciContent = buildICIContent(currentICI);

    try {
        await InstallApp(iciContent, installDir);
    } catch (err) {
        showToast(err, 'error');
        wizardShowStep('confirm');
    }
}

async function checkPendingProtocol() {
    try {
        const protocolURL = await GetPendingProtocolURL();
        if (protocolURL) {
            showToast('Loading package from: ' + protocolURL, 'warning');
            const ici = await LoadRegistryICI(protocolURL);
            if (ici) {
                showWizardConfirm(ici, protocolURL);
            } else {
                showToast('Failed to load package: no data returned', 'error');
            }
        }
    } catch (e) {
        showToast('Failed to load package: ' + (e.message || e), 'error');
    }
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
                <label>Shortcuts (one per line: exe.exe=Display Name)</label>
                <textarea id="ici-shortcuts" placeholder="mpv.exe=MPV Player&#10;tools/recorder.exe=Recorder" style="width:100%;height:80px;padding:12px;background:var(--bg-tertiary);border:1px solid var(--border);border-radius:var(--radius-sm);color:var(--text-primary);font-size:12px;font-family:monospace;resize:vertical;outline:none"></textarea>
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
            <div id="install-shortcuts-section" style="display:none;margin-top:12px">
                <label style="font-size:12px;font-weight:600;color:var(--text-secondary);text-transform:uppercase;letter-spacing:0.5px">Shortcuts</label>
                <div id="install-shortcuts-list" style="margin-top:6px;font-size:13px;color:var(--text-primary)"></div>
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

function screenStore() {
    return `
    <div id="screen-store" class="screen">
        <div class="page-header">
            <h2>Store</h2>
            <div style="display:flex;gap:8px">
                <button class="btn btn-secondary" id="btn-store-settings">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14"><circle cx="12" cy="12" r="3"/><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06A1.65 1.65 0 0 0 4.68 15a1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06A1.65 1.65 0 0 0 9 4.68a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06A1.65 1.65 0 0 0 19.32 9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"/></svg>
                    Settings
                </button>
                <button class="btn btn-secondary" id="btn-store-refresh">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" width="14" height="14"><polyline points="23 4 23 10 17 10"/><polyline points="1 20 1 14 7 14"/><path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/></svg>
                    Refresh
                </button>
            </div>
        </div>
        <div id="store-status" style="margin-bottom:12px;font-size:12px;color:var(--text-secondary)"></div>
        <input type="text" id="store-search" placeholder="Search packages..." style="width:100%;padding:10px 14px;background:var(--bg-tertiary);border:1px solid var(--border);border-radius:var(--radius-sm);color:var(--text-primary);font-size:13px;outline:none;margin-bottom:16px;box-sizing:border-box"/>
        <div id="store-list" class="installed-list"></div>
        <div id="store-settings-panel" style="display:none;margin-top:16px;padding:16px;background:var(--bg-secondary);border:1px solid var(--border);border-radius:var(--radius-md)">
            <label style="font-size:12px;font-weight:600;color:var(--text-secondary);text-transform:uppercase;letter-spacing:0.5px">Registry URL</label>
            <div style="display:flex;gap:8px;margin-top:8px">
                <input type="text" id="store-url-input" style="flex:1;padding:10px 14px;background:var(--bg-tertiary);border:1px solid var(--border);border-radius:var(--radius-sm);color:var(--text-primary);font-size:13px;outline:none"/>
                <button class="btn btn-primary" id="btn-store-url-save">Save</button>
            </div>
        </div>
    </div>`;
}

function setupHome() {
    document.getElementById('btn-open-ici').addEventListener('click', async () => {
        try {
            const filePath = await OpenFile();
            if (filePath) {
                pendingSource = filePath;
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
        } catch (e) {
            console.error(e);
        }
    });

    document.getElementById('btn-back-install').addEventListener('click', () => {
        switchScreen('home');
    });

    document.getElementById('btn-do-install').addEventListener('click', doInstall);
}

function setupCreate() {
    const fields = ['ici-name', 'ici-version', 'ici-desc', 'ici-url', 'ici-type', 'ici-install-dir', 'ici-shortcuts', 'ici-startup'];
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
    const shortcutsText = document.getElementById('ici-shortcuts').value.trim();
    const startup = document.getElementById('ici-startup').checked;

    let lines = [];
    if (name) lines.push(`name: ${name}`);
    if (version) lines.push(`version: ${version}`);
    if (desc) lines.push(`desc: ${desc}`);
    if (url) lines.push(`url: ${url}`);
    if (type) lines.push(`type: ${type}`);
    lines.push(`install_dir: ${installDir}`);
    if (shortcutsText) {
        const entries = shortcutsText.split('\n').map(s => s.trim()).filter(s => s);
        if (entries.length) lines.push(`shortcuts: ${entries.join(', ')}`);
    }
    lines.push(`startup: ${startup}`);

    document.getElementById('ici-preview').value = lines.join('\n');
}

function setupInstalled() {
    document.getElementById('btn-refresh-installed').addEventListener('click', refreshInstalled);
}

function setupStore() {
    document.getElementById('btn-store-refresh').addEventListener('click', () => loadStore());
    document.getElementById('btn-store-settings').addEventListener('click', () => {
        const panel = document.getElementById('store-settings-panel');
        panel.style.display = panel.style.display === 'none' ? 'block' : 'none';
    });
    document.getElementById('btn-store-url-save').addEventListener('click', async () => {
        const url = document.getElementById('store-url-input').value.trim();
        if (!url) return;
        try {
            await SetRegistryURL(url);
            showToast('Registry URL saved', 'success');
            document.getElementById('store-settings-panel').style.display = 'none';
            loadStore();
        } catch (e) {
            showToast('Failed to save URL: ' + e, 'error');
        }
    });
    document.getElementById('store-search').addEventListener('input', filterStore);

    loadStore();
}

async function loadStore() {
    const list = document.getElementById('store-list');
    const status = document.getElementById('store-status');
    list.innerHTML = '<div style="padding:20px;text-align:center;color:var(--text-secondary)">Loading...</div>';

    try {
        const result = await GetRegistryApps();
        storeData = result.index?.packages || [];
        status.textContent = result.error
            ? result.error
            : '';
        renderStore(storeData);
    } catch (e) {
        list.innerHTML = '<div class="empty-state"><h3>Failed to load store</h3><p>' + (e.message || e) + '</p></div>';
    }

    try {
        const url = await GetRegistryURL();
        document.getElementById('store-url-input').value = url;
    } catch (e) {}
}

function filterStore() {
    const query = document.getElementById('store-search').value.toLowerCase();
    const filtered = storeData.filter(p =>
        p.name.toLowerCase().includes(query) ||
        (p.desc && p.desc.toLowerCase().includes(query))
    );
    renderStore(filtered);
}

function renderStore(packages) {
    const list = document.getElementById('store-list');
    if (!packages || packages.length === 0) {
        list.innerHTML = `
        <div class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5"><path d="M6 2L3 6v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2V6l-3-4z"/><line x1="3" y1="6" x2="21" y2="6"/><path d="M16 10a4 4 0 0 1-8 0"/></svg>
            <h3>No packages found</h3>
            <p>Check your registry URL or try a different search.</p>
        </div>`;
        return;
    }
    list.innerHTML = packages.map(pkg => `
        <div class="installed-item">
            <div class="app-info">
                <span class="name">${pkg.name}${pkg.version ? ' v' + pkg.version : ''}</span>
                <span class="meta">${pkg.desc || ''}${pkg.homepage ? ' · ' + pkg.homepage : ''}</span>
            </div>
            <div class="app-actions">
                <button class="btn btn-primary btn-sm" onclick="window.storeInstall('${pkg.ici ? pkg.ici.replace(/'/g, "\\'") : ''}')">
                    Install
                </button>
            </div>
        </div>
    `).join('');
}

window.storeInstall = async function(iciURL) {
    if (!iciURL) return;
    try {
        const ici = await LoadRegistryICI(iciURL);
        if (ici) {
            pendingSource = iciURL;
            if (wizardMode) {
                showWizardConfirm(ici, iciURL);
            } else {
                await EnterWizardMode();
                wizardMode = true;
                document.querySelector('#app').innerHTML = wizardLayout();
                setupWizard();
                showWizardConfirm(ici, iciURL);
            }
        }
    } catch (e) {
        showToast('Failed to load package: ' + e, 'error');
    }
};

function showInstallScreen(ici) {
    currentICI = ici;
    document.getElementById('install-name').textContent = ici.name;
    document.getElementById('install-version').textContent = ici.version || '';
    document.getElementById('install-desc').textContent = ici.desc || '';
    document.getElementById('install-url').textContent = ici.url;
    document.getElementById('install-type').textContent = (ici.type || 'zip').toUpperCase();

    const shortcutsSection = document.getElementById('install-shortcuts-section');
    const shortcutsList = document.getElementById('install-shortcuts-list');
    if (ici.shortcuts && ici.shortcuts.length > 0) {
        shortcutsSection.style.display = 'block';
        shortcutsList.innerHTML = ici.shortcuts.map(s => {
            const name = s.name || s.Name || s;
            return `<div style="padding:4px 0">  ${name}</div>`;
        }).join('');
    } else if (ici.shortcut) {
        shortcutsSection.style.display = 'block';
        shortcutsList.innerHTML = `<div style="padding:4px 0">  ${ici.shortcut}</div>`;
    } else {
        shortcutsSection.style.display = 'none';
        shortcutsList.innerHTML = '';
    }

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

    const iciContent = buildICIContent(currentICI);

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
        const customInput = document.getElementById('custom-dir-input') || document.getElementById('wiz-custom-dir-input');
        lines.push(`custom_dir: ${customInput?.value || ''}`);
    }
    if (ici.shortcuts && ici.shortcuts.length > 0) {
        const entries = ici.shortcuts.map(s => {
            const exe = s.exe || s.Exe || s;
            const name = s.name || s.Name || s;
            return exe === name ? exe : `${exe}=${name}`;
        });
        lines.push(`shortcuts: ${entries.join(', ')}`);
    } else if (ici.shortcut) {
        lines.push(`shortcut: ${ici.shortcut}`);
    }
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
    if (eventsBound) return;
    eventsBound = true;

    EventsOn('ici-loaded', async (ici) => {
        if (wizardMode) {
            showWizardConfirm(ici, pendingSource);
        } else {
            showInstallScreen(ici);
        }
        try {
            const auto = await IsAutoInstall();
            if (auto) {
                setTimeout(() => {
                    if (wizardMode) doWizardInstall();
                    else doInstall();
                }, 500);
            }
        } catch (e) {}
    });

    EventsOn('protocol-ici', async (iciURL) => {
        try {
            const ici = await LoadRegistryICI(iciURL);
            if (ici) {
                if (wizardMode) {
                    showWizardConfirm(ici, iciURL);
                } else {
                    wizardMode = true;
                    document.querySelector('#app').innerHTML = wizardLayout();
                    setupWizard();
                    showWizardConfirm(ici, iciURL);
                }
            }
        } catch (e) {
            showToast('Failed to load package: ' + e, 'error');
        }
    });

    EventsOn('download-progress', (data) => {
        const percent = Math.round(data.percent || 0);
        const bar = document.getElementById('wiz-progress-bar') || document.getElementById('progress-bar');
        const text = document.getElementById('wiz-progress-text') || document.getElementById('progress-text');
        const pct = document.getElementById('wiz-progress-percent') || document.getElementById('progress-percent');
        if (bar) bar.style.width = percent + '%';
        if (text) text.textContent = `Downloading... ${formatBytes(data.downloaded)} / ${formatBytes(data.total)}`;
        if (pct) pct.textContent = percent + '%';
    });

    EventsOn('install-progress', (data) => {
        const statusMessages = {
            downloading: 'Downloading...',
            extracting: 'Extracting files...',
            shortcut: 'Creating shortcuts...',
            saving: 'Saving installation record...',
        };
        const msg = data.message || statusMessages[data.status] || data.status;
        const text = document.getElementById('wiz-progress-text') || document.getElementById('progress-text');
        if (text) text.textContent = msg;
        if (data.status === 'shortcut-warning') {
            showToast(msg, 'warning');
            if (wizardMode) {
                const warnDiv = document.getElementById('wiz-progress-warnings');
                if (warnDiv) warnDiv.innerHTML += `<div style="color:var(--warning);font-size:12px;margin-top:4px">  ${msg}</div>`;
            }
        }
    });

    EventsOn('extract-progress', (data) => {
        const text = document.getElementById('wiz-progress-text') || document.getElementById('progress-text');
        if (text) text.textContent = `Extracting: ${data.file} (${data.current}/${data.total})`;
    });

    EventsOn('install-complete', (data) => {
        if (wizardMode) {
            document.getElementById('wiz-done-name').textContent = data.name;
            document.getElementById('wiz-done-path').textContent = data.path;
            wizardShowStep('done');
        } else {
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
        }
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

async function checkPendingFile() {
    try {
        const pending = await GetPendingFile();
        if (pending && pending.path) {
            const ici = await LoadICIFile(pending.path);
            if (ici) {
                if (pending.auto) {
                    setTimeout(() => {
                        if (wizardMode) doWizardInstall();
                        else doInstall();
                    }, 500);
                }
            } else {
                showToast('Failed to load .ici file: no data returned', 'error');
            }
        }
    } catch (e) {
        showToast('Failed to load .ici file: ' + (e.message || e), 'error');
    }
}

init();
