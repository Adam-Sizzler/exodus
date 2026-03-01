import React, { useEffect, useMemo, useState } from 'react';
import { squadsApi, configProfilesApi, usersApi } from '../api';
import TilesToolbar from '../components/TilesToolbar';
import useAdaptiveEntityGridColumns from '../components/useAdaptiveEntityGridColumns';
import SideDrawerPanel from '../components/SideDrawerPanel';
import ConfirmActionModal from '../components/ConfirmActionModal';

const INBOUNDS_TAB = {
  profiles: 'profiles',
  list: 'list',
};

const FLAT_FILTER = {
  all: 'all',
  selected: 'selected',
  unselected: 'unselected',
};

const getInboundPort = (inbound) => {
  return inbound?.port ?? inbound?.listen_port ?? inbound?.raw_inbound?.listen_port ?? 'N/A';
};

const getInboundProtocol = (inbound) => {
  return inbound?.type ?? inbound?.protocol ?? inbound?.raw_inbound?.protocol ?? 'unknown';
};

const getInboundSecurity = (inbound) => {
  return inbound?.security ?? inbound?.raw_inbound?.security ?? 'none';
};

function InternalSquads() {
  const { gridRef, columns: gridColumns } = useAdaptiveEntityGridColumns();
  const [squads, setSquads] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [apiNotAvailable, setApiNotAvailable] = useState(false);

  // Modal states
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showInboundsDrawer, setShowInboundsDrawer] = useState(false);
  const [showMembersModal, setShowMembersModal] = useState(false);
  const [showRenameModal, setShowRenameModal] = useState(false);

  // Selected squad for editing
  const [selectedSquad, setSelectedSquad] = useState(null);

  // Form data
  const [formData, setFormData] = useState({ name: '', view_position: 0 });
  const [renameName, setRenameName] = useState('');

  // Inbounds and members data
  const [configProfiles, setConfigProfiles] = useState([]);
  const [availableMembers, setAvailableMembers] = useState([]);
  const [selectedInbounds, setSelectedInbounds] = useState(new Set());
  const [selectedMemberIds, setSelectedMemberIds] = useState(new Set());

  // Drawer state
  const [inboundsTab, setInboundsTab] = useState(INBOUNDS_TAB.profiles);
  const [inboundsSearchQuery, setInboundsSearchQuery] = useState('');
  const [flatFilter, setFlatFilter] = useState(FLAT_FILTER.all);
  const [loadingInboundsDrawer, setLoadingInboundsDrawer] = useState(false);
  const [savingInbounds, setSavingInbounds] = useState(false);
  const [expandedProfiles, setExpandedProfiles] = useState(new Set());

  // Dropdown state
  const [openDropdown, setOpenDropdown] = useState(null);
  const [dragIndex, setDragIndex] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [deleteTarget, setDeleteTarget] = useState(null);
  const [deleting, setDeleting] = useState(false);

  // Load squads
  const loadSquads = async () => {
    try {
      setLoading(true);
      const data = await squadsApi.getAllSummary();
      setSquads(data.squads || []);
      setError(null);
      setApiNotAvailable(false);
    } catch (err) {
      console.error('Failed to load squads:', err);
      setError(err.message);
      setApiNotAvailable(true);
    } finally {
      setLoading(false);
    }
  };

  // Load available members (users)
  const loadAvailableMembers = async () => {
    try {
      const data = await usersApi.getAll();
      setAvailableMembers(data.users || []);
    } catch (err) {
      console.error('Failed to load members:', err);
    }
  };

  // Load config profiles with inbounds
  const loadConfigProfiles = async () => {
    try {
      const data = await configProfilesApi.getAllWithInbounds();
      setConfigProfiles(data.profiles || []);
    } catch (err) {
      console.error('Failed to load config profiles:', err);
    }
  };

  useEffect(() => {
    loadSquads();
    loadConfigProfiles();
    loadAvailableMembers();
  }, []);

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = () => setOpenDropdown(null);
    if (openDropdown) {
      document.addEventListener('click', handleClickOutside);
    }
    return () => document.removeEventListener('click', handleClickOutside);
  }, [openDropdown]);

  // Open create modal
  const handleCreate = () => {
    setFormData({ name: '', view_position: 0 });
    setShowCreateModal(true);
  };

  // Open rename modal
  const handleRename = (squad) => {
    setSelectedSquad(squad);
    setRenameName(squad.name);
    setShowRenameModal(true);
  };

  const closeInboundsDrawer = () => {
    if (savingInbounds) {
      return;
    }

    setShowInboundsDrawer(false);
    setLoadingInboundsDrawer(false);
    setInboundsTab(INBOUNDS_TAB.profiles);
    setFlatFilter(FLAT_FILTER.all);
    setInboundsSearchQuery('');
    setExpandedProfiles(new Set());
  };

  // Open inbounds drawer
  const handleEditInbounds = async (squad) => {
    setOpenDropdown(null);
    setSelectedSquad(squad);
    setInboundsTab(INBOUNDS_TAB.profiles);
    setFlatFilter(FLAT_FILTER.all);
    setInboundsSearchQuery('');
    setExpandedProfiles(new Set());
    setShowInboundsDrawer(true);
    setLoadingInboundsDrawer(true);

    try {
      const details = await squadsApi.getDetails(squad.uuid);
      const currentInbounds = details.squad?.inbounds || [];
      setSelectedInbounds(new Set(currentInbounds.map((inbound) => inbound.uuid)));
    } catch (err) {
      alert(`Failed to load squad inbounds: ${err.message}`);
      setShowInboundsDrawer(false);
    } finally {
      setLoadingInboundsDrawer(false);
    }
  };

  // Open members modal
  const handleEditMembers = async (squad) => {
    setSelectedSquad(squad);
    try {
      const members = await squadsApi.getMembers(squad.uuid);
      const currentMemberIds = new Set(members.squad_members.map((member) => member.user_id || member.t_id));
      setSelectedMemberIds(currentMemberIds);
      setShowMembersModal(true);
    } catch (err) {
      alert(`Failed to load squad members: ${err.message}`);
    }
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
      await squadsApi.delete(deleteTarget.uuid);
      setDeleteTarget(null);
      await loadSquads();
    } catch (err) {
      alert(`Failed to delete: ${err.message}`);
    } finally {
      setDeleting(false);
    }
  };

  // Handle create submit
  const handleCreateSubmit = async (event) => {
    event.preventDefault();
    try {
      await squadsApi.create(formData);
      setShowCreateModal(false);
      await loadSquads();
    } catch (err) {
      alert(`Failed to create: ${err.message}`);
    }
  };

  // Handle rename submit
  const handleRenameSubmit = async (event) => {
    event.preventDefault();
    try {
      await squadsApi.update(selectedSquad.uuid, { name: renameName });
      setShowRenameModal(false);
      await loadSquads();
    } catch (err) {
      alert(`Failed to rename: ${err.message}`);
    }
  };

  // Handle save inbounds
  const handleSaveInbounds = async () => {
    if (!selectedSquad || savingInbounds) {
      return;
    }

    try {
      setSavingInbounds(true);
      await squadsApi.setInbounds(selectedSquad.uuid, Array.from(selectedInbounds));
      setShowInboundsDrawer(false);
      await loadSquads();
    } catch (err) {
      alert(`Failed to save inbounds: ${err.message}`);
    } finally {
      setSavingInbounds(false);
    }
  };

  // Handle save members
  const handleSaveMembers = async () => {
    try {
      await squadsApi.setMembers(selectedSquad.uuid, Array.from(selectedMemberIds));
      setShowMembersModal(false);
      await loadSquads();
    } catch (err) {
      alert(`Failed to save members: ${err.message}`);
    }
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

  // Toggle inbound selection
  const toggleInbound = (inboundUuid) => {
    setSelectedInbounds((previous) => {
      const next = new Set(previous);
      if (next.has(inboundUuid)) {
        next.delete(inboundUuid);
      } else {
        next.add(inboundUuid);
      }
      return next;
    });
  };

  const setProfileSelection = (profileInbounds, shouldSelect) => {
    setSelectedInbounds((previous) => {
      const next = new Set(previous);
      profileInbounds.forEach((inbound) => {
        if (shouldSelect) {
          next.add(inbound.uuid);
        } else {
          next.delete(inbound.uuid);
        }
      });
      return next;
    });
  };

  const toggleProfileAccordion = (profileUuid) => {
    setExpandedProfiles((previous) => {
      const next = new Set(previous);
      if (next.has(profileUuid)) {
        next.delete(profileUuid);
      } else {
        next.add(profileUuid);
      }
      return next;
    });
  };

  // Toggle member selection
  const toggleMember = (userId) => {
    const next = new Set(selectedMemberIds);
    if (next.has(userId)) {
      next.delete(userId);
    } else {
      next.add(userId);
    }
    setSelectedMemberIds(next);
  };

  // Toggle dropdown
  const toggleDropdown = (event, squadUuid) => {
    event.stopPropagation();
    setOpenDropdown(openDropdown === squadUuid ? null : squadUuid);
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

    const reordered = [...squads];
    const [moved] = reordered.splice(dragIndex, 1);
    reordered.splice(dropIndex, 0, moved);

    setSquads(reordered);
    setDragIndex(null);

    try {
      await squadsApi.reorder(reordered.map((squad) => squad.uuid));
      setSquads((previous) => previous.map((squad, index) => ({ ...squad, view_position: index })));
    } catch (err) {
      alert(`Failed to reorder squads: ${err.message}`);
      await loadSquads();
    }
  };

  const normalizedSearch = searchQuery.trim().toLowerCase();
  const canReorder = normalizedSearch === '';
  const filteredSquads = squads.filter((squad) => {
    if (!normalizedSearch) return true;
    const name = squad.name?.toLowerCase() || '';
    const uuid = squad.uuid?.toLowerCase() || '';
    return name.includes(normalizedSearch) || uuid.includes(normalizedSearch);
  });

  const normalizedInboundsSearch = inboundsSearchQuery.trim().toLowerCase();
  const allInbounds = useMemo(() => {
    return configProfiles.flatMap((profile) => {
      const profileInbounds = Array.isArray(profile.inbounds) ? profile.inbounds : [];
      return profileInbounds.map((inbound) => ({
        ...inbound,
        profileUuid: profile.uuid,
        profileName: profile.name,
      }));
    });
  }, [configProfiles]);

  const inboundMatchesSearch = (inbound, profileName = '') => {
    if (!normalizedInboundsSearch) {
      return true;
    }

    const searchable = [
      inbound?.tag,
      inbound?.type,
      String(getInboundPort(inbound)),
      profileName,
      inbound?.security,
    ]
      .filter(Boolean)
      .join(' ')
      .toLowerCase();

    return searchable.includes(normalizedInboundsSearch);
  };

  const profileGroups = useMemo(() => {
    return configProfiles
      .map((profile) => {
        const profileInbounds = Array.isArray(profile.inbounds) ? profile.inbounds : [];
        const visibleInbounds = profileInbounds.filter((inbound) =>
          inboundMatchesSearch(inbound, profile.name)
        );
        const selectedCount = profileInbounds.reduce(
          (accumulator, inbound) => accumulator + (selectedInbounds.has(inbound.uuid) ? 1 : 0),
          0
        );

        return {
          ...profile,
          profileInbounds,
          visibleInbounds,
          selectedCount,
        };
      })
      .filter((profile) => {
        if (!normalizedInboundsSearch) {
          return true;
        }

        return (
          profile.visibleInbounds.length > 0 ||
          profile.name?.toLowerCase().includes(normalizedInboundsSearch)
        );
      });
  }, [configProfiles, normalizedInboundsSearch, selectedInbounds]);

  const flatInbounds = useMemo(() => {
    let result = allInbounds.filter((inbound) =>
      inboundMatchesSearch(inbound, inbound.profileName)
    );

    if (flatFilter === FLAT_FILTER.selected) {
      result = result.filter((inbound) => selectedInbounds.has(inbound.uuid));
    } else if (flatFilter === FLAT_FILTER.unselected) {
      result = result.filter((inbound) => !selectedInbounds.has(inbound.uuid));
    }

    return result;
  }, [allInbounds, flatFilter, normalizedInboundsSearch, selectedInbounds]);

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="page">
      <TilesToolbar
        title="Сквады"
        searchValue={searchQuery}
        onSearchChange={setSearchQuery}
        searchPlaceholder="Поиск сквадов..."
        onReload={loadSquads}
        onCreate={handleCreate}
        createTitle="Добавить сквад"
      />

      {error && (
        <div className="alert alert-error">
          {error}
          <button type="button" className="btn-icon" onClick={loadSquads}>↻</button>
        </div>
      )}

      {apiNotAvailable && (
        <div className="alert alert-error">
          <div>
            <strong>Internal Squads API is not available</strong>
            <p className="alert-details">
              The backend endpoints for Internal Squads are not implemented yet.
            </p>
          </div>
        </div>
      )}

      {squads.length === 0 ? (
        <div className="empty-state">
          <p>No internal squads found</p>
          <button className="btn btn-primary" onClick={handleCreate}>
            Create your first squad
          </button>
        </div>
      ) : filteredSquads.length === 0 ? (
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
          {filteredSquads.map((squad, index) => (
            <div key={squad.uuid} className="templates-grid-item">
              <div
                className={`templates-item-wrapper card-tile templates-window ${dragIndex === index ? 'dragging' : ''} ${openDropdown === squad.uuid ? 'menu-open' : ''}`}
                draggable={canReorder}
                onDragStart={(event) => handleDragStart(event, index)}
                onDragOver={(event) => canReorder && event.preventDefault()}
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
                      <h3 className="card-tile-title" title={squad.name}>{squad.name}</h3>
                      <p className="templates-subtitle">Internal Squad</p>
                    </div>
                  </div>

                  <div className="templates-window-stats">
                    <div className="templates-window-stat">
                      <span className="templates-window-stat-value">{squad.inbounds_count || 0}</span>
                      <span className="templates-window-stat-label">Inbounds</span>
                    </div>
                    <div className="templates-window-stat">
                      <span className="templates-window-stat-value">{squad.members_count || 0}</span>
                      <span className="templates-window-stat-label">Users</span>
                    </div>
                  </div>

                  <div className="templates-window-body">
                    <div className="card-tile-actions templates-actions">
                      <button className="btn btn-primary btn-block templates-edit-btn" onClick={() => handleEditInbounds(squad)}>
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                          <path d="M7 7h-1a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-1"></path>
                          <path d="M20.385 6.585a2.1 2.1 0 0 0 -2.97 -2.97l-8.415 8.385v3h3l8.385 -8.415z"></path>
                          <path d="M16 5l3 3"></path>
                        </svg>
                        <span>Изменить</span>
                      </button>
                      <div className="templates-menu-wrap" onClick={(event) => event.stopPropagation()}>
                        <button
                          className={`templates-menu-control ${openDropdown === squad.uuid ? 'open' : ''}`}
                          onClick={(event) => toggleDropdown(event, squad.uuid)}
                          title="More actions"
                          aria-label="More actions"
                        >
                          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                            <path d="M6 9l6 6l6 -6"></path>
                          </svg>
                        </button>
                        {openDropdown === squad.uuid && (
                          <div className="templates-dropdown-menu">
                            <button className="templates-dropdown-item" onClick={() => handleEditMembers(squad)}>
                              Edit Members
                            </button>
                            <button className="templates-dropdown-item" onClick={() => handleRename(squad)}>
                              Rename
                            </button>
                            <button className="templates-dropdown-item" onClick={() => copyUuid(squad.uuid)}>
                              Copy UUID
                            </button>
                            <button className="templates-dropdown-item danger" onClick={() => handleDelete(squad.uuid, squad.name)}>
                              Delete Squad
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

      {/* Create Modal */}
      {showCreateModal && (
        <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
          <div className="modal" onClick={(event) => event.stopPropagation()}>
            <div className="modal-header">
              <h2>Create Squad</h2>
              <button className="btn-icon" onClick={() => setShowCreateModal(false)}>✕</button>
            </div>
            <form onSubmit={handleCreateSubmit}>
              <div className="form-group">
                <label>Name *</label>
                <input
                  type="text"
                  value={formData.name}
                  onChange={(event) => setFormData((previous) => ({ ...previous, name: event.target.value }))}
                  required
                  placeholder="Squad name"
                />
              </div>
              <div className="modal-footer">
                <button type="button" className="btn btn-secondary" onClick={() => setShowCreateModal(false)}>
                  Cancel
                </button>
                <button type="submit" className="btn btn-primary">
                  Create Squad
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Rename Modal */}
      {showRenameModal && selectedSquad && (
        <div className="modal-overlay" onClick={() => setShowRenameModal(false)}>
          <div className="modal" onClick={(event) => event.stopPropagation()}>
            <div className="modal-header">
              <h2>Rename Squad</h2>
              <button className="btn-icon" onClick={() => setShowRenameModal(false)}>✕</button>
            </div>
            <form onSubmit={handleRenameSubmit}>
              <div className="form-group">
                <label>Name *</label>
                <input
                  type="text"
                  value={renameName}
                  onChange={(event) => setRenameName(event.target.value)}
                  required
                  placeholder="Squad name"
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

      <SideDrawerPanel
        open={showInboundsDrawer}
        onClose={closeInboundsDrawer}
        width="calc(100vw / 3)"
        title="Изменить сквад"
        icon={(
          <svg stroke="currentColor" fill="none" strokeWidth="2" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round" xmlns="http://www.w3.org/2000/svg">
            <path d="M9.183 6.117a6 6 0 1 0 4.511 3.986"></path>
            <path d="M14.813 17.883a6 6 0 1 0 -4.496 -3.954"></path>
          </svg>
        )}
      >
        <div className="squad-editor-stack">
          <div className="squad-editor-summary">
            <div className="squad-editor-summary-header">
              <div className="squad-editor-summary-identity">
                <button type="button" className="squad-editor-summary-icon" tabIndex={-1} aria-hidden="true">
                  <svg stroke="currentColor" fill="none" strokeWidth="2" viewBox="0 0 24 24" strokeLinecap="round" strokeLinejoin="round" xmlns="http://www.w3.org/2000/svg">
                    <path d="M9.183 6.117a6 6 0 1 0 4.511 3.986"></path>
                    <path d="M14.813 17.883a6 6 0 1 0 -4.496 -3.954"></path>
                  </svg>
                </button>
                <div className="squad-editor-summary-title-wrap">
                  <h3 className="squad-editor-summary-title" title={selectedSquad?.name || ''}>
                    {selectedSquad?.name || 'Squad'}
                  </h3>
                </div>
              </div>
              <div className="squad-editor-summary-stats">
                <span className="squad-editor-pill">Выбрано: {selectedInbounds.size}</span>
                <span className="squad-editor-pill muted">Всего: {allInbounds.length}</span>
              </div>
            </div>
            <div className="squad-editor-summary-actions">
              <button
                type="button"
                className="squad-editor-summary-btn cancel"
                onClick={closeInboundsDrawer}
                disabled={savingInbounds}
              >
                Отмена
              </button>
              <button
                type="button"
                className="squad-editor-summary-btn save"
                onClick={handleSaveInbounds}
                disabled={!selectedSquad || savingInbounds}
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                  <path d="M6 4h10l4 4v10a2 2 0 0 1 -2 2h-12a2 2 0 0 1 -2 -2v-12a2 2 0 0 1 2 -2"></path>
                  <path d="M12 14m-2 0a2 2 0 1 0 4 0a2 2 0 1 0 -4 0"></path>
                  <path d="M14 4l0 4l-6 0l0 -4"></path>
                </svg>
                <span>{savingInbounds ? 'Сохранение...' : 'Сохранить'}</span>
              </button>
            </div>
          </div>

          <div className="search-box search-box-compact squad-editor-search">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
              <circle cx="11" cy="11" r="8" />
              <line x1="21" y1="21" x2="16.65" y2="16.65" />
            </svg>
            <input
              type="text"
              value={inboundsSearchQuery}
              onChange={(event) => setInboundsSearchQuery(event.target.value)}
              placeholder="Поиск по профилям или инбаундам..."
            />
          </div>

          <div className="squad-editor-tabs" role="tablist" aria-label="Inbounds view mode">
            <button
              type="button"
              className={`squad-editor-tab ${inboundsTab === INBOUNDS_TAB.profiles ? 'active' : ''}`}
              onClick={() => setInboundsTab(INBOUNDS_TAB.profiles)}
              role="tab"
              aria-selected={inboundsTab === INBOUNDS_TAB.profiles}
            >
              Профили
            </button>
            <button
              type="button"
              className={`squad-editor-tab ${inboundsTab === INBOUNDS_TAB.list ? 'active' : ''}`}
              onClick={() => setInboundsTab(INBOUNDS_TAB.list)}
              role="tab"
              aria-selected={inboundsTab === INBOUNDS_TAB.list}
            >
              Список
            </button>
          </div>

          {loadingInboundsDrawer ? (
            <div className="squad-editor-empty">Загрузка данных сквада...</div>
          ) : inboundsTab === INBOUNDS_TAB.profiles ? (
            <div className="squad-editor-list">
              {profileGroups.length === 0 ? (
                <div className="squad-editor-empty">Профили или инбаунды не найдены.</div>
              ) : (
                profileGroups.map((profile) => (
                  <section
                    key={profile.uuid}
                    className={`squad-editor-group ${expandedProfiles.has(profile.uuid) ? 'expanded' : ''}`}
                  >
                    <div className="squad-editor-group-row">
                      <button
                        type="button"
                        className="squad-editor-group-toggle"
                        aria-expanded={expandedProfiles.has(profile.uuid)}
                        onClick={() => toggleProfileAccordion(profile.uuid)}
                      >
                        <span className={`squad-editor-group-chevron ${expandedProfiles.has(profile.uuid) ? 'open' : ''}`}>
                          <svg viewBox="0 0 15 15" fill="currentColor" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                            <path
                              fillRule="evenodd"
                              clipRule="evenodd"
                              d="M3.13523 6.15803C3.3241 5.95657 3.64052 5.94637 3.84197 6.13523L7.5 9.56464L11.158 6.13523C11.3595 5.94637 11.6759 5.95657 11.8648 6.15803C12.0536 6.35949 12.0434 6.67591 11.842 6.86477L7.84197 10.6148C7.64964 10.7951 7.35036 10.7951 7.15803 10.6148L3.15803 6.86477C2.95657 6.67591 2.94637 6.35949 3.13523 6.15803Z"
                            />
                          </svg>
                        </span>
                        <span className="squad-editor-group-toggle-content">
                          <span className="squad-editor-group-title">{profile.name}</span>
                          <span className="squad-editor-group-badges">
                            <span className="squad-editor-badge">
                              {profile.selectedCount} / {profile.profileInbounds.length}
                            </span>
                            <span className="squad-editor-badge muted">{profile.visibleInbounds.length}</span>
                          </span>
                        </span>
                      </button>
                      <div className="squad-editor-group-actions">
                        <button
                          type="button"
                          className="squad-editor-group-icon-btn"
                          onClick={(event) => {
                            event.stopPropagation();
                            setProfileSelection(profile.profileInbounds, true);
                          }}
                          disabled={profile.profileInbounds.length === 0}
                          aria-label="Выбрать все инбаунды профиля"
                          title="Выбрать все"
                        >
                          <svg viewBox="0 0 256 256" fill="currentColor" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                            <path d="M232.49,80.49l-128,128a12,12,0,0,1-17,0l-56-56a12,12,0,1,1,17-17L96,183,215.51,63.51a12,12,0,0,1,17,17Z" />
                          </svg>
                        </button>
                        <button
                          type="button"
                          className="squad-editor-group-icon-btn"
                          onClick={(event) => {
                            event.stopPropagation();
                            setProfileSelection(profile.profileInbounds, false);
                          }}
                          disabled={profile.profileInbounds.length === 0}
                          aria-label="Снять все инбаунды профиля"
                          title="Снять все"
                        >
                          <svg viewBox="0 0 256 256" fill="currentColor" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                            <path d="M208.49,191.51a12,12,0,0,1-17,17L128,145,64.49,208.49a12,12,0,0,1-17-17L111,128,47.51,64.49a12,12,0,0,1,17-17L128,111l63.51-63.52a12,12,0,0,1,17,17L145,128Z" />
                          </svg>
                        </button>
                      </div>
                    </div>

                    {expandedProfiles.has(profile.uuid) ? (
                      <div className="squad-editor-group-panel">
                        {profile.profileInbounds.length === 0 ? (
                          <p className="squad-editor-empty small">В этом профиле нет инбаундов.</p>
                        ) : profile.visibleInbounds.length === 0 ? (
                          <p className="squad-editor-empty small">По запросу ничего не найдено.</p>
                        ) : (
                          <div className="squad-editor-items">
                            {profile.visibleInbounds.map((inbound) => (
                              <label key={inbound.uuid} className="squad-editor-checkbox">
                                <input
                                  type="checkbox"
                                  checked={selectedInbounds.has(inbound.uuid)}
                                  onChange={() => toggleInbound(inbound.uuid)}
                                />
                                <span className="squad-editor-checkbox-body">
                                  <span className="squad-editor-inbound-title">{inbound.tag || 'untagged'}</span>
                                  <span className="squad-editor-inbound-badges">
                                    <span className="squad-editor-inline-pill protocol">{getInboundProtocol(inbound)}</span>
                                    <span className="squad-editor-inline-pill port">{String(getInboundPort(inbound))}</span>
                                    <span className="squad-editor-inline-pill security">{getInboundSecurity(inbound)}</span>
                                  </span>
                                </span>
                              </label>
                            ))}
                          </div>
                        )}
                      </div>
                    ) : null}
                  </section>
                ))
              )}
            </div>
          ) : (
            <>
              <div className="squad-editor-flat-filters" role="radiogroup" aria-label="Inbound list filters">
                <button
                  type="button"
                  className={`squad-editor-filter ${flatFilter === FLAT_FILTER.all ? 'active' : ''}`}
                  onClick={() => setFlatFilter(FLAT_FILTER.all)}
                >
                  Все
                </button>
                <button
                  type="button"
                  className={`squad-editor-filter ${flatFilter === FLAT_FILTER.selected ? 'active' : ''}`}
                  onClick={() => setFlatFilter(FLAT_FILTER.selected)}
                >
                  Выбранные
                </button>
                <button
                  type="button"
                  className={`squad-editor-filter ${flatFilter === FLAT_FILTER.unselected ? 'active' : ''}`}
                  onClick={() => setFlatFilter(FLAT_FILTER.unselected)}
                >
                  Невыбранные
                </button>
              </div>

              <div className="squad-editor-list">
                {flatInbounds.length === 0 ? (
                  <div className="squad-editor-empty">Список инбаундов пуст для текущего фильтра.</div>
                ) : (
                  <div className="squad-editor-items flat">
                    {flatInbounds.map((inbound) => (
                      <label key={inbound.uuid} className="squad-editor-checkbox squad-editor-checkbox-flat">
                        <input
                          type="checkbox"
                          checked={selectedInbounds.has(inbound.uuid)}
                          onChange={() => toggleInbound(inbound.uuid)}
                        />
                        <span className="squad-editor-checkbox-body">
                          <span className="squad-editor-inbound-main">
                            <span className="squad-editor-inbound-title">{inbound.tag || 'untagged'}</span>
                            <span className="squad-editor-inbound-subtitle">{inbound.profileName || 'profile'}</span>
                          </span>
                          <span className="squad-editor-inbound-badges">
                            <span className="squad-editor-inline-pill protocol">{getInboundProtocol(inbound)}</span>
                            <span className="squad-editor-inline-pill port">{String(getInboundPort(inbound))}</span>
                            <span className="squad-editor-inline-pill security">{getInboundSecurity(inbound)}</span>
                          </span>
                        </span>
                      </label>
                    ))}
                  </div>
                )}
              </div>
            </>
          )}
        </div>
      </SideDrawerPanel>

      {/* Members Modal */}
      {showMembersModal && (
        <div className="modal-overlay" onClick={() => setShowMembersModal(false)}>
          <div className="modal" onClick={(event) => event.stopPropagation()}>
            <div className="modal-header">
              <h2>Edit Squad Members - {selectedSquad?.name}</h2>
              <button className="btn-icon" onClick={() => setShowMembersModal(false)}>✕</button>
            </div>
            <div className="members-list">
              <p className="members-instruction">
                Check the users you want to add to this squad:
              </p>
              {availableMembers.length === 0 ? (
                <p className="no-members">No users available</p>
              ) : (
                <div className="member-checkboxes">
                  {availableMembers.map((user) => {
                    const userId = user.t_id || user.user_id;
                    return (
                      <label key={userId} className="checkbox-label">
                        <input
                          type="checkbox"
                          checked={selectedMemberIds.has(userId)}
                          onChange={() => toggleMember(userId)}
                        />
                        <span className="checkbox-text">
                          <strong>{user.username}</strong>
                          <span className="member-details">
                            {user.email || 'No email'} • {user.status || 'UNKNOWN'}
                          </span>
                        </span>
                      </label>
                    );
                  })}
                </div>
              )}
            </div>
            <div className="modal-footer">
              <button type="button" className="btn btn-secondary" onClick={() => setShowMembersModal(false)}>
                Cancel
              </button>
              <button type="button" className="btn btn-primary" onClick={handleSaveMembers}>
                Save Members
              </button>
            </div>
          </div>
        </div>
      )}

      <ConfirmActionModal
        open={Boolean(deleteTarget)}
        title="Удалить"
        message={`Вы уверены, что хотите удалить сквад "${deleteTarget?.name || ''}"? Это действие нельзя отменить.`}
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

export default InternalSquads;
