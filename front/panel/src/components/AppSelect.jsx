import { Children, isValidElement, useEffect, useMemo, useRef, useState } from 'react';

const normalizeValue = (value) => {
  if (value === null || value === undefined) {
    return '';
  }
  return String(value);
};

function AppSelect({
  value,
  onChange,
  children,
  className = '',
  disabled = false,
  placeholder = '',
  id,
  name,
  'aria-labelledby': ariaLabelledBy,
  'aria-label': ariaLabel,
}) {
  const [open, setOpen] = useState(false);
  const rootRef = useRef(null);
  const triggerRef = useRef(null);

  const options = useMemo(
    () =>
      Children.toArray(children)
        .filter((child) => isValidElement(child) && child.type === 'option')
        .map((child) => ({
          value: normalizeValue(child.props.value),
          disabled: !!child.props.disabled,
          label: child.props.children,
        })),
    [children]
  );

  const selectedValue = normalizeValue(value);
  const selectedOption = options.find((option) => option.value === selectedValue) || null;

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
        triggerRef.current?.focus();
      }
    };

    document.addEventListener('mousedown', handlePointerDown);
    document.addEventListener('keydown', handleEscape);
    return () => {
      document.removeEventListener('mousedown', handlePointerDown);
      document.removeEventListener('keydown', handleEscape);
    };
  }, [open]);

  const emitChange = (nextValue) => {
    if (!onChange) {
      return;
    }
    onChange({
      target: {
        value: nextValue,
        name,
      },
    });
  };

  const handleSelect = (optionValue) => {
    emitChange(optionValue);
    setOpen(false);
    triggerRef.current?.focus();
  };

  const displayLabel = selectedOption?.label || placeholder || '';

  return (
    <div
      ref={rootRef}
      className={`app-select ${className} ${open ? 'open' : ''} ${disabled ? 'disabled' : ''}`.trim()}
    >
      <button
        ref={triggerRef}
        id={id}
        type="button"
        className="app-select-trigger"
        onClick={() => !disabled && setOpen((prev) => !prev)}
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-labelledby={ariaLabelledBy}
        aria-label={ariaLabel}
        disabled={disabled}
      >
        <span className={`app-select-trigger-label ${selectedOption ? '' : 'placeholder'}`.trim()}>
          {displayLabel}
        </span>
        <span className={`app-select-chevron ${open ? 'open' : ''}`.trim()} aria-hidden="true">
          <svg viewBox="0 0 15 15" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path
              d="M4.93179 5.43179C4.75605 5.60753 4.75605 5.89245 4.93179 6.06819C5.10753 6.24392 5.39245 6.24392 5.56819 6.06819L7.49999 4.13638L9.43179 6.06819C9.60753 6.24392 9.89245 6.24392 10.0682 6.06819C10.2439 5.89245 10.2439 5.60753 10.0682 5.43179L7.81819 3.18179C7.73379 3.0974 7.61933 3.04999 7.49999 3.04999C7.38064 3.04999 7.26618 3.0974 7.18179 3.18179L4.93179 5.43179ZM10.0682 9.56819C10.2439 9.39245 10.2439 9.10753 10.0682 8.93179C9.89245 8.75606 9.60753 8.75606 9.43179 8.93179L7.49999 10.8636L5.56819 8.93179C5.39245 8.75606 5.10753 8.75606 4.93179 8.93179C4.75605 9.10753 4.75605 9.39245 4.93179 9.56819L7.18179 11.8182C7.35753 11.9939 7.64245 11.9939 7.81819 11.8182L10.0682 9.56819Z"
              fill="currentColor"
              fillRule="evenodd"
              clipRule="evenodd"
            />
          </svg>
        </span>
      </button>

      {open ? (
        <div className="app-select-dropdown" role="presentation">
          <div className="app-select-options" role="listbox" aria-labelledby={ariaLabelledBy}>
            {options.map((option) => {
              const isSelected = option.value === selectedValue;
              return (
                <button
                  key={`${option.value}-${String(option.label)}`}
                  type="button"
                  className={`app-select-option ${isSelected ? 'selected' : ''}`.trim()}
                  onClick={() => handleSelect(option.value)}
                  disabled={option.disabled}
                  role="option"
                  aria-selected={isSelected}
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
            })}
          </div>
        </div>
      ) : null}
    </div>
  );
}

export default AppSelect;
