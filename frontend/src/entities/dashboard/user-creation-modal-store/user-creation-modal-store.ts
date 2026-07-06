import { devtools } from 'zustand/middleware'

import { queryClient } from '@shared/api'
import { QueryKeys } from '@shared/api/hooks'
import { create } from '@shared/hocs/store-wrapper'

import { IActions, IState } from './interfaces'

const initialState: IState = {
    isLoading: false,
    isModalOpen: false
}

export const useUserCreationModalStore = create<IActions & IState>()(
    devtools(
        (set) => ({
            ...initialState,
            actions: {
                changeModalState: async (modalState: boolean) => {
                    set(() => ({ isModalOpen: modalState }))
                },
                getInitialState: () => {
                    return initialState
                },
                resetState: async (): Promise<void> => {
                    await queryClient.refetchQueries({
                        queryKey: QueryKeys.users.getAllUsers._def
                    })
                    await queryClient.refetchQueries({ queryKey: QueryKeys.system._def })

                    set({ ...initialState })
                }
            }
        }),
        {
            name: 'userCreationModalStore',
            anonymousActionType: 'userCreationModalStore'
        }
    )
)

export const useUserCreationModalStoreIsModalOpen = () =>
    useUserCreationModalStore((state) => state.isModalOpen)

export const useUserCreationModalStoreActions = () =>
    useUserCreationModalStore((store) => store.actions)
