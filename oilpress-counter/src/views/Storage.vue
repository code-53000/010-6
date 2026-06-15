<template>
  <div>
    <h2>油桶寄存管理</h2>
    <el-row :gutter="20">
      <el-col :span="10">
        <el-card>
          <template #header>新建寄存</template>
          <el-form :model="form" label-width="100px" ref="formRef">
            <el-form-item label="关联订单号" prop="order_id" required>
              <el-select v-model="form.order_id" placeholder="请选择关联订单" style="width: 100%" filterable>
                <el-option
                  v-for="order in availableOrders"
                  :key="order.id"
                  :value="order.id"
                  :label="`#${order.queue_number} · ${order.farmer.name} · ${order.seed_type === 'peanut' ? '花生' : '菜籽'} · 实油${order.actual_oil}kg`"
                />
              </el-select>
              <div v-if="orderInfo" style="margin-top:6px; color:#67c23a; font-size:13px">
                {{ orderInfo.farmer?.name }} · {{ orderInfo.seed_type === 'peanut' ? '花生' : '菜籽' }} · 实油 {{ orderInfo.actual_oil }}kg
              </div>
              <div v-if="orderInfo" style="margin-top:4px; color:#e6a23c; font-size:13px">
                已取 {{ orderInfo.picked_up_oil }}kg · 剩余可取 {{ (orderInfo.actual_oil - orderInfo.picked_up_oil).toFixed(2) }}kg
              </div>
            </el-form-item>
            <el-form-item label="桶数" prop="barrel_count" required>
              <el-input-number v-model="form.barrel_count" :min="1" style="width: 100%" />
            </el-form-item>
            <el-form-item label="每桶油量(kg)" prop="oil_per_barrel" required>
              <el-input-number v-model="form.oil_per_barrel" :min="0" :precision="2" :step="0.5" style="width: 100%" />
            </el-form-item>
            <el-form-item label="合计油量">
              <span style="color: #409eff; font-weight: bold">{{ (form.barrel_count * form.oil_per_barrel).toFixed(2) }} 公斤</span>
            </el-form-item>
            <el-form-item label="预存天数" prop="stored_days">
              <el-input-number v-model="form.stored_days" :min="0" :step="1" style="width: 100%" />
              <span style="color:#999; font-size:12px; margin-left:8px">每天每桶 0.5 元</span>
            </el-form-item>
            <el-form-item label="寄存费">
              <span style="color: #e6a23c; font-weight: bold">￥{{ (form.stored_days * 0.5 * form.barrel_count).toFixed(2) }}</span>
            </el-form-item>
            <el-form-item label="备注">
              <el-input v-model="form.remark" type="textarea" :rows="2" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submit" :loading="loading">登记寄存</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
      <el-col :span="14">
        <el-card>
          <template #header>
            <div style="display:flex; justify-content: space-between; align-items: center">
              <span>寄存列表</span>
              <el-radio-group v-model="filter" size="small">
                <el-radio-button :value="null">全部</el-radio-button>
                <el-radio-button :value="false">寄存中</el-radio-button>
                <el-radio-button :value="true">已取走</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <el-table :data="list" stripe>
            <el-table-column prop="id" label="寄存号" width="80" />
            <el-table-column label="农户">
              <template #default="{ row }">{{ row.farmer?.name }}</template>
            </el-table-column>
            <el-table-column prop="barrel_count" label="桶数" width="70" />
            <el-table-column prop="oil_per_barrel" label="每桶(kg)" width="90" />
            <el-table-column prop="total_oil" label="总油量(kg)" width="100" />
            <el-table-column prop="stored_days" label="预存天" width="80" />
            <el-table-column prop="storage_fee" label="寄存费" width="90">
              <template #default="{ row }">￥{{ row.storage_fee }}</template>
            </el-table-column>
            <el-table-column label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="row.is_picked_up ? 'success' : 'warning'">{{ row.is_picked_up ? '已取走' : '寄存中' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="存入时间" width="160">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="110">
              <template #default="{ row }">
                <el-button v-if="!row.is_picked_up" type="success" size="small" @click="pickup(row)">取走</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { createStorage, getStorageList, pickupStorage, getOrder, getOrderList } from '../api'

const formRef = ref()
const loading = ref(false)
const list = ref([])
const filter = ref(null)
const orderInfo = ref(null)
const availableOrders = ref([])

const form = reactive({
  order_id: null,
  barrel_count: 1,
  oil_per_barrel: 5,
  stored_days: 7,
  remark: ''
})

const remainingOil = computed(() => {
  if (!orderInfo.value) return 0
  return orderInfo.value.actual_oil - orderInfo.value.picked_up_oil
})

const totalStorageOil = computed(() => {
  return form.barrel_count * form.oil_per_barrel
})

const loadAvailableOrders = async () => {
  try {
    const res = await getOrderList({ status: 'completed' })
    availableOrders.value = res.data.filter(o => o.pickup_status !== 'done')
  } catch (e) {}
}

watch(() => form.order_id, async (id) => {
  orderInfo.value = null
  if (!id) return
  try {
    const res = await getOrder(id)
    orderInfo.value = res.data
  } catch (e) {
    orderInfo.value = null
  }
}, { immediate: true })

const submit = async () => {
  await formRef.value.validate()
  if (totalStorageOil.value > remainingOil.value + 0.001) {
    ElMessage.warning(`寄存油量(${totalStorageOil.value.toFixed(2)}kg)超出剩余可取油量(${remainingOil.value.toFixed(2)}kg)`)
    return
  }
  loading.value = true
  try {
    await createStorage(form)
    ElMessage.success('寄存登记成功')
    loadList()
    loadAvailableOrders()
    form.barrel_count = 1
    form.oil_per_barrel = 5
    form.stored_days = 7
    form.remark = ''
    if (form.order_id) {
      const res = await getOrder(form.order_id)
      orderInfo.value = res.data
    }
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '操作失败')
  } finally {
    loading.value = false
  }
}

const loadList = async () => {
  try {
    const params = {}
    if (filter.value !== null) params.is_picked_up = filter.value
    const res = await getStorageList(params)
    list.value = res.data
  } catch (e) {}
}

const pickup = async (row) => {
  try {
    await ElMessageBox.confirm(`确认取走寄存号 ${row.id}？超期将加收寄存费`, '提示', { type: 'warning' })
    const res = await pickupStorage(row.id)
    ElMessage.success(`取走成功，实收寄存费 ￥${res.data.storage_fee}`)
    loadList()
  } catch (e) {
    if (e !== 'cancel') ElMessage.error(e.response?.data?.error || '操作失败')
  }
}

const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : ''

watch(filter, loadList)

onMounted(() => {
  loadAvailableOrders()
  const orderData = sessionStorage.getItem('storage_order')
  if (orderData) {
    const o = JSON.parse(orderData)
    form.order_id = o.id
    orderInfo.value = o
    sessionStorage.removeItem('storage_order')
  }
  loadList()
})
</script>
