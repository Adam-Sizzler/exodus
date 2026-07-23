"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EXODUS_BYPASS_HTTPS_RESTRCTIONS = exports.EXODUS_REAL_IP_HEADER = exports.EXODUS_CLIENT_TYPE_BROWSER = exports.EXODUS_CLIENT_TYPE_HEADER = void 0;
exports.EXODUS_CLIENT_TYPE_HEADER = 'X-Exodus-Client-Type';
exports.EXODUS_CLIENT_TYPE_BROWSER = 'browser';
exports.EXODUS_REAL_IP_HEADER = 'x-exodus-real-ip';
exports.EXODUS_BYPASS_HTTPS_RESTRCTIONS = {
    'x-forwarded-proto': 'https',
    'x-forwarded-for': '127.0.0.1',
};
