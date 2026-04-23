import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
        host: '0.0.0.0',          // слушаем на всех интерфейсах (Tailscale и т.д.)
        strictPort: true,
        allowedHosts: [
            'higu.su', 
            'app.higu.su'],
        },
    }
  )
