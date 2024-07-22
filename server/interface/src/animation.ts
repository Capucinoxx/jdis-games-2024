import { gsap } from 'gsap';

const toggle_toast = (el: HTMLElement) => {
  console.log(el);
  if (el.style.visibility === 'visible') {
    gsap.to(el, {
      duration: 0.5, opacity: 0, y: 20,
      onComplete: () => { el.style.visibility = 'hidden' },
    })
  } else {
    el.style.visibility = 'visible';
    gsap.fromTo(el, {opacity: 0, y: 20}, {duration: 0.5, opacity: 1, y: 0});
  }
};

const toggle_btn = (btn: HTMLElement, toggle: (() => void) | null = null): void => {
  let active = false;

  const open = gsap.timeline({ paused: true });
  const close = gsap.timeline({ paused: true });

  const btn_bg = btn.querySelector('.btn--bg');
  const btn_open = btn.querySelector('.open');
  const btn_close = btn.querySelector('.close');

  open
    .set(btn,       { pointerEvents: 'none' })
    .to(btn_close,  { opacity: 1, yPercent: -125, duration: 1, ease: 'power4.out' }, 0)
    .to(btn_open,   { opacity: 0, yPercent: -125, duration: 1, ease: 'power4.out' }, 0)
    .to(btn,        { 'box-shadow': '0 4px 8px #c2ccff33', duration: 1.5, ease: 'power4.out' }, 0)
    .set(btn,       { pointerEvents: 'all' });

  close 
    .set(btn,       { pointerEvents: 'none' })
    .to(btn_bg,     { scale: 0.9, duration: 0.25, ease: 'elastic.out' }, 0)
    .to(btn_bg,     { scale: 1, duration: 0.25, ease: 'elastic.out' }, '+=0.25')
    .to(btn_close,  { opacity: 0, yPercent: 125, duration: 1, ease: 'power4.out' }, 0)
    .to(btn_open,   { opacity: 1, yPercent: 0, duration: 1, ease: 'power4.out' }, 0)
    .to(btn,        { 'box-shadow': '0 4px 8px #c2ccff00', duration: 1.5, ease: 'power4.out' }, 0)
    .set(btn,       { pointerEvents: 'all' });

  const on_toggle = () => {
    if (active) close.seek(0).play();
    else open.seek(0).play();
    active = !active;

    if (toggle)
        toggle();
  };

  btn.addEventListener('click', on_toggle);
};

const animate_number = (li: HTMLLIElement, finale_score: number): void => {
  const el = li.querySelector('.item-score')!;
  let current_score = parseInt(li.dataset.score!, 10);
  li.dataset.score = finale_score.toString();

  const duration = 30;
  let count = 0;
  const step = (finale_score - current_score) / duration;

  const run = () => {
    if (count === duration) {
      el.textContent = finale_score.toLocaleString('fr-FR');
      return;
    }

    el.textContent = Math.floor((current_score += step)).toLocaleString('fr-FR');
    count++;
    requestAnimationFrame(run);
  };

  requestAnimationFrame(run);
};

export { toggle_toast, toggle_btn, animate_number };