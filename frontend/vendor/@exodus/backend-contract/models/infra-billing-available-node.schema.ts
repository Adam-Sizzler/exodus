import { NodesSchema } from './nodes.schema';

export const InfraBillingAvailableNodeSchema = NodesSchema.pick({
    uuid: true,
    name: true,
    countryCode: true,
});
