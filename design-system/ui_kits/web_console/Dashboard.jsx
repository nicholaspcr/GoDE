// Dashboard — 3 quick-action cards + recent executions list
function Dashboard({ onNew, onViewAll, onViewExecution, executions }) {
  const total = executions.length;
  const recent = executions.slice(0, 5);
  const statusLabel = (s) => ({ running:'Running', pending:'Pending', completed:'Completed', failed:'Failed', cancelled:'Cancelled' }[s] || 'Unknown');
  return (
    <div>
      <div style={{ display:'grid', gridTemplateColumns:'repeat(3, 1fr)', gap:24 }}>
        <Card>
          <CardHeader><CardTitle>New Execution</CardTitle><CardDescription>Configure and run a new DE optimization</CardDescription></CardHeader>
          <CardContent><Button style={{ width:'100%' }} onClick={onNew}>Create Execution</Button></CardContent>
        </Card>
        <Card>
          <CardHeader><CardTitle>All Executions</CardTitle><CardDescription>View and manage your optimization runs</CardDescription></CardHeader>
          <CardContent><Button variant="outline" style={{ width:'100%' }} onClick={onViewAll}>View All</Button></CardContent>
        </Card>
        <Card>
          <CardHeader><CardTitle>Quick Stats</CardTitle><CardDescription>Your optimization activity</CardDescription></CardHeader>
          <CardContent>
            <div style={{ fontSize:30, fontWeight:700 }}>{total}</div>
            <p style={{ fontSize:14, color:'hsl(215.4 16.3% 46.9%)', margin:0 }}>Total executions</p>
          </CardContent>
        </Card>
      </div>
      <div style={{ marginTop:32 }}>
        <h2 style={{ margin:'0 0 16px', fontSize:20, fontWeight:600 }}>Recent Executions</h2>
        {recent.length === 0 ? (
          <Card><CardContent style={{ padding:32, textAlign:'center', color:'hsl(215.4 16.3% 46.9%)' }}>No executions yet. Create your first optimization run!</CardContent></Card>
        ) : (
          <div style={{ display:'flex', flexDirection:'column', gap:8 }}>
            {recent.map(ex => (
              <Card key={ex.id}>
                <CardContent style={{ padding:'16px 24px', display:'flex', alignItems:'center', justifyContent:'space-between' }}>
                  <div>
                    <div style={{ fontWeight:500, fontFamily:'JetBrains Mono, ui-monospace, monospace', fontSize:13 }}>{ex.id}</div>
                    <div style={{ fontSize:14, color:'hsl(215.4 16.3% 46.9%)' }}>Status: {statusLabel(ex.status)}</div>
                  </div>
                  <Button variant="outline" size="sm" onClick={()=>onViewExecution(ex.id)}>View</Button>
                </CardContent>
              </Card>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
Object.assign(window, { Dashboard });
