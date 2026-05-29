<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  size?: number | string
  animated?: boolean
}>(), {
  size: 32,
  animated: true
})

const sizeStyle = computed(() => {
  const s = typeof props.size === 'number' ? `${props.size}px` : props.size
  return { width: s, height: s }
})
</script>

<template>
  <div class="ccx-web-logo" :style="sizeStyle">
    <svg
      viewBox="0 0 100 100"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
      class="ccx-logo-svg"
      aria-hidden="true"
    >
      <defs>
        <!-- 与 App Icon 同源的终端网关流光渐变；颜色由 Vuetify 主题类驱动 -->
        <linearGradient id="ccx-web-logo-flow" x1="18%" y1="28%" x2="82%" y2="72%">
          <stop offset="0%" stop-color="var(--ccx-logo-flow-start)" />
          <stop offset="48%" stop-color="var(--ccx-logo-flow-mid)" />
          <stop offset="100%" stop-color="var(--ccx-logo-flow-end)" />
        </linearGradient>

        <linearGradient id="ccx-web-logo-panel" x1="15%" y1="20%" x2="85%" y2="82%">
          <stop offset="0%" stop-color="var(--ccx-logo-panel-start)" stop-opacity="var(--ccx-logo-panel-start-opacity)" />
          <stop offset="52%" stop-color="var(--ccx-logo-panel-mid)" stop-opacity="var(--ccx-logo-panel-mid-opacity)" />
          <stop offset="100%" stop-color="var(--ccx-logo-panel-end)" stop-opacity="var(--ccx-logo-panel-end-opacity)" />
        </linearGradient>

        <radialGradient id="ccx-web-logo-bg" cx="70%" cy="70%" r="86%">
          <stop offset="0%" stop-color="var(--ccx-logo-bg-start)" />
          <stop offset="40%" stop-color="var(--ccx-logo-bg-mid)" />
          <stop offset="100%" stop-color="var(--ccx-logo-bg-end)" />
        </radialGradient>

        <filter id="ccx-web-logo-glow" x="-28%" y="-28%" width="156%" height="156%">
          <feGaussianBlur stdDeviation="2.2" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>

      <!-- 1. App 图标同源圆角底 -->
      <rect x="5" y="5" width="90" height="90" rx="22" fill="url(#ccx-web-logo-bg)" />
      <rect
        x="7.5" y="7.5" width="85" height="85" rx="20"
        fill="none"
        stroke="var(--ccx-logo-border)"
        stroke-width="0.9"
        opacity="var(--ccx-logo-border-opacity)"
      />

      <!-- 2. 玻璃终端窗口 -->
      <rect
        x="15" y="20" width="70" height="62" rx="14"
        fill="url(#ccx-web-logo-panel)"
        stroke="var(--ccx-logo-panel-border)"
        stroke-width="1.4"
        opacity="0.98"
      />
      <path d="M 19 32 H 81" stroke="var(--ccx-logo-divider)" stroke-width="0.8" opacity="var(--ccx-logo-divider-opacity)" />
      <circle cx="25" cy="26.5" r="2.3" fill="var(--ccx-logo-dot-green)" />
      <circle cx="32" cy="26.5" r="2.3" fill="var(--ccx-logo-dot-blue)" opacity="var(--ccx-logo-dot-opacity)" />
      <circle cx="39" cy="26.5" r="2.3" fill="var(--ccx-logo-dot-indigo)" opacity="var(--ccx-logo-dot-opacity)" />

      <!-- 3. 终端网关提示符与 X 路由束 -->
      <g filter="url(#ccx-web-logo-glow)" stroke-linecap="round" stroke-linejoin="round">
        <path d="M 28 39 L 42 51 L 28 63" stroke="url(#ccx-web-logo-flow)" stroke-width="8" />
        <path d="M 52 38 L 73 64" stroke="url(#ccx-web-logo-flow)" stroke-width="8" />
        <path d="M 73 38 L 52 64" stroke="url(#ccx-web-logo-flow)" stroke-width="8" />
      </g>

      <!-- 4. 底部网关状态线与在线节点 -->
      <path d="M 22 74 H 50" stroke="var(--ccx-logo-status-green)" stroke-width="2.6" stroke-linecap="round" opacity="var(--ccx-logo-status-green-opacity)" />
      <path d="M 55 74 H 68" stroke="var(--ccx-logo-status-blue)" stroke-width="2.6" stroke-linecap="round" opacity="var(--ccx-logo-status-blue-opacity)" />
      <g :class="{ 'animate-gateway-pulse': animated }">
        <circle cx="76" cy="74" r="2.4" fill="var(--ccx-logo-online)" />
        <circle cx="76" cy="74" r="5.5" stroke="var(--ccx-logo-online)" stroke-width="1.1" opacity="var(--ccx-logo-online-ring-opacity)" />
      </g>
    </svg>
  </div>
