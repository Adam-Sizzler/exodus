import { useContext } from 'react'

import { IsMobileContext } from '@shared/hocs/is-mobile-provider'

export function useIsMobile(): boolean {
    return useContext(IsMobileContext)
}
