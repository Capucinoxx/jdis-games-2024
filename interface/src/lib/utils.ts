type CartesianCoordinates = { x: number, y: number };
type IsometricCoordinates = { x: number, y: number };

const cartesian_to_isometric = (cart: CartesianCoordinates): IsometricCoordinates => {
  return { x: cart.x - cart.y, y: (cart.x + cart.y) / 2 };
};

const isometric_to_cartesian = (iso: IsometricCoordinates): CartesianCoordinates => {
  return { x: (iso.x + 2 * iso.y) / 2, y: iso.y - iso.x / 2 };
};


export {
  CartesianCoordinates,
  IsometricCoordinates,
  cartesian_to_isometric,
  isometric_to_cartesian
};