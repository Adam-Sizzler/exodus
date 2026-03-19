import { GetConfigProfileByUuidCommand, GetSnippetsCommand } from '@cerberus/backend-contract'

export interface IProps {
    configProfile: GetConfigProfileByUuidCommand.Response['response']
    snippets: GetSnippetsCommand.Response['response']
}
