import { NodePluginsHeaderActionButtonsFeature } from '@features/ui/dashboard/node-plugins/header-action-buttons'
import { GetNodesCommand, GetNodePluginsCommand } from '@exodus/backend-contract'
import { NodePluginsGridWidget } from '@widgets/dashboard/node-plugins/node-plugins-grid/node-plugins-grid.widget'
import { NodePluginsSpotlightWidget } from '@widgets/dashboard/node-plugins/node-plugins-spotlight'
import { motion } from 'motion/react'
import { useTranslation } from 'react-i18next'
import { TbPackage } from 'react-icons/tb'

import { Page, PageHeaderShared } from '@shared/ui'

interface Props {
    nodes: GetNodesCommand.Response['response']
    plugins: GetNodePluginsCommand.Response['response']['nodePlugins']
}

export const NodePluginsBasePageComponent = (props: Props) => {
    const { nodes, plugins } = props
    const { t } = useTranslation()

    return (
        <Page title={t('constants.node-plugins')}>
            <PageHeaderShared
                actions={<NodePluginsHeaderActionButtonsFeature />}
                icon={<TbPackage size={24} />}
                title={`${t('constants.node-plugins')} β`}
            />

            <motion.div
                animate={{ opacity: 1 }}
                initial={{ opacity: 0 }}
                transition={{ duration: 0.5 }}
            >
                <NodePluginsGridWidget nodes={nodes} plugins={plugins} />
            </motion.div>

            <NodePluginsSpotlightWidget plugins={plugins} />
        </Page>
    )
}
