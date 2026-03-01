import React, { useState, useEffect } from 'react';
import { highlight, languages } from 'prismjs';
import 'prismjs/components/prism-json';
import Editor from 'react-simple-code-editor';
import { configProfilesApi } from '../api';
import { validateWithWasm, formatWithWasm } from '../utils/wasmValidator';
import TilesToolbar from '../components/TilesToolbar';
import useAdaptiveEntityGridColumns from '../components/useAdaptiveEntityGridColumns';
import ConfirmActionModal from '../components/ConfirmActionModal';

function ConfigProfiles() {
  const { gridRef, columns: gridColumns } = useAdaptiveEntityGridColumns();
  const [profiles, setProfiles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [editingProfile, setEditingProfile] = useState(null);
  const [showRenameModal, setShowRenameModal] = useState(false);
  const [renameProfile, setRenameProfile] = useState(null);
  const [renameName, setRenameName] = useState('');
  const [configJson, setConfigJson] = useState('');
  const [jsonError, setJsonError] = useState(null);
  const [openDropdown, setOpenDropdown] = useState(null);
  const [dragIndex, setDragIndex] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [deleteTarget, setDeleteTarget] = useState(null);
  const [deleting, setDeleting] = useState(false);

  // Load profiles
  const loadProfiles = async () => {
    try {
      setLoading(true);
      const data = await configProfilesApi.getAll();
      setProfiles(data.profiles || []);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProfiles();
  }, []);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = () => setOpenDropdown(null);
    if (openDropdown) {
      document.addEventListener('click', handleClickOutside);
    }
    return () => document.removeEventListener('click', handleClickOutside);
  }, [openDropdown]);

  // Open modal for create
  const handleCreate = () => {
    setEditingProfile(null);
    setConfigJson('{}');
    setJsonError(null);
    setShowModal(true);
  };

  // Open modal for edit
  const handleEdit = (profile) => {
    setEditingProfile(profile);
    setConfigJson(JSON.stringify(profile.config, null, 2));
    setJsonError(null);
    setShowModal(true);
  };

  // Open rename modal
  const handleRename = (profile) => {
    setRenameProfile(profile);
    setRenameName(profile.name);
    setShowRenameModal(true);
  };

  // Handle delete
  const handleDelete = (uuid, name) => {
    setOpenDropdown(null);
    setDeleteTarget({ uuid, name });
  };

  const handleConfirmDelete = async () => {
    if (!deleteTarget || deleting) {
      return;
    }

    try {
      setDeleting(true);
      await configProfilesApi.delete(deleteTarget.uuid);
      setDeleteTarget(null);
      await loadProfiles();
    } catch (err) {
      alert(`Failed to delete: ${err.message}`);
    } finally {
      setDeleting(false);
    }
  };

  // Download config as JSON file
  const downloadConfig = (profile) => {
    const jsonString = JSON.stringify(profile.config, null, 2);
    const blob = new Blob([jsonString], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `${profile.name || 'config'}.json`;
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    URL.revokeObjectURL(url);
  };

  // Copy UUID to clipboard
  const copyUuid = async (uuid) => {
    try {
      await navigator.clipboard.writeText(uuid);
      alert('UUID copied to clipboard!');
    } catch (err) {
      console.error('Failed to copy UUID:', err);
    }
  };

  // Handle rename submit
  const handleRenameSubmit = async (e) => {
    e.preventDefault();
    try {
      await configProfilesApi.update(renameProfile.uuid, { name: renameName });
      setShowRenameModal(false);
      await loadProfiles();
    } catch (err) {
      alert(`Failed to rename: ${err.message}`);
    }
  };

  // Validate and parse JSON
  const validateJson = (jsonString) => {
    try {
      return JSON.parse(jsonString);
    } catch (err) {
      setJsonError(`Invalid JSON: ${err.message}`);
      return null;
    }
  };

  // Handle JSON editor change
  const handleJsonChange = (value) => {
    setConfigJson(value);
    const parsed = validateJson(value);
    if (parsed) {
      setJsonError(null);
    }
  };

  const formatConfigJson = async () => {
    try {
      const wasmFormatted = await formatWithWasm('xray', configJson || '{}');
      setConfigJson(wasmFormatted);
      setJsonError(null);
      return;
    } catch (err) {
      console.warn('[wasm:xray] format fallback', err);
    }

    try {
      const normalized = JSON.stringify(JSON.parse(configJson || '{}'), null, 2);
      setConfigJson(normalized);
      setJsonError(null);
    } catch (err) {
      setJsonError(`Invalid JSON: ${err.message}`);
      alert(`Invalid JSON: ${err.message}`);
    }
  };

  const resetConfigEditor = () => {
    if (editingProfile) {
      setConfigJson(JSON.stringify(editingProfile.config || {}, null, 2));
      setJsonError(null);
      return;
    }

    setConfigJson('{}');
    setJsonError(null);
  };

  const validateConfigEditor = async () => {
    try {
      const result = await validateWithWasm('xray', configJson || '{}');
      if (result.ok) {
        if (result.formatted) {
          setConfigJson(result.formatted);
        }
        setJsonError(null);
        alert(result.message || 'JSON is valid (WASM).');
      } else {
        const message = result.message || 'Config validation failed.';
        setJsonError(message);
        alert(message);
      }
      return;
    } catch (err) {
      console.warn('[wasm:xray] validate fallback', err);
    }

    try {
      JSON.parse(configJson || '{}');
      setJsonError(null);
      alert('JSON is valid.');
    } catch (err) {
      setJsonError(`Invalid JSON: ${err.message}`);
      alert(`Invalid JSON: ${err.message}`);
    }
  };

  const highlightJsonWithIndentGuides = (code) => {
    try {
      const highlighted = highlight(code, languages.json, 'json');
      return highlighted
        .split('\n')
        .map((line) =>
          line.replace(/^([ \t]+)/, (indent) => {
            const normalizedIndent = indent.replace(/\t/g, '  ');
            const levels = Math.floor(normalizedIndent.length / 2);
            let guides = '';
            for (let i = 0; i < levels; i += 1) {
              const colorClass = jsonError ? 'indent-guide-error' : `indent-guide-level-${i % 6}`;
              guides += `<span class="indent-guide ${colorClass}">|</span><span class="indent-guide-gap"> </span>`;
            }
            return `<span class="indent-guides">${guides}</span>`;
          })
        )
        .join('\n');
    } catch (err) {
      console.error('Syntax highlighting error:', err);
      return code;
    }
  };

  // Handle form submit
  const handleSubmit = async (e) => {
    e.preventDefault();

    const parsedConfig = validateJson(configJson);
    if (!parsedConfig) {
      return;
    }

    const submitData = {
      name: editingProfile ? editingProfile.name : 'Unnamed Profile',
      config: parsedConfig,
    };

    try {
      if (editingProfile) {
        await configProfilesApi.update(editingProfile.uuid, submitData);
      } else {
        await configProfilesApi.create(submitData);
      }
      setShowModal(false);
      await loadProfiles();
    } catch (err) {
      alert(`Failed to save: ${err.message}`);
    }
  };

  // Toggle dropdown
  const toggleDropdown = (e, profileUuid) => {
    e.stopPropagation();
    setOpenDropdown(openDropdown === profileUuid ? null : profileUuid);
  };

  const prepareDragPreview = (event) => {
    const clone = event.currentTarget.cloneNode(true);
    const { width, height } = event.currentTarget.getBoundingClientRect();
    const toRem = (value) => `${value / 16}rem`;
    const previewHeight = Math.min(height, 352);
    clone.style.position = 'fixed';
    clone.style.top = '-62.5rem';
    clone.style.left = '-62.5rem';
    clone.style.width = toRem(width);
    clone.style.height = toRem(previewHeight);
    clone.style.maxHeight = toRem(previewHeight);
    clone.style.overflow = 'hidden';
    clone.style.background = 'var(--card-background)';
    clone.style.boxShadow = '0 0.875rem 2.25rem rgba(0, 0, 0, 0.28)';
    clone.style.border = '0.0625rem solid var(--border-color)';
    clone.style.opacity = '1';
    clone.style.pointerEvents = 'none';
    document.body.appendChild(clone);
    event.dataTransfer.setDragImage(clone, 20, 20);
    window.setTimeout(() => {
      if (clone.parentNode) {
        clone.parentNode.removeChild(clone);
      }
    }, 120);
  };

  const handleDragStart = (event, index) => {
    if (!canReorder) return;
    setDragIndex(index);
    event.dataTransfer.effectAllowed = 'move';
    prepareDragPreview(event);
  };

  const handleDrop = async (dropIndex) => {
    if (!canReorder || dragIndex === null || dragIndex === dropIndex) {
      setDragIndex(null);
      return;
    }

    const reordered = [...profiles];
    const [moved] = reordered.splice(dragIndex, 1);
    reordered.splice(dropIndex, 0, moved);

    setProfiles(reordered);
    setDragIndex(null);

    try {
      await configProfilesApi.reorder(reordered.map((profile) => profile.uuid));
      setProfiles((prev) => prev.map((profile, index) => ({ ...profile, view_position: index })));
    } catch (err) {
      alert(`Failed to reorder profiles: ${err.message}`);
      await loadProfiles();
    }
  };

  // Count inbounds in config
  const countInbounds = (config) => {
    if (!config || !config.inbounds) return 0;
    return Array.isArray(config.inbounds) ? config.inbounds.length : 0;
  };

  const normalizedSearch = searchQuery.trim().toLowerCase();
  const canReorder = normalizedSearch === '';
  const filteredProfiles = profiles.filter((profile) => {
    if (!normalizedSearch) return true;
    const name = profile.name?.toLowerCase() || '';
    const uuid = profile.uuid?.toLowerCase() || '';
    return name.includes(normalizedSearch) || uuid.includes(normalizedSearch);
  });

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="page">
      {!showModal ? (
        <>
          <TilesToolbar
            title="Профили"
            searchValue={searchQuery}
            onSearchChange={setSearchQuery}
            searchPlaceholder="Поиск профилей..."
            onReload={loadProfiles}
            onCreate={handleCreate}
            createTitle="Добавить профиль"
          />

          {error && (
            <div className="alert alert-error">
              {error}
              <button type="button" className="btn-icon" onClick={() => loadProfiles()}>↻</button>
            </div>
          )}

          {profiles.length === 0 ? (
            <div className="empty-state">
              <p>No config profiles found</p>
              <button className="btn btn-primary" onClick={handleCreate}>
                Create your first profile
              </button>
            </div>
          ) : filteredProfiles.length === 0 ? (
            <div className="empty-state">
              <p>По запросу ничего не найдено</p>
              <button type="button" className="btn btn-secondary" onClick={() => setSearchQuery('')}>
                Сбросить поиск
              </button>
            </div>
          ) : (
            <div
              ref={gridRef}
              className="templates-windows-grid entity-windows-grid"
              style={{
                '--entity-grid-columns': gridColumns,
              }}
            >
              {filteredProfiles.map((profile, index) => (
                <div key={profile.uuid} className="templates-grid-item">
                  <div
                    className={`templates-item-wrapper card-tile templates-window ${dragIndex === index ? 'dragging' : ''} ${openDropdown === profile.uuid ? 'menu-open' : ''}`}
                    draggable={canReorder}
                    onDragStart={(e) => handleDragStart(e, index)}
                    onDragOver={(e) => canReorder && e.preventDefault()}
                    onDrop={() => handleDrop(index)}
                    onDragEnd={() => setDragIndex(null)}
                  >
                    <button
                      type="button"
                      className={`templates-drag-handle ${canReorder ? '' : 'disabled'}`}
                      aria-label="Reorder item"
                      title={canReorder ? 'Перетащите для сортировки' : 'Сортировка отключена при поиске'}
                    >
                      <svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                        <path d="M8.5 7A1.5 1.5 0 1 0 8.5 4a1.5 1.5 0 0 0 0 3zm0 6.5a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zm0 6.5a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zM15.5 7a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zm0 6.5a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zm0 6.5a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3z" />
                      </svg>
                    </button>
                    <div className="templates-top-accent"></div>
                    <div className="templates-glow-effect"></div>
                    <div className="templates-window-main">
                      <div className="templates-window-header">
                        <div className="templates-icon-wrapper">
                          <button type="button" className="templates-icon-button" tabIndex={-1} aria-hidden="true">
                            <svg fill="currentColor" preserveAspectRatio="xMidYMid meet" viewBox="0 0 35 35" xmlns="http://www.w3.org/2000/svg">
                              <path d="M16.6961 15.2606L16.5825 3.49701C16.5718 2.38439 15.025 2.11843 14.6433 3.16356L11.7279 11.1447C11.6384 11.3898 11.4566 11.5902 11.2213 11.7031L5.66765 14.3687C4.70841 14.8291 5.03635 16.2703 6.10036 16.2703H15.6962C16.2522 16.2703 16.7015 15.8166 16.6961 15.2606Z" />
                              <path d="M18.6471 15.2703V5.88936C18.6471 4.84679 20.0428 4.49998 20.5308 5.4213L23.5833 11.1845C23.7 11.4049 23.8948 11.5737 24.1296 11.6578L31.5829 14.3289C32.6388 14.7073 32.3671 16.2703 31.2455 16.2703H19.6471C19.0948 16.2703 18.6471 15.8226 18.6471 15.2703Z" />
                              <path d="M18.6471 31.4643V19.3784C18.6471 18.8261 19.0948 18.3784 19.6471 18.3784H29.2853C30.3376 18.3784 30.676 19.7947 29.7374 20.2704L24.1129 23.1208C23.889 23.2343 23.716 23.4278 23.6281 23.663L20.5839 31.8141C20.1941 32.8578 18.6471 32.5783 18.6471 31.4643Z" />
                              <path d="M16.7059 28.9873V19.3784C16.7059 18.8261 16.2582 18.3784 15.7059 18.3784H3.83963C2.71522 18.3784 2.44656 19.9473 3.50691 20.3214L11.5457 23.1578C11.7987 23.247 12.0052 23.4342 12.1188 23.6772L14.8 29.4109C15.2531 30.3797 16.7059 30.0568 16.7059 28.9873Z" />
                            </svg>
                          </button>
                        </div>
                        <div className="templates-title-wrap">
                          <h3 className="card-tile-title" title={profile.name}>{profile.name}</h3>
                          <p className="templates-subtitle">Config Profile</p>
                        </div>
                      </div>

                      <div className="templates-window-stats">
                        <div className="templates-window-stat">
                          <span className="templates-window-stat-value">{countInbounds(profile.config)}</span>
                          <span className="templates-window-stat-label">Inbounds</span>
                        </div>
                      </div>

                      <div className="templates-window-body">
                        <div className="card-tile-actions templates-actions">
                          <button className="btn btn-primary btn-block templates-edit-btn" onClick={() => handleEdit(profile)}>
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                              <path d="M7 7h-1a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-1"></path>
                              <path d="M20.385 6.585a2.1 2.1 0 0 0 -2.97 -2.97l-8.415 8.385v3h3l8.385 -8.415z"></path>
                              <path d="M16 5l3 3"></path>
                            </svg>
                            <span>Изменить</span>
                          </button>
                          <div className="templates-menu-wrap" onClick={(e) => e.stopPropagation()}>
                            <button
                              className={`templates-menu-control ${openDropdown === profile.uuid ? 'open' : ''}`}
                              onClick={(e) => toggleDropdown(e, profile.uuid)}
                              title="More actions"
                              aria-label="More actions"
                            >
                              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                                <path d="M6 9l6 6l6 -6"></path>
                              </svg>
                            </button>
                            {openDropdown === profile.uuid && (
                              <div className="templates-dropdown-menu">
                                <button className="templates-dropdown-item" onClick={() => handleRename(profile)}>
                                  Rename
                                </button>
                                <button className="templates-dropdown-item" onClick={() => downloadConfig(profile)}>
                                  Download Config
                                </button>
                                <button className="templates-dropdown-item" onClick={() => copyUuid(profile.uuid)}>
                                  Copy UUID
                                </button>
                                <button className="templates-dropdown-item danger" onClick={() => handleDelete(profile.uuid, profile.name)}>
                                  Delete Profile
                                </button>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </>
      ) : (
        <div className="right-editor-inline" aria-label="Config editor panel">
          <div className="right-editor-modal right-editor-inline-window">
            <div className="modal-body-fullscreen right-editor-layout">
              <div className="right-editor-card right-editor-head">
                <div className="right-editor-head-left">
                  <button type="button" className="editor-round-icon editor-round-icon-cyan" title="Profile">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M14 2H6a2 2 0 0 0 -2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2 -2V8z" />
                      <polyline points="14 2 14 8 20 8" />
                      <line x1="16" y1="13" x2="8" y2="13" />
                      <line x1="16" y1="17" x2="8" y2="17" />
                    </svg>
                  </button>
                  <div className="right-editor-title-stack">
                    <h4 className="right-editor-title">{editingProfile?.name || 'Default'}</h4>
                    <span className="right-editor-subtitle">Config Profile</span>
                  </div>
                </div>
                <div className="right-editor-head-actions">
                  <button type="button" className="editor-round-icon editor-round-icon-lime" title="Validate" onClick={validateConfigEditor}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M12 17h.01" />
                      <path d="M12 3a9 9 0 0 1 9 9a9 9 0 1 1 -9 -9z" />
                      <path d="M9.09 9a3 3 0 1 1 5.82 1c0 2-3 3-3 3" />
                    </svg>
                  </button>
                  <button type="button" className="editor-round-icon editor-round-icon-gray" title="Reset" onClick={resetConfigEditor}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M3 10h11a4 4 0 1 1 0 8h-1" />
                      <path d="M6 14l-4-4l4-4" />
                    </svg>
                  </button>
                  <button type="button" className="editor-round-icon editor-round-icon-default" title="Свернуть" onClick={() => setShowModal(false)}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <line x1="5" y1="12" x2="19" y2="12" />
                    </svg>
                  </button>
                </div>
              </div>

              <div className="right-editor-card right-editor-main">
                <form id="config-profile-editor-form" className="config-editor-form right-editor-form" onSubmit={handleSubmit}>
                  <div className="form-group config-editor-group">
                    <label className="form-label">Config JSON *</label>
                    <div className={`code-editor-container ${jsonError ? 'error' : ''}`}>
                      <Editor
                        value={configJson}
                        onValueChange={handleJsonChange}
                        highlight={highlightJsonWithIndentGuides}
                        padding={12}
                        style={{
                          fontFamily: '"Monaco", "Consolas", monospace',
                          fontSize: '0.875rem',
                          lineHeight: '1.5',
                          height: '100%',
                          backgroundColor: 'var(--editor-background)',
                          width: '100%',
                          boxSizing: 'border-box'
                        }}
                        textareaClassName="code-editor-textarea"
                        preClassName="code-editor-pre"
                        placeholder='{"log": {"level": "info"}, ...}'
                      />
                    </div>
                    {jsonError && <div className="form-error config-editor-error">{jsonError}</div>}
                  </div>
                </form>
              </div>

              <div className="right-editor-card right-editor-footer">
                <button
                  type="submit"
                  className="btn btn-primary"
                  form="config-profile-editor-form"
                  disabled={!!jsonError}
                >
                  {editingProfile ? 'Сохранить' : 'Создать'}
                </button>
                <div className="right-editor-footer-actions">
                  <button type="button" className="editor-round-icon editor-round-icon-default" title="Actions">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <line x1="4" y1="6" x2="20" y2="6" />
                      <line x1="7" y1="12" x2="20" y2="12" />
                      <line x1="10" y1="18" x2="20" y2="18" />
                    </svg>
                  </button>
                  <button type="button" className="btn btn-secondary" onClick={formatConfigJson}>
                    Форматировать
                  </button>
                  <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>
                    Закрыть
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Rename Modal */}
      {showRenameModal && renameProfile && (
        <div className="modal-overlay" onClick={() => setShowRenameModal(false)}>
          <div className="modal" onClick={(e) => e.stopPropagation()}>
            <div className="modal-header">
              <h2>Rename Profile</h2>
              <button className="btn-icon" onClick={() => setShowRenameModal(false)}>✕</button>
            </div>
            <form onSubmit={handleRenameSubmit}>
              <div className="form-group">
                <label>Name *</label>
                <input
                  type="text"
                  value={renameName}
                  onChange={(e) => setRenameName(e.target.value)}
                  required
                  placeholder="Profile name"
                  autoFocus
                />
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowRenameModal(false)}>
                  Cancel
                </button>
                <button type="submit" className="btn btn-primary">
                  Save Changes
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      <ConfirmActionModal
        open={Boolean(deleteTarget)}
        title="Удалить"
        message={`Вы уверены, что хотите удалить профиль "${deleteTarget?.name || ''}"? Это действие нельзя отменить.`}
        confirmLabel="Удалить"
        cancelLabel="Закрыть"
        loading={deleting}
        onClose={() => {
          if (!deleting) {
            setDeleteTarget(null);
          }
        }}
        onConfirm={handleConfirmDelete}
      />
    </div>
  );
}

export default ConfigProfiles;
