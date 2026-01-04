import { useState } from 'react'
import { Card, Button } from '@/components/ui'
import { ParetoPlot2D } from './ParetoPlot2D'
import { ParetoPlot3D } from './ParetoPlot3D'
import { AxisSelector } from './AxisSelector'
import { ObjectiveTable } from './ObjectiveTable'
import type { ApiV1Vector } from '@/api/generated'

type ViewMode = '2d' | '3d' | 'table'

interface ParetoVisualizationProps {
  vectors: ApiV1Vector[]
}

export function ParetoVisualization({ vectors }: ParetoVisualizationProps) {
  const objectivesCount = vectors[0]?.objectives?.length ?? 2

  const [viewMode, setViewMode] = useState<ViewMode>('2d')
  const [xAxis, setXAxis] = useState(0)
  const [yAxis, setYAxis] = useState(Math.min(1, objectivesCount - 1))
  const [zAxis, setZAxis] = useState(Math.min(2, objectivesCount - 1))

  if (vectors.length === 0) {
    return (
      <Card className="p-6 text-center text-muted-foreground">
        No Pareto-optimal solutions to display
      </Card>
    )
  }

  const can3D = objectivesCount >= 3

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap justify-between items-center gap-4">
        <div className="flex gap-2">
          <Button
            variant={viewMode === '2d' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setViewMode('2d')}
          >
            2D Plot
          </Button>
          {can3D && (
            <Button
              variant={viewMode === '3d' ? 'default' : 'outline'}
              size="sm"
              onClick={() => setViewMode('3d')}
            >
              3D Plot
            </Button>
          )}
          <Button
            variant={viewMode === 'table' ? 'default' : 'outline'}
            size="sm"
            onClick={() => setViewMode('table')}
          >
            Table
          </Button>
        </div>

        {(viewMode === '2d' || viewMode === '3d') && (
          <AxisSelector
            objectivesCount={objectivesCount}
            xAxis={xAxis}
            yAxis={yAxis}
            zAxis={zAxis}
            onXAxisChange={setXAxis}
            onYAxisChange={setYAxis}
            onZAxisChange={setZAxis}
            show3D={viewMode === '3d'}
          />
        )}
      </div>

      <Card className="p-4">
        {viewMode === '2d' && (
          <ParetoPlot2D vectors={vectors} xAxis={xAxis} yAxis={yAxis} />
        )}
        {viewMode === '3d' && (
          <ParetoPlot3D
            vectors={vectors}
            xAxis={xAxis}
            yAxis={yAxis}
            zAxis={zAxis}
          />
        )}
        {viewMode === 'table' && <ObjectiveTable vectors={vectors} />}
      </Card>
    </div>
  )
}
