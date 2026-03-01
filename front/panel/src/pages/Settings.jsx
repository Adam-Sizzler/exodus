import { useEffect, useMemo, useState } from 'react';
import { panelSettingsApi } from '../api';

const createDefaultSettings = () => ({
  passkey_settings: {
    enabled: false,
    rpId: null,
    origin: null,
  },
  oauth2_settings: {
    github: { enabled: false, clientId: null, clientSecret: null, allowedEmails: [] },
    yandex: { enabled: false, clientId: null, clientSecret: null, allowedEmails: [] },
    pocketid: { enabled: false, clientId: null, clientSecret: null, plainDomain: null, allowedEmails: [] },
    keycloak: {
      enabled: false,
      clientId: null,
      clientSecret: null,
      realm: null,
      frontendDomain: null,
      keycloakDomain: null,
      allowedEmails: [],
    },
    generic: {
      enabled: false,
      clientId: null,
      clientSecret: null,
      frontendDomain: null,
      authorizationUrl: null,
      tokenUrl: null,
      withPkce: false,
      allowedEmails: [],
    },
  },
  tg_auth_settings: {
    enabled: false,
    botToken: null,
    adminIds: [],
  },
  password_settings: {
    enabled: true,
  },
  branding_settings: {
    title: 'V2RS',
    logoUrl: null,
  },
});

function deepMerge(base, patch) {
  if (!patch || typeof patch !== 'object') {
    return base;
  }
  const merged = Array.isArray(base) ? [...base] : { ...base };
  Object.entries(patch).forEach(([key, value]) => {
    if (value && typeof value === 'object' && !Array.isArray(value) && typeof merged[key] === 'object' && merged[key] !== null && !Array.isArray(merged[key])) {
      merged[key] = deepMerge(merged[key], value);
      return;
    }
    merged[key] = value;
  });
  return merged;
}

function toCsv(arrayValue) {
  if (!Array.isArray(arrayValue)) {
    return '';
  }
  return arrayValue.join(', ');
}

function fromCsv(csvValue) {
  return csvValue
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean);
}

