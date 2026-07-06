import { modals } from '@mantine/modals'
import { IconCrownFilled } from '@tabler/icons-react'

import { PrimeModalContent } from '../prime-modal/prime-modal.shared'
import { HeaderControl } from './HeaderControl'
import classes from './PrimeControl.module.css'

export function PrimeControl() {
    const openPrimeModal = () => {
        modals.open({
            centered: true,
            size: 460,
            withCloseButton: false,
            children: <PrimeModalContent />
        })
    }

    return (
        <HeaderControl className={classes.support} onClick={openPrimeModal}>
            <IconCrownFilled />
        </HeaderControl>
    )
}
