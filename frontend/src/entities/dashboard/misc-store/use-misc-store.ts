import { create } from 'zustand'
import { createJSONStorage, devtools, persist } from 'zustand/middleware'

import { IActions, IState } from './interfaces'

export const useMiscStore = create<IActions & IState>()(
    persist(
        devtools(
            (set) => ({
                disclaimerAccepted: false,
                mobileWarningClosed: false,
                srrAdvancedModalClosed: false,
                actions: {
                    setDisclaimerAccepted: (accepted: boolean) => {
                        set({ disclaimerAccepted: accepted })
                    },
                    setMobileWarningClosed: (closed: boolean) => {
                        set({ mobileWarningClosed: closed })
                    },
                    setSrrAdvancedModalClosed: (closed: boolean) => {
                        set({ srrAdvancedModalClosed: closed })
                    }
                }
            }),
            { name: 'miscStore', anonymousActionType: 'miscStore' }
        ),
        {
            name: 'miscStore',
            partialize: (state) => ({
                disclaimerAccepted: state.disclaimerAccepted,
                mobileWarningClosed: state.mobileWarningClosed,
                srrAdvancedModalClosed: state.srrAdvancedModalClosed
            }),
            migrate: (persistedState: unknown, version: number) => {
                if (version === 1) {
                    return {
                        disclaimerAccepted: false,
                        mobileWarningClosed:
                            (persistedState as IState).mobileWarningClosed ?? false,
                        srrAdvancedModalClosed: false
                    }
                }
                if (version === 2) {
                    return {
                        ...(persistedState as IState),
                        disclaimerAccepted: false
                    }
                }
                return persistedState
            },
            storage: createJSONStorage(() => localStorage),
            version: 3
        }
    )
)

export const useDisclaimerAccepted = () => useMiscStore((state) => state.disclaimerAccepted)
export const useMobileWarningClosed = () => useMiscStore((state) => state.mobileWarningClosed)
export const useSrrAdvancedModalClosed = () => useMiscStore((state) => state.srrAdvancedModalClosed)
export const useMiscStoreActions = () => useMiscStore((state) => state.actions)
