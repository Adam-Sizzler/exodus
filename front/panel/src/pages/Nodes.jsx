import { useState, useEffect, useMemo, useRef } from 'react';
import NodeModal from '../components/NodeModal';
import { nodesApi } from '../api';
import SharedRichTable from '../components/SharedRichTable';

const PAGE_SIZE_OPTIONS = [5, 10, 15, 20, 25, 30, 50, 100];

const NODE_COLUMNS = [
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
    key: 'name',
    label: 'Нода',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 230,
    minWidth: 180,
    defaultPin: 'left',
  },
  {
    key: 'address',
    label: 'Адрес',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 208,
    minWidth: 170,
    defaultPin: 'left',
  },
  {
    key: 'apiSchema',
    label: 'API Schema',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 130,
    minWidth: 110,
  },
  {
    key: 'apiPath',
    label: 'API Path',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 180,
    minWidth: 140,
  },
  {
    key: 'country',
    label: 'Страна',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 112,
    minWidth: 94,
  },
  {
    key: 'usersOnline',
    label: 'Пользователи онлайн',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 170,
    minWidth: 150,
  },
  {
    key: 'trafficUsage',
    label: 'Расход трафика',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 300,
    minWidth: 240,
  },
  {
    key: 'status',
    label: 'Статус',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 148,
    minWidth: 128,
  },
  {
    key: 'profile',
    label: 'Config Profile',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 210,
    minWidth: 170,
  },
  {
    key: 'xrayVersion',
    label: 'Xray Version',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 160,
    minWidth: 130,
  },
  {
    key: 'nodeVersion',
    label: 'Node Version',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 150,
    minWidth: 130,
  },
  {
    key: 'uptime',
    label: 'Uptime',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 150,
    minWidth: 120,
  },
  {
    key: 'trafficLimit',
    label: 'Лимит трафика',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 168,
    minWidth: 140,
  },
  {
    key: 'trafficUsed',
    label: 'Использовано',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 168,
    minWidth: 140,
  },
  {
    key: 'tags',
    label: 'Tags',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 220,
    minWidth: 170,
  },
  {
    key: 'createdAt',
    label: 'Дата создания',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 190,
    minWidth: 160,
  },
  {
    key: 'updatedAt',
    label: 'Дата обновления',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 190,
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

const mapApiNodeToUi = (node) => ({
  uuid: node.uuid,
  name: node.name ?? '',
  address: node.address ?? '',
  port: node.port ?? null,
  api_schema: node.apiSchema ?? 'grpc',
  api_path: node.apiPath ?? '',
  is_connected: !!node.isConnected,
  is_connecting: !!node.isConnecting,
  is_disabled: !!node.isDisabled,
  last_status_change: node.lastStatusChange ?? null,
  last_status_message: node.lastStatusMessage ?? null,
  xray_version: node.xrayVersion ?? null,
  node_version: node.nodeVersion ?? null,
  xray_uptime: node.xrayUptime ?? '',
  is_traffic_tracking_active: !!node.isTrafficTrackingActive,
  traffic_reset_day: node.trafficResetDay ?? null,
  traffic_limit_bytes: node.trafficLimitBytes ?? 0,
  traffic_used_bytes: node.trafficUsedBytes ?? 0,
  notify_percent: node.notifyPercent ?? 0,
  users_online: node.usersOnline ?? 0,
  view_position: node.viewPosition ?? 0,
  country_code: node.countryCode ?? 'XX',
  consumption_multiplier: node.consumptionMultiplier ?? 1,
  tags: Array.isArray(node.tags) ? node.tags : [],
  cpu_count: node.cpuCount ?? null,
  cpu_model: node.cpuModel ?? null,
  total_ram: node.totalRam ?? null,
  created_at: node.createdAt ?? null,
  updated_at: node.updatedAt ?? null,
  active_config_profile_uuid: node.configProfile?.activeConfigProfileUuid ?? null,
  active_inbounds: Array.isArray(node.configProfile?.activeInbounds) ? node.configProfile.activeInbounds : [],
  provider_uuid: node.providerUuid ?? null,
  provider: node.provider ?? null,
});

function Nodes() {
  const [nodes, setNodes] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalOpen, setModalOpen] = useState(false);
  const [editingNode, setEditingNode] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedNodeUUIDs, setSelectedNodeUUIDs] = useState(new Set());
  const [compactRows, setCompactRows] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [sortState, setSortState] = useState({ key: null, direction: 'asc' });

  const searchInputRef = useRef(null);

  useEffect(() => {
    fetchNodes();
  }, []);

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

  const fetchNodes = async () => {
    try {
      const data = await nodesApi.getAll();
      const list = Array.isArray(data?.response) ? data.response.map(mapApiNodeToUi) : [];

      setNodes(list);
      setSelectedNodeUUIDs((prev) => {
        const valid = new Set(list.map((node) => node.uuid));
        return new Set(Array.from(prev).filter((uuid) => valid.has(uuid)));
      });
      setLoading(false);
    } catch (err) {
      console.error('Error fetching nodes:', err);
      setLoading(false);
    }
  };

  const handleAdd = () => {
    setEditingNode(null);
    setModalOpen(true);
  };

  const handleEdit = (node) => {
    setEditingNode(node);
    setModalOpen(true);
  };

  const handleSave = async ({ payload, desiredDisabled }) => {
    try {
      if (editingNode) {
        await nodesApi.update({ uuid: editingNode.uuid, ...payload });
        if (typeof desiredDisabled === 'boolean' && desiredDisabled !== Boolean(editingNode.is_disabled)) {
          if (desiredDisabled) {
            await nodesApi.disable(editingNode.uuid);
          } else {
            await nodesApi.enable(editingNode.uuid);
          }
        }
      } else {
        const created = await nodesApi.create(payload);
        const createdUUID = created?.response?.uuid;
        if (createdUUID && desiredDisabled) {
          await nodesApi.disable(createdUUID);
        }
      }
      setModalOpen(false);
      fetchNodes();
    } catch (err) {
      console.error('Error saving node:', err);
      alert(`Failed to save node: ${err.message}`);
    }
  };

  const handleDeleteSelected = async () => {
    if (selectedNodeUUIDs.size === 0) return;
    if (!confirm(`Delete ${selectedNodeUUIDs.size} selected node(s)?`)) return;

    try {
      await Promise.all(Array.from(selectedNodeUUIDs).map((uuid) => nodesApi.delete(uuid)));
      await fetchNodes();
    } catch (err) {
      console.error('Error deleting selected nodes:', err);
      alert(`Failed to delete selected nodes: ${err.message}`);
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

  const formatSpeed = (bytesPerSec) => `${formatBytes(bytesPerSec)}/s`;

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

  const toCountryFlag = (countryCode) => {
    const code = String(countryCode || '').trim().toUpperCase();
    if (!/^[A-Z]{2}$/.test(code)) {
      return '';
    }
    return Array.from(code).map((char) => String.fromCodePoint(127397 + char.charCodeAt(0))).join('');
  };

  const isNodeOnline = (node) => {
    return (node.is_connected || node.is_active) && !node.is_disabled;
  };

  const getNodeStatusLabel = (node) => {
    if (node.is_disabled) return 'DISABLED';
    return isNodeOnline(node) ? 'ONLINE' : 'OFFLINE';
  };

  const getNodeStatusClass = (node) => {
    if (node.is_disabled) return 'users-status-disabled';
    return isNodeOnline(node) ? 'users-status-active' : 'users-status-expired';
  };

  const getNodeTags = (node) => {
    if (Array.isArray(node.tags)) {
      return node.tags.filter(Boolean).join(', ');
    }
    if (typeof node.tags === 'string') {
      try {
        const parsed = JSON.parse(node.tags);
        if (Array.isArray(parsed)) {
          return parsed.filter(Boolean).join(', ');
        }
      } catch (_) {
        return node.tags;
      }
      return node.tags;
    }
    return '';
  };

  const getNodeApiSchema = (node) => {
    const directValue = (
      node.api_schema
      ?? node.apiSchema
      ?? node.api_scheme
      ?? node.apiScheme
      ?? node.api_protocol
      ?? node.apiProtocol
      ?? node.api_type
      ?? node.apiType
      ?? node.schema
    );

    if (directValue !== null && directValue !== undefined && String(directValue).trim() !== '') {
      return String(directValue).trim().toLowerCase();
    }

    return '';
  };

  const getNodeUploadBps = (node) => {
    return Number(
      node.upload_bytes_per_second
        ?? node.upload_bps
        ?? node.uplink_bps
        ?? node.up_speed_bps
        ?? 0,
    );
  };

  const getNodeDownloadBps = (node) => {
    return Number(
      node.download_bytes_per_second
        ?? node.download_bps
        ?? node.downlink_bps
        ?? node.down_speed_bps
        ?? 0,
    );
  };

  const getNodeTrafficUsage = (node) => {
    const used = Math.max(0, Number(node.traffic_used_bytes ?? 0));
    const limit = Math.max(0, Number(node.traffic_limit_bytes ?? 0));
    const hasLimit = limit > 0;
    const percent = hasLimit ? Math.min(100, (used / limit) * 100) : 0;
    return {
      used,
      limit,
      hasLimit,
      percent,
      freePercent: hasLimit ? Math.max(0, 100 - percent) : 100,
    };
  };

  const getSortValue = (node, key) => {
    switch (key) {
      case 'name':
        return String(node.name || '').toLowerCase();
      case 'address':
        return String(`${node.address || ''}:${node.port || ''}`).toLowerCase();
      case 'apiSchema':
        return getNodeApiSchema(node);
      case 'apiPath':
        return String(node.api_path || '').toLowerCase();
      case 'country':
        return String(node.country_code || '').toLowerCase();
      case 'usersOnline':
        return Number(node.users_online ?? 0);
      case 'trafficUsage':
        return Number(node.traffic_used_bytes ?? 0);
      case 'status':
        return isNodeOnline(node) ? 0 : (node.is_disabled ? 2 : 1);
      case 'profile':
        return String(node.config_profile_name || node.active_config_profile_uuid || '').toLowerCase();
      case 'xrayVersion':
        return String(node.xray_version || '').toLowerCase();
      case 'nodeVersion':
        return String(node.node_version || '').toLowerCase();
      case 'uptime':
        return String(node.xray_uptime || '').toLowerCase();
      case 'trafficLimit':
        return Number(node.traffic_limit_bytes ?? 0);
      case 'trafficUsed':
        return Number(node.traffic_used_bytes ?? 0);
      case 'tags':
        return getNodeTags(node).toLowerCase();
      case 'createdAt':
        return node.created_at ? new Date(node.created_at).getTime() : 0;
      case 'updatedAt':
        return node.updated_at ? new Date(node.updated_at).getTime() : 0;
      case 'uuid':
        return String(node.uuid || '').toLowerCase();
      default:
        return '';
    }
  };

  const query = searchQuery.trim().toLowerCase();

  const filteredNodes = useMemo(() => {
    if (!query) return nodes;

    return nodes.filter((node) => {
      const pool = [
        node.name,
        node.address,
        getNodeApiSchema(node),
        node.api_path,
        node.country_code,
        node.config_profile_name,
        node.active_config_profile_uuid,
        node.xray_version,
        node.node_version,
        node.uuid,
        getNodeTags(node),
      ]
        .filter(Boolean)
        .map((value) => String(value).toLowerCase());

      return pool.some((value) => value.includes(query));
    });
  }, [nodes, query]);

  const sortedNodes = useMemo(() => {
    if (!sortState.key) {
      return filteredNodes;
    }

    const directionFactor = sortState.direction === 'asc' ? 1 : -1;

    return filteredNodes
      .map((node, index) => ({ node, index }))
      .sort((a, b) => {
        const aValue = getSortValue(a.node, sortState.key);
        const bValue = getSortValue(b.node, sortState.key);

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
      .map((item) => item.node);
  }, [filteredNodes, sortState]);

  const nodesStats = useMemo(() => {
    const total = nodes.length;
    const online = nodes.filter((node) => isNodeOnline(node)).length;
    const offline = total - online;
    const usersOnline = nodes.reduce((sum, node) => sum + Number(node.users_online ?? 0), 0);
    const totalTraffic = nodes.reduce((sum, node) => sum + Number(node.traffic_used_bytes ?? 0), 0);
    const totalUploadBps = nodes.reduce((sum, node) => sum + getNodeUploadBps(node), 0);
    const totalDownloadBps = nodes.reduce((sum, node) => sum + getNodeDownloadBps(node), 0);
    const averageSpeedBps = total > 0 ? Math.round((totalUploadBps + totalDownloadBps) / total) : 0;
    const activeNodes = nodes.filter((node) => !node.is_disabled && (node.is_connected || node.is_active)).length;

    return {
      total,
      online,
      offline,
      usersOnline,
      totalTraffic,
      totalUploadBps,
      totalDownloadBps,
      averageSpeedBps,
      activeNodes,
    };
  }, [nodes]);

  const renderCell = (column, node) => {
    const countryFlag = toCountryFlag(node.country_code);
    const trafficUsage = getNodeTrafficUsage(node);

    switch (column.key) {
      case 'name':
        return (
          <td key={column.key}>
            <div className="users-name-cell">
              <span className={`users-presence-dot ${isNodeOnline(node) ? 'is-online' : 'is-idle'}`}></span>
              <div className="users-name-text">
                <strong>{node.name || 'Unnamed'}</strong>
              </div>
            </div>
          </td>
        );
      case 'address':
        return <td key={column.key} className="users-cell-mono">{node.address || '—'}:{node.port || '80'}</td>;
      case 'apiSchema':
        return <td key={column.key} className="users-cell-center users-cell-mono">{getNodeApiSchema(node) || '—'}</td>;
      case 'apiPath':
        return <td key={column.key} className="users-cell-mono">{node.api_path || '—'}</td>;
      case 'country':
        return <td key={column.key} className="users-cell-center">{countryFlag ? `${countryFlag} ${node.country_code}` : (node.country_code || '—')}</td>;
      case 'usersOnline':
        return <td key={column.key} className="users-cell-center users-cell-strong">{node.users_online ?? 0}</td>;
      case 'trafficUsage':
        return (
          <td key={column.key}>
            <div className="users-traffic-cell">
              <div className="users-traffic-meta">
                <span className="users-traffic-used-label">{trafficUsage.percent.toFixed(2)}% {trafficUsage.hasLimit ? '' : '∞'}</span>
                <span className="users-traffic-free-label">Σ {formatBytes(trafficUsage.used)} {trafficUsage.freePercent.toFixed(2)}%</span>
              </div>
              <div className="users-traffic-bar" aria-label="Использование трафика">
                <span className="users-traffic-used" style={{ width: `${trafficUsage.percent}%` }}></span>
                <span className="users-traffic-free" style={{ width: `${trafficUsage.freePercent}%` }}></span>
              </div>
              <div className="users-traffic-foot">
                <span>{formatBytes(trafficUsage.used)}</span>
                <span>{trafficUsage.hasLimit ? formatBytes(trafficUsage.limit) : '∞'}</span>
              </div>
            </div>
          </td>
        );
      case 'status':
        return (
          <td key={column.key} className="users-cell-center">
            <span className={`users-status-badge ${getNodeStatusClass(node)}`}>
              <span className="users-status-dot"></span>
              {getNodeStatusLabel(node)}
            </span>
          </td>
        );
      case 'profile':
        return <td key={column.key}>{node.config_profile_name || node.active_config_profile_uuid || '—'}</td>;
      case 'xrayVersion':
        return <td key={column.key} className="users-cell-center users-cell-mono">{node.xray_version || '—'}</td>;
      case 'nodeVersion':
        return <td key={column.key} className="users-cell-center users-cell-mono">{node.node_version || '—'}</td>;
      case 'uptime':
        return <td key={column.key} className="users-cell-center">{node.xray_uptime || '—'}</td>;
      case 'trafficLimit':
        return <td key={column.key} className="users-cell-center users-cell-mono">{node.traffic_limit_bytes ? formatBytes(node.traffic_limit_bytes) : '∞'}</td>;
      case 'trafficUsed':
        return <td key={column.key} className="users-cell-center users-cell-mono">{formatBytes(node.traffic_used_bytes ?? 0)}</td>;
      case 'tags':
        return <td key={column.key}>{getNodeTags(node) || '—'}</td>;
      case 'createdAt':
        return <td key={column.key} className="users-cell-center">{formatDateTime(node.created_at)}</td>;
      case 'updatedAt':
        return <td key={column.key} className="users-cell-center">{formatDateTime(node.updated_at)}</td>;
      case 'uuid':
        return <td key={column.key} className="users-cell-mono users-cell-dim">{node.uuid || '—'}</td>;
      default:
        return <td key={column.key}>—</td>;
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
      <div className="nodes-overview-grid">
        <div className="nodes-overview-card theme-cyan">
          <div className="nodes-overview-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M9 7m-4 0a4 4 0 1 0 8 0a4 4 0 1 0 -8 0"></path>
              <path d="M3 21v-2a4 4 0 0 1 4 -4h4a4 4 0 0 1 4 4v2"></path>
              <path d="M16 3.13a4 4 0 0 1 0 7.75"></path>
              <path d="M21 21v-2a4 4 0 0 0 -3 -3.85"></path>
            </svg>
          </div>
          <div className="nodes-overview-content">
            <p className="nodes-overview-label">Пользователи онлайн</p>
            <p className="nodes-overview-value">{nodesStats.usersOnline}</p>
          </div>
        </div>

        <div className="nodes-overview-card theme-teal">
          <div className="nodes-overview-icon">
            <svg viewBox="0 0 256 256" fill="currentColor">
              <path d="M240,128a8,8,0,0,1-8,8H204.94l-37.78,75.58A8,8,0,0,1,160,216h-.4a8,8,0,0,1-7.08-5.14L95.35,60.76,63.28,131.31A8,8,0,0,1,56,136H24a8,8,0,0,1,0-16H50.85L88.72,36.69a8,8,0,0,1,14.76.46l57.51,151,31.85-63.71A8,8,0,0,1,200,120h32A8,8,0,0,1,240,128Z"></path>
            </svg>
          </div>
          <div className="nodes-overview-content">
            <p className="nodes-overview-label">Ноды онлайн</p>
            <p className="nodes-overview-value">{nodesStats.online}</p>
          </div>
        </div>

        <div className="nodes-overview-card theme-red">
          <div className="nodes-overview-icon">
            <svg viewBox="0 0 256 256" fill="currentColor">
              <path d="M128,24A104,104,0,1,0,232,128,104.11,104.11,0,0,0,128,24Zm0,192a88,88,0,1,1,88-88A88.1,88.1,0,0,1,128,216Zm-8-80V80a8,8,0,0,1,16,0v56a8,8,0,0,1-16,0Zm20,36a12,12,0,1,1-12-12A12,12,0,0,1,140,172Z"></path>
            </svg>
          </div>
          <div className="nodes-overview-content">
            <p className="nodes-overview-label">Ноды оффлайн</p>
            <p className="nodes-overview-value">{nodesStats.offline}</p>
          </div>
        </div>

        <div className="nodes-overview-card theme-cyan">
          <div className="nodes-overview-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M18 16v2a1 1 0 0 1 -1 1h-11l6 -7l-6 -7h11a1 1 0 0 1 1 1v2"></path>
            </svg>
          </div>
          <div className="nodes-overview-content">
            <p className="nodes-overview-label">Общее количество трафика</p>
            <p className="nodes-overview-value">{formatBytes(nodesStats.totalTraffic)}</p>
          </div>
        </div>

        <div className="nodes-overview-card theme-blue">
          <div className="nodes-overview-icon">
            <svg viewBox="0 0 256 256" fill="currentColor">
              <path d="M200,112H56l72-72Z" opacity="0.2"></path>
              <path d="M205.66,106.34l-72-72a8,8,0,0,0-11.32,0l-72,72A8,8,0,0,0,56,120h64v96a8,8,0,0,0,16,0V120h64a8,8,0,0,0,5.66-13.66ZM75.31,104,128,51.31,180.69,104Z"></path>
            </svg>
          </div>
          <div className="nodes-overview-content">
            <p className="nodes-overview-label">Всего загружено</p>
            <p className="nodes-overview-value">{formatSpeed(nodesStats.totalUploadBps)}</p>
            <p className="nodes-overview-subtitle">Текущий час</p>
          </div>
        </div>

        <div className="nodes-overview-card theme-teal">
          <div className="nodes-overview-icon">
            <svg viewBox="0 0 256 256" fill="currentColor">
              <path d="M200,144l-72,72L56,144Z" opacity="0.2"></path>
              <path d="M207.39,140.94A8,8,0,0,0,200,136H136V40a8,8,0,0,0-16,0v96H56a8,8,0,0,0-5.66,13.66l72,72a8,8,0,0,0,11.32,0l72-72A8,8,0,0,0,207.39,140.94ZM128,204.69,75.31,152H180.69Z"></path>
            </svg>
          </div>
          <div className="nodes-overview-content">
            <p className="nodes-overview-label">Всего скачано</p>
            <p className="nodes-overview-value">{formatSpeed(nodesStats.totalDownloadBps)}</p>
            <p className="nodes-overview-subtitle">Текущий час</p>
          </div>
        </div>

        <div className="nodes-overview-card theme-indigo">
          <div className="nodes-overview-icon">
            <svg viewBox="0 0 256 256" fill="currentColor">
              <path d="M114.34,154.34l96-96a8,8,0,0,1,11.32,11.32l-96,96a8,8,0,0,1-11.32-11.32ZM128,88a63.9,63.9,0,0,1,20.44,3.33,8,8,0,1,0,5.11-15.16A80,80,0,0,0,48.49,160.88,8,8,0,0,0,56.43,168c.29,0,.59,0,.89-.05a8,8,0,0,0,7.07-8.83A64.92,64.92,0,0,1,64,152,64.07,64.07,0,0,1,128,88Zm99.74,13a8,8,0,0,0-14.24,7.3,96.27,96.27,0,0,1,5,75.71l-181.1-.07A96.24,96.24,0,0,1,128,56h.88a95,95,0,0,1,42.82,10.5A8,8,0,1,0,179,52.27a112,112,0,0,0-156.66,137A16.07,16.07,0,0,0,37.46,200H218.53a16,16,0,0,0,15.11-10.71,112.35,112.35,0,0,0-5.9-88.3Z"></path>
            </svg>
          </div>
          <div className="nodes-overview-content">
            <p className="nodes-overview-label">Средняя скорость</p>
            <p className="nodes-overview-value">{formatSpeed(nodesStats.averageSpeedBps)}</p>
            <p className="nodes-overview-subtitle">Текущий час</p>
          </div>
        </div>

        <div className="nodes-overview-card theme-indigo">
          <div className="nodes-overview-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M3 4m0 3a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v2a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3z"></path>
              <path d="M3 12m0 3a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v2a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3z"></path>
              <path d="M7 8l0 .01"></path>
              <path d="M7 16l0 .01"></path>
              <path d="M11 8h6"></path>
              <path d="M11 16h6"></path>
            </svg>
          </div>
          <div className="nodes-overview-content">
            <p className="nodes-overview-label">Активные ноды</p>
            <p className="nodes-overview-value">{nodesStats.activeNodes}</p>
            <p className="nodes-overview-subtitle">Текущий час</p>
          </div>
        </div>
      </div>

      <div className={`card users-list-card ${isFullscreen ? 'users-list-card-fullscreen' : ''}`}>
        <div className="users-list-header">
          <div className="users-list-header-main">
            <div className="users-action-icon users-action-cyan users-section-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10" />
                <line x1="2" y1="12" x2="22" y2="12" />
                <path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z" />
              </svg>
            </div>
            <div className="users-list-heading-stack">
              <h2 className="card-title">Ноды</h2>
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
                placeholder="Поиск нод..."
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
              title={`Удалить выбранных: ${selectedNodeUUIDs.size}`}
              aria-label={`Удалить выбранных: ${selectedNodeUUIDs.size}`}
              disabled={selectedNodeUUIDs.size === 0}
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"></path>
              </svg>
              {selectedNodeUUIDs.size > 0 && <span className="users-action-count">{selectedNodeUUIDs.size}</span>}
            </button>
            <button
              type="button"
              className="users-action-icon users-action-info"
              onClick={fetchNodes}
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
              title="Добавить ноду"
              aria-label="Добавить ноду"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 5l0 14"></path>
                <path d="M5 12l14 0"></path>
              </svg>
            </button>
          </div>
        </div>

        <SharedRichTable
          columns={NODE_COLUMNS}
          rows={sortedNodes}
          getRowId={(node) => node.uuid}
          renderCell={renderCell}
          onRowClick={handleEdit}
          getRowAriaLabel={(node) => `Редактировать ноду ${node.name || ''}`}
          selectedRowIds={selectedNodeUUIDs}
          setSelectedRowIds={setSelectedNodeUUIDs}
          sortState={sortState}
          setSortState={setSortState}
          compactRows={compactRows}
          setCompactRows={setCompactRows}
          isFullscreen={isFullscreen}
          setIsFullscreen={setIsFullscreen}
          searchInputRef={searchInputRef}
          pageResetKey={searchQuery}
          pageSizeOptions={PAGE_SIZE_OPTIONS}
          pageSizeLabelId="nodes-rpp-label"
          tableAriaLabel="Таблица нод с горизонтальной прокруткой"
          emptyState={(
            <div className="empty-state users-empty-state">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="2" y1="12" x2="22" y2="12"></line>
                <path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"></path>
              </svg>
              <h3>Ноды не найдены</h3>
              <p>Проверьте фильтр или добавьте новую ноду</p>
            </div>
          )}
        />
      </div>

      {modalOpen && (
        <NodeModal
          node={editingNode}
          onClose={() => setModalOpen(false)}
          onSave={handleSave}
        />
      )}
    </>
  );
}

export default Nodes;
