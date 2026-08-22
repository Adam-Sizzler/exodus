import { create } from 'zustand'
import { createJSONStorage, devtools, persist } from 'zustand/middleware'

import { HOSTS_VIEW_MODE, IActions, IState, LAYOUT_STYLE, NODES_VIEW_MODE } from './interfaces'

const initialState: IState = {
    nodesViewMode: NODES_VIEW_MODE.CARDS,
    nodesActiveTag: null,
    hostsViewMode: HOSTS_VIEW_MODE.CARDS,
    hostsActiveTag: null,
    layoutStyle: LAYOUT_STYLE.COMPACT
}

export const useViewPreferencesStore = create<IActions & IState>()(
    devtools(
        persist(
            (set) => ({
                ...initialState,
                actions: {
                    setNodesViewMode: (mode) => set({ nodesViewMode: mode }),
                    setNodesActiveTag: (tag) => set({ nodesActiveTag: tag }),
                    setHostsViewMode: (mode) => set({ hostsViewMode: mode }),
                    setHostsActiveTag: (tag) => set({ hostsActiveTag: tag }),
                    toggleLayoutStyle: () =>
                        set((state) => ({
                            layoutStyle:
                                state.layoutStyle === LAYOUT_STYLE.SIDEBAR
                                    ? LAYOUT_STYLE.COMPACT
                                    : LAYOUT_STYLE.SIDEBAR
                        })),
                    resetState: () => set(initialState)
                }
            }),
            {
                name: 'view-preferences-storage',
                storage: createJSONStorage(() => localStorage),
                version: 1,
                partialize: (state) => ({
                    nodesViewMode: state.nodesViewMode,
                    nodesActiveTag: state.nodesActiveTag,
                    hostsViewMode: state.hostsViewMode,
                    hostsActiveTag: state.hostsActiveTag,
                    layoutStyle: state.layoutStyle
                })
            }
        )
    )
)

export const useNodesViewMode = () => useViewPreferencesStore((state) => state.nodesViewMode)
export const useNodesActiveTag = () => useViewPreferencesStore((state) => state.nodesActiveTag)
export const useViewPreferencesStoreActions = () =>
    useViewPreferencesStore((state) => state.actions)
export const useHostsViewMode = () => useViewPreferencesStore((state) => state.hostsViewMode)
export const useHostsActiveTag = () => useViewPreferencesStore((state) => state.hostsActiveTag)
export const useLayoutStyle = () => useViewPreferencesStore((state) => state.layoutStyle)
export const useToggleLayoutStyleAction = () =>
    useViewPreferencesStore((state) => state.actions.toggleLayoutStyle)
