"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.PasskeySettingsSchema = void 0;
const zod_1 = __importDefault(require("zod"));
exports.PasskeySettingsSchema = zod_1.default.object({
    enabled: zod_1.default.boolean(),
    rpId: zod_1.default.nullable(zod_1.default.string().refine((val) => {
        if (val === 'localhost' || val === '127.0.0.1') {
            return true;
        }
        const fqdnRegex = /(?=^.{4,253}$)(^((?!-)[a-zA-Z0-9-]{0,62}[a-zA-Z0-9]\.)+[a-zA-Z]{2,63}$)/;
        if (fqdnRegex.test(val)) {
            return true;
        }
        return false;
    }, {
        message: 'Must be a valid fully qualified domain name (FQDN), e.g. "docs.exodus.dev"',
    })),
    origin: zod_1.default.nullable(zod_1.default.string().refine((value) => {
        try {
            const parsed = new URL(value);
            if (parsed.protocol !== 'http:' && parsed.protocol !== 'https:') {
                return false;
            }
            if (parsed.hostname === 'localhost' || parsed.hostname === '127.0.0.1') {
                return true;
            }
            if (parsed.protocol === 'https:' && /(?=.*\.[a-z]{2,})[^\s\/?#]+$/i.test(parsed.hostname)) {
                return true;
            }
        }
        catch {
            return false;
        }
        return false;
    }, {
        message: 'Must be a valid plain URL, e.g. "https://docs.exodus.dev".',
    })),
});
