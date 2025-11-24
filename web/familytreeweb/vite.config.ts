import { defineConfig } from 'vite';
import react from '@vitejs/plugin-react';
import fs from 'fs';
import path from 'path';

// Paths to certificates
const keyPath = path.resolve(__dirname, '../../hack/certs/server.key');
const certPath = path.resolve(__dirname, '../../hack/certs/server.crt');

// Only enable HTTPS if certificates exist (for local development)
const httpsConfig =
  fs.existsSync(keyPath) && fs.existsSync(certPath)
    ? {
        key: fs.readFileSync(keyPath),
        cert: fs.readFileSync(certPath),
      }
    : undefined;

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    https: httpsConfig,
  },
});
