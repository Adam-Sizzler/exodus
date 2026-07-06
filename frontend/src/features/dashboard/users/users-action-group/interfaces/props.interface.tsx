import { MRT_TableInstance } from '@kastov/mantine-react-table-open'
/* eslint-disable camelcase */
import { GetAllUsersCommand } from '@exodus/backend-contract'

export interface IProps {
    isLoading: boolean
    refetch: () => void

    table: MRT_TableInstance<GetAllUsersCommand.Response['response']['users'][0]>
}
