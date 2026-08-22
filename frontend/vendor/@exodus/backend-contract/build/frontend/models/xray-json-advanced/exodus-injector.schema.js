"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExodusInjectorSchema = void 0;
const zod_1 = require("zod");
const HostSelectorSchema = zod_1.z.discriminatedUnion('type', [
    zod_1.z.object({
        type: zod_1.z.literal('uuids'),
        values: zod_1.z.array(zod_1.z.uuid()).min(1),
    }),
    zod_1.z.object({
        type: zod_1.z.literal('remarkRegex'),
        pattern: zod_1.z.string().min(1),
    }),
    zod_1.z.object({
        type: zod_1.z.literal('tagRegex'),
        pattern: zod_1.z.string().min(1),
    }),
    zod_1.z.object({
        type: zod_1.z.literal('sameTagAsRecipient'),
    }),
]);
const InjectHostsEntrySchema = zod_1.z.object({
    selector: HostSelectorSchema,
    selectFrom: zod_1.z.enum(['ALL', 'HIDDEN', 'NOT_HIDDEN']).optional(),
    tagPrefix: zod_1.z.string().min(1).optional(),
    useHostRemarkAsTag: zod_1.z.boolean().optional(),
    useHostTagAsTag: zod_1.z.boolean().optional(),
});
exports.ExodusInjectorSchema = zod_1.z.object({
    injectHosts: zod_1.z.array(InjectHostsEntrySchema).optional(),
    addVirtualHostAsOutbound: zod_1.z.boolean().optional(),
});
