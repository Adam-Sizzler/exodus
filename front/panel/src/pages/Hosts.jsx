import React, { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { configProfilesApi, hostsApi, nodesApi } from '../api';
import AppMultiSelect from '../components/AppMultiSelect';
import AppSelect from '../components/AppSelect';
import SharedRichTable from '../components/SharedRichTable';
import SideDrawerPanel from '../components/SideDrawerPanel';

const SECURITY_LAYER_OPTIONS = [
  { value: 'DEFAULT', label: 'По умолчанию' },
  { value: 'TLS', label: 'TLS' },
  { value: 'NONE', label: 'Без TLS' },
];
const FINGERPRINT_OPTIONS = ['chrome', 'firefox', 'safari', 'ios', 'android', 'edge', 'qq', 'random', 'randomized'];
const ALPN_OPTIONS = ['h3', 'h2', 'http/1.1', 'h2,http/1.1', 'h3,h2,http/1.1', 'h3,h2'];
const PAGE_SIZE_OPTIONS = [5, 10, 15, 20, 25, 30, 50, 100];

const HOST_COLUMNS = [
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
    key: 'remark',
    label: 'Remark',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 250,
    minWidth: 190,
    defaultPin: 'left',
  },
  {
    key: 'address',
    label: 'Address',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 224,
    minWidth: 180,
    defaultPin: 'left',
  },
  {
    key: 'port',
    label: 'Port',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 90,
    minWidth: 74,
  },
  {
    key: 'status',
    label: 'Статус',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 148,
    minWidth: 126,
  },
  {
    key: 'tag',
    label: 'Tag',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 146,
    minWidth: 110,
  },
  {
    key: 'inbound',
    label: 'Inbound',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 190,
    minWidth: 150,
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
    key: 'nodes',
    label: 'Nodes',
    sortable: true,
    defaultVisible: true,
    defaultWidth: 320,
    minWidth: 240,
  },
  {
    key: 'securityLayer',
    label: 'Security',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 132,
    minWidth: 110,
  },
  {
    key: 'allowInsecure',
    label: 'Allow Insecure',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 176,
    minWidth: 145,
  },
  {
    key: 'isHidden',
    label: 'Hidden',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 118,
    minWidth: 100,
  },
  {
    key: 'viewPosition',
    label: 'Позиция',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 116,
    minWidth: 96,
  },
  {
    key: 'sni',
    label: 'SNI',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 220,
    minWidth: 160,
  },
  {
    key: 'host',
    label: 'Host Header',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 220,
    minWidth: 160,
  },
  {
    key: 'path',
    label: 'Path',
    sortable: true,
    defaultVisible: false,
    defaultWidth: 220,
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

const defaultHostForm = () => ({
  view_position: 0,
  remark: '',
  tag: '',
  address: '',
  port: 443,
  config_profile_uuid: '',
  config_profile_inbound_uuid: '',
  path: '',
  sni: '',
  host: '',
  alpn: '',
  fingerprint: '',
  security_layer: 'DEFAULT',
  xhttp_extra_params: '',
  mux_params: '',
  sockopt_params: '',
  is_disabled: false,
  server_description: '',
  vless_route_id: '',
  allow_insecure: false,
  shuffle_host: false,
  mihomo_x25519: false,
  xray_json_template_uuid: '',
  keep_sni_blank: false,
  is_hidden: false,
  override_sni_from_address: false,
  nodes: [],
  excluded_internal_squads: [],
  exclude_from_subscription_types: [],
});

const toNullable = (value) => {
  if (value === null || value === undefined) return '';
  return String(value);
};

const normalizeSecurityLayer = (value) => {
  const normalized = String(value || '').trim().toUpperCase();
  if (SECURITY_LAYER_OPTIONS.some((item) => item.value === normalized)) {
    return normalized;
  }
  return 'DEFAULT';
};

const toCountryFlag = (countryCode) => {
  const code = String(countryCode || '').trim().toUpperCase();
  if (!/^[A-Z]{2}$/.test(code)) {
    return '';
  }
  return Array.from(code).map((char) => String.fromCodePoint(127397 + char.charCodeAt(0))).join('');
};

const buildHostNodeMap = (hostsList = []) => {
  const map = new Map();
  hostsList.forEach((host) => {
    if (!host?.uuid) {
      return;
    }
    map.set(host.uuid, new Set(host.nodes || []));
  });
  return map;
};

const toJSONString = (value) => {
  if (value === null || value === undefined) return '';
  if (typeof value === 'string') return value;
  try {
    return JSON.stringify(value, null, 2);
  } catch (err) {
    return String(value);
  }
};

const parseJSONInput = (value, label) => {
  const trimmed = String(value || '').trim();
  if (!trimmed) return null;
  try {
    return JSON.parse(trimmed);
  } catch (err) {
    throw new Error(`Некорректный JSON в поле ${label}`);
  }
};

const mapApiHostToUi = (host) => ({
  uuid: host.uuid,
  view_position: host.viewPosition ?? 0,
  remark: host.remark ?? '',
  tag: host.tag ?? '',
  address: host.address ?? '',
  port: host.port ?? 443,
  config_profile_uuid: host.inbound?.configProfileUuid ?? '',
  config_profile_inbound_uuid: host.inbound?.configProfileInboundUuid ?? '',
  path: host.path ?? '',
  sni: host.sni ?? '',
  host: host.host ?? '',
  alpn: host.alpn ?? '',
  fingerprint: host.fingerprint ?? '',
  security_layer: normalizeSecurityLayer(host.securityLayer),
  xhttp_extra_params: toJSONString(host.xHttpExtraParams),
  mux_params: toJSONString(host.muxParams),
  sockopt_params: toJSONString(host.sockoptParams),
  is_disabled: !!host.isDisabled,
  server_description: host.serverDescription ?? '',
  vless_route_id: host.vlessRouteId === null || host.vlessRouteId === undefined ? '' : String(host.vlessRouteId),
  allow_insecure: !!host.allowInsecure,
  shuffle_host: !!host.shuffleHost,
  mihomo_x25519: !!host.mihomoX25519,
  xray_json_template_uuid: host.xrayJsonTemplateUuid ?? '',
  keep_sni_blank: !!host.keepSniBlank,
  is_hidden: !!host.isHidden,
  override_sni_from_address: !!host.overrideSniFromAddress,
  nodes: Array.isArray(host.nodes) ? host.nodes : [],
  excluded_internal_squads: Array.isArray(host.excludedInternalSquads) ? host.excludedInternalSquads : [],
  exclude_from_subscription_types: Array.isArray(host.excludeFromSubscriptionTypes) ? host.excludeFromSubscriptionTypes : [],
});

function Hosts() {
  const [hosts, setHosts] = useState([]);
  const [profiles, setProfiles] = useState([]);
  const [nodes, setNodes] = useState([]);
  const [hostNodesMap, setHostNodesMap] = useState(new Map());
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [editingHost, setEditingHost] = useState(null);
  const [mode, setMode] = useState('basic');
  const [formData, setFormData] = useState(defaultHostForm());
  const [selectedHostUUIDs, setSelectedHostUUIDs] = useState(new Set());
  const [savingHost, setSavingHost] = useState(false);
  const [searchQuery, setSearchQuery] = useState('');
  const [compactRows, setCompactRows] = useState(false);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [sortState, setSortState] = useState({ key: null, direction: 'asc' });

  const searchInputRef = useRef(null);

  const profileNameByUuid = useMemo(() => {
    const map = new Map();
    profiles.forEach((profile) => map.set(profile.uuid, profile.name));
    return map;
  }, [profiles]);

  const inboundList = useMemo(() => {
    const list = [];
    profiles.forEach((profile) => {
      (profile.inbounds || []).forEach((inbound) => {
        list.push({
          uuid: inbound.uuid,
          tag: inbound.tag,
          profile_uuid: profile.uuid,
          profile_name: profile.name,
        });
      });
    });
    return list;
  }, [profiles]);

  const inboundByUuid = useMemo(() => {
    const map = new Map();
    inboundList.forEach((item) => map.set(item.uuid, item));
    return map;
  }, [inboundList]);

  const visibleInboundOptions = useMemo(() => {
    if (!formData.config_profile_uuid) return inboundList;
    return inboundList.filter((inbound) => inbound.profile_uuid === formData.config_profile_uuid);
  }, [formData.config_profile_uuid, inboundList]);

  const selectedInbound = useMemo(() => {
    if (!formData.config_profile_inbound_uuid) return null;
    return inboundByUuid.get(formData.config_profile_inbound_uuid) || null;
  }, [formData.config_profile_inbound_uuid, inboundByUuid]);

  const nodeByUuid = useMemo(() => {
    const map = new Map();
    nodes.forEach((node) => map.set(node.uuid, node));
    return map;
  }, [nodes]);

  const nodeOptions = useMemo(
    () =>
      nodes.map((node) => {
        const flag = toCountryFlag(node.country_code);
        const nodeName = node.name || node.uuid;
        return {
          value: node.uuid,
          label: flag ? `${flag} ${nodeName}` : nodeName,
          searchText: `${nodeName} ${node.country_code || ''} ${node.address || ''}`,
        };
      }),
    [nodes]
  );

  const loadHosts = async () => {
    try {
      setLoading(true);
      const hostsData = await hostsApi.getAll();
      const mappedHosts = (hostsData.response || []).map(mapApiHostToUi);
      setHosts(mappedHosts);
      setHostNodesMap(buildHostNodeMap(mappedHosts));
      setSelectedHostUUIDs(new Set());
      setError(null);
    } catch (err) {
      console.error('Failed to load hosts:', err);
      setError(err.message || 'Failed to load hosts');
    } finally {
      setLoading(false);
    }
  };

  const loadNodes = async () => {
    try {
      const data = await nodesApi.getAll();
      setNodes(data.nodes || []);
    } catch (err) {
      console.error('Failed to load nodes:', err);
      setNodes([]);
    }
  };

  const loadProfiles = async () => {
    try {
      const data = await configProfilesApi.getAllWithInbounds();
      setProfiles(data.profiles || []);
    } catch (err) {
      console.error('Failed to load config profiles:', err);
      setProfiles([]);
    }
  };

  useEffect(() => {
    loadHosts();
    loadProfiles();
    loadNodes();
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

  const fillFormFromHost = (host, nodeUUIDs = []) => ({
    view_position: host.view_position ?? 0,
    remark: toNullable(host.remark),
    tag: toNullable(host.tag),
    address: toNullable(host.address),
    port: host.port ?? 443,
    config_profile_uuid: toNullable(host.config_profile_uuid),
    config_profile_inbound_uuid: toNullable(host.config_profile_inbound_uuid),
    path: toNullable(host.path),
    sni: toNullable(host.sni),
    host: toNullable(host.host),
    alpn: toNullable(host.alpn),
    fingerprint: toNullable(host.fingerprint),
    security_layer: normalizeSecurityLayer(host.security_layer),
    xhttp_extra_params: toNullable(host.xhttp_extra_params),
    mux_params: toNullable(host.mux_params),
    sockopt_params: toNullable(host.sockopt_params),
    is_disabled: !!host.is_disabled,
    server_description: toNullable(host.server_description),
    vless_route_id: host.vless_route_id === null || host.vless_route_id === undefined ? '' : String(host.vless_route_id),
    allow_insecure: !!host.allow_insecure,
    shuffle_host: !!host.shuffle_host,
    mihomo_x25519: !!host.mihomo_x25519,
    xray_json_template_uuid: toNullable(host.xray_json_template_uuid),
    keep_sni_blank: !!host.keep_sni_blank,
    is_hidden: !!host.is_hidden,
    override_sni_from_address: !!host.override_sni_from_address,
    nodes: [...nodeUUIDs],
    excluded_internal_squads: Array.isArray(host.excluded_internal_squads) ? host.excluded_internal_squads : [],
    exclude_from_subscription_types: Array.isArray(host.exclude_from_subscription_types) ? host.exclude_from_subscription_types : [],
  });

  const openCreate = () => {
    setEditingHost(null);
    setFormData({ ...defaultHostForm(), view_position: hosts.length });
    setMode('basic');
    setShowModal(true);
  };

  const openEdit = (host) => {
    setEditingHost(host);
    const nodeUUIDs = Array.from(hostNodesMap.get(host.uuid) || []);
    setFormData(fillFormFromHost(host, nodeUUIDs));
    setMode('basic');
    setShowModal(true);
  };

  const closeEditor = () => {
    if (savingHost) return;
    setShowModal(false);
  };

  const setField = (key, value) => {
    setFormData((previous) => ({ ...previous, [key]: value }));
  };

  const setALPNValues = (options) => {
    setFormData((previous) => ({ ...previous, alpn: options.join(',') }));
  };

  const normalizePayload = () => {
    const toNull = (value) => {
      if (value === undefined || value === null) return null;
      const normalized = String(value).trim();
      return normalized === '' ? null : normalized;
    };
    const toStringOrEmpty = (value) => String(value ?? '');

    const inbound = (() => {
      const configProfileUuid = toNull(formData.config_profile_uuid);
      const configProfileInboundUuid = toNull(formData.config_profile_inbound_uuid);
      if (!configProfileUuid || !configProfileInboundUuid) return null;
      return {
        configProfileUuid,
        configProfileInboundUuid,
      };
    })();

    const payload = {
      remark: String(formData.remark || ''),
      tag: toNull(formData.tag),
      address: String(formData.address || '').trim(),
      port: Number(formData.port),
      path: toStringOrEmpty(formData.path),
      sni: toStringOrEmpty(formData.sni),
      host: toStringOrEmpty(formData.host),
      alpn: toNull(formData.alpn),
      fingerprint: toNull(formData.fingerprint),
      securityLayer: normalizeSecurityLayer(formData.security_layer),
      xHttpExtraParams: parseJSONInput(formData.xhttp_extra_params, 'xhttp_extra_params'),
      muxParams: parseJSONInput(formData.mux_params, 'mux_params'),
      sockoptParams: parseJSONInput(formData.sockopt_params, 'sockopt_params'),
      isDisabled: !!formData.is_disabled,
      serverDescription: toNull(formData.server_description),
      vlessRouteId: formData.vless_route_id === '' ? null : Number(formData.vless_route_id),
      allowInsecure: !!formData.allow_insecure,
      shuffleHost: !!formData.shuffle_host,
      mihomoX25519: !!formData.mihomo_x25519,
      xrayJsonTemplateUuid: toNull(formData.xray_json_template_uuid),
      keepSniBlank: !!formData.keep_sni_blank,
      isHidden: !!formData.is_hidden,
      overrideSniFromAddress: !!formData.override_sni_from_address,
      nodes: Array.from(new Set((formData.nodes || []).filter(Boolean))),
      excludedInternalSquads: Array.from(new Set((formData.excluded_internal_squads || []).filter(Boolean))),
      excludeFromSubscriptionTypes: Array.from(new Set((formData.exclude_from_subscription_types || []).filter(Boolean))),
    };

    if (inbound) {
      payload.inbound = inbound;
    }

    return payload;
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    let payload;
    try {
      payload = normalizePayload();
    } catch (err) {
      alert(err.message || 'Некорректные данные');
      return;
    }

    if (!payload.remark) {
      alert('Поле Remark обязательно');
      return;
    }
    if (!payload.address) {
      alert('Поле Address обязательно');
      return;
    }
    if (!payload.port || payload.port < 1 || payload.port > 65535) {
      alert('Порт должен быть в диапазоне 1..65535');
      return;
    }
    const profileSelected = String(formData.config_profile_uuid || '').trim();
    const inboundSelected = String(formData.config_profile_inbound_uuid || '').trim();
    if (!profileSelected || !inboundSelected) {
      alert('Выберите Config Profile и Inbound');
      return;
    }

    try {
      setSavingHost(true);
      if (editingHost) {
        await hostsApi.update({ uuid: editingHost.uuid, ...payload });
      } else {
        await hostsApi.create(payload);
      }

      setShowModal(false);
      await loadHosts();
    } catch (err) {
      alert(`Failed to save host: ${err.message}`);
    } finally {
      setSavingHost(false);
    }
  };

  const handleDeleteSelected = async () => {
    if (selectedHostUUIDs.size === 0) return;
    if (!confirm(`Delete ${selectedHostUUIDs.size} selected host(s)?`)) return;

    try {
      await hostsApi.deleteMany(Array.from(selectedHostUUIDs));
      await loadHosts();
    } catch (err) {
      alert(`Failed to delete selected hosts: ${err.message}`);
    }
  };

  const resolveProfileName = (host) => {
    if (host.config_profile_name) return host.config_profile_name;
    if (host.config_profile_uuid) return profileNameByUuid.get(host.config_profile_uuid) || '-';
    return '-';
  };

  const resolveInboundTag = (host) => {
    if (host.config_profile_inbound_tag) return host.config_profile_inbound_tag;
    if (host.config_profile_inbound_uuid) return inboundByUuid.get(host.config_profile_inbound_uuid)?.tag || '-';
    return '-';
  };

  const handleProfileChange = (profileUUID) => {
    setFormData((previous) => {
      const next = { ...previous, config_profile_uuid: profileUUID };
      if (
        next.config_profile_inbound_uuid &&
        profileUUID &&
        inboundByUuid.get(next.config_profile_inbound_uuid)?.profile_uuid !== profileUUID
      ) {
        next.config_profile_inbound_uuid = '';
      }
      return next;
    });
  };

  const handleInboundChange = (inboundUUID) => {
    setFormData((previous) => {
      const next = { ...previous, config_profile_inbound_uuid: inboundUUID };
      if (inboundUUID) {
        const inbound = inboundByUuid.get(inboundUUID);
        if (inbound) {
          next.config_profile_uuid = inbound.profile_uuid;
        }
      }
      return next;
    });
  };

  const handleNodesChange = (nodeUUIDs) => {
    setField('nodes', nodeUUIDs);
  };

  const resolveHostNodes = (hostUUID) => {
    const assignedNodeUUIDs = Array.from(hostNodesMap.get(hostUUID) || []);
    if (assignedNodeUUIDs.length === 0) {
      return '-';
    }
    return assignedNodeUUIDs
      .map((nodeUUID) => nodeByUuid.get(nodeUUID)?.name || nodeUUID)
      .join(', ');
  };

  const getHostStatusLabel = (host) => (host.is_disabled ? 'DISABLED' : 'ACTIVE');

  const getHostStatusClass = (host) => (host.is_disabled ? 'users-status-disabled' : 'users-status-active');

  const getSortValue = (host, key) => {
    switch (key) {
      case 'remark':
        return String(host.remark || '').toLowerCase();
      case 'address':
        return String(host.address || '').toLowerCase();
      case 'port':
        return Number(host.port ?? 0);
      case 'status':
        return host.is_disabled ? 1 : 0;
      case 'tag':
        return String(host.tag || '').toLowerCase();
      case 'inbound':
        return String(resolveInboundTag(host) || '').toLowerCase();
      case 'profile':
        return String(resolveProfileName(host) || '').toLowerCase();
      case 'nodes':
        return String(resolveHostNodes(host.uuid) || '').toLowerCase();
      case 'securityLayer':
        return String(host.security_layer || 'DEFAULT').toUpperCase();
      case 'allowInsecure':
        return host.allow_insecure ? 1 : 0;
      case 'isHidden':
        return host.is_hidden ? 1 : 0;
      case 'viewPosition':
        return Number(host.view_position ?? 0);
      case 'sni':
        return String(host.sni || '').toLowerCase();
      case 'host':
        return String(host.host || '').toLowerCase();
      case 'path':
        return String(host.path || '').toLowerCase();
      case 'uuid':
        return String(host.uuid || '').toLowerCase();
      default:
        return '';
    }
  };

  const query = searchQuery.trim().toLowerCase();

  const filteredHosts = useMemo(() => {
    if (!query) return hosts;

    return hosts.filter((host) => {
      const pool = [
        host.remark,
        host.address,
        host.tag,
        host.uuid,
        host.security_layer,
        host.path,
        host.sni,
        host.host,
        resolveProfileName(host),
        resolveInboundTag(host),
        resolveHostNodes(host.uuid),
      ]
        .filter(Boolean)
        .map((value) => String(value).toLowerCase());

      return pool.some((value) => value.includes(query));
    });
  }, [hosts, query, profileNameByUuid, inboundByUuid, hostNodesMap, nodeByUuid]);

  const sortedHosts = useMemo(() => {
    if (!sortState.key) {
      return filteredHosts;
    }

    const directionFactor = sortState.direction === 'asc' ? 1 : -1;

    return filteredHosts
      .map((host, index) => ({ host, index }))
      .sort((a, b) => {
        const aValue = getSortValue(a.host, sortState.key);
        const bValue = getSortValue(b.host, sortState.key);

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
      .map((item) => item.host);
  }, [filteredHosts, sortState, profileNameByUuid, inboundByUuid, hostNodesMap, nodeByUuid]);

  const renderCell = (column, host) => {
    switch (column.key) {
      case 'remark':
        return (
          <td key={column.key}>
            <div className="users-name-cell">
              <span className={`users-presence-dot ${host.is_disabled ? 'is-idle' : 'is-online'}`}></span>
              <div className="users-name-text">
                <strong>{host.remark || '—'}</strong>
              </div>
            </div>
          </td>
        );
      case 'address':
        return <td key={column.key} className="users-cell-mono">{host.address || '—'}</td>;
      case 'port':
        return <td key={column.key} className="users-cell-center users-cell-mono">{host.port ?? '—'}</td>;
      case 'status':
        return (
          <td key={column.key} className="users-cell-center">
            <span className={`users-status-badge ${getHostStatusClass(host)}`}>
              <span className="users-status-dot"></span>
              {getHostStatusLabel(host)}
            </span>
          </td>
        );
      case 'tag':
        return <td key={column.key}>{host.tag || '—'}</td>;
      case 'inbound':
        return <td key={column.key}>{resolveInboundTag(host) || '—'}</td>;
      case 'profile':
        return <td key={column.key}>{resolveProfileName(host) || '—'}</td>;
      case 'nodes':
        return <td key={column.key} title={resolveHostNodes(host.uuid)}>{resolveHostNodes(host.uuid)}</td>;
      case 'securityLayer':
        return <td key={column.key} className="users-cell-center users-cell-mono">{(host.security_layer || 'DEFAULT').toUpperCase()}</td>;
      case 'allowInsecure':
        return <td key={column.key} className="users-cell-center">{host.allow_insecure ? 'Да' : 'Нет'}</td>;
      case 'isHidden':
        return <td key={column.key} className="users-cell-center">{host.is_hidden ? 'Да' : 'Нет'}</td>;
      case 'viewPosition':
        return <td key={column.key} className="users-cell-center users-cell-strong">{host.view_position ?? 0}</td>;
      case 'sni':
        return <td key={column.key} className="users-cell-mono">{host.sni || '—'}</td>;
      case 'host':
        return <td key={column.key} className="users-cell-mono">{host.host || '—'}</td>;
      case 'path':
        return <td key={column.key} className="users-cell-mono">{host.path || '—'}</td>;
      case 'uuid':
        return <td key={column.key} className="users-cell-mono users-cell-dim">{host.uuid || '—'}</td>;
      default:
        return <td key={column.key}>—</td>;
    }
  };

  const selectedALPNValues = useMemo(
    () => (formData.alpn || '').split(',').map((value) => value.trim()).filter(Boolean),
    [formData.alpn]
  );
  const alpnOptions = useMemo(() => {
    const options = new Set(ALPN_OPTIONS);
    selectedALPNValues.forEach((value) => options.add(value));
    return Array.from(options).map((value) => ({ value, label: value }));
  }, [selectedALPNValues]);
  const visibleHostEnabled = !formData.is_disabled;
  const currentHostName = editingHost?.remark || formData.remark || formData.address || 'Хост';

  if (loading) {
    return (
      <div className="empty-state">
        <div className="spinner"></div>
      </div>
    );
  }

  return (
    <div className="page">
      {error ? (
        <div className="alert alert-error">
          {error}
          <button className="btn-icon" onClick={loadHosts}>↻</button>
        </div>
      ) : null}

      <div className={`card users-list-card ${isFullscreen ? 'users-list-card-fullscreen' : ''}`}>
        <div className="users-list-header">
          <div className="users-list-header-main">
            <div className="users-action-icon users-action-cyan users-section-icon" aria-hidden="true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="3" />
                <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z" />
              </svg>
            </div>
            <div className="users-list-heading-stack">
              <h2 className="card-title">Hosts</h2>
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
                placeholder="Поиск hosts..."
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
              title={`Удалить выбранных: ${selectedHostUUIDs.size}`}
              aria-label={`Удалить выбранных: ${selectedHostUUIDs.size}`}
              disabled={selectedHostUUIDs.size === 0}
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <polyline points="3 6 5 6 21 6"></polyline>
                <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"></path>
              </svg>
              {selectedHostUUIDs.size > 0 && <span className="users-action-count">{selectedHostUUIDs.size}</span>}
            </button>
            <button
              type="button"
              className="users-action-icon users-action-info"
              onClick={loadHosts}
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
              onClick={openCreate}
              title="Добавить хост"
              aria-label="Добавить хост"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M12 5l0 14"></path>
                <path d="M5 12l14 0"></path>
              </svg>
            </button>
          </div>
        </div>

        <SharedRichTable
          columns={HOST_COLUMNS}
          rows={sortedHosts}
          getRowId={(host) => host.uuid}
          renderCell={renderCell}
          onRowClick={openEdit}
          getRowAriaLabel={(host) => `Редактировать host ${host.remark || ''}`}
          selectedRowIds={selectedHostUUIDs}
          setSelectedRowIds={setSelectedHostUUIDs}
          sortState={sortState}
          setSortState={setSortState}
          compactRows={compactRows}
          setCompactRows={setCompactRows}
          isFullscreen={isFullscreen}
          setIsFullscreen={setIsFullscreen}
          searchInputRef={searchInputRef}
          pageResetKey={searchQuery}
          pageSizeOptions={PAGE_SIZE_OPTIONS}
          pageSizeLabelId="hosts-rpp-label"
          tableAriaLabel="Таблица hosts с горизонтальной прокруткой"
          emptyState={(
            <div className="empty-state users-empty-state">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <circle cx="12" cy="12" r="10"></circle>
                <line x1="2" y1="12" x2="22" y2="12"></line>
                <path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"></path>
              </svg>
              <h3>Hosts не найдены</h3>
              <p>Проверьте фильтр или добавьте новый host</p>
              <button className="btn btn-primary" onClick={openCreate}>Создать host</button>
            </div>
          )}
        />
      </div>

      <SideDrawerPanel
        open={showModal}
        onClose={closeEditor}
        width="calc(100vw / 3)"
        className="host-editor-drawer"
        title={editingHost ? 'Редактирование хоста' : 'Новый хост'}
        subtitle={editingHost ? `UUID: ${editingHost.uuid}` : 'Создание нового хоста'}
        icon={(
          <svg stroke="currentColor" fill="currentColor" strokeWidth="0" viewBox="0 0 256 256" xmlns="http://www.w3.org/2000/svg">
            <path d="M224,128a8,8,0,0,1-8,8H128a8,8,0,0,1,0-16h88A8,8,0,0,1,224,128ZM128,72h88a8,8,0,0,0,0-16H128a8,8,0,0,0,0,16Zm88,112H128a8,8,0,0,0,0,16h88a8,8,0,0,0,0-16ZM82.34,42.34,56,68.69,45.66,58.34A8,8,0,0,0,34.34,69.66l16,16a8,8,0,0,0,11.32,0l32-32A8,8,0,0,0,82.34,42.34Zm0,64L56,132.69,45.66,122.34a8,8,0,0,0-11.32,11.32l16,16a8,8,0,0,0,11.32,0l32-32a8,8,0,0,0-11.32-11.32Zm0,64L56,196.69,45.66,186.34a8,8,0,0,0-11.32,11.32l16,16a8,8,0,0,0,11.32,0l32-32a8,8,0,0,0-11.32-11.32Z"></path>
          </svg>
        )}
        footer={(
          <div className="host-editor-footer-row">
            <button
              type="button"
              className="squad-editor-summary-btn cancel"
              onClick={closeEditor}
              disabled={savingHost}
            >
              Отмена
            </button>
            <button
              type="submit"
              form="host-editor-form"
              className="squad-editor-summary-btn save"
              disabled={savingHost}
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                <path d="M6 4h10l4 4v10a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12a2 2 0 0 1 2 -2"></path>
                <path d="M12 14m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0"></path>
                <path d="M14 4l0 4l-6 0l0 -4"></path>
              </svg>
              <span>{savingHost ? 'Сохранение...' : (editingHost ? 'Сохранить' : 'Создать')}</span>
            </button>
          </div>
        )}
      >
        <form id="host-editor-form" className="host-editor-form" onSubmit={handleSubmit}>
          <div className="host-editor-stack">
            <section className="host-editor-summary">
              <div className="host-editor-summary-main">
                <div className="host-editor-summary-identity">
                  <button type="button" className="host-editor-summary-icon" tabIndex={-1} aria-hidden="true">
                    <svg stroke="currentColor" fill="none" strokeWidth="2" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
                      <path d="M14 3v4a1 1 0 0 0 1 1h4"></path>
                      <path d="M17 21h-10a2 2 0 0 1 -2 -2v-14a2 2 0 0 1 2 -2h7l5 5v11a2 2 0 0 1 -2 2z"></path>
                    </svg>
                  </button>
                  <div className="host-editor-summary-title-wrap">
                    <h3 className="host-editor-summary-title" title={currentHostName}>{currentHostName}</h3>
                    <p className="host-editor-summary-subtitle">
                      {selectedInbound ? `${selectedInbound.profile_name} / ${selectedInbound.tag}` : 'Инбаунд не выбран'}
                    </p>
                  </div>
                </div>
                <div className="host-editor-summary-pills">
                  <span className="squad-editor-pill">{mode === 'basic' ? 'Основные' : 'Расширенные'}</span>
                  {formData.is_hidden ? <span className="squad-editor-pill muted">Скрыт</span> : null}
                </div>
              </div>

              <div className="host-editor-switch-row">
                <div className="host-editor-switch-copy">
                  <p className="host-editor-switch-title">Видимость хоста</p>
                  <p className="host-editor-switch-hint">
                    {visibleHostEnabled ? 'Хост видим и участвует в маршрутизации.' : 'Хост отключен и не участвует в маршрутизации.'}
                  </p>
                </div>
                <label className="host-switch" aria-label="Видимость хоста">
                  <input
                    type="checkbox"
                    checked={visibleHostEnabled}
                    onChange={(event) => setField('is_disabled', !event.target.checked)}
                  />
                  <span className="host-switch-slider"></span>
                </label>
              </div>
            </section>

            <div className="host-editor-tabs" role="tablist" aria-label="Режим редактирования хоста">
              <button
                type="button"
                className={`host-editor-tab ${mode === 'basic' ? 'active' : ''}`}
                onClick={() => setMode('basic')}
                role="tab"
                aria-selected={mode === 'basic'}
              >
                Основные
              </button>
              <button
                type="button"
                className={`host-editor-tab ${mode === 'advanced' ? 'active' : ''}`}
                onClick={() => setMode('advanced')}
                role="tab"
                aria-selected={mode === 'advanced'}
              >
                Расширенные
              </button>
            </div>

            {mode === 'basic' ? (
              <div className="host-editor-scroll">
                <section className="host-editor-card">
                  <div className="host-editor-card-header">
                    <div className="host-editor-card-icon host-editor-card-icon-teal">
                      <svg stroke="currentColor" fill="currentColor" strokeWidth="0" viewBox="0 0 256 256">
                        <path d="M243.31,136,144,36.69A15.86,15.86,0,0,0,132.69,32H40a8,8,0,0,0-8,8v92.69A15.86,15.86,0,0,0,36.69,144L136,243.31a16,16,0,0,0,22.63,0l84.68-84.68a16,16,0,0,0,0-22.63Zm-96,96L48,132.69V48h84.69L232,147.31ZM96,84A12,12,0,1,1,84,72,12,12,0,0,1,96,84Z"></path>
                      </svg>
                    </div>
                    <div className="host-editor-card-copy">
                      <h4 className="host-editor-card-title">Базовые параметры</h4>
                      <p className="host-editor-card-subtitle">Основные поля для идентификации и маршрутизации хоста.</p>
                    </div>
                  </div>

                  <div className="host-editor-grid">
                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">Примечание *</label>
                      <input
                        type="text"
                        className="form-input"
                        value={formData.remark}
                        onChange={(event) => setField('remark', event.target.value)}
                        required
                      />
                    </div>

                    <div className="host-editor-inline span-2">
                      <div className="form-group host-editor-field">
                        <label className="form-label">Адрес *</label>
                        <input
                          type="text"
                          className="form-input"
                          value={formData.address}
                          onChange={(event) => setField('address', event.target.value)}
                          placeholder="например, example.com"
                          required
                        />
                      </div>
                      <div className="form-group host-editor-field">
                        <label className="form-label">Порт *</label>
                        <input
                          type="number"
                          className="form-input"
                          value={formData.port}
                          min="1"
                          max="65535"
                          onChange={(event) => setField('port', event.target.value)}
                          required
                        />
                      </div>
                    </div>

                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">Tag</label>
                      <input
                        type="text"
                        className="form-input"
                        value={formData.tag}
                        onChange={(event) => setField('tag', event.target.value)}
                        placeholder="ROUTING_HOST"
                      />
                    </div>

                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">Ноды</label>
                      <AppMultiSelect
                        value={formData.nodes}
                        onChange={handleNodesChange}
                        options={nodeOptions}
                        placeholder={nodes.length === 0 ? 'Ноды отсутствуют' : 'Выберите ноды'}
                        description="Выберите ноды, которые направлены на этот хост."
                        leftSection={(
                          <svg stroke="currentColor" fill="none" strokeWidth="2" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M3 4m0 3a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v2a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3z"></path>
                            <path d="M3 12m0 3a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v2a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3z"></path>
                            <path d="M7 8l0 .01"></path>
                            <path d="M7 16l0 .01"></path>
                            <path d="M11 8h6"></path>
                            <path d="M11 16h6"></path>
                          </svg>
                        )}
                      />
                    </div>
                  </div>
                </section>

                <section className="host-editor-card">
                  <div className="host-editor-card-header">
                    <div className="host-editor-card-icon host-editor-card-icon-indigo">
                      <svg stroke="currentColor" fill="none" strokeWidth="2" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M3 4m0 3a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v2a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3z"></path>
                        <path d="M3 12m0 3a3 3 0 0 1 3 -3h12a3 3 0 0 1 3 3v2a3 3 0 0 1 -3 3h-12a3 3 0 0 1 -3 -3z"></path>
                        <path d="M7 8l0 .01"></path>
                        <path d="M7 16l0 .01"></path>
                        <path d="M11 8h6"></path>
                        <path d="M11 16h6"></path>
                      </svg>
                    </div>
                    <div className="host-editor-card-copy">
                      <h4 className="host-editor-card-title">Связи с профилем</h4>
                      <p className="host-editor-card-subtitle">Выбери профиль и инбаунд, которые будут связаны с хостом.</p>
                    </div>
                  </div>

                  <div className="host-editor-grid">
                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">Config Profile</label>
                      <AppSelect
                        value={formData.config_profile_uuid}
                        onChange={(event) => handleProfileChange(event.target.value)}
                      >
                        <option value="">No profile</option>
                        {profiles.map((profile) => (
                          <option key={profile.uuid} value={profile.uuid}>
                            {profile.name}
                          </option>
                        ))}
                      </AppSelect>
                    </div>

                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">Config Profile Inbound</label>
                      <AppSelect
                        value={formData.config_profile_inbound_uuid}
                        onChange={(event) => handleInboundChange(event.target.value)}
                      >
                        <option value="">No inbound</option>
                        {visibleInboundOptions.map((inbound) => (
                          <option key={inbound.uuid} value={inbound.uuid}>
                            {inbound.profile_name} / {inbound.tag}
                          </option>
                        ))}
                      </AppSelect>
                    </div>
                  </div>
                </section>
              </div>
            ) : (
              <div className="host-editor-scroll">
                <section className="host-editor-card">
                  <div className="host-editor-card-header">
                    <div className="host-editor-card-icon host-editor-card-icon-teal">
                      <svg stroke="currentColor" fill="currentColor" strokeWidth="0" viewBox="0 0 256 256">
                        <path d="M230.1,108.76,198.25,90.62c-.64-1.16-1.31-2.29-2-3.41l-.12-36A104.61,104.61,0,0,0,162,32L130,49.89c-1.34,0-2.69,0-4,0L94,32A104.58,104.58,0,0,0,59.89,51.25l-.16,36c-.7,1.12-1.37,2.26-2,3.41l-31.84,18.1a99.15,99.15,0,0,0,0,38.46l31.85,18.14c.64,1.16,1.31,2.29,2,3.41l.12,36A104.61,104.61,0,0,0,94,224l32-17.87c1.34,0,2.69,0,4,0L162,224a104.58,104.58,0,0,0,34.08-19.25l.16-36c.7-1.12,1.37-2.26,2-3.41l31.84-18.1A99.15,99.15,0,0,0,230.1,108.76ZM128,168a40,40,0,1,1,40-40A40,40,0,0,1,128,168Z"></path>
                      </svg>
                    </div>
                    <div className="host-editor-card-copy">
                      <h4 className="host-editor-card-title">Переопределения соединений</h4>
                      <p className="host-editor-card-subtitle">Настройки SNI, хоста, пути и TLS-параметров.</p>
                    </div>
                  </div>

                  <div className="host-editor-grid">
                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">SNI</label>
                      <input
                        type="text"
                        className="form-input"
                        value={formData.sni}
                        onChange={(event) => setField('sni', event.target.value)}
                        placeholder="SNI (например, example.com)"
                      />
                    </div>

                    <label className="host-editor-flag span-2">
                      <span className="host-editor-flag-copy">
                        <span className="host-editor-flag-title">Переопределить SNI из адреса</span>
                      </span>
                      <span className="host-switch host-switch-sm">
                        <input
                          type="checkbox"
                          checked={!!formData.override_sni_from_address}
                          onChange={(event) => setField('override_sni_from_address', event.target.checked)}
                        />
                        <span className="host-switch-slider"></span>
                      </span>
                    </label>

                    <label className="host-editor-flag span-2">
                      <span className="host-editor-flag-copy">
                        <span className="host-editor-flag-title">Оставить SNI пустым</span>
                      </span>
                      <span className="host-switch host-switch-sm">
                        <input
                          type="checkbox"
                          checked={!!formData.keep_sni_blank}
                          onChange={(event) => setField('keep_sni_blank', event.target.checked)}
                        />
                        <span className="host-switch-slider"></span>
                      </span>
                    </label>

                    <div className="host-editor-inline span-2">
                      <div className="form-group host-editor-field">
                        <label className="form-label">Host</label>
                        <input
                          type="text"
                          className="form-input"
                          value={formData.host}
                          onChange={(event) => setField('host', event.target.value)}
                          placeholder="example.com"
                        />
                      </div>
                      <div className="form-group host-editor-field">
                        <label className="form-label">Path</label>
                        <input
                          type="text"
                          className="form-input"
                          value={formData.path}
                          onChange={(event) => setField('path', event.target.value)}
                          placeholder="/ws"
                        />
                      </div>
                    </div>

                    <div className="form-group host-editor-field">
                      <label className="form-label">Security Layer</label>
                      <AppSelect
                        value={formData.security_layer}
                        onChange={(event) => setField('security_layer', event.target.value)}
                      >
                        {SECURITY_LAYER_OPTIONS.map((option) => (
                          <option key={option.value} value={option.value}>{option.label}</option>
                        ))}
                      </AppSelect>
                    </div>

                    <div className="form-group host-editor-field">
                      <label className="form-label">Fingerprint</label>
                      <AppSelect
                        value={formData.fingerprint}
                        onChange={(event) => setField('fingerprint', event.target.value)}
                      >
                        <option value="">No fingerprint</option>
                        {FINGERPRINT_OPTIONS.map((option) => (
                          <option key={option} value={option}>{option}</option>
                        ))}
                      </AppSelect>
                    </div>

                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">ALPN</label>
                      <AppMultiSelect
                        value={selectedALPNValues}
                        onChange={setALPNValues}
                        options={alpnOptions}
                        placeholder="Выберите ALPN"
                      />
                    </div>

                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">Vless Route ID</label>
                      <input
                        type="number"
                        className="form-input"
                        value={formData.vless_route_id}
                        onChange={(event) => setField('vless_route_id', event.target.value)}
                        placeholder="От 1 до 65535"
                      />
                    </div>
                  </div>
                </section>

                <section className="host-editor-card">
                  <div className="host-editor-card-header">
                    <div className="host-editor-card-icon host-editor-card-icon-violet">
                      <svg stroke="currentColor" fill="none" strokeWidth="2" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M3 20h7"></path>
                        <path d="M14 20h7"></path>
                        <path d="M10 20a2 2 0 1 0 4 0a2 2 0 0 0 -4 0"></path>
                        <path d="M12 16v2"></path>
                        <path d="M8 16.004h-1.343c-2.572 -.004 -4.657 -2.011 -4.657 -4.487c0 -2.475 2.085 -4.482 4.657 -4.482c.393 -1.762 1.794 -3.2 3.675 -3.773c1.88 -.572 3.956 -.193 5.444 1c1.488 1.19 2.162 3.007 1.77 4.769h.99c1.913 0 3.464 1.56 3.464 3.486c0 1.927 -1.551 3.487 -3.465 3.487h-2.535"></path>
                      </svg>
                    </div>
                    <div className="host-editor-card-copy">
                      <h4 className="host-editor-card-title">Xray Json &amp; Raw</h4>
                      <p className="host-editor-card-subtitle">Параметры raw-полей и скрытия хоста.</p>
                    </div>
                  </div>

                  <div className="host-editor-grid">
                    <label className="host-editor-flag span-2">
                      <span className="host-editor-flag-copy">
                        <span className="host-editor-flag-title">Скрыть хост</span>
                      </span>
                      <span className="host-switch host-switch-sm">
                        <input
                          type="checkbox"
                          checked={!!formData.is_hidden}
                          onChange={(event) => setField('is_hidden', event.target.checked)}
                        />
                        <span className="host-switch-slider"></span>
                      </span>
                    </label>

                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">Шаблон Xray JSON</label>
                      <input
                        type="text"
                        className="form-input"
                        value={formData.xray_json_template_uuid}
                        onChange={(event) => setField('xray_json_template_uuid', event.target.value)}
                        placeholder="UUID шаблона"
                      />
                    </div>

                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">xhttp_extra_params</label>
                      <input
                        type="text"
                        className="form-input"
                        value={formData.xhttp_extra_params}
                        onChange={(event) => setField('xhttp_extra_params', event.target.value)}
                      />
                    </div>

                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">mux_params</label>
                      <input
                        type="text"
                        className="form-input"
                        value={formData.mux_params}
                        onChange={(event) => setField('mux_params', event.target.value)}
                      />
                    </div>

                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">sockopt_params</label>
                      <input
                        type="text"
                        className="form-input"
                        value={formData.sockopt_params}
                        onChange={(event) => setField('sockopt_params', event.target.value)}
                      />
                    </div>
                  </div>
                </section>

                <section className="host-editor-card">
                  <div className="host-editor-card-header">
                    <div className="host-editor-card-icon host-editor-card-icon-indigo">
                      <svg stroke="currentColor" fill="none" strokeWidth="2" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round">
                        <path d="M9.183 6.117a6 6 0 1 0 4.511 3.986"></path>
                        <path d="M14.813 17.883a6 6 0 1 0 -4.496 -3.954"></path>
                      </svg>
                    </div>
                    <div className="host-editor-card-copy">
                      <h4 className="host-editor-card-title">Прочие настройки</h4>
                      <p className="host-editor-card-subtitle">Флаги безопасности и параметры для клиентов.</p>
                    </div>
                  </div>

                  <div className="host-editor-grid">
                    <div className="form-group host-editor-field span-2">
                      <label className="form-label">Server Description</label>
                      <input
                        type="text"
                        className="form-input"
                        value={formData.server_description}
                        onChange={(event) => setField('server_description', event.target.value)}
                        placeholder="Максимум 30 символов"
                      />
                    </div>

                    <label className="host-editor-flag span-2">
                      <span className="host-editor-flag-copy">
                        <span className="host-editor-flag-title">Перемешать хост</span>
                      </span>
                      <span className="host-switch host-switch-sm">
                        <input
                          type="checkbox"
                          checked={!!formData.shuffle_host}
                          onChange={(event) => setField('shuffle_host', event.target.checked)}
                        />
                        <span className="host-switch-slider"></span>
                      </span>
                    </label>

                    <label className="host-editor-flag span-2">
                      <span className="host-editor-flag-copy">
                        <span className="host-editor-flag-title">Разрешить небезопасные</span>
                      </span>
                      <span className="host-switch host-switch-sm">
                        <input
                          type="checkbox"
                          checked={!!formData.allow_insecure}
                          onChange={(event) => setField('allow_insecure', event.target.checked)}
                        />
                        <span className="host-switch-slider"></span>
                      </span>
                    </label>

                    <label className="host-editor-flag span-2">
                      <span className="host-editor-flag-copy">
                        <span className="host-editor-flag-title">Включение x25519mlkem768</span>
                      </span>
                      <span className="host-switch host-switch-sm">
                        <input
                          type="checkbox"
                          checked={!!formData.mihomo_x25519}
                          onChange={(event) => setField('mihomo_x25519', event.target.checked)}
                        />
                        <span className="host-switch-slider"></span>
                      </span>
                    </label>
                  </div>
                </section>
              </div>
            )}
          </div>
        </form>
      </SideDrawerPanel>
    </div>
  );
}

export default Hosts;
