import {
    GetNodeIntegrationsCommand,
    GetNodePluginsCommand,
    GetNodesCommand
} from '@exodus/backend-contract'

export interface IProps {
    nodes: GetNodesCommand.Response['response'] | undefined
    nodePlugins: GetNodePluginsCommand.Response['response'] | undefined
    nodeIntegrations: GetNodeIntegrationsCommand.Response['response'] | undefined
}
