import axios from 'axios'

export async function exportLabaFile() {
  try {
    const response = await axios.get('http://localhost:8080/laba/export', {
      responseType: 'blob',
    })

    const blob = new Blob([response.data])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')

    let fileName = 'laba_export.xlsx'
    const disposition = response.headers['content-disposition']
    const match = disposition?.match(/filename="?(.+)"?/)
    if (match && match[1]) fileName = match[1]

    link.href = url
    link.setAttribute('download', fileName)
    document.body.appendChild(link)
    link.click()
    link.remove()
    window.URL.revokeObjectURL(url)
  } catch (error: any) {
    if (axios.isAxiosError(error)) {
      const reader = new FileReader()
      const blob = error.response?.data

      return new Promise((_, reject) => {
        reader.onload = () => {
          try {
            const json = JSON.parse(reader.result as string)
            reject(new Error(json?.message || 'Failed export data'))
          } catch {
            reject(new Error('failed export data'))
          }
        }

        if (blob instanceof Blob) {
          reader.readAsText(blob)
        } else {
          reject(new Error('failed export file'))
        }
      })
    }

    throw new Error('failed export file')
  }
}
