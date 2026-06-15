<template>
  <el-container style="height: 100vh">
    <el-aside width="200px" style="background: #001529">
      <div style="color: #fff; padding: 20px; font-size: 18px; font-weight: bold; text-align: center; border-bottom: 1px solid #1f3a5c">
        榨油坊柜台
      </div>
      <el-menu
        :default-active="activeMenu"
        background-color="#001529"
        text-color="#b9c4cf"
        active-text-color="#409eff"
        router
      >
        <el-menu-item index="/queue">
          <el-icon><Tickets /></el-icon>
          <span>排队登记</span>
        </el-menu-item>
        <el-menu-item index="/call">
          <el-icon><Bell /></el-icon>
          <span>叫号/榨油</span>
        </el-menu-item>
        <el-menu-item index="/pickup">
          <el-icon><Goods /></el-icon>
          <span>取油管理</span>
        </el-menu-item>
        <el-menu-item index="/storage">
          <el-icon><Box /></el-icon>
          <span>油桶寄存</span>
        </el-menu-item>
        <el-menu-item index="/pricing">
          <el-icon><Money /></el-icon>
          <span>价格配置</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-container>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Tickets, Bell, Goods, Box, Money } from '@element-plus/icons-vue'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
}
.el-menu-item {
  height: 56px !important;
  line-height: 56px !important;
}
</style>
