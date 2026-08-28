import { AppShell, Box, Burger, Container, Group, ScrollArea } from '@mantine/core'
import clsx from 'clsx'

import { LayoutBrand, LayoutMain } from '../layout-shared'
import classes from '../layout.module.css'
import { MobileNavigation } from '../navbar/mobile-navigation.layout'

interface IProps {
    closedSide: 'desktop' | 'mobile'
    footer?: React.ReactNode
    headerControls: React.ReactNode
    navbarRef?: React.Ref<HTMLDivElement>
    onNavClose?: () => void
    opened: boolean
    padding: string
    toggle: () => void
    withFadeIn?: boolean
}

export const SidebarShellLayout = (props: IProps) => {
    const {
        closedSide,
        footer,
        headerControls,
        navbarRef,
        onNavClose,
        opened,
        padding,
        toggle,
        withFadeIn
    } = props

    const navbarCollapsed =
        closedSide === 'mobile'
            ? { mobile: !opened, desktop: true }
            : { mobile: false, desktop: !opened }

    const closedClass =
        closedSide === 'mobile'
            ? classes.sidebarWrapperClosedMobile
            : classes.sidebarWrapperClosedDesktop

    return (
        <AppShell
            className={withFadeIn ? classes.appShellFadeIn : undefined}
            header={{ height: 64, collapsed: false, offset: false }}
            layout="alt"
            navbar={{ width: 300, breakpoint: 'lg', collapsed: navbarCollapsed }}
            padding={padding}
            transitionDuration={500}
            transitionTimingFunction="ease-in-out"
        >
            <AppShell.Header className={classes.header} withBorder={false}>
                <Container fluid px="lg" py="xs">
                    <Group justify="space-between" style={{ flexWrap: 'nowrap' }}>
                        <Group style={{ flex: 1, justifyContent: 'flex-start' }}>
                            <Burger onClick={toggle} opened={opened} size="md" />
                        </Group>
                        <Group style={{ flexShrink: 0 }}>{headerControls}</Group>
                    </Group>
                </Container>
            </AppShell.Header>

            <AppShell.Navbar
                className={clsx(classes.sidebarWrapper, { [closedClass]: !opened })}
                p="md"
                pb={0}
                ref={navbarRef}
                w={300}
                withBorder={false}
            >
                <AppShell.Section className={classes.logoSection}>
                    <Box style={{ position: 'absolute', left: '0' }}>
                        <Burger hiddenFrom="lg" onClick={toggle} opened={opened} size="sm" />
                    </Box>

                    <LayoutBrand gap="xs" justify="center" wrap="nowrap" />
                </AppShell.Section>

                <AppShell.Section
                    className={classes.scrollArea}
                    component={ScrollArea}
                    flex={1}
                    scrollbarSize="0.2rem"
                >
                    <MobileNavigation onClose={onNavClose} />
                </AppShell.Section>

                {footer && (
                    <AppShell.Section className={classes.footerSection}>{footer}</AppShell.Section>
                )}
            </AppShell.Navbar>

            <LayoutMain
                pb="var(--mantine-spacing-md)"
                pt="calc(var(--app-shell-header-height) + 10px)"
            />
        </AppShell>
    )
}
