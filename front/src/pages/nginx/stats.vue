<template>
  <a-row :gutter="[16, 16]" class="stats">
    <a-col :xs="24" :sm="12" :lg="6">
      <a-card :bordered="false" class="card">
        <div class="head">
          <span class="title">日志数量</span>
          <span class="meta">{{ windowText }}</span>
        </div>
        <a-statistic :value="data.total" />
        <div class="foot">
          <span class="hint">最后刷新</span>
          <span class="time">{{ data.updatedAt || '-' }}</span>
        </div>
      </a-card>
    </a-col>

    <a-col :xs="24" :sm="12" :lg="6">
      <a-card :bordered="false" class="card">
        <div class="head">
          <span class="title">独立 IP</span>
          <span class="meta">unique</span>
        </div>
        <a-statistic :value="data.ipCount" />
        <div class="foot">
          <span class="hint">统计范围</span>
          <span class="time">当前筛选结果</span>
        </div>
      </a-card>
    </a-col>

    <a-col :xs="24" :sm="12" :lg="6">
      <a-card :bordered="false" class="card">
        <div class="head">
          <span class="title">状态码统计</span>
          <span class="meta">2xx/3xx/4xx/5xx</span>
        </div>

        <div class="status-grid">
          <div class="sitem">
            <span class="label ok">2xx</span>
            <span class="val">{{ data.status['2xx'] }}</span>
          </div>
          <div class="sitem">
            <span class="label info">3xx</span>
            <span class="val">{{ data.status['3xx'] }}</span>
          </div>
          <div class="sitem">
            <span class="label warn">4xx</span>
            <span class="val">{{ data.status['4xx'] }}</span>
          </div>
          <div class="sitem">
            <span class="label err">5xx</span>
            <span class="val">{{ data.status['5xx'] }}</span>
          </div>
        </div>

        <a-tooltip :title="barTip">
          <div class="bar" aria-label="status distribution">
            <span class="seg ok" :style="{ width: pct(data.status['2xx']) }" />
            <span class="seg info" :style="{ width: pct(data.status['3xx']) }" />
            <span class="seg warn" :style="{ width: pct(data.status['4xx']) }" />
            <span class="seg err" :style="{ width: pct(data.status['5xx']) }" />
          </div>
        </a-tooltip>
      </a-card>
    </a-col>

    <a-col :xs="24" :sm="12" :lg="6">
      <a-card :bordered="false" class="card">
        <div class="head">
          <span class="title">Top 路由</span>
          <span class="meta">Top {{ data.topRoutes.length || 0 }}</span>
        </div>

        <div class="routes">
          <div v-if="!data.topRoutes.length" class="empty">
            暂无数据
          </div>
          <div
            v-for="(item, idx) in data.topRoutes"
            :key="item.key"
            class="route"
          >
            <span class="rank">{{ idx + 1 }}</span>
            <a-tooltip :title="item.route">
              <span class="rtext">{{ item.route }}</span>
            </a-tooltip>
            <span class="count">{{ item.count }}</span>
          </div>
        </div>
      </a-card>
    </a-col>
  </a-row>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  data: {
    required: true,
    type: Object,
  },
})

const totalCount = computed(() => Number(props.data.total || 0))

const windowText = computed(() => {
  const w = props.data.window || ''
  return w ? `窗口：${w}` : ''
})

function pct(n) {
  const total = totalCount.value || 0
  const v = Number(n || 0)
  if (!total) return '0%'
  return `${Math.max(0, Math.min(100, (v / total) * 100)).toFixed(2)}%`
}

const barTip = computed(() => {
  const s = props.data.status || {}
  return `2xx: ${pct(s['2xx'])} / 3xx: ${pct(s['3xx'])} / 4xx: ${pct(s['4xx'])} / 5xx: ${pct(s['5xx'])}`
})
</script>

<style scoped lang="scss">

.card {
  border-radius: 10px;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.04);
}

.head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  margin-bottom: 8px;

  .title {
    font-weight: 600;
    color: $text-color-1;
  }
  .meta {
    color: $text-color-4;
    font-size: $font-size-5;
  }
}

.foot {
  margin-top: 8px;
  display: flex;
  justify-content: space-between;
  color: $text-color-4;
  font-size: $font-size-5;
}

.status-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 10px;
  margin: 6px 0 10px;
}

.sitem {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.label {
  font-size: 12px;
  font-weight: 600;
}

.val {
  font-size: 18px;
  font-weight: 700;
  color: $text-color-1;
  line-height: 1.2;
}

.bar {
  height: 10px;
  width: 100%;
  border-radius: 999px;
  background: rgba(0, 0, 0, 0.05);
  overflow: hidden;
  display: flex;
}

.seg {
  height: 100%;
  display: block;
}

.ok {
  color: #10b981;
  background: #10b981;
}

.info {
  color: #3b82f6;
  background: #3b82f6;
}

.warn {
  color: #f59e0b;
  background: #f59e0b;
}

.err {
  color: #ef4444;
  background: #ef4444;
}

.routes {
  margin-top: 4px;
}

.route {
  display: grid;
  grid-template-columns: 18px 1fr auto;
  gap: 10px;
  align-items: center;
  padding: 6px 0;
  border-bottom: 1px dashed rgba(0, 0, 0, 0.06);

  &:last-child {
    border-bottom: none;
  }
}

.rank {
  width: 18px;
  height: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  font-size: 12px;
  color: $text-color-3;
  background: rgba(0, 0, 0, 0.04);
}

.rtext {
  @include overflow;
  max-width: 100%;
  color: $text-color-2;
}

.count {
  font-variant-numeric: tabular-nums;
  color: $text-color-3;
}

.empty {
  padding: 18px 0;
  text-align: center;
  color: $text-color-4;
}
</style>

