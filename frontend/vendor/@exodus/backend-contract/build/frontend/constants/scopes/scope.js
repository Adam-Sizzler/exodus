"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.buildEndpointScope = exports.buildActionScope = exports.buildResourceScope = exports.normalizeControllerUrl = exports.SCOPE_WILDCARD = exports.SCOPE_ACTION = void 0;
exports.SCOPE_ACTION = {
    READ: 'read',
    WRITE: 'write',
};
exports.SCOPE_WILDCARD = '*';
const normalizeControllerUrl = (controllerUrl) => {
    const segments = controllerUrl
        .split('/')
        .filter(Boolean)
        .map((segment) => (segment.startsWith(':') ? `{${segment.slice(1)}}` : segment));
    return `/${segments.join('/')}`;
};
exports.normalizeControllerUrl = normalizeControllerUrl;
const buildResourceScope = (resource) => `${resource}:${exports.SCOPE_WILDCARD}`;
exports.buildResourceScope = buildResourceScope;
const buildActionScope = (resource, action) => `${resource}:${action}`;
exports.buildActionScope = buildActionScope;
const buildEndpointScope = (resource, scopeSlug) => `${resource}:${scopeSlug}`;
exports.buildEndpointScope = buildEndpointScope;
