import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import AppSelect from './AppSelect';

const DEFAULT_PAGE_SIZE_OPTIONS = [5, 10, 15, 20, 25, 30, 50, 100];
const EMPTY_SET = new Set();

const buildColumnMaps = (columns) => {
  const byKey = {};
  const defaultWidths = {};
  const minWidths = {};
  const defaultVisibility = {};
  const defaultPinning = {};

  columns.forEach((column) => {
    byKey[column.key] = column;
    defaultWidths[column.key] = column.defaultWidth ?? 120;
    minWidths[column.key] = column.minWidth ?? 80;
    defaultVisibility[column.key] = Boolean(column.defaultVisible || column.alwaysVisible);
    defaultPinning[column.key] = column.defaultPin || null;
  });

  return {
    byKey,
    defaultWidths,
    minWidths,
    defaultVisibility,
    defaultPinning,
  };
};

function SharedRichTable({
  columns,
  rows,
  getRowId,
  renderCell,
  onRowClick,
  getRowAriaLabel,
  selectedRowIds = EMPTY_SET,
  setSelectedRowIds,
  sortState,
  setSortState,
  compactRows,
  setCompactRows,
  isFullscreen,
  setIsFullscreen,
  searchInputRef,
  pageResetKey,
  pageSizeOptions = DEFAULT_PAGE_SIZE_OPTIONS,
  defaultPageSize = 25,
  pageSizeLabelId,
  pageSizeLabel = 'Rows per page',
  tableAriaLabel = 'Таблица',
  emptyState = null,
}) {
  const columnMaps = useMemo(() => buildColumnMaps(columns), [columns]);
  const [columnWidths, setColumnWidths] = useState(columnMaps.defaultWidths);
  const [columnVisibility, setColumnVisibility] = useState(columnMaps.defaultVisibility);
  const [columnPinning, setColumnPinning] = useState(columnMaps.defaultPinning);
  const [isColumnsMenuOpen, setIsColumnsMenuOpen] = useState(false);
  const [pageSize, setPageSize] = useState(defaultPageSize);
  const [pageIndex, setPageIndex] = useState(0);
  const [resizeGuideLeft, setResizeGuideLeft] = useState(null);

  const resizeStateRef = useRef(null);
  const columnsMenuRef = useRef(null);
  const columnsMenuButtonRef = useRef(null);
  const tableScrollRef = useRef(null);

  useEffect(() => {
    if (!isColumnsMenuOpen) return undefined;

    const handleOutsideClick = (event) => {
      if (columnsMenuRef.current?.contains(event.target)) return;
      if (columnsMenuButtonRef.current?.contains(event.target)) return;
      setIsColumnsMenuOpen(false);
    };

    const handleEscape = (event) => {
      if (event.key === 'Escape') {
        setIsColumnsMenuOpen(false);
      }
    };

    document.addEventListener('mousedown', handleOutsideClick);
    window.addEventListener('keydown', handleEscape);

    return () => {
      document.removeEventListener('mousedown', handleOutsideClick);
      window.removeEventListener('keydown', handleEscape);
    };
  }, [isColumnsMenuOpen]);

  useEffect(() => {
    if (pageResetKey === undefined) return;
    setPageIndex(0);
  }, [pageResetKey]);

  useEffect(() => {
    if (!sortState?.key || !setSortState) return;
    if (columnVisibility[sortState.key]) return;
    setSortState({ key: null, direction: 'asc' });
  }, [columnVisibility, sortState?.key, setSortState]);

  const visibleColumns = useMemo(() => {
    return columns.filter((column) => columnVisibility[column.key]);
  }, [columns, columnVisibility]);

  const renderedColumns = useMemo(() => {
    const left = [];
    const center = [];
    const right = [];

    visibleColumns.forEach((column) => {
      const pin = column.disablePinning ? 'left' : columnPinning[column.key];
      if (pin === 'left') {
        left.push(column);
      } else if (pin === 'right') {
        right.push(column);
      } else {
        center.push(column);
      }
    });

    return [...left, ...center, ...right];
  }, [visibleColumns, columnPinning]);

  const toggleSort = useCallback((key) => {
    if (!setSortState) return;
    const column = columnMaps.byKey[key];
    if (!column?.sortable) {
      return;
    }

    setSortState((prev) => {
      if (prev.key !== key) {
        return { key, direction: 'asc' };
      }
      if (prev.direction === 'asc') {
        return { key, direction: 'desc' };
      }
      return { key: null, direction: 'asc' };
    });
  }, [columnMaps.byKey, setSortState]);

  const getSortDirection = useCallback(
    (key) => (sortState?.key === key ? sortState.direction : null),
    [sortState],
  );

  const updateResizeGuide = useCallback((clientX) => {
    const scrollEl = tableScrollRef.current;
    if (!scrollEl) return;

    const rect = scrollEl.getBoundingClientRect();
    const left = clientX - rect.left + scrollEl.scrollLeft;
    setResizeGuideLeft(Math.max(0, left));
  }, []);

  const handleColumnResizeMove = useCallback((event) => {
    const state = resizeStateRef.current;
    if (!state) return;

    const delta = event.clientX - state.startX;
    const minWidth = columnMaps.minWidths[state.key] ?? 80;
    const nextWidth = Math.max(minWidth, Math.round(state.startWidth + delta));

    setColumnWidths((prev) => {
      if (prev[state.key] === nextWidth) return prev;
      return { ...prev, [state.key]: nextWidth };
    });
    updateResizeGuide(event.clientX);
  }, [columnMaps.minWidths, updateResizeGuide]);

  const stopColumnResize = useCallback(() => {
    resizeStateRef.current = null;
    setResizeGuideLeft(null);
    document.body.classList.remove('users-col-resizing');
    window.removeEventListener('mousemove', handleColumnResizeMove);
    window.removeEventListener('mouseup', stopColumnResize);
  }, [handleColumnResizeMove]);

  const startColumnResize = useCallback((event, key) => {
    event.preventDefault();
    event.stopPropagation();

    const startWidth = columnWidths[key] ?? columnMaps.defaultWidths[key] ?? 120;
    resizeStateRef.current = {
      key,
      startX: event.clientX,
      startWidth,
    };

    document.body.classList.add('users-col-resizing');
    updateResizeGuide(event.clientX);
    window.addEventListener('mousemove', handleColumnResizeMove);
    window.addEventListener('mouseup', stopColumnResize);
  }, [columnWidths, columnMaps.defaultWidths, handleColumnResizeMove, stopColumnResize, updateResizeGuide]);

  useEffect(() => () => stopColumnResize(), [stopColumnResize]);

  const tableMinWidth = useMemo(() => {
    return renderedColumns.reduce((sum, column) => {
      return sum + (columnWidths[column.key] ?? columnMaps.defaultWidths[column.key] ?? 120);
    }, 0);
  }, [columnMaps.defaultWidths, columnWidths, renderedColumns]);

  const rowHeightRem = compactRows ? '7rem' : '8rem';

  const totalRows = rows.length;
  const totalPages = Math.max(1, Math.ceil(totalRows / pageSize));
  const normalizedPageIndex = Math.min(pageIndex, totalPages - 1);
  const pageStartOffset = normalizedPageIndex * pageSize;
  const paginatedRows = useMemo(
    () => rows.slice(pageStartOffset, pageStartOffset + pageSize),
    [pageSize, pageStartOffset, rows],
  );
  const pageRangeStart = totalRows === 0 ? 0 : pageStartOffset + 1;
  const pageRangeEnd = totalRows === 0 ? 0 : Math.min(pageStartOffset + pageSize, totalRows);

  useEffect(() => {
    if (pageIndex !== normalizedPageIndex) {
      setPageIndex(normalizedPageIndex);
    }
  }, [pageIndex, normalizedPageIndex]);

  const toggleRowSelection = useCallback((rowId) => {
    if (!setSelectedRowIds) return;

    setSelectedRowIds((prev) => {
      const next = new Set(prev);
      if (next.has(rowId)) {
        next.delete(rowId);
      } else {
        next.add(rowId);
      }
      return next;
    });
  }, [setSelectedRowIds]);

  const toggleSelectAllVisible = useCallback((checked) => {
    if (!setSelectedRowIds) return;

    setSelectedRowIds((prev) => {
      const next = new Set(prev);
      paginatedRows.forEach((row) => {
        const rowId = getRowId(row);
        if (checked) {
          next.add(rowId);
        } else {
          next.delete(rowId);
        }
      });
      return next;
    });
  }, [getRowId, paginatedRows, setSelectedRowIds]);

  const allVisibleSelected = paginatedRows.length > 0 && paginatedRows.every((row) => selectedRowIds.has(getRowId(row)));

  const allToggleableVisible = useMemo(() => {
    return columns.every((column) => column.alwaysVisible || columnVisibility[column.key]);
  }, [columnVisibility, columns]);

  const hasVisibleToggleableColumns = useMemo(() => {
    return columns.some((column) => !column.alwaysVisible && columnVisibility[column.key]);
  }, [columnVisibility, columns]);

  const hasPinnedColumns = useMemo(() => {
    return columns.some((column) => !column.disablePinning && columnPinning[column.key]);
  }, [columnPinning, columns]);

  const toggleColumnVisibility = (key, checked) => {
    const column = columnMaps.byKey[key];
    if (!column || column.alwaysVisible) return;

    setColumnVisibility((prev) => ({
      ...prev,
      [key]: checked,
    }));
  };

  const hideAllColumns = () => {
    setColumnVisibility((prev) => {
      const next = { ...prev };
      columns.forEach((column) => {
        if (!column.alwaysVisible) {
          next[column.key] = false;
        }
      });
      return next;
    });
  };

  const showAllColumns = () => {
    setColumnVisibility((prev) => {
      const next = { ...prev };
      columns.forEach((column) => {
        next[column.key] = true;
      });
      return next;
    });
  };

  const setColumnPin = (key, position) => {
    const column = columnMaps.byKey[key];
    if (!column || column.disablePinning) return;

    setColumnPinning((prev) => {
      const current = prev[key];
      return {
        ...prev,
        [key]: current === position ? null : position,
      };
    });
  };

  const clearColumnPin = (key) => {
    const column = columnMaps.byKey[key];
    if (!column || column.disablePinning) return;

    setColumnPinning((prev) => ({
      ...prev,
      [key]: null,
    }));
  };

  const unpinAllColumns = () => {
    setColumnPinning((prev) => {
      const next = { ...prev };
      columns.forEach((column) => {
        if (!column.disablePinning) {
          next[column.key] = null;
        }
      });
      return next;
    });
  };

  const renderSortIcon = (key) => {
    const direction = getSortDirection(key);
    return (
      <span className={`users-th-sort-icon ${direction ? `is-${direction}` : 'is-none'}`} aria-hidden="true">
        {direction === 'asc' ? '↑' : direction === 'desc' ? '↓' : '↕'}
      </span>
    );
  };

  const renderHeaderCell = (column) => {
    if (column.key === 'select' && setSelectedRowIds) {
      return (
        <th key={column.key} className="users-cell-center users-cell-checkbox">
          <div className="users-th-wrap users-th-wrap-checkbox">
            <div className="users-checkbox-wrap">
              <input
                type="checkbox"
                checked={allVisibleSelected}
                onChange={(event) => toggleSelectAllVisible(event.target.checked)}
                aria-label="Выбрать все видимые"
                onClick={(event) => event.stopPropagation()}
              />
            </div>
            <button
              type="button"
              className="users-col-resizer"
              onMouseDown={(event) => startColumnResize(event, column.key)}
              onClick={(event) => event.stopPropagation()}
              aria-label="Изменить ширину колонки выбора"
              title="Изменить ширину"
            />
          </div>
        </th>
      );
    }

    return (
      <th key={column.key}>
        <div className="users-th-wrap">
          {column.sortable ? (
            <button
              type="button"
              className={`users-th-sort ${getSortDirection(column.key) ? 'is-active' : ''}`}
              onClick={() => toggleSort(column.key)}
            >
              {column.label} {renderSortIcon(column.key)}
            </button>
          ) : (
            <span className="users-th-label">{column.label}</span>
          )}
          <button
            type="button"
            className="users-col-resizer"
            onMouseDown={(event) => startColumnResize(event, column.key)}
            onClick={(event) => event.stopPropagation()}
            aria-label={`Изменить ширину колонки ${column.label}`}
            title={`Изменить ширину: ${column.label}`}
          />
        </div>
      </th>
    );
  };

  return (
    <>
      <div className="users-table-toolbar">
        <button
          type="button"
          className="users-action-icon users-action-neutral"
          onClick={() => searchInputRef?.current?.focus()}
          title="Фокус на поиск"
          aria-label="Фокус на поиск"
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M4 4h16v2.172a2 2 0 0 1 -.586 1.414l-4.414 4.414v7l-6 2v-8.5l-4.48 -4.928a2 2 0 0 1 -.52 -1.345v-2.227"></path>
          </svg>
        </button>

        <div className="users-columns-menu-anchor">
          <button
            ref={columnsMenuButtonRef}
            type="button"
            className={`users-action-icon users-action-neutral ${isColumnsMenuOpen ? 'users-action-active' : ''}`}
            onClick={() => setIsColumnsMenuOpen((prev) => !prev)}
            title="Показать/скрыть колонки"
            aria-label="Показать/скрыть колонки"
            aria-expanded={isColumnsMenuOpen}
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M4 6l5.5 0"></path>
              <path d="M4 10l5.5 0"></path>
              <path d="M4 14l5.5 0"></path>
              <path d="M4 18l5.5 0"></path>
              <path d="M14.5 6l5.5 0"></path>
              <path d="M14.5 10l5.5 0"></path>
              <path d="M14.5 14l5.5 0"></path>
              <path d="M14.5 18l5.5 0"></path>
            </svg>
          </button>

          {isColumnsMenuOpen && (
            <div className="users-columns-menu" ref={columnsMenuRef} role="menu" aria-label="Управление колонками">
              <div className="users-columns-menu-controls">
                <button
                  type="button"
                  className="users-columns-menu-control"
                  onClick={hideAllColumns}
                  disabled={!hasVisibleToggleableColumns}
                >
                  Hide all
                </button>
                <button
                  type="button"
                  className="users-columns-menu-control"
                  onClick={unpinAllColumns}
                  disabled={!hasPinnedColumns}
                >
                  Unpin all
                </button>
                <button
                  type="button"
                  className="users-columns-menu-control"
                  onClick={showAllColumns}
                  disabled={allToggleableVisible}
                >
                  Show all
                </button>
              </div>
              <div className="users-columns-menu-divider"></div>
              <div className="users-columns-menu-list">
                {columns.map((column) => {
                  const pin = column.disablePinning ? 'left' : columnPinning[column.key];

                  return (
                    <div key={column.key} className="users-columns-menu-item" role="menuitemcheckbox" aria-checked={columnVisibility[column.key]}>
                      <div className="users-columns-menu-pins">
                        {column.disablePinning ? (
                          <span className="users-columns-menu-pin-lock" title="Колонка закреплена">🔒</span>
                        ) : (
                          <>
                            <button
                              type="button"
                              className={`users-columns-menu-pin-btn ${pin === 'left' ? 'is-active' : ''}`}
                              onClick={() => setColumnPin(column.key, 'left')}
                              title="Закрепить слева"
                              aria-label={`Закрепить слева: ${column.label}`}
                            >
                              L
                            </button>
                            <button
                              type="button"
                              className={`users-columns-menu-pin-btn ${pin === 'right' ? 'is-active' : ''}`}
                              onClick={() => setColumnPin(column.key, 'right')}
                              title="Закрепить справа"
                              aria-label={`Закрепить справа: ${column.label}`}
                            >
                              R
                            </button>
                            <button
                              type="button"
                              className={`users-columns-menu-pin-btn ${!pin ? 'is-active' : ''}`}
                              onClick={() => clearColumnPin(column.key)}
                              title="Снять закрепление"
                              aria-label={`Снять закрепление: ${column.label}`}
                            >
                              ×
                            </button>
                          </>
                        )}
                      </div>
                      <label className={`users-columns-menu-toggle ${column.alwaysVisible ? 'is-locked' : ''}`}>
                        <input
                          type="checkbox"
                          checked={columnVisibility[column.key]}
                          disabled={column.alwaysVisible}
                          onChange={(event) => toggleColumnVisibility(column.key, event.target.checked)}
                        />
                        <span>{column.label}</span>
                      </label>
                    </div>
                  );
                })}
              </div>
            </div>
          )}
        </div>

        <button
          type="button"
          className={`users-action-icon users-action-neutral ${compactRows ? 'users-action-active' : ''}`}
          onClick={() => setCompactRows((prev) => !prev)}
          title={compactRows ? 'Обычная плотность' : 'Компактный режим'}
          aria-label={compactRows ? 'Обычная плотность' : 'Компактный режим'}
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M4 3h16"></path>
            <path d="M4 9h16"></path>
            <path d="M4 15h16"></path>
            <path d="M4 21h16"></path>
          </svg>
        </button>
        <button
          type="button"
          className={`users-action-icon users-action-neutral ${isFullscreen ? 'users-action-active' : ''}`}
          onClick={() => setIsFullscreen((prev) => !prev)}
          title={isFullscreen ? 'Выйти из расширенного режима' : 'Расширенный режим'}
          aria-label={isFullscreen ? 'Выйти из расширенного режима' : 'Расширенный режим'}
        >
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M4 8v-2a2 2 0 0 1 2 -2h2"></path>
            <path d="M4 16v2a2 2 0 0 0 2 2h2"></path>
            <path d="M16 4h2a2 2 0 0 1 2 2v2"></path>
            <path d="M16 20h2a2 2 0 0 0 2 -2v-2"></path>
          </svg>
        </button>
      </div>

      <div className="card-body users-list-body">
        <div ref={tableScrollRef} className="users-table-scroll" role="region" aria-label={tableAriaLabel}>
          {resizeGuideLeft !== null ? (
            <div className="users-col-resize-guide" style={{ left: `${resizeGuideLeft}px` }} aria-hidden="true"></div>
          ) : null}
          <table
            className={`table users-rich-table ${compactRows ? 'users-table-compact' : ''}`}
            style={{ minWidth: `${tableMinWidth}px`, width: '100%', '--users-row-height': rowHeightRem }}
          >
            <colgroup>
              {renderedColumns.map((column) => (
                <col
                  key={column.key}
                  style={{ width: `${columnWidths[column.key] ?? columnMaps.defaultWidths[column.key] ?? 120}px` }}
                />
              ))}
              <col className="users-col-spacer" />
            </colgroup>
            <thead>
              <tr>
                {renderedColumns.map((column) => renderHeaderCell(column))}
                <th className="users-spacer-head" aria-hidden="true"></th>
              </tr>
            </thead>
            <tbody>
              {rows.length === 0 ? (
                <tr>
                  <td colSpan={renderedColumns.length + 1}>
                    {emptyState}
                  </td>
                </tr>
              ) : (
                paginatedRows.map((row) => {
                  const rowId = getRowId(row);
                  return (
                    <tr
                      key={rowId}
                      className={onRowClick ? 'users-row-clickable' : ''}
                      onClick={onRowClick ? () => onRowClick(row) : undefined}
                      onKeyDown={onRowClick
                        ? (event) => {
                          if (event.key === 'Enter' || event.key === ' ') {
                            event.preventDefault();
                            onRowClick(row);
                          }
                        }
                        : undefined}
                      role={onRowClick ? 'button' : undefined}
                      tabIndex={onRowClick ? 0 : undefined}
                      aria-label={onRowClick && getRowAriaLabel ? getRowAriaLabel(row) : undefined}
                    >
                      {renderedColumns.map((column) => {
                        if (column.key === 'select' && setSelectedRowIds) {
                          return (
                            <td key={column.key} className="users-cell-center users-cell-checkbox" onClick={(event) => event.stopPropagation()}>
                              <div className="users-checkbox-wrap">
                                <input
                                  type="checkbox"
                                  checked={selectedRowIds.has(rowId)}
                                  onChange={() => toggleRowSelection(rowId)}
                                  onClick={(event) => event.stopPropagation()}
                                  aria-label="Выбрать строку"
                                />
                              </div>
                            </td>
                          );
                        }

                        return renderCell(column, row);
                      })}
                      <td className="users-spacer-cell" aria-hidden="true"></td>
                    </tr>
                  );
                })
              )}
            </tbody>
          </table>
        </div>
      </div>

      <div className="users-table-bottom-toolbar">
        <div className="users-table-bottom-container">
          <span></span>
          <div className="users-table-pagination-container">
            <div className="users-table-pagination">
              <div className="users-table-pagination-group">
                <p className="users-table-pagination-label" id={pageSizeLabelId}>{pageSizeLabel}</p>
                <div className="users-table-pagination-select-wrap">
                  <AppSelect
                    className="users-table-pagination-select"
                    aria-labelledby={pageSizeLabelId}
                    value={pageSize}
                    onChange={(event) => {
                      setPageSize(parseInt(event.target.value, 10));
                      setPageIndex(0);
                    }}
                  >
                    {pageSizeOptions.map((option) => (
                      <option key={option} value={option}>{option}</option>
                    ))}
                  </AppSelect>
                </div>
                <p className="users-table-pagination-range">{pageRangeStart}-{pageRangeEnd} of {totalRows}</p>
                <div className="users-table-pagination-nav">
                  <button
                    type="button"
                    className="users-pagination-nav-btn"
                    aria-label="Go to previous page"
                    disabled={normalizedPageIndex <= 0}
                    onClick={() => setPageIndex((prev) => Math.max(0, prev - 1))}
                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M15 6l-6 6l6 6"></path>
                    </svg>
                  </button>
                  <button
                    type="button"
                    className="users-pagination-nav-btn"
                    aria-label="Go to next page"
                    disabled={normalizedPageIndex >= totalPages - 1}
                    onClick={() => setPageIndex((prev) => Math.min(totalPages - 1, prev + 1))}
                  >
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M9 6l6 6l-6 6"></path>
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

export default SharedRichTable;
