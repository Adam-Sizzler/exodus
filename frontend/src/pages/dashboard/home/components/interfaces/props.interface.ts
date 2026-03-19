import {
    GetBandwidthStatsCommand,
    GetCerberusHealthCommand,
    GetStatsCommand
} from '@cerberus/backend-contract'

export interface IProps {
    bandwidthStats: GetBandwidthStatsCommand.Response['response']
    cerberusHealth: GetCerberusHealthCommand.Response['response']
    systemInfo: GetStatsCommand.Response['response']
}
