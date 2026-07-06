import { Button, Checkbox, Divider, Group, Popover, ScrollArea, Stack, Text } from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { type DataTableColumnToggle } from 'mantine-datatable'
import { useTranslation } from 'react-i18next'
import { TbColumns3, TbEye, TbEyeOff } from 'react-icons/tb'

interface IProps {
    columnsToggle: DataTableColumnToggle[]
    labelByAccessor: Record<string, string>
    setColumnsToggle: (toggle: DataTableColumnToggle[]) => void
}

export function ColumnsVisibilityPopover(props: IProps) {
    const { t } = useTranslation()
    const { columnsToggle, labelByAccessor, setColumnsToggle } = props
    const [opened, { close, toggle }] = useDisclosure(false)

    const toggleableColumns = columnsToggle.filter((column) => column.toggleable)
    const visibleCount = toggleableColumns.filter((column) => column.toggled).length

    const setToggled = (accessor: string, toggled: boolean) => {
        setColumnsToggle(
            columnsToggle.map((column) =>
                column.accessor === accessor ? { ...column, toggled } : column
            )
        )
    }

    const setAll = (toggled: boolean) => {
        setColumnsToggle(
            columnsToggle.map((column) => (column.toggleable ? { ...column, toggled } : column))
        )
    }

    return (
        <Popover
            keepMounted={false}
            onClose={close}
            opened={opened}
            position="top-end"
            shadow="md"
            width={280}
            withArrow
        >
            <Popover.Target>
                <Button
                    leftSection={<TbColumns3 size={16} />}
                    onClick={toggle}
                    size="sm"
                    variant="default"
                >
                    {visibleCount}/{toggleableColumns.length}
                </Button>
            </Popover.Target>
            <Popover.Dropdown>
                <Group justify="space-between" mb="xs">
                    <Button
                        leftSection={<TbEye size={14} />}
                        onClick={() => setAll(true)}
                        size="compact-xs"
                        variant="subtle"
                    >
                        {t('columns-visibility-popover.show-all')}
                    </Button>
                    <Button
                        color="gray"
                        leftSection={<TbEyeOff size={14} />}
                        onClick={() => setAll(false)}
                        size="compact-xs"
                        variant="subtle"
                    >
                        {t('columns-visibility-popover.hide-all')}
                    </Button>
                </Group>
                <Divider mb="xs" />
                <ScrollArea.Autosize mah={320} type="auto">
                    <Stack gap="xs" pr="sm">
                        {toggleableColumns.map((column) => (
                            <Checkbox
                                checked={column.toggled}
                                key={column.accessor}
                                label={
                                    <Text size="sm">
                                        {labelByAccessor[column.accessor] ?? column.accessor}
                                    </Text>
                                }
                                onChange={(event) =>
                                    setToggled(column.accessor, event.currentTarget.checked)
                                }
                            />
                        ))}
                    </Stack>
                </ScrollArea.Autosize>
            </Popover.Dropdown>
        </Popover>
    )
}
