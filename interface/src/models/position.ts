
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

export { Vector };