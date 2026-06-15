<template>
  <div>
    <h2>价格配置</h2>
    <el-alert
      type="info"
      :closable="false"
      style="margin-bottom: 20px"
      title="修改价格只会影响新登记的订单，已登记的订单不受影响。旺季临时调价无需担心已登记订单混乱。"
    />
    <el-card>
      <template #header>原料价格与参数</template>
      <el-table :data="pricingList" border>
        <el-table-column label="原料类型" width="150">
          <template #default="{ row }">
            <strong>{{ row.id === 'peanut' ? '花生' : '菜籽' }}</strong>
          </template>
        </el-table-column>
        <el-table-column label="原料单价(元/百公斤)" width="220">
          <template #default="{ row }">
            <el-input-number v-model="row.price_per_kg" :min="0" :precision="2" :step="10" />
          </template>
        </el-table-column>
        <el-table-column label="出油率(%)" width="180">
          <template #default="{ row }">
            <el-input-number v-model="row.oil_rate" :min="0" :max="100" :precision="1" :step="1" />
          </template>
        </el-table-column>
        <el-table-column label="基础加工费(元/批)" width="200">
          <template #default="{ row }">
            <el-input-number v-model="row.processing_fee" :min="0" :precision="2" :step="5" />
          </template>
        </el-table-column>
        <el-table-column label="饼粕带走加价(元/百公斤)" width="240">
          <template #default="{ row }">
            <el-input-number v-model="row.cake_take_fee" :min="0" :precision="2" :step="5" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" size="small" @click="save(row)">保存</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-card style="margin-top: 20px">
      <template #header>计费说明</template>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="加工费计算">
          原料毛重 × 原料单价 ÷ 100 + 基础加工费
          <span v-if="cakeTakenDemo" style="color:#e6a23c"> + 饼粕带走费（毛重 × 饼粕带走加价 ÷ 100）</span>
        </el-descriptions-item>
        <el-descriptions-item label="预计出油">原料毛重 × 出油率 ÷ 100</el-descriptions-item>
        <el-descriptions-item label="寄存费">每桶每天 0.5 元，超期自动累加</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getPricingList, updatePricing } from '../api'

const pricingList = ref([])
const cakeTakenDemo = ref(true)

const load = async () => {
  try {
    const res = await getPricingList()
    pricingList.value = res.data
  } catch (e) {}
}

const save = async (row) => {
  try {
    await updatePricing(row)
    ElMessage.success('保存成功，新登记的订单将使用此价格')
  } catch (e) {
    ElMessage.error(e.response?.data?.error || '保存失败')
  }
}

onMounted(load)
</script>
