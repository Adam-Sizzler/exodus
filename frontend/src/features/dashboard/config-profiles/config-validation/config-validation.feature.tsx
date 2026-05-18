import type { editor } from 'monaco-editor'

import { Monaco } from '@monaco-editor/react'
import { RefObject } from 'react'
import dayjs from 'dayjs'

export const ConfigValidationFeature = {
    validate: (
        editorRef: RefObject<editor.IStandaloneCodeEditor | null>,
        monacoRef: RefObject<Monaco | null>,
        setResult: (message: string) => void,
        setIsConfigValid: (isValid: boolean) => void
    ) => {
        try {
            if (!editorRef.current) return
            if (!monacoRef.current) return

            const currentValue = editorRef.current.getValue()

            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            let clonedCurrentValue: any
            try {
                clonedCurrentValue = JSON.parse(currentValue)
            } catch {
                setResult(`${dayjs().format('HH:mm:ss')} | Invalid JSON.`)
                setIsConfigValid(false)
                return
            }

            if (typeof window.SingboxParseConfig !== 'function') {
                setResult(
                    `${dayjs().format('HH:mm:ss')} | WASM validator is unavailable (SingboxParseConfig is not initialized).`
                )
                setIsConfigValid(false)
                return
            }

            const validationResult = window.SingboxParseConfig(JSON.stringify(clonedCurrentValue))

            setResult(
                `${dayjs().format('HH:mm:ss')} | ${validationResult || 'Sing-box config is valid.'}`
            )
            setIsConfigValid(!validationResult)
        } catch (err: unknown) {
            setResult(`${dayjs().format('HH:mm:ss')} | Validation error: ${(err as Error).message}`)
            setIsConfigValid(false)
        }
    }
}
