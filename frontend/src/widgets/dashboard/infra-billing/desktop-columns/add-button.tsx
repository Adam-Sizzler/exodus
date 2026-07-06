import { ActionIcon } from '@mantine/core'
import { TbPlus } from 'react-icons/tb'

interface AddButtonProps {
    onClick: () => void
}

export function AddButton(props: AddButtonProps) {
    const { onClick } = props

    return (
        <ActionIcon color="teal" onClick={onClick} size="input-xs" variant="soft">
            <TbPlus size={18} />
        </ActionIcon>
    )
}
