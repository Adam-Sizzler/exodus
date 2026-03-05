import { useEffect, useState } from 'react';
import { nodesApi, usersApi } from '../api';

const ICONS = {
  nodes: (
    <svg viewBox="0 0 256 256" fill="currentColor">
      <path d="M232,56H24A16,16,0,0,0,8,72V200a8,8,0,0,0,16,0V184H40v16a8,8,0,0,0,16,0V184H72v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V184h16v16a8,8,0,0,0,16,0V72A16,16,0,0,0,232,56ZM24,72H232v96H24Zm88,80a8,8,0,0,0,8-8V96a8,8,0,0,0-8-8H48a8,8,0,0,0-8,8v48a8,8,0,0,0,8,8ZM56,104h48v32H56Zm88,48h64a8,8,0,0,0,8-8V96a8,8,0,0,0-8-8H144a8,8,0,0,0-8,8v48A8,8,0,0,0,144,152Zm8-48h48v32H152Z"></path>
    </svg>
  ),
  users: (
    <svg viewBox="0 0 256 256" fill="currentColor">
      <path d="M117.25,157.92a60,60,0,1,0-66.5,0A95.83,95.83,0,0,0,3.53,195.63a8,8,0,1,0,13.4,8.74,80,80,0,0,1,134.14,0,8,8,0,0,0,13.4-8.74A95.83,95.83,0,0,0,117.25,157.92ZM40,108a44,44,0,1,1,44,44A44.05,44.05,0,0,1,40,108Zm210.14,98.7a8,8,0,0,1-11.07-2.33A79.83,79.83,0,0,0,172,168a8,8,0,0,1,0-16,44,44,0,1,0-16.34-84.87,8,8,0,1,1-5.94-14.85,60,60,0,0,1,55.53,105.64,95.83,95.83,0,0,1,47.22,37.71A8,8,0,0,1,250.14,206.7Z"></path>
    </svg>
  ),
  active: (
    <svg viewBox="0 0 256 256" fill="currentColor">
      <path d="M232,128a104,104,0,1,1-41.42-83.14,8,8,0,1,1-9.62,12.78A88,88,0,1,0,216,128a8,8,0,0,1,16,0Zm-54.34-50.34-64,64a8,8,0,0,1-11.32,0l-24-24a8,8,0,0,1,11.32-11.32L108,124.69l58.34-58.35a8,8,0,0,1,11.32,11.32Z"></path>
    </svg>
  ),
  traffic: (
    <svg viewBox="0 0 256 256" fill="currentColor">
      <path d="M240,128a8,8,0,0,1-8,8H204.94l-37.78,75.58A8,8,0,0,1,160,216h-.4a8,8,0,0,1-7.08-5.14L95.35,60.76,63.28,131.31A8,8,0,0,1,56,136H24a8,8,0,0,1,0-16H50.85L88.72,36.69a8,8,0,0,1,14.76.46l57.51,151,31.85-63.71A8,8,0,0,1,200,120h32A8,8,0,0,1,240,128Z"></path>
    </svg>
  ),
  system: (
    <svg viewBox="0 0 256 256" fill="currentColor">
      <path d="M237.94,107.21a8,8,0,0,0-3.89-5.4l-29.83-17-.12-33.62a8,8,0,0,0-2.83-6.08,111.91,111.91,0,0,0-36.72-20.67,8,8,0,0,0-6.46.59L128,41.85,97.88,25a8,8,0,0,0-6.47-.6A111.92,111.92,0,0,0,54.73,45.15a8,8,0,0,0-2.83,6.07l-.15,33.65-29.83,17a8,8,0,0,0-3.89,5.4,106.47,106.47,0,0,0,0,41.56,8,8,0,0,0,3.89,5.4l29.83,17,.12,33.63a8,8,0,0,0,2.83,6.08,111.91,111.91,0,0,0,36.72,20.67,8,8,0,0,0,6.46-.59L128,214.15,158.12,231a7.91,7.91,0,0,0,3.9,1,8.09,8.09,0,0,0,2.57-.42,112.1,112.1,0,0,0,36.68-20.73,8,8,0,0,0,2.83-6.07l.15-33.65,29.83-17a8,8,0,0,0,3.89-5.4A106.47,106.47,0,0,0,237.94,107.21ZM128,168a40,40,0,1,1,40-40A40,40,0,0,1,128,168Z"></path>
    </svg>
  ),
};

