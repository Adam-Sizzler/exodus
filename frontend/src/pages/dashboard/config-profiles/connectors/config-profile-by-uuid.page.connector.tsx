import { Navigate, useParams } from 'react-router-dom'
import { useLayoutEffect, useState } from 'react'
import { consola } from 'consola/browser'

import { useGetConfigProfile } from '@shared/api/hooks'
import { fetchWithProgress } from '@shared/utils/fetch-with-progress'
import { ROUTES } from '@shared/constants'
import { LoadingScreen } from '@shared/ui'
import { app } from 'src/config'

import { ConfigProfileByUuidPageComponent } from '../components/config-profile-by-uuid.page.component'

export function ConfigProfileByUuidPageConnector() {
    const { uuid } = useParams()

    const [downloadProgress, setDownloadProgress] = useState(0)
    const [isLoading, setIsLoading] = useState(true)

    const { data: configProfile, isLoading: isConfigProfileLoading } = useGetConfigProfile({
        route: { uuid: uuid! },
        rQueryParams: {
            enabled: !!uuid,
            refetchOnWindowFocus: false
        }
    })

    useLayoutEffect(() => {
        const initWasm = async () => {
            try {
                if (typeof window.Go !== 'function') {
                    throw new Error('Go WASM runtime is not loaded. Check /assets/wasm_exec.js')
                }

                const go = new window.Go()
                const wasmInitialized = new Promise<void>((resolve) => {
                    const timeoutId = window.setTimeout(() => {
                        resolve()
                    }, 5000)

                    window.onWasmInitialized = () => {
                        window.clearTimeout(timeoutId)
                        consola.info('WASM module initialized')
                        resolve()
                    }
                })

                const wasmBytes: ArrayBuffer = await fetchWithProgress(
                    app.configEditor.wasmUrl,
                    setDownloadProgress
                ) as ArrayBuffer
                const wasmUint8 = new Uint8Array(wasmBytes)
                if (wasmUint8.byteLength < 8) {
                    throw new Error('Downloaded WASM is too small')
                }
                const wasmMagicOk =
                    wasmUint8[0] === 0x00 &&
                    wasmUint8[1] === 0x61 &&
                    wasmUint8[2] === 0x73 &&
                    wasmUint8[3] === 0x6d
                if (!wasmMagicOk) {
                    throw new Error(
                        'Downloaded file is not a valid wasm binary (magic number mismatch). Check /assets/main.wasm source.'
                    )
                }

                const { instance } = await WebAssembly.instantiate(wasmUint8, go.importObject)
                go.run(instance)
                await wasmInitialized

                if (typeof window.SingboxParseConfig === 'function') {
                    setIsLoading(false)
                } else {
                    throw new Error(
                        'SingboxParseConfig is not initialized. Ensure your wasm exports this function and calls window.onWasmInitialized().'
                    )
                }
            } catch (err: unknown) {
                consola.error('WASM initialization error:', err)
                setIsLoading(false)
            }
        }

        initWasm()

        return () => {
            delete window.onWasmInitialized
        }
    }, [])

    if (!uuid) {
        return <Navigate to={ROUTES.DASHBOARD.MANAGEMENT.CONFIG_PROFILES} />
    }

    if (isLoading || isConfigProfileLoading || !configProfile) {
        return <LoadingScreen text="WASM module is loading..." value={downloadProgress} />
    }

    return <ConfigProfileByUuidPageComponent configProfile={configProfile} />
}
