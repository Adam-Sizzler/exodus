/* =================================================================
    КОНФИГУРАЦИЯ ПРИЛОЖЕНИЙ И ПЛАТФОРМ
================================================================= */

const platformData = {
    windows: {
        apps: [
            { 
                name: "V2RayN", 
                url: "https://github.com/2dust/v2rayN/releases/latest",
                client: "xray"
            },
            { 
                name: "FlClash", 
                url: "https://github.com/chen08209/FlClash/releases/latest", 
                scheme: "flclash://install-config?url=",
                client: "mihomo"
            },
            { 
                name: "Sing-box", 
                url: "https://github.com/SagerNet/sing-box/releases/download/v1.11.13/sing-box-1.11.13-windows-amd64.zip", 
                scheme: "sing-box://import-remote-profile?url=", 
                client: "singbox"
            },
            {
                name: "Happ",
                url: "https://github.com/Happ-proxy/happ-desktop/releases",
                scheme: "happ://add/",
                client: "xray"
            }
        ]
    },
    android: {
        apps: [
            { 
                name: "V2RayNG", 
                url: "https://github.com/2dust/v2rayNG/releases/latest", 
                scheme: "v2rayng://install-sub/?url=", 
                needsEncode: true,
                client: "xray"
            },
            { 
                name: "FlClash", 
                url: "https://github.com/chen08209/FlClash/releases/latest", 
                scheme: "flclash://install-config?url=",
                client: "mihomo"
            },
            { 
                name: "Sing-box", 
                url: "https://github.com/SagerNet/sing-box/releases/download/v1.11.13/SFA-1.11.13-universal.apk", 
                scheme: "sing-box://import-remote-profile?url=", 
                client: "singbox"
            },
            {
                name: "Happ",
                url: "https://play.google.com/store/apps/details?id=com.happproxy&hl=en",
                scheme: "happ://add/",
                client: "xray"
            }
        ]
    },
    ios: {
        apps: [
            { 
                name: "V2Box", 
                url: "https://apps.apple.com/us/app/v2box-v2ray-client/id6446814690", 
                scheme: "v2box://install-sub?url=", 
                useHash: true,
                client: "xray"
            },
            { 
                name: "Streisand", 
                url: "https://apps.apple.com/us/app/streisand/id6450534064", 
                scheme: "streisand://import/", 
                client: "xray"
            },
            { 
                name: "Sing-box", 
                url: "https://apps.apple.com/us/app/sing-box-vt/id6451272673", 
                scheme: "sing-box://import-remote-profile?url=", 
                client: "singbox"
            },
            {
                name: "Happ",
                url: "https://apps.apple.com/us/app/happ-proxy-utility/id6504287215",
                scheme: "happ://add/",
                client: "xray"
            }
        ]
    }
};

// Глобальные переменные состояния
let currentApps = [];
let activeAppIndex = 0;

