import { UnstyledButton } from '@mantine/core'
import clsx from 'clsx'
import { MouseEvent } from 'react'
import { TbCheck, TbMinus } from 'react-icons/tb'

import classes from '../api-token-card.module.css'
import { KindState } from './scopes.utils'

interface IProps {
    label: string
    onClick?: () => void
    readOnly?: boolean
    state: KindState
}

export const KindPill = (props: IProps) => {
    const { label, onClick, readOnly, state } = props

    const handleClick = (event: MouseEvent) => {
        event.stopPropagation()
        onClick?.()
    }

    const isActive = state === 'on' || state === 'partial'

    return (
        <UnstyledButton
            className={clsx(classes.kindPill, {
                [classes.kindPillActive]: isActive
            })}
            data-readonly={readOnly || undefined}
            onClick={readOnly ? undefined : handleClick}
        >
            {state === 'on' && <TbCheck size={12} />}
            {state === 'partial' && <TbMinus size={12} />}
            {label}
        </UnstyledButton>
    )
}
