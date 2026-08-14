import { SubscriptionSettingsSchema } from '../subscription-settings.schema';

export const ExternalSquadSubscriptionSettingsSchema = SubscriptionSettingsSchema.pick({
    serveJsonAtBaseSubscription: true,
    isShowCustomRemarks: true,
    randomizeHosts: true,
}).partial();
