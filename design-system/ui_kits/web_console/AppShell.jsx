// AppShell — top header with wordmark + sign-out
function AppShell({ username, onLogout, children }) {
  return (
    <div style={{ minHeight:'100vh', background:'hsl(0 0% 100%)', color:'hsl(222.2 84% 4.9%)' }}>
      <header style={{ borderBottom:'1px solid hsl(214.3 31.8% 91.4%)' }}>
        <div style={{ maxWidth:1200, margin:'0 auto', padding:'16px 24px', display:'flex', alignItems:'center', justifyContent:'space-between' }}>
          <h1 style={{ margin:0, fontSize:24, fontWeight:700, letterSpacing:'-0.025em' }}>GoDE</h1>
          <div style={{ display:'flex', alignItems:'center', gap:16 }}>
            <span style={{ fontSize:14, color:'hsl(215.4 16.3% 46.9%)' }}>Welcome, {username}</span>
            <Button variant="outline" size="sm" onClick={onLogout}>Sign Out</Button>
          </div>
        </div>
      </header>
      <main style={{ maxWidth:1200, margin:'0 auto', padding:'32px 24px' }}>{children}</main>
    </div>
  );
}
Object.assign(window, { AppShell });
