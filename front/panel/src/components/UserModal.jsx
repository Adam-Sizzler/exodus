import { useCallback, useEffect, useMemo, useState } from 'react';
import { configProfilesApi, squadsApi } from '../api';
import AppSelect from './AppSelect';
import SideDrawerPanel from './SideDrawerPanel';

const STATUS_OPTIONS = ['ACTIVE', 'DISABLED', 'LIMITED', 'EXPIRED'];
const TRAFFIC_STRATEGY_OPTIONS = ['NO_RESET', 'DAY', 'WEEK', 'MONTH'];
const INBOUNDS_TAB = {
  profiles: 'profiles',
  list: 'list',
};
const FLAT_FILTER = {
  all: 'all',
  selected: 'selected',
  unselected: 'unselected',
};

const TRAFFIC_STRATEGY_LABELS = {
  NO_RESET: 'Никогда не сбрасывать',
  DAY: 'Каждый день',
  WEEK: 'Каждую неделю',
  MONTH: 'Каждый месяц'
};

const createDefaultFormData = () => ({
  username: '',
  status: 'ACTIVE',
  traffic_limit_bytes: 0,
  traffic_limit_strategy: 'NO_RESET',
  expire_at: '',
  description: '',
  tag: '',
  telegram_id: '',
  email: '',
  hwid_device_limit: 0,
  last_triggered_threshold: 0
});

const getInboundPort = (inbound) => {
  return inbound?.port ?? inbound?.listen_port ?? inbound?.raw_inbound?.listen_port ?? 'N/A';
};

const getInboundProtocol = (inbound) => {
  return inbound?.type ?? inbound?.protocol ?? inbound?.raw_inbound?.protocol ?? 'unknown';
};

const getInboundSecurity = (inbound) => {
  return inbound?.security ?? inbound?.raw_inbound?.security ?? 'none';
};

const createDefaultSelectedSquads = (user) => {
  const next = new Set();
  if (user?.internal_squad_uuid) {
    next.add(user.internal_squad_uuid);
  }
  return next;
};

const formatDateParts = (isoValue) => {
  if (!isoValue) {
    return { date: '', time: '23:59' };
  }

  const parsed = new Date(isoValue);
  if (Number.isNaN(parsed.getTime())) {
    return { date: '', time: '23:59' };
  }

  const year = parsed.getFullYear();
  const month = String(parsed.getMonth() + 1).padStart(2, '0');
  const day = String(parsed.getDate()).padStart(2, '0');
  const hours = String(parsed.getHours()).padStart(2, '0');
  const minutes = String(parsed.getMinutes()).padStart(2, '0');

  return {
    date: `${year}-${month}-${day}`,
    time: `${hours}:${minutes}`
  };
};

const buildExpireAt = (date, time) => {
  if (!date || !time) return '';
  const parsed = new Date(`${date}T${time}:00`);
  if (Number.isNaN(parsed.getTime())) return '';
  return parsed.toISOString();
};

