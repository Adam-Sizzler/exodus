import { useLayoutEffect } from 'react'

import { preventBackScrollTables } from '@shared/utils/misc'

export function usePreventTableBackScroll() {
    useLayoutEffect(() => {
        document.body.addEventListener('wheel', preventBackScrollTables, { passive: false })
        return () => {
            document.body.removeEventListener('wheel', preventBackScrollTables)
        }
    }, [])
}
