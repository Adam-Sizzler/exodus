import { useDisclosure } from '@mantine/hooks'

import { SidebarShellLayout } from './sidebar-shell.layout'

interface IProps {
    headerControls: React.ReactNode
}

export const SidebarLayout = (props: IProps) => {
    const { headerControls } = props

    const [opened, { toggle }] = useDisclosure(true)

    return (
        <SidebarShellLayout
            closedSide="desktop"
            headerControls={headerControls}
            opened={opened}
            padding="xl"
            toggle={toggle}
            withFadeIn
        />
    )
}
