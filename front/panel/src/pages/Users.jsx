import { useState, useEffect, useRef, useMemo } from 'react';
import UserModal from '../components/UserModal';
import { squadsApi, usersApi } from '../api';
import AppSelect from '../components/AppSelect';
import SharedRichTable from '../components/SharedRichTable';

const STATUS_ORDER = {
  ACTIVE: 0,
  LIMITED: 1,
  DISABLED: 2,
  EXPIRED: 3,
};

const PAGE_SIZE_OPTIONS = [5, 10, 15, 20, 25, 30, 50, 100];

const USER_COLUMNS = [
  {
    key: 'select',
    label: 'Select',
    sortable: false,
    alwaysVisible: true,
    defaultVisible: true,
    defaultWidth: 64,
    minWidth: 56,
    defaultPin: 'left',
    disablePinning: true,
  },
  {
    key: 'username',
    label: 'Имя пользователя',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 270,
    minWidth: 210,
    defaultPin: 'left',
  },
  {
    key: 'id',
    label: 'ID',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 96,
    minWidth: 80,
    defaultPin: 'left',
  },
  {
    key: 'status',
    label: 'Статус',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 148,
    minWidth: 130,
  },
  {
    key: 'lastConnection',
    label: 'Последнее подключение',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 198,
    minWidth: 160,
  },
  {
    key: 'expireAt',
    label: 'Истекает',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 176,
    minWidth: 148,
  },
  {
    key: 'traffic',
    label: 'Расход трафика',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 342,
    minWidth: 280,
  },
  {
    key: 'shortUuid',
    label: 'Ссылка-подписка',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 224,
    minWidth: 188,
  },
  {
    key: 'description',
    label: 'Описание',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 220,
    minWidth: 170,
  },
  {
    key: 'telegramId',
    label: 'Telegram ID',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 150,
    minWidth: 125,
  },
  {
    key: 'tag',
    label: 'Tag',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 132,
    minWidth: 100,
  },
  {
    key: 'internalSquad',
    label: 'Внутренние сквады',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 220,
    minWidth: 170,
  },
  {
    key: 'externalSquad',
    label: 'Внешние сквады',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 210,
    minWidth: 170,
  },
  {
    key: 'email',
    label: 'Email',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 220,
    minWidth: 160,
  },
  {
    key: 'firstConnection',
    label: 'Первое подключение',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 220,
    minWidth: 176,
  },
  {
    key: 'lastTrafficReset',
    label: 'Сброс трафика',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 176,
    minWidth: 150,
  },
  {
    key: 'onlineAt',
    label: 'Был в сети',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 176,
    minWidth: 150,
  },
  {
    key: 'lastUserAgent',
    label: 'Последний UA',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 270,
    minWidth: 200,
  },
  {
    key: 'lastSubscriptionRequest',
    label: 'Последний запрос подписки',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 250,
    minWidth: 200,
  },
  {
    key: 'lifetimeTraffic',
    label: 'Трафик за все время',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 192,
    minWidth: 160,
  },
  {
    key: 'subRevokedAt',
    label: 'Сброс ссылки-подписки',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 220,
    minWidth: 180,
  },
  {
    key: 'trafficStrategy',
    label: 'Стратегия лимита',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 174,
    minWidth: 140,
  },
  {
    key: 'hwidLimit',
    label: 'Лимит устройств',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 148,
    minWidth: 128,
  },
  {
    key: 'threshold',
    label: 'Порог триггера',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 150,
    minWidth: 130,
  },
  {
    key: 'createdAt',
    label: 'Дата создания',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 188,
    minWidth: 160,
  },
  {
    key: 'updatedAt',
    label: 'Дата обновления',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 188,
    minWidth: 160,
  },
  {
    key: 'uuid',
    label: 'UUID',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 300,
    minWidth: 220,
  },
];

