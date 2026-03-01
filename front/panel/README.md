# v2ray-stat Frontend

React-based dashboard for v2ray-stat backend.

## Configuration

Create a `.env` file in the root directory:

```bash
cp .env.example .env
```

Edit `.env` with your backend settings:

```env
PORT=9000
VITE_API_URL=http://localhost:9243
VITE_API_TOKEN=your-api-token-here
```

## Available Scripts

```bash
# Install dependencies
npm install

# Start development server (runs on http://localhost:9000)
npm run dev

# Start production server (runs on http://localhost:9000)
npm run server

# Build for production
npm run build

# Preview production build
npm run preview
```

## API Endpoints

The frontend communicates with the Go backend via REST API.

### Config Profiles

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/config-profiles` | Get all profiles |
| GET | `/api/v1/config-profiles/{uuid}` | Get profile by UUID |
| POST | `/api/v1/config-profiles` | Create new profile |
| PATCH | `/api/v1/config-profiles/{uuid}` | Update profile |
| DELETE | `/api/v1/config-profiles/{uuid}` | Delete profile |

### Nodes

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/nodes` | Get all nodes |
| GET | `/api/v1/nodes/{uuid}` | Get node by UUID |
| PATCH | `/api/v1/nodes/{uuid}` | Update node |

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/users-list` | Get all users |
| GET | `/api/v1/users-list/{uuid}` | Get user by UUID |
| POST | `/api/v1/users-list/create` | Create user |
| PATCH | `/api/v1/users-list/{uuid}` | Update user |

## Example API Requests

### Get Config Profiles

```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:9243/api/v1/config-profiles
```

### Create Config Profile

```bash
curl -X POST \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "singbox-profile",
    "view_position": 0,
    "config": {
      "log": {"level": "info"},
      "dns": {"servers": ["8.8.8.8"]},
      "inbounds": [{"type": "mixed", "tag": "mixed-in", "listen": "0.0.0.0", "listen_port": 2080}],
      "outbounds": [{"type": "direct", "tag": "direct"}],
      "route": {"rules": []}
    }
  }' \
  http://localhost:9243/api/v1/config-profiles
```

### Update Config Profile

```bash
curl -X PATCH \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "config": {
      "log": {"level": "debug"},
      "dns": {"servers": ["1.1.1.1", "8.8.8.8"]},
      "inbounds": [{"type": "mixed", "tag": "mixed-in", "listen": "0.0.0.0", "listen_port": 2080}],
      "outbounds": [{"type": "direct", "tag": "direct"}],
      "route": {"rules": [], "auto_detect_interface": true}
    }
  }' \
  http://localhost:9243/api/v1/config-profiles/550e8400-e29b-41d4-a716-446655440000
```

## Pages

- **Dashboard** - Overview and statistics
- **Nodes** - Manage server nodes
- **Users** - Manage users
- **Config Profiles** - Manage sing-box configuration profiles (JSON editor)
- **Settings** - Application settings

## Config Profiles Feature

The Config Profiles page allows you to:

1. **View** all sing-box configuration profiles
2. **Create** new profiles with JSON editor
3. **Edit** existing profiles with live JSON validation
4. **Delete** profiles

### JSON Editor

When creating or editing a profile, you can directly edit the JSON configuration. The editor validates JSON in real-time and shows errors if the JSON is invalid.

Example sing-box configuration:

```json
{
  "log": {
    "level": "info",
    "timestamp": true
  },
  "dns": {
    "servers": [
      {"tag": "google", "address": "8.8.8.8"},
      {"tag": "local", "address": "223.5.5.5", "detour": "direct"}
    ]
  },
  "inbounds": [
    {
      "type": "mixed",
      "tag": "mixed-in",
      "listen": "0.0.0.0",
      "listen_port": 2080
    }
  ],
  "outbounds": [
    {"type": "direct", "tag": "direct"}
  ],
  "route": {
    "rules": [],
    "auto_detect_interface": true
  }
}
```

## Tech Stack

- React 18
- Vite
- Fetch API (no external HTTP libraries)
- CSS3
