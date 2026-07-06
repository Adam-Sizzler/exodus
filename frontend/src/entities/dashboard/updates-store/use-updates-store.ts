import axios from 'axios'
import { create } from 'zustand'
import { createJSONStorage, devtools, persist } from 'zustand/middleware'

import { sToMs } from '@shared/utils/time-utils'
import { findLatestStableVersionTag } from '@shared/utils/version-utils'

const CACHE_TIME = sToMs(24 * 60 * 60)
const EXODUS_GITHUB_OWNER = 'Adam-Sizzler'
const EXODUS_GITHUB_REPO = 'exodus'
const EXODUS_GITHUB_API_REPO = `https://api.github.com/repos/${EXODUS_GITHUB_OWNER}/${EXODUS_GITHUB_REPO}`
const GITHUB_TAGS_PER_PAGE = 100
const GITHUB_TAGS_MAX_PAGES = 5

interface GithubTag {
    name: string
}

async function getGithubTags() {
    const tags: GithubTag[] = []

    for (let page = 1; page <= GITHUB_TAGS_MAX_PAGES; page += 1) {
        const response = await axios.get<GithubTag[]>(`${EXODUS_GITHUB_API_REPO}/tags`, {
            headers: {
                Accept: 'application/vnd.github+json'
            },
            params: {
                page,
                per_page: GITHUB_TAGS_PER_PAGE
            }
        })

        tags.push(...response.data)

        if (response.data.length < GITHUB_TAGS_PER_PAGE) {
            break
        }
    }

    return tags
}

export interface IExodusInfo {
    latestVersion: string
    starsCount: number
}

interface IState {
    isLoading: boolean
    lastUpdateTimestamp: number
    exodusInfo: IExodusInfo
}

interface IActions {
    actions: {
        getExodusInfo: () => Promise<void>
        resetState: () => void
        setExodusInfo: (info: IExodusInfo) => void
    }
}

const initialState: IState = {
    isLoading: false,
    lastUpdateTimestamp: 0,
    exodusInfo: {
        latestVersion: '0.0.0',
        starsCount: 0
    }
}

export const useUpdatesStore = create<IActions & IState>()(
    persist(
        devtools(
            (set, get) => ({
                ...initialState,
                actions: {
                    getExodusInfo: async () => {
                        const { lastUpdateTimestamp, exodusInfo } = get()
                        const now = Date.now()

                        if (
                            lastUpdateTimestamp &&
                            now - lastUpdateTimestamp < CACHE_TIME &&
                            exodusInfo.latestVersion &&
                            exodusInfo.starsCount > 0
                        ) {
                            return
                        }

                        try {
                            set({ isLoading: true })

                            const [starsResponse, tags] = await Promise.all([
                                axios.get<{
                                    totalStars: number
                                }>('https://ungh.cc/stars/Adam-Sizzler/exodus'),
                                getGithubTags()
                            ])

                            const latestVersion = findLatestStableVersionTag(
                                tags.map((tag) => tag.name)
                            )

                            set({
                                exodusInfo: {
                                    latestVersion: latestVersion ?? exodusInfo.latestVersion,
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

                    setExodusInfo: (info: IExodusInfo) => {
                        set({ exodusInfo: info, lastUpdateTimestamp: Date.now() })
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
            version: 2,
            migrate: (persistedState) => ({
                ...initialState,
                ...(typeof persistedState === 'object' && persistedState !== null
                    ? persistedState
                    : {})
            }),
            partialize: (state) => ({
                lastUpdateTimestamp: state.lastUpdateTimestamp,
                exodusInfo: state.exodusInfo
            })
        }
    )
)

export const useExodusInfo = () => useUpdatesStore((state) => state.exodusInfo)
export const useLastUpdateTimestamp = () => useUpdatesStore((state) => state.lastUpdateTimestamp)
export const useIsLoadingExodusUpdates = () => useUpdatesStore((state) => state.isLoading)
export const useUpdatesStoreActions = () => useUpdatesStore((state) => state.actions)