function StatCard({ theme, icon, title, value, trend }) {
  return (
    <article className="dashboard-stat-card">
      <div className={`dashboard-stat-icon ${theme}`}>{icon}</div>
      <div className="dashboard-stat-content">
        <p className="dashboard-stat-title">{title}</p>
        <p className="dashboard-stat-value">{value}</p>
        {trend ? <p className="dashboard-stat-trend">{trend}</p> : null}
      </div>
    </article>
  );
}

function Dashboard() {
  const [stats, setStats] = useState({
    totalUsers: 0,
    activeUsers: 0,
    totalUpload: 0,
    totalDownload: 0,
    totalNodes: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchStats();
    const interval = setInterval(fetchStats, 30000);
    return () => clearInterval(interval);
  }, []);

  const fetchStats = async () => {
    try {
      const [usersData, nodesData] = await Promise.all([
        usersApi.getAll().catch(() => ({ users: [] })),
        nodesApi.getAllWithConfig().catch(() => ({ response: [] })),
      ]);

      const users = usersData.users || [];
      const nodes = Array.isArray(nodesData?.response) ? nodesData.response : [];

      const totalUsers = users.length;
      const activeUsers = users.filter((u) => u.status === 'ACTIVE').length;
      const totalUpload = 0;
      const totalDownload = nodes.reduce((sum, node) => sum + (node.trafficUsedBytes || 0), 0);

      setStats({ totalUsers, activeUsers, totalUpload, totalDownload, totalNodes: nodes.length });
      setLoading(false);
    } catch (err) {
      console.error('Error fetching stats:', err);
      setLoading(false);
    }
  };

  const formatBytes = (bytes) => {
    if (!bytes) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
  };

  if (loading) {
    return (
      <div className="empty-state">
        <div className="spinner"></div>
      </div>
    );
  }

  const inactiveUsers = Math.max(stats.totalUsers - stats.activeUsers, 0);
  const activityRate = stats.totalUsers > 0 ? `${Math.round((stats.activeUsers / stats.totalUsers) * 100)}%` : '0%';
  const totalTraffic = stats.totalUpload + stats.totalDownload;

  return (
    <div className="dashboard-shell">
      <section className="dashboard-section">
        <h4 className="dashboard-section-title">Обзор</h4>
        <div className="dashboard-grid dashboard-grid-4">
          <StatCard theme="theme-blue" icon={ICONS.nodes} title="Ноды" value={stats.totalNodes} />
          <StatCard theme="theme-cyan" icon={ICONS.users} title="Пользователи" value={stats.totalUsers} />
          <StatCard theme="theme-green" icon={ICONS.active} title="Активные" value={stats.activeUsers} trend={`Доля активности: ${activityRate}`} />
          <StatCard theme="theme-orange" icon={ICONS.system} title="Трафик всего" value={formatBytes(totalTraffic)} />
        </div>
      </section>

      <section className="dashboard-section">
        <h4 className="dashboard-section-title">Трафик</h4>
        <div className="dashboard-grid dashboard-grid-3">
          <StatCard theme="theme-blue" icon={ICONS.traffic} title="Upload" value={formatBytes(stats.totalUpload)} trend="Обновляется каждые 30 секунд" />
          <StatCard theme="theme-green" icon={ICONS.traffic} title="Download" value={formatBytes(stats.totalDownload)} trend="Сумма по всем нодам" />
          <StatCard theme="theme-teal" icon={ICONS.traffic} title="Баланс" value={formatBytes(Math.max(stats.totalDownload - stats.totalUpload, 0))} trend="Download - Upload" />
        </div>
      </section>

      <section className="dashboard-section">
        <h4 className="dashboard-section-title">Пользователи</h4>
        <div className="dashboard-grid dashboard-grid-4">
          <StatCard theme="theme-cyan" icon={ICONS.users} title="Всего" value={stats.totalUsers} />
          <StatCard theme="theme-green" icon={ICONS.active} title="Active" value={stats.activeUsers} />
          <StatCard theme="theme-red" icon={ICONS.users} title="Inactive" value={inactiveUsers} />
          <StatCard theme="theme-indigo" icon={ICONS.users} title="Нагрузка на ноду" value={stats.totalNodes ? `${(stats.totalUsers / stats.totalNodes).toFixed(1)}` : '0'} trend="Пользователей на одну ноду" />
        </div>
      </section>
    </div>
  );
}

export default Dashboard;
