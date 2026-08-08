import { MRT_Column } from '@kastov/mantine-react-table-open'
/* eslint-disable camelcase */
import { Group, NativeSelect, NumberInput } from '@mantine/core'
import { useDebouncedValue } from '@mantine/hooks'
import { GetUsersCommand } from '@exodus/backend-contract'
import { useEffect, useRef, useState } from 'react'

import {
    bestFitIecUnitUtil,
    bytesToUnitUtil,
    IEC_UNITS,
    TIecUnit,
    unitToBytesUtil
} from '@shared/utils/bytes'

type TUser = GetUsersCommand.Response['response']['users'][number]
type TRange = [null | string, null | string]
type TValues = [number | string, number | string]

interface IProps {
    column: MRT_Column<TUser>
}

const FILTER_DEBOUNCE_MS = 400

const toBytes = (bound: null | string): null | number => {
    if (bound === null || bound === '') return null
    const num = Number(bound)
    return Number.isFinite(num) ? num : null
}

export const TrafficLimitRangeFilter = ({ column }: IProps) => {
    const filterValue = column.getFilterValue() as TRange | undefined

    const [unit, setUnit] = useState<TIecUnit>(() => {
        const [min, max] = filterValue ?? [null, null]
        return bestFitIecUnitUtil(toBytes(min) ?? toBytes(max))
    })

    const [values, setValues] = useState<TValues>(() => {
        const [min, max] = filterValue ?? [null, null]
        const fitUnit = bestFitIecUnitUtil(toBytes(min) ?? toBytes(max))
        return [
            bytesToUnitUtil(toBytes(min), fitUnit) ?? '',
            bytesToUnitUtil(toBytes(max), fitUnit) ?? ''
        ]
    })

    const [prevFilterValue, setPrevFilterValue] = useState(filterValue)
    if (filterValue !== prevFilterValue) {
        setPrevFilterValue(filterValue)
        if (filterValue === undefined && (values[0] !== '' || values[1] !== '')) {
            setValues(['', ''])
        }
    }

    const [debouncedValues] = useDebouncedValue(values, FILTER_DEBOUNCE_MS)
    const [debouncedUnit] = useDebouncedValue(unit, FILTER_DEBOUNCE_MS)
    const isMounted = useRef(false)

    useEffect(() => {
        if (!isMounted.current) {
            isMounted.current = true
            return
        }

        const toBound = (value: number | string): null | string => {
            const bytes = unitToBytesUtil(value, debouncedUnit)
            return bytes === undefined ? null : String(bytes)
        }

        const next: TRange = [toBound(debouncedValues[0]), toBound(debouncedValues[1])]
        const normalized = next[0] === null && next[1] === null ? undefined : next

        const current = column.getFilterValue() as TRange | undefined
        const unchanged =
            (current === undefined && normalized === undefined) ||
            (Array.isArray(current) &&
                Array.isArray(normalized) &&
                current[0] === normalized[0] &&
                current[1] === normalized[1])

        if (!unchanged) {
            column.setFilterValue(normalized)
        }
    }, [debouncedValues, debouncedUnit])

    return (
        <Group align="center" gap={6} style={{ gridColumn: '1 / -1' }} w="100%" wrap="nowrap">
            <NumberInput
                allowNegative={false}
                flex={1}
                hideControls
                miw={0}
                onChange={(value) => setValues(([, currentMax]) => [value, currentMax])}
                placeholder="Min"
                size="xs"
                thousandSeparator=","
                value={values[0]}
            />
            <NumberInput
                allowNegative={false}
                flex={1}
                hideControls
                miw={0}
                onChange={(value) => setValues(([currentMin]) => [currentMin, value])}
                placeholder="Max"
                size="xs"
                thousandSeparator=","
                value={values[1]}
            />
            <NativeSelect
                aria-label="Unit"
                data={IEC_UNITS}
                onChange={(event) => setUnit(event.currentTarget.value as TIecUnit)}
                size="xs"
                styles={{ input: { fontWeight: 600 } }}
                value={unit}
                w={72}
            />
        </Group>
    )
}
