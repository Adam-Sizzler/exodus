import React from 'react';

function TilesToolbar({
  title,
  searchValue,
  onSearchChange,
  searchPlaceholder = 'Поиск...',
  onReload,
  onCreate,
  reloadTitle = 'Обновить список',
  createTitle = 'Добавить',
}) {
  return (
    <div className="profiles-toolbar">
      <div className="profiles-toolbar-shell">
        <div className="profiles-toolbar-header">
          <div className="profiles-toolbar-title">
            <span className="profiles-toolbar-badge" aria-hidden="true">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M9.183 6.117a6 6 0 1 0 4.511 3.986" />
                <path d="M14.813 17.883a6 6 0 1 0 -4.496 -3.954" />
              </svg>
            </span>
            <div className="profiles-toolbar-heading">
              <h1 className="page-title">{title}</h1>
            </div>
          </div>
        </div>
        <div className="profiles-toolbar-actions">
          <div className="search-box search-box-compact profiles-search-box">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="11" cy="11" r="8" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
            <input
              type="text"
              placeholder={searchPlaceholder}
              value={searchValue}
              onChange={(event) => onSearchChange(event.target.value)}
            />
          </div>
          <button
            type="button"
            className="profiles-action-icon profiles-action-refresh"
            onClick={onReload}
            title={reloadTitle}
            aria-label={reloadTitle}
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M20 11a8.1 8.1 0 0 0 -15.5 -2m-.5 -4v4h4" />
              <path d="M4 13a8.1 8.1 0 0 0 15.5 2m.5 4v-4h-4" />
            </svg>
          </button>
          <button
            type="button"
            className="profiles-action-icon profiles-action-add"
            onClick={onCreate}
            title={createTitle}
            aria-label={createTitle}
          >
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M12 5l0 14" />
              <path d="M5 12l14 0" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  );
}

export default TilesToolbar;
