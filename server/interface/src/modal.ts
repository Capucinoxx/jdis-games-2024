import { gsap } from 'gsap';

const handle_modal_animation = (game: Phaser.Game, container: HTMLElement, btn: HTMLElement, toggle_validation: ((e: KeyboardEvent, active: boolean) => boolean) | null = null) => {
  let active = false;

  const btn_bg = btn.querySelector('.btn--bg');
  const btn_open = btn.querySelector('.open');
  const btn_close = btn.querySelector('.close');

  const open = gsap.timeline({ paused: true });
  const close = gsap.timeline({ paused: true });

  open
    .set(btn,       { pointerEvents: 'none', zIndex: 11 })
    .to(container,  { 'clip-path': 'circle(200% at 100% 20px)', duration: 1.5, ease: 'power4.out' }, 0)
    .to(btn_close,  { opacity: 1, yPercent: -125, duration: 1, ease: 'power4.out' }, 0)
    .to(btn_open,   { opacity: 0, yPercent: -125, duration: 1, ease: 'power4.out' }, 0)
    .to(btn,        { 'box-shadow': '0 4px 8px #c2ccff33', duration: 1.5, ease: 'power4.out' }, 0)
    .set(btn,       { pointerEvents: 'all' });

  close 
    .set(btn,       { pointerEvents: 'none', zIndex: 9 })
    .to(btn_bg,     { scale: 0.9, duration: 0.25, ease: 'elastic.out' }, 0)
    .to(btn_bg,     { scale: 1, duration: 0.25, ease: 'elastic.out' }, '+=0.25')
    .to(container,  { 'clip-path': 'circle(0% at 100% 20px)', duration: 1.2, ease: 'power4.out' }, '-=0.5')
    .to(btn_close,  { opacity: 0, yPercent: 125, duration: 1, ease: 'power4.out' }, 0)
    .to(btn_open,   { opacity: 1, yPercent: 0, duration: 1, ease: 'power4.out' }, 0)
    .to(btn,        { 'box-shadow': '0 4px 8px #c2ccff00', duration: 1.5, ease: 'power4.out' }, 0)
    .set(btn,       { pointerEvents: 'all' });

    const toggle_menu = () => {
      if (active) close.seek(0).play();
      else open.seek(0).play();
      game && game.input.keyboard && (game.input.keyboard.enabled = active);
      active = !active;
    };

    btn.addEventListener('click', toggle_menu);

    if (toggle_validation)
      document.addEventListener('keydown', (e) => {
        if (toggle_validation(e, active))
          toggle_menu();
      });
};

const handle_modals = (game: Phaser.Game) => {
  const navbar = document.querySelector('nav');
  if (!navbar) return;

  const navbar_btn = document.querySelector('.nav-btn') as HTMLElement | undefined;
  if (!navbar_btn) return;

  const leaderboard = document.querySelector('#leaderboard') as HTMLElement | undefined;
  if (!leaderboard) return;

  const leaderboard_btn = document.querySelector('.leaderboard-btn') as HTMLElement | undefined;
  if (!leaderboard_btn) return;

  handle_modal_animation(game, navbar, navbar_btn, (e, b) => (e.key =='m' && !b) || (e.key === 'Escape' && b));
  handle_modal_animation(game, leaderboard, leaderboard_btn);
};



export { handle_modals };