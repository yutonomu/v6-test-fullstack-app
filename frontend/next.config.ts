import type { NextConfig } from 'next';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/notes',
        destination: `${API_BASE_URL}/notes`,
      },
    ];
  },
};

export default nextConfig;
