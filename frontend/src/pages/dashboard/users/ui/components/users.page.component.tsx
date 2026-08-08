import { Stack } from '@mantine/core'
import { useMediaQuery } from '@mantine/hooks'
import { UsersMetrics } from '@widgets/dashboard/users/users-metrics'
import { UserTableWidget } from '@widgets/dashboard/users/users-table'
import { motion } from 'framer-motion'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { LoadingScreen } from '@shared/ui'
import { MobileWarningOverlay } from '@shared/ui/mobile-warning-overlay/mobile-warning-overlay'
import { Page } from '@shared/ui/page'

const DeferredUserTableWidget = () => {
    const [shouldRender, setShouldRender] = useState(false)

    useEffect(() => {
        const timer = setTimeout(() => setShouldRender(true), 300)
        return () => clearTimeout(timer)
    }, [])

    if (!shouldRender) {
        return <LoadingScreen height="80vh" />
    }

    return (
        <motion.div
            animate={{ opacity: 1 }}
            initial={{ opacity: 0 }}
            transition={{ duration: 0.5, ease: [0, 0.71, 0.2, 1.01] }}
        >
            <UserTableWidget />
        </motion.div>
    )
}

export default function UsersPageComponent() {
    const { t } = useTranslation()
    const isMobile = useMediaQuery('(max-width: 500px)')

    return (
        <Page title={t('constants.users')}>
            <Stack>
                {isMobile && <MobileWarningOverlay />}
                <UsersMetrics />

                <DeferredUserTableWidget />
            </Stack>
        </Page>
    )
}
