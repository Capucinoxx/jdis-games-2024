import { Chart, LineController, LineElement, PointElement, LinearScale, Title, CategoryScale, ChartConfiguration } from 'chart.js';

Chart.register(LineController, LineElement, PointElement, LinearScale, Title, CategoryScale);

interface UpdateOptions {
    data: number[][];
    colors: string[];
}

class LineChart {
    private ctx: HTMLCanvasElement;
    private chart: Chart;

    constructor(canvasId: string) {
        this.ctx = document.getElementById(canvasId) as HTMLCanvasElement;
        this.chart = new Chart(this.ctx, this.getConfig());
    }

    private getConfig(): ChartConfiguration {
        const datasets = Array.from({ length: 10 }, (_, i) => ({
            label: `Dataset ${i + 1}`,
            data: [],
            borderWidth: 1,
            fill: false,
            pointRadius: 0
        }));

        const config = {
            type: 'line',
            data: { datasets: datasets },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: { display: false },
                    y: { display: false }
                },
                plugins: { legend: { display: true } }
            }
        } as ChartConfiguration;

        return config;
    }

    public update(options: UpdateOptions) {
        this.chart.data.datasets.forEach((dataset, index) => {
            dataset.data = options.data[index] || dataset.data;
        });


        options.colors.forEach((color, index) => {
            if (this.chart.data.datasets[index]) {
                this.chart.data.datasets[index].borderColor = color;
            }
        });


        this.chart.update();
    }
}

export { LineChart, UpdateOptions };
