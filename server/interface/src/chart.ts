import { Chart, LineController, LineElement, PointElement, LinearScale, Title, CategoryScale, ChartConfiguration } from 'chart.js';

Chart.register(LineController, LineElement, PointElement, LinearScale, Title, CategoryScale);

interface UpdateOptions {
    data?: number[];
    colors?: string[];
}

class LineChart {
    private ctx: HTMLCanvasElement;
    private chart: Chart;

    constructor(canvasId: string) {
        this.ctx = document.getElementById(canvasId) as HTMLCanvasElement;
        this.chart = new Chart(this.ctx, this.getConfig());
    }

    private ddd(): number[] {
        let v = Math.floor(Math.random() * 1000);

        const data = [];
        for (let i = 0; i < 100; i++) {
            v += Math.floor(Math.random() * 1000) * i;
            data.push(v);
        }
        return data;
    }

    private getConfig(): ChartConfiguration {
        const datasets = Array.from({ length: 10 }, (_, i) => ({
            label: `Dataset ${i + 1}`,
            data: this.ddd(),
            borderColor: `rgba(${Math.floor(Math.random() * 255)}, ${Math.floor(Math.random() * 255)}, ${Math.floor(Math.random() * 255)}, 1)`,
            borderWidth: 1,
            fill: false,
            pointRadius: 0
        }));

        const data = {
            labels: new Array(10).fill(''),
            datasets: datasets
        };

        const config = {
            type: 'line',
            data: data,
            options: {
                scales: {
                    x: {
                        display: false
                    },
                    y: {
                        display: false
                    }
                },
                plugins: {
                    legend: {
                        display: false
                    }
                }
            }
        } as ChartConfiguration;

        return config;
    }

    public update(options: UpdateOptions) {
        if (options.data) {
            this.chart.data.datasets[0].data = options.data;
        }

        if (options.colors) {
            options.colors.forEach((color, index) => {
                if (this.chart.data.datasets[index]) {
                    this.chart.data.datasets[index].borderColor = color;
                }
            });
        }

        this.chart.update();
    }
}

export { LineChart };
