import {
    GetBandwidthStatsCommand,
    GetExodusHealthCommand,
    GetStatsCommand
} from '@exodus/backend-contract'

export interface IProps {
    bandwidthStats: GetBandwidthStatsCommand.Response['response']
    exodusHealth: GetExodusHealthCommand.Response['response']
    systemInfo: GetStatsCommand.Response['response']
}