document.addEventListener('DOMContentLoaded', () => {
    setTimeout(() => {
        document.documentElement.classList.remove('is-loading');
    }, 500);
    // 1. Инициализация данных пользователя
    const valState = window.validationState || { base: false, advanced: false };
    const userData = window.userData || {};

    // Элементы DOM
    const platformSelect = document.getElementById("platformSelect");
    const customSelect = document.getElementById('customPlatformSelect');
    const nameDisplay = document.getElementById("nameDisplay");
    const displayTitle = document.getElementById("displayTitle");
    const trafficUsed = document.getElementById("trafficUsed");
    const expireDateStat = document.getElementById("expireDateStat");
    
    // Получение параметров URL
    const params = new URLSearchParams(window.location.search);
    const urlName = params.get('user') || 'user';

    // Заполнение статистики и заголовков
    if (nameDisplay) nameDisplay.textContent = `${urlName}`;
    if (displayTitle && userData.profileTitle) displayTitle.textContent = userData.profileTitle;
    
    if (trafficUsed && userData.total) {
        const used = (userData.upload + userData.download);
        trafficUsed.textContent = `${formatBytes(used)} / ${formatBytes(userData.total)}`;
    }
    
    if (expireDateStat && userData.expire) {
        const date = new Date(userData.expire * 1000);
        expireDateStat.textContent = date.toLocaleDateString('ru-RU');
    }

    // 2. Логика кастомного селектора (Выбор ОС)
    if (customSelect) {
        const trigger = customSelect.querySelector('.select-trigger');
        const options = customSelect.querySelectorAll('.option');

        trigger.addEventListener('click', () => customSelect.classList.toggle('active'));

        options.forEach(opt => {
            opt.addEventListener('click', () => {
                const val = opt.getAttribute('data-value');
                updateSelectedPlatformUI(val);
                
                if (platformSelect) platformSelect.value = val;
                localStorage.setItem('selectedPlatform', val);
                updatePlatformUI();
                customSelect.classList.remove('active');
            });
        });

        document.addEventListener('click', (e) => {
            if (!customSelect.contains(e.target)) customSelect.classList.remove('active');
        });

        const savedPlatform = localStorage.getItem('selectedPlatform') || 'android';
        updateSelectedPlatformUI(savedPlatform);
        if (platformSelect) platformSelect.value = savedPlatform;
    }

    // 3. Обработчики глобальных кнопок (Support, Copy, QR)
    const supportBtn = document.getElementById("supportBtn");
    if (supportBtn) {
        supportBtn.addEventListener('click', () => {
            if (userData.supportUrl) window.open(userData.supportUrl, '_blank', 'noopener,noreferrer');
        });
    }

    const copyLinkBtn = document.getElementById("copyLinkBtn");
    if(copyLinkBtn) {
        copyLinkBtn.addEventListener('click', () => {
            copyToClipboard(buildSubscriptionUrl(urlName, 'advanced', 'xray'));
        });
    }

    const showQrBtn = document.getElementById("showQrBtn");
    if(showQrBtn) {
        showQrBtn.addEventListener('click', () => {
            const initialMode = valState.base ? 'base' : 'advanced';
            openQrModal(urlName, initialMode, valState);
        });
    }

    // Первичная отрисовка интерфейса
    updatePlatformUI();
});

/* =================================================================
    ЛОГИКА ИНТЕРФЕЙСА (Табы и Кнопки)
    =================================================================
*/

function updateSelectedPlatformUI(platform) {
    const valDisplay = document.getElementById('selectedValue');
    const iconContainer = document.getElementById('selectedIcon');

    if (valDisplay) {
        // Ищем в списке опций ту, у которой data-value совпадает с выбранной платформой
        const option = document.querySelector(`.option[data-value="${platform}"]`);
        if (option) {
            // Берем текст прямо из HTML (там у вас написано "iOS", "Android" и т.д.)
            valDisplay.textContent = option.textContent;
        }
    }

    // Безопасная проверка: существует ли объект и есть ли в нем ключ
    if (iconContainer) {
        if (typeof platformIcons !== 'undefined' && platformIcons[platform]) {
            iconContainer.innerHTML = platformIcons[platform];
        } else {
            // Если иконок нет, можно либо очистить контейнер, либо оставить пустым
            iconContainer.innerHTML = ''; 
        }
    }
}

function updatePlatformUI() {
    const select = document.getElementById("platformSelect");
    const platform = select ? select.value : 'android';
    const clientTabs = document.getElementById("clientTabs");
    const userData = window.userData || {};
    const availableClients = userData.availableClients ?? [];
    if (!clientTabs) return;
    
    // Получаем приложения для выбранной платформы или пустой массив
    let allApps = platformData[platform] ? platformData[platform].apps : [];
    currentApps = allApps.filter(app => availableClients.includes(app.client));
    activeAppIndex = 0; 
    
    clientTabs.innerHTML = "";
    
    if (currentApps.length > 0) {
        currentApps.forEach((app, index) => {
            const btn = document.createElement("button");
            btn.className = `tab-btn ${index === 0 ? 'active' : ''}`;

            let iconSvg = '';
            if (typeof APP_ICONS !== 'undefined' && APP_ICONS[app.name]) {
                iconSvg = APP_ICONS[app.name];
            }

            btn.innerHTML = `
                <span class="featured-dot"></span>
                <span class="tab-name">${app.name}</span>
                <span class="tab-bg-icon">${iconSvg}</span>
            `;
            
            btn.onclick = () => switchTab(index);
            clientTabs.appendChild(btn);
        });
        clientTabs.style.display = "grid";
    } else {
        clientTabs.style.display = "none";
    }
    
    renderStepActions();
}

function switchTab(index) {
    activeAppIndex = index;
    const tabs = document.querySelectorAll('.tab-btn');
    tabs.forEach((t, i) => {
        t.classList.toggle('active', i === index);
    });
    renderStepActions();
}

