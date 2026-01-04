import { Label, Select } from '@/components/ui'

interface AxisSelectorProps {
  objectivesCount: number
  xAxis: number
  yAxis: number
  zAxis?: number
  onXAxisChange: (value: number) => void
  onYAxisChange: (value: number) => void
  onZAxisChange?: (value: number) => void
  show3D?: boolean
}

export function AxisSelector({
  objectivesCount,
  xAxis,
  yAxis,
  zAxis,
  onXAxisChange,
  onYAxisChange,
  onZAxisChange,
  show3D = false,
}: AxisSelectorProps) {
  const options = Array.from({ length: objectivesCount }, (_, i) => ({
    value: String(i),
    label: `Objective ${i + 1}`,
  }))

  return (
    <div className="flex flex-wrap gap-4 items-end">
      <div className="space-y-1">
        <Label htmlFor="x-axis" className="text-sm">
          X Axis
        </Label>
        <Select
          id="x-axis"
          value={String(xAxis)}
          onChange={(e) => onXAxisChange(Number(e.target.value))}
          className="w-36"
        >
          {options.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </Select>
      </div>

      <div className="space-y-1">
        <Label htmlFor="y-axis" className="text-sm">
          Y Axis
        </Label>
        <Select
          id="y-axis"
          value={String(yAxis)}
          onChange={(e) => onYAxisChange(Number(e.target.value))}
          className="w-36"
        >
          {options.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </Select>
      </div>

      {show3D && onZAxisChange && zAxis !== undefined && (
        <div className="space-y-1">
          <Label htmlFor="z-axis" className="text-sm">
            Z Axis
          </Label>
          <Select
            id="z-axis"
            value={String(zAxis)}
            onChange={(e) => onZAxisChange(Number(e.target.value))}
            className="w-36"
          >
            {options.map((opt) => (
              <option key={opt.value} value={opt.value}>
                {opt.label}
              </option>
            ))}
          </Select>
        </div>
      )}
    </div>
  )
}
