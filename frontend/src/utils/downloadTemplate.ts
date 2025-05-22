import * as XLSX from 'xlsx'

export function downloadTemplate() {
  const wsData = [
    ['NO', 'Label Rekonsiliasi', '2019-08-31', '2019-09-30'],
    [1, 'Laba Sebelum Pajak', 123, 456],
    [1, 'Laba Sesudah Pajak', 78910, 11234]
  ]

  const ws = XLSX.utils.aoa_to_sheet(wsData)
  const wb = XLSX.utils.book_new()
  XLSX.utils.book_append_sheet(wb, ws, 'Template')

  const wbout = XLSX.write(wb, { bookType: 'xlsx', type: 'array' })

  const blob = new Blob([wbout], { type: 'application/octet-stream' })
  const url = window.URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'template-laba.xlsx'
  a.click()
  window.URL.revokeObjectURL(url)
}
