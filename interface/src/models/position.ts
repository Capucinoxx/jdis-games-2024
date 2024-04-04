import { Graphics } from 'pixi.js';

class Polygon {
  public graphics: Graphics;
  private vertices: Vector[];

  constructor(...vertices: Vector[]) {
    this.graphics = new Graphics();
    this.vertices = vertices;

    this.draw();
  }

  /**
   * draw Dessine le polygone sur le canvas.
   */
  private draw(): void {
    this.graphics.clear();
    this.graphics.moveTo(this.vertices[0].x, this.vertices[0].y);
    for (let i = 1; i < this.vertices.length; i++) {
      this.graphics.lineTo(this.vertices[i].x, this.vertices[i].y);
    }
    this.graphics.lineTo(this.vertices[0].x, this.vertices[0].y);
    this.graphics.fill(0x000000);
    this.graphics.stroke({ color: 0xffffff, width: 1  });
  }

  /**
   * multiply Multiplie les vecteurs du polygone par un scalaire.
   * 
   * @param scalar {number} Scalaire par lequel multiplier le vecteur
   * @returns Polygone résultant de la multiplication par scalaire
   */
  public multiply(scalar: number): Polygon {
    let vertices = this.vertices.map(v => v.multiply(scalar));
    return new Polygon(...vertices);
  }
};


class Vector {
  constructor(public x: number, public y: number) {}

  /**
   * add Additionne un vecteur à un autre vecteur ou un nombre à un vecteur.
   * Lorsque l'on ajoute un nombre à un vecteur, on ajoute ce nombre à chaque 
   * composante du vecteur.
   * 
   * @param v {Vector | number} Vecteur ou nombre à ajouter
   * @param y {number} Composante y du vecteur à ajouter si le premier argument est un nombre.
   * @returns Vecteur résultant de l'addition
   */
  public add(v: Vector | number, y?: number): Vector {
    if (v instanceof Vector) {
        return new Vector(this.x + v.x, this.y + v.y);
    } else if (typeof v === 'number' && typeof y === 'number') {
        return new Vector(this.x + v, this.y + y);
    }
    throw new Error('Invalid arguments for Vector.add');
  }

  /**
   * subtract Soustrait un vecteur à un autre vecteur ou un nombre à un vecteur.
   * Lorsque l'on soustrait un nombre à un vecteur, on soustrait ce nombre à chaque
   * composante du vecteur.
   * 
   * @param v {Vector | number} Vecteur ou nombre à soustraire 
   * @param y {number} Composante y du vecteur à soustraire si le premier argument est un nombre.
   * @returns Vecteur résultant de la soustraction
   */
  public subtract(v: Vector | number, y?: number): Vector {
    if (v instanceof Vector) {
        return new Vector(this.x - v.x, this.y - v.y);
    } else if (typeof v === 'number' && typeof y === 'number') {
        return new Vector(this.x - v, this.y - y);
    }
    throw new Error('Invalid arguments for Vector.subtract');
  }

  /**
   * multiply Multiplie un vecteur par un scalaire.
   * 
   * @param scalar {number} Scalaire par lequel multiplier le vecteur
   * @returns Vecteur résultant de la multiplication par scalaire
   */
  public multiply(scalar: number): Vector {
    return new Vector(this.x * scalar, this.y * scalar);
  }
};

export { Polygon, Vector };