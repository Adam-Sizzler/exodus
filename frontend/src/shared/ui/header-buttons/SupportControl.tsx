import { TbHeartFilled } from 'react-icons/tb'

import { HeaderControl } from './HeaderControl'
import classes from './SupportControl.module.css'

export function SupportControl() {
    return (
        <HeaderControl
            className={classes.support}
            component="a"
            href="https://docs.ex/docs/donate"
            rel="noopener noreferrer"
            target="_blank"
        >
            <TbHeartFilled />
        </HeaderControl>
    )
}
