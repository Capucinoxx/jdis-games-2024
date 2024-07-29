const generate_particles = (canvas_container: string): void => {
  const canvas = document.getElementById(canvas_container) as HTMLCanvasElement;
  const ctx = canvas.getContext('2d');

  if (!ctx)
    throw new Error("Cannot get canvas context.");

  canvas.width = window.innerWidth;
  canvas.height = window.innerHeight;



  class Particle {
    private x: number = 0;
    private y: number;
    private speed: number = 0;
    private opacity: number = 0;
    private fade_delay: number;
    private fade_start: number;
    private fading_out: boolean;

    constructor() {
      this.reset();
      this.y = Math.random() * canvas.height;
      this.fade_delay = Math.random() * 600 + 100;
      this.fade_start = Date.now() + this.fade_delay;
      this.fading_out = false;
    }

    private reset(): void {
      this.x = Math.random() * canvas.width;
      this.y = Math.random() * canvas.height;
      this.speed = Math.random() / 5 + 0.1;
      this.opacity = 1;
      this.fade_delay = Math.random() * 600 + 100;
      this.fade_start = Date.now() + this.fade_delay;
      this.fading_out = false;
    }

    public update(): void {
      this.y -= this.speed;
      if (this.y < 0) {
        this.reset();
      }

      if (!this.fading_out && Date.now() > this.fade_start) {
        this.fading_out = true;
      }

      if (this.fading_out) {
        this.opacity -= 0.008;
        if (this.opacity <= 0) {
          this.reset();
        }
      }
    }

    public draw(): void {
      ctx!.fillStyle = `rgba(${255 - (Math.random() * 127.5)}, 255, 255, ${this.opacity})`;
      ctx!.fillRect(this.x, this.y, 1, Math.random() * 2 + 1);
    }
  }

  const calculate_particle_count = (): number => {
    return Math.floor((canvas.width * canvas.height) / 20_000);
  }

  let particles: Particle[] = [];
  let particle_count = calculate_particle_count();

  const init_particles = (): void => {
    particles = [];
    for (let i = 0; i < particle_count; i++) {
      particles.push(new Particle());
    }
  }

  const animate = (): void => {
    ctx!.clearRect(0, 0, canvas.width, canvas.height);
    particles.forEach(particle => {
      particle.update();
      particle.draw();
    });
    requestAnimationFrame(animate);
  }



  const on_resize = (): void => {
    canvas.width = window.innerWidth;
    canvas.height = window.innerHeight;
    particle_count = calculate_particle_count();
    init_particles();
  }

  window.addEventListener('resize', on_resize);

  init_particles();
  animate();
}

export { generate_particles };