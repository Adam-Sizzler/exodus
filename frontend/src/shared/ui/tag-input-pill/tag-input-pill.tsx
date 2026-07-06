import { Badge, CloseButton } from '@mantine/core'

interface IProps {
    onRemove?: () => void
    value: string | undefined
}

export function TagInputPill(props: IProps) {
    const { onRemove, value } = props

    return (
        <Badge
            leftSection={
                onRemove ? (
                    <CloseButton
                        c="cyan"
                        onMouseDown={(event) => {
                            event.preventDefault()
                            event.stopPropagation()
                            onRemove()
                        }}
                        radius="xl"
                        size={18}
                        variant="transparent"
                    />
                ) : null
            }
            pl={0}
            radius="sm"
            size="md"
            styles={{
                label: { overflow: 'visible' },
                root: { textTransform: 'none' },
                section: { margin: 0 }
            }}
            variant="soft"
        >
            {value}
        </Badge>
    )
}
