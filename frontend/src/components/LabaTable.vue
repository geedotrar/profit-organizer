<template>
  <div>
    <h2 class="mb-3">Laba</h2>

    <div v-if="loading" class="alert alert-success">Loading...</div>
    <div v-if="error" class="alert alert-danger">{{ error }}</div>

    <div class="d-flex justify-content-end gap-2 mb-3">
      <input 
        type="file" 
        ref="fileInput" 
        @change="handleFileChange" 
        accept=".xlsx,.xls"
        class="form-control"
        style="max-width: 300px;"
      />
      <button 
      class="btn btn-primary cursor-pointer"
      :disabled="loadingImport"
        @click="uploadFile"
      >
        <span v-if="loadingImport" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
        <span v-if="!loadingImport">Import Data</span>
        <span v-else>Importing...</span>
      </button>

      <button 
        class="btn btn-info"
        @click="downloadTemplate"
        :disabled="loadingExport || loadingImport"
      >
        Download Template
      </button>

      <button 
        @click="exportData"
        :disabled="loadingExport"
        style="background-color: white; color: black; border: 1px solid #ccc;"
        class="btn">
        <span v-if="loadingExport" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
        <span v-if="!loadingExport">Export Data</span>
        <span v-else>Exporting...</span>
      </button>
    </div>

    <div style="overflow-x: auto;">
      <table v-if="labas.length" class="table table-bordered table-striped text-center" style="font-size: 1.25rem; min-width: 700px;">
        <thead class="table-success">
          <tr>
            <th style="width: 60px;">NO</th>
            <th style="min-width: 200px;">Label Rekonsiliasi</th>
            <th v-for="periode in sortedPeriods" :key="periode" style="min-width: 140px;">
              {{ periode }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(item, index) in labas" :key="item.label_rekonsiliasi_fiskal" class="align-middle">
            <td>{{ index + 1 }}</td>
            <td class="text-start">{{ item.label_rekonsiliasi_fiskal }}</td>
            <td v-for="periode in sortedPeriods" :key="periode">
              {{ formatNumber(item.periode[periode]) }}
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else-if="!loading && !error" class="alert alert-warning">
        Tidak ada data
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'
import { useLabaStore } from '../stores/laba'
import { importLabaFile } from '../stores/importLaba'
import { exportLabaFile } from '../stores/exportLaba'
import { downloadTemplate } from '../utils/downloadTemplate'

const store = useLabaStore()

const selectedFile = ref<File | null>(null)
const fileInput = ref<HTMLInputElement | null>(null)

onMounted(() => {
  store.fetchAll()
})

const loading = computed(() => store.loading)
const error = computed(() => store.error)
const loadingImport = computed(() => store.loadingImport)
const loadingExport = computed(() => store.loadingExport)

const sortedPeriods = computed(() => {
  const set = new Set<string>()
  Object.values(store.labas).forEach(p => {
    Object.keys(p).forEach(key => set.add(key))
  })
  return Array.from(set).sort()
})

const labas = computed(() =>
  Object.entries(store.labas).map(([label, periode]) => ({
    label_rekonsiliasi_fiskal: label,
    periode
  }))
)

function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  selectedFile.value = target.files?.[0] ?? null
}

async function uploadFile() {
  if (!selectedFile.value) {
    alert('choose file first')
    return
  }

  try {
    await importLabaFile(selectedFile.value)
    alert('Success Import New Data!')
    selectedFile.value = null
    if (fileInput.value) fileInput.value.value = ''
  } catch (err: any) {
    alert(err.message)
  }
}

async function exportData() {
  try {
    await exportLabaFile()
  } catch (err: any) {
    alert(err.message)
  }
}

function formatNumber(value: number | undefined) {
  if (value === undefined || value === null) return '0'
  return new Intl.NumberFormat('en-US').format(value)
}
</script>
