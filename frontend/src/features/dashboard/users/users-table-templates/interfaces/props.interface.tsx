/* eslint-disable camelcase */
import { MRT_TableInstance } from '@kastov/mantine-react-table-open'
import { GetAllUsersCommand } from '@exodus/backend-contract'

export interface IProps {
    table: MRT_TableInstance<GetAllUsersCommand.Response['response']['users'][0]>
}
