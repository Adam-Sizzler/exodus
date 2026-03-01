// State
let usersData = [];
let searchQuery = '';
let activeTab = localStorage.getItem('v2ray-active-tab') || 'dashboard';
let sidebarOpen = false;

// Format bytes to human readable
function formatBytes(bytes) {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Format timestamp to date
function formatDate(timestamp) {
  if (!timestamp || timestamp === 0) return '-';
  return new Date(timestamp * 1000).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
}

// Get user initials for avatar
function getInitials(name) {
  return name.substring(0, 2).toUpperCase();
}

// Calculate total traffic
function calculateTotalTraffic(data) {
  let totalUpload = 0;
  let totalDownload = 0;
  let totalUsers = 0;
  let activeUsers = 0;

  data.forEach(node => {
    if (node.users) {
      node.users.forEach(user => {
        totalUsers++;
        totalUpload += user.uplink || 0;
        totalDownload += user.downlink || 0;
        if (user.enabled === 'true' || user.enabled === true) {
          activeUsers++;
        }
      });
    }
  });

  return { totalUpload, totalDownload, totalUsers, activeUsers };
}

// Render summary cards
function renderSummaryCards(data) {
  const stats = calculateTotalTraffic(data);
  
  document.getElementById('totalUsers').textContent = stats.totalUsers;
  document.getElementById('activeUsers').textContent = stats.activeUsers;
  document.getElementById('totalUpload').textContent = formatBytes(stats.totalUpload);
  document.getElementById('totalDownload').textContent = formatBytes(stats.totalDownload);
}

// Render users table
function renderUsersTable(data) {
  const tbody = document.getElementById('usersTableBody');
  const filteredData = filterUsers(data);
  
  if (filteredData.length === 0) {
    tbody.innerHTML = `
      <tr>
        <td colspan="7">
          <div class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
              <circle cx="9" cy="7" r="4"/>
              <path d="M23 21v-2a4 4 0 00-3-3.87"/>
              <path d="M16 3.13a4 4 0 010 7.75"/>
            </svg>
            <h3>No users found</h3>
            <p>Try adjusting your search query</p>
          </div>
        </td>
      </tr>
    `;
    return;
  }

  let html = '';
  filteredData.forEach(item => {
    const user = item.user;
    const nodeName = item.nodeName;
    const status = (user.enabled === 'true' || user.enabled === true) ? 'active' : 'inactive';
    const statusText = status === 'active' ? 'Active' : 'Inactive';
    const total = (user.uplink || 0) + (user.downlink || 0);

    html += `
      <tr data-user="${user.user}">
        <td>
          <div style="display: flex; align-items: center; gap: 0.75rem;">
            <div class="user-avatar" style="width: 2.25rem; height: 2.25rem; font-size: 0.875rem;">
              ${getInitials(user.user)}
            </div>
            <strong>${user.user}</strong>
          </div>
        </td>
        <td>${nodeName}</td>
        <td>
          <span class="status-badge ${status}">
            <span class="status-dot"></span>
            ${statusText}
          </span>
        </td>
        <td class="upload">${formatBytes(user.uplink || 0)}</td>
        <td class="download">${formatBytes(user.downlink || 0)}</td>
        <td><strong>${formatBytes(total)}</strong></td>
        <td>${formatDate(user.created)}</td>
      </tr>
    `;
  });

  tbody.innerHTML = html;
}

// Render user detail cards
function renderUserDetails(data) {
  const container = document.getElementById('userDetails');
  const filteredData = filterUsers(data);

  if (filteredData.length === 0) {
    container.innerHTML = '';
    return;
  }

  let html = '';
  filteredData.forEach(item => {
    const user = item.user;
    const nodeName = item.nodeName;
    const status = (user.enabled === 'true' || user.enabled === true) ? 'active' : 'inactive';
    const statusText = status === 'active' ? 'Active' : 'Inactive';
    const total = (user.uplink || 0) + (user.downlink || 0);
    
    // Calculate traffic percentage for progress bar (assuming 100GB cap as example)
    const trafficCap = user.traffic_cap || 107374182400; // 100GB default
    const trafficPercent = Math.min((total / trafficCap) * 100, 100);

    html += `
      <div class="user-detail-card" data-user="${user.user}">
        <div class="user-detail-header">
          <div class="user-detail-name">
            <div class="user-avatar">${getInitials(user.user)}</div>
            <div>
              <div class="user-name">${user.user}</div>
              <div class="user-node">${nodeName}</div>
            </div>
          </div>
          <span class="status-badge ${status}">
            <span class="status-dot"></span>
            ${statusText}
          </span>
        </div>
        <div class="user-stats">
          <div class="user-stat-row">
            <span class="user-stat-label">Upload</span>
            <span class="user-stat-value upload">${formatBytes(user.uplink || 0)}</span>
          </div>
          <div class="user-stat-row">
            <span class="user-stat-label">Download</span>
            <span class="user-stat-value download">${formatBytes(user.downlink || 0)}</span>
          </div>
          <div class="user-stat-row">
            <span class="user-stat-label">Total Traffic</span>
            <span class="user-stat-value total">${formatBytes(total)}</span>
          </div>
          <div class="user-stat-row">
            <span class="user-stat-label">Created</span>
            <span class="user-stat-value">${formatDate(user.created)}</span>
          </div>
          ${user.inbounds && user.inbounds.length > 0 ? `
            <div class="user-stat-row">
              <span class="user-stat-label">Inbound</span>
              <span class="user-stat-value">${user.inbounds[0].inbound_tag || '-'}</span>
            </div>
          ` : ''}
          <div class="progress-bar">
            <div class="progress-fill" style="width: ${trafficPercent}%"></div>
          </div>
        </div>
      </div>
    `;
  });

  container.innerHTML = html;
}

// Filter users based on search query
function filterUsers(data) {
  const result = [];
  
  data.forEach(node => {
    if (node.users) {
      node.users.forEach(user => {
        if (searchQuery === '' || user.user.toLowerCase().includes(searchQuery.toLowerCase())) {
          result.push({
            user: user,
            nodeName: node.node_name || 'Unknown'
          });
        }
      });
    }
  });

  return result;
}

// Fetch data from API
async function fetchData() {
  const refreshBtn = document.getElementById('refreshBtn');
  const tbody = document.getElementById('usersTableBody');
  
  // Add spinning animation to refresh button
  refreshBtn.style.transform = 'rotate(360deg)';
  refreshBtn.style.transition = 'transform 0.5s ease';
  setTimeout(() => {
    refreshBtn.style.transform = 'rotate(0deg)';
  }, 500);

  try {
    const response = await fetch('/api/users');
    
    if (!response.ok) {
      throw new Error('Network response was not ok');
    }

    const data = await response.json();
    usersData = data;

    renderSummaryCards(data);
    renderUsersTable(data);
    renderUserDetails(data);

  } catch (error) {
    console.error('Error fetching data:', error);
    tbody.innerHTML = `
      <tr>
        <td colspan="7">
          <div class="error-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <line x1="12" y1="8" x2="12" y2="12"/>
              <line x1="12" y1="16" x2="12.01" y2="16"/>
            </svg>
            <h3>Failed to load data</h3>
            <p>${error.message}</p>
          </div>
        </td>
      </tr>
    `;
  }
}

// Initialize
document.addEventListener('DOMContentLoaded', () => {
  // Initial fetch
  fetchData();

  // Refresh button
  document.getElementById('refreshBtn').addEventListener('click', fetchData);

  // Search input
  const searchInput = document.getElementById('searchInput');
  searchInput.addEventListener('input', (e) => {
    searchQuery = e.target.value;
    renderUsersTable(usersData);
    renderUserDetails(usersData);
  });

  // Auto refresh every 30 seconds
  setInterval(fetchData, 30000);

  // Sidebar toggle functionality
  const menuToggle = document.getElementById('menuToggle');
  const sidebar = document.querySelector('.sidebar');
  const sidebarOverlay = document.getElementById('sidebarOverlay');
  const navLinks = document.querySelectorAll('.nav-link');

  // Show menu toggle on smaller screens
  const checkScreenSize = () => {
    if (window.innerWidth <= 1024) {
      menuToggle.style.display = 'flex';
    } else {
      menuToggle.style.display = 'none';
      sidebar.classList.remove('open');
      sidebarOverlay.classList.remove('open');
      sidebarOpen = false;
    }
  };

  checkScreenSize();
  window.addEventListener('resize', checkScreenSize);

  // Toggle sidebar
  menuToggle.addEventListener('click', (e) => {
    e.preventDefault();
    sidebarOpen = !sidebarOpen;
    if (sidebarOpen) {
      sidebar.classList.add('open');
      sidebarOverlay.classList.add('open');
    } else {
      sidebar.classList.remove('open');
      sidebarOverlay.classList.remove('open');
    }
  });

  // Close sidebar when clicking overlay
  sidebarOverlay.addEventListener('click', () => {
    sidebarOpen = false;
    sidebar.classList.remove('open');
    sidebarOverlay.classList.remove('open');
  });

  // Handle nav link clicks - persist active tab
  navLinks.forEach((link, index) => {
    link.addEventListener('click', (e) => {
      e.preventDefault();
      const tabNames = ['dashboard', 'users', 'statistics', 'settings'];
      activeTab = tabNames[index];
      localStorage.setItem('v2ray-active-tab', activeTab);

      // Update active state
      document.querySelectorAll('.nav-item').forEach((item, i) => {
        item.classList.toggle('active', i === index);
      });

      // Close sidebar on mobile
      if (window.innerWidth <= 1024) {
        sidebarOpen = false;
        sidebar.classList.remove('open');
        sidebarOverlay.classList.remove('open');
      }
    });
  });

  // Restore active tab on page load
  const tabNames = ['dashboard', 'users', 'statistics', 'settings'];
  const activeIndex = tabNames.indexOf(activeTab);
  if (activeIndex !== -1) {
    document.querySelectorAll('.nav-item').forEach((item, i) => {
      item.classList.toggle('active', i === activeIndex);
    });
  }
});
