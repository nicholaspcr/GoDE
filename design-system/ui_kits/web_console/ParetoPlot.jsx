// ParetoPlot — inline SVG approximation of Plotly's 2D scatter w/ Viridis crowding color
function ParetoPlot({ problem='zdt1', width=720, height=420 }) {
  const pts = React.useMemo(() => {
    // Synthesize a curve resembling ZDT1 pareto front f2 = 1 - sqrt(f1)
    const out = [];
    for (let i = 0; i < 60; i++) {
      const f1 = i / 59;
      const noise = (Math.random() - 0.5) * 0.04;
      const f2 = 1 - Math.sqrt(f1) + noise;
      out.push({ f1, f2, cd: Math.random() });
    }
    return out;
  }, [problem]);

  const pad = { l:60, r:80, t:40, b:50 };
  const iw = width - pad.l - pad.r, ih = height - pad.t - pad.b;
  const sx = (x) => pad.l + x * iw;
  const sy = (y) => pad.t + (1 - y) * ih;

  // Viridis stops
  const vir = ['#440154','#3b528b','#21918c','#5ec962','#fde725'];
  const mix = (a,b,t) => {
    const ah = parseInt(a.slice(1),16), bh = parseInt(b.slice(1),16);
    const ar=(ah>>16)&255, ag=(ah>>8)&255, ab=ah&255;
    const br=(bh>>16)&255, bg=(bh>>8)&255, bb=bh&255;
    const r=Math.round(ar+(br-ar)*t), g=Math.round(ag+(bg-ag)*t), bl=Math.round(ab+(bb-ab)*t);
    return `rgb(${r},${g},${bl})`;
  };
  const viridis = (t) => {
    const s = Math.max(0, Math.min(0.9999, t)) * (vir.length - 1);
    const i = Math.floor(s);
    return mix(vir[i], vir[i+1], s - i);
  };

  const gridLines = [0, 0.25, 0.5, 0.75, 1];
  return (
    <svg viewBox={`0 0 ${width} ${height}`} style={{ width:'100%', height:'auto' }} xmlns="http://www.w3.org/2000/svg">
      <rect x="0" y="0" width={width} height={height} fill="transparent"/>
      {/* grid */}
      {gridLines.map(g => (
        <g key={g}>
          <line x1={sx(g)} x2={sx(g)} y1={pad.t} y2={pad.t+ih} stroke="rgba(128,128,128,0.2)" />
          <line y1={sy(g)} y2={sy(g)} x1={pad.l} x2={pad.l+iw} stroke="rgba(128,128,128,0.2)" />
          <text x={sx(g)} y={pad.t+ih+18} fontSize="11" fill="hsl(215.4 16.3% 46.9%)" textAnchor="middle" fontFamily="Inter, system-ui">{g.toFixed(2)}</text>
          <text x={pad.l-10} y={sy(g)+4} fontSize="11" fill="hsl(215.4 16.3% 46.9%)" textAnchor="end" fontFamily="Inter, system-ui">{g.toFixed(2)}</text>
        </g>
      ))}
      {/* axis lines */}
      <line x1={pad.l} x2={pad.l+iw} y1={pad.t+ih} y2={pad.t+ih} stroke="hsl(222.2 84% 4.9%)" strokeWidth="1"/>
      <line x1={pad.l} x2={pad.l} y1={pad.t} y2={pad.t+ih} stroke="hsl(222.2 84% 4.9%)" strokeWidth="1"/>
      {/* title */}
      <text x={width/2} y={24} fontSize="16" fontWeight="600" fill="hsl(222.2 84% 4.9%)" textAnchor="middle" fontFamily="Inter, system-ui">Pareto Front</text>
      {/* axis labels */}
      <text x={pad.l+iw/2} y={height-10} fontSize="12" fill="hsl(222.2 84% 4.9%)" textAnchor="middle" fontFamily="Inter, system-ui">Objective 1</text>
      <text transform={`rotate(-90 18 ${pad.t+ih/2})`} x={18} y={pad.t+ih/2} fontSize="12" fill="hsl(222.2 84% 4.9%)" textAnchor="middle" fontFamily="Inter, system-ui">Objective 2</text>
      {/* points */}
      {pts.map((p,i) => (
        <circle key={i} cx={sx(p.f1)} cy={sy(p.f2)} r="5" fill={viridis(p.cd)} stroke="white" strokeWidth="0.5"/>
      ))}
      {/* colorbar */}
      <defs>
        <linearGradient id="vir" x1="0" y1="1" x2="0" y2="0">
          {vir.map((c,i)=><stop key={i} offset={i/(vir.length-1)} stopColor={c}/>)}
        </linearGradient>
      </defs>
      <rect x={width-pad.r+24} y={pad.t} width="14" height={ih} fill="url(#vir)" stroke="hsl(214.3 31.8% 91.4%)"/>
      <text x={width-pad.r+31} y={pad.t-8} fontSize="10" fill="hsl(215.4 16.3% 46.9%)" textAnchor="middle" fontFamily="Inter, system-ui">CD</text>
    </svg>
  );
}
Object.assign(window, { ParetoPlot });
