import { GetAllUsersCommand } from '@cerberus/backend-contract'

export type User = GetAllUsersCommand.Response['response']['users'][number]
