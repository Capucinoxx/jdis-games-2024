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
            borderColor: 'rgba(0,0,0,1)',
            borderWidth: 2,
            fill: false,
            pointRadius: 0
        }));

        const config = {
            type: 'line',
            data: { datasets: datasets, labels: [] },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    x: { display: false, type: 'linear' },
                    y: { display: false }
                },
                plugins: { legend: { display: true } }
            }
        } as ChartConfiguration;

        return config;
    }

    public update(options: UpdateOptions) {
        console.log(options);

        const length = Math.max(...options.data.map(arr => arr.length));
        this.chart.data.labels = Array.from({ length: length }, (_, i) => i);

        this.chart.data.datasets.forEach((dataset, index) => {
            dataset.data = options.data[index] || dataset.data;
        });


        options.colors.forEach((color, index) => {
            if (this.chart.data.datasets[index])
                this.chart.data.datasets[index].borderColor = color;
        });
        this.chart.update();
    }
}

export { LineChart, UpdateOptions };
