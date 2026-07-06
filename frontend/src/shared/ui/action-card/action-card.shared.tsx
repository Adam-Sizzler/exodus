import { Group, Paper, Stack, Text, ThemeIcon, ThemeIconProps, UnstyledButton } from '@mantine/core'
import { useClickOutside } from '@mantine/hooks'
import clsx from 'clsx'
import { useState } from 'react'
import { TbArrowRight, TbChevronRight } from 'react-icons/tb'

import styles from './action-card.module.css'

interface IProps {
    description: string
    icon: React.ReactNode
    iconColor?: ThemeIconProps['color']
    isLoading?: boolean
    onClick: () => void
    title: string
    variant: ThemeIconProps['variant']
    withConfirmation?: boolean
}

export function ActionCardShared(props: IProps) {
    const { description, icon, iconColor, isLoading, onClick, title, variant, withConfirmation } =
        props
    const [isConfirming, setIsConfirming] = useState(false)
    const ref = useClickOutside(() => setIsConfirming(false))

    const handleClick = () => {
        if (!withConfirmation) {
            onClick()
            return
        }
        if (isConfirming) {
            setIsConfirming(false)
            onClick()
        } else {
            setIsConfirming(true)
        }
    }

    return (
        <UnstyledButton disabled={isLoading} onClick={handleClick} ref={ref} w="100%">
            <Paper
                className={clsx(
                    styles.card,
                    isLoading && styles.loading,
                    isConfirming && styles.confirming
                )}
                p="md"
                radius="md"
            >
                <Group gap="md" justify="space-between" wrap="nowrap">
                    <Group gap="md" wrap="nowrap">
                        <ThemeIcon color={iconColor} radius="md" size="xl" variant={variant}>
                            {icon}
                        </ThemeIcon>
                        <Stack gap={2}>
                            <Text fw={600} size="sm">
                                {title}
                            </Text>
                            <Text c="dimmed" size="xs">
                                {description}
                            </Text>
                        </Stack>
                    </Group>
                    {isConfirming ? (
                        <TbArrowRight color="var(--mantine-color-red-5)" size={20} />
                    ) : (
                        <TbChevronRight color="var(--mantine-color-dimmed)" size={20} />
                    )}
                </Group>
            </Paper>
        </UnstyledButton>
    )
}
