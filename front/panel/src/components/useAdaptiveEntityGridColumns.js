import { useCallback, useEffect, useState } from 'react';

const SINGLE_COLUMN_BREAKPOINT_PX = 2560 / 3;
const WIDTH_PER_COLUMN_PX = 380;
const MIN_COLUMNS = 2;
const MAX_COLUMNS = 12;

function getAdaptiveColumns(containerWidth) {
  if (!containerWidth || containerWidth <= SINGLE_COLUMN_BREAKPOINT_PX) {
    return 1;
  }

  const adaptiveColumns = Math.floor(containerWidth / WIDTH_PER_COLUMN_PX);
  return Math.max(MIN_COLUMNS, Math.min(MAX_COLUMNS, adaptiveColumns));
}

export default function useAdaptiveEntityGridColumns() {
  const [gridElement, setGridElement] = useState(null);
  const [columns, setColumns] = useState(1);
  const gridRef = useCallback((node) => {
    setGridElement(node);
  }, []);

  useEffect(() => {
    if (!gridElement) {
      return undefined;
    }

    const updateColumns = () => {
      const nextColumns = getAdaptiveColumns(gridElement.clientWidth);
      setColumns((prevColumns) => (prevColumns === nextColumns ? prevColumns : nextColumns));
    };

    updateColumns();

    const resizeObserver = new ResizeObserver(updateColumns);
    resizeObserver.observe(gridElement);
    window.addEventListener('resize', updateColumns);

    return () => {
      window.removeEventListener('resize', updateColumns);
      resizeObserver.disconnect();
    };
  }, [gridElement]);

  return { gridRef, columns };
}
