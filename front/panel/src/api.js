// API Client for v2ray-stat backend
// Use relative path for API requests - Vite proxy will forward to backend.
const API_BASE_URL = '';
const API_TOKEN = import.meta.env.VITE_API_TOKEN || '';

// Helper function for API requests
async function apiRequest(endpoint, options = {}) {
  const url = `${API_BASE_URL}${endpoint}`;
  const { suppressAuthEvent = false, ...fetchOptions } = options;

  const headers = {
    'Content-Type': 'application/json',
    ...fetchOptions.headers,
  };

  // Add authorization token if set
  if (API_TOKEN) {
    headers['Authorization'] = `Bearer ${API_TOKEN}`;
  }

  const config = {
    credentials: 'include',
    ...fetchOptions,
    headers,
  };

  try {
    const response = await fetch(url, config);

    const contentType = response.headers.get('content-type') || '';
    let data = null;

    if (contentType.includes('application/json')) {
      data = await response.json();
    } else if (contentType.includes('text/html')) {
      throw new Error(`Endpoint ${endpoint} returned HTML instead of JSON. This endpoint may not exist on the backend.`);
    } else {
      const text = await response.text();
      data = text ? { message: text } : {};
    }

    if (!response.ok) {
      const error = new Error(data?.error || data?.message || `HTTP ${response.status}`);
      error.status = response.status;
      error.endpoint = endpoint;

      if (response.status === 401 && !suppressAuthEvent) {
        window.dispatchEvent(new CustomEvent('auth-expired', { detail: { endpoint } }));
      }

      throw error;
    }

    return data ?? {};
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
}

export const authApi = {
  bootstrap: () => apiRequest('/api/v1/auth/bootstrap', { method: 'GET', suppressAuthEvent: true }),
  setup: (payload) => apiRequest('/api/v1/auth/setup', {
    method: 'POST',
    body: JSON.stringify(payload),
    suppressAuthEvent: true,
  }),
  login: (payload) => apiRequest('/api/v1/auth/login', {
    method: 'POST',
    body: JSON.stringify(payload),
    suppressAuthEvent: true,
  }),
  me: () => apiRequest('/api/v1/auth/me', { method: 'GET', suppressAuthEvent: true }),
  logout: () => apiRequest('/api/v1/auth/logout', {
    method: 'POST',
    suppressAuthEvent: true,
  }),
};

// Config Profiles API
export const configProfilesApi = {
  // Get all config profiles with inbounds
  getAllWithInbounds: async () => {
    try {
      const data = await apiRequest('/api/config-profiles', { method: 'GET' });
      return {
        profiles: Array.isArray(data?.response?.configProfiles) ? data.response.configProfiles : [],
        response: data?.response,
      };
    } catch (err) {
      // Fallback to regular profiles endpoint
      return await apiRequest('/api/v1/config-profiles-with-inbounds', { method: 'GET' });
    }
  },

  // Get all config profiles
  getAll: () => apiRequest('/api/config-profiles', { method: 'GET' }),

  // Get single config profile by UUID
  getById: (uuid) => apiRequest(`/api/config-profiles/${uuid}`, { method: 'GET' }),

  // Create new config profile
  create: (data) => apiRequest('/api/config-profiles', {
    method: 'POST',
    body: JSON.stringify(data),
  }),

  // Update config profile (partial update)
  update: (uuid, data) => apiRequest('/api/config-profiles', {
    method: 'PATCH',
    body: JSON.stringify({ uuid, ...data }),
  }),

  // Delete config profile
  delete: (uuid) => apiRequest(`/api/config-profiles/${uuid}`, {
    method: 'DELETE',
  }),
  reorder: (items) => apiRequest('/api/config-profiles/actions/reorder', {
    method: 'POST',
    body: JSON.stringify({ items }),
  }),
};

export const configProfileSnippetsApi = {
  getAll: () => apiRequest('/api/config-profiles/snippets', { method: 'GET' }),
  create: (data) => apiRequest('/api/config-profiles/snippets', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (data) => apiRequest('/api/config-profiles/snippets', {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
  delete: (name) => apiRequest('/api/config-profiles/snippets', {
    method: 'DELETE',
    body: JSON.stringify({ name }),
  }),
};

// Nodes API
export const nodesApi = {
  getAll: () => apiRequest('/api/nodes', { method: 'GET' }),
  getById: (uuid) => apiRequest(`/api/nodes/${uuid}`, { method: 'GET' }),
  create: (data) => apiRequest('/api/nodes', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (data) => apiRequest('/api/nodes', {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
  delete: (uuid) => apiRequest(`/api/nodes/${uuid}`, {
    method: 'DELETE',
  }),
  enable: (uuid) => apiRequest(`/api/nodes/${uuid}/actions/enable`, {
    method: 'POST',
  }),
  disable: (uuid) => apiRequest(`/api/nodes/${uuid}/actions/disable`, {
    method: 'POST',
  }),
  restart: (uuid) => apiRequest(`/api/nodes/${uuid}/actions/restart`, {
    method: 'POST',
  }),
  resetTraffic: (uuid) => apiRequest(`/api/nodes/${uuid}/actions/reset-traffic`, {
    method: 'POST',
  }),
  restartAll: (forceRestart = false) => apiRequest('/api/nodes/actions/restart-all', {
    method: 'POST',
    body: JSON.stringify({ forceRestart }),
  }),
  deleteMany: (uuids, action = 'DISABLE') => apiRequest('/api/nodes/bulk-actions', {
    method: 'POST',
    body: JSON.stringify({ uuids, action }),
  }),
  bulkAction: (uuids, action) => apiRequest('/api/nodes/bulk-actions', {
    method: 'POST',
    body: JSON.stringify({ uuids, action }),
  }),
  profileModification: (uuids, configProfile) => apiRequest('/api/nodes/bulk-actions/profile-modification', {
    method: 'POST',
    body: JSON.stringify({ uuids, configProfile }),
  }),
  reorder: (items) => apiRequest('/api/nodes/actions/reorder', {
    method: 'POST',
    body: JSON.stringify({ nodes: items }),
  }),
  getTags: () => apiRequest('/api/nodes/tags', { method: 'GET' }),
  getAllWithConfig: async () => {
    return await apiRequest('/api/nodes', { method: 'GET' });
  },
  getInboundAssignments: (nodeUuid) => apiRequest(`/api/v1/inbound-assignments?node_uuid=${nodeUuid}`, { method: 'GET' }),
  setInboundAssignments: (nodeUuid, inboundUuids) => apiRequest('/api/v1/inbound-assignments', {
    method: 'POST',
    body: JSON.stringify({
      node_uuid: nodeUuid,
      inbound_uuids: inboundUuids,
    }),
  }),
  deleteManyLegacy: (uuids) => apiRequest(`/api/v1/nodes?uuids=${encodeURIComponent(uuids.join(','))}`, {
    method: 'DELETE',
  }),
};

// Internal Squads API
export const squadsApi = {
  // Get all squads summary
  getAllSummary: () => apiRequest('/api/v1/squads-summary', { method: 'GET' }),

  // Get squad details (with inbounds and members)
  getDetails: (uuid) => apiRequest(`/api/v1/squad-details/${uuid}`, { method: 'GET' }),

  // Get all squads
  getAll: () => apiRequest('/api/v1/internal-squads', { method: 'GET' }),

  // Create squad
  create: (data) => apiRequest('/api/v1/internal-squads', {
    method: 'POST',
    body: JSON.stringify(data),
  }),

  // Update squad
  update: (uuid, data) => apiRequest(`/api/v1/internal-squads/${uuid}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),

  // Delete squad
  delete: (uuid) => apiRequest(`/api/v1/internal-squads/${uuid}`, {
    method: 'DELETE',
  }),
  reorder: (orderedUuids) => apiRequest('/api/v1/internal-squads/reorder', {
    method: 'POST',
    body: JSON.stringify({ ordered_uuids: orderedUuids }),
  }),

  // Squad inbounds
  getInbounds: (squadUuid) => apiRequest(`/api/v1/squad-inbounds?squad_uuid=${squadUuid}`, { method: 'GET' }),
  setInbounds: (squadUuid, inboundUuids) => apiRequest('/api/v1/squad-inbounds', {
    method: 'POST',
    body: JSON.stringify({
      squad_uuid: squadUuid,
      inbound_uuids: inboundUuids,
    }),
  }),

  // Squad members
  getMembers: (squadUuid) => apiRequest(`/api/v1/squad-members?squad_uuid=${squadUuid}`, { method: 'GET' }),
  setMembers: (squadUuid, userIds) => apiRequest('/api/v1/squad-members', {
    method: 'POST',
    body: JSON.stringify({
      squad_uuid: squadUuid,
      user_ids: userIds,
    }),
  }),
};

// Users API
export const usersApi = {
  getAll: () => apiRequest('/api/users', { method: 'GET' }),
  getById: (uuid) => apiRequest(`/api/users/${uuid}`, { method: 'GET' }),
  create: (data) => apiRequest('/api/users', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (uuid, data) => apiRequest('/api/users', {
    method: 'PATCH',
    body: JSON.stringify({ uuid, ...data }),
  }),
  delete: (uuid) => apiRequest(`/api/users/${uuid}`, {
    method: 'DELETE',
  }),
  deleteMany: (uuids) => apiRequest('/api/users/bulk/delete', {
    method: 'POST',
    body: JSON.stringify({ uuids }),
  }),
  enable: (uuid) => apiRequest(`/api/users/${uuid}/actions/enable`, { method: 'POST' }),
  disable: (uuid) => apiRequest(`/api/users/${uuid}/actions/disable`, { method: 'POST' }),
  resetTraffic: (uuid) => apiRequest(`/api/users/${uuid}/actions/reset-traffic`, { method: 'POST' }),
  revokeSubscription: (uuid) => apiRequest(`/api/users/${uuid}/actions/revoke`, { method: 'POST' }),
  getTags: () => apiRequest('/api/users/tags', { method: 'GET' }),
};

// Hosts API
export const hostsApi = {
  getAll: () => apiRequest('/api/hosts', { method: 'GET' }),
  getById: (uuid) => apiRequest(`/api/hosts/${uuid}`, { method: 'GET' }),
  create: (data) => apiRequest('/api/hosts', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (data) => apiRequest('/api/hosts', {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
  delete: (uuid) => apiRequest(`/api/hosts/${uuid}`, {
    method: 'DELETE',
  }),
  deleteMany: (uuids) => apiRequest('/api/hosts/bulk/delete', {
    method: 'POST',
    body: JSON.stringify({ uuids }),
  }),
  bulkEnable: (uuids) => apiRequest('/api/hosts/bulk/enable', {
    method: 'POST',
    body: JSON.stringify({ uuids }),
  }),
  bulkDisable: (uuids) => apiRequest('/api/hosts/bulk/disable', {
    method: 'POST',
    body: JSON.stringify({ uuids }),
  }),
  setInbound: (uuids, configProfileUuid, configProfileInboundUuid) => apiRequest('/api/hosts/bulk/set-inbound', {
    method: 'POST',
    body: JSON.stringify({ uuids, configProfileUuid, configProfileInboundUuid }),
  }),
  setPort: (uuids, port) => apiRequest('/api/hosts/bulk/set-port', {
    method: 'POST',
    body: JSON.stringify({ uuids, port }),
  }),
  reorder: (items) => apiRequest('/api/hosts/actions/reorder', {
    method: 'POST',
    body: JSON.stringify({ hosts: items }),
  }),
  getTags: () => apiRequest('/api/hosts/tags', { method: 'GET' }),
};

// Subscription settings API
export const subscriptionSettingsApi = {
  get: () => apiRequest('/api/subscription-settings', { method: 'GET' }),
  update: (data) => apiRequest('/api/subscription-settings', {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
};

// Subscription templates API
export const templatesApi = {
  getAll: () =>
    apiRequest('/api/subscription-templates', { method: 'GET' }),
  getById: (uuid) => apiRequest(`/api/subscription-templates/${uuid}`, { method: 'GET' }),
  create: (data) => apiRequest('/api/subscription-templates', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (data) => apiRequest('/api/subscription-templates', {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
  delete: (uuid) => apiRequest(`/api/subscription-templates/${uuid}`, {
    method: 'DELETE',
  }),
  reorder: (items) => apiRequest('/api/subscription-templates/actions/reorder', {
    method: 'POST',
    body: JSON.stringify({ items }),
  }),
};

export const panelSettingsApi = {
  getSettings: () => apiRequest('/api/v1/panel/settings', { method: 'GET' }),
  updateSettings: (payload) => apiRequest('/api/v1/panel/settings', {
    method: 'PATCH',
    body: JSON.stringify(payload),
  }),
  getApiTokens: () => apiRequest('/api/v1/panel/api-tokens', { method: 'GET' }),
  createApiToken: (tokenName) => apiRequest('/api/v1/panel/api-tokens', {
    method: 'POST',
    body: JSON.stringify({ token_name: tokenName }),
  }),
  deleteApiToken: (uuid) => apiRequest(`/api/v1/panel/api-tokens/${uuid}`, {
    method: 'DELETE',
  }),
};

// Health check with timeout to avoid long hangs on unreachable backend
export const healthCheck = async (timeoutMs = 3000) => {
  const controller = new AbortController();
  const timeoutId = globalThis.setTimeout(() => controller.abort(), timeoutMs);

  try {
    return await apiRequest('/api/health', {
      method: 'GET',
      signal: controller.signal,
    });
  } finally {
    globalThis.clearTimeout(timeoutId);
  }
};
