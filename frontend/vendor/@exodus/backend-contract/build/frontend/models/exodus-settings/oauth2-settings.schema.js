"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.Oauth2SettingsSchema = void 0;
const zod_1 = __importDefault(require("zod"));
exports.Oauth2SettingsSchema = zod_1.default.object({
    github: zod_1.default.object({
        enabled: zod_1.default.boolean(),
        clientId: zod_1.default.nullable(zod_1.default.string()),
        clientSecret: zod_1.default.nullable(zod_1.default.string()),
        allowedEmails: zod_1.default.array(zod_1.default.string()),
    }),
    pocketid: zod_1.default.object({
        enabled: zod_1.default.boolean(),
        clientId: zod_1.default.nullable(zod_1.default.string()),
        clientSecret: zod_1.default.nullable(zod_1.default.string()),
        frontendDomain: zod_1.default.nullable(zod_1.default.string().refine((val) => {
            const fqdnRegex = /(?=^.{4,253}$)(^((?!-)[a-zA-Z0-9-]{0,62}[a-zA-Z0-9]\.)+[a-zA-Z]{2,63}$)/;
            if (fqdnRegex.test(val)) {
                return true;
            }
            return false;
        }, {
            message: 'Must be a valid fully qualified domain name (FQDN), e.g. "docs.ex"',
        })),
        plainDomain: zod_1.default.nullable(zod_1.default.string().refine((val) => {
            const fqdnRegex = /(?=^.{4,253}$)(^((?!-)[a-zA-Z0-9-]{0,62}[a-zA-Z0-9]\.)+[a-zA-Z]{2,63}$)/;
            if (fqdnRegex.test(val)) {
                return true;
            }
            return false;
        }, {
            message: 'Must be a valid fully qualified domain name (FQDN), e.g. "docs.ex"',
        })),
        allowedEmails: zod_1.default.array(zod_1.default.string()),
    }),
    yandex: zod_1.default.object({
        enabled: zod_1.default.boolean(),
        clientId: zod_1.default.nullable(zod_1.default.string()),
        clientSecret: zod_1.default.nullable(zod_1.default.string()),
        allowedEmails: zod_1.default.array(zod_1.default.string()),
    }),
    keycloak: zod_1.default
        .object({
        enabled: zod_1.default.boolean(),
        realm: zod_1.default.nullable(zod_1.default.string()),
        clientId: zod_1.default.nullable(zod_1.default.string()),
        clientSecret: zod_1.default.nullable(zod_1.default.string()),
        frontendDomain: zod_1.default.nullable(zod_1.default.string().refine((val) => {
            const fqdnRegex = /(?=^.{4,253}$)(^((?!-)[a-zA-Z0-9-]{0,62}[a-zA-Z0-9]\.)+[a-zA-Z]{2,63}$)/;
            if (fqdnRegex.test(val)) {
                return true;
            }
            return false;
        }, {
            message: 'Must be a valid fully qualified domain name (FQDN), e.g. "docs.ex"',
        })),
        keycloakDomain: zod_1.default.nullable(zod_1.default.string().refine((val) => {
            const fqdnRegex = /(?=^.{4,253}$)(^((?!-)[a-zA-Z0-9-]{0,62}[a-zA-Z0-9]\.)+[a-zA-Z]{2,63}$)/;
            if (fqdnRegex.test(val)) {
                return true;
            }
            return false;
        }, {
            message: 'Must be a valid fully qualified domain name (FQDN), e.g. "docs.ex"',
        })),
        allowedEmails: zod_1.default.array(zod_1.default.string()),
    })
        .default({
        enabled: false,
        realm: null,
        frontendDomain: null,
        keycloakDomain: null,
        clientId: null,
        clientSecret: null,
        allowedEmails: [],
    }),
    generic: zod_1.default
        .object({
        enabled: zod_1.default.boolean(),
        clientId: zod_1.default.nullable(zod_1.default.string()),
        clientSecret: zod_1.default.nullable(zod_1.default.string()),
        withPkce: zod_1.default.boolean(),
        authorizationUrl: zod_1.default.nullable(zod_1.default.string()),
        tokenUrl: zod_1.default.nullable(zod_1.default.string()),
        frontendDomain: zod_1.default.nullable(zod_1.default.string().refine((val) => {
            if (val.startsWith('127.0.0.1') || val.startsWith('localhost')) {
                return true;
            }
            const fqdnRegex = /(?=^.{4,253}$)(^((?!-)[a-zA-Z0-9-]{0,62}[a-zA-Z0-9]\.)+[a-zA-Z]{2,63}$)/;
            if (fqdnRegex.test(val)) {
                return true;
            }
            return false;
        }, {
            message: 'Must be a valid fully qualified domain name (FQDN), e.g. "docs.ex"',
        })),
        allowedEmails: zod_1.default.array(zod_1.default.string()),
    })
        .default({
        enabled: false,
        frontendDomain: null,
        tokenUrl: null,
        clientId: null,
        clientSecret: null,
        withPkce: false,
        authorizationUrl: null,
        allowedEmails: [],
    }),
    telegram: zod_1.default
        .object({
        enabled: zod_1.default.boolean(),
        clientId: zod_1.default.nullable(zod_1.default.string()),
        clientSecret: zod_1.default.nullable(zod_1.default.string()),
        allowedIds: zod_1.default.array(zod_1.default.string()),
        frontendDomain: zod_1.default.nullable(zod_1.default.string().refine((val) => {
            const fqdnRegex = /(?=^.{4,253}$)(^((?!-)[a-zA-Z0-9-]{0,62}[a-zA-Z0-9]\.)+[a-zA-Z]{2,63}$)/;
            if (fqdnRegex.test(val)) {
                return true;
            }
            return false;
        }, {
            message: 'Must be a valid fully qualified domain name (FQDN), e.g. "docs.ex"',
        })),
    })
        .default({
        enabled: false,
        clientId: null,
        clientSecret: null,
        allowedIds: [],
        frontendDomain: null,
    }),
});
