<template>
  <div class="print-receipt">
    <div class="receipt-container" :class="{ 'a5-size': paperSize === 'a5' }">
      <div class="receipt-header">
        <div class="receipt-shop-name">{{ shopConfig.shopName }}</div>
        <div class="receipt-title">排队登记小票</div>
        <div class="receipt-queue-number">No.{{ orderData.queue_number }}</div>
      </div>

      <div class="receipt-info">
        <div class="receipt-row">
          <span class="receipt-label">登记时间：</span>
          <span class="receipt-value">{{ formatTime(orderData.created_at) }}</span>
        </div>
        <div class="receipt-row">
          <span class="receipt-label">农户姓名：</span>
          <span class="receipt-value">{{ orderData.farmer?.name || '-' }}</span>
        </div>
        <div v-if="orderData.farmer?.phone" class="receipt-row">
          <span class="receipt-label">联系电话：</span>
          <span class="receipt-value">{{ orderData.farmer.phone }}</span>
        </div>
        <div class="receipt-row">
          <span class="receipt-label">原料类型：</span>
          <span class="receipt-value">{{ seedTypeText }}</span>
        </div>
        <div class="receipt-row">
          <span class="receipt-label">毛重：</span>
          <span class="receipt-value">{{ orderData.gross_weight }} 公斤</span>
        </div>
        <div v-if="orderData.net_weight > 0" class="receipt-row">
          <span class="receipt-label">净重：</span>
          <span class="receipt-value">{{ orderData.net_weight }} 公斤</span>
        </div>
        <div class="receipt-row">
          <span class="receipt-label">带走饼粕：</span>
          <span class="receipt-value">{{ orderData.cake_taken ? '是' : '否' }}</span>
        </div>
      </div>

      <div class="receipt-total-section">
        <div class="receipt-total-row">
          <span>预计出油：</span>
          <span class="receipt-oil">{{ orderData.expected_oil || 0 }} 公斤</span>
        </div>
        <div class="receipt-total-row bold">
          <span>加工费：</span>
          <span class="receipt-fee">￥{{ orderData.processing_fee || 0 }}</span>
        </div>
      </div>

      <div v-if="orderData.remark" class="receipt-remark">
        <div class="receipt-row">
          <span class="receipt-label">备注：</span>
        </div>
        <div class="receipt-remark-text">{{ orderData.remark }}</div>
      </div>

      <div class="receipt-footer">
        <p>{{ shopConfig.thankYouMessage }}</p>
        <p v-if="shopConfig.shopAddress">地址：{{ shopConfig.shopAddress }}</p>
        <p v-if="shopConfig.shopPhone">电话：{{ shopConfig.shopPhone }}</p>
        <p>{{ formatDate(orderData.created_at) }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { SHOP_CONFIG, SEED_TYPE_MAP } from '../config'

const props = defineProps({
  orderData: {
    type: Object,
    required: true,
    default: () => ({})
  },
  paperSize: {
    type: String,
    default: '80mm',
    validator: (v) => ['80mm', 'a5'].includes(v)
  }
})

const shopConfig = SHOP_CONFIG

const seedTypeText = computed(() => {
  return SEED_TYPE_MAP[props.orderData.seed_type] || props.orderData.seed_type || '-'
})

const formatTime = (t) => {
  if (!t) return ''
  const d = new Date(t)
  return d.toLocaleString('zh-CN', { hour12: false })
}

const formatDate = (t) => {
  if (!t) return ''
  const d = new Date(t)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}年${m}月${day}日`
}
</script>