function renderStepActions() {
    const downloadArea = document.getElementById("downloadButtonsArea");
    const importArea = document.getElementById("importButtonsArea");
    const valState = window.validationState || { base: false, advanced: false };
    const nameDisplay = document.getElementById("nameDisplay");
    if (!downloadArea || !importArea) return;
    
    // Очистка областей
    downloadArea.innerHTML = "";
    importArea.innerHTML = "";
    
    if (currentApps.length === 0) {
        downloadArea.innerHTML = "<p>Для этой платформы пока нет инструкций.</p>";
        return;
    }
    
    const app = currentApps[activeAppIndex];
    const rawUser = nameDisplay ? nameDisplay.textContent.replace('@', '') : 'user'; 

    // 1. Рендер кнопки скачивания
    downloadArea.innerHTML = `
        <a href="${app.url}" target="_blank" rel="noopener noreferrer" class="btn btn-secondary">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 8px; flex-shrink: 0;">
                <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                <path d="M12 6h-6a2 2 0 0 0 -2 2v10a2 2 0 0 0 2 2h10a2 2 0 0 0 2 -2v-6"></path>
                <path d="M11 13l9 -9"></path>
                <path d="M15 4h5v5"></path>
            </svg>
            Скачать ${app.name}
        </a>`;

    // 2. Рендер кнопок импорта (подписки)
    if (!valState.base && !valState.advanced) {
        importArea.innerHTML = `<div style="color:#ef4444; padding: 10px;">Нет доступных подписок</div>`;
        return;
    }

    if (app.client === 'xray') {
        // Для Xray показываем выбор, если доступны оба режима
        if (valState.base && valState.advanced) {
            importArea.appendChild(createImportButton(`Добавить (Базовая)`, () => handleImport(rawUser, app, 'base'), true));
            importArea.appendChild(createImportButton(`Добавить (Расширенная)`, () => handleImport(rawUser, app, 'advanced')));
        } else {
            const mode = valState.base ? 'base' : 'advanced';
            importArea.appendChild(createImportButton(`Добавить подписку`, () => handleImport(rawUser, app, mode)));
        }
    } else {
        // Для НЕ xray (Sing-box, Mihomo) — только одна кнопка, режим всегда advanced (но в URL не пишем)
        importArea.appendChild(createImportButton(`Добавить подписку`, () => handleImport(rawUser, app, 'advanced')));
    }
}

/* =================================================================
    ВСПОМОГАТЕЛЬНЫЕ ФУНКЦИИ (Импорт, QR, Utils)
    =================================================================
*/

function createImportButton(text, onClick, outline = false) {
    const btn = document.createElement('button');
    btn.className = outline ? "btn btn-primary btn-outline" : "btn btn-primary";
    btn.style.marginBottom = "10px"; 
    
    btn.innerHTML = `
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" style="margin-right: 4px;">
            <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
            <path d="M12 5v14"></path>
            <path d="M5 12h14"></path>
        </svg>
        ${text}
    `;
    
    btn.onclick = onClick;
    return btn;
}

function buildSubscriptionUrl(user, mode, client) {
    const url = new URL(window.location.origin + window.location.pathname);
    url.searchParams.set('user', user);
    
    if (client === 'xray' && mode === 'base') {
        url.searchParams.set('mode', 'base');
    } else {
        url.searchParams.delete('mode');
    }
    
    return url.toString();
}

function handleImport(userName, app, mode) {
    const SUB_NAME = "tlw";
    const subUrl = buildSubscriptionUrl(userName, mode, app.client);

    if (!app.scheme) {
        copyToClipboard(subUrl);
        return;
    }

    let deepLink = "";
    const encodedUrl = encodeURIComponent(subUrl);

    if (app.name === "Sing-box") {
        deepLink = `${app.scheme}${subUrl}`;
    } else if (app.name === "V2Box") {
        // V2Box ожидает urlencoded и name параметр
        deepLink = `${app.scheme}${encodedUrl}&name=${encodeURIComponent(SUB_NAME)}`;
    } else if (app.name === "Happ") {
        deepLink = `${app.scheme}${subUrl}`;
    } else if (app.needsEncode) {
        deepLink = `${app.scheme}${encodeURIComponent(subUrl + "#" + SUB_NAME)}`;
    } else if (app.useHash) {
        deepLink = `${app.scheme}${subUrl}#${encodeURIComponent(SUB_NAME)}`;
    } else {
        deepLink = `${app.scheme}${encodedUrl}`;
    }

    openDeepLink(deepLink, subUrl, app.name);
}

