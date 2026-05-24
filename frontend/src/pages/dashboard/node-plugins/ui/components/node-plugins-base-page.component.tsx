import { TbPackage } from 'react-icons/tb'
import { motion } from 'motion/react'

import { NodePluginsGridWidget } from '@widgets/dashboard/node-plugins/node-plugins-grid/node-plugins-grid.widget'
import { NodePluginsHeaderActionButtonsFeature } from '@features/ui/dashboard/node-plugins/header-action-buttons'
import { NodePluginsSpotlightWidget } from '@widgets/dashboard/node-plugins/node-plugins-spotlight'
import { NodePluginExecutorDrawer } from '@widgets/dashboard/node-plugins/node-plugin-executor'
import { RenameModalShared } from '@shared/ui/modals/rename-modal.shared'
import { NodePluginResponse, NodeResponse } from '@shared/api/hooks'
import { Page, PageHeaderShared } from '@shared/ui'

interface IProps {
    nodes: NodeResponse[]
    plugins: NodePluginResponse[]
}

export function NodePluginsBasePageComponent(props: IProps) {
    const { nodes, plugins } = props

    return (
        <Page title="Plugins">
            <PageHeaderShared
                actions={<NodePluginsHeaderActionButtonsFeature />}
                icon={<TbPackage size={24} />}
                title="Plugins β"
                wrapActions
            />

            <motion.div
                animate={{ opacity: 1 }}
                initial={{ opacity: 0 }}
                transition={{ duration: 0.5 }}
            >
                <NodePluginsGridWidget nodes={nodes} plugins={plugins} />
            </motion.div>

            <NodePluginsSpotlightWidget plugins={plugins} />

            <RenameModalShared key="rename-node-plugin-modal" renameFrom="nodePlugin" />
            <NodePluginExecutorDrawer />
        </Page>
    )
}
