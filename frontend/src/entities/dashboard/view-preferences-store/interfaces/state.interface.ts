import { CONFIG_PROFILES_VIEW_MODE, HOSTS_VIEW_MODE, LAYOUT_STYLE, NODES_VIEW_MODE } from './enums'

export interface IState {
    configProfilesViewMode: CONFIG_PROFILES_VIEW_MODE
    hostsActiveTag: null | string
    hostsViewMode: HOSTS_VIEW_MODE
    nodesActiveTag: null | string
    nodesViewMode: NODES_VIEW_MODE
    layoutStyle: LAYOUT_STYLE
}