function UserModal({ user, onClose, onSave }) {
  const [formData, setFormData] = useState(createDefaultFormData);
  const [expireDate, setExpireDate] = useState('');
  const [expireTime, setExpireTime] = useState('23:59');
  const [squads, setSquads] = useState([]);
  const [configProfiles, setConfigProfiles] = useState([]);
  const [loadingSquads, setLoadingSquads] = useState(false);
  const [loadingUserSquads, setLoadingUserSquads] = useState(false);
  const [savingUser, setSavingUser] = useState(false);
  const [squadSearch, setSquadSearch] = useState('');
  const [selectedInternalSquads, setSelectedInternalSquads] = useState(new Set());
  const [showInboundsDrawer, setShowInboundsDrawer] = useState(false);
  const [selectedSquad, setSelectedSquad] = useState(null);
  const [selectedInbounds, setSelectedInbounds] = useState(new Set());
  const [inboundsTab, setInboundsTab] = useState(INBOUNDS_TAB.profiles);
  const [inboundsSearchQuery, setInboundsSearchQuery] = useState('');
  const [flatFilter, setFlatFilter] = useState(FLAT_FILTER.all);
  const [expandedProfiles, setExpandedProfiles] = useState(new Set());
  const [loadingInboundsDrawer, setLoadingInboundsDrawer] = useState(false);
  const [savingInbounds, setSavingInbounds] = useState(false);

  useEffect(() => {
    if (user) {
      const nextForm = {
        username: user.username || '',
        status: user.status || 'ACTIVE',
        traffic_limit_bytes: user.traffic_limit_bytes || 0,
        traffic_limit_strategy: user.traffic_limit_strategy || 'NO_RESET',
        expire_at: user.expire_at || '',
        description: user.description || '',
        tag: user.tag || '',
        telegram_id: user.telegram_id || '',
        email: user.email || '',
        hwid_device_limit: user.hwid_device_limit || 0,
        last_triggered_threshold: user.last_triggered_threshold || 0
      };
      setFormData(nextForm);
      const parts = formatDateParts(nextForm.expire_at);
      setExpireDate(parts.date);
      setExpireTime(parts.time);
      setSelectedInternalSquads(createDefaultSelectedSquads(user));
      return;
    }

    setFormData(createDefaultFormData());
    setExpireDate('');
    setExpireTime('23:59');
    setSelectedInternalSquads(new Set());
  }, [user]);

  const loadSquads = async () => {
    setLoadingSquads(true);
    try {
      const data = await squadsApi.getAll();
      setSquads(data.squads || []);
    } catch (err) {
      console.error('Error fetching squads:', err);
    } finally {
      setLoadingSquads(false);
    }
  };

  const loadConfigProfiles = async () => {
    try {
      const data = await configProfilesApi.getAllWithInbounds();
      setConfigProfiles(data.profiles || []);
    } catch (err) {
      console.error('Error fetching config profiles:', err);
      setConfigProfiles([]);
    }
  };

  useEffect(() => {
    loadSquads();
    loadConfigProfiles();
  }, []);

  useEffect(() => {
    let cancelled = false;

    if (!user) {
      setLoadingUserSquads(false);
      return () => {
        cancelled = true;
      };
    }

    if (squads.length === 0) {
      setLoadingUserSquads(false);
      return () => {
        cancelled = true;
      };
    }

    const baseSelection = createDefaultSelectedSquads(user);
    const userId = user.t_id ?? user.user_id;

    if (userId === undefined || userId === null) {
      setSelectedInternalSquads(baseSelection);
      setLoadingUserSquads(false);
      return () => {
        cancelled = true;
      };
    }

    const loadUserSquadMemberships = async () => {
      setLoadingUserSquads(true);
      try {
        const memberships = await Promise.all(
          squads.map(async (squad) => {
            const data = await squadsApi.getMembers(squad.uuid);
            const members = data?.squad_members || [];
            const isMember = members.some((member) => String(member.user_id ?? member.t_id) === String(userId));
            return isMember ? squad.uuid : null;
          })
        );

        if (cancelled) return;

        const next = new Set(baseSelection);
        memberships.forEach((squadUuid) => {
          if (squadUuid) {
            next.add(squadUuid);
          }
        });
        setSelectedInternalSquads(next);
      } catch (err) {
        console.error('Error fetching squad memberships for user:', err);
        if (!cancelled) {
          setSelectedInternalSquads(baseSelection);
        }
      } finally {
        if (!cancelled) {
          setLoadingUserSquads(false);
        }
      }
    };

    loadUserSquadMemberships();

    return () => {
      cancelled = true;
    };
  }, [user, squads]);

  const filteredSquads = useMemo(() => {
    const query = squadSearch.trim().toLowerCase();
    if (!query) return squads;
    return squads.filter((squad) => {
      const name = String(squad.name || '').toLowerCase();
      const uuid = String(squad.uuid || '').toLowerCase();
      return name.includes(query) || uuid.includes(query);
    });
  }, [squadSearch, squads]);

  const handleChange = (field, value) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const toggleInternalSquad = (squadUuid) => {
    setSelectedInternalSquads((previous) => {
      const next = new Set(previous);
      if (next.has(squadUuid)) {
        next.delete(squadUuid);
      } else {
        next.add(squadUuid);
      }
      return next;
    });
  };

  const closeInboundsDrawer = () => {
    if (savingInbounds) {
      return;
    }
    setShowInboundsDrawer(false);
    setSelectedSquad(null);
    setLoadingInboundsDrawer(false);
    setInboundsTab(INBOUNDS_TAB.profiles);
    setFlatFilter(FLAT_FILTER.all);
    setInboundsSearchQuery('');
    setExpandedProfiles(new Set());
  };

  const handleOpenSquadEditor = async (squad) => {
    setSelectedSquad(squad);
    setShowInboundsDrawer(true);
    setLoadingInboundsDrawer(true);
    setInboundsTab(INBOUNDS_TAB.profiles);
    setFlatFilter(FLAT_FILTER.all);
    setInboundsSearchQuery('');
    setExpandedProfiles(new Set());

    try {
      const details = await squadsApi.getDetails(squad.uuid);
      const currentInbounds = details?.squad?.inbounds || [];
      setSelectedInbounds(new Set(currentInbounds.map((inbound) => inbound.uuid)));
    } catch (err) {
      alert(`Не удалось загрузить инбаунды сквада: ${err.message}`);
      setShowInboundsDrawer(false);
      setSelectedSquad(null);
    } finally {
      setLoadingInboundsDrawer(false);
    }
  };

  const handleSaveInbounds = async () => {
    if (!selectedSquad || savingInbounds) {
      return;
    }

    try {
      setSavingInbounds(true);
      await squadsApi.setInbounds(selectedSquad.uuid, Array.from(selectedInbounds));
      await loadSquads();
      closeInboundsDrawer();
    } catch (err) {
      alert(`Не удалось сохранить инбаунды: ${err.message}`);
    } finally {
      setSavingInbounds(false);
    }
  };

  const toggleInbound = (inboundUuid) => {
    setSelectedInbounds((previous) => {
      const next = new Set(previous);
      if (next.has(inboundUuid)) {
        next.delete(inboundUuid);
      } else {
        next.add(inboundUuid);
      }
      return next;
    });
  };

  const setProfileSelection = (profileInbounds, shouldSelect) => {
    setSelectedInbounds((previous) => {
      const next = new Set(previous);
      profileInbounds.forEach((inbound) => {
        if (shouldSelect) {
          next.add(inbound.uuid);
        } else {
          next.delete(inbound.uuid);
        }
      });
      return next;
    });
  };

  const toggleProfileAccordion = (profileUuid) => {
    setExpandedProfiles((previous) => {
      const next = new Set(previous);
      if (next.has(profileUuid)) {
        next.delete(profileUuid);
      } else {
        next.add(profileUuid);
      }
      return next;
    });
  };

  const normalizedInboundsSearch = inboundsSearchQuery.trim().toLowerCase();
  const allInbounds = useMemo(() => {
    return configProfiles.flatMap((profile) => {
      const profileInbounds = Array.isArray(profile.inbounds) ? profile.inbounds : [];
      return profileInbounds.map((inbound) => ({
        ...inbound,
        profileUuid: profile.uuid,
        profileName: profile.name,
      }));
    });
  }, [configProfiles]);

  const inboundMatchesSearch = useCallback((inbound, profileName = '') => {
    if (!normalizedInboundsSearch) {
      return true;
    }

    const searchable = [
      inbound?.tag,
      inbound?.type,
      String(getInboundPort(inbound)),
      profileName,
      inbound?.security,
    ]
      .filter(Boolean)
      .join(' ')
      .toLowerCase();

    return searchable.includes(normalizedInboundsSearch);
  }, [normalizedInboundsSearch]);

  const profileGroups = useMemo(() => {
    return configProfiles
      .map((profile) => {
        const profileInbounds = Array.isArray(profile.inbounds) ? profile.inbounds : [];
        const visibleInbounds = profileInbounds.filter((inbound) =>
          inboundMatchesSearch(inbound, profile.name)
        );
        const selectedCount = profileInbounds.reduce(
          (accumulator, inbound) => accumulator + (selectedInbounds.has(inbound.uuid) ? 1 : 0),
          0
        );

        return {
          ...profile,
          profileInbounds,
          visibleInbounds,
          selectedCount,
        };
      })
      .filter((profile) => {
        if (!normalizedInboundsSearch) {
          return true;
        }

        return (
          profile.visibleInbounds.length > 0 ||
          profile.name?.toLowerCase().includes(normalizedInboundsSearch)
        );
      });
  }, [configProfiles, inboundMatchesSearch, normalizedInboundsSearch, selectedInbounds]);

  const flatInbounds = useMemo(() => {
    let result = allInbounds.filter((inbound) =>
      inboundMatchesSearch(inbound, inbound.profileName)
    );

    if (flatFilter === FLAT_FILTER.selected) {
      result = result.filter((inbound) => selectedInbounds.has(inbound.uuid));
    } else if (flatFilter === FLAT_FILTER.unselected) {
      result = result.filter((inbound) => !selectedInbounds.has(inbound.uuid));
    }

    return result;
  }, [allInbounds, flatFilter, inboundMatchesSearch, selectedInbounds]);

  const handleSubmit = async (event) => {
    event.preventDefault();

    const usernameRegex = /^[a-zA-Z0-9_-]+$/;
    if (!usernameRegex.test(formData.username)) {
      alert('Имя пользователя может содержать только буквы, цифры, _ и -');
      return;
    }

    if (!STATUS_OPTIONS.includes(formData.status)) {
      alert('Некорректный статус пользователя');
      return;
    }

    if (savingUser) {
      return;
    }

    const expireAt = buildExpireAt(expireDate, expireTime);
    if (!expireAt) {
      alert('Укажите корректную дату и время истечения');
      return;
    }

    const submitData = {
      ...formData,
      expire_at: expireAt,
      traffic_limit_bytes: Number(formData.traffic_limit_bytes) || 0,
      hwid_device_limit: Number(formData.hwid_device_limit) || 0,
      last_triggered_threshold: Number(formData.last_triggered_threshold) || 0
    };

    if (submitData.telegram_id === '') {
      submitData.telegram_id = null;
    } else {
      submitData.telegram_id = Number(submitData.telegram_id) || null;
    }
    if (submitData.description === '') {
      submitData.description = null;
    }
    if (submitData.tag === '') {
      submitData.tag = null;
    }
    if (submitData.email === '') {
      submitData.email = null;
    }

    try {
      setSavingUser(true);
      await onSave(submitData, {
        selectedInternalSquadUuids: Array.from(selectedInternalSquads),
      });
    } finally {
      setSavingUser(false);
    }
  };

  return (
    <>
      <div className="modal-overlay user-create-overlay" onClick={onClose}>
        <div className="modal modal-large user-create-modal" onClick={(event) => event.stopPropagation()}>
          <header className="modal-header user-create-header">
            <div className="user-create-title-group">
              <div className="user-create-title-icon" aria-hidden="true">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M8 7a4 4 0 1 0 8 0a4 4 0 0 0 -8 0" />
                  <path d="M6 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2" />
                </svg>
              </div>
              <div className="user-create-title-stack">
                <h3 className="modal-title">{user ? 'Редактировать пользователя' : 'Создать пользователя'}</h3>
                <p className="user-create-title-subtitle">
                  {user ? 'Измените параметры и доступ пользователя' : 'Заполните данные для нового пользователя'}
                </p>
              </div>
            </div>
            <button type="button" className="modal-close" onClick={onClose} aria-label="Закрыть">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            </button>
          </header>

          <form className="user-create-form" onSubmit={handleSubmit}>
            <div className="modal-body modal-body-columns user-create-body">
              <div className="user-create-column">
                <section className="form-section-card user-create-card">
                  <div className="form-section-header user-create-card-header">
                    <span className="user-create-card-icon blue" aria-hidden="true">
                      <svg viewBox="0 0 256 256" fill="currentColor">
                        <path d="M192,96a64,64,0,1,1-64-64A64,64,0,0,1,192,96Z" opacity="0.2"></path>
                        <path d="M230.92,212c-15.23-26.33-38.7-45.21-66.09-54.16a72,72,0,1,0-73.66,0C63.78,166.78,40.31,185.66,25.08,212a8,8,0,1,0,13.85,8c18.84-32.56,52.14-52,89.07-52s70.23,19.44,89.07,52a8,8,0,1,0,13.85-8ZM72,96a56,56,0,1,1,56,56A56.06,56.06,0,0,1,72,96Z"></path>
                      </svg>
                    </span>
                    <h4 className="form-section-title">Личность</h4>
                  </div>
                  <div className="form-section-content single-column user-create-card-content">
                    <div className="form-group form-group-last">
                      <label className="form-label">Имя пользователя *</label>
                      <input
                        type="text"
                        className="form-input monospace-field"
                        value={formData.username}
                        onChange={(event) => handleChange('username', event.target.value)}
                        placeholder="Введите имя пользователя"
                        required
                        disabled={Boolean(user)}
                      />
                      <div className="form-hint">Имя пользователя нельзя изменить после создания</div>
                    </div>
                  </div>
                </section>

                <section className="form-section-card user-create-card">
                  <div className="form-section-header user-create-card-header">
                    <span className="user-create-card-icon teal" aria-hidden="true">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <path d="M3 7a2 2 0 0 1 2 -2h14a2 2 0 0 1 2 2v10a2 2 0 0 1 -2 2h-14a2 2 0 0 1 -2 -2v-10z"></path>
                        <path d="M3 7l9 6l9 -6"></path>
                      </svg>
                    </span>
                    <h4 className="form-section-title">Contact Information</h4>
                  </div>
                  <div className="form-section-content single-column user-create-card-content">
                    <div className="form-group">
                      <label className="form-label">Telegram ID</label>
                      <input
                        type="number"
                        className="form-input"
                        value={formData.telegram_id}
                        onChange={(event) => handleChange('telegram_id', event.target.value)}
                        placeholder="Enter user's Telegram ID (optional)"
                      />
                    </div>
                    <div className="form-group form-group-last">
                      <label className="form-label">Email</label>
                      <input
                        type="email"
                        className="form-input"
                        value={formData.email}
                        onChange={(event) => handleChange('email', event.target.value)}
                        placeholder="Enter user's email (optional)"
                      />
                    </div>
                  </div>
                </section>

                <section className="form-section-card user-create-card">
                  <div className="form-section-header user-create-card-header">
                    <span className="user-create-card-icon orange" aria-hidden="true">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <path d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z"></path>
                        <path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0"></path>
                      </svg>
                    </span>
                    <h4 className="form-section-title">Device & Tag Settings</h4>
                  </div>
                  <div className="form-section-content single-column user-create-card-content">
                    <div className="form-group">
                      <label className="form-label">Ограничение HWID устройств</label>
                      <input
                        type="number"
                        className="form-input"
                        value={formData.hwid_device_limit}
                        onChange={(event) => handleChange('hwid_device_limit', parseInt(event.target.value, 10) || 0)}
                        min="0"
                        placeholder="Fallback Device Limit in use"
                      />
                    </div>

                    <div className="form-group">
                      <label className="form-label">Tag</label>
                      <input
                        type="text"
                        className="form-input"
                        value={formData.tag}
                        onChange={(event) => handleChange('tag', event.target.value)}
                        placeholder="EXAMPLE_TAG_1"
                      />
                    </div>

                    <div className="form-group form-group-last">
                      <label className="form-label">Описание</label>
                      <textarea
                        className="form-textarea"
                        value={formData.description}
                        onChange={(event) => handleChange('description', event.target.value)}
                        placeholder="Описание пользователя"
                      />
                    </div>
                  </div>
                </section>
              </div>

              <div className="user-create-column">
                <section className="form-section-card user-create-card">
                  <div className="form-section-header user-create-card-header">
                    <span className="user-create-card-icon violet" aria-hidden="true">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <path d="M4 19l16 0"></path>
                        <path d="M4 15l4 -6l4 2l4 -5l4 4"></path>
                      </svg>
                    </span>
                    <h4 className="form-section-title">Traffic & Limits</h4>
                  </div>
                  <div className="form-section-content single-column user-create-card-content">
                    <div className="form-group">
                      <label className="form-label">Лимит трафика</label>
                      <input
                        type="number"
                        className="form-input"
                        value={formData.traffic_limit_bytes}
                        onChange={(event) => handleChange('traffic_limit_bytes', parseInt(event.target.value, 10) || 0)}
                        min="0"
                      />
                      <div className="form-hint">Введите лимит в байтах, 0 — безлимит</div>
                    </div>

                    <div className="form-group">
                      <label className="form-label">Стратегия сброса трафика</label>
                      <AppSelect
                        value={formData.traffic_limit_strategy}
                        onChange={(event) => handleChange('traffic_limit_strategy', event.target.value)}
                      >
                        {TRAFFIC_STRATEGY_OPTIONS.map((strategy) => (
                          <option key={strategy} value={strategy}>
                            {TRAFFIC_STRATEGY_LABELS[strategy] || strategy}
                          </option>
                        ))}
                      </AppSelect>
                    </div>

                    <div className="form-group">
                      <label className="form-label">Статус пользователя</label>
                      <AppSelect
                        value={formData.status}
                        onChange={(event) => handleChange('status', event.target.value)}
                      >
                        {STATUS_OPTIONS.map((status) => (
                          <option key={status} value={status}>
                            {status}
                          </option>
                        ))}
                      </AppSelect>
                    </div>

                    <div className="form-group form-group-last">
                      <label className="form-label">Порог уведомления (%)</label>
                      <input
                        type="number"
                        className="form-input"
                        value={formData.last_triggered_threshold}
                        onChange={(event) => handleChange('last_triggered_threshold', parseInt(event.target.value, 10) || 0)}
                        min="0"
                        max="100"
                      />
                    </div>
                  </div>
                </section>

                <section className="form-section-card user-create-card">
                  <div className="form-section-header user-create-card-header">
                    <span className="user-create-card-icon indigo" aria-hidden="true">
                      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                        <path d="M12 3a12 12 0 0 0 8.5 3a12 12 0 0 1 -8.5 15a12 12 0 0 1 -8.5 -15a12 12 0 0 0 8.5 -3"></path>
                      </svg>
                    </span>
                    <h4 className="form-section-title">Настройки доступа</h4>
                  </div>
                  <div className="form-section-content single-column user-create-card-content">
                    <div className="form-group">
                      <label className="form-label">Дата истечения подписки *</label>
                      <div className="user-create-inline-row">
                        <input
                          type="date"
                          className="form-input"
                          value={expireDate}
                          onChange={(event) => setExpireDate(event.target.value)}
                          required
                        />
                        <input
                          type="time"
                          className="form-input"
                          value={expireTime}
                          onChange={(event) => setExpireTime(event.target.value)}
                          required
                        />
                      </div>
                    </div>

                    <div className="form-group form-group-last">
                      <label className="form-label">Внутренние сквады</label>
                      <div className="user-squad-selector">
                        <div className="search-box search-box-compact squad-editor-search user-squad-drawer-search">
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                            <circle cx="11" cy="11" r="8"></circle>
                            <line x1="21" y1="21" x2="16.65" y2="16.65"></line>
                          </svg>
                          <input
                            type="text"
                            placeholder="Поиск по внутренним сквадам..."
                            value={squadSearch}
                            onChange={(event) => setSquadSearch(event.target.value)}
                          />
                        </div>

                        <div className="user-squad-selector-list">
                          {loadingSquads || loadingUserSquads ? (
                            <div className="squad-editor-empty">Загрузка сквадов...</div>
                          ) : filteredSquads.length === 0 ? (
                            <div className="squad-editor-empty">По вашему запросу сквады не найдены.</div>
                          ) : (
                            filteredSquads.map((squad) => {
                              const isSelected = selectedInternalSquads.has(squad.uuid);
                              const inboundsCount = squad.inbounds_count ?? 0;
                              const membersCount = squad.members_count ?? 0;

                              return (
                                <div key={squad.uuid} className={`user-squad-selector-row ${isSelected ? 'selected' : ''}`}>
                                  <button
                                    type="button"
                                    className="user-squad-selector-main"
                                    onClick={() => toggleInternalSquad(squad.uuid)}
                                    aria-pressed={isSelected}
                                  >
                                    <span className={`user-squad-selector-check ${isSelected ? 'checked' : ''}`} aria-hidden="true">
                                      {isSelected ? '✓' : ''}
                                    </span>
                                    <span className="user-squad-selector-content">
                                      <span className="user-squad-selector-title">{squad.name || 'Unnamed squad'}</span>
                                      <span className="user-squad-option-badges">
                                        <span className="user-squad-option-badge">{inboundsCount} инбаундов</span>
                                        <span className="user-squad-option-badge muted">{membersCount} участников</span>
                                      </span>
                                    </span>
                                  </button>
                                  <button
                                    type="button"
                                    className="user-squad-selector-edit"
                                    onClick={() => handleOpenSquadEditor(squad)}
                                    title="Редактировать сквад"
                                    aria-label={`Редактировать сквад ${squad.name || ''}`}
                                  >
                                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                                      <path d="M7 7h-1a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-1"></path>
                                      <path d="M20.385 6.585a2.1 2.1 0 0 0 -2.97 -2.97l-8.415 8.385v3h3l8.385 -8.415z"></path>
                                      <path d="M16 5l3 3"></path>
                                    </svg>
                                  </button>
                                </div>
                              );
                            })
                          )}
                        </div>

                        <div className="form-hint">
                          Выбрано сквадов: {selectedInternalSquads.size}. Можно выбрать все, один или ни одного.
                        </div>
                      </div>
                    </div>
                  </div>
                </section>
              </div>
            </div>

            <footer className="modal-footer user-create-footer">
              <button type="button" className="squad-editor-summary-btn cancel user-create-btn-cancel" onClick={onClose} disabled={savingUser}>
                Отмена
              </button>
              <button type="submit" className="squad-editor-summary-btn save user-create-btn-save" disabled={savingUser}>
                <svg viewBox="0 0 256 256" fill="currentColor" aria-hidden="true">
                  <path d="M216,83.31V208a8,8,0,0,1-8,8H176V152a8,8,0,0,0-8-8H88a8,8,0,0,0-8,8v64H48a8,8,0,0,1-8-8V48a8,8,0,0,1,8-8H172.69a8,8,0,0,1,5.65,2.34l35.32,35.32A8,8,0,0,1,216,83.31Z" opacity="0.2"></path>
                  <path d="M219.31,72,184,36.69A15.86,15.86,0,0,0,172.69,32H48A16,16,0,0,0,32,48V208a16,16,0,0,0,16,16H208a16,16,0,0,0,16-16V83.31A15.86,15.86,0,0,0,219.31,72ZM168,208H88V152h80Zm40,0H184V152a16,16,0,0,0-16-16H88a16,16,0,0,0-16,16v56H48V48H172.69L208,83.31ZM160,72a8,8,0,0,1-8,8H96a8,8,0,0,1,0-16h56A8,8,0,0,1,160,72Z"></path>
                </svg>
                <span>{savingUser ? 'Сохранение...' : user ? 'Сохранить' : 'Создать'}</span>
              </button>
            </footer>
          </form>
        </div>
      </div>

      <SideDrawerPanel
        open={showInboundsDrawer}
        onClose={closeInboundsDrawer}
        width="56rem"
        title="Изменить сквад"
        icon={(
          <svg stroke="currentColor" fill="none" strokeWidth="2" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round" xmlns="http://www.w3.org/2000/svg">
            <path d="M9.183 6.117a6 6 0 1 0 4.511 3.986"></path>
            <path d="M14.813 17.883a6 6 0 1 0 -4.496 -3.954"></path>
          </svg>
        )}
      >
        <div className="squad-editor-stack">
          <div className="squad-editor-summary">
            <div className="squad-editor-summary-header">
              <div className="squad-editor-summary-identity">
                <button type="button" className="squad-editor-summary-icon" tabIndex={-1} aria-hidden="true">
                  <svg stroke="currentColor" fill="none" strokeWidth="2" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round" xmlns="http://www.w3.org/2000/svg">
                    <path d="M9.183 6.117a6 6 0 1 0 4.511 3.986"></path>
                    <path d="M14.813 17.883a6 6 0 1 0 -4.496 -3.954"></path>
                  </svg>
                </button>
                <div className="squad-editor-summary-title-wrap">
                  <h3 className="squad-editor-summary-title" title={selectedSquad?.name || ''}>
                    {selectedSquad?.name || 'Squad'}
                  </h3>
                </div>
              </div>
              <div className="squad-editor-summary-stats">
                <span className="squad-editor-pill">Выбрано: {selectedInbounds.size}</span>
                <span className="squad-editor-pill muted">Всего: {allInbounds.length}</span>
              </div>
            </div>
            <div className="squad-editor-summary-actions">
              <button
                type="button"
                className="squad-editor-summary-btn cancel"
                onClick={closeInboundsDrawer}
                disabled={savingInbounds}
              >
                Отмена
              </button>
              <button
                type="button"
                className="squad-editor-summary-btn save"
                onClick={handleSaveInbounds}
                disabled={!selectedSquad || savingInbounds}
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                  <path d="M6 4h10l4 4v10a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12a2 2 0 0 1 2 -2"></path>
                  <path d="M12 14m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0"></path>
                  <path d="M14 4l0 4l-6 0l0 -4"></path>
                </svg>
                <span>{savingInbounds ? 'Сохранение...' : 'Сохранить'}</span>
              </button>
            </div>
          </div>

          <div className="search-box search-box-compact squad-editor-search">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
              <circle cx="11" cy="11" r="8" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
            <input
              type="text"
              value={inboundsSearchQuery}
              onChange={(event) => setInboundsSearchQuery(event.target.value)}
              placeholder="Поиск по профилям или инбаундам..."
            />
          </div>

          <div className="squad-editor-tabs" role="tablist" aria-label="Inbounds view mode">
            <button
              type="button"
              className={`squad-editor-tab ${inboundsTab === INBOUNDS_TAB.profiles ? 'active' : ''}`}
              onClick={() => setInboundsTab(INBOUNDS_TAB.profiles)}
              role="tab"
              aria-selected={inboundsTab === INBOUNDS_TAB.profiles}
            >
              Профили
            </button>
            <button
              type="button"
              className={`squad-editor-tab ${inboundsTab === INBOUNDS_TAB.list ? 'active' : ''}`}
              onClick={() => setInboundsTab(INBOUNDS_TAB.list)}
              role="tab"
              aria-selected={inboundsTab === INBOUNDS_TAB.list}
            >
              Список
            </button>
          </div>

          {loadingInboundsDrawer ? (
            <div className="squad-editor-empty">Загрузка данных сквада...</div>
          ) : inboundsTab === INBOUNDS_TAB.profiles ? (
            <div className="squad-editor-list">
              {profileGroups.length === 0 ? (
                <div className="squad-editor-empty">Профили или инбаунды не найдены.</div>
              ) : (
                profileGroups.map((profile) => (
                  <section
                    key={profile.uuid}
                    className={`squad-editor-group ${expandedProfiles.has(profile.uuid) ? 'expanded' : ''}`}
                  >
                    <div className="squad-editor-group-row">
                      <button
                        type="button"
                        className="squad-editor-group-toggle"
                        aria-expanded={expandedProfiles.has(profile.uuid)}
                        onClick={() => toggleProfileAccordion(profile.uuid)}
                      >
                        <span className={`squad-editor-group-chevron ${expandedProfiles.has(profile.uuid) ? 'open' : ''}`}>
                          <svg viewBox="0 0 15 15" fill="currentColor" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                            <path
                              fillRule="evenodd"
                              clipRule="evenodd"
                              d="M3.13523 6.15803C3.3241 5.95657 3.64052 5.94637 3.84197 6.13523L7.5 9.56464L11.158 6.13523C11.3595 5.94637 11.6759 5.95657 11.8648 6.15803C12.0536 6.35949 12.0434 6.67591 11.842 6.86477L7.84197 10.6148C7.64964 10.7951 7.35036 10.7951 7.15803 10.6148L3.15803 6.86477C2.95657 6.67591 2.94637 6.35949 3.13523 6.15803Z"
                            />
                          </svg>
                        </span>
                        <span className="squad-editor-group-toggle-content">
                          <span className="squad-editor-group-title">{profile.name}</span>
                          <span className="squad-editor-group-badges">
                            <span className="squad-editor-badge">
                              {profile.selectedCount} / {profile.profileInbounds.length}
                            </span>
                            <span className="squad-editor-badge muted">{profile.visibleInbounds.length}</span>
                          </span>
                        </span>
                      </button>
                      <div className="squad-editor-group-actions">
                        <button
                          type="button"
                          className="squad-editor-group-icon-btn"
                          onClick={(event) => {
                            event.stopPropagation();
                            setProfileSelection(profile.profileInbounds, true);
                          }}
                          disabled={profile.profileInbounds.length === 0}
                          aria-label="Выбрать все инбаунды профиля"
                          title="Выбрать все"
                        >
                          <svg viewBox="0 0 256 256" fill="currentColor" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                            <path d="M232.49,80.49l-128,128a12,12,0,0,1-17,0l-56-56a12,12,0,1,1,17-17L96,183,215.51,63.51a12,12,0,0,1,17,17Z" />
                          </svg>
                        </button>
                        <button
                          type="button"
                          className="squad-editor-group-icon-btn"
                          onClick={(event) => {
                            event.stopPropagation();
                            setProfileSelection(profile.profileInbounds, false);
                          }}
                          disabled={profile.profileInbounds.length === 0}
                          aria-label="Снять все инбаунды профиля"
                          title="Снять все"
                        >
                          <svg viewBox="0 0 256 256" fill="currentColor" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                            <path d="M208.49,191.51a12,12,0,0,1-17,17L128,145,64.49,208.49a12,12,0,0,1-17-17L111,128,47.51,64.49a12,12,0,0,1,17-17L128,111l63.51-63.52a12,12,0,0,1,17,17L145,128Z" />
                          </svg>
                        </button>
                      </div>
                    </div>

                    {expandedProfiles.has(profile.uuid) ? (
                      <div className="squad-editor-group-panel">
                        {profile.profileInbounds.length === 0 ? (
                          <p className="squad-editor-empty small">В этом профиле нет инбаундов.</p>
                        ) : profile.visibleInbounds.length === 0 ? (
                          <p className="squad-editor-empty small">По запросу ничего не найдено.</p>
                        ) : (
                          <div className="squad-editor-items">
                            {profile.visibleInbounds.map((inbound) => (
                              <label key={inbound.uuid} className="squad-editor-checkbox">
                                <input
                                  type="checkbox"
                                  checked={selectedInbounds.has(inbound.uuid)}
                                  onChange={() => toggleInbound(inbound.uuid)}
                                />
                                <span className="squad-editor-checkbox-body">
                                  <span className="squad-editor-inbound-title">{inbound.tag || 'untagged'}</span>
                                  <span className="squad-editor-inbound-badges">
                                    <span className="squad-editor-inline-pill protocol">{getInboundProtocol(inbound)}</span>
                                    <span className="squad-editor-inline-pill port">{String(getInboundPort(inbound))}</span>
                                    <span className="squad-editor-inline-pill security">{getInboundSecurity(inbound)}</span>
                                  </span>
                                </span>
                              </label>
                            ))}
                          </div>
                        )}
                      </div>
                    ) : null}
                  </section>
                ))
              )}
            </div>
          ) : (
            <>
              <div className="squad-editor-flat-filters" role="radiogroup" aria-label="Inbound list filters">
                <button
                  type="button"
                  className={`squad-editor-filter ${flatFilter === FLAT_FILTER.all ? 'active' : ''}`}
                  onClick={() => setFlatFilter(FLAT_FILTER.all)}
                >
                  Все
                </button>
                <button
                  type="button"
                  className={`squad-editor-filter ${flatFilter === FLAT_FILTER.selected ? 'active' : ''}`}
                  onClick={() => setFlatFilter(FLAT_FILTER.selected)}
                >
                  Выбранные
                </button>
                <button
                  type="button"
                  className={`squad-editor-filter ${flatFilter === FLAT_FILTER.unselected ? 'active' : ''}`}
                  onClick={() => setFlatFilter(FLAT_FILTER.unselected)}
                >
                  Невыбранные
                </button>
              </div>

              <div className="squad-editor-list">
                {flatInbounds.length === 0 ? (
                  <div className="squad-editor-empty">Список инбаундов пуст для текущего фильтра.</div>
                ) : (
                  <div className="squad-editor-items flat">
                    {flatInbounds.map((inbound) => (
                      <label key={inbound.uuid} className="squad-editor-checkbox squad-editor-checkbox-flat">
                        <input
                          type="checkbox"
                          checked={selectedInbounds.has(inbound.uuid)}
                          onChange={() => toggleInbound(inbound.uuid)}
                        />
                        <span className="squad-editor-checkbox-body">
                          <span className="squad-editor-inbound-main">
                            <span className="squad-editor-inbound-title">{inbound.tag || 'untagged'}</span>
                            <span className="squad-editor-inbound-subtitle">{inbound.profileName || 'profile'}</span>
                          </span>
                          <span className="squad-editor-inbound-badges">
                            <span className="squad-editor-inline-pill protocol">{getInboundProtocol(inbound)}</span>
                            <span className="squad-editor-inline-pill port">{String(getInboundPort(inbound))}</span>
                            <span className="squad-editor-inline-pill security">{getInboundSecurity(inbound)}</span>
                          </span>
                        </span>
                      </label>
                    ))}
                  </div>
                )}
              </div>
            </>
          )}
        </div>
      </SideDrawerPanel>
    </>
  );
}

export default UserModal;
