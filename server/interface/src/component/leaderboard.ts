import { animate_number, toggle_btn, toggle_fullscreen, toggle_toast } from "../animation";
import { LineChart, UpdateOptions } from "../chart";

interface GlobalData {
  top: UpdateOptions;
  leaderboard: RowData[];
};

interface RowData {
  name: string;
  score: number;
  ranking: number;
};

class Leaderboard {
  private root: HTMLUListElement;
  private items: Map<string, HTMLLIElement>;

  constructor(root: HTMLUListElement) {
    this.root = root;
    this.items = new Map();
  }

  public update(data: RowData[]): void {
    const fragment = document.createDocumentFragment();

    data.forEach(row => {
      let el = this.items.get(row.name);
      if (!el) {
        el = this.create_item(row);
        fragment.appendChild(el);
        this.items.set(row.name, el);
      } else {
        animate_number(el, row.score);
        this.update_ranking(el, row.ranking);
      }
    });

    this.root.appendChild(fragment);
  }

  private create_item(row: RowData): HTMLLIElement {
    const li = document.createElement('li');
    li.dataset.uid = row.name;
    li.dataset.score = row.score.toString();
    li.dataset.rank = row.ranking.toString();

    li.classList.add('item');
    this.transform_position(li, row.ranking);

    const ranking = document.createElement('div');
    ranking.classList.add('item-ranking');
    ranking.textContent = `${row.ranking}`;

    const name = document.createElement('div');
    name.classList.add('item-name');
    name.textContent = row.name;

    const score = document.createElement('div');
    score.classList.add('item-score');
    score.textContent = `${row.score}`;

    li.appendChild(ranking);
    li.appendChild(name);
    li.appendChild(score);

    return li;
  }

  private transform_position(li: HTMLLIElement, rank: number): void {
    li.style.transform = `translate(0, ${(rank - 1) * 100}%)`;
  }

  private update_ranking(li: HTMLLIElement, rank: number): void {
    if (parseInt(li.dataset.rank!, 10) !== rank) {
      li.dataset.rank = rank.toString();

      const el = li.querySelector('.item-ranking')!;
      el.textContent = `${rank}`;

      this.transform_position(li, rank);
    }
  }
};

class Leaderboards {
  private current_leaderboard: Leaderboard | null = null;
  private global_leaderboard: Leaderboard | null = null;
  private chart: LineChart;

  constructor(root: HTMLElement, current_root: HTMLUListElement | null, global_root: HTMLUListElement | null) {
    if (current_root) this.current_leaderboard = new Leaderboard(current_root);
    if (global_root)  this.global_leaderboard = new Leaderboard(global_root);

    this.chart = new LineChart('leaderboard-graph');

    this.handle_expansion(root);
    this.handle_open();
  }

  private handle_open(): void {
    toggle_btn(document.querySelector('.leaderboard-btn')!, () => {
      toggle_toast(document.getElementById('leaderboard')!);
    });
  }

  private handle_expansion(root: HTMLElement) {
    const btn = root.querySelector('.expand');
    if (btn)
      toggle_fullscreen(root, btn as HTMLElement);
  }

  public set current(data: RowData[]) {
    if(this.current_leaderboard) this.current_leaderboard.update(data);
  }

  public set global(data: GlobalData) {
    if (this.global_leaderboard) this.global_leaderboard.update(data.leaderboard);
    this.chart.update(data.top);
  }
};

export { Leaderboards };