function Settings() {
  const [settings, setSettings] = useState(createDefaultSettings);
  const [loading, setLoading] = useState(true);
  const [savingSettings, setSavingSettings] = useState(false);
  const [settingsError, setSettingsError] = useState('');
  const [settingsSaved, setSettingsSaved] = useState(false);

  const [apiTokens, setApiTokens] = useState([]);
  const [tokensLoading, setTokensLoading] = useState(true);
  const [tokensError, setTokensError] = useState('');
  const [newTokenName, setNewTokenName] = useState('');
  const [creatingToken, setCreatingToken] = useState(false);
  const [createdTokenValue, setCreatedTokenValue] = useState('');

  const apiTokenCount = useMemo(() => apiTokens.length, [apiTokens]);

  const updatePath = (path, value) => {
    const parts = path.split('.');
    setSettings((prev) => {
      const next = { ...prev };
      let cursor = next;
      for (let i = 0; i < parts.length - 1; i += 1) {
        const part = parts[i];
        cursor[part] = { ...(cursor[part] || {}) };
        cursor = cursor[part];
      }
      cursor[parts[parts.length - 1]] = value;
      return next;
    });
  };

  const loadSettings = async () => {
    setLoading(true);
    setSettingsError('');
    try {
      const response = await panelSettingsApi.getSettings();
      const merged = deepMerge(createDefaultSettings(), response?.settings || {});
      setSettings(merged);
    } catch (err) {
      setSettingsError(err.message || 'Failed to load panel settings');
    } finally {
      setLoading(false);
    }
  };

  const loadTokens = async () => {
    setTokensLoading(true);
    setTokensError('');
    try {
      const response = await panelSettingsApi.getApiTokens();
      setApiTokens(response?.tokens || []);
    } catch (err) {
      setTokensError(err.message || 'Failed to load API tokens');
    } finally {
      setTokensLoading(false);
    }
  };

  useEffect(() => {
    Promise.all([loadSettings(), loadTokens()]).catch((error) => {
      console.error('Failed to load settings page state:', error);
    });
  }, []);

  const saveSettings = async (event) => {
    event.preventDefault();
    setSavingSettings(true);
    setSettingsSaved(false);
    setSettingsError('');

    try {
      const payload = {
        passkey_settings: settings.passkey_settings,
        oauth2_settings: settings.oauth2_settings,
        tg_auth_settings: settings.tg_auth_settings,
        password_settings: settings.password_settings,
        branding_settings: settings.branding_settings,
      };
      const response = await panelSettingsApi.updateSettings(payload);
      const merged = deepMerge(createDefaultSettings(), response?.settings || {});
      setSettings(merged);
      setSettingsSaved(true);
      window.dispatchEvent(new CustomEvent('panel-branding-updated', { detail: merged.branding_settings }));
      setTimeout(() => setSettingsSaved(false), 2200);
    } catch (err) {
      setSettingsError(err.message || 'Failed to save panel settings');
    } finally {
      setSavingSettings(false);
    }
  };

  const createToken = async () => {
    const tokenName = newTokenName.trim();
    if (!tokenName) {
      setTokensError('Введите имя токена');
      return;
    }

    setCreatingToken(true);
    setTokensError('');
    setCreatedTokenValue('');
    try {
      const response = await panelSettingsApi.createApiToken(tokenName);
      setNewTokenName('');
      setCreatedTokenValue(response?.token?.token || '');
      await loadTokens();
    } catch (err) {
      setTokensError(err.message || 'Failed to create API token');
    } finally {
      setCreatingToken(false);
    }
  };

  const deleteToken = async (tokenUuid) => {
    if (!window.confirm('Удалить этот API токен?')) {
      return;
    }
    try {
      await panelSettingsApi.deleteApiToken(tokenUuid);
      await loadTokens();
    } catch (err) {
      setTokensError(err.message || 'Failed to delete API token');
    }
  };

  const copyText = async (value) => {
    try {
      await navigator.clipboard.writeText(value);
      alert('Скопировано');
    } catch (err) {
      alert(`Не удалось скопировать: ${err.message}`);
    }
  };

  if (loading) {
    return <div className="loading">Loading panel settings...</div>;
  }

  return (
    <div className="settings-layout">
      <div className="card settings-card">
        <div className="card-header">
          <h2 className="card-title">Настройки Remnawave</h2>
        </div>
        <form className="card-body settings-form-grid" onSubmit={saveSettings}>
          <div className="settings-column">
            <section className="settings-section card">
              <div className="card-header">
                <h3 className="card-title">Способы аутентификации</h3>
              </div>
              <div className="card-body">
                <label className="checkbox-label">
                  <input
                    type="checkbox"
                    checked={Boolean(settings.password_settings?.enabled)}
                    onChange={(event) => updatePath('password_settings.enabled', event.target.checked)}
                  />
                  <span className="checkbox-text"><strong>Пароль</strong></span>
                </label>

                <details className="settings-details">
                  <summary>Passkey</summary>
                  <div className="settings-details-body">
                    <label className="checkbox-label">
                      <input
                        type="checkbox"
                        checked={Boolean(settings.passkey_settings?.enabled)}
                        onChange={(event) => updatePath('passkey_settings.enabled', event.target.checked)}
                      />
                      <span className="checkbox-text"><strong>Включено</strong></span>
                    </label>
                    <div className="form-group">
                      <label className="form-label">Frontend Domain (rpId)</label>
                      <input
                        className="form-input"
                        placeholder="example.com"
                        value={settings.passkey_settings?.rpId || ''}
                        onChange={(event) => updatePath('passkey_settings.rpId', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Backend Origin</label>
                      <input
                        className="form-input"
                        placeholder="https://api.example.com"
                        value={settings.passkey_settings?.origin || ''}
                        onChange={(event) => updatePath('passkey_settings.origin', event.target.value || null)}
                      />
                    </div>
                  </div>
                </details>

                <details className="settings-details">
                  <summary>GitHub OAuth2</summary>
                  <div className="settings-details-body">
                    <label className="checkbox-label">
                      <input
                        type="checkbox"
                        checked={Boolean(settings.oauth2_settings?.github?.enabled)}
                        onChange={(event) => updatePath('oauth2_settings.github.enabled', event.target.checked)}
                      />
                      <span className="checkbox-text"><strong>Включено</strong></span>
                    </label>
                    <div className="form-group">
                      <label className="form-label">Client ID</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.github?.clientId || ''}
                        onChange={(event) => updatePath('oauth2_settings.github.clientId', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Client Secret</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.github?.clientSecret || ''}
                        onChange={(event) => updatePath('oauth2_settings.github.clientSecret', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Allowed Emails (comma separated)</label>
                      <input
                        className="form-input"
                        value={toCsv(settings.oauth2_settings?.github?.allowedEmails)}
                        onChange={(event) => updatePath('oauth2_settings.github.allowedEmails', fromCsv(event.target.value))}
                      />
                    </div>
                  </div>
                </details>

                <details className="settings-details">
                  <summary>Yandex OAuth2</summary>
                  <div className="settings-details-body">
                    <label className="checkbox-label">
                      <input
                        type="checkbox"
                        checked={Boolean(settings.oauth2_settings?.yandex?.enabled)}
                        onChange={(event) => updatePath('oauth2_settings.yandex.enabled', event.target.checked)}
                      />
                      <span className="checkbox-text"><strong>Включено</strong></span>
                    </label>
                    <div className="form-group">
                      <label className="form-label">Client ID</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.yandex?.clientId || ''}
                        onChange={(event) => updatePath('oauth2_settings.yandex.clientId', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Client Secret</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.yandex?.clientSecret || ''}
                        onChange={(event) => updatePath('oauth2_settings.yandex.clientSecret', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Allowed Emails (comma separated)</label>
                      <input
                        className="form-input"
                        value={toCsv(settings.oauth2_settings?.yandex?.allowedEmails)}
                        onChange={(event) => updatePath('oauth2_settings.yandex.allowedEmails', fromCsv(event.target.value))}
                      />
                    </div>
                  </div>
                </details>

                <details className="settings-details">
                  <summary>PocketID OAuth2</summary>
                  <div className="settings-details-body">
                    <label className="checkbox-label">
                      <input
                        type="checkbox"
                        checked={Boolean(settings.oauth2_settings?.pocketid?.enabled)}
                        onChange={(event) => updatePath('oauth2_settings.pocketid.enabled', event.target.checked)}
                      />
                      <span className="checkbox-text"><strong>Включено</strong></span>
                    </label>
                    <div className="form-group">
                      <label className="form-label">Client ID</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.pocketid?.clientId || ''}
                        onChange={(event) => updatePath('oauth2_settings.pocketid.clientId', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Client Secret</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.pocketid?.clientSecret || ''}
                        onChange={(event) => updatePath('oauth2_settings.pocketid.clientSecret', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Domain</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.pocketid?.plainDomain || ''}
                        onChange={(event) => updatePath('oauth2_settings.pocketid.plainDomain', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Allowed Emails (comma separated)</label>
                      <input
                        className="form-input"
                        value={toCsv(settings.oauth2_settings?.pocketid?.allowedEmails)}
                        onChange={(event) => updatePath('oauth2_settings.pocketid.allowedEmails', fromCsv(event.target.value))}
                      />
                    </div>
                  </div>
                </details>

                <details className="settings-details">
                  <summary>Keycloak OAuth2</summary>
                  <div className="settings-details-body">
                    <label className="checkbox-label">
                      <input
                        type="checkbox"
                        checked={Boolean(settings.oauth2_settings?.keycloak?.enabled)}
                        onChange={(event) => updatePath('oauth2_settings.keycloak.enabled', event.target.checked)}
                      />
                      <span className="checkbox-text"><strong>Включено</strong></span>
                    </label>
                    <div className="form-group">
                      <label className="form-label">Client ID</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.keycloak?.clientId || ''}
                        onChange={(event) => updatePath('oauth2_settings.keycloak.clientId', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Client Secret</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.keycloak?.clientSecret || ''}
                        onChange={(event) => updatePath('oauth2_settings.keycloak.clientSecret', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Realm</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.keycloak?.realm || ''}
                        onChange={(event) => updatePath('oauth2_settings.keycloak.realm', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Frontend Domain</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.keycloak?.frontendDomain || ''}
                        onChange={(event) => updatePath('oauth2_settings.keycloak.frontendDomain', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Keycloak Domain</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.keycloak?.keycloakDomain || ''}
                        onChange={(event) => updatePath('oauth2_settings.keycloak.keycloakDomain', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Allowed Emails (comma separated)</label>
                      <input
                        className="form-input"
                        value={toCsv(settings.oauth2_settings?.keycloak?.allowedEmails)}
                        onChange={(event) => updatePath('oauth2_settings.keycloak.allowedEmails', fromCsv(event.target.value))}
                      />
                    </div>
                  </div>
                </details>

                <details className="settings-details">
                  <summary>Generic OAuth2</summary>
                  <div className="settings-details-body">
                    <label className="checkbox-label">
                      <input
                        type="checkbox"
                        checked={Boolean(settings.oauth2_settings?.generic?.enabled)}
                        onChange={(event) => updatePath('oauth2_settings.generic.enabled', event.target.checked)}
                      />
                      <span className="checkbox-text"><strong>Включено</strong></span>
                    </label>
                    <div className="form-group">
                      <label className="form-label">Client ID</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.generic?.clientId || ''}
                        onChange={(event) => updatePath('oauth2_settings.generic.clientId', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Client Secret</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.generic?.clientSecret || ''}
                        onChange={(event) => updatePath('oauth2_settings.generic.clientSecret', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Frontend Domain</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.generic?.frontendDomain || ''}
                        onChange={(event) => updatePath('oauth2_settings.generic.frontendDomain', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Authorization URL</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.generic?.authorizationUrl || ''}
                        onChange={(event) => updatePath('oauth2_settings.generic.authorizationUrl', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Token URL</label>
                      <input
                        className="form-input"
                        value={settings.oauth2_settings?.generic?.tokenUrl || ''}
                        onChange={(event) => updatePath('oauth2_settings.generic.tokenUrl', event.target.value || null)}
                      />
                    </div>
                    <label className="checkbox-label">
                      <input
                        type="checkbox"
                        checked={Boolean(settings.oauth2_settings?.generic?.withPkce)}
                        onChange={(event) => updatePath('oauth2_settings.generic.withPkce', event.target.checked)}
                      />
                      <span className="checkbox-text"><strong>With PKCE</strong></span>
                    </label>
                    <div className="form-group">
                      <label className="form-label">Allowed Emails (comma separated)</label>
                      <input
                        className="form-input"
                        value={toCsv(settings.oauth2_settings?.generic?.allowedEmails)}
                        onChange={(event) => updatePath('oauth2_settings.generic.allowedEmails', fromCsv(event.target.value))}
                      />
                    </div>
                  </div>
                </details>

                <details className="settings-details">
                  <summary>Telegram Auth</summary>
                  <div className="settings-details-body">
                    <label className="checkbox-label">
                      <input
                        type="checkbox"
                        checked={Boolean(settings.tg_auth_settings?.enabled)}
                        onChange={(event) => updatePath('tg_auth_settings.enabled', event.target.checked)}
                      />
                      <span className="checkbox-text"><strong>Включено</strong></span>
                    </label>
                    <div className="form-group">
                      <label className="form-label">Bot Token</label>
                      <input
                        className="form-input"
                        value={settings.tg_auth_settings?.botToken || ''}
                        onChange={(event) => updatePath('tg_auth_settings.botToken', event.target.value || null)}
                      />
                    </div>
                    <div className="form-group">
                      <label className="form-label">Allowed Telegram IDs (comma separated)</label>
                      <input
                        className="form-input"
                        value={toCsv(settings.tg_auth_settings?.adminIds)}
                        onChange={(event) => updatePath('tg_auth_settings.adminIds', fromCsv(event.target.value))}
                      />
                    </div>
                  </div>
                </details>
              </div>
            </section>
          </div>

          <div className="settings-column">
            <section className="settings-section card">
              <div className="card-header">
                <h3 className="card-title">Настройки кастомизации</h3>
              </div>
              <div className="card-body">
                <div className="form-group">
                  <label className="form-label">Название бренда</label>
                  <input
                    className="form-input"
                    value={settings.branding_settings?.title || ''}
                    onChange={(event) => updatePath('branding_settings.title', event.target.value || 'V2RS')}
                    placeholder="V2RS"
                  />
                </div>
                <div className="form-group">
                  <label className="form-label">Ссылка на логотип</label>
                  <input
                    className="form-input"
                    value={settings.branding_settings?.logoUrl || ''}
                    onChange={(event) => updatePath('branding_settings.logoUrl', event.target.value || null)}
                    placeholder="https://example.com/logo.png"
                  />
                </div>
              </div>
            </section>

            <section className="settings-section card">
              <div className="card-header">
                <h3 className="card-title">API токены</h3>
              </div>
              <div className="card-body">
                <div className="settings-token-create">
                  <input
                    className="form-input"
                    value={newTokenName}
                    onChange={(event) => setNewTokenName(event.target.value)}
                    placeholder="Имя токена (например, dev)"
                  />
                  <button className="btn btn-primary" type="button" onClick={createToken} disabled={creatingToken}>
                    {creatingToken ? 'Создание...' : 'Создать'}
                  </button>
                </div>

                {createdTokenValue && (
                  <div className="settings-token-created">
                    <div className="form-hint">Новый токен (сохраните):</div>
                    <div className="settings-token-created-row">
                      <code>{createdTokenValue}</code>
                      <button className="btn btn-secondary" type="button" onClick={() => copyText(createdTokenValue)}>
                        Копировать
                      </button>
                    </div>
                  </div>
                )}
                {tokensLoading ? (
                  <div className="loading-inline">Loading API tokens...</div>
                ) : (
                  <div className="settings-token-list">
                    {apiTokenCount === 0 ? (
                      <div className="form-hint">Токены пока не созданы.</div>
                    ) : (
                      apiTokens.map((token) => (
                        <div key={token.uuid} className="settings-token-item">
                          <div className="settings-token-head">
                            <strong>{token.token_name}</strong>
                            <span>{new Date(token.created_at).toLocaleString()}</span>
                          </div>
                          <code>{token.token}</code>
                          <div className="settings-token-actions">
                            <button className="btn btn-secondary" type="button" onClick={() => copyText(token.token)}>
                              Копировать
                            </button>
                            <button className="btn btn-danger" type="button" onClick={() => deleteToken(token.uuid)}>
                              Удалить
                            </button>
                          </div>
                        </div>
                      ))
                    )}
                  </div>
                )}

              </div>
            </section>
          </div>

          {(settingsError || tokensError) && (
            <div className="alert alert-error">
              {settingsError || tokensError}
            </div>
          )}

          <div className="settings-form-actions">
            <button className="btn btn-primary" type="submit" disabled={savingSettings}>
              {savingSettings ? 'Сохранение...' : 'Сохранить'}
            </button>
            {settingsSaved && <span className="save-success-text">Настройки сохранены</span>}
          </div>
        </form>
      </div>

    </div>
  );
}

export default Settings;
