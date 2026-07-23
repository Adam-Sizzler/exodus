"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.REMNAWAVE_BYPASS_HTTPS_RESTRCTIONS = exports.REMNAWAVE_REAL_IP_HEADER = exports.REMNAWAVE_CLIENT_TYPE_BROWSER = exports.REMNAWAVE_CLIENT_TYPE_HEADER = void 0;
exports.REMNAWAVE_CLIENT_TYPE_HEADER = 'X-Remnawave-Client-Type';
exports.REMNAWAVE_CLIENT_TYPE_BROWSER = 'browser';
exports.REMNAWAVE_REAL_IP_HEADER = 'x-remnawave-real-ip';
exports.REMNAWAVE_BYPASS_HTTPS_RESTRCTIONS = {
    'x-forwarded-proto': 'https',
    'x-forwarded-for': '127.0.0.1',
};
