import { Chart, LineController, LineElement, PointElement, LinearScale, Title, CategoryScale, TimeScale, ChartConfiguration } from 'chart.js';
import 'chartjs-adapter-date-fns';

Chart.register(LineController, LineElement, PointElement, LinearScale, Title, CategoryScale, TimeScale);


interface UpdateOptions {
    data: { x: number, y: number }[][];
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
                    x: {
                        type: 'time',
                        min: 0,
                        max: 1000,
                        display: false, 
                        time: {
                            unit: 'millisecond',
                            tooltipFormat: 'Pp',
                            displayFormats: {
                                second: 'HH:mm:ss'
                            }
                        }
                    },
                    y: { display: false }
                },
                plugins: { legend: { display: true } }
            }
        } as ChartConfiguration;

        return config;
    }

    public update(options: UpdateOptions) {
        const length = Math.max(...options.data.map(arr => arr.length));
        const minimum = Math.min(...options.data.map(arr => Math.min(...arr.map(v => v.x))));
        const maximum = Math.max(...options.data.map(arr => Math.max(...arr.map(v => v.x))));
        
        this.chart.options.scales!.x!.min = minimum;
        this.chart.options.scales!.x!.max = maximum;

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
