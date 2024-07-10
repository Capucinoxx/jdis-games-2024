
const slide_down = (el: HTMLElement): Promise<void> => {
  return new Promise<void>((resolve) => {
    el.style.transition = 'height .5s, margin .5s, padding .5s';
    el.style.height = '0';
    el.style.margin = '0';
    el.style.padding = '0';
    window.setTimeout(() => { resolve() }, 500);
  });
};

interface ToastOptions {
  type?: string;
  subject?: string;
  message?: string;
};

class Toast extends HTMLElement {
  private type: string;
  private message?: string;
  private subject: string;

  constructor({ type, subject, message }: ToastOptions = {}) {
    super();
    this.type = type ?? 'alert';
    this.subject = subject ?? 'information';
    this.message = message;
    this.close = this.close.bind(this);
  }

  connectedCallback() {
    this.subject = this.getAttribute('subject') ?? this.subject;
    this.type = this.getAttribute('type') ?? this.type;

    const msg = this.innerHTML;

    this.classList.add('alert', `${this.type}`);
    this.innerHTML = `
      <h3>${this.subject}</h3>
      <button class='alert-close'>x</button>
      <div>${this.message || msg}</div>
    `;

    this.querySelector('.alert-close')?.addEventListener('click', (e) => {
      e.preventDefault();
      this.close();
    });
  }

  private close() {
    this.classList.add('out');
    window.setTimeout(async () => {
      await slide_down(this);
      this.parentElement?.removeChild(this);
    });
  }
};

customElements.define('toast-alert', Toast);

const handle_forms = () => {
  const forms = document.querySelectorAll<HTMLFormElement>('form');
  forms.forEach(form => {
    form.querySelector('button')?.addEventListener('click', async (e) => {
      e.preventDefault();
      e.stopPropagation();

      const data = Array.from(new FormData(form))
        .reduce((acc, [key, val]) => ({ ...acc, [key]: val }), {});

      const response = await fetch(form.action, {
        method: form.method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
      });

      if (!response.ok) {
        flash({ subject: 'error', message: 'Server error' });
        return;
      }

      const result = await response.json();
      flash(result);
    })
  }
  );
};

const flash = (data: { [key: string]: any }) => {
  const toast = document.createElement('toast-alert');
  data.type && toast.setAttribute('type', data.type);
  data.subject && toast.setAttribute('subject', data.subject);
  toast.innerText = data.message;

  document.body.appendChild(toast);
  
};

export { handle_forms };
