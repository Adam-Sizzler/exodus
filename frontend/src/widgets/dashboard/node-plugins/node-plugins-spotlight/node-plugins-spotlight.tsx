import { GetNodePluginsCommand } from '@exodus/backend-contract'
import { TbPackage } from 'react-icons/tb'
import { generatePath, useNavigate } from 'react-router'

import { ROUTES } from '@shared/constants'
import { UniversalSpotlightContentShared } from '@shared/ui/universal-spotlight'

interface IProps {
    plugins: GetNodePluginsCommand.Response['response']['nodePlugins']
}

export const NodePluginsSpotlightWidget = (props: IProps) => {
    const { plugins } = props

    const navigate = useNavigate()

    const handleViewNodePlugin = (nodePluginUuid: string) => {
        navigate(
            generatePath(ROUTES.DASHBOARD.MANAGEMENT.NODE_PLUGINS.NODE_PLUGIN_BY_UUID, {
                uuid: nodePluginUuid
            })
        )
    }

    return (
        <UniversalSpotlightContentShared
            actions={plugins.map((plugin) => ({
                label: plugin.name,
                id: plugin.uuid,
                leftSection: <TbPackage size={24} />,

                onClick: () => handleViewNodePlugin(plugin.uuid)
            }))}
        />
    )
}
