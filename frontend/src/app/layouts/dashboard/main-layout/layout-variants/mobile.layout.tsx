import { Group } from '@mantine/core'
import { useClickOutside, useDisclosure } from '@mantine/hooks'

import { HeaderControls } from '@shared/ui'

import { IExodusInfo } from '@entities/dashboard/updates-store'

import { DASHBOARD_LINKS } from '../layout-shared'
import { SidebarShellLayout } from './sidebar-shell.layout'

interface IProps {
    headerControls: React.ReactNode
    isLoadingUpdates: boolean
    isSocialButtons: boolean
    exodusInfo: IExodusInfo
}

export const MobileLayout = (props: IProps) => {
    const { headerControls, isLoadingUpdates, isSocialButtons, exodusInfo } = props

    const [opened, { toggle }] = useDisclosure()

    const ref = useClickOutside(() => {
        if (opened) {
            toggle()
        }
    })

    return (
        <SidebarShellLayout
            closedSide="mobile"
            footer={
                isSocialButtons && (
                    <Group justify="center" mt="md" style={{ flexShrink: 0 }}>
                        <HeaderControls
                            {...DASHBOARD_LINKS}
                            isGithubLoading={isLoadingUpdates}
                            stars={exodusInfo.starsCount || undefined}
                            withLanguage={false}
                            withLogout={false}
                            withVersion={false}
                        />
                    </Group>
                )
            }
            headerControls={headerControls}
            navbarRef={ref}
            onNavClose={toggle}
            opened={opened}
            padding="md"
            toggle={toggle}
        />
    )
}
