import { animate_number } from "../animation";

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
    ranking.textContent = `# ${row.ranking}`;

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
      el.textContent = `# ${rank}`;

      this.transform_position(li, rank);
    }
  }
};

export { Leaderboard };