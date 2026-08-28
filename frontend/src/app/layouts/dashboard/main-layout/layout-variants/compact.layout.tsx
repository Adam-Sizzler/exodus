import { AppShell, Group, Divider } from '@mantine/core'

import { LayoutBrand, LayoutMain } from '../layout-shared'
import classes from '../layout.module.css'
import { DesktopNavigation } from '../navbar/desktop-navigation.layout'

interface IProps {
    headerControls: React.ReactNode
    isHiResDesktop: boolean
}

export const CompactLayout = (props: IProps) => {
    const { headerControls, isHiResDesktop } = props

    return (
        <AppShell
            className={classes.appShellFadeIn}
            header={{ height: isHiResDesktop ? 64 : 116, offset: false }}
            padding="xl"
        >
            <AppShell.Header className={classes.header}>
                <div className={classes.brandRow}>
                    <Group align="stretch" gap="xs" h="100%" style={{ minWidth: 0 }} wrap="nowrap">
                        <LayoutBrand gap="xs" style={{ flexShrink: 0 }} wrap="nowrap" />

                        {isHiResDesktop && (
                            <>
                                <Divider
                                    h="50%"
                                    orientation="vertical"
                                    style={{
                                        alignSelf: 'center',
                                        marginLeft: '10px',
                                        marginRight: '10px'
                                    }}
                                />

                                <DesktopNavigation />
                            </>
                        )}
                    </Group>

                    <Group gap="xs" style={{ flexShrink: 0 }} wrap="nowrap">
                        {headerControls}
                    </Group>
                </div>
                {!isHiResDesktop && (
                    <div className={classes.navRowDesktop}>
                        <DesktopNavigation />
                    </div>
                )}
            </AppShell.Header>

            <LayoutMain
                pb="calc(var(--mantine-spacing-md) + env(safe-area-inset-bottom, 0px))"
                pl={isHiResDesktop ? '10vw' : 'max(var(--mantine-spacing-xl), env(safe-area-inset-left, 0px))'}
                pr={isHiResDesktop ? '10vw' : 'max(var(--mantine-spacing-xl), env(safe-area-inset-right, 0px))'}
                pt="calc(var(--app-shell-header-height) + var(--mantine-spacing-md))"
            />
        </AppShell>
    )
}
