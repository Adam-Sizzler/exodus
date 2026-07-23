import {
    APP_CONFIG_ROUTE_LEADING_PATH,
    SubscriptionPageRawConfigSchema
} from '@exodus/subscription-page-types'
import { GetSubscriptionInfoByShortUuidCommand } from '@exodus/backend-contract'
import { Outlet } from 'react-router'
import { useLayoutEffect } from 'react'
import consola from 'consola/browser'
import { ofetch } from 'ofetch'

import {
    useSubscriptionInfoStoreActions,
    useSubscriptionInfoStoreInfo
} from '@entities/subscription-info-store'
import { useAppConfigStoreActions, useIsConfigLoaded } from '@entities/app-config-store'
import { LoadingScreen } from '@shared/ui'

import classes from './root.module.css'

export function RootLayout() {
    const subscriptionActions = useSubscriptionInfoStoreActions()
    const configActions = useAppConfigStoreActions()

    const { subscription } = useSubscriptionInfoStoreInfo()
    const isConfigLoaded = useIsConfigLoaded()

    useLayoutEffect(() => {
        const subPageDiv = document.getElementById('sbpg')

        if (subPageDiv) {
            const subscriptionUrl = subPageDiv.dataset.panel

            if (subscriptionUrl) {
                try {
                    const subscription: GetSubscriptionInfoByShortUuidCommand.Response = JSON.parse(
                        atob(subscriptionUrl)
                    )

                    subscriptionActions.setSubscriptionInfo({
                        subscription: subscription.response
                    })
                } catch (error) {
                    consola.log(error)
                } finally {
                    subPageDiv.remove()
                }
            }
        }

        const fetchConfig = async () => {
            try {
                // APP_CONFIG_ROUTE_LEADING_PATH ('/assets/.app-config-v2.json') is an
                // absolute, domain-root path baked into the vendor package. Passed
                // directly to fetch, it always resolves against the origin and
                // ignores whatever sub-path this app is actually deployed under
                // (e.g. /subscription/). Resolve it instead relative to this
                // module's own resolved URL, which already lives under the correct
                // assets/ folder for the current deployment.
                const configFileName = APP_CONFIG_ROUTE_LEADING_PATH.replace(/^\/assets\//, '')
                const configUrl = new URL(configFileName, import.meta.url)
                configUrl.searchParams.set('v', String(Date.now()))

                const tempConfig = await ofetch<unknown>(configUrl.toString(), {
                    parseResponse: (response) => JSON.parse(response)
                })

                const parsedConfig =
                    await SubscriptionPageRawConfigSchema.safeParseAsync(tempConfig)

                if (!parsedConfig.success) {
                    consola.error('Failed to parse app config:', parsedConfig.error)
                    return
                }

                configActions.setConfig(parsedConfig.data)
            } catch (error) {
                consola.error('Failed to fetch app config:', error)
            }
        }

        fetchConfig()
    }, [])

    if (!isConfigLoaded || !subscription) {
        return (
            <div className={classes.root}>
                <div className="animated-background"></div>
                <div className={classes.content}>
                    <main className={classes.main}>
                        <LoadingScreen height="100vh" />
                    </main>
                </div>
            </div>
        )
    }

    return (
        <div className={classes.root}>
            <div className="animated-background"></div>
            <div className={classes.content}>
                <main className={classes.main}>
                    <Outlet />
                </main>
            </div>
        </div>
    )
}
