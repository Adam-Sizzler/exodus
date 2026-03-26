import { devtools } from 'zustand/middleware'

import { create } from '@shared/hocs/store-wrapper'

import { IActions, IState } from './interfaces'

const initialState: IState = {
    createModal: {
        isOpen: false,
        isLoading: false
    }
}

export const useSubscriptionConnectionsStore = create<IActions & IState>()(
    devtools(
        (set) => ({
            ...initialState,
            actions: {
                toggleCreateModal: (isOpen: boolean) => {
                    set((state) => ({
                        createModal: { ...state.createModal, isOpen }
                    }))
                    if (!isOpen) {
                        set((state) => ({
                            createModal: { ...state.createModal, isLoading: false }
                        }))
                    }
                },
                getInitialState: () => {
                    return initialState
                },
                resetState: async () => {
                    set({ ...initialState })
                }
            }
        }),
        {
            name: 'subscriptionConnectionsStore',
            anonymousActionType: 'subscriptionConnectionsStore'
        }
    )
)

export const useSubscriptionConnectionsStoreActions = () => useSubscriptionConnectionsStore((store) => store.actions)

// Create Modal
export const useSubscriptionConnectionsStoreCreateModalIsOpen = () =>
    useSubscriptionConnectionsStore((state) => state.createModal.isOpen)
export const useSubscriptionConnectionsStoreCreateModalIsLoading = () =>
    useSubscriptionConnectionsStore((state) => state.createModal.isLoading)
