<template>
  <div>
    <h2>排队登记</h2>
    <el-row :gutter="20">
      <el-col :span="10">
        <el-card>
          <template #header>新登记</template>
          <el-form :model="form" :rules="rules" label-width="100px" ref="formRef">
            <el-form-item label="农户姓名" prop="farmer_name" required>
              <el-input v-model="form.farmer_name" placeholder="请输入姓名" />
            </el-form-item>
            <el-form-item label="联系电话">
              <el-input v-model="form.farmer_phone" placeholder="可选" />
            </el-form-item>
            <el-form-item label="原料类型" prop="seed_type" required>
              <el-select v-model="form.seed_type" placeholder="请选择" style="width: 100%" @change="calcFee">
                <el-option label="花生" value="peanut" />
                <el-option label="菜籽" value="rapeseed" />
              </el-select>
            </el-form-item>
            <el-form-item label="毛重(公斤)" prop="gross_weight" required>
              <el-input-number v-model="form.gross_weight" :min="0" :precision="2" :step="1" style="width: 100%" @change="calcFee" />
            </el-form-item>
            <el-form-item label="净重(公斤)">
              <el-input-number v-model="form.net_weight" :min="0" :precision="2" :step="1" style="width: 100%" />
            </el-form-item>
            <el-form-item label="是否带走饼粕">
              <el-switch v-model="form.cake_taken" @change="calcFee" />
            </el-form-item>
            <el-form-item label="备注">
              <el-input v-model="form.remark" type="textarea" :rows="2" />
            </el-form-item>
            <el-divider />
            <el-form-item label="预计出油">
              <span style="color: #409eff; font-size: 18px; font-weight: bold">{{ preview.expected_oil || 0 }} 公斤</span>
            </el-form-item>
            <el-form-item label="加工费">
              <span style="color: #e6a23c; font-size: 18px; font-weight: bold">￥{{ preview.processing_fee || 0 }}</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submit" :loading="loading">登记取号</el-button>
              <el-button @click="reset">重置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
      <el-col :span="14">
        <el-card>
          <template #header>
            <div style="display:flex; justify-content: space-between; align-items: center">
              <span>今日排队</span>
              <el-button size="small" @click="loadList">刷新</el-button>
            </div>
          </template>
          <el-table :data="orders" stripe>
            <el-table-column prop="queue_number" label="排队号" width="100" />
            <el-table-column label="农户">
              <template #default="{ row }">
                {{ row.farmer.name }}
                <div v-if="row.farmer.phone" style="color:#999; font-size:12px">{{ row.farmer.phone }}</div>
              </template>
            </el-table-column>
            <el-table-column label="原料" width="80">
              <template #default="{ row }">
                {{ row.seed_type === 'peanut' ? '花生' : '菜籽' }}
              </template>
            </el-table-column>
            <el-table-column prop="gross_weight" label="毛重(kg)" width="90" />
            <el-table-column prop="processing_fee" label="加工费" width="90">
              <template #default="{ row }">￥{{ row.processing_fee }}</template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="statusTag(row.status)">{{ statusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="登记时间" width="170">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="100">
              <template #default="{ row }">
                <el-button size="small" type="primary" link @click="openPrintDialog(row)">打印小票</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog
      v-model="printDialogVisible"
      title="打印小票"
      :width="paperSize === 'a5' ? '620px' : '480px'"
      :close-on-click-modal="false"
      class="no-print-header"
    >
      <div class="paper-size-select no-print" style="margin-bottom: 16px;">
        <span style="margin-right: 12px;">纸张规格：</span>
        <el-radio-group v-model="paperSize">
          <el-radio-button value="80mm">80mm 小票纸</el-radio-button>
          <el-radio-button value="a5">A5 半页</el-radio-button>
        </el-radio-group>
      </div>

      <div style="border: 1px solid #e4e7ed; padding: 10px; background: #fafafa; border-radius: 4px;">
        <PrintReceipt :order-data="currentOrder" :paper-size="paperSize" ref="printReceiptRef" />
      </div>

      <template #footer>
        <div class="no-print">
          <el-button @click="printDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="doPrint">
            <el-icon><Printer /></el-icon>
            <span>打印小票</span>
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Printer } from '@element-plus/icons-vue'
import { createOrder, getOrderList, calculateFee, getOrder } from '../api'
import PrintReceipt from '../components/PrintReceipt.vue'

