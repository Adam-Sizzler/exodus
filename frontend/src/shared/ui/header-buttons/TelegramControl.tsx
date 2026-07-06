import { PiTelegramLogoBold } from 'react-icons/pi'

import { HeaderControl } from './HeaderControl'
import classes from './TelegramControl.module.css'

interface TelegramControlProps {
    link: string
}

export function TelegramControl({ link, ...others }: TelegramControlProps) {
    return (
        <HeaderControl
            className={classes.telegram}
            component="a"
            href={link}
            rel="noopener noreferrer"
            target="_blank"
            {...others}
        >
            <PiTelegramLogoBold />
        </HeaderControl>
    )
}
