/**
 * clamp Contrain une valeur entre un minimum et un maximum.
 * 
 * @param value {number} Valeur à contraindre.
 * @param min {number} Valeur minimale.
 * @param max {number} Valeur maximale.
 * @returns {number}
 */
const clamp = (value: number, min: number, max: number): number => {
  return Math.min(Math.max(value, min), max);
}

export { clamp };