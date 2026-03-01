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
      return await apiRequest('/api/v1/config-profiles-with-inbounds', { method: 'GET' });
    } catch (err) {
      // Fallback to regular profiles endpoint
      return await apiRequest('/api/v1/config-profiles', { method: 'GET' });
    }
  },

  // Get all config profiles
  getAll: () => apiRequest('/api/v1/config-profiles', { method: 'GET' }),

  // Get single config profile by UUID
  getById: (uuid) => apiRequest(`/api/v1/config-profiles/${uuid}`, { method: 'GET' }),

  // Create new config profile
  create: (data) => apiRequest('/api/v1/config-profiles', {
    method: 'POST',
    body: JSON.stringify(data),
  }),

  // Update config profile (partial update)
  update: (uuid, data) => apiRequest(`/api/v1/config-profiles/${uuid}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),

  // Delete config profile
  delete: (uuid) => apiRequest(`/api/v1/config-profiles/${uuid}`, {
    method: 'DELETE',
  }),
  reorder: (orderedUuids) => apiRequest('/api/v1/config-profiles/reorder', {
    method: 'POST',
    body: JSON.stringify({ ordered_uuids: orderedUuids }),
  }),
};

// Nodes API
export const nodesApi = {
  // Get nodes with config profile info (fallback to regular nodes if not available)
  getAllWithConfig: async () => {
    try {
      return await apiRequest('/api/v1/nodes-with-config', { method: 'GET' });
    } catch (err) {
      // Fallback to regular nodes endpoint
      return await apiRequest('/api/v1/nodes', { method: 'GET' });
    }
  },

  getAll: () => apiRequest('/api/v1/nodes', { method: 'GET' }),
  getById: (uuid) => apiRequest(`/api/v1/nodes/${uuid}`, { method: 'GET' }),
  create: (data) => apiRequest('/api/v1/nodes', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (uuid, data) => apiRequest(`/api/v1/nodes/${uuid}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
  delete: (uuid) => apiRequest(`/api/v1/nodes/${uuid}`, {
    method: 'DELETE',
  }),
  deleteMany: (uuids) => apiRequest(`/api/v1/nodes?uuids=${encodeURIComponent(uuids.join(','))}`, {
    method: 'DELETE',
  }),
  reorder: (orderedUuids) => apiRequest('/api/v1/nodes/reorder', {
    method: 'POST',
    body: JSON.stringify({ ordered_uuids: orderedUuids }),
  }),

  // Inbound assignments for nodes
  getInboundAssignments: (nodeUuid) => apiRequest(`/api/v1/inbound-assignments?node_uuid=${nodeUuid}`, { method: 'GET' }),
  setInboundAssignments: (nodeUuid, inboundUuids) => apiRequest('/api/v1/inbound-assignments', {
    method: 'POST',
    body: JSON.stringify({
      node_uuid: nodeUuid,
      inbound_uuids: inboundUuids,
    }),
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
  getAll: () => apiRequest('/api/v1/users-list', { method: 'GET' }),
  getById: (uuid) => apiRequest(`/api/v1/users-list/${uuid}`, { method: 'GET' }),
  create: (data) => apiRequest('/api/v1/users-list/create', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (uuid, data) => apiRequest(`/api/v1/users-list/${uuid}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
  delete: (uuid) => apiRequest(`/api/v1/users-list/${uuid}`, {
    method: 'DELETE',
  }),
  deleteMany: (uuids) => apiRequest(`/api/v1/users-list?uuids=${encodeURIComponent(uuids.join(','))}`, {
    method: 'DELETE',
  }),
  reorder: (orderedUuids) => apiRequest('/api/v1/users-list/reorder', {
    method: 'POST',
    body: JSON.stringify({ ordered_uuids: orderedUuids }),
  }),
};

// Hosts API
export const hostsApi = {
  getAll: () => apiRequest('/api/v1/hosts', { method: 'GET' }),
  getById: (uuid) => apiRequest(`/api/v1/hosts/${uuid}`, { method: 'GET' }),
  create: (data) => apiRequest('/api/v1/hosts', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (uuid, data) => apiRequest(`/api/v1/hosts/${uuid}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
  delete: (uuid) => apiRequest(`/api/v1/hosts/${uuid}`, {
    method: 'DELETE',
  }),
  deleteMany: (uuids) => apiRequest(`/api/v1/hosts?uuids=${encodeURIComponent(uuids.join(','))}`, {
    method: 'DELETE',
  }),
  getNodeAssignments: (hostUuid = '', nodeUuid = '') => {
    const params = new URLSearchParams();
    if (hostUuid) {
      params.set('host_uuid', hostUuid);
    }
    if (nodeUuid) {
      params.set('node_uuid', nodeUuid);
    }
    const query = params.toString();
    return apiRequest(`/api/v1/hosts-to-nodes${query ? `?${query}` : ''}`, { method: 'GET' });
  },
  setNodeAssignments: (hostUuid, nodeUuids) => apiRequest('/api/v1/hosts-to-nodes', {
    method: 'POST',
    body: JSON.stringify({
      host_uuid: hostUuid,
      node_uuids: nodeUuids,
    }),
  }),
  deleteNodeAssignments: (hostUuid, nodeUuids = []) => apiRequest('/api/v1/hosts-to-nodes', {
    method: 'DELETE',
    body: JSON.stringify({
      host_uuid: hostUuid,
      node_uuids: nodeUuids,
    }),
  }),
  reorder: (orderedUuids) => apiRequest('/api/v1/hosts/reorder', {
    method: 'POST',
    body: JSON.stringify({ ordered_uuids: orderedUuids }),
  }),
};

// Subscription settings API
export const subscriptionSettingsApi = {
  get: () => apiRequest('/api/v1/subscription-settings', { method: 'GET' }),
  getById: (uuid) => apiRequest(`/api/v1/subscription-settings/${uuid}`, { method: 'GET' }),
  update: (uuid, data) => apiRequest(`/api/v1/subscription-settings/${uuid}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
};

// Subscription templates API
export const templatesApi = {
  getAll: (templateType) =>
    apiRequest(
      templateType
        ? `/api/v1/templates?template_type=${encodeURIComponent(templateType)}`
        : '/api/v1/templates',
      { method: 'GET' }
    ),
  getById: (uuid) => apiRequest(`/api/v1/templates/${uuid}`, { method: 'GET' }),
  create: (data) => apiRequest('/api/v1/templates', {
    method: 'POST',
    body: JSON.stringify(data),
  }),
  update: (uuid, data) => apiRequest(`/api/v1/templates/${uuid}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  }),
  delete: (uuid) => apiRequest(`/api/v1/templates/${uuid}`, {
    method: 'DELETE',
  }),
  reorder: (orderedUuids) => apiRequest('/api/v1/templates/reorder', {
    method: 'POST',
    body: JSON.stringify({ ordered_uuids: orderedUuids }),
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
