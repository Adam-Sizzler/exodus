import { createJSONStorage, devtools, persist } from 'zustand/middleware'
import { create } from 'zustand'
import axios from 'axios'

import { sToMs } from '@shared/utils/time-utils'

const CACHE_TIME = sToMs(24 * 60 * 60)

interface ICerberusInfo {
    latestVersion: string
    starsCount: number
}

interface IState {
    isLoading: boolean
    lastUpdateTimestamp: number
    cerberusInfo: ICerberusInfo
}

interface IActions {
    actions: {
        getCerberusInfo: () => Promise<void>
        resetState: () => void
        setCerberusInfo: (info: ICerberusInfo) => void
    }
}

const initialState: IState = {
    isLoading: false,
    lastUpdateTimestamp: 0,
    cerberusInfo: {
        latestVersion: '1.12.0',
        starsCount: 0
    }
}

export const useUpdatesStore = create<IActions & IState>()(
    persist(
        devtools(
            (set, get) => ({
                ...initialState,
                actions: {
                    getCerberusInfo: async () => {
                        const { lastUpdateTimestamp, cerberusInfo } = get()
                        const now = Date.now()

                        if (
                            lastUpdateTimestamp &&
                            now - lastUpdateTimestamp < CACHE_TIME &&
                            cerberusInfo.latestVersion &&
                            cerberusInfo.starsCount > 0
                        ) {
                            return
                        }

                        try {
                            set({ isLoading: true })

                            const starsResponse = await axios.get<{
                                totalStars: number
                            }>('https://ungh.cc/stars/SagerNet/sing-box')

                            const versionResponse = await axios.get<{
                                release: {
                                    tag: string
                                }
                            }>('https://ungh.cc/repos/SagerNet/sing-box/releases/latest')

                            set({
                                cerberusInfo: {
                                    latestVersion: versionResponse.data.release.tag,
                                    starsCount: starsResponse.data.totalStars
                                },
                                lastUpdateTimestamp: now
                            })
                        } catch {
                            // silent error
                        } finally {
                            set({ isLoading: false })
                        }
                    },

                    setCerberusInfo: (info: ICerberusInfo) => {
                        set({ cerberusInfo: info, lastUpdateTimestamp: Date.now() })
                    },
                    resetState: () => {
                        set({ ...initialState })
                    }
                }
            }),
            { name: 'updatesStore', anonymousActionType: 'updatesStore' }
        ),
        {
            name: 'updatesStore',
            storage: createJSONStorage(() => localStorage),
            version: 1,
            partialize: (state) => ({
                lastUpdateTimestamp: state.lastUpdateTimestamp,
                cerberusInfo: state.cerberusInfo
            })
        }
    )
)

export const useCerberusInfo = () => useUpdatesStore((state) => state.cerberusInfo)
export const useLastUpdateTimestamp = () => useUpdatesStore((state) => state.lastUpdateTimestamp)
export const useIsLoadingCerberusUpdates = () => useUpdatesStore((state) => state.isLoading)
export const useUpdatesStoreActions = () => useUpdatesStore((state) => state.actions)
