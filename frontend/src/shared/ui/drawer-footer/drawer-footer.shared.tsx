import { Drawer, Group } from '@mantine/core'

import { useIsMobile } from '@shared/hooks'

import styles from './DrawerFooter.module.css'

interface IProps {
    children: React.ReactNode
}

export function DrawerFooter(props: IProps) {
    const { children } = props
    const isMobile = useIsMobile()

    return (
        <Drawer.Header
            bottom={10}
            className={styles.modalFooter}
            component="footer"
            h="auto"
            mt="md"
            pos="sticky"
        >
            <Group
                gap="md"
                grow={!!isMobile}
                justify="flex-end"
                preventGrowOverflow={false}
                w="100%"
                wrap="wrap"
            >
                {children}
            </Group>
        </Drawer.Header>
    )
}
