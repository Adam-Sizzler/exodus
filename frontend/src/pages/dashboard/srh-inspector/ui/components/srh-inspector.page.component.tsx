import { Stack } from '@mantine/core'
import { useMediaQuery } from '@mantine/hooks'
import { SrhInspectorMetrics } from '@widgets/dashboard/srh-inspector/srh-inspector-metrics'
import { SrhInspectorTableWidget } from '@widgets/dashboard/srh-inspector/srh-inspector-table'
import { motion } from 'framer-motion'
import { useTranslation } from 'react-i18next'

import { Page } from '@shared/ui'
import { MobileWarningOverlay } from '@shared/ui/mobile-warning-overlay/mobile-warning-overlay'

export default function SrhInspectorPageComponent() {
    const { t } = useTranslation()
    const isMobile = useMediaQuery('(max-width: 500px)')

    return (
        <Page title={t('constants.srh-inspector')}>
            <Stack>
                {isMobile && <MobileWarningOverlay />}
                <SrhInspectorMetrics />

                <motion.div
                    animate={{ opacity: 1 }}
                    initial={{ opacity: 0 }}
                    transition={{ duration: 0.5, ease: [0, 0.71, 0.2, 1.01] }}
                >
                    <SrhInspectorTableWidget />
                </motion.div>
            </Stack>
        </Page>
    )
}