function openDeepLink(url, fallbackUrl, appName) {
    const fallbackTimer = setTimeout(() => {
        copyToClipboard(fallbackUrl);
        showNotification("Не удалось открыть приложение. Ссылка скопирована", "warning");
    }, 1400);

    const onVisibilityChange = () => {
        if (document.visibilityState === 'hidden') {
            clearTimeout(fallbackTimer);
            document.removeEventListener('visibilitychange', onVisibilityChange);
        }
    };
    document.addEventListener('visibilitychange', onVisibilityChange);

    const a = document.createElement('a');
    a.href = url;
    a.rel = 'noopener noreferrer';
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
    showNotification(`Открываю ${appName}...`, 'success');
}


// === QR Code Logic ===

function openQrModal(user, initialMode, valState) {
    const modal = document.getElementById('qrModal');
    const qrContainer = document.getElementById("qrcode");
    if (!modal || !qrContainer) return;
    
    const currentApp = currentApps[activeAppIndex];
    const clientType = currentApp ? currentApp.client : 'xray';

    // Функция рендера содержимого модалки
    const renderContent = (mode) => {
        const url = buildSubscriptionUrl(user, mode, clientType);
        
        qrContainer.innerHTML = `
            <div class="qr-main-card">
                <div class="qr-header">
                    <span class="qr-title">QR код подписки</span>
                    <button class="close-icon-btn" onclick="document.getElementById('qrModal').style.display='none'">
                        <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line></svg>
                    </button>
                </div>

                ${(valState.base) ? `
                <div class="qr-modes-container">
                    <div class="qr-modes">
                        <button class="mode-btn ${mode === 'base' ? 'active' : ''}" id="btnModeBase">Базовая</button>
                        <button class="mode-btn ${mode === 'advanced' ? 'active' : ''}" id="btnModeAdv">Расширенная</button>
                    </div>
                </div>
                ` : ''}

                <div class="qr-display-area">
                    <canvas id="qrCanvas"></canvas>
                </div>

                <div style="text-align: center;">
                    <div class="qr-instruction">Отсканируйте код в приложении</div>
                    <div class="qr-sub-instruction">
                        Наведите камеру в приложении VPN или скопируйте ссылку вручную.
                    </div>
                </div>

                <button class="btn btn-full" id="copyLinkActionBtn">
                    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path></svg>
                    Скопировать ссылку
                </button>
            </div>
        `;

        // Слушатели переключения режимов
        if (valState.base && valState.advanced) {
            const btnBase = document.getElementById('btnModeBase');
            const btnAdv = document.getElementById('btnModeAdv');
            if (btnBase) btnBase.onclick = () => renderContent('base');
            if (btnAdv) btnAdv.onclick = () => renderContent('advanced');
        }

        // Кнопка копирования внутри модалки
        const cpBtn = document.getElementById('copyLinkActionBtn');
        if (cpBtn) cpBtn.onclick = () => copyToClipboard(url);

        // Генерация Canvas
        if (typeof QRCode !== 'undefined' && QRCode.toCanvas) {
            const canvas = document.getElementById('qrCanvas');
            QRCode.toCanvas(canvas, url, {
                width: 380,
                margin: 0,
                color: { dark: '#22d3ee', light: '#00000000' }
            }, function (error) {
                if (error) console.error(error);
            });
        }
    };

    renderContent(initialMode);
    modal.style.display = 'flex';
}

// === Utilities ===

function copyToClipboard(text) {
    if (navigator.clipboard && navigator.clipboard.writeText) {
        navigator.clipboard.writeText(text).then(() => {
            showNotification("Ссылка скопирована!", "success");
        }).catch(err => {
            console.error(err);
            showNotification("Ошибка копирования", "error");
        });
    } else {
        // Fallback
        const textArea = document.createElement("textarea");
        textArea.value = text;
        document.body.appendChild(textArea);
        textArea.select();
        try {
            document.execCommand('copy');
            showNotification("Ссылка скопирована!", "success");
        } catch (err) {
            showNotification("Ошибка копирования", "error");
        }
        document.body.removeChild(textArea);
    }
}

function showNotification(msg, type) {
    const area = document.getElementById('notificationArea');
    if (!area) return; 
    
    const note = document.createElement('div');
    note.className = 'notification';
    note.style.borderColor = type === 'error' ? '#ef4444' : '#1d9bf0';
    note.textContent = msg;
    
    area.appendChild(note);
    setTimeout(() => {
        note.style.opacity = '0';
        setTimeout(() => note.remove(), 300);
    }, 3000);
}

function formatBytes(bytes, decimals = 2) {
    if (!+bytes) return '0 B';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`;
}
