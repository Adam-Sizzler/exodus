import React, { useEffect, useState } from 'react';
import { subscriptionSettingsApi } from '../api';

const EMPTY_FORM = {
  profile_title: '',
  support_link: '',
  profile_update_interval: 12,
  address: '',
  port: 9263,
  api_schema: 'grpc',
  api_path: '',
  happ_announce: '',
  happ_routing: '',
  is_profile_webpage_url_enabled: true,
  serve_json_at_base_subscription: false,
  is_show_custom_remarks: true,
  randomize_hosts: false,
  response_rules: '{}',
};

const DEFAULT_HWID_SETTINGS = {
  enabled: false,
  fallbackDeviceLimit: 999,
  maxDevicesAnnounce: '',
};

const SUBSCRIPTION_API_SCHEMAS = ['grpc', 'https', 'http'];

const REMARK_GROUPS = [
  {
    key: 'HWIDMaxDevicesExceeded',
    title: 'HWID: Макс. число устройств превышено',
    tone: 'red',
    defaults: ['Limit of devices reached'],
  },
  {
    key: 'expiredUsers',
    title: 'Статус пользователя: EXPIRED',
    tone: 'red',
    defaults: ['⌛ Subscription expired', 'Contact support'],
  },
  {
    key: 'disabledUsers',
    title: 'Статус пользователя: DISABLED',
    tone: 'gray',
    defaults: ['🚫 Subscription disabled', 'Contact support'],
  },
  {
    key: 'HWIDNotSupported',
    title: 'HWID: Не поддерживается',
    tone: 'red',
    defaults: ['App not supported'],
  },
  {
    key: 'limitedUsers',
    title: 'Статус пользователя: LIMITED',
    tone: 'orange',
    defaults: ['🚧 Subscription limited', 'Contact support'],
  },
  {
    key: 'emptyHosts',
    title: 'Отсутствуют хосты',
    tone: 'blue',
    defaults: ['→ Remnawave', '→ No hosts found', '→ Check Hosts tab', '→ Check Internal Squads tab'],
  },
];

let localRowId = 0;

const nextRowId = () => {
  localRowId += 1;
  return `row-${Date.now()}-${localRowId}`;
};

const parseJsonObject = (value, fallback = {}) => {
  if (!value || typeof value !== 'string') {
    return fallback;
  }

  try {
    const parsed = JSON.parse(value);
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) {
      return fallback;
    }
    return parsed;
  } catch {
    return fallback;
  }
};

const parseHwidSettings = (value) => {
  const parsed = parseJsonObject(value);
  const fallbackDeviceLimit = Number(parsed.fallbackDeviceLimit);

  return {
    enabled: !!parsed.enabled,
    fallbackDeviceLimit:
      Number.isFinite(fallbackDeviceLimit) && fallbackDeviceLimit >= 0
        ? fallbackDeviceLimit
        : DEFAULT_HWID_SETTINGS.fallbackDeviceLimit,
    maxDevicesAnnounce:
      parsed.maxDevicesAnnounce === null || parsed.maxDevicesAnnounce === undefined
        ? ''
        : String(parsed.maxDevicesAnnounce),
  };
};

const parseCustomRemarks = (value) => {
  const parsed = parseJsonObject(value);
  const result = {};

  REMARK_GROUPS.forEach((group) => {
    if (Object.prototype.hasOwnProperty.call(parsed, group.key) && Array.isArray(parsed[group.key])) {
      result[group.key] = parsed[group.key].map((item) => String(item ?? ''));
      return;
    }

    result[group.key] = [...group.defaults];
  });

  return result;
};

const parseCustomHeaders = (value) => {
  const parsed = parseJsonObject(value);
  return Object.entries(parsed).map(([key, rawValue]) => ({
    id: nextRowId(),
    key: String(key),
    value:
      typeof rawValue === 'string'
        ? rawValue
        : rawValue === null || rawValue === undefined
          ? ''
          : JSON.stringify(rawValue),
  }));
};