</template>

<style>
.ccx-web-logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;

  /* 亮色主题默认值 */
  --ccx-logo-bg-start: #d1fae5;
  --ccx-logo-bg-mid: #dbeafe;
  --ccx-logo-bg-end: #f8fafc;
  --ccx-logo-border: #3b82f6;
  --ccx-logo-border-opacity: 0.3;
  --ccx-logo-panel-start: #f0f9ff;
  --ccx-logo-panel-mid: #eff6ff;
  --ccx-logo-panel-end: #ecfdf5;
  --ccx-logo-panel-start-opacity: 0.98;
  --ccx-logo-panel-mid-opacity: 0.95;
  --ccx-logo-panel-end-opacity: 0.98;
  --ccx-logo-panel-border: #3b82f6;
  --ccx-logo-divider: #94a3b8;
  --ccx-logo-divider-opacity: 0.25;
  --ccx-logo-flow-start: #0284c7;
  --ccx-logo-flow-mid: #4f46e5;
  --ccx-logo-flow-end: #059669;
  --ccx-logo-dot-green: #059669;
  --ccx-logo-dot-blue: #0284c7;
  --ccx-logo-dot-indigo: #4f46e5;
  --ccx-logo-dot-opacity: 0.85;
  --ccx-logo-status-green: #059669;
  --ccx-logo-status-blue: #0284c7;
  --ccx-logo-status-green-opacity: 0.55;
  --ccx-logo-status-blue-opacity: 0.42;
  --ccx-logo-online: #14b8a6;
  --ccx-logo-online-ring-opacity: 0.3;
}

.v-theme--dark .ccx-web-logo {
  /* 暗色主题：深底 + 蓝靛绿霓虹 */
  --ccx-logo-bg-start: #064e3b;
  --ccx-logo-bg-mid: #082f49;
  --ccx-logo-bg-end: #020617;
  --ccx-logo-border: #93c5fd;
  --ccx-logo-border-opacity: 0.32;
  --ccx-logo-panel-start: #102a56;
  --ccx-logo-panel-mid: #06142a;
  --ccx-logo-panel-end: #042f2e;
  --ccx-logo-panel-start-opacity: 0.95;
  --ccx-logo-panel-mid-opacity: 0.92;
  --ccx-logo-panel-end-opacity: 0.95;
  --ccx-logo-panel-border: #93c5fd;
  --ccx-logo-divider: #bae6fd;
  --ccx-logo-divider-opacity: 0.18;
  --ccx-logo-flow-start: #38bdf8;
  --ccx-logo-flow-mid: #6366f1;
  --ccx-logo-flow-end: #10b981;
  --ccx-logo-dot-green: #10b981;
  --ccx-logo-dot-blue: #38bdf8;
  --ccx-logo-dot-indigo: #6366f1;
  --ccx-logo-dot-opacity: 0.78;
  --ccx-logo-status-green: #10b981;
  --ccx-logo-status-blue: #38bdf8;
  --ccx-logo-status-green-opacity: 0.46;
  --ccx-logo-status-blue-opacity: 0.34;
  --ccx-logo-online: #5eead4;
  --ccx-logo-online-ring-opacity: 0.24;
}

.ccx-logo-svg {
  width: 100%;
  height: 100%;
}

/* 在线网关节点呼吸脉冲 */
@keyframes gateway-pulse {
  0%, 100% {
    transform: scale(0.92);
    transform-origin: 76px 74px;
    opacity: 0.82;
  }
  50% {
    transform: scale(1.12);
    transform-origin: 76px 74px;
    opacity: 1;
  }
}

.animate-gateway-pulse {
  animation: gateway-pulse 2.4s infinite ease-in-out;
}
</style>