function Users() {
  const [users, setUsers] = useState([]);
  const [squadNamesByUuid, setSquadNamesByUuid] = useState({});
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [refreshInterval, setRefreshInterval] = useState(30);
  const [selectedUserUUIDs, setSelectedUserUUIDs] = useState(new Set());
  const [compactRows, setCompactRows] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [sortState, setSortState] = useState({ key: null, direction: 'asc' });

  const searchInputRef = useRef(null);

  useEffect(() => {
    fetchUsers();
  }, []);

  useEffect(() => {
    const interval = setInterval(fetchUsers, refreshInterval * 1000);
    return () => clearInterval(interval);
  }, [refreshInterval]);

  useEffect(() => {
    document.body.classList.toggle('users-fullscreen-open', isFullscreen);
    return () => document.body.classList.remove('users-fullscreen-open');
  }, [isFullscreen]);

  useEffect(() => {
    if (!isFullscreen) return undefined;
    const handleEscape = (event) => {
      if (event.key === 'Escape') {
        setIsFullscreen(false);
      }
    };
    window.addEventListener('keydown', handleEscape);
    return () => window.removeEventListener('keydown', handleEscape);
  }, [isFullscreen]);

  const fetchUsers = async () => {
    try {
      const [usersResult, squadsResult] = await Promise.allSettled([
        usersApi.getAll(),
        squadsApi.getAllSummary(),
      ]);

      if (usersResult.status !== 'fulfilled') {
        throw usersResult.reason;
      }

      const usersList = Array.isArray(usersResult.value?.users) ? usersResult.value.users : [];
      setUsers(usersList);
      setSelectedUserUUIDs((prev) => {
        const valid = new Set(usersList.map((user) => user.uuid));
        return new Set(Array.from(prev).filter((uuid) => valid.has(uuid)));
      });

      if (squadsResult.status === 'fulfilled') {
        const nextMap = {};
        const squads = Array.isArray(squadsResult.value?.squads) ? squadsResult.value.squads : [];
        squads.forEach((squad) => {
          if (squad?.uuid) {
            nextMap[squad.uuid] = squad.name || squad.uuid;
          }
        });
        setSquadNamesByUuid(nextMap);
      }

      setLoading(false);
    } catch (err) {
      console.error('Error fetching users:', err);
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setEditingUser(null);
    setModalOpen(true);
  };

  const handleEdit = (user) => {
    setEditingUser(user);
    setModalOpen(true);
  };

  const extractUserIdentity = (payload) => {
    const candidates = [payload, payload?.user, payload?.data, payload?.result].filter(Boolean);
    for (const candidate of candidates) {
      const userId = candidate?.t_id ?? candidate?.user_id ?? candidate?.id;
      const userUuid = candidate?.uuid ?? null;
      if (userId !== undefined && userId !== null) {
        return { userId, userUuid };
      }
      if (userUuid) {
        return { userId: null, userUuid };
      }
    }
    return { userId: null, userUuid: null };
  };

  const resolveSavedUserIdentity = async (responseData, fallbackUser, username) => {
    const fallbackId = fallbackUser?.t_id ?? fallbackUser?.user_id ?? null;
    const fallbackUuid = fallbackUser?.uuid ?? null;
    if (fallbackId !== null && fallbackId !== undefined) {
      return { userId: fallbackId, userUuid: fallbackUuid };
    }

    const fromResponse = extractUserIdentity(responseData);
    if (fromResponse.userId !== null && fromResponse.userId !== undefined) {
      return fromResponse;
    }

    if (fromResponse.userUuid) {
      try {
        const userDetails = await usersApi.getById(fromResponse.userUuid);
        const details = extractUserIdentity(userDetails);
        if (details.userId !== null && details.userId !== undefined) {
          return details;
        }
      } catch (err) {
        console.warn('Failed to resolve user id by uuid:', err);
      }
    }

    const normalizedUsername = String(username || '').trim();
    if (!normalizedUsername) {
      return { userId: null, userUuid: fromResponse.userUuid || fallbackUuid };
    }

    const allUsersResponse = await usersApi.getAll();
    const usersList = Array.isArray(allUsersResponse?.users) ? allUsersResponse.users : [];
    const exactMatches = usersList.filter((item) => item?.username === normalizedUsername);
    const userCandidate = exactMatches[exactMatches.length - 1] || null;

    if (!userCandidate) {
      return { userId: null, userUuid: fromResponse.userUuid || fallbackUuid };
    }

    const userId = userCandidate.t_id ?? userCandidate.user_id ?? null;
    return {
      userId,
      userUuid: userCandidate.uuid ?? fromResponse.userUuid ?? fallbackUuid,
    };
  };

  const syncInternalSquadMemberships = async (userId, selectedInternalSquadUuids) => {
    const squadsResponse = await squadsApi.getAllSummary();
    const squadsList = Array.isArray(squadsResponse?.squads) ? squadsResponse.squads : [];
    const selectedSet = new Set((selectedInternalSquadUuids || []).filter(Boolean).map((value) => String(value)));

    const numericUserId = Number(userId);
    const targetUserId = Number.isFinite(numericUserId) ? numericUserId : userId;
    const targetUserKey = String(targetUserId);

    for (const squad of squadsList) {
      const membersResponse = await squadsApi.getMembers(squad.uuid);
      const currentIds = Array.isArray(membersResponse?.squad_members)
        ? membersResponse.squad_members
          .map((member) => member?.user_id ?? member?.t_id)
          .filter((value) => value !== undefined && value !== null)
        : [];

      const hasUser = currentIds.some((value) => String(value) === targetUserKey);
      const shouldHaveUser = selectedSet.has(String(squad.uuid));

      if (hasUser === shouldHaveUser) {
        continue;
      }

      const nextIds = shouldHaveUser
        ? Array.from(new Set([...currentIds, targetUserId]))
        : currentIds.filter((value) => String(value) !== targetUserKey);

      await squadsApi.setMembers(squad.uuid, nextIds);
    }
  };

  const handleSave = async (userData, options = {}) => {
    const { selectedInternalSquadUuids } = options;
    try {
      let saveResponse;
      if (editingUser) {
        saveResponse = await usersApi.update(editingUser.uuid, userData);
      } else {
        saveResponse = await usersApi.create(userData);
      }

      if (Array.isArray(selectedInternalSquadUuids)) {
        try {
          const { userId } = await resolveSavedUserIdentity(saveResponse, editingUser, userData.username);
          if (userId === null || userId === undefined) {
            throw new Error('Не удалось определить идентификатор пользователя');
          }
          await syncInternalSquadMemberships(userId, selectedInternalSquadUuids);
        } catch (membershipErr) {
          console.error('Error syncing user internal squads:', membershipErr);
          alert(`Пользователь сохранен, но привязка к внутренним сквадам не выполнена: ${membershipErr.message}`);
        }
      }

      setModalOpen(false);
      await fetchUsers();
    } catch (err) {
      console.error('Error saving user:', err);
      alert(err.message || 'Failed to save user');
      throw err;
    }
  };

  const handleDeleteSelected = async () => {
    if (selectedUserUUIDs.size === 0) return;
    if (!confirm(`Удалить выбранных пользователей: ${selectedUserUUIDs.size}?`)) return;

    try {
      await usersApi.deleteMany(Array.from(selectedUserUUIDs));
      await fetchUsers();
    } catch (err) {
      console.error('Error deleting selected users:', err);
      alert(`Failed to delete selected users: ${err.message}`);
    }
  };

  const formatBytes = (bytes) => {
    const normalized = Number(bytes);
    if (!Number.isFinite(normalized) || normalized <= 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.min(Math.floor(Math.log(normalized) / Math.log(k)), sizes.length - 1);
    return `${parseFloat((normalized / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
  };

  const normalizeStatus = (status) => String(status || '').toUpperCase();

  const formatDateTime = (value) => {
    if (!value) return '–';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '–';
    return date.toLocaleString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const formatDuration = (milliseconds) => {
    const abs = Math.max(0, Math.abs(milliseconds));
    const minute = 60 * 1000;
    const hour = 60 * minute;
    const day = 24 * hour;

    if (abs >= day) {
      const value = Math.round(abs / day);
      return `${value} дн.`;
    }
    if (abs >= hour) {
      const value = Math.round(abs / hour);
      return `${value} ч.`;
    }
    const value = Math.max(1, Math.round(abs / minute));
    return `${value} мин.`;
  };

  const formatRelativeTime = (value) => {
    if (!value) return 'Никогда';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return 'Никогда';

    const diff = date.getTime() - Date.now();
    const duration = formatDuration(diff);
    if (diff >= 0) {
      return `через ${duration}`;
    }
    return `${duration} назад`;
  };

  const formatExpireAt = (value) => {
    if (!value) return '–';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '–';

    const diff = date.getTime() - Date.now();
    const duration = formatDuration(diff);
    if (diff >= 0) {
      return `Через ${duration}`;
    }
    return `Истек ${duration} назад`;
  };

  const getInternalSquadLabel = (user) => {
    if (!user.internal_squad_uuid) return null;
    return squadNamesByUuid[user.internal_squad_uuid] || user.internal_squad_uuid.slice(0, 8);
  };

  const getExternalSquadLabel = (user) => {
    if (!user.external_squad_uuid) return null;
    return squadNamesByUuid[user.external_squad_uuid] || user.external_squad_uuid.slice(0, 8);
  };

  const getLastConnectionValue = (user) => {
    return user.last_connected_node_name
      || user.node_name
      || user.nodeName
      || user.last_node_name
      || null;
  };

  const getFirstConnectionAt = (user) => {
    return user.user_traffic_first_connected_at
      || user.first_connected_at
      || null;
  };

  const getLastSubscriptionRequestAt = (user) => {
    return user.last_subscription_request_at
      || user.subscription_last_requested_at
      || user.sub_last_opened_at
      || null;
  };

  const getTrafficDetails = (user) => {
    const usedBytes = Math.max(
      0,
      Number(
        user.used_traffic_bytes
          ?? user.user_traffic_lifetime_used_traffic_bytes
          ?? user.lifetime_used_traffic_bytes
          ?? 0,
      ),
    );
    const lifetimeBytes = Math.max(
      0,
      Number(user.lifetime_used_traffic_bytes ?? user.user_traffic_lifetime_used_traffic_bytes ?? 0),
    );
    const limitBytes = Math.max(0, Number(user.traffic_limit_bytes ?? 0));
    const hasLimit = limitBytes > 0;
    const usedPercent = hasLimit ? Math.min(100, (usedBytes / limitBytes) * 100) : 0;
    const freePercent = hasLimit ? Math.max(0, 100 - usedPercent) : 100;

    return {
      usedBytes,
      lifetimeBytes,
      limitBytes,
      hasLimit,
      usedPercent,
      freePercent,
    };
  };

  const getSubscriptionUrl = (shortUuid) => {
    if (!shortUuid) return null;
    return `/api/v1/sub?mode=advanced&user=${encodeURIComponent(String(shortUuid))}`;
  };

  const getSortValue = (user, key) => {
    const internalSquadLabel = getInternalSquadLabel(user) || '';
    const externalSquadLabel = getExternalSquadLabel(user) || '';

    switch (key) {
      case 'username':
        return String(user.username || '').toLowerCase();
      case 'id':
        return Number(user.t_id ?? user.user_id ?? 0);
      case 'status':
        return STATUS_ORDER[normalizeStatus(user.status)] ?? 99;
      case 'lastConnection':
        return String(getLastConnectionValue(user) || '').toLowerCase();
      case 'expireAt':
        return user.expire_at ? new Date(user.expire_at).getTime() : 0;
      case 'traffic':
        return Number(user.traffic_limit_bytes ?? 0);
      case 'shortUuid':
        return String(user.short_uuid || '').toLowerCase();
      case 'description':
        return String(user.description || '').toLowerCase();
      case 'telegramId':
        return Number(user.telegram_id ?? 0);
      case 'tag':
        return String(user.tag || '').toLowerCase();
      case 'internalSquad':
        return internalSquadLabel.toLowerCase();
      case 'externalSquad':
        return externalSquadLabel.toLowerCase();
      case 'email':
        return String(user.email || '').toLowerCase();
      case 'firstConnection':
        return getFirstConnectionAt(user) ? new Date(getFirstConnectionAt(user)).getTime() : 0;
      case 'lastTrafficReset':
        return user.last_traffic_reset_at ? new Date(user.last_traffic_reset_at).getTime() : 0;
      case 'onlineAt':
        return user.sub_last_opened_at ? new Date(user.sub_last_opened_at).getTime() : 0;
      case 'lastUserAgent':
        return String(user.sub_last_user_agent || '').toLowerCase();
      case 'lastSubscriptionRequest':
        return getLastSubscriptionRequestAt(user) ? new Date(getLastSubscriptionRequestAt(user)).getTime() : 0;
      case 'lifetimeTraffic':
        return Number(user.lifetime_used_traffic_bytes ?? user.user_traffic_lifetime_used_traffic_bytes ?? 0);
      case 'subRevokedAt':
        return user.sub_revoked_at ? new Date(user.sub_revoked_at).getTime() : 0;
      case 'trafficStrategy':
        return String(user.traffic_limit_strategy || '').toUpperCase();
      case 'hwidLimit':
        return Number(user.hwid_device_limit ?? 0);
      case 'threshold':
        return Number(user.last_triggered_threshold ?? 0);
      case 'createdAt':
        return user.created_at ? new Date(user.created_at).getTime() : 0;
      case 'updatedAt':
        return user.updated_at ? new Date(user.updated_at).getTime() : 0;
      case 'uuid':
        return String(user.uuid || '').toLowerCase();
      default:
        return '';
    }
  };

  const query = searchQuery.trim().toLowerCase();

  const filteredUsers = useMemo(() => {
    if (!query) {
      return users;
    }

    return users.filter((user) => {
      const pool = [
        user.username,
        user.email,
        user.tag,
        user.status,
        user.short_uuid,
        user.description,
        user.sub_last_user_agent,
        String(user.telegram_id ?? ''),
        String(user.t_id ?? user.user_id ?? ''),
        user.uuid,
      ]
        .filter(Boolean)
        .map((value) => String(value).toLowerCase());

      return pool.some((value) => value.includes(query));
    });
  }, [users, query]);

  const sortedUsers = useMemo(() => {
    if (!sortState.key) {
      return filteredUsers;
    }

    const directionFactor = sortState.direction === 'asc' ? 1 : -1;

    return filteredUsers
      .map((user, index) => ({ user, index }))
      .sort((a, b) => {
        const aValue = getSortValue(a.user, sortState.key);
        const bValue = getSortValue(b.user, sortState.key);

        if (typeof aValue === 'string' || typeof bValue === 'string') {
          const compare = String(aValue).localeCompare(String(bValue), 'ru', { sensitivity: 'base' });
          if (compare !== 0) {
            return compare * directionFactor;
          }
        } else {
          const compare = Number(aValue) - Number(bValue);
          if (compare !== 0) {
            return compare * directionFactor;
          }
        }

        return a.index - b.index;
      })
      .map((item) => item.user);
  }, [filteredUsers, sortState, squadNamesByUuid]);

  const usersStats = useMemo(() => {
    const stats = {
      total: users.length,
      active: 0,
      expired: 0,
      limited: 0,
      disabled: 0,
    };

    users.forEach((user) => {
      const status = normalizeStatus(user.status);
      if (status === 'ACTIVE') stats.active += 1;
      if (status === 'EXPIRED') stats.expired += 1;
      if (status === 'LIMITED') stats.limited += 1;
      if (status === 'DISABLED') stats.disabled += 1;
    });

    return stats;
  }, [users]);

  const openEditFromRow = (user) => {
    handleEdit(user);
  };

  const renderCell = (column, user) => {
    const status = normalizeStatus(user.status);
    const statusKey = status.toLowerCase();
    const internalSquadLabel = getInternalSquadLabel(user);
    const externalSquadLabel = getExternalSquadLabel(user);
    const traffic = getTrafficDetails(user);
    const lifetimeTrafficLabel = traffic.lifetimeBytes > 0 ? formatBytes(traffic.lifetimeBytes) : '–';
    const shortUuid = user.short_uuid || '';
    const subscriptionUrl = getSubscriptionUrl(shortUuid);

    switch (column.key) {
      case 'username':
        return (
          <td key={column.key}>
            <div className="users-name-cell">
              <span className={`users-presence-dot ${user.sub_last_opened_at ? 'is-online' : 'is-idle'}`}></span>
              <div className="users-name-text">
                <strong>{user.username || '—'}</strong>
                <span>{user.sub_last_opened_at ? formatRelativeTime(user.sub_last_opened_at) : 'Не подключался'}</span>
              </div>
            </div>
          </td>
        );
      case 'id':
        return <td key={column.key} className="users-cell-center users-cell-strong">{user.t_id ?? user.user_id ?? '–'}</td>;
      case 'status':
        return (
          <td key={column.key} className="users-cell-center">
            <span className={`users-status-badge users-status-${statusKey}`}>
              <span className="users-status-dot"></span>
              {status || 'UNKNOWN'}
            </span>
          </td>
        );
      case 'lastConnection':
        return <td key={column.key} className="users-cell-center users-cell-dim">{getLastConnectionValue(user) || '–'}</td>;
      case 'expireAt':
        return <td key={column.key} className="users-cell-center users-cell-dim">{formatExpireAt(user.expire_at)}</td>;
      case 'traffic':
        return (
          <td key={column.key}>
            <div className="users-traffic-cell">
              <div className="users-traffic-meta">
                <span className="users-traffic-used-label">{traffic.usedPercent.toFixed(2)}% {traffic.hasLimit ? '' : '∞'}</span>
                <span className="users-traffic-free-label">Σ {formatBytes(traffic.usedBytes)} {traffic.freePercent.toFixed(2)}%</span>
              </div>
              <div className="users-traffic-bar" aria-label="Использование трафика">
                <span className="users-traffic-used" style={{ width: `${traffic.usedPercent}%` }}></span>
                <span className="users-traffic-free" style={{ width: `${traffic.freePercent}%` }}></span>
              </div>
              <div className="users-traffic-foot">
                <span>{formatBytes(traffic.usedBytes)}</span>
                <span>{traffic.hasLimit ? formatBytes(traffic.limitBytes) : '∞'}</span>
              </div>
            </div>
          </td>
        );
      case 'shortUuid':
        return (
          <td key={column.key} className="users-cell-mono users-cell-center">
            {subscriptionUrl ? (
              <a
                href={subscriptionUrl}
                target="_blank"
                rel="noreferrer"
                className="users-sub-link"
                onClick={(event) => event.stopPropagation()}
                title={`Открыть подписку: ${shortUuid}`}
              >
                {shortUuid}
              </a>
            ) : (
              '–'
            )}
          </td>
        );
      case 'description':
        return <td key={column.key} className="users-cell-dim">{user.description || '–'}</td>;
      case 'telegramId':
        return <td key={column.key} className="users-cell-center users-cell-mono">{user.telegram_id ?? '–'}</td>;
      case 'tag':
        return <td key={column.key} className="users-cell-center users-cell-mono">{user.tag || '–'}</td>;
      case 'internalSquad':
        return (
          <td key={column.key}>
            {internalSquadLabel ? (
              <span className="users-squad-pill">{String(internalSquadLabel).toUpperCase()}</span>
            ) : (
              <span className="users-dash">–</span>
            )}
          </td>
        );
      case 'externalSquad':
        return (
          <td key={column.key}>
            {externalSquadLabel ? (
              <span className="users-squad-pill">{String(externalSquadLabel).toUpperCase()}</span>
            ) : (
              <span className="users-dash">–</span>
            )}
          </td>
        );
      case 'email':
        return <td key={column.key}>{user.email || <span className="users-dash">–</span>}</td>;
      case 'firstConnection':
        return <td key={column.key} className="users-cell-center">{formatDateTime(getFirstConnectionAt(user))}</td>;
      case 'lastTrafficReset':
        return <td key={column.key} className="users-cell-center">{user.last_traffic_reset_at ? formatDateTime(user.last_traffic_reset_at) : 'Никогда'}</td>;
      case 'onlineAt':
        return <td key={column.key} className="users-cell-center">{user.sub_last_opened_at ? formatRelativeTime(user.sub_last_opened_at) : 'Никогда'}</td>;
      case 'lastUserAgent':
        return <td key={column.key} className="users-cell-dim">{user.sub_last_user_agent || '–'}</td>;
      case 'lastSubscriptionRequest':
        return <td key={column.key} className="users-cell-center">{formatDateTime(getLastSubscriptionRequestAt(user))}</td>;
      case 'lifetimeTraffic':
        return <td key={column.key} className="users-cell-center users-cell-mono">{lifetimeTrafficLabel}</td>;
      case 'subRevokedAt':
        return <td key={column.key} className="users-cell-center">{user.sub_revoked_at ? formatDateTime(user.sub_revoked_at) : 'Никогда'}</td>;
      case 'trafficStrategy':
        return <td key={column.key} className="users-cell-center users-cell-mono">{user.traffic_limit_strategy || 'NO_RESET'}</td>;
      case 'hwidLimit':
        return <td key={column.key} className="users-cell-center users-cell-mono">{user.hwid_device_limit ?? '–'}</td>;
      case 'threshold':
        return <td key={column.key} className="users-cell-center users-cell-mono">{user.last_triggered_threshold ?? '–'}</td>;
      case 'createdAt':
        return <td key={column.key} className="users-cell-center">{formatDateTime(user.created_at)}</td>;
      case 'updatedAt':
        return <td key={column.key} className="users-cell-center">{formatDateTime(user.updated_at)}</td>;
      case 'uuid':
        return <td key={column.key} className="users-cell-mono users-cell-dim">{user.uuid || '–'}</td>;
      default:
        return <td key={column.key} className="users-cell-center">–</td>;
    }
  };

  if (loading) {
    return (
      <div className="empty-state">
        <div className="spinner"></div>
      </div>
    );
  }

  return (
    <>
      <div className="users-overview-grid">
        <div className="users-overview-card theme-total">
          <div className="users-overview-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2" />
              <circle cx="9" cy="7" r="4" />
              <path d="M23 21v-2a4 4 0 0 0-3-3.87" />
              <path d="M16 3.13a4 4 0 0 1 0 7.75" />
            </svg>
          </div>
          <div className="users-overview-content">
            <span className="users-overview-label">Всего</span>
            <span className="users-overview-value">{usersStats.total}</span>
          </div>
        </div>
        <div className="users-overview-card theme-active">
          <div className="users-overview-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M3 12h4l3-9 4 18 3-9h4" />
            </svg>
          </div>
          <div className="users-overview-content">
            <span className="users-overview-label">Active</span>
            <span className="users-overview-value">{usersStats.active}</span>
          </div>
        </div>
        <div className="users-overview-card theme-expired">
          <div className="users-overview-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="9" />
              <path d="M12 7v6l4 2" />
            </svg>
          </div>
          <div className="users-overview-content">
            <span className="users-overview-label">Expired</span>
            <span className="users-overview-value">{usersStats.expired}</span>
          </div>
        </div>
        <div className="users-overview-card theme-limited">
          <div className="users-overview-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="9" />
              <path d="M12 7v5l3 2" />
            </svg>
          </div>
          <div className="users-overview-content">
            <span className="users-overview-label">Limited</span>
            <span className="users-overview-value">{usersStats.limited}</span>
          </div>
        </div>
        <div className="users-overview-card theme-disabled">
          <div className="users-overview-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="9" />
              <path d="M5 5l14 14" />
            </svg>
          </div>
          <div className="users-overview-content">
            <span className="users-overview-label">Disabled</span>
            <span className="users-overview-value">{usersStats.disabled}</span>
          </div>
        </div>
      </div>

      <div className={`card users-list-card ${isFullscreen ? 'users-list-card-fullscreen' : ''}`}>
        <div className="users-list-header">
          <div className="users-list-header-main">
            <div className="users-action-icon users-action-cyan users-section-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" />
                <circle cx="9" cy="7" r="4" />
                <path d="M23 21v-2a4 4 0 00-3-3.87" />
                <path d="M16 3.13a4 4 0 010 7.75" />
              </svg>
            </div>
            <div className="users-list-heading-stack">
              <h2 className="card-title">Пользователи</h2>
            </div>
          </div>
          <div className="users-list-header-actions">
            <div className="search-box search-box-compact users-list-search">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="11" cy="11" r="8" />
                <line x1="21" y1="21" x2="16.65" y2="16.65" />
              </svg>
              <input
                ref={searchInputRef}
                type="text"
                placeholder="Поиск пользователей..."
                value={searchQuery}
                onChange={(event) => setSearchQuery(event.target.value)}
              />
            </div>
            <button
              type="button"
              className="users-action-icon users-action-help"
              onClick={() => alert('Кликните по заголовку колонки для сортировки. Потяните разделитель справа у заголовка, чтобы изменить ширину колонки.')}
              title="Подсказка"
              aria-label="Подсказка"
            >
              ?
            </button>
            <button
              type="button"
              className="users-action-icon users-action-neutral"
              onClick={() => setSearchQuery('')}
              title="Сбросить фильтр"
              aria-label="Сбросить фильтр"
              disabled={!searchQuery.trim()}
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M3.06 13a9 9 0 1 0 .49 -4.087"></path>
                <path d="M3 4.001v5h5"></path>
                <path d="M12 12m-1 0a1 1 0 1 0 2 0a1 1 0 1 0 -2 0"></path>
              </svg>
            </button>
            <button
              type="button"
              className="users-action-icon users-action-neutral"
              onClick={() => setSortState({ key: null, direction: 'asc' })}
              title="Сбросить сортировку"
              aria-label="Сбросить сортировку"
              disabled={!sortState.key}
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M8 4h12v2.172a2 2 0 0 1 -.586 1.414l-3.914 3.914m-.5 3.5v4l-6 2v-8.5l-4.48 -4.928a2 2 0 0 1 -.52 -1.345v-2.227"></path>
                <path d="M3 3l18 18"></path>
              </svg>
            </button>
            <button
              type="button"
              className="users-action-icon users-action-danger"
              onClick={handleDeleteSelected}
              title={`Удалить выбранных: ${selectedUserUUIDs.size}`}
              aria-label={`Удалить выбранных: ${selectedUserUUIDs.size}`}
              disabled={selectedUserUUIDs.size === 0}
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"></path>
              </svg>
              {selectedUserUUIDs.size > 0 && <span className="users-action-count">{selectedUserUUIDs.size}</span>}
            </button>
            <button
              type="button"
              className="users-action-icon users-action-info"
              onClick={fetchUsers}
              title="Обновить список"
              aria-label="Обновить список"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4"></path>
                <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4"></path>
              </svg>
            </button>
            <button
              type="button"
              className="users-action-icon users-action-success"
              onClick={handleAdd}
              title="Добавить пользователя"
              aria-label="Добавить пользователя"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 5l0 14"></path>
                <path d="M5 12l14 0"></path>
              </svg>
            </button>
          </div>
        </div>

        <SharedRichTable
          columns={USER_COLUMNS}
          rows={sortedUsers}
          getRowId={(user) => user.uuid}
          renderCell={renderCell}
          onRowClick={openEditFromRow}
          getRowAriaLabel={(user) => `Редактировать пользователя ${user.username || ''}`}
          selectedRowIds={selectedUserUUIDs}
          setSelectedRowIds={setSelectedUserUUIDs}
          sortState={sortState}
          setSortState={setSortState}
          compactRows={compactRows}
          setCompactRows={setCompactRows}
          isFullscreen={isFullscreen}
          setIsFullscreen={setIsFullscreen}
          searchInputRef={searchInputRef}
          pageResetKey={searchQuery}
          pageSizeOptions={PAGE_SIZE_OPTIONS}
          pageSizeLabelId="users-rpp-label"
          tableAriaLabel="Таблица пользователей с горизонтальной прокруткой"
          emptyState={(
            <div className="empty-state users-empty-state">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2" />
                <circle cx="9" cy="7" r="4" />
              </svg>
              <h3>Пользователи не найдены</h3>
              <p>Проверьте фильтр или добавьте нового пользователя</p>
            </div>
          )}
        />
      </div>

      {!isFullscreen && (
        <div className="auto-refresh-row">
          <span className="auto-refresh-label">Auto-refresh:</span>
          <AppSelect
            className="auto-refresh-select"
            value={refreshInterval}
            onChange={(event) => setRefreshInterval(parseInt(event.target.value, 10))}
          >
            <option value="10">10 seconds</option>
            <option value="30">30 seconds</option>
            <option value="60">1 minute</option>
            <option value="300">5 minutes</option>
          </AppSelect>
        </div>
      )}

      {modalOpen && (
        <UserModal
          user={editingUser}
          onClose={() => setModalOpen(false)}
          onSave={handleSave}
        />
      )}
    </>
  );
}

export default Users;
