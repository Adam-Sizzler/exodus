import { Box, BoxProps } from '@mantine/core'
import { nprogress } from '@mantine/nprogress'
import { AnimatePresence, motion } from 'framer-motion'
import { forwardRef, ReactNode, useEffect, useMemo } from 'react'
import { app } from 'src/config'

import { useGetAuthStatus } from '@shared/api/hooks/auth/auth.query.hooks'
import { parseColoredTextUtil } from '@shared/utils/misc'

interface PageProps extends BoxProps {
    children: ReactNode
    meta?: ReactNode
    title: string
}

export const Page = forwardRef<HTMLDivElement, PageProps>(
    ({ children, title = '', meta, ...other }, ref) => {
        const { data: authStatus } = useGetAuthStatus()

        useEffect(() => {
            nprogress.complete()
            return () => nprogress.start()
        }, [])

        const titleParts = useMemo(() => {
            if (authStatus?.branding.title) {
                return parseColoredTextUtil(authStatus.branding.title)
                    .map((part) => part.text)
                    .join('')
            }

            return app.name
        }, [authStatus])

        const pageTitle = `${title} | ${titleParts}`

        return (
            <>
                <title>{pageTitle}</title>
                {meta}

                <AnimatePresence mode="wait">
                    <motion.div
                        animate={{ opacity: 1 }}
                        exit={{ opacity: 0 }}
                        initial={{ opacity: 0 }}
                        transition={{
                            duration: 0.3,
                            ease: 'easeInOut'
                        }}
                    >
                        <Box ref={ref} {...other}>
                            {children}
                        </Box>
                    </motion.div>
                </AnimatePresence>
            </>
        )
    }
)
