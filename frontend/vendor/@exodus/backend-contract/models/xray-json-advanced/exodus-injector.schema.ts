import { z } from 'zod';

const HostSelectorSchema = z.discriminatedUnion('type', [
    z.object({
        type: z.literal('uuids'),
        values: z.array(z.uuid()).min(1),
    }),
    z.object({
        type: z.literal('remarkRegex'),
        pattern: z.string().min(1),
    }),
    z.object({
        type: z.literal('tagRegex'),
        pattern: z.string().min(1),
    }),
    z.object({
        type: z.literal('sameTagAsRecipient'),
    }),
]);

const InjectHostsEntrySchema = z.object({
    selector: HostSelectorSchema,
    selectFrom: z.enum(['ALL', 'HIDDEN', 'NOT_HIDDEN']).optional(),
    tagPrefix: z.string().min(1).optional(),
    useHostRemarkAsTag: z.boolean().optional(),
    useHostTagAsTag: z.boolean().optional(),
});

export const ExodusInjectorSchema = z.object({
    injectHosts: z.array(InjectHostsEntrySchema).optional(),
    addVirtualHostAsOutbound: z.boolean().optional(),
});

export type TExodusInjector = z.infer<typeof ExodusInjectorSchema>;
export type TExodusInjectorSelector = z.infer<typeof HostSelectorSchema>;
export type TExodusInjectorSelectFrom = z.infer<typeof InjectHostsEntrySchema>['selectFrom'];
