import { ScrollArea } from '@mantine/core'
import { useCallback, useEffect, useRef, useState } from 'react'

import styles from '../mobile/fade-mask.module.css'

interface FadeScrollAreaProps {
    children: React.ReactNode
    height: string
}

export function FadeScrollArea(props: FadeScrollAreaProps) {
    const { children, height } = props
    const viewportRef = useRef<HTMLDivElement>(null)
    const [fade, setFade] = useState({ bottom: false, top: false })

    const updateFade = useCallback(() => {
        const el = viewportRef.current
        if (!el) {
            return
        }

        const { scrollTop, scrollHeight, clientHeight } = el
        const isScrollable = scrollHeight - clientHeight > 1

        setFade({
            top: isScrollable && scrollTop > 4,
            bottom: isScrollable && scrollTop + clientHeight < scrollHeight - 4
        })
    }, [])

    useEffect(() => {
        updateFade()

        const el = viewportRef.current
        if (!el) {
            return undefined
        }

        const observer = new ResizeObserver(updateFade)
        observer.observe(el)
        if (el.firstElementChild) {
            observer.observe(el.firstElementChild)
        }

        return () => observer.disconnect()
    }, [updateFade, height])

    const viewportClassName =
        [fade.top && styles.fadeTop, fade.bottom && styles.fadeBottom].filter(Boolean).join(' ') ||
        undefined

    return (
        <ScrollArea
            classNames={{ viewport: viewportClassName }}
            h={height}
            offsetScrollbars
            onScrollPositionChange={updateFade}
            type="auto"
            viewportRef={viewportRef}
        >
            {children}
        </ScrollArea>
    )
}
