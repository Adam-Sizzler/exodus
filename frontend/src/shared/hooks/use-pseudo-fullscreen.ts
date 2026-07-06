import { useDisclosure, useWindowEvent } from '@mantine/hooks'

export interface UsePseudoFullscreenReturn {
    close: () => void
    isFullscreen: boolean
    open: () => void
    toggle: () => void
}

export function usePseudoFullscreen(initial = false): UsePseudoFullscreenReturn {
    const [isFullscreen, { toggle, open, close }] = useDisclosure(initial)

    useWindowEvent('keydown', (event) => {
        if (isFullscreen && event.key === 'Escape') close()
    })

    return { close, isFullscreen, open, toggle }
}
