import { OverflowList, OverflowListProps } from '@mantine/core'

interface IProps<T> extends Omit<OverflowListProps<T>, 'maxRows'> {
    rowHeight?: number
}

export function SingleRowOverflowList<T>({ rowHeight = 20, ...props }: IProps<T>) {
    return (
        <OverflowList<T>
            {...props}
            maxRows={1}
            styles={{ root: { maxHeight: rowHeight, overflow: 'hidden' } }}
        />
    )
}
