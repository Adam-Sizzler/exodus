import { useEffect, useMemo, useState } from 'react';
import { configProfilesApi } from '../api';
import AppSelect from './AppSelect';

const API_SCHEMAS = ['http', 'https', 'grpc'];

const parseConfigProfileInbounds = (profile) => {
  if (!profile || typeof profile !== 'object') {
    return [];
  }

  if (Array.isArray(profile.inbounds)) {
    return profile.inbounds;
  }

  const rawConfig = profile.config;
  if (!rawConfig) {
    return [];
  }

  try {
    const parsedConfig = typeof rawConfig === 'string' ? JSON.parse(rawConfig) : rawConfig;
    if (parsedConfig && Array.isArray(parsedConfig.inbounds)) {
      return parsedConfig.inbounds;
    }
  } catch {
    // ignore parse errors and fallback to empty list
  }

  return [];
};

const DEFAULT_FORM_DATA = {
  name: '',
  address: '',
  port: 9253,
  api_schema: 'grpc',
  api_path: '',
  is_disabled: false,
  consumption_multiplier: 100,
  is_traffic_tracking_active: true,
  traffic_reset_day: 1,
  traffic_limit_bytes: 0,
  notify_percent: 80,
  country_code: '',
  tags: [],
  active_config_profile_uuid: '',
  provider_uuid: ''
};

const parseTags = (value) => {
  if (Array.isArray(value)) {
    return value.map((tag) => String(tag).trim()).filter(Boolean);
  }

  if (typeof value === 'string') {
    return value.split(',').map((tag) => tag.trim()).filter(Boolean);
  }

  return [];
};

const buildFormDataFromNode = (node) => {
  if (!node) {
    return { ...DEFAULT_FORM_DATA };
  }

  return {
    ...DEFAULT_FORM_DATA,
    name: node.name ?? '',
    address: node.address ?? '',
    port: node.port ?? 9253,
    api_schema: node.api_schema ?? 'grpc',
    api_path: node.api_path ?? '',
    is_disabled: Boolean(node.is_disabled),
    consumption_multiplier: node.consumption_multiplier ?? 100,
    is_traffic_tracking_active: node.is_traffic_tracking_active ?? true,
    traffic_reset_day: node.traffic_reset_day ?? 1,
    traffic_limit_bytes: node.traffic_limit_bytes ?? 0,
    notify_percent: node.notify_percent ?? 80,
    country_code: node.country_code ?? '',
    tags: parseTags(node.tags),
    active_config_profile_uuid: node.active_config_profile_uuid ?? '',
    provider_uuid: node.provider_uuid ?? ''
  };
};

