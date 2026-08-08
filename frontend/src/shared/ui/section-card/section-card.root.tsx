import { Card, CardProps, Divider, MantineSpacing, Stack } from '@mantine/core'
import { Children, ReactNode, RefObject } from 'react'

interface ISectionCardRootProps extends Omit<CardProps, 'children'> {
    children: ReactNode
    dividerOpacity?: number
    onlyFirstDivider?: boolean
    allDividers?: boolean
    gap?: MantineSpacing
    ref?: RefObject<HTMLDivElement | null>
}

export function SectionCardRoot({
    children,
    dividerOpacity = 0.3,
    onlyFirstDivider = false,
    allDividers = true,
    gap = 'md',
    p = 'md',
    radius = 'md',
    style,
    ref,
    ...props
}: ISectionCardRootProps) {
    const childArray = Children.toArray(children).filter(Boolean)

    const childrenWithDividers = childArray.flatMap((child, index) => {
        if (onlyFirstDivider && index === 0 && childArray.length > 1) {
            return [child, <Divider key={`divider-${index}`} style={{ opacity: dividerOpacity }} />]
        }
        if (allDividers && index < childArray.length - 1) {
            return [child, <Divider key={`divider-${index}`} style={{ opacity: dividerOpacity }} />]
        }
        return [child]
    })

    return (
        <Card
            p={p}
            radius={radius}
            ref={ref}
            style={{
                background: 'rgba(255, 255, 255, 0.02)',
                border: '1px solid rgba(255, 255, 255, 0.08)',
                ...style
            }}
            {...props}
        >
            <Stack gap={gap}>{childrenWithDividers}</Stack>
        </Card>
    )
}
