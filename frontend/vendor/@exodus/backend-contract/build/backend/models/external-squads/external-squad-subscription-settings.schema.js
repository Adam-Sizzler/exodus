"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.ExternalSquadSubscriptionSettingsSchema = void 0;
const subscription_settings_schema_1 = require("../subscription-settings.schema");
exports.ExternalSquadSubscriptionSettingsSchema = subscription_settings_schema_1.SubscriptionSettingsSchema.pick({
    serveJsonAtBaseSubscription: true,
    isShowCustomRemarks: true,
    randomizeHosts: true,
}).partial();
