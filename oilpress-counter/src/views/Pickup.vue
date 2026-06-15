<template>
  <div>
    <h2>取油管理</h2>
    <el-card>
      <template #header>
        <div style="display:flex; justify-content: space-between; align-items: center">
          <span>可取油订单</span>
          <el-tabs v-model="statusFilter" size="small">
            <el-tab-pane label="全部" name="all" />
            <el-tab-pane label="未取" name="pending" />
            <el-tab-pane label="部分取" name="partial" />
            <el-tab-pane label="已取完" name="done" />
          </el-tabs>
        </div>
      </template>
      <el-table :data="filteredList" stripe>
        <el-table-column prop="queue_number" label="排队号" width="100" />
        <el-table-column label="农户">
          <template #default="{ row }">
            {{ row.farmer.name }}
            <div v-if="row.farmer.phone" style="color:#999; font-size:12px">{{ row.farmer.phone }}</div>
          </template>
        </el-table-column>
        <el-table-column label="原料" width="80">
          <template #default="{ row }">{{ row.seed_type === 'peanut' ? '花生' : '菜籽' }}</template>
        </el-table-column>
        <el-table-column prop="actual_oil" label="实际出油(kg)" width="110" />
        <el-table-column prop="picked_up_oil" label="已取(kg)" width="100" />
        <el-table-column label="剩余(kg)" width="100">
          <template #default="{ row }">
            <span style="color: #e6a23c; font-weight: bold">{{ (row.actual_oil - row.picked_up_oil).toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="取油状态" width="100">
          <template #default="{ row }">
            <el-tag :type="pickupTag(row.pickup_status)">{{ pickupText(row.pickup_status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="completed_at" label="榨完时间" width="160">
          <template #default="{ row }">{{ formatTime(row.completed_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button v-if="row.pickup_status !== 'done'" type="primary" size="small" @click="showPickup(row)">取油</el-button>
            <el-button type="success" size="small" @click="showStorage(row)" :disabled="row.pickup_status === 'done'">
              {{ row.pickup_status === 'done' ? '已取完' : '寄存' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="pickupVisible" title="取油登记" width="420px">
      <el-form label-width="100px">
        <el-form-item label="排队号">{{ pickupOrder?.queue_number }}</el-form-item>
        <el-form-item label="农户">{{ pickupOrder?.farmer.name }}</el-form-item>
        <el-form-item label="实际出油">{{ pickupOrder?.actual_oil }} 公斤</el-form-item>
        <el-form-item label="已取">
          <span style="color:#67c23a">{{ pickupOrder?.picked_up_oil }} 公斤</span>
        </el-form-item>
        <el-form-item label="本次取油">
          <el-input-number v-model="pickupAmount" :min="0" :max="remaining" :precision="2" :step="0.5" style="width: 100%" />
          <span style="color:#999; margin-left: 8px">公斤（剩余 {{ remaining }} 公斤）</span>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pickupVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmPickup" :loading="picking">确认取油</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getOrderList, pickupOil } from '../api'
import { useRouter } from 'vue-router'

const router = useRouter()
const list = ref([])
const statusFilter = ref('all')
const pickupVisible = ref(false)
const pickupOrder = ref(null)
const pickupAmount = ref(0)
const picking = ref(false)

const remaining = computed(() => pickupOrder.value ? pickupOrder.value.actual_oil - pickupOrder.value.picked_up_oil : 0)

const filteredList = computed(() => {
  if (statusFilter.value === 'all') return list.value
  return list.value.filter(o => o.pickup_status === statusFilter.value)
})

const loadList = async () => {
  try {
    const res = await getOrderList({ status: 'completed' })
    list.value = res.data
  } catch (e) {}
}

const showPickup = (row) => {
  pickupOrder.value = row
  pickupAmount.value = remaining.value
  pickupVisible.value = true
}

const confirmPickup = async () => {
  if (pickupAmount.value <= 0) {
    ElMessage.warning('请输入取油量')
    return
  }
  try {
    picking.value = true
    await pickupOil(pickupOrder.value.id, { amount: pickupAmount.value })
    ElMessage.success('取油登记成功')
    pickupVisible.value = false
    loadList()
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '操作失败')
  } finally {
    picking.value = false
  }
}

const showStorage = async (row) => {
  try {
    await ElMessageBox.confirm(`是否将此订单的油转至寄存管理？\n点击确认将跳转到寄存页面`, '提示', { type: 'info' })
    sessionStorage.setItem('storage_order', JSON.stringify(row))
    router.push('/storage')
  } catch (e) {}
}

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : ''
const pickupText = (s) => ({ pending: '未取', partial: '部分取', done: '已取完' }[s] || s)
const pickupTag = (s) => ({ pending: 'warning', partial: 'info', done: 'success' }[s] || '')

onMounted(loadList)
</script>
