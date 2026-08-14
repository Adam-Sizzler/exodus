import z from 'zod';

export const NodePluginSchema = z.object({
    uuid: z.uuid(),
    viewPosition: z.number().int(),
    name: z.string(),
    pluginConfig: z.nullable(z.unknown()),
});
