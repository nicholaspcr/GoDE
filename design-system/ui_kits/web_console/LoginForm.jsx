// LoginForm — centered card with username/password
function LoginForm({ onSuccess }) {
  const [username, setUsername] = React.useState('researcher');
  const [password, setPassword] = React.useState('••••••••');
  const [pending, setPending] = React.useState(false);
  const submit = (e) => {
    e.preventDefault();
    setPending(true);
    setTimeout(() => { setPending(false); onSuccess?.(username); }, 400);
  };
  return (
    <div style={{ minHeight:'100vh', display:'flex', alignItems:'center', justifyContent:'center', background:'white', padding:16 }}>
      <div style={{ width:'100%', maxWidth:420, display:'flex', flexDirection:'column', gap:24 }}>
        <div style={{ textAlign:'center' }}>
          <h1 style={{ margin:0, fontSize:30, fontWeight:700, letterSpacing:'-0.025em' }}>GoDE</h1>
          <p style={{ margin:'6px 0 0', color:'hsl(215.4 16.3% 46.9%)', fontSize:14 }}>Differential Evolution Optimization Framework</p>
        </div>
        <Card>
          <CardHeader>
            <CardTitle>Login</CardTitle>
            <CardDescription>Enter your credentials to access your account</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={submit} style={{ display:'flex', flexDirection:'column', gap:16 }}>
              <Field label="Username"><Input value={username} onChange={e=>setUsername(e.target.value)} placeholder="Enter your username" /></Field>
              <Field label="Password"><Input type="password" value={password} onChange={e=>setPassword(e.target.value)} placeholder="Enter your password" /></Field>
              <Button type="submit" style={{ width:'100%' }} disabled={pending}>{pending ? 'Signing in...' : 'Sign In'}</Button>
              <p style={{ textAlign:'center', fontSize:14, color:'hsl(215.4 16.3% 46.9%)', margin:0 }}>
                Don't have an account? <a href="#" style={{ color:'hsl(222.2 47.4% 11.2%)', textUnderlineOffset:'4px' }}>Register</a>
              </p>
            </form>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
Object.assign(window, { LoginForm });