const formRef = ref()
const loading = ref(false)
const orders = ref([])
const preview = reactive({ expected_oil: 0, processing_fee: 0 })

const printDialogVisible = ref(false)
const printReceiptRef = ref(null)
const paperSize = ref('80mm')
const currentOrder = ref({})

const rules = {
  farmer_name: [
    { required: true, message: '请输入农户姓名', trigger: 'blur' },
    { min: 1, max: 20, message: '姓名长度在 1 到 20 个字符', trigger: 'blur' }
  ],
  seed_type: [
    { required: true, message: '请选择原料类型', trigger: 'change' }
  ],
  gross_weight: [
    { required: true, message: '请输入毛重', trigger: 'blur' },
    { type: 'number', min: 0.01, message: '毛重必须大于 0', trigger: 'blur' }
  ]
}

const form = reactive({
  farmer_name: '',
  farmer_phone: '',
  seed_type: '',
  gross_weight: 0,
  net_weight: 0,
  cake_taken: false,
  remark: ''
})

const calcFee = async () => {
  if (!form.seed_type || !form.gross_weight) {
    preview.expected_oil = 0
    preview.processing_fee = 0
    return
  }
  try {
    const res = await calculateFee({
      seed_type: form.seed_type,
      gross_weight: form.gross_weight,
      cake_taken: form.cake_taken
    })
    preview.expected_oil = res.data.expected_oil
    preview.processing_fee = res.data.processing_fee
  } catch (e) {}
}

const openPrintDialog = async (row) => {
  try {
    const res = await getOrder(row.id)
    currentOrder.value = res.data
  } catch (e) {
    currentOrder.value = { ...row }
  }
  paperSize.value = '80mm'
  printDialogVisible.value = true
}

const doPrint = async () => {
  await nextTick()
  setTimeout(() => {
    window.print()
  }, 100)
}

const submit = async () => {
  try {
    await formRef.value.validate()
  } catch (e) {
    ElMessage.warning('请完整填写必填项')
    return
  }
  loading.value = true
  try {
    const res = await createOrder(form)
    const newOrder = res.data
    ElMessage.success(`登记成功！排队号：${newOrder.queue_number}`)
    reset()
    loadList()

    try {
      await ElMessageBox.confirm(
        `登记成功！排队号：${newOrder.queue_number}\n是否打印小票？`,
        '打印小票',
        {
          confirmButtonText: '打印小票',
          cancelButtonText: '不打印',
          type: 'success',
          distinguishCancelAndClose: true
        }
      )
      try {
        const detailRes = await getOrder(newOrder.id)
        currentOrder.value = detailRes.data
      } catch (e) {
        currentOrder.value = newOrder
      }
      paperSize.value = '80mm'
      printDialogVisible.value = true
    } catch (e) {
    }
  } catch (e) {
    if (e === 'cancel' || e === 'close') return
    const msg = e?.response?.data?.error
    if (msg && msg.includes('required')) {
      ElMessage.error('必填项未填完整：' + msg)
    } else if (msg && msg.includes('gt=0')) {
      ElMessage.error('毛重必须大于 0 公斤')
    } else if (msg) {
      ElMessage.error(msg)
    } else {
      ElMessage.error('登记失败，请检查网络或稍后重试')
    }
  } finally {
    loading.value = false
  }
}

const reset = () => {
  formRef.value?.resetFields()
  Object.assign(form, { farmer_name: '', farmer_phone: '', seed_type: '', gross_weight: 0, net_weight: 0, cake_taken: false, remark: '' })
  preview.expected_oil = 0
  preview.processing_fee = 0
}

const loadList = async () => {
  try {
    const res = await getOrderList({ today_only: true })
    orders.value = res.data
  } catch (e) {}
}

const statusText = (s) => ({ waiting: '排队中', processing: '榨油中', completed: '已完成', cancelled: '已取消' }[s] || s)
const statusTag = (s) => ({ waiting: 'info', processing: 'warning', completed: 'success', cancelled: 'danger' }[s] || '')
const formatTime = (t) => t ? new Date(t).toLocaleString('zh-CN') : ''

onMounted(loadList)
</script>
