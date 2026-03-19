"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.CerberusInjectorSchema = void 0;
const zod_1 = __importDefault(require("zod"));
const HostSelectorSchema = zod_1.default.discriminatedUnion('type', [
    zod_1.default.object({
        type: zod_1.default.literal('uuids'),
        values: zod_1.default.array(zod_1.default.string().uuid()).min(1),
    }),
    zod_1.default.object({
        type: zod_1.default.literal('remarkRegex'),
        pattern: zod_1.default.string().min(1),
    }),
    zod_1.default.object({
        type: zod_1.default.literal('tagRegex'),
        pattern: zod_1.default.string().min(1),
    }),
    zod_1.default.object({
        type: zod_1.default.literal('sameTagAsRecipient'),
    }),
]);
const InjectHostsEntrySchema = zod_1.default.object({
    selector: HostSelectorSchema,
    selectFrom: zod_1.default.enum(['ALL', 'HIDDEN', 'NOT_HIDDEN']).optional(),
    tagPrefix: zod_1.default.string().min(1).optional(),
    useHostRemarkAsTag: zod_1.default.boolean().optional(),
    useHostTagAsTag: zod_1.default.boolean().optional(),
});
exports.CerberusInjectorSchema = zod_1.default.object({
    injectHosts: zod_1.default.array(InjectHostsEntrySchema).optional(),
    addVirtualHostAsOutbound: zod_1.default.boolean().optional(),
});
