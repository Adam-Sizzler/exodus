import { inferQueryKeyStore, mergeQueryKeys } from '@lukemorales/query-key-factory'

import { subscriptionRequestHistoryQueryKeys } from './subscription-request-history/subscription-request-history.query.hooks'
import { subscriptionTemplateQueryKeys } from './subscription-template/subscription-template.query.hooks'
import { subscriptionSettingsQueryKeys } from './subscription-settings/subscription-settings.query.hooks'
import { exodusSettingsQueryKeys } from './exodus-settings/exodus-settings.query.hooks'
import { hwidUserDevicesQueryKeys } from './hwid-user-devices/hwid-user-devices.query.hooks'
import { configProfilesQueryKeys } from './config-profiles/config-profiles.query.hooks'
import { internalSquadsQueryKeys } from './internal-squads/internal-squads.query.hooks'
import { externalSquadsQueryKeys } from './external-squads/external-squads.query.hooks'
import { subpageConfigsQueryKeys } from './subpage-configs/subpage-configs.query.hooks'
import { bandwidthStatsQueryKeys } from './bandwidth-stats/bandwidth-stats.query.hooks'
import { infraBillingQueryKeys } from './infra-billing/infra-billing.query.hooks'
import { apiTokensQueryKeys } from './api-tokens/api-tokens.query.hooks'
import { ipControlQueryKeys } from './ip-control/ip-control.query.hooks'
import { passkeysQueryKeys } from './passkeys/passkeys.query.hooks'
import { srsListsQueryKeys } from './srs-lists/srs-lists.query.hooks'
import { systemQueryKeys } from './system/system.query.hooks'
import { usersQueryKeys } from './users/users.query.hooks'
import { hostsQueryKeys } from './hosts/hosts.query.hooks'
import { nodesQueryKeys } from './nodes/nodes.query.hooks'
import { subscriptionConnectionsQueryKeys } from './subscription-connections/subscription-connections.query.hooks'
import { authQueryKeys } from './auth/auth.query.hooks'

export const QueryKeys = mergeQueryKeys(
    usersQueryKeys,
    systemQueryKeys,
    hostsQueryKeys,
    nodesQueryKeys,
    subscriptionConnectionsQueryKeys,
    apiTokensQueryKeys,
    authQueryKeys,
    subscriptionTemplateQueryKeys,
    subscriptionSettingsQueryKeys,
    hwidUserDevicesQueryKeys,
    configProfilesQueryKeys,
    internalSquadsQueryKeys,
    infraBillingQueryKeys,
    subscriptionRequestHistoryQueryKeys,
    externalSquadsQueryKeys,
    exodusSettingsQueryKeys,
    passkeysQueryKeys,
    srsListsQueryKeys,
    subpageConfigsQueryKeys,
    bandwidthStatsQueryKeys,
    ipControlQueryKeys
)

export type TQueryKeys = inferQueryKeyStore<typeof QueryKeys>
