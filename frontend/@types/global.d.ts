declare global {
    interface Window {
        Go: typeof window.Go
        onWasmInitialized?: () => void

        SingboxParseConfig: (config: string) => null | string
    }
}

export {}
