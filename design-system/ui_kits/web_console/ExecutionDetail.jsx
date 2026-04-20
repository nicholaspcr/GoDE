// ExecutionDetail — info + config grids, progress if running, Pareto plot if completed
function ExecutionDetail({ execution, onBack, onCancel, onDelete }) {
  const status = execution.status;
  const isRunning = status === 'running' || status === 'pending';
  const isCompleted = status === 'completed';
  const canDelete = isCompleted || status === 'failed' || status === 'cancelled';
  const statusMap = {
    pending: { label:'Pending', variant:'secondary' },
    running: { label:'Running', variant:'default' },
    completed: { label:'Completed', variant:'outline' },
    failed: { label:'Failed', variant:'destructive' },
    cancelled: { label:'Cancelled', variant:'secondary' },
  };
  const s = statusMap[status] || { label:'Unknown', variant:'outline' };
  const [prog, setProg] = React.useState(isRunning ? 0 : 100);
  const [gen, setGen] = React.useState(0);
  React.useEffect(() => {
    if (!isRunning) return;
    const id = setInterval(() => {
      setProg(p => Math.min(100, p + 3));
      setGen(g => Math.min(execution.generations, g + Math.ceil(execution.generations/33)));
    }, 300);
    return () => clearInterval(id);
  }, [isRunning]);

  return (
    <div>
      <div style={{ marginBottom:24 }}>
        <a href="#" onClick={(e)=>{e.preventDefault();onBack();}} style={{ fontSize:14, color:'hsl(215.4 16.3% 46.9%)', textDecoration:'none' }}>&larr; Back to Executions</a>
      </div>
      <div style={{ display:'flex', alignItems:'center', justifyContent:'space-between', marginBottom:24 }}>
        <div style={{ display:'flex', alignItems:'center', gap:16 }}>
          <h1 style={{ margin:0, fontSize:24, fontWeight:700, letterSpacing:'-0.025em' }}>Execution Details</h1>
          <Badge variant={s.variant}>{s.label}</Badge>
        </div>
        <div style={{ display:'flex', gap:8 }}>
          {isRunning && <Button variant="outline" onClick={onCancel}>Cancel</Button>}
          {canDelete && <Button variant="destructive" onClick={onDelete}>Delete</Button>}
        </div>
      </div>

      <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:24 }}>
        <Card><CardContent style={{ padding:24 }}>
          <h2 style={{ margin:'0 0 16px', fontSize:18, fontWeight:600 }}>Information</h2>
          <dl style={{ margin:0, display:'flex', flexDirection:'column', gap:8, fontSize:14 }}>
            <DLRow k="ID" v={<span style={{ fontFamily:'JetBrains Mono, monospace' }}>{execution.id}</span>} />
            <DLRow k="Status" v={s.label} />
            <DLRow k="Problem" v={execution.problem} />
            <DLRow k="Algorithm" v={execution.algorithm} />
            <DLRow k="Variant" v={execution.variant} />
            <DLRow k="Created" v="Apr 20, 2026, 2:32:11 PM" />
          </dl>
        </CardContent></Card>
        <Card><CardContent style={{ padding:24 }}>
          <h2 style={{ margin:'0 0 16px', fontSize:18, fontWeight:600 }}>Configuration</h2>
          <dl style={{ margin:0, display:'grid', gridTemplateColumns:'1fr 1fr', gap:8, fontSize:14 }}>
            <DLRow k="Executions" v={execution.executions} />
            <DLRow k="Generations" v={execution.generations} />
            <DLRow k="Population" v={execution.populationSize} />
            <DLRow k="Dimensions" v={execution.dimensionsSize} />
            <DLRow k="Objectives" v={execution.objectivesSize} />
            <DLRow k="Floor" v={execution.floorLimiter} />
            <DLRow k="Ceiling" v={execution.ceilLimiter} />
            <DLRow k="CR" v={execution.cr} />
            <DLRow k="F" v={execution.f} />
            <DLRow k="P" v={execution.p} />
          </dl>
        </CardContent></Card>
      </div>

      {isRunning && (
        <Card style={{ marginTop:24 }}><CardContent style={{ padding:24 }}>
          <h2 style={{ margin:'0 0 16px', fontSize:18, fontWeight:600 }}>Progress</h2>
          <div style={{ display:'flex', flexDirection:'column', gap:16 }}>
            <div>
              <div style={{ display:'flex', justifyContent:'space-between', fontSize:14, marginBottom:8 }}>
                <span>Overall Progress</span><span>{Math.round(prog)}%</span>
              </div>
              <Progress value={prog} />
            </div>
            <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr', gap:16, fontSize:14 }}>
              <div><span style={{ color:'hsl(215.4 16.3% 46.9%)' }}>Generation: </span>{gen} / {execution.generations}</div>
              <div><span style={{ color:'hsl(215.4 16.3% 46.9%)' }}>Execution: </span>1 / {execution.executions}</div>
            </div>
          </div>
        </CardContent></Card>
      )}

      {isCompleted && (
        <div style={{ marginTop:24 }}>
          <h2 style={{ margin:'0 0 16px', fontSize:18, fontWeight:600 }}>Results</h2>
          <Card><CardContent style={{ padding:24 }}>
            <ParetoPlot problem={execution.problem} />
          </CardContent></Card>
        </div>
      )}
    </div>
  );
}

function DLRow({ k, v }) {
  return (
    <div style={{ display:'flex', justifyContent:'space-between', gap:12 }}>
      <dt style={{ color:'hsl(215.4 16.3% 46.9%)' }}>{k}</dt>
      <dd style={{ margin:0, textAlign:'right' }}>{v}</dd>
    </div>
  );
}
Object.assign(window, { ExecutionDetail });
