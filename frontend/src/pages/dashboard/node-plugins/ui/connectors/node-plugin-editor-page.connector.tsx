import { Navigate, useParams } from 'react-router-dom'
import { useTranslation } from 'react-i18next'

import { useGetNodePlugin } from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui'
import { ROUTES } from '@shared/constants'

import { NodePluginEditorPageComponent } from '../components/node-plugin-editor-page.component'

export function NodePluginEditorPageConnector() {
    const { t } = useTranslation()
    const { uuid } = useParams<{ uuid: string }>()

    const { data: plugin, isLoading } = useGetNodePlugin({
        route: {
            uuid: uuid ?? ''
        },
        rQueryParams: {
            enabled: Boolean(uuid)
        }
    })

    if (!uuid) {
        return <Navigate replace to={ROUTES.DASHBOARD.MANAGEMENT.NODE_PLUGINS.ROOT} />
    }

    if (isLoading || !plugin) {
        return <LoadingScreen text={t('node-plugins-page.connector.loading-node-plugin')} />
    }

    return <NodePluginEditorPageComponent plugin={plugin} />
}
