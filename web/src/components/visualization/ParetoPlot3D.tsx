import Plot from 'react-plotly.js'
import type { ApiV1Vector } from '@/api/generated'
import type { Data, Layout } from 'plotly.js'

interface ParetoPlot3DProps {
  vectors: ApiV1Vector[]
  xAxis: number
  yAxis: number
  zAxis: number
  title?: string
}

export function ParetoPlot3D({
  vectors,
  xAxis,
  yAxis,
  zAxis,
  title = 'Pareto Front (3D)',
}: ParetoPlot3DProps) {
  const xValues = vectors.map((v) => v.objectives?.[xAxis] ?? 0)
  const yValues = vectors.map((v) => v.objectives?.[yAxis] ?? 0)
  const zValues = vectors.map((v) => v.objectives?.[zAxis] ?? 0)
  const crowdingDistances = vectors.map((v) => v.crowdingDistance ?? 0)

  const data: Data[] = [
    {
      x: xValues,
      y: yValues,
      z: zValues,
      mode: 'markers',
      type: 'scatter3d',
      marker: {
        size: 4,
        color: crowdingDistances,
        colorscale: 'Viridis',
        showscale: true,
        colorbar: {
          title: { text: 'Crowding Distance', side: 'right' },
        },
      },
      hovertemplate:
        `Obj ${xAxis + 1}: %{x:.4f}<br>` +
        `Obj ${yAxis + 1}: %{y:.4f}<br>` +
        `Obj ${zAxis + 1}: %{z:.4f}<br>` +
        `<extra></extra>`,
    },
  ]

  const layout: Partial<Layout> = {
    title: {
      text: title,
      font: { size: 16 },
    },
    scene: {
      xaxis: { title: { text: `Objective ${xAxis + 1}` } },
      yaxis: { title: { text: `Objective ${yAxis + 1}` } },
      zaxis: { title: { text: `Objective ${zAxis + 1}` } },
      camera: {
        eye: { x: 1.5, y: 1.5, z: 1.5 },
      },
    },
    paper_bgcolor: 'transparent',
    autosize: true,
    margin: { l: 0, r: 0, t: 60, b: 0 },
  }

  return (
    <Plot
      data={data}
      layout={layout}
      useResizeHandler
      className="w-full h-full min-h-[500px]"
      config={{
        displaylogo: false,
        toImageButtonOptions: {
          format: 'png',
          filename: 'pareto_front_3d',
          scale: 2,
        },
      }}
    />
  )
}
