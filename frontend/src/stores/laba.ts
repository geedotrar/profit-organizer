import { defineStore } from 'pinia'
import axios from 'axios'

export const useLabaStore = defineStore('laba', {
  state: () => ({
    labas: {} as Record<string, Record<string, number>>,
    loading: false,
    error: null as string | null,
    loadingImport: false,
    loadingExport: false,
  }),
  actions: {
    async fetchAll() {
      this.loading = true
      this.error = null
      try {
        const res = await axios.get('http://localhost:8080/laba')
        if (res.data.status === 200 && !res.data.error) {
          this.labas = res.data.data
        } else {
          this.error = res.data.message
        }
      } catch (err: any) {
        this.error = err.message
      } finally {
        this.loading = false
      }
    }
  }
})