function NodeModal({ node, onClose, onSave }) {
  const initialFormData = buildFormDataFromNode(node);
  const [formData, setFormData] = useState(initialFormData);
  const [tagsInput, setTagsInput] = useState(initialFormData.tags.join(', '));
  const [configProfiles, setConfigProfiles] = useState([]);
  const [profilesLoading, setProfilesLoading] = useState(true);
  const [savingNode, setSavingNode] = useState(false);

  useEffect(() => {
    const fetchConfigProfiles = async () => {
      try {
        setProfilesLoading(true);
        const data = await configProfilesApi.getAllWithInbounds();
        setConfigProfiles(data.profiles || []);
      } catch (err) {
        console.error('Error fetching config profiles:', err);
        setConfigProfiles([]);
      } finally {
        setProfilesLoading(false);
      }
    };

    fetchConfigProfiles();
  }, []);

  useEffect(() => {
    const nextFormData = buildFormDataFromNode(node);
    setFormData(nextFormData);
    setTagsInput(nextFormData.tags.join(', '));
  }, [node]);

  const handleChange = (field, value) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handleCheckboxChange = (field, event) => {
    handleChange(field, event.target.checked);
  };

  const handleTagsChange = (event) => {
    const value = event.target.value;
    const tags = parseTags(value);
    setTagsInput(value);
    handleChange('tags', tags);
  };

  const handleSubmit = async (event) => {
    event.preventDefault();

    const submitData = { ...formData };
    if (submitData.active_config_profile_uuid === '') {
      submitData.active_config_profile_uuid = null;
    }
    if (submitData.provider_uuid === '') {
      submitData.provider_uuid = null;
    }

    let payload = submitData;
    if (node) {
      payload = {};
      Object.keys(submitData).forEach((key) => {
        let previousValue = node[key];
        if (key === 'is_traffic_tracking_active' && previousValue === undefined) {
          previousValue = true;
        }
        if (key === 'tags') {
          previousValue = parseTags(previousValue);
        }
        if (key === 'active_config_profile_uuid' || key === 'provider_uuid') {
          previousValue = previousValue || null;
        }

        const currentValue = submitData[key];
        const changed = Array.isArray(currentValue)
          ? JSON.stringify(currentValue) !== JSON.stringify(previousValue)
          : currentValue !== previousValue;

        if (changed) {
          payload[key] = currentValue;
        }
      });
    }

    try {
      setSavingNode(true);
      await onSave(payload);
    } finally {
      setSavingNode(false);
    }
  };

  const selectedProfile = useMemo(() => {
    if (!formData.active_config_profile_uuid) {
      return null;
    }
    return configProfiles.find((profile) => profile.uuid === formData.active_config_profile_uuid) || null;
  }, [configProfiles, formData.active_config_profile_uuid]);

  const selectedProfileInbounds = useMemo(
    () => parseConfigProfileInbounds(selectedProfile),
    [selectedProfile]
  );

  const selectedProfileBadges = useMemo(() => {
    return selectedProfileInbounds
      .map((inbound) => {
        if (inbound?.tag) return String(inbound.tag);
        if (inbound?.port !== undefined && inbound?.port !== null) return String(inbound.port);
        return '';
      })
      .filter(Boolean)
      .slice(0, 8);
  }, [selectedProfileInbounds]);

  const handleProfileChangeClick = () => {
    // intentionally empty by request
  };

  return (
    <div className="modal-overlay user-create-overlay" onClick={onClose}>
      <div className="modal modal-large user-create-modal node-create-modal" onClick={(event) => event.stopPropagation()}>
        <header className="modal-header user-create-header">
          <div className="user-create-title-group">
            <div className="user-create-title-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M5 12l2.5 -2.5a2.121 2.121 0 0 1 3 0l3 3a2.121 2.121 0 0 0 3 0l2.5 -2.5" />
                <path d="M3 5m0 2a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v10a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2z" />
              </svg>
            </div>
            <div className="user-create-title-stack">
              <h3 className="modal-title">{node ? 'Редактировать ноду' : 'Создать ноду'}</h3>
              <p className="user-create-title-subtitle">
                {node ? 'Измените параметры подключения и мониторинга' : 'Заполните параметры новой ноды'}
              </p>
            </div>
          </div>
          <button type="button" className="modal-close" onClick={onClose} aria-label="Закрыть" disabled={savingNode}>
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          </button>
        </header>

        <form className="user-create-form node-modal-form" onSubmit={handleSubmit}>
          <div className="modal-body modal-body-columns user-create-body">
            <div className="user-create-column">
              <section className="form-section-card user-create-card">
                <div className="form-section-header user-create-card-header">
                  <span className="user-create-card-icon blue" aria-hidden="true">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M3 12h6l3 -9l4 18l3 -9h2" />
                    </svg>
                  </span>
                  <h4 className="form-section-title">Параметры ноды</h4>
                </div>
                <div className="form-section-content single-column user-create-card-content">
                  <div className="form-group">
                    <label className="form-label">Имя ноды *</label>
                    <input
                      type="text"
                      className="form-input monospace-field"
                      value={formData.name}
                      onChange={(event) => handleChange('name', event.target.value)}
                      placeholder="DE-Frankfurt-01"
                      required
                    />
                  </div>

                  <div className="form-group">
                    <label className="form-label">Адрес *</label>
                    <input
                      type="text"
                      className="form-input monospace-field"
                      value={formData.address}
                      onChange={(event) => handleChange('address', event.target.value)}
                      placeholder="192.168.1.1 или node.example.com"
                      required
                    />
                  </div>

                  <div className="user-create-inline-row">
                    <div className="form-group">
                      <label className="form-label">Порт</label>
                      <input
                        type="number"
                        className="form-input"
                        value={formData.port}
                        onChange={(event) => handleChange('port', parseInt(event.target.value, 10) || 0)}
                        min="1"
                        max="65535"
                      />
                    </div>

                    <div className="form-group">
                      <label className="form-label">API Schema</label>
                      <AppSelect
                        value={formData.api_schema}
                        onChange={(event) => handleChange('api_schema', event.target.value)}
                      >
                        {API_SCHEMAS.map((schema) => (
                          <option key={schema} value={schema}>
                            {schema}
                          </option>
                        ))}
                      </AppSelect>
                    </div>
                  </div>

                  <div className="form-group">
                    <label className="form-label">API Path</label>
                    <input
                      type="text"
                      className="form-input monospace-field"
                      value={formData.api_path}
                      onChange={(event) => handleChange('api_path', event.target.value)}
                      placeholder="/api/v1"
                    />
                  </div>

                </div>
              </section>

              <section className="form-section-card user-create-card">
                <div className="form-section-header user-create-card-header">
                  <span className="user-create-card-icon teal" aria-hidden="true">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M6 20h12" />
                      <path d="M7 20v-14a1 1 0 0 1 1 -1h8a1 1 0 0 1 1 1v14" />
                      <path d="M9 9h6" />
                      <path d="M9 13h6" />
                    </svg>
                  </span>
                  <h4 className="form-section-title">Идентификация</h4>
                </div>
                <div className="form-section-content single-column user-create-card-content">
                  <div className="form-group">
                    <label className="form-label">Country Code</label>
                    <input
                      type="text"
                      className="form-input monospace-field"
                      value={formData.country_code}
                      onChange={(event) => handleChange('country_code', event.target.value.toUpperCase())}
                      placeholder="DE, US, NL"
                      maxLength={2}
                    />
                  </div>

                  <div className="form-group">
                    <label className="form-label">Tags</label>
                    <input
                      type="text"
                      className="form-input"
                      value={tagsInput}
                      onChange={handleTagsChange}
                      placeholder="premium, fast, grpc"
                    />
                  </div>

                  <div className="form-group form-group-last">
                    <label className="form-label">Provider UUID</label>
                    <input
                      type="text"
                      className="form-input monospace-field"
                      value={formData.provider_uuid}
                      onChange={(event) => handleChange('provider_uuid', event.target.value)}
                      placeholder="Оставьте пустым"
                    />
                  </div>
                </div>
              </section>
            </div>

            <div className="user-create-column">
              <section className="form-section-card user-create-card">
                <div className="form-section-header user-create-card-header">
                  <span className="user-create-card-icon teal" aria-hidden="true">
                    <svg viewBox="0 0 24 24" fill="currentColor">
                      <path d="M16.3696 2.5006 12.0006 5 7.6303 7.5006v-5L12.0006 0Zm6.1177 3.499.0028 4.986-8.7282-4.9929 4.3564-2.4923Zm-4.369 9.5085-.0014 4.9972 4.3774-2.5007-.0028-5.018-4.3732-2.502zM7.6274 21.502 12.0006 24l4.369-2.4952v-4.9972zM7.6303 9.5v5.0014l4.3703 2.4992 4.369-2.4937V9.5001l-4.369-2.4993Zm-6.1248 8.5044.0028-5.0055 8.7464 5.0027-4.376 2.5008Zm4.376-14.504L1.5125 6.001l-.0028 4.9985 4.3718 2.502z" />
                    </svg>
                  </span>
                  <h4 className="form-section-title">Конфигурация ядра</h4>
                </div>
                <div className="form-section-content single-column user-create-card-content">
                  <div className="node-profile-intro">
                    <div className="node-profile-intro-title-row">
                      <svg viewBox="0 0 35 35" fill="currentColor" aria-hidden="true">
                        <path d="M16.6961 15.2606L16.5825 3.49701C16.5718 2.38439 15.025 2.11843 14.6433 3.16356L11.7279 11.1447C11.6384 11.3898 11.4566 11.5902 11.2213 11.7031L5.66765 14.3687C4.70841 14.8291 5.03635 16.2703 6.10036 16.2703H15.6962C16.2522 16.2703 16.7015 15.8166 16.6961 15.2606Z"></path>
                        <path d="M18.6471 15.2703V5.88936C18.6471 4.84679 20.0428 4.49998 20.5308 5.4213L23.5833 11.1845C23.7 11.4049 23.8948 11.5737 24.1296 11.6578L31.5829 14.3289C32.6388 14.7073 32.3671 16.2703 31.2455 16.2703H19.6471C19.0948 16.2703 18.6471 15.8226 18.6471 15.2703Z"></path>
                        <path d="M18.6471 31.4643V19.3784C18.6471 18.8261 19.0948 18.3784 19.6471 18.3784H29.2853C30.3376 18.3784 30.676 19.7947 29.7374 20.2704L24.1129 23.1208C23.889 23.2343 23.716 23.4278 23.6281 23.663L20.5839 31.8141C20.1941 32.8578 18.6471 32.5783 18.6471 31.4643Z"></path>
                        <path d="M16.7059 28.9873V19.3784C16.7059 18.8261 16.2582 18.3784 15.7059 18.3784H3.83963C2.71522 18.3784 2.44656 19.9473 3.50691 20.3214L11.5457 23.1578C11.7987 23.247 12.0052 23.4342 12.1188 23.6772L14.8 29.4109C15.2531 30.3797 16.7059 30.0568 16.7059 28.9873Z"></path>
                      </svg>
                      <p className="node-profile-intro-title">Профили</p>
                    </div>
                    <p className="node-profile-intro-subtitle">Выберите профиль, который будет применен для этой ноды.</p>
                  </div>

                  <div className="node-profile-selected-card">
                    {profilesLoading ? (
                      <p className="node-profile-empty">Загрузка профилей...</p>
                    ) : selectedProfile ? (
                      <>
                        <div className="node-profile-selected-head">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                            <path d="M14 3v4a1 1 0 0 0 1 1h4"></path>
                            <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"></path>
                          </svg>
                          <p className="node-profile-selected-name">{selectedProfile.name || 'Без названия'}</p>
                          <span className="node-profile-selected-count">{selectedProfileInbounds.length}</span>
                        </div>
                        {selectedProfileBadges.length > 0 ? (
                          <div className="node-profile-pill-list">
                            {selectedProfileBadges.map((badge, index) => (
                              <span key={`${badge}-${index}`} className="node-profile-pill">
                                {badge}
                              </span>
                            ))}
                          </div>
                        ) : (
                          <p className="node-profile-empty">В профиле нет инбаундов</p>
                        )}
                      </>
                    ) : (
                      <p className="node-profile-empty">Профиль не выбран</p>
                    )}
                  </div>

                  <button type="button" className="node-profile-edit-btn" onClick={handleProfileChangeClick}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                      <path d="M14 3v4a1 1 0 0 0 1 1h4"></path>
                      <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"></path>
                    </svg>
                    <span>Изменить</span>
                  </button>
                </div>
              </section>

              <section className="form-section-card user-create-card">
                <div className="form-section-header user-create-card-header">
                  <span className="user-create-card-icon orange" aria-hidden="true">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M4 19l16 0" />
                      <path d="M4 15l4 -6l4 2l4 -5l4 4" />
                    </svg>
                  </span>
                  <h4 className="form-section-title">Трафик и лимиты</h4>
                </div>
                <div className="form-section-content single-column user-create-card-content">
                  <div className="user-create-inline-row">
                    <div className="form-group">
                      <label className="form-label">Лимит трафика (bytes)</label>
                      <input
                        type="number"
                        className="form-input"
                        value={formData.traffic_limit_bytes}
                        onChange={(event) => handleChange('traffic_limit_bytes', parseInt(event.target.value, 10) || 0)}
                        min="0"
                      />
                    </div>

                    <div className="form-group">
                      <label className="form-label">Traffic Reset Day</label>
                      <input
                        type="number"
                        className="form-input"
                        value={formData.traffic_reset_day}
                        onChange={(event) => handleChange('traffic_reset_day', parseInt(event.target.value, 10) || 1)}
                        min="1"
                        max="31"
                      />
                    </div>
                  </div>

                  <div className="form-group">
                    <label className="form-label">Notify at (%)</label>
                    <input
                      type="number"
                      className="form-input"
                      value={formData.notify_percent}
                      onChange={(event) => handleChange('notify_percent', parseInt(event.target.value, 10) || 0)}
                      min="0"
                      max="100"
                    />
                  </div>

                  <div className="form-group form-group-last">
                    <label className="form-label">Коэффициент потребления (%)</label>
                    <input
                      type="number"
                      className="form-input"
                      value={formData.consumption_multiplier}
                      onChange={(event) => handleChange('consumption_multiplier', parseInt(event.target.value, 10) || 0)}
                      min="0"
                    />
                  </div>
                </div>
              </section>

              <section className="form-section-card user-create-card">
                <div className="form-section-header user-create-card-header">
                  <span className="user-create-card-icon violet" aria-hidden="true">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M9 12l2 2l4 -4" />
                      <path d="M12 3c4.97 0 9 4.03 9 9s-4.03 9 -9 9s-9 -4.03 -9 -9s4.03 -9 9 -9z" />
                    </svg>
                  </span>
                  <h4 className="form-section-title">Статус ноды</h4>
                </div>
                <div className="form-section-content single-column user-create-card-content">
                  <div className="form-checkbox">
                    <input
                      type="checkbox"
                      id="is_disabled"
                      checked={formData.is_disabled}
                      onChange={(event) => handleCheckboxChange('is_disabled', event)}
                    />
                    <label htmlFor="is_disabled">Нода отключена</label>
                  </div>

                  <div className="form-checkbox form-group-last">
                    <input
                      type="checkbox"
                      id="is_traffic_tracking_active"
                      checked={formData.is_traffic_tracking_active}
                      onChange={(event) => handleCheckboxChange('is_traffic_tracking_active', event)}
                    />
                    <label htmlFor="is_traffic_tracking_active">Отслеживать трафик</label>
                  </div>
                </div>
              </section>

            </div>
          </div>

          <footer className="modal-footer user-create-footer">
            <button
              type="button"
              className="squad-editor-summary-btn cancel user-create-btn-cancel"
              onClick={onClose}
              disabled={savingNode}
            >
              Отмена
            </button>
            <button type="submit" className="squad-editor-summary-btn save user-create-btn-save" disabled={savingNode}>
              <svg viewBox="0 0 256 256" fill="currentColor" aria-hidden="true">
                <path d="M216,83.31V208a8,8,0,0,1-8,8H176V152a8,8,0,0,0-8-8H88a8,8,0,0,0-8,8v64H48a8,8,0,0,1-8-8V48a8,8,0,0,1,8-8H172.69a8,8,0,0,1,5.65,2.34l35.32,35.32A8,8,0,0,1,216,83.31Z" opacity="0.2"></path>
                <path d="M219.31,72,184,36.69A15.86,15.86,0,0,0,172.69,32H48A16,16,0,0,0,32,48V208a16,16,0,0,0,16,16H208a16,16,0,0,0,16-16V83.31A15.86,15.86,0,0,0,219.31,72ZM168,208H88V152h80Zm40,0H184V152a16,16,0,0,0-16-16H88a16,16,0,0,0-16,16v56H48V48H172.69L208,83.31ZM160,72a8,8,0,0,1-8,8H96a8,8,0,0,1,0-16h56A8,8,0,0,1,160,72Z"></path>
              </svg>
              <span>{savingNode ? 'Сохранение...' : node ? 'Сохранить' : 'Создать'}</span>
            </button>
          </footer>
        </form>
      </div>
    </div>
  );
}

export default NodeModal;
