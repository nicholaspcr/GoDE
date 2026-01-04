import Plot from 'react-plotly.js'
import type { ApiV1Vector } from '@/api/generated'
import type { Data, Layout } from 'plotly.js'

interface ParetoPlot2DProps {
  vectors: ApiV1Vector[]
  xAxis: number
  yAxis: number
  title?: string
}

export function ParetoPlot2D({
  vectors,
  xAxis,
  yAxis,
  title = 'Pareto Front',
}: ParetoPlot2DProps) {
  const xValues = vectors.map((v) => v.objectives?.[xAxis] ?? 0)
  const yValues = vectors.map((v) => v.objectives?.[yAxis] ?? 0)
  const crowdingDistances = vectors.map((v) => v.crowdingDistance ?? 0)

  const data: Data[] = [
    {
      x: xValues,
      y: yValues,
      mode: 'markers',
      type: 'scatter',
      marker: {
        size: 8,
        color: crowdingDistances,
        colorscale: 'Viridis',
        showscale: true,
        colorbar: {
          title: { text: 'Crowding Distance', side: 'right' },
        },
      },
      hovertemplate:
        `Objective ${xAxis + 1}: %{x:.4f}<br>` +
        `Objective ${yAxis + 1}: %{y:.4f}<br>` +
        `<extra></extra>`,
    },
  ]

  const layout: Partial<Layout> = {
    title: {
      text: title,
      font: { size: 16 },
    },
    xaxis: {
      title: { text: `Objective ${xAxis + 1}` },
      gridcolor: 'rgba(128, 128, 128, 0.2)',
    },
    yaxis: {
      title: { text: `Objective ${yAxis + 1}` },
      gridcolor: 'rgba(128, 128, 128, 0.2)',
    },
    paper_bgcolor: 'transparent',
    plot_bgcolor: 'transparent',
    autosize: true,
    margin: { l: 60, r: 40, t: 60, b: 60 },
  }

  return (
    <Plot
      data={data}
      layout={layout}
      useResizeHandler
      className="w-full h-full min-h-[400px]"
      config={{
        displaylogo: false,
        modeBarButtonsToRemove: ['lasso2d', 'select2d'],
        toImageButtonOptions: {
          format: 'png',
          filename: 'pareto_front_2d',
          scale: 2,
        },
      }}
    />
  )
}
