import { useEffect, useMemo, useRef, useState } from 'react';

const normalizeValue = (value) => {
  if (value === null || value === undefined) {
    return '';
  }
  return String(value);
};

function AppMultiSelect({
  value = [],
  onChange,
  options = [],
  className = '',
  placeholder = 'Выберите значения',
  description = '',
  disabled = false,
  leftSection = null,
  noResultsLabel = 'Ничего не найдено',
  searchable = true,
  id,
  ariaDescribedBy,
}) {
  const [open, setOpen] = useState(false);
  const [search, setSearch] = useState('');
  const rootRef = useRef(null);
  const inputRef = useRef(null);
  const selectedValues = useMemo(() => value.map((item) => normalizeValue(item)), [value]);
  const selectedSet = useMemo(() => new Set(selectedValues), [selectedValues]);

  const normalizedOptions = useMemo(
    () =>
      options.map((option) => ({
        ...option,
        value: normalizeValue(option.value),
        label: option.label ?? option.value,
        searchText: String(option.searchText ?? option.label ?? option.value).toLowerCase(),
      })),
    [options]
  );

  const selectedOptions = useMemo(
    () => normalizedOptions.filter((option) => selectedSet.has(option.value)),
    [normalizedOptions, selectedSet]
  );

  const filteredOptions = useMemo(() => {
    const needle = search.trim().toLowerCase();
    if (!needle) {
      return normalizedOptions;
    }
    return normalizedOptions.filter((option) => option.searchText.includes(needle));
  }, [normalizedOptions, search]);

  useEffect(() => {
    if (!open) {
      return undefined;
    }

    const handlePointerDown = (event) => {
      if (!rootRef.current?.contains(event.target)) {
        setOpen(false);
      }
    };

    const handleEscape = (event) => {
      if (event.key === 'Escape') {
        setOpen(false);
        inputRef.current?.blur();
      }
    };

    document.addEventListener('mousedown', handlePointerDown);
    document.addEventListener('keydown', handleEscape);
    return () => {
      document.removeEventListener('mousedown', handlePointerDown);
      document.removeEventListener('keydown', handleEscape);
    };
  }, [open]);

  const emitChange = (nextValues) => {
    onChange?.(nextValues);
  };

  const toggleValue = (optionValue) => {
    if (disabled) {
      return;
    }
    const exists = selectedSet.has(optionValue);
    const nextValues = exists
      ? selectedValues.filter((valueItem) => valueItem !== optionValue)
      : [...selectedValues, optionValue];
    emitChange(nextValues);
  };

  const removeValue = (event, optionValue) => {
    event.stopPropagation();
    toggleValue(optionValue);
  };

  const handleShellClick = () => {
    if (disabled) {
      return;
    }
    setOpen(true);
    inputRef.current?.focus();
  };

  return (
    <div ref={rootRef} className={`app-multiselect ${className} ${open ? 'open' : ''} ${disabled ? 'disabled' : ''}`.trim()}>
      <div className="app-multiselect-control" onClick={handleShellClick} aria-disabled={disabled}>
        {leftSection ? <div className="app-multiselect-section left">{leftSection}</div> : null}
        <div className="app-multiselect-input-shell">
          <div className="app-multiselect-pills">
            {selectedOptions.map((option) => (
              <span key={option.value} className="app-multiselect-pill">
                <span className="app-multiselect-pill-label">{option.label}</span>
                <button
                  type="button"
                  className="app-multiselect-pill-remove"
                  onClick={(event) => removeValue(event, option.value)}
                  aria-label={`Удалить ${option.label}`}
                  disabled={disabled}
                >
                  <svg viewBox="0 0 15 15" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path
                      d="M11.7816 4.03157C12.0062 3.80702 12.0062 3.44295 11.7816 3.2184C11.5571 2.99385 11.193 2.99385 10.9685 3.2184L7.50005 6.68682L4.03164 3.2184C3.80708 2.99385 3.44301 2.99385 3.21846 3.2184C2.99391 3.44295 2.99391 3.80702 3.21846 4.03157L6.68688 7.49999L3.21846 10.9684C2.99391 11.193 2.99391 11.557 3.21846 11.7816C3.44301 12.0061 3.80708 12.0061 4.03164 11.7816L7.50005 8.31316L10.9685 11.7816C11.193 12.0061 11.5571 12.0061 11.7816 11.7816C12.0062 11.557 12.0062 11.193 11.7816 10.9684L8.31322 7.49999L11.7816 4.03157Z"
                      fill="currentColor"
                      fillRule="evenodd"
                      clipRule="evenodd"
                    />
                  </svg>
                </button>
              </span>
            ))}
            <input
              ref={inputRef}
              id={id}
              type="text"
              className="app-multiselect-input"
              value={search}
              onChange={(event) => setSearch(event.target.value)}
              onFocus={() => setOpen(true)}
              placeholder={selectedOptions.length === 0 ? placeholder : ''}
              disabled={disabled || !searchable}
              aria-describedby={ariaDescribedBy}
            />
          </div>
        </div>
        <div className="app-multiselect-section right" aria-hidden="true">
          <svg className={`app-select-chevron ${open ? 'open' : ''}`.trim()} viewBox="0 0 15 15" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M4.93179 5.43179C4.75605 5.60753 4.75605 5.89245 4.93179 6.06819C5.10753 6.24392 5.39245 6.24392 5.56819 6.06819L7.49999 4.13638L9.43179 6.06819C9.60753 6.24392 9.89245 6.24392 10.0682 6.06819C10.2439 5.89245 10.2439 5.60753 10.0682 5.43179L7.81819 3.18179C7.73379 3.0974 7.61933 3.04999 7.49999 3.04999C7.38064 3.04999 7.26618 3.0974 7.18179 3.18179L4.93179 5.43179ZM10.0682 9.56819C10.2439 9.39245 10.2439 9.10753 10.0682 8.93179C9.89245 8.75606 9.60753 8.75606 9.43179 8.93179L7.49999 10.8636L5.56819 8.93179C5.39245 8.75606 5.10753 8.75606 4.93179 8.93179C4.75605 9.10753 4.75605 9.39245 4.93179 9.56819L7.18179 11.8182C7.35753 11.9939 7.64245 11.9939 7.81819 11.8182L10.0682 9.56819Z"
              fill="currentColor"
              fillRule="evenodd"
              clipRule="evenodd"
            />
          </svg>
        </div>
      </div>

      {description ? (
        <p id={ariaDescribedBy} className="app-multiselect-description">
          {description}
        </p>
      ) : null}

      {open ? (
        <div className="app-select-dropdown app-multiselect-dropdown" role="presentation">
          <div className="app-select-options" role="listbox">
            {filteredOptions.length === 0 ? (
              <div className="app-select-empty">{noResultsLabel}</div>
            ) : (
              filteredOptions.map((option) => {
                const isSelected = selectedSet.has(option.value);
                return (
                  <button
                    key={option.value}
                    type="button"
                    className={`app-select-option ${isSelected ? 'selected' : ''}`.trim()}
                    onClick={() => toggleValue(option.value)}
                    role="option"
                    aria-selected={isSelected}
                    disabled={disabled || option.disabled}
                  >
                    <span>{option.label}</span>
                    {isSelected ? (
                      <svg viewBox="0 0 10 7" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                        <path
                          d="M4 4.586L1.707 2.293A1 1 0 1 0 .293 3.707l3 3a.997.997 0 0 0 1.414 0l5-5A1 1 0 1 0 8.293.293L4 4.586z"
                          fill="currentColor"
                          fillRule="evenodd"
                          clipRule="evenodd"
                        />
                      </svg>
                    ) : null}
                  </button>
                );
              })
            )}
          </div>
        </div>
      ) : null}
    </div>
  );
}

export default AppMultiSelect;
