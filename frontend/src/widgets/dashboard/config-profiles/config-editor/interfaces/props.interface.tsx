import { GetConfigProfileByUuidCommand, GetSnippetsCommand } from '@exodus/backend-contract'

export interface IProps {
    configProfile: GetConfigProfileByUuidCommand.Response['response']
    snippets: GetSnippetsCommand.Response['response']
}