const buildCustomRemarksJson = (remarks) => {
  const payload = {};

  REMARK_GROUPS.forEach((group) => {
    const list = Array.isArray(remarks[group.key]) ? remarks[group.key] : [];
    payload[group.key] = list.map((entry) => String(entry ?? ''));
  });

  return JSON.stringify(payload);
};

const buildCustomHeadersJson = (rows) => {
  const payload = {};

  rows.forEach((row) => {
    const key = String(row.key ?? '').trim();
    if (!key) {
      return;
    }
    payload[key] = String(row.value ?? '');
  });

  return JSON.stringify(payload);
};

const buildHwidSettingsJson = (hwidSettings) => {
  const limit = Number(hwidSettings.fallbackDeviceLimit);
  return JSON.stringify({
    enabled: !!hwidSettings.enabled,
    fallbackDeviceLimit:
      Number.isFinite(limit) && limit >= 0 ? limit : DEFAULT_HWID_SETTINGS.fallbackDeviceLimit,
    maxDevicesAnnounce: hwidSettings.maxDevicesAnnounce.trim()
      ? hwidSettings.maxDevicesAnnounce.trim()
      : null,
  });
};

function SubscriptionSettings() {
  const [activeTab, setActiveTab] = useState('general');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState(null);
  const [saveSuccess, setSaveSuccess] = useState('');

  const [settings, setSettings] = useState(null);
  const [formData, setFormData] = useState(EMPTY_FORM);
  const [hwidSettings, setHwidSettings] = useState(DEFAULT_HWID_SETTINGS);
  const [remarks, setRemarks] = useState(() => parseCustomRemarks('{}'));
  const [customHeaders, setCustomHeaders] = useState([]);

  const loadSettings = async () => {
    try {
      setLoading(true);
      const data = await subscriptionSettingsApi.get();
      const next = data.settings;

      setSettings(next);
      setFormData({
        profile_title: next.profile_title ?? '',
        support_link: next.support_link ?? '',
        profile_update_interval: Number(next.profile_update_interval ?? 12),
        address: next.address ?? '',
        port: Number(next.port ?? 9263),
        api_schema: String(next.api_schema ?? 'grpc').toLowerCase(),
        api_path: next.api_path ?? '',
        happ_announce: next.happ_announce ?? '',
        happ_routing: next.happ_routing ?? '',
        is_profile_webpage_url_enabled: !!next.is_profile_webpage_url_enabled,
        serve_json_at_base_subscription: !!next.serve_json_at_base_subscription,
        is_show_custom_remarks: !!next.is_show_custom_remarks,
        randomize_hosts: !!next.randomize_hosts,
        response_rules: next.response_rules ?? '{}',
      });

      setHwidSettings(parseHwidSettings(next.hwid_settings ?? '{}'));
      setRemarks(parseCustomRemarks(next.custom_remarks ?? '{}'));
      setCustomHeaders(parseCustomHeaders(next.custom_response_headers ?? '{}'));

      setError(null);
    } catch (err) {
      setError(err.message || 'Не удалось загрузить настройки подписки');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadSettings();
  }, []);

  const setField = (field, value) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const addHeaderRow = () => {
    setCustomHeaders((prev) => [...prev, { id: nextRowId(), key: '', value: '' }]);
  };

  const updateHeaderRow = (id, field, value) => {
    setCustomHeaders((prev) =>
      prev.map((row) => (row.id === id ? { ...row, [field]: value } : row))
    );
  };

  const removeHeaderRow = (id) => {
    setCustomHeaders((prev) => prev.filter((row) => row.id !== id));
  };

  const addRemarkRow = (groupKey) => {
    setRemarks((prev) => ({
      ...prev,
      [groupKey]: [...(prev[groupKey] ?? []), ''],
    }));
  };

  const updateRemarkRow = (groupKey, index, value) => {
    setRemarks((prev) => {
      const current = [...(prev[groupKey] ?? [])];
      current[index] = value;
      return { ...prev, [groupKey]: current };
    });
  };

  const removeRemarkRow = (groupKey, index) => {
    setRemarks((prev) => {
      const current = [...(prev[groupKey] ?? [])];
      current.splice(index, 1);
      return { ...prev, [groupKey]: current };
    });
  };

  const handleSave = async () => {
    if (!settings?.uuid) {
      return;
    }

    if (!formData.profile_title.trim()) {
      alert('Поле "Заголовок профиля" обязательно');
      return;
    }

    if (!formData.support_link.trim()) {
      alert('Поле "Ссылка на поддержку" обязательно');
      return;
    }

    if (!Number.isFinite(Number(formData.profile_update_interval)) || Number(formData.profile_update_interval) < 1) {
      alert('Интервал авто-обновления должен быть не меньше 1 часа');
      return;
    }
    if (!formData.address.trim()) {
      alert('Поле "Address" обязательно');
      return;
    }
    if (!Number.isFinite(Number(formData.port)) || Number(formData.port) < 1 || Number(formData.port) > 65535) {
      alert('Поле "Port" должно быть в диапазоне 1-65535');
      return;
    }
    if (!SUBSCRIPTION_API_SCHEMAS.includes(String(formData.api_schema ?? '').toLowerCase())) {
      alert('Поле "API Schema" должно быть одним из: grpc, https, http');
      return;
    }

    const usedHeaderKeys = new Set();
    for (const row of customHeaders) {
      const key = row.key.trim();
      if (!key) {
        continue;
      }
      const normalized = key.toLowerCase();
      if (usedHeaderKeys.has(normalized)) {
        alert(`Дублирующийся хэдер: ${key}`);
        return;
      }
      usedHeaderKeys.add(normalized);
    }

    const payload = {
      profile_title: formData.profile_title.trim(),
      support_link: formData.support_link.trim(),
      profile_update_interval: Number(formData.profile_update_interval),
      address: formData.address.trim(),
      port: Number(formData.port),
      api_schema: String(formData.api_schema).toLowerCase(),
      api_path: formData.api_path.trim(),
      happ_announce: formData.happ_announce,
      happ_routing: formData.happ_routing,
      is_profile_webpage_url_enabled: !!formData.is_profile_webpage_url_enabled,
      serve_json_at_base_subscription: !!formData.serve_json_at_base_subscription,
      is_show_custom_remarks: !!formData.is_show_custom_remarks,
      randomize_hosts: !!formData.randomize_hosts,
      response_rules: String(formData.response_rules ?? '{}'),
      custom_response_headers: buildCustomHeadersJson(customHeaders),
      hwid_settings: buildHwidSettingsJson(hwidSettings),
      custom_remarks: buildCustomRemarksJson(remarks),
    };

    try {
      setSaving(true);
      const data = await subscriptionSettingsApi.update(settings.uuid, payload);
      if (data.settings) {
        setSettings(data.settings);
      }
      await loadSettings();
      setSaveSuccess('Настройки сохранены');
    } catch (err) {
      alert(`Не удалось сохранить настройки: ${err.message}`);
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return <div className="loading">Загрузка...</div>;
  }

  return (
    <div className="subscription-settings-shell">
      <div className="subscription-settings-top">
        <div className="subscription-settings-top-main">
          <span className="subscription-settings-top-icon" aria-hidden="true">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="3" />
              <path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 1 1 -2.83 2.83l-.06-.06a1.65 1.65 0 0 0 -1.82 -.33a1.65 1.65 0 0 0 -1 1.51V21a2 2 0 1 1 -4 0v-.09a1.65 1.65 0 0 0 -1 -1.51a1.65 1.65 0 0 0 -1.82 .33l-.06 .06a2 2 0 1 1 -2.83 -2.83l.06-.06a1.65 1.65 0 0 0 .33 -1.82a1.65 1.65 0 0 0 -1.51 -1H3a2 2 0 1 1 0 -4h.09a1.65 1.65 0 0 0 1.51 -1a1.65 1.65 0 0 0 -.33 -1.82l-.06 -.06a2 2 0 1 1 2.83 -2.83l.06 .06a1.65 1.65 0 0 0 1.82 .33H9a1.65 1.65 0 0 0 1 -1.51V3a2 2 0 1 1 4 0v.09a1.65 1.65 0 0 0 1 1.51a1.65 1.65 0 0 0 1.82 -.33l.06 -.06a2 2 0 1 1 2.83 2.83l-.06 .06a1.65 1.65 0 0 0 -.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 1 1 0 4h-.09a1.65 1.65 0 0 0 -1.51 1z" />
            </svg>
          </span>
          <div>
            <h2 className="subscription-settings-title">Настройки подписки</h2>
            <p className="subscription-settings-subtitle">Информация о подписке, примечания и дополнительные хэдеры</p>
          </div>
        </div>
        <div className="subscription-settings-actions">
          <button className="btn btn-secondary" onClick={loadSettings} disabled={saving} type="button">
            Обновить
          </button>
          <button className="btn btn-primary" onClick={handleSave} disabled={saving || !settings?.uuid} type="button">
            {saving ? 'Сохранение...' : 'Сохранить'}
          </button>
        </div>
      </div>

      <div className="subscription-tabs" role="tablist" aria-label="Subscription settings tabs">
        <button
          className={`subscription-tab ${activeTab === 'general' ? 'active' : ''}`}
          type="button"
          role="tab"
          aria-selected={activeTab === 'general'}
          onClick={() => setActiveTab('general')}
        >
          Информация о подписке
        </button>
        <button
          className={`subscription-tab ${activeTab === 'connection' ? 'active' : ''}`}
          type="button"
          role="tab"
          aria-selected={activeTab === 'connection'}
          onClick={() => setActiveTab('connection')}
        >
          Подключение
        </button>
        <button
          className={`subscription-tab ${activeTab === 'remarks' ? 'active' : ''}`}
          type="button"
          role="tab"
          aria-selected={activeTab === 'remarks'}
          onClick={() => setActiveTab('remarks')}
        >
          Кастомные примечания
        </button>
        <button
          className={`subscription-tab ${activeTab === 'headers' ? 'active' : ''}`}
          type="button"
          role="tab"
          aria-selected={activeTab === 'headers'}
          onClick={() => setActiveTab('headers')}
        >
          Доп. хэдеры
        </button>
      </div>

      {error ? <div className="alert alert-error">{error}</div> : null}
      {saveSuccess ? <div className="subscription-success">{saveSuccess}</div> : null}

      {!settings ? (
        <div className="empty-state">
          <p>Настройки подписки не найдены</p>
        </div>
      ) : null}

      {settings && activeTab === 'general' ? (
        <div className="subscription-grid">
          <section className="subscription-card">
            <header className="subscription-card-header">
              <h3 className="subscription-card-title">Информация о подписке</h3>
              <p className="subscription-card-desc">
                Эти настройки поддерживаются клиентами вроде Happ, v2rayNG, Streisand и другими.
              </p>
            </header>
            <div className="subscription-card-body">
              <div className="subscription-field">
                <label className="subscription-field-label" htmlFor="subscription-profile-title">Заголовок профиля</label>
                <input
                  id="subscription-profile-title"
                  className="subscription-field-input"
                  type="text"
                  placeholder="Введите заголовок профиля"
                  value={formData.profile_title}
                  onChange={(e) => setField('profile_title', e.target.value)}
                />
              </div>

              <div className="subscription-field">
                <label className="subscription-field-label" htmlFor="subscription-profile-update-interval">Интервал авто-обновления (часы)</label>
                <input
                  id="subscription-profile-update-interval"
                  className="subscription-field-input"
                  type="number"
                  min="1"
                  placeholder="12"
                  value={formData.profile_update_interval}
                  onChange={(e) => setField('profile_update_interval', e.target.value)}
                />
              </div>

              <div className="subscription-field">
                <label className="subscription-field-label" htmlFor="subscription-support-link">Ссылка на поддержку</label>
                <input
                  id="subscription-support-link"
                  className="subscription-field-input"
                  type="text"
                  placeholder="https://support.example.com"
                  value={formData.support_link}
                  onChange={(e) => setField('support_link', e.target.value)}
                />
              </div>
            </div>
          </section>

          <section className="subscription-card">
            <header className="subscription-card-header">
              <h3 className="subscription-card-title">Доп. опции</h3>
              <p className="subscription-card-desc">Настройка поведения базовой подписки.</p>
            </header>
            <div className="subscription-card-body">
              <label className="subscription-check">
                <input
                  type="checkbox"
                  checked={formData.serve_json_at_base_subscription}
                  onChange={(e) => setField('serve_json_at_base_subscription', e.target.checked)}
                />
                <div>
                  <span className="subscription-check-title">Использовать JSON в базовой подписке</span>
                  <span className="subscription-check-text">Если клиент поддерживает JSON, отдавать JSON вместо обычного формата.</span>
                </div>
              </label>

              <label className="subscription-check">
                <input
                  type="checkbox"
                  checked={formData.randomize_hosts}
                  onChange={(e) => setField('randomize_hosts', e.target.checked)}
                />
                <div>
                  <span className="subscription-check-title">Перемешивать хосты</span>
                  <span className="subscription-check-text">Выдавать хосты в случайном порядке в содержимом подписки.</span>
                </div>
              </label>

              <label className="subscription-check">
                <input
                  type="checkbox"
                  checked={formData.is_profile_webpage_url_enabled}
                  onChange={(e) => setField('is_profile_webpage_url_enabled', e.target.checked)}
                />
                <div>
                  <span className="subscription-check-title">URL страницы профиля</span>
                  <span className="subscription-check-text">Включить публикацию URL страницы профиля для поддерживаемых клиентов.</span>
                </div>
              </label>

              <label className="subscription-check">
                <input
                  type="checkbox"
                  checked={formData.is_show_custom_remarks}
                  onChange={(e) => setField('is_show_custom_remarks', e.target.checked)}
                />
                <div>
                  <span className="subscription-check-title">Отображать кастомные примечания</span>
                  <span className="subscription-check-text">Показывать специальные сообщения для EXPIRED/LIMITED/DISABLED.</span>
                </div>
              </label>
            </div>
          </section>

          <section className="subscription-card">
            <header className="subscription-card-header">
              <h3 className="subscription-card-title">Настройки HWID</h3>
              <p className="subscription-card-desc">Параметры лимитов устройств для подписки.</p>
            </header>
            <div className="subscription-card-body">
              <label className="subscription-check">
                <input
                  type="checkbox"
                  checked={hwidSettings.enabled}
                  onChange={(e) => setHwidSettings((prev) => ({ ...prev, enabled: e.target.checked }))}
                />
                <div>
                  <span className="subscription-check-title">Лимит HWID</span>
                  <span className="subscription-check-text">Включить ограничение числа устройств по HWID.</span>
                </div>
              </label>

              <div className="subscription-field">
                <label className="subscription-field-label" htmlFor="subscription-hwid-device-limit">Лимит устройств по умолчанию</label>
                <input
                  id="subscription-hwid-device-limit"
                  className="subscription-field-input"
                  type="number"
                  min="0"
                  placeholder="999"
                  value={hwidSettings.fallbackDeviceLimit}
                  onChange={(e) =>
                    setHwidSettings((prev) => ({
                      ...prev,
                      fallbackDeviceLimit: e.target.value,
                    }))
                  }
                />
              </div>

              <div className="subscription-field">
                <label className="subscription-field-label" htmlFor="subscription-hwid-announce">Объявление при лимите устройств</label>
                <textarea
                  id="subscription-hwid-announce"
                  className="subscription-field-input subscription-field-textarea"
                  rows={3}
                  maxLength={200}
                  placeholder="Макс. 200 символов"
                  value={hwidSettings.maxDevicesAnnounce}
                  onChange={(e) =>
                    setHwidSettings((prev) => ({
                      ...prev,
                      maxDevicesAnnounce: e.target.value,
                    }))
                  }
                />
              </div>
            </div>
          </section>

          <section className="subscription-card">
            <header className="subscription-card-header">
              <h3 className="subscription-card-title">Announce и Routing</h3>
              <p className="subscription-card-desc">Настройка announce-сообщения и Happ routing.</p>
            </header>
            <div className="subscription-card-body">
              <div className="subscription-field">
                <label className="subscription-field-label" htmlFor="subscription-happ-announce">Announce</label>
                <textarea
                  id="subscription-happ-announce"
                  className="subscription-field-input subscription-field-textarea"
                  rows={3}
                  placeholder="Введите announce-сообщение"
                  value={formData.happ_announce}
                  onChange={(e) => setField('happ_announce', e.target.value)}
                />
              </div>

              <div className="subscription-field">
                <label className="subscription-field-label" htmlFor="subscription-happ-routing">Happ routing</label>
                <textarea
                  id="subscription-happ-routing"
                  className="subscription-field-input subscription-field-textarea"
                  rows={4}
                  placeholder="happ://routing/add/..."
                  value={formData.happ_routing}
                  onChange={(e) => setField('happ_routing', e.target.value)}
                />
              </div>
            </div>
          </section>
        </div>
      ) : null}

      {settings && activeTab === 'connection' ? (
        <section className="subscription-card">
          <header className="subscription-card-header">
            <h3 className="subscription-card-title">Подключение к сервису подписки</h3>
            <p className="subscription-card-desc">Настройки адреса и транспорта для связи с subscription service.</p>
          </header>
          <div className="subscription-card-body">
            <div className="subscription-field">
              <label className="subscription-field-label" htmlFor="subscription-connection-address">Address</label>
              <input
                id="subscription-connection-address"
                className="subscription-field-input"
                type="text"
                placeholder="subscription.local"
                value={formData.address}
                onChange={(e) => setField('address', e.target.value)}
              />
            </div>

            <div className="subscription-inline-row subscription-inline-row-2cols">
              <div className="subscription-field">
                <label className="subscription-field-label" htmlFor="subscription-connection-port">Port</label>
                <input
                  id="subscription-connection-port"
                  className="subscription-field-input"
                  type="number"
                  min="1"
                  max="65535"
                  placeholder="9263"
                  value={formData.port}
                  onChange={(e) => setField('port', e.target.value)}
                />
              </div>

              <div className="subscription-field">
                <label className="subscription-field-label" htmlFor="subscription-connection-schema">API Schema</label>
                <select
                  id="subscription-connection-schema"
                  className="subscription-field-input"
                  value={formData.api_schema}
                  onChange={(e) => setField('api_schema', e.target.value)}
                >
                  {SUBSCRIPTION_API_SCHEMAS.map((schema) => (
                    <option key={schema} value={schema}>
                      {schema}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="subscription-field">
              <label className="subscription-field-label" htmlFor="subscription-connection-api-path">API Path</label>
              <input
                id="subscription-connection-api-path"
                className="subscription-field-input"
                type="text"
                placeholder="/"
                value={formData.api_path}
                onChange={(e) => setField('api_path', e.target.value)}
              />
            </div>
          </div>
        </section>
      ) : null}

      {settings && activeTab === 'remarks' ? (
        <section className="subscription-card">
          <header className="subscription-card-header">
            <h3 className="subscription-card-title">Кастомные примечания</h3>
            <p className="subscription-card-desc">
              Эти значения будут отображаться, когда пользователь ограничен, просрочен или отключен.
            </p>
          </header>
          <div className="subscription-card-body">
            <div className="subscription-remarks-grid">
              {REMARK_GROUPS.map((group) => {
                const entries = remarks[group.key] ?? [];

                return (
                  <div key={group.key} className={`subscription-remark-card tone-${group.tone}`}>
                    <div className="subscription-remark-card-header">
                      <h4>{group.title}</h4>
                      <button className="btn btn-secondary subscription-small-btn" type="button" onClick={() => addRemarkRow(group.key)}>
                        Добавить
                      </button>
                    </div>

                    <div className="subscription-remark-card-body">
                      {entries.length === 0 ? <div className="subscription-empty-row">Примечаний нет</div> : null}

                      {entries.map((entry, index) => (
                        <div className="subscription-inline-row" key={`${group.key}-${index}`}>
                          <input
                            className="subscription-field-input"
                            type="text"
                            placeholder="Введите примечание"
                            value={entry}
                            onChange={(e) => updateRemarkRow(group.key, index, e.target.value)}
                          />
                          <button
                            className="subscription-icon-btn danger"
                            type="button"
                            onClick={() => removeRemarkRow(group.key, index)}
                            aria-label="Удалить примечание"
                          >
                            <svg viewBox="0 0 256 256" fill="currentColor">
                              <path d="M216,48H176V40a24,24,0,0,0-24-24H104A24,24,0,0,0,80,40v8H40a8,8,0,0,0,0,16h8V208a16,16,0,0,0,16,16H192a16,16,0,0,0,16-16V64h8a8,8,0,0,0,0-16ZM96,40a8,8,0,0,1,8-8h48a8,8,0,0,1,8,8v8H96Zm96,168H64V64H192ZM112,104v64a8,8,0,0,1-16,0V104a8,8,0,0,1,16,0Zm48,0v64a8,8,0,0,1-16,0V104a8,8,0,0,1,16,0Z" />
                            </svg>
                          </button>
                        </div>
                      ))}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>
        </section>
      ) : null}

      {settings && activeTab === 'headers' ? (
        <section className="subscription-card">
          <header className="subscription-card-header">
            <h3 className="subscription-card-title">Доп. хэдеры</h3>
            <p className="subscription-card-desc">Хэдеры отправляются вместе с содержимым подписки.</p>
          </header>
          <div className="subscription-card-body">
            <div className="subscription-inline-actions">
              <button className="btn btn-secondary" type="button" onClick={addHeaderRow}>
                Добавить хэдер
              </button>
            </div>

            <div className="subscription-headers-list">
              {customHeaders.length === 0 ? (
                <div className="subscription-empty-row">Список дополнительных хэдеров пуст</div>
              ) : null}

              {customHeaders.map((row) => (
                <div className="subscription-header-row" key={row.id}>
                  <input
                    className="subscription-field-input"
                    type="text"
                    placeholder="Header-Name"
                    value={row.key}
                    onChange={(e) => updateHeaderRow(row.id, 'key', e.target.value)}
                  />
                  <input
                    className="subscription-field-input"
                    type="text"
                    placeholder="Header value"
                    value={row.value}
                    onChange={(e) => updateHeaderRow(row.id, 'value', e.target.value)}
                  />
                  <button
                    className="subscription-icon-btn danger"
                    type="button"
                    onClick={() => removeHeaderRow(row.id)}
                    aria-label="Удалить хэдер"
                  >
                    <svg viewBox="0 0 256 256" fill="currentColor">
                      <path d="M216,48H176V40a24,24,0,0,0-24-24H104A24,24,0,0,0,80,40v8H40a8,8,0,0,0,0,16h8V208a16,16,0,0,0,16,16H192a16,16,0,0,0,16-16V64h8a8,8,0,0,0,0-16ZM96,40a8,8,0,0,1,8-8h48a8,8,0,0,1,8,8v8H96Zm96,168H64V64H192ZM112,104v64a8,8,0,0,1-16,0V104a8,8,0,0,1,16,0Zm48,0v64a8,8,0,0,1-16,0V104a8,8,0,0,1,16,0Z" />
                    </svg>
                  </button>
                </div>
              ))}
            </div>
          </div>
        </section>
      ) : null}

      {settings ? (
        <div className="subscription-meta-row">
          <span>UUID: {settings.uuid}</span>
          <span>Updated: {settings.updated_at ? new Date(settings.updated_at).toLocaleString() : '-'}</span>
        </div>
      ) : null}
    </div>
  );
}

export default SubscriptionSettings;
