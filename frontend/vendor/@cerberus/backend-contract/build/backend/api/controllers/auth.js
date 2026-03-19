"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AUTH_ROUTES = exports.AUTH_CONTROLLER = void 0;
exports.AUTH_CONTROLLER = 'auth';
exports.AUTH_ROUTES = {
    LOGIN: 'login',
    REGISTER: 'register',
    GET_STATUS: 'status',
    OAUTH2: {
        TELEGRAM_CALLBACK: 'oauth2/tg/callback',
        AUTHORIZE: 'oauth2/authorize',
        CALLBACK: 'oauth2/callback',
    },
    PASSKEY: {
        GET_AUTHENTICATION_OPTIONS: 'passkey/authentication/options', // get
        VERIFY_AUTHENTICATION: 'passkey/authentication/verify', // post
    },
};
