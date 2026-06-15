<template>
  <div>
    <h2>叫号 / 榨油</h2>
    <el-row :gutter="20">
      <el-col :span="10">
        <el-card>
          <template #header>
            <div style="display:flex; justify-content: space-between; align-items: center">
              <span>当前叫号</span>
            </div>
          </template>
          <div v-if="currentOrder" style="text-align: center; padding: 20px 0">
            <div style="font-size: 72px; font-weight: bold; color: #f56c6c; margin: 20px 0">
              {{ currentOrder.queue_number }}
            </div>
            <div style="font-size: 24px; margin: 10px 0">
              {{ currentOrder.farmer.name }}
            </div>
            <el-descriptions :column="1" border style="margin-top: 20px">
              <el-descriptions-item label="原料">{{ currentOrder.seed_type === 'peanut' ? '花生' : '菜籽' }}</el-descriptions-item>
              <el-descriptions-item label="毛重">{{ currentOrder.gross_weight }} 公斤</el-descriptions-item>
              <el-descriptions-item label="饼粕带走">{{ currentOrder.cake_taken ? '是' : '否' }}</el-descriptions-item>
              <el-descriptions-item label="预计出油">{{ currentOrder.expected_oil }} 公斤</el-descriptions-item>
              <el-descriptions-item label="加工费">￥{{ currentOrder.processing_fee }}</el-descriptions-item>
            </el-descriptions>
            <el-divider />
            <el-form label-width="100px">
              <el-form-item label="实际出油">
                <el-input-number v-model="actualOil" :min="0" :precision="2" :step="0.5" style="width: 100%" />
                <span style="color:#999; margin-left:8px">公斤</span>
              </el-form-item>
            </el-form>
            <div style="margin-top: 20px">
              <el-button type="success" size="large" @click="complete" :loading="loading">
                完成榨油，录入实际出油
              </el-button>
            </div>
          </div>
          <el-empty v-else description="暂无正在榨油的订单" />
        </el-card>
      </el-col>
      <el-col :span="14">
        <el-card>
          <template #header>
            <div style="display:flex; justify-content: space-between; align-items: center">
              <span>排队等待中</span>
              <el-button type="primary" @click="callNextFn" :loading="calling">
                <el-icon><Bell /></el-icon> 叫下一位
              </el-button>
            </div>
          </template>
          <el-table :data="waitingList" stripe>
            <el-table-column prop="queue_number" label="排队号" width="100" />
            <el-table-column label="农户">
              <template #default="{ row }">{{ row.farmer.name }}</template>
            </el-table-column>
            <el-table-column label="原料" width="80">
              <template #default="{ row }">{{ row.seed_type === 'peanut' ? '花生' : '菜籽' }}</template>
            </el-table-column>
            <el-table-column prop="gross_weight" label="毛重(kg)" width="100" />
            <el-table-column prop="created_at" label="登记时间" width="170">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button type="primary" size="small" @click="callOne(row.id)">优先叫号</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <el-card style="margin-top: 20px">
          <template #header>今日已完成</template>
          <el-table :data="completedList" stripe size="small">
            <el-table-column prop="queue_number" label="排队号" width="90" />
            <el-table-column label="农户" width="100">
              <template #default="{ row }">{{ row.farmer.name }}</template>
            </el-table-column>
            <el-table-column prop="actual_oil" label="实际出油(kg)" width="110" />
            <el-table-column label="取油状态" width="100">
              <template #default="{ row }">
                <el-tag :type="pickupTag(row.pickup_status)">{{ pickupText(row.pickup_status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="completed_at" label="完成时间" width="160">
              <template #default="{ row }">{{ formatTime(row.completed_at) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { callNext, callSpecific, completeOrder, getOrderList } from '../api'

const waitingList = ref([])
const completedList = ref([])
const currentOrder = ref(null)
const actualOil = ref(0)
const loading = ref(false)
const calling = ref(false)

const loadList = async () => {
  try {
    const [w, c, p] = await Promise.all([
      getOrderList({ status: 'waiting' }),
      getOrderList({ status: 'completed' }),
      getOrderList({ status: 'processing' })
    ])
    waitingList.value = w.data
    completedList.value = c.data
    currentOrder.value = p.data[0] || null
    if (currentOrder.value) {
      actualOil.value = currentOrder.value.expected_oil
    }
  } catch (e) {}
}

const callNextFn = async () => {
  calling.value = true
  try {
    const res = await callNext()
    ElMessage.success(`已叫号：${res.data.queue_number} - ${res.data.farmer.name}`)
    loadList()
  } catch (e) {
    const msg = e?.response?.data?.error
    if (msg && msg.includes('no waiting')) {
      ElMessage.warning('当前没有排队中的订单，请先登记')
    } else if (msg) {
      ElMessage.warning(msg)
    } else if (e?.message) {
      ElMessage.error('网络异常：' + e.message)
    } else {
      ElMessage.warning('没有等待中的订单')
    }
  } finally {
    calling.value = false
  }
}

const callOne = async (id) => {
  try {
    await ElMessageBox.confirm('确认优先叫此号？', '提示', { type: 'warning' })
  } catch (e) {
    return
  }
  try {
    const res = await callSpecific(id)
    ElMessage.success(`已叫号：${res.data.queue_number}`)
    loadList()
  } catch (e) {
    const msg = e?.response?.data?.error
    ElMessage.error(msg || '叫号失败，请稍后重试')
  }
}

const complete = async () => {
  if (!actualOil.value || actualOil.value <= 0) {
    ElMessage.warning('请输入实际出油量')
    return
  }
  try {
    await ElMessageBox.confirm(`确认录入实际出油 ${actualOil.value} 公斤？`, '提示', { type: 'warning' })
  } catch (e) {
    return
  }
  try {
    loading.value = true
    await completeOrder(currentOrder.value.id, { actual_oil: actualOil.value })
    ElMessage.success('已完成')
    currentOrder.value = null
    loadList()
  } catch (e) {
    const msg = e?.response?.data?.error
    ElMessage.error(msg || '操作失败，请稍后重试')
  } finally {
    loading.value = false
  }
}

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : ''
const pickupText = (s) => ({ pending: '未取', partial: '部分取', done: '已取完' }[s] || s)
const pickupTag = (s) => ({ pending: 'warning', partial: 'info', done: 'success' }[s] || '')

watch(currentOrder, (v) => {
  if (v) actualOil.value = v.expected_oil
})

onMounted(loadList)
</script>
