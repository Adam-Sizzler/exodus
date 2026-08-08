/* eslint-disable camelcase */
import { MRT_ColumnSizingState, MRT_VisibilityState } from '@kastov/mantine-react-table-open'

import { createMrtTableStore } from '@shared/lib/mrt-table-store'

const DEFAULT_VISIBILITY: MRT_VisibilityState = {
    shortUuid: false,
    createdAt: false,
    subRevokedAt: false,
    description: false,
    telegramId: false,
    email: false,
    externalSquadUuid: false,
    vlessUuid: false,
    trojanPassword: false,
    ssPassword: false,
    naivePassword: false,
    shadowtlsPassword: false,
    hysteria2Password: false,
    anytlsPassword: false,
    hwidDeviceLimit: false,
    usedTrafficPercentage: false,
    trafficLimitBytes: false
}

const DEFAULT_COLUMN_SIZE: MRT_ColumnSizingState = {
    shortUuid: 100,
    createdAt: 100,
    subRevokedAt: 100,
    totalUsedBytes: 100,
    onlineAt: 100
}

export const useUsersTableStore = createMrtTableStore({
    name: 'x-rmnw-users-table',
    version: 13,
    defaults: {
        columnVisibility: DEFAULT_VISIBILITY,
        columnSize: DEFAULT_COLUMN_SIZE,
        sorting: [{ id: 'id', desc: true }]
    }
})

export const useUsersTableStoreActions = () => useUsersTableStore((store) => store.actions)
