import { useState } from 'react'
import { Button } from '@/components/ui'
import type { ApiV1Vector } from '@/api/generated'

interface ObjectiveTableProps {
  vectors: ApiV1Vector[]
  pageSize?: number
}

export function ObjectiveTable({ vectors, pageSize = 20 }: ObjectiveTableProps) {
  const [page, setPage] = useState(0)
  const [sortColumn, setSortColumn] = useState<number | 'crowding'>(-1)
  const [sortAsc, setSortAsc] = useState(true)

  const objectivesCount = vectors[0]?.objectives?.length ?? 0

  const sortedVectors = [...vectors].sort((a, b) => {
    if (sortColumn === -1) return 0
    if (sortColumn === 'crowding') {
      const aVal = a.crowdingDistance ?? 0
      const bVal = b.crowdingDistance ?? 0
      return sortAsc ? aVal - bVal : bVal - aVal
    }
    const aVal = a.objectives?.[sortColumn] ?? 0
    const bVal = b.objectives?.[sortColumn] ?? 0
    return sortAsc ? aVal - bVal : bVal - aVal
  })

  const totalPages = Math.ceil(vectors.length / pageSize)
  const paginatedVectors = sortedVectors.slice(
    page * pageSize,
    (page + 1) * pageSize
  )

  const handleSort = (column: number | 'crowding') => {
    if (sortColumn === column) {
      setSortAsc(!sortAsc)
    } else {
      setSortColumn(column)
      setSortAsc(true)
    }
  }

  const exportCSV = () => {
    const headers = [
      '#',
      ...Array.from({ length: objectivesCount }, (_, i) => `Objective ${i + 1}`),
      'Crowding Distance',
    ]
    const rows = vectors.map((v, idx) => [
      idx + 1,
      ...(v.objectives?.map((o) => o.toFixed(6)) ?? []),
      v.crowdingDistance?.toFixed(6) ?? '',
    ])
    const csv = [headers.join(','), ...rows.map((r) => r.join(','))].join('\n')
    const blob = new Blob([csv], { type: 'text/csv' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'pareto_solutions.csv'
    a.click()
    URL.revokeObjectURL(url)
  }

  return (
    <div className="space-y-4">
      <div className="flex justify-between items-center">
        <p className="text-sm text-muted-foreground">
          {vectors.length} Pareto-optimal solutions
        </p>
        <Button variant="outline" size="sm" onClick={exportCSV}>
          Export CSV
        </Button>
      </div>

      <div className="overflow-x-auto rounded-md border">
        <table className="w-full text-sm">
          <thead className="bg-muted/50">
            <tr>
              <th className="text-left py-3 px-3 font-medium">#</th>
              {Array.from({ length: objectivesCount }, (_, i) => (
                <th
                  key={i}
                  className="text-left py-3 px-3 font-medium cursor-pointer hover:bg-muted"
                  onClick={() => handleSort(i)}
                >
                  Objective {i + 1}
                  {sortColumn === i && (
                    <span className="ml-1">{sortAsc ? '↑' : '↓'}</span>
                  )}
                </th>
              ))}
              <th
                className="text-left py-3 px-3 font-medium cursor-pointer hover:bg-muted"
                onClick={() => handleSort('crowding')}
              >
                Crowding
                {sortColumn === 'crowding' && (
                  <span className="ml-1">{sortAsc ? '↑' : '↓'}</span>
                )}
              </th>
            </tr>
          </thead>
          <tbody>
            {paginatedVectors.map((vector, idx) => (
              <tr key={idx} className="border-t hover:bg-muted/30">
                <td className="py-2 px-3 text-muted-foreground">
                  {page * pageSize + idx + 1}
                </td>
                {vector.objectives?.map((obj, objIdx) => (
                  <td key={objIdx} className="py-2 px-3 font-mono text-xs">
                    {obj.toFixed(6)}
                  </td>
                ))}
                <td className="py-2 px-3 font-mono text-xs">
                  {vector.crowdingDistance?.toFixed(4) ?? '-'}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {totalPages > 1 && (
        <div className="flex justify-center items-center gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            disabled={page === 0}
          >
            Previous
          </Button>
          <span className="text-sm text-muted-foreground">
            Page {page + 1} of {totalPages}
          </span>
          <Button
            variant="outline"
            size="sm"
            onClick={() => setPage((p) => Math.min(totalPages - 1, p + 1))}
            disabled={page === totalPages - 1}
          >
            Next
          </Button>
        </div>
      )}
    </div>
  )
}
