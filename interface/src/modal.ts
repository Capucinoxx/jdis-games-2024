import { gsap } from 'gsap';

const handle_modal_interraction = () => {
  let active = false;

  const navbar = document.querySelector('nav');
  if (!navbar)
    return;

  const nav_els = navbar!.querySelectorAll('li');
  const btn = document.querySelector('.nav-btn');
  const btn_bg = btn!.querySelector('.nav-btn--bg');
  const btn_open = btn!.querySelector('.line');
  const btn_close = btn!.querySelector('.close');

 
  const title_sections = Array.from(nav_els)
    .map((el) => el.querySelector("input[type='checkbox']"))
    .filter((el) => el !== null && el !== undefined) as HTMLInputElement[];



  title_sections.forEach(s => 
    s.addEventListener('click', (e) => {
      const cur = e.target as HTMLInputElement;
      if (cur.checked) {
        title_sections.forEach((ss) => (ss !== cur && (ss.checked = false)));
      }
    })
  );

  const open = gsap.timeline({ paused: true });
  const close = gsap.timeline({ paused: true });

  open
    .set(btn,       { pointerEvents: 'none' })
    .to(navbar,     { 'clip-path': 'circle(200% at 60px 60px)', duration: 1.5, ease: 'power4.out' }, 0)
    .to(nav_els,    { x: 0, opacity: 1, pointerEvents: 'all', duration: 1.25, stagger: 0.1, ease: 'elastic.out(1.15, .95)' }, 0)
    .to(btn_close,  { opacity: 1, yPercent: -125, duration: 1, ease: 'power4.out' }, 0)
    .to(btn_open,   { opacity: 0, yPercent: -125, duration: 1, ease: 'power4.out' }, 0)
    .set(btn,       { pointerEvents: 'all' });

  close 
    .set(btn,       { pointerEvents: 'none' })
    .to(btn_bg,     { scale: 0.9, duration: 0.25, ease: 'elastic.out' }, 0)
    .to(btn_bg,     { scale: 1, duration: 0.25, ease: 'elastic.out' }, '+=0.25')
    .to(navbar,     { 'clip-path': 'circle(0% at 60px 60px)', duration: 1.2, ease: 'power4.out' }, '-=0.5')
    .to(nav_els,    { x: -200, opacity: 0, pointerEvents: 'none', duration: 1, stagger: 0.1, ease: 'power4.out' }, 0)
    .to(btn_close,  { opacity: 0, yPercent: 125, duration: 1, ease: 'power4.out' }, 0)
    .to(btn_open,   { opacity: 1, yPercent: 0, duration: 1, ease: 'power4.out' }, 0)
    .set(btn,       { pointerEvents: 'all' });

  btn?.addEventListener('click', () => { 
    if (active) close.seek(0).play();
    else open.seek(0).play(); 
    active = !active; 
  });
};

export { handle_modal_interraction };
