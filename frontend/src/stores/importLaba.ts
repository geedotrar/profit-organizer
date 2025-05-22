import axios from 'axios'
import { useLabaStore } from './laba'

export async function importLabaFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)

  const store = useLabaStore()

  try {
     await axios.post('http://localhost:8080/laba/import', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })

    await store.fetchAll()
  } catch (error: any) {
    if (axios.isAxiosError(error)) {
      const apiError = error.response?.data
      throw new Error(apiError?.message || 'failed import data')
    }
    throw new Error('failed import file')
  }
}
