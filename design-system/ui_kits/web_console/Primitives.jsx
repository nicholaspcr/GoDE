// Primitives — Button, Badge, Card, Input, Select, Label, Progress
const { forwardRef } = React;

const cls = (...a) => a.filter(Boolean).join(' ');

const btnStyles = {
  base: { display:'inline-flex', alignItems:'center', justifyContent:'center', whiteSpace:'nowrap', borderRadius:'6px', fontSize:14, fontWeight:500, fontFamily:'inherit', cursor:'pointer', transition:'background-color 150ms cubic-bezier(.4,0,.2,1), color 150ms', border:'1px solid transparent' },
  sizes: { default:{ height:40, padding:'0 16px'}, sm:{ height:36, padding:'0 12px'}, lg:{ height:44, padding:'0 32px'}, icon:{ height:40, width:40, padding:0 } },
  variants: {
    default:{ background:'hsl(222.2 47.4% 11.2%)', color:'hsl(210 40% 98%)'},
    destructive:{ background:'hsl(0 84.2% 60.2%)', color:'white'},
    outline:{ background:'white', color:'hsl(222.2 84% 4.9%)', borderColor:'hsl(214.3 31.8% 91.4%)'},
    secondary:{ background:'hsl(210 40% 96.1%)', color:'hsl(222.2 47.4% 11.2%)'},
    ghost:{ background:'transparent', color:'hsl(222.2 84% 4.9%)'},
    link:{ background:'transparent', color:'hsl(222.2 47.4% 11.2%)', textDecoration:'underline', textUnderlineOffset:'4px', padding:0, height:'auto'},
  }
};

function Button({ variant='default', size='default', children, style, disabled, onClick, type='button', asChild, ...rest }) {
  const [hover, setHover] = React.useState(false);
  const v = btnStyles.variants[variant];
  const hv = hover && !disabled ? (
    variant === 'default' ? { background:'hsl(222.2 47.4% 18%)'} :
    variant === 'destructive' ? { background:'hsl(0 84% 55%)' } :
    variant === 'outline' || variant === 'ghost' ? { background:'hsl(210 40% 96.1%)' } :
    variant === 'secondary' ? { background:'hsl(210 40% 92%)' } : {}
  ) : {};
  const merged = { ...btnStyles.base, ...btnStyles.sizes[size], ...v, ...hv, opacity: disabled?0.5:1, pointerEvents: disabled?'none':'auto', ...style };
  return <button type={type} style={merged} onMouseEnter={()=>setHover(true)} onMouseLeave={()=>setHover(false)} onClick={onClick} {...rest}>{children}</button>;
}

function Badge({ variant='default', children, style }) {
  const vs = {
    default:{ background:'hsl(222.2 47.4% 11.2%)', color:'white'},
    secondary:{ background:'hsl(210 40% 96.1%)', color:'hsl(222.2 47.4% 11.2%)'},
    destructive:{ background:'hsl(0 84.2% 60.2%)', color:'white'},
    outline:{ color:'hsl(222.2 84% 4.9%)', border:'1px solid hsl(214.3 31.8% 91.4%)'},
    success:{ background:'hsl(142 71% 45%)', color:'white'},
    warning:{ background:'hsl(48 96% 53%)', color:'#111'},
  }[variant] || {};
  return <span style={{ display:'inline-flex', alignItems:'center', borderRadius:9999, padding:'2px 10px', fontSize:12, fontWeight:600, border:'1px solid transparent', ...vs, ...style }}>{children}</span>;
}

function Card({ children, style, className }) {
  return <div className={className} style={{ border:'1px solid hsl(214.3 31.8% 91.4%)', background:'white', color:'hsl(222.2 84% 4.9%)', borderRadius:8, boxShadow:'0 1px 2px 0 rgb(0 0 0 / .05)', ...style }}>{children}</div>;
}
function CardHeader({ children, style }) { return <div style={{ padding:'24px 24px 0', display:'flex', flexDirection:'column', gap:6, ...style }}>{children}</div>; }
function CardTitle({ children, style }) { return <div style={{ fontSize:24, fontWeight:600, letterSpacing:'-0.025em', lineHeight:1, ...style }}>{children}</div>; }
function CardDescription({ children, style }) { return <div style={{ fontSize:14, color:'hsl(215.4 16.3% 46.9%)', ...style }}>{children}</div>; }
function CardContent({ children, style }) { return <div style={{ padding:'16px 24px 24px', ...style }}>{children}</div>; }

const inputBase = { display:'flex', height:40, width:'100%', borderRadius:6, border:'1px solid hsl(214.3 31.8% 91.4%)', background:'white', padding:'0 12px', fontSize:14, fontFamily:'inherit', color:'inherit', outline:'none' };
function Input(props) { return <input {...props} style={{ ...inputBase, ...(props.style||{}) }} />; }
function Select({ children, ...props }) { return <select {...props} style={{ ...inputBase, ...(props.style||{}) }}>{children}</select>; }
function Label({ children, htmlFor, style }) { return <label htmlFor={htmlFor} style={{ fontSize:14, fontWeight:500, lineHeight:1, display:'block', ...style }}>{children}</label>; }

function Progress({ value=0, max=100, style }) {
  const p = Math.min(100, Math.max(0, (value/max)*100));
  return (
    <div style={{ position:'relative', height:16, width:'100%', overflow:'hidden', borderRadius:9999, background:'hsl(210 40% 96.1%)', ...style }}>
      <div style={{ height:'100%', background:'hsl(222.2 47.4% 11.2%)', width:`${p}%`, transition:'width 300ms cubic-bezier(.4,0,.2,1)' }}/>
    </div>
  );
}

function Field({ label, error, children }) {
  return (
    <div style={{ display:'flex', flexDirection:'column', gap:8 }}>
      <Label>{label}</Label>
      {children}
      {error && <div style={{ fontSize:14, color:'hsl(0 84.2% 60.2%)' }}>{error}</div>}
    </div>
  );
}

Object.assign(window, { Button, Badge, Card, CardHeader, CardTitle, CardDescription, CardContent, Input, Select, Label, Progress, Field, cls });
