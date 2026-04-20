// ExecutionForm — stacked cards: Algorithm, Parameters, GDE3
const ALGOS = ['gde3'];
const VARIANTS = ['rand/1','rand/2','best/1','best/2','pbest','current-to-best/1'];
const PROBLEMS = ['zdt1','zdt2','zdt3','zdt4','zdt6','dtlz1','dtlz2','dtlz3','dtlz4','wfg1','wfg2','wfg3','wfg4','wfg5'];

function ExecutionForm({ onCancel, onStart }) {
  const [f, setF] = React.useState({ algorithm:'gde3', variant:'rand/1', problem:'zdt1', executions:1, generations:100, populationSize:100, dimensionsSize:30, objectivesSize:2, floorLimiter:0, ceilLimiter:1, cr:0.9, f:0.5, p:0.1 });
  const set = (k) => (e) => setF({...f, [k]: e.target.type === 'number' ? Number(e.target.value) : e.target.value });
  const submit = (e) => { e.preventDefault(); onStart({ ...f, id: 'exec-' + Math.random().toString(36).slice(2,10) }); };
  return (
    <form onSubmit={submit} style={{ display:'flex', flexDirection:'column', gap:24 }}>
      <Card>
        <CardContent style={{ padding:24 }}>
          <h3 style={{ margin:'0 0 16px', fontSize:18, fontWeight:600 }}>Algorithm Configuration</h3>
          <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr 1fr', gap:16 }}>
            <Field label="Algorithm"><Select value={f.algorithm} onChange={set('algorithm')}>{ALGOS.map(a=><option key={a}>{a}</option>)}</Select></Field>
            <Field label="Variant"><Select value={f.variant} onChange={set('variant')}>{VARIANTS.map(v=><option key={v}>{v}</option>)}</Select></Field>
            <Field label="Problem"><Select value={f.problem} onChange={set('problem')}>{PROBLEMS.map(p=><option key={p}>{p}</option>)}</Select></Field>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent style={{ padding:24 }}>
          <h3 style={{ margin:'0 0 16px', fontSize:18, fontWeight:600 }}>Execution Parameters</h3>
          <div style={{ display:'grid', gridTemplateColumns:'repeat(4, 1fr)', gap:16 }}>
            <Field label="Executions"><Input type="number" value={f.executions} onChange={set('executions')} min={1} /></Field>
            <Field label="Generations"><Input type="number" value={f.generations} onChange={set('generations')} min={1} /></Field>
            <Field label="Population Size"><Input type="number" value={f.populationSize} onChange={set('populationSize')} min={1} /></Field>
            <Field label="Dimensions"><Input type="number" value={f.dimensionsSize} onChange={set('dimensionsSize')} min={1} /></Field>
            <Field label="Objectives"><Input type="number" value={f.objectivesSize} onChange={set('objectivesSize')} min={2} /></Field>
            <Field label="Floor Limiter"><Input type="number" step="0.1" value={f.floorLimiter} onChange={set('floorLimiter')} /></Field>
            <Field label="Ceil Limiter"><Input type="number" step="0.1" value={f.ceilLimiter} onChange={set('ceilLimiter')} /></Field>
          </div>
        </CardContent>
      </Card>
      <Card>
        <CardContent style={{ padding:24 }}>
          <h3 style={{ margin:'0 0 16px', fontSize:18, fontWeight:600 }}>GDE3 Parameters</h3>
          <div style={{ display:'grid', gridTemplateColumns:'1fr 1fr 1fr', gap:16 }}>
            <Field label="CR (Crossover Rate)"><Input type="number" step="0.01" value={f.cr} onChange={set('cr')} min={0} max={1} /></Field>
            <Field label="F (Scaling Factor)"><Input type="number" step="0.01" value={f.f} onChange={set('f')} min={0} max={2} /></Field>
            <Field label="P (Selection Parameter)"><Input type="number" step="0.01" value={f.p} onChange={set('p')} min={0} max={1} /></Field>
          </div>
        </CardContent>
      </Card>
      <div style={{ display:'flex', justifyContent:'flex-end', gap:16 }}>
        <Button type="button" variant="outline" onClick={onCancel}>Cancel</Button>
        <Button type="submit">Start Execution</Button>
      </div>
    </form>
  );
}
Object.assign(window, { ExecutionForm });
