import { Drawer, DrawerProps, Group } from '@mantine/core'
import { ReactNode } from 'react'

interface IProps {
    drawerProps: Omit<DrawerProps, 'title' | 'children'>
    children: ReactNode
    title: ReactNode
    buttons?: ReactNode
}

export const CompoundDrawerShared = (props: IProps) => {
    const { drawerProps, children, title, buttons } = props

    return (
        <Drawer.Root {...drawerProps}>
            <Drawer.Overlay />
            <Drawer.Content>
                <Drawer.Header>
                    <Drawer.Title>{title}</Drawer.Title>
                    <Group gap="xs" justify="flex-end" wrap="nowrap">
                        {buttons}
                        <Drawer.CloseButton />
                    </Group>
                </Drawer.Header>
                <Drawer.Body>{children}</Drawer.Body>
            </Drawer.Content>
        </Drawer.Root>
    )
}
