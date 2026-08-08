import { Box } from '@mantine/core'
import clsx from 'clsx'
import { useState } from 'react'
import { Virtuoso } from 'react-virtuoso'

import { NodeActiveSessionItem } from './node-active-session.item.widget'
import classes from './node-active-sessions.module.css'
import { ActiveSessionUser } from './use-node-active-sessions'

interface IProps {
    users: ActiveSessionUser[]
}

export const SessionsListWidget = ({ users }: IProps) => {
    const [isScrolled, setIsScrolled] = useState(false)

    return (
        <Box
            className={clsx(
                classes.listContainer,
                classes.fadeBottom,
                isScrolled && classes.fadeTop
            )}
        >
            <Virtuoso
                data={users}
                itemContent={(_index, user) => (
                    <Box className={classes.itemWrapper}>
                        <NodeActiveSessionItem user={user} />
                    </Box>
                )}
                onScroll={(e) => {
                    const target = e.target as HTMLElement
                    setIsScrolled(target.scrollTop > 0)
                }}
                style={{
                    height: '100%'
                }}
                totalCount={users.length}
                useWindowScroll={false}
            />
        </Box>
    )
}
