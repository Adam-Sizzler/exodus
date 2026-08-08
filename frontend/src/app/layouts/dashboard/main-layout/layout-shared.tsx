import { AppShell, Group, GroupProps } from '@mantine/core'
import { Outlet, ScrollRestoration } from 'react-router'

import { SidebarLogoShared, SidebarTitleShared } from '@shared/ui/sidebar'

import { app } from '../../../../config'

export const DASHBOARD_LINKS = {
    githubLink: app.githubRepo,
    telegramLink: 'https://t.me/exodus'
} as const

type LayoutMainProps = Omit<React.ComponentProps<typeof AppShell.Main>, 'children'>

export const LayoutMain = (props: LayoutMainProps) => (
    <AppShell.Main {...props}>
        <Outlet />
        <ScrollRestoration />
    </AppShell.Main>
)

export const LayoutBrand = (props: GroupProps) => (
    <Group {...props}>
        <SidebarLogoShared />
        <SidebarTitleShared />
    </Group>
)
