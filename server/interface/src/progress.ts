class ProgressBar {
    public static MAX_VALUE: number = 900;

    private container: SVGRectElement | null = null;
    private value: number = 0;

    constructor(container: SVGRectElement) {
        this.container = container;
    }

    private update() {
        const percentage = Math.min(this.value / ProgressBar.MAX_VALUE * 100, 100);
        this.animate(percentage);
    }

    private animate(percentage: number) {
        const total_length = 4 * 800;
        const progress_length = (percentage / 100) * total_length;
        const dash_array = [progress_length, total_length - progress_length];
        this.container!.style.strokeDasharray = dash_array.join(' ');
        this.container!.style.transition = 'stroke-dasharray 0.03s linear';
    } 

    public reset() {
        if (!this.container)
            return;

        this.container.style.transition = 'none';
        this.container.style.strokeDasharray = '0 3200';
        setTimeout(() => {
            this.container!.style.transition = 'stroke-dasharray 0.03s linear';
        }, 5);
    }

    public set current_value(v: number) {
        if (v == 0)
            this.reset();

        this.value = v;
        this.update();
    }
};

export { ProgressBar };

