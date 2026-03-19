"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.CERBERUS_BYPASS_HTTPS_RESTRCTIONS = exports.CERBERUS_REAL_IP_HEADER = exports.CERBERUS_CLIENT_TYPE_BROWSER = exports.CERBERUS_CLIENT_TYPE_HEADER = void 0;
exports.CERBERUS_CLIENT_TYPE_HEADER = 'X-Cerberus-Client-Type';
exports.CERBERUS_CLIENT_TYPE_BROWSER = 'browser';
exports.CERBERUS_REAL_IP_HEADER = 'x-cerberus-real-ip';
exports.CERBERUS_BYPASS_HTTPS_RESTRCTIONS = {
    'x-forwarded-proto': 'https',
    'x-forwarded-for': '127.0.0.1',
};
