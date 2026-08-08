import { Stack } from '@mantine/core'
import { useMediaQuery } from '@mantine/hooks'
import { HwidInspectorLeaderboardWidget } from '@widgets/dashboard/hwid-inspector/hwid-inspector-leaderboard'
import { HwidInspectorMetrics } from '@widgets/dashboard/hwid-inspector/hwid-inspector-metrics'
import { HwidInspectorTableWidget } from '@widgets/dashboard/hwid-inspector/hwid-inspector-table'
import { motion } from 'framer-motion'
import { useTranslation } from 'react-i18next'

import { Page } from '@shared/ui'
import { MobileWarningOverlay } from '@shared/ui/mobile-warning-overlay/mobile-warning-overlay'

export default function HwidInspectorPageComponent() {
    const { t } = useTranslation()
    const isMobile = useMediaQuery('(max-width: 500px)')

    return (
        <Page title={t('constants.hwid-inspector')}>
            <Stack>
                {isMobile && <MobileWarningOverlay />}
                <HwidInspectorMetrics />

                <motion.div
                    animate={{ opacity: 1 }}
                    initial={{ opacity: 0 }}
                    transition={{ duration: 0.5, ease: [0, 0.71, 0.2, 1.01] }}
                >
                    <HwidInspectorLeaderboardWidget />
                </motion.div>

                <motion.div
                    animate={{ opacity: 1 }}
                    initial={{ opacity: 0 }}
                    transition={{ duration: 0.5, ease: [0, 0.71, 0.2, 1.01] }}
                >
                    <HwidInspectorTableWidget />
                </motion.div>
            </Stack>
        </Page>
    )
}